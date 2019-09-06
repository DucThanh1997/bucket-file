package utils

import (
	"fmt"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenUUID() string {
	id := uuid.NewV4()
	return strings.Replace(fmt.Sprintf("%v", id), "-", "", -1)
}

func IsStringInArray(str string, arr []string) bool {
	for _, i := range arr {
		if i == str {
			return true
		}
	}
	return false
}

func Mkdir(folder string) error {
	_, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
