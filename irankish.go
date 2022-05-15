package irankish

import (
	"fmt"
	"net/http"
)

type IranKish struct {
	merchantID string
	sha1Key    string
}

func New(merchantID, sha1Key string) *IranKish {
	return &IranKish{merchantID: merchantID, sha1Key: sha1Key}
}

func (i *IranKish) NewToken(invoiceID string, amount int64, callbackUrl string) *token {
	return &token{
		merchantID:    i.merchantID,
		invoiceNumber: invoiceID,
		amount:        amount,
		callbackUrl:   callbackUrl,
	}
}

func (i *IranKish) SimpleFromRedirectingToGateway(token string) string {
	form := `<form id="redirectform" action="%s" method="POST"><input type="hidden" name="token" value="%s"/><input type="hidden" name="merchantId" value="%s"/></form><script>document.forms["redirectform"].submit();</script>`
	return fmt.Sprintf(form, gatewayUrl, token, i.merchantID)
}

func (i *IranKish) CallbackHandler(data chan<- *CallbackData) (handler func(w http.ResponseWriter, r *http.Request)) {
	handler = func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		cd := &CallbackData{
			Token:        r.Form.Get("token"),
			MerchantID:   r.Form.Get("merchantId"),
			AcceptorID:   r.Form.Get("acceptorId"),
			ResponseCode: r.Form.Get("responseCode"),
			RequestID:    r.Form.Get("requestId"),
			PaymentID:    r.Form.Get("paymentId"),
			Amount:       r.Form.Get("amount"),
			ReferenceID:  r.Form.Get("referenceId"),
			response:     make(chan []byte, 1),
		}

		go func() {
			data <- cd
		}()

		w.Write(<-cd.response)
	}
	return
}

func (i *IranKish) VerifyPayment(token, referenceNumber string) (*VerifyPaymentResult, error) {
	p := &payment{
		merchantID:      i.merchantID,
		sha1Key:         i.sha1Key,
		token:           token,
		referenceNumber: referenceNumber,
	}
	return p.verify()
}
