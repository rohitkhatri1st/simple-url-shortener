package app

import (
	"fmt"
	"simple-url-shortener/model"
)

func GetServerAddress() string {
	return fmt.Sprintf("%s://%s:%s", model.ServerScheme, model.ServerHost, model.ServerPort)
}
