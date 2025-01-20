package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"zzz_helper/internal/config"
	"zzz_helper/internal/utils"
	"zzz_helper/internal/utils/lznet/lzhttp"
)

type MiyousheGameRole struct {
	Data struct {
		List []struct {
			GameBiz    string `json:"game_biz"`
			GameUid    string `json:"game_uid"`
			IsChosen   bool   `json:"is_chosen"`
			IsOfficial bool   `json:"is_official"`
			Level      int    `json:"level"`
			Nickname   string `json:"nickname"`
			Region     string `json:"region"`
			RegionName string `json:"region_name"`
		} `json:"list"`
	} `json:"data"`
	Message string `json:"message"`
	Retcode int    `json:"retcode"`
}

type GameInfo struct {
	GameBiz string `json:"game_biz"`
	Region  string `json:"region"`
	GameUid string `json:"game_uid"`
}

func NewMihoyoClient(miyousheCookie, mihoyoCookie string) (*MihoyoClient, error) {
	mihoyoClient := &MihoyoClient{
		miyousheCookie: miyousheCookie,
		mihoyoCookie:   mihoyoCookie,
	}
	err := mihoyoClient.Init()
	if err != nil {
		return nil, err
	}
	return mihoyoClient, nil
}

type MihoyoClient struct {
	client         *lzhttp.HttpClient
	miyousheCookie string
	mihoyoCookie   string
}

func (this *MihoyoClient) Init() error {
	this.client = lzhttp.NewHttpClient(lzhttp.HttpClientInitParam{
		Charset:      "UTF-8",
		ConnTimeout:  time.Second * 5,
		TransTimeout: time.Second * 10,
	})
	return nil
}

func dsGen() string {
	salt := "AbuxbruiFDIgxLXksUNMAMvDyciznofM"
	t := fmt.Sprintf("%v", time.Now().Unix())
	r := utils.RandStrAndDigital(6)
	meta := fmt.Sprintf("salt=%s&t=%s&r=%s", salt, t, r)
	hash := md5.Sum([]byte(meta))
	hashStr := hex.EncodeToString(hash[:])

	ds := strings.Join([]string{t, r, hashStr}, ",")
	return ds
}

func (this *MihoyoClient) GetInfo() ([]GameInfo, error) {
	version := "2.79.0"
	clientType := "5"

	headers := lzhttp.Header{}
	headers.Set("Cookie", this.miyousheCookie)
	headers.Set("x-rpc-app_version", version)
	headers.Set("x-rpc-client_type", clientType)
	headers.Set("DS", dsGen())
	resp := this.client.HttpGet("https://api-takumi.miyoushe.com/binding/api/getUserGameRolesByStoken", headers)
	if resp.Err != nil {
		return nil, resp.Err
	}
	if resp.Resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status code %d", resp.Resp.StatusCode)
	}
	role := MiyousheGameRole{}
	err := json.Unmarshal([]byte(resp.Body), &role)
	if err != nil {
		return nil, err
	}
	if role.Retcode != 0 {
		return nil, fmt.Errorf("%s", role.Message)
	}
	infos := make([]GameInfo, 0)
	for _, s := range role.Data.List {
		infos = append(infos, GameInfo{
			GameBiz: s.GameBiz,
			Region:  s.Region,
			GameUid: s.GameUid,
		})
	}
	return infos, nil
}
func (this *MihoyoClient) Sign() error {
	infos, err := this.GetInfo()
	if err != nil {
		return err
	}
	for _, info := range infos {
		game := ""
		actId := ""
		switch info.GameBiz {
		case "hk4e_cn":
			game = "hk4e"
			actId = "e202311201442471"
		case "nap_cn":
			game = "zzz"
			actId = "e202406242138391"
		}

		headers := lzhttp.Header{}
		headers.Set("Cookie", this.mihoyoCookie)
		headers.Set("Content-Type", "application/json")
		headers.Set("x-rpc-signgame", game)

		data := fmt.Sprintf(`{"act_id":"%s","region":"%s","uid":"%s","lang":"zh-cn"}`, actId, info.Region, info.GameUid)
		resp := this.client.HttpPost(fmt.Sprintf("https://api-takumi.mihoyo.com/event/luna/%s/sign", game), data, headers, false)
		if resp.Err != nil {
			log.Printf("[-] %s\n", resp.Err.Error())
			continue
		}
		if resp.Resp.StatusCode != 200 {
			log.Printf("[-] http status code %d\n", resp.Resp.StatusCode)
			continue

		}
		log.Printf("[+] %s\n", resp.Body)
	}
	return nil
}
func (this *MihoyoClient) Close() error {
	if this.client != nil {
		return this.client.Close()
	}
	return nil
}

func subSign() error {
	client, err := NewMihoyoClient(config.GlobalConfig.MiyousheCookie, config.GlobalConfig.MihoyoCookie)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Sign()
	if err != nil {
		return err
	}
	return nil
}
func main() {
	for {
		err := subSign()
		if err != nil {
			log.Printf("[-] %s\n", err.Error())
		}
		time.Sleep(time.Hour * 6)
	}
}
