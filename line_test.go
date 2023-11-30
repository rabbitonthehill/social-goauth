package oauth

import (
	"fmt"
	"testing"
)

func TestLine(t *testing.T) {
	service, err := NewService("2000596845", "d8b512a384a343465202763eeea1a0e9", AuthLine, WithProxyURL("http://127.0.0.1:8001"))
	if nil != err {
		panic(err)
	}
	line, err := NewLine(service).UserInformation("aaaaa")
	if nil != err {
		panic(err)
	}
	fmt.Println(line)
}
