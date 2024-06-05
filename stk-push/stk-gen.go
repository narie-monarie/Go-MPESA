package stkpush

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"narie-monarie/token"
	"net/http"
	"strings"
	"time"
)

type RequestData struct {
	PhoneNumber int `json:"phone_number"`
	Amount      int `json:"amount"`
}

func GetSTKPush(w http.ResponseWriter, r *http.Request) {
	accessToken, err := token.GetAccessToken()
	if err != nil {
		http.Error(w, "Error getting access token", http.StatusInternalServerError)
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	phoneNumber := requestData.PhoneNumber
	amount := requestData.Amount

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

	businessShortCode := 174379
	// Enter your phoneNumber below
	partyA := 254712345678
	partyB := 174379

	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)

	passKey := ""
	password := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%s%s", businessShortCode, passKey, timestamp)))

	payload := strings.NewReader(fmt.Sprintf(`{
      "BusinessShortCode": %d,
      "Password": "%s",
      "Timestamp": "%s",
      "TransactionType": "CustomerPayBillOnline",
      "Amount": %d,
      "PartyA": %d,
      "PartyB": %d,
      "PhoneNumber": %d,
    	"CallBackURL": "",
      "AccountReference": "CompanyXLTD",
      "TransactionDesc": "Payment of Boogs" 
    }`, businessShortCode, password, timestamp, amount, partyA, partyB, phoneNumber))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}
