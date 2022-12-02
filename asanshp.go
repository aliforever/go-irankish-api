package irankish

type asanShp struct {
	PrepaymentAmount int64 `json:"prepaymentAmount"`
	LoanAmount       int64 `json:"loanAmount"`
	LoadCount        int64 `json:"loadCount"`
}

func NewAsanShp(prepaymentAmount, loanAmount, loanCount int64) asanShp {
	return asanShp{
		PrepaymentAmount: prepaymentAmount,
		LoanAmount:       loanAmount,
		LoadCount:        loanCount,
	}
}
