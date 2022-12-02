package irankish

type billInfo struct {
	BillID        string `json:"billid"`
	BillPaymentID string `json:"billpaymentid"`
}

func NewBillInfo(billID string, billPaymentID string) billInfo {
	return billInfo{
		BillID:        billID,
		BillPaymentID: billPaymentID,
	}
}
