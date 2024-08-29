package irankish

type inquiry struct {
	PassPhrase               string  `json:"passPhrase"`
	TerminalID               string  `json:"terminalId"`
	RetrievalReferenceNumber *string `json:"retrievalReferenceNumber,omitempty"`
	TokenIdentity            *string `json:"tokenIdentity,omitempty"`
	RequestID                *string `json:"requestId,omitempty"`
	FindOption               int     `json:"findOption"`
}

func newInquiry() *inquiry {
	return &inquiry{}
}

func (m *inquiry) SetPassPhrase(passPhrase string) *inquiry {
	m.PassPhrase = passPhrase

	return m
}

func (m *inquiry) SetByReferenceNumber(retrievalReferenceNumber string) *inquiry {
	m.FindOption = 1

	m.RetrievalReferenceNumber = &retrievalReferenceNumber

	return m
}

func (m *inquiry) SetByPaymentToken(paymentToken string) *inquiry {
	m.FindOption = 2

	m.TokenIdentity = &paymentToken

	return m
}

func (m *inquiry) SetByRequestID(requestID string) *inquiry {
	m.FindOption = 3

	m.RequestID = &requestID

	return m
}

func (m *inquiry) SetTerminalID(terminalID string) *inquiry {
	m.TerminalID = terminalID

	return m
}

func (m *inquiry) Build() inquiry {
	return *m
}
