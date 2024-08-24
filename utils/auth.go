package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetAccessToken(authCode string) ([]byte, error) {
	exURL := os.Getenv("EXCHANGE_URL")
	redirectURI := os.Getenv("REDIRECT_URI")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	formData := map[string]string{
		"grant_type":   "authorization_code",
		"code":         authCode,
		"redirect_uri": redirectURI,
	}

	//form object.
	form := url.Values{}

	for k, v := range formData {
		form.Set(k, v)
	}

	encodedForm := form.Encode()
	body := strings.NewReader(encodedForm)

	//Make http request to url.
	req, err := http.NewRequest(http.MethodPost, exURL, body)
	if err != nil {
		return nil, err
	}

	//Set headers.
	req.Header.Set("Content_Type", " application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	//send request.
	client := &http.Client{}
	response, err1 := client.Do(req)
	if err1 != nil {
		return nil, err1
	}
	defer response.Body.Close()

	//Read response body.
	respBody, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		return nil, err2
	}

	return respBody, nil
}
