package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	EmailValidatorKey = "EMAIL_VALIDATOR_KEY"
)

type ResponseParameter struct {
	EmailAddress string `json:"email_address"`
	Domain       string `json:"domain"`
	IsSyntax     string `json:"is_syntax"`
	IsDomain     string `json:"is_domain"`
	ISVerified   string `json:"is_verified"`
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func main() {
	var email string
	if _, err := fmt.Scan(&email); err != nil {
		log.Fatalf("reading user input error: %v\n", err)
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.mailboxvalidator.com/v1/validation/single", nil)
	if err != nil {
		log.Fatalf("creating request error: %v\n", err)
	}

	query := req.URL.Query()
	query.Add("email", email)
	query.Add("key", os.Getenv(EmailValidatorKey))
	query.Add("format", "json")
	req.URL.RawQuery = query.Encode()
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("calling api error: %v\n", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("reading response body error: %v\n", err)
	}

	var response ResponseParameter
	if err := json.Unmarshal(b, &response); err != nil {
		log.Fatalf(err.Error())
	}

	if strings.Compare(response.ErrorCode, "") != 0 || strings.Compare(response.ErrorMessage, "") != 0 {
		log.Fatalf("Error in calling api")
	}

	if strings.Compare(response.IsSyntax, "True") == 0 &&
		strings.Compare(response.Status, "True") == 0 &&
		strings.Compare(response.IsDomain, "True") == 0 &&
		strings.Compare(response.ISVerified, "True") == 0 {
		fmt.Printf("%s is valid email\n", response.EmailAddress)
	} else {
		fmt.Printf("%s is not valid email\n", response.EmailAddress)
	}
}
