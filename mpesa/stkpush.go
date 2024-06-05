package mpesa

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// STKPushRequest represents the parameters for an STK push request.
type STKPushRequest struct {
	BusinessShortCode string
	TransactionType   string
	Amount            string
	PartyA            string
	PartyB            string
	PhoneNumber       string
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

// MPESAExpress performs an MPESA Express (STK Push) request.
func (c *Config) MPESAExpress(params STKPushRequest, w http.ResponseWriter, r *http.Request) {
	accessToken, err := c.GetAuth()
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		now.Year(),
		int(now.Month()),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)

	password := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s%s", params.BusinessShortCode, c.PassKey, timestamp)))

	payload := strings.NewReader(fmt.Sprintf(`{
		"BusinessShortCode": "%s",
		"Password": "%s",
		"Timestamp": "%s",
		"TransactionType": "%s",
		"Amount": "%s",
		"PartyA": "%s",
		"PartyB": "%s",
		"PhoneNumber": "%s",
		"CallBackURL": "%s",
		"AccountReference": "%s",
		"TransactionDesc": "%s"
	}`, params.BusinessShortCode, password, timestamp, params.TransactionType, params.Amount, params.PartyA, params.PartyB, params.PhoneNumber, params.CallBackURL, params.AccountReference, params.TransactionDesc))

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest", payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
