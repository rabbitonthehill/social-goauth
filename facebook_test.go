package oauth

import (
	"fmt"
	"testing"
)

func TestFacebook(t *testing.T) {
	service, err := NewService("238852502115217", "742f1cbd03c3a68f89821a3bdc777a7d", AuthFacebook)
	if nil != err {
		panic(err)
	}
	fmt.Println(service)
}
