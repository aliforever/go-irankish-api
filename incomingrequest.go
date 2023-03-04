package irankish

import "net/http"

type IncomingRequest struct {
	Request        *http.Request
	responseWriter http.ResponseWriter

	Done chan bool
}

func (i *IncomingRequest) WriteResponse(statusCode int, data []byte) {
	i.responseWriter.WriteHeader(statusCode)
	i.responseWriter.Write(data)

	i.Done <- true
}

func (i *IncomingRequest) ParseUserInput() (*UserInput, error) {
	if err := i.Request.ParseForm(); err != nil {
		return nil, err
	}

	return &UserInput{
		Token:                    i.Request.Form.Get("token"),
		AcceptorID:               i.Request.Form.Get("acceptorId"),
		MerchantID:               i.Request.Form.Get("merchantId"),
		ResponseCode:             i.Request.Form.Get("responseCode"),
		PaymentID:                i.Request.Form.Get("paymentId"),
		RequestID:                i.Request.Form.Get("requestId"),
		Sha256OfPan:              i.Request.Form.Get("sha256OfPan"),
		RetrievalReferenceNumber: i.Request.Form.Get("retrievalReferenceNumber"),
		Amount:                   i.Request.Form.Get("amount"),
		MaskedPan:                i.Request.Form.Get("maskedPan"),
		SystemTraceAuditNumber:   i.Request.Form.Get("systemTraceAuditNumber"),
		Ttl:                      i.Request.Form.Get("ttl"),
		Sha1OfPan:                i.Request.Form.Get("sha1OfPan"),
	}, nil
}

type UserInput struct {
	Token                    string
	AcceptorID               string
	MerchantID               string
	ResponseCode             string
	PaymentID                string
	RequestID                string
	Sha256OfPan              string
	RetrievalReferenceNumber string
	SystemTraceAuditNumber   string
	Amount                   string
	MaskedPan                interface{}
	Ttl                      string
	Sha1OfPan                interface{}
}
