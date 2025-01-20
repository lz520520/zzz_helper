package lzhttp

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"zzz_helper/internal/utils/lznet/lang"
)

type HttpResp struct {
	Resp        *http.Response
	Header      string
	Feature     string
	Title       string
	Body        string
	RawFullResp string
	RAWFullReq  string
	Err         error
}

type OptPayload struct {
	URI     string
	Headers Header
	Data    string
}

func ParseRequest(s string) (*http.Request, error) {
	reader := bufio.NewReader(strings.NewReader(s))
	return http.ReadRequest(reader)
}

func Request2String(r *http.Request) string {
	headerStr := fmt.Sprintf("Host: %s\r\n", r.Host)
	for k, v := range r.Header {
		for _, i := range v {
			headerStr += fmt.Sprintf("%s: %s\r\n", k, i)
		}
	}
	if r.ContentLength > 0 {
		headerStr += fmt.Sprintf("%s: %v\r\n", "Content-Length", r.ContentLength)
	}
	if r.TransferEncoding != nil {
		v := ""
		for _, vv := range r.TransferEncoding {
			v += vv + "; "
		}
		v = strings.TrimSuffix(v, "; ")
		headerStr += fmt.Sprintf("%s: %v\r\n", "TransferEncoding", v)
	}

	bodyStr := ""
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			bodyStr = string(bodyBytes)
			// 恢复body的缓冲区
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		}
	}
	s := fmt.Sprintf("%s %s %s\r\n%s\r\n%s", r.Method, r.URL.RequestURI(), r.Proto, headerStr, bodyStr)
	return s
}

func Response2String(r *http.Response, charset string) (httpFullData, headers, body string) {
	for k, v := range r.Header {
		for _, i := range v {
			headers += fmt.Sprintf("%s: %s\r\n", k, i)
		}
	}
	if r.ContentLength > 0 {
		headers += fmt.Sprintf("%s: %v\r\n", "Content-Length", r.ContentLength)
	}
	if r.TransferEncoding != nil {
		v := ""
		for _, vv := range r.TransferEncoding {
			v += vv + "; "
		}
		v = strings.TrimSuffix(v, "; ")
		headers += fmt.Sprintf("%s: %v\r\n", "TransferEncoding", v)
	}

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil || (bodyBytes != nil && len(bodyBytes) != 0) {
			body = string(bodyBytes)
			body = lang.LanguageCode(body, charset, lang.LangUTF8)
			r.ContentLength = int64(len(bodyBytes))
		}
	}
	httpFullData = fmt.Sprintf("%s %s\r\n%s\r\n%s", r.Proto, r.Status, headers, body)

	return
}

func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()
		undatas, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return undatas, nil
	}
}

func ParseHeader(s string) Header {
	s = strings.TrimSpace(s)
	reader := bufio.NewReader(strings.NewReader(s + "\r\n\r\n"))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil
	}

	// lzhttp.Header and lztextproto.MIMEHeader are both just a map[string][]string
	httpHeader := Header(mimeHeader)
	return httpHeader
}

func GetUrl2Host(target string) (host string, err error) {
	u, err := url.Parse(target)
	if err != nil {
		return
	}
	port := u.Port()
	if port <= "0" {
		if u.Scheme == "http" {
			port = "80"
		} else if u.Scheme == "https" {
			port = "443"
		}
	}
	host = fmt.Sprintf("%s://%s:%s", u.Scheme, u.Host, port)
	return
}
