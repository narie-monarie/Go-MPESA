# Go-MPESA

NB: This project is still in development and works with golang 1.22+


```sh
go get -u github.com/narie-monarie/Go-MPESA@v1.0.1
```
```go

package main

import (
	"fmt"
	"net/http"
	"github.com/narie-monarie/Go-MPESA"
)

func main() {
	conf := mpesa.NewConfig(mpesa.Config{
			ConsumerKey:    "",
			ConsumerSecret: "",
			PassKey:        "",
		})

	http.HandleFunc("/stkPush", func(w http.ResponseWriter, r *http.Request) {
		params := mpesa.STKPushRequest{
			BusinessShortCode: "174379",
			TransactionType:   "CustomerPayBillOnline",
			Amount:            "20",
			PartyA:            "254712345678", // Your phone number
			PartyB:            "174379",
			PhoneNumber:       "254712345678", // Your phone number
			CallBackURL:       "https://yourwebsite.com", //u can use localhost and ngrok to get this one
			AccountReference:  "Test",
			TransactionDesc:   "PAYMENT OF GOODS",
		}
		conf.MPESAExpress(params, w, r)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
```
