package oauth

import (
	"fmt"
	"testing"
)

func TestGoogle(t *testing.T) {
	service, err := NewService("276609465331-6b865bo5hn43tirug01ef895n45vh001.apps.googleusercontent.com",
		"AIzaSyDH5G_xq7T44GZ7xtwOIAzV1X_zLgaFV1s",
		AuthGoogle,
	)
	if nil != err {
		panic(err)
	}
	fmt.Println(service)
}
