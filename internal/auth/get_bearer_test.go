package auth

import (
	"testing"
	"net/http"
)

func TestGetBearerToken(t *testing.T){
	headers := http.Header{
		"Authorization": []string{"Bearer a;ldj;lfkjdsojaoseflajwel;fja;lsjd",
						},
	}
	
	_, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("Unable to get token")
	}


}