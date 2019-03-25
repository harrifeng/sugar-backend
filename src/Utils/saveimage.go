package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func DecodeWebImageBase64(imageBase64 string) ([]byte, error) {
	var decodeResult []byte
	items := strings.Split(imageBase64, ",")
	if len(items) != 2 {
		return decodeResult, errors.New("the format of image base64 is invalid")
	}
	headerPart := items[0]
	imagePart := items[1]
	fmt.Println("header:")
	fmt.Println(headerPart)
	fmt.Println("image:")
	fmt.Println(imagePart)
	var data, dataType, fileFormat, encodeMethod string
	fmt.Println("split:")
	fmt.Printf("%s %s %s %s\n", data, dataType, fileFormat, encodeMethod)
	return base64.StdEncoding.DecodeString(imageBase64)
}

func SaveImageFromBase64(imageBase64 string, FilePath string) error {
	imageData, err := DecodeWebImageBase64(imageBase64)
	if err != nil {
		return err
	}
	fmt.Println(imageData)
	//err = ioutil.WriteFile(FilePath,imageData,os.ModeAppend)
	//if err!=nil{
	//	return err
	//}
	return nil
}
