package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func sendForm(ctx context.Context, endpoint string, data url.Values, headers map[string]string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return string(b), fmt.Errorf("sms api error: %s", string(b))
	}
	return string(b), nil
}

func SendSMS(to, message string) error {
	username := os.Getenv("AT_USERNAME")
	apiKey := os.Getenv("AT_API_KEY")
	base := os.Getenv("AT_BASEURL")
	if base == "" {
		base = "https://sandbox.africastalking.com/version1/messaging"
	}
	if username == "" || apiKey == "" {
		return fmt.Errorf("africastalking not configured")
	}
	data := url.Values{}
	data.Set("username", username)
	data.Set("to", to)
	data.Set("message", message)

	headers := map[string]string{
		"apiKey": apiKey, // some examples use apiKey header
	}
	_, err := sendForm(context.Background(), base, data, headers)
	return err
}
