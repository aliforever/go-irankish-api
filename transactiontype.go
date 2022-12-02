package irankish

type transactionType string

const (
	TransactionTypePurchase  transactionType = "Purchase"
	TransactionTypeBill      transactionType = "Bill"
	TransactionTypeAsnShpWPP transactionType = "AsanShpWPP"
)
