package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveImageFromBase64(imageBase64 string, dirPath string, fileName string) error {
	decodeResult, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return err
	}
	filePath := filepath.Join(dirPath, fileName)
	filePath = fmt.Sprintf("%s.%s", filePath, "jpg")
	err = ioutil.WriteFile(filePath, decodeResult, os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}
