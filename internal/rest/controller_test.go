package rest

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rntrp/bimg-rest/internal/config"
)

func init() {
	config.Load()
}

//go:embed test.png
var img []byte

func testInitPostReq(query string) *http.Request {
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	formPart, _ := multipartWriter.CreateFormFile("img", "test.png")
	formPart.Write(img)
	multipartWriter.Close()
	req := httptest.NewRequest("POST", query, body)
	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	return req
}

func TestScale(t *testing.T) {
	rec := httptest.NewRecorder()
	req := testInitPostReq("/convert?width=2&height=2&format=jpeg")

	Convert(rec, req)

	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	outMagic := make([]byte, 3)
	rec.Body.Read(outMagic)
	jpgMagic, _ := hex.DecodeString("ffd8ff")
	if !bytes.Equal(outMagic, jpgMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body JPEG Magic Number = %v; want ffd8ff", got)
	}
}
