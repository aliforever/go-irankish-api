package irankish

type CallbackData struct {
	Token        string
	MerchantID   string
	AcceptorID   string
	ResponseCode string
	RequestID    string
	PaymentID    string
	Amount       string
	ReferenceID  string
	response     chan []byte
}

func (c *CallbackData) WriteResponse(data []byte) {
	c.response <- data
}

func (c *CallbackData) TranslateResponseCode() string {
	if message, ok := callbackCodes[c.ResponseCode]; ok {
		return message
	}

	return "وضعیت نامعلوم"
}
