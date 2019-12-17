package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var port = flag.Int("port", 9102, "port to listen on")

type authToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {
	flag.Parse()
	http.HandleFunc("/auth_callback", authServer)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func authServer(w http.ResponseWriter, r *http.Request) {
	username, password, err := getGCENodeCredentials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "username=%s\npassword=%s", username, password)
}

func getGCENodeCredentials() (string, string, error) {
	emailRequest := "instance/service-accounts/default/email"
	email, err := doMetadataGet(emailRequest)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch email, err: %v", err)
	}

	tokenRequest := "instance/service-accounts/default/token"
	token, err := doMetadataGet(tokenRequest)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch token, err: %v", err)
	}

	at := authToken{}
	if err = json.Unmarshal([]byte(token), &at); err != nil {
		return "", "", fmt.Errorf("unmarshal json failed: %s, err: %v", token, err)
	}
	return email, at.AccessToken, nil
}

func doMetadataGet(suffix string) (string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 1,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", "http://metadata/computeMetadata/v1/"+suffix, nil)
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := netClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata %s, err: %v", suffix, err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code from metadata, resp: %v", resp)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from metadata %s, err: %v", suffix, err)
	}
	return string(body), nil
}
