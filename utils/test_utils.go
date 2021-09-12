package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func HackResponseWriter(c *gin.Context) *responseBodyWriter {
	w := &responseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
	return w
}

// LoadFileAsStringFromPath load file data from file path as string
func LoadFileAsStringFromPath(filePath string) string {
	path := GetRelativePath(filePath)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open file  %s, error : %s", filePath, err.Error())
		return ""
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Unable to read  %s, error : %s", filePath, err.Error())
		return ""
	}

	return string(content)
}

func GetRelativePath(filePath string) string {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../..")
	return root + filePath
}

func FormatJsonResponse(response string) string {
	stringResponse := strings.Replace(response, "\n", "", -1)
	stringResponse = strings.Replace(stringResponse, "  ", "", -1)
	return stringResponse
}
