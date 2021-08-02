package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func concatURL(baseURL, endpoint string) string {
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = fmt.Sprintf("/%s", endpoint)
	}
	return fmt.Sprintf("%s%s", baseURL, endpoint)
}

func RandomUA() string {
	userAgentList := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.60",
	}
	return userAgentList[rand.Intn(len(userAgentList))]
}

func CreateJSONPayload(data map[string]interface{}) (body *bytes.Buffer, err error) {
	if data == nil {
		return
	}
	body = &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(data)
	return
}

func CreateFormPayload(textFields, fileFields map[string]string) (*bytes.Buffer, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	defer writer.Close()
	for key, value := range textFields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}
	for key, value := range fileFields {
		file, err := os.Open(value)
		if err != nil {
			return nil, err
		}
		fw, err := writer.CreateFormFile(key, filepath.Base(value))
		if _, err := io.Copy(fw, file); err != nil {
			return nil, err
		}
		if err := file.Close(); err != nil {
			return nil, err
		}
	}
	return payload, nil
}

func CopyBody(resp *http.Response) ([]byte, error) {
	b, err := ReadBody(resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return b, nil
}

func ReadBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Error(resp *http.Response) error {
	b, _ := ReadBody(resp)
	return fmt.Errorf("%s %s %d: %s", resp.Request.Method, resp.Request.URL.RawPath, resp.StatusCode, string(b))
}
