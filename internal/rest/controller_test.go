package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math"
	"net/http/httptest"
	"testing"
)

func TestScale(t *testing.T) {
	rec := httptest.NewRecorder()

	const pngB64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
	pngBytes, _ := base64.StdEncoding.DecodeString(pngB64)
	body := []byte("--boundary\r\nContent-Disposition: form-data; name=\"image\"; filename=\"image\"\r\n\r\n")
	body = append(body, pngBytes...)
	body = append(body, []byte("\r\n--boundary--")...)
	req := httptest.NewRequest("POST", "/scale?width=2&height=2&format=png", bytes.NewReader(body))
	req.Header.Add("Content-Type", "multipart/form-data; boundary=\"boundary\"")

	Scale(rec, req)

	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	out := rec.Body.Bytes()
	outMagic := out[:int(math.Min(8, float64(len(out))))]
	pngMagic, _ := hex.DecodeString("89504e470d0a1a0a")
	if !bytes.Equal(outMagic, pngMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body PNG Magic Number = %v; want 89504e470d0a1a0a", got)
	}
}
