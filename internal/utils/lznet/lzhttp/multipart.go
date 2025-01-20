package lzhttp

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// --------------------multipart-------------------------------------------
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func CreateFormFile(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	if filename != "" {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(fieldname), escapeQuotes(filename)))
	} else {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"`,
				escapeQuotes(fieldname)))
	}
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}
	return w.CreatePart(h)
}

type PostMultiPart struct {
	FieldName   string
	FileName    string
	ContentType string
	Content     []byte
}

func (this *HttpClient) HttpPostMulti(target string, postMultiParts []PostMultiPart, headers Header, chunked bool) (httpresp HttpResp) {

	//if chunked {
	//	data = ChunkData(data, 10)
	//	headers["Transfer-Encoding"] = []string{"chunked"}
	//}
	//strings.NewReader()
	//reqBody := bytes.NewBufferString(data)
	// 构造 http client

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
