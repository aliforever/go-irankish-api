package irankish

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	ResponseCode string          `json:"responseCode"`
	Description  string          `json:"description"`
	Status       bool            `json:"status"`
	Result       json.RawMessage `json:"result"`
}

type MakeTokenResult struct {
	Token             string          `json:"token"`
	InitiateTimestamp int64           `json:"initiateTimeStamp"`
	ExpiryTimestamp   int64           `json:"expiryTimeStamp"`
	TransactionType   transactionType `json:"transactionType"`
	BillInfo          interface{}     `json:"billInfo"`

	merchantID string
}

func (m *MakeTokenResult) RedirectForm() string {
	form :=
		`<form id="redirectform" action="%s" method="POST">
				<input type="hidden" name="tokenIdentity" value="%s"/>
		</form>
		<script>document.forms["redirectform"].submit();</script>`

	return fmt.Sprintf(form, RedirectUrl, m.Token)
}

type TransactionResult struct {
	ResponseCode             string `json:"responseCode"`
	SystemTraceAuditNumber   string `json:"systemTraceAuditNumber"`
	RetrievalReferenceNumber string `json:"retrievalReferenceNumber"`
	TransactionDate          int64  `json:"transactionDate"`
	TransactionTime          int64  `json:"transactionTime"`
	ProcessCode              string `json:"processCode"`
	BillType                 string `json:"billType"`
	BillId                   string `json:"billId"`
	PaymentId                string `json:"paymentId"`
	Amount                   string `json:"amount"`
	DuplicateVerify          bool   `json:"duplicateVerify"` // indicates if payment is already verified
}
