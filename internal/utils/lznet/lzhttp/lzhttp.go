package lzhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"time"
	_ "unsafe"
)

//go:linkname reqWriteExcludeHeader net/http.reqWriteExcludeHeader
var reqWriteExcludeHeader map[string]bool

type RemoteInfo struct {
	IP string
}
type HttpClient struct {
	client *http.Client

	//target string
	param HttpClientInitParam

	remoteInfo RemoteInfo
}

type RequestAuth struct {
	AuthEnable bool   `json:"auth_enable" yaml:"auth_enable"`
	AuthType   string `json:"auth_type" yaml:"auth_type"` // Basic/NTLMv1/NTLMv2
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Domain     string `json:"domain" yaml:"domain"`
}
type HttpClientInitParam struct {
	Charset        string
	ConnTimeout    time.Duration
	TransTimeout   time.Duration
	Proxy          func(*http.Request) (*url.URL, error)
	ShortConn      bool
	RequestAuth    RequestAuth
	Redirect       bool
	DialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

func NewHttpClient(param HttpClientInitParam) *HttpClient {

	var transport http.RoundTripper
	if param.ShortConn {
		transport = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             param.Proxy,
			DisableKeepAlives: true,
			DialTLSContext:    param.DialTLSContext,
		}
	} else {
		transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   param.ConnTimeout,      // 设置连接超时为10秒
				KeepAlive: param.TransTimeout * 3, // 待定
			}).DialContext,
			MaxIdleConns:        25,
			MaxIdleConnsPerHost: 50,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			Proxy:               param.Proxy,
			DialTLSContext:      param.DialTLSContext,
		}
	}
	if param.RequestAuth.AuthEnable {
		switch param.RequestAuth.AuthType {
		case "Basic":
		case "NTLMv1":
		case "NTLMv2":
			transport = ntlmssp.Negotiator{
				RoundTripper: transport,
			}
		}
	}

	httpClient := &HttpClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   param.TransTimeout,
		},
		param: param,
	}
	if !param.Redirect {
		httpClient.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return httpClient
}

func (this *HttpClient) GetRemoteInfo() RemoteInfo {
	return this.remoteInfo
}

func (this *HttpClient) auth(req *http.Request) {
	if this.param.RequestAuth.AuthEnable {
		username := ""
		if this.param.RequestAuth.Domain != "" {
			username = fmt.Sprintf("%s\\%s", this.param.RequestAuth.Domain, this.param.RequestAuth.Username)
		} else {
			username = this.param.RequestAuth.Username
		}
		req.SetBasicAuth(username, this.param.RequestAuth.Password)
	}
}
func (this *HttpClient) HttpHead(target string, headers Header, chunked bool) (httpresp HttpResp) {
	req, err := http.NewRequest("HEAD", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	this.auth(req)

	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}

	defer resp.Body.Close()
	httpresp.Resp = resp

	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)
	return
}

func (this *HttpClient) HttpGet(target string, headers Header) (httpresp HttpResp) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	this.auth(req)
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}

	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)

	return
}

func (this *HttpClient) HttpDelete(target string, headers Header, chunked bool) (httpresp HttpResp) {
	req, err := http.NewRequest("DELETE", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	this.auth(req)

	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)

	return
}

func (this *HttpClient) HttpPost(target string, data string, headers Header, chunked bool) (httpresp HttpResp) {

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)

	req, err := http.NewRequest("POST", target, reqBody)

	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	this.auth(req)
	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}

	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)

	return
}

func (this *HttpClient) Close() error {
	if this.client != nil {
		this.client.CloseIdleConnections()
	}
	return nil
}
func (this *HttpClient) CustomHttp(req *http.Request) (httpresp HttpResp) {
	// 构造 http client
	req.RequestURI = ""

	httpresp.RAWFullReq = Request2String(req)
	this.auth(req)

	// 发送请求
	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)

	return
}

func (this *HttpClient) HttpPut(target string, data string, headers Header, chunked bool) (httpresp HttpResp) {
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)

	req, err := http.NewRequest("PUT", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	this.auth(req)

	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := this.client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, this.param.Charset)
	return
}

func HttpGetWithoutRedirect(target string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}

	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func HttpGet(target string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func HttpPostWithoutRedirect(target string, data string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", target, reqBody)

	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}

	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}
func HttpPost(target string, data string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func HttpDeleteWithoutRedirect(target string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("DELETE", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func HttpDelete(target string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("DELETE", target, nil)
	if err != nil {
		httpresp.Err = err
		return
	}
	//if chunked { req.TransferEncoding = []string{"ChunkEd"} }

	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}
	// 发送请求
	httpresp.RAWFullReq = Request2String(req)

	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}
func HttpPut(target string, data string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("PUT", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)
	return
}

func HttpPutWithoutRedirect(target string, data string, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	//strings.NewReader()
	reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("PUT", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	httpresp.RAWFullReq = Request2String(req)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)
	return
}

func HttpPostMulti(target string, postMultiParts []PostMultiPart, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	//if chunked {
	//	data = ChunkData(data, 10)
	//	headers["Transfer-Encoding"] = []string{"chunked"}
	//}
	//strings.NewReader()
	//reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
	// multipart构造
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	for _, pmp := range postMultiParts {
		if pmp.ContentType == "" {
			pmp.ContentType = "application/octet-stream"
		}
		part, err := CreateFormFile(writer, pmp.FieldName, pmp.FileName, pmp.ContentType)
		if err != nil {
			httpresp.Err = err
			return
		}
		if _, err = part.Write(pmp.Content); err != nil {
			httpresp.Err = err
			return
		}
	}
	writer.Close()
	headers.Set("Content-Type", writer.FormDataContentType())

	req, err := http.NewRequest("POST", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func HttpPostMultiWithoutRedirect(target string, postMultiParts []PostMultiPart, headers Header, charset string, timeout time.Duration, proxy func(*http.Request) (*url.URL, error), chunked bool) (httpresp HttpResp) {

	//if chunked {
	//	data = ChunkData(data, 10)
	//	headers["Transfer-Encoding"] = []string{"chunked"}
	//}
	//strings.NewReader()
	//reqBody := bytes.NewBufferString(data)
	// 构造 http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			Proxy:             proxy,
			DisableKeepAlives: true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// multipart构造
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	for _, pmp := range postMultiParts {
		if pmp.ContentType == "" {
			pmp.ContentType = "application/octet-stream"
		}
		part, err := CreateFormFile(writer, pmp.FieldName, pmp.FileName, pmp.ContentType)
		if err != nil {
			httpresp.Err = err
			return
		}
		if _, err = part.Write(pmp.Content); err != nil {
			httpresp.Err = err
			return
		}
	}
	writer.Close()
	headers.Set("Content-Type", writer.FormDataContentType())

	req, err := http.NewRequest("POST", target, reqBody)
	if err != nil {
		httpresp.Err = err
		return
	}

	//if chunked { req.TransferEncoding = []string{"chunked"} }
	//req.Header.Set()
	req.Header = http.Header(headers)
	if tmpHost := headers.Get("Host"); tmpHost != "" {
		req.Host = tmpHost
	}

	httpresp.RAWFullReq = Request2String(req)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		httpresp.Err = err
		return
	}
	defer resp.Body.Close()
	httpresp.Resp = resp

	// 接收请求
	httpresp.RawFullResp, httpresp.Header, httpresp.Body = Response2String(resp, charset)

	return
}

func init() {
	reqWriteExcludeHeader["Transfer-Encoding"] = false
}
