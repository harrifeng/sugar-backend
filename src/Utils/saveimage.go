package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var validImageFormat = []string{
	"jpeg",
}

func checkImageFormatValid(format string) bool {
	for i := 0; i < len(validImageFormat); i++ {
		if validImageFormat[i] == format {
			return true
		}
	}
	return false
}

func getHeaders(header string) (fileFormat string, err error) {
	var data, dataType, encodeMethod string
	i := 0
	k := i
	for ; i < len(header) && header[i] != ':'; i++ {
	}
	if i == len(header) {
		return fileFormat, errors.New("the format of image base64 is invalid")
	}
	data = header[k:i]
	k = i + 1
	for ; i < len(header) && header[i] != '/'; i++ {
	}
	if i == len(header) {
		return fileFormat, errors.New("the format of image base64 is invalid")
	}
	dataType = header[k:i]
	k = i + 1
	for ; i < len(header) && header[i] != ';'; i++ {
	}
	if i == len(header) {
		return fileFormat, errors.New("the format of image base64 is invalid")
	}
	fileFormat = header[k:i]
	k = i + 1
	for ; i < len(header); i++ {
	}
	encodeMethod = header[k:i]
	if data != "data" || dataType != "image" || encodeMethod != "base64" || !checkImageFormatValid(fileFormat) {
		return fileFormat, errors.New("the format of image base64 is invalid")
	}
	return fileFormat, nil
}

func DecodeWebImageBase64(imageBase64 string) (decodeResult []byte, fileFormat string, err error) {
	items := strings.Split(imageBase64, ",")
	if len(items) != 2 {
		err = errors.New("the format of image base64 is invalid")
		return decodeResult, fileFormat, err
	}
	headerPart := items[0]
	imagePart := items[1]
	fileFormat, err = getHeaders(headerPart)
	if err != nil {
		return decodeResult, fileFormat, err
	}
	decodeResult, err = base64.StdEncoding.DecodeString(imagePart)
	return decodeResult, fileFormat, err
}

func SaveImageFromBase64(imageBase64 string, dirPath string, fileName string) error {
	imageData, fileFormat, err := DecodeWebImageBase64(imageBase64)
	if err != nil {
		return err
	}
	filePath := filepath.Join(dirPath, fileName)
	filePath = fmt.Sprintf("%s.%s", filePath, fileFormat)
	err = ioutil.WriteFile(filePath, imageData, os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}
