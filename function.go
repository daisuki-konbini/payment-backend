package p

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
)

type Address struct {
	City  []string
	Line1 string
	Line2 string
	Name  string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	buf := make([]byte, 1024)
	n, _ := r.Body.Read(buf)
	address := new(Address)
	json.Unmarshal(buf[0:n], &address)
	session := getSession(address)
	ret, _ := json.Marshal(session)
	w.Write(ret)
}

func getSession(address *Address) *stripe.CheckoutSession {
	city := strings.Join(address.City, " ")
	stripe.Key = os.Getenv("StripeKey")
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(os.Getenv("StripePrice")),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Shipping: &stripe.ShippingDetailsParams{
				Name: stripe.String(address.Name),
				Address: &stripe.AddressParams{
					City:       stripe.String(city),
					Country:    stripe.String("日本"),
					Line1:      stripe.String(address.Line1),
					Line2:      stripe.String(address.Line2),
					PostalCode: stripe.String(""),
					State:      stripe.String(""),
				},
			},
		},
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String(os.Getenv("SuccessURL")),
		CancelURL:  stripe.String(os.Getenv("CancelURL")),
	}
	s, _ := session.New(params)
	return s
}
