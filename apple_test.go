package oauth

import (
	"fmt"
	"testing"
)

func TestApple(t *testing.T) {
	service, err := NewService("com.short.roll", "RGGHW6A8T4", AuthApple)
	if nil != err {
		panic(err)
	}
	apple := NewApple(service)
	resp, err := apple.IdToken("eyJraWQiOiJXNldjT0tCIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLnNob3J0LnJvbGwiLCJleHAiOjE2OTg5OTYxODMsImlhdCI6MTY5ODkwOTc4Mywic3ViIjoiMDAxNTk3LjNlZmNlMjc5ZTc0ODQ5Y2U5MzZmMzhiNzI2YTliM2U1LjAzNDUiLCJjX2hhc2giOiJpQUx5Z2JrUXNwTHdENHA3dGp1dVR3IiwiZW1haWwiOiJyb2xsc2hvcnRAaWNsb3VkLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImF1dGhfdGltZSI6MTY5ODkwOTc4Mywibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.hLEA0WuaIQZp8MopfiHZ0p0VEwMM3tp7LdIaTcXTJbZiecf_G_ZG-mspgOAQ623T8YodA1BmBhnVs6wdhpIkR-eA0D6MmDj-djEB814NtjDRQ7v7SK2jAQXcJG5DifDCuXY-BJO73gnVdV80OBphqjtGwCtFUr7qJDb2syYeJhVJAq-un3yH-7UrOprmEEJ3zMMVoCIDS2zrsygj8Jrwmsgc1HnVXU816as7q8FXywZEnadOtKp2704RvNjMjHh_qBpWYzlApvJU4CsrwFXQqUa0uMwevoZS6QaEpLqYsGltt9xnHHOD_4j_tXVQxJ0s-qnZpA2mQ1ekSfzpVrygMQ")
	if nil != err {
		panic(err)
	}
	fmt.Println(resp)
}
