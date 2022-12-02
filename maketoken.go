package irankish

import "time"

type makeToken struct {
	/*
		// "billInfo":             "",
		// "multiplexParameters":  "",
		// "asanShp":              "",
	*/
	TransactionType      transactionType       `json:"transactionType"`
	TerminalId           string                `json:"terminalId"`
	AcceptorId           string                `json:"acceptorId"`
	Amount               int64                 `json:"amount"`
	RevertUri            string                `json:"revertUri"`
	RequestId            string                `json:"requestId"`
	PaymentId            string                `json:"paymentId"`
	CmsPreservationId    string                `json:"cmsPreservationId"`
	RequestTimestamp     int64                 `json:"requestTimestamp"`
	AdditionalParameters *additionalParameters `json:"additionalParameters,omitempty"`
}

func newMakeToken() *makeToken {
	return &makeToken{}
}

func (m *makeToken) SetTransactionTypePurchase() *makeToken {
	m.TransactionType = TransactionTypePurchase

	return m
}

func (m *makeToken) SetTerminalID(terminalID string) *makeToken {
	m.TerminalId = terminalID

	return m
}

func (m *makeToken) SetAcceptorID(acceptorID string) *makeToken {
	m.AcceptorId = acceptorID

	return m
}

func (m *makeToken) SetAmount(amount int64) *makeToken {
	m.Amount = amount

	return m
}

func (m *makeToken) SetRevertUri(revertUri string) *makeToken {
	m.RevertUri = revertUri

	return m
}

func (m *makeToken) SetRequestID(requestId string) *makeToken {
	m.RequestId = requestId

	return m
}

func (m *makeToken) SetPaymentID(paymentID string) *makeToken {
	m.PaymentId = paymentID

	return m
}

func (m *makeToken) SetRequestTimestamp(t int64) *makeToken {
	m.RequestTimestamp = t

	return m
}

func (m *makeToken) SetRequestTimestampNow() *makeToken {
	m.RequestTimestamp = time.Now().Unix()

	return m
}

func (m *makeToken) SetAdditionalParameters(parameters additionalParameters) *makeToken {
	m.AdditionalParameters = &parameters

	return m
}

func (m *makeToken) Build() makeToken {
	return *m
}
