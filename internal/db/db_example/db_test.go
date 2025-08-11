package db_example

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
	"zzz_helper/internal/db/accelerator/db_cache"
	"zzz_helper/internal/db/db_model"
	"zzz_helper/internal/global_config"
	"zzz_helper/internal/utils"
	"zzz_helper/internal/utils/byte2"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/reflect2"
	"zzz_helper/internal/utils/serial"
	"zzz_helper/pkg/crypto"
)

func TestDB(t *testing.T) {
	var a = reflect2.MustToStructPtr(IpPoolDB{IP: "666"})
	t.Log(a)

	db := refGetIpPoolDB("")
	result, err := db.ReadOne(&IpPoolDB{
		IP: "1.1.1.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.(*IpPoolDB).Protocol != "TCP" {
		t.Fatal("protocol not equall")
	}

	err = db.Delete()
	if err != nil {
		t.Fatal(err)
	}
	db.Insert(&IpPoolDB{
		IP:            "1111",
		Protocol:      "",
		Describe:      "",
		OnlineCount:   0,
		OnlineClient:  "",
		CompileCount:  0,
		CompileClient: "",
	})
	results, err := db.Read(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].(*IpPoolDB).IP != "1111" {
		t.Fatal("not equal")
	}
	err = db.Update(&IpPoolDB{IP: "12345"}, false, &IpPoolDB{IP: "1111"})
	if err != nil {
		t.Fatal(err)
	}
	result, err = db.ReadOne(&IpPoolDB{IP: "12345"})
	if err != nil {
		t.Fatal(err)
	}
	if result.(*IpPoolDB).IP != "12345" {
		t.Fatal("not equal 12345")
	}

	results, err = db.ReadWithFilter(-1, 01, []map[string]int{
		{
			"ip": db_model.FilterLike,
		},
	}, &IpPoolDB{IP: "121"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)

}

func TestMft(t *testing.T) {

	data, _ := file2.ReadFileBytes(`E:\code\go\Level6WorkSpace\Level6\tmpout\server_acc\cache\file_download_ccc8768ceca0fbb8954c10467f3b2cf2`)
	data = utils.XorEncode(data, global_config.CacheKey)
	data = utils.XorEncode(data, []byte("ce69440c-402a-4667-af8c-eadf588f32a3"))
	data = crypto.GzipDecompress(data)

	fileParams := serial.NewControlParams()
	err := fileParams.FormatUnMarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	driverStr := fileParams.MustGetStringParam("drivers")
	drivers := strings.Split(driverStr, ":")
	mftFiles := make([]*db_cache.AgentMftCacheDB, 0)

	for _, driver := range drivers {
		volBytes := fileParams.MustGetBytesParam(driver)
		readBuffer := bytes.NewBuffer(volBytes)
		volMftFiles := map[uint64]*db_cache.AgentMftCacheDB{}
		for {
			if readBuffer.Len() == 0 {
				break
			}
			indexBytes := make([]byte, 8)
			_, err = readBuffer.Read(indexBytes)
			if err != nil {
				break
			}
			index := byte2.BytesMustToUint(indexBytes, false)

			parentIndexBytes := make([]byte, 8)
			_, err = readBuffer.Read(parentIndexBytes)
			if err != nil {
				break
			}
			parentIndex := byte2.BytesMustToUint(parentIndexBytes, false)

			fileNameLengthBytes := make([]byte, 2)
			_, err = readBuffer.Read(fileNameLengthBytes)
			if err != nil {
				break
			}
			fileNameLength := byte2.BytesMustToUint(fileNameLengthBytes, false)

			fileNameBytes := make([]byte, fileNameLength)
			_, err = readBuffer.Read(fileNameBytes)
			if err != nil {
				break
			}
			fileName := string(fileNameBytes)
			volMftFiles[index] = &db_cache.AgentMftCacheDB{
				Index:       fmt.Sprintf("%d", index),
				ParentIndex: fmt.Sprintf("%d", parentIndex),
				FileName:    fileName,
			}
		}
		for _, file := range volMftFiles {
			//getPath(volMftFiles, index)
			mftFiles = append(mftFiles, file)
		}
	}

	t.Logf("mft length: %d, start write database", len(mftFiles))

	db := db_cache.GetAgentMftCacheDB("test")
	db.Delete()
	err = db.InsertMulti(mftFiles)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("write database over")
}

func TestMft2(t *testing.T) {

	data, _ := file2.ReadFileBytes(`E:\code\go\Level6WorkSpace\Level6\tmpout\server_acc\cache\file_download_3a6aa5a3f79182cf7472120a173da02c`)
	data = utils.XorEncode(data, global_config.CacheKey)
	data = utils.XorEncode(data, []byte("ce69440c-402a-4667-af8c-eadf588f32a3"))
	data = crypto.GzipDecompress(data)

	fileParams := serial.NewControlParams()
	err := fileParams.FormatUnMarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	driverStr := fileParams.MustGetStringParam("drivers")
	drivers := strings.Split(driverStr, ":")
	mftFiles := make([]*db_cache.AgentMftCacheDB, 0)

	for _, driver := range drivers {
		volBytes := fileParams.MustGetBytesParam(driver)
		readBuffer := bytes.NewBuffer(volBytes)
		volMftFiles := map[uint64]*db_cache.AgentMftCacheDB{}
		for {
			if readBuffer.Len() == 0 {
				break
			}
			indexBytes := make([]byte, 8)
			_, err = readBuffer.Read(indexBytes)
			if err != nil {
				break
			}
			index := byte2.BytesMustToUint(indexBytes, false)

			parentIndexBytes := make([]byte, 8)
			_, err = readBuffer.Read(parentIndexBytes)
			if err != nil {
				break
			}
			parentIndex := byte2.BytesMustToUint(parentIndexBytes, false)

			fileNameLengthBytes := make([]byte, 2)
			_, err = readBuffer.Read(fileNameLengthBytes)
			if err != nil {
				break
			}
			fileNameLength := byte2.BytesMustToUint(fileNameLengthBytes, false)

			fileNameBytes := make([]byte, fileNameLength)
			_, err = readBuffer.Read(fileNameBytes)
			if err != nil {
				break
			}
			fileName := string(fileNameBytes)

			createTimeBytes := make([]byte, 8)
			_, err = readBuffer.Read(createTimeBytes)
			if err != nil {
				break
			}
			createTime := byte2.BytesMustToUint(createTimeBytes, false)

			modifiedTimeBytes := make([]byte, 8)
			_, err = readBuffer.Read(modifiedTimeBytes)
			if err != nil {
				break
			}
			modifiedTime := byte2.BytesMustToUint(modifiedTimeBytes, false)

			fileSizeBytes := make([]byte, 8)
			_, err = readBuffer.Read(fileSizeBytes)
			if err != nil {
				break
			}
			fileSize := byte2.BytesMustToUint(fileSizeBytes, false)
			var isDir byte
			isDir, err = readBuffer.ReadByte()
			if err != nil {
				break
			}
			isDirBool := false
			if isDir != 0 {
				isDirBool = true
			}

			volMftFiles[index] = &db_cache.AgentMftCacheDB{
				Index:        fmt.Sprintf("%d", index),
				ParentIndex:  fmt.Sprintf("%d", parentIndex),
				FileName:     fileName,
				IsDir:        isDirBool,
				CreateTime:   fmt.Sprintf("%d", createTime),
				ModifiedTime: fmt.Sprintf("%d", modifiedTime),
				FileSize:     int64(fileSize),
			}
		}
		for _, file := range volMftFiles {
			//getPath(volMftFiles, index)
			mftFiles = append(mftFiles, file)
		}
	}
	start := time.Now()
	t.Logf("mft length: %d, start write database", len(mftFiles))
	db := db_cache.GetAgentMftCacheDB("test")
	// 关闭旧的
	db.Close()
	db = db_cache.GetAgentMftCacheDB("test")

	db.Delete()

	db.Exec(`
PRAGMA journal_mode = OFF; 
PRAGMA synchronous = OFF; 
PRAGMA temp_store = MEMORY; 
PRAGMA cache_size = -200; `)

	defer db.Exec(`
PRAGMA journal_mode = DELETE;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -20;
`)
	err = db.InsertMulti(mftFiles)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("write database over, use time: %v", time.Since(start))
}
