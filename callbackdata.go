package irankish

type CallbackData struct {
	Token         string
	MerchantID    string
	AcceptorID    string
	ResultCode    string
	InvoiceNumber string
	PaymentID     string
	Amount        string
	ReferenceID   string
	CardNo        string
	response      chan []byte
}

func (c *CallbackData) WasSuccessful() bool {
	return c.ResultCode == "100"
}

func (c *CallbackData) WriteResponse(data []byte) {
	go func() {
		c.response <- data
	}()
}

func (c *CallbackData) TranslateResultCode() string {
	if message, ok := callbackCodes[c.ResultCode]; ok {
		return message
	}

	return "وضعیت نامعلوم"
}
