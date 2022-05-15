package irankish

type VerifyPaymentResult struct {
	Result string
}

func (v *VerifyPaymentResult) TranslateResultCode() string {
	if message, ok := verifyErrors[v.Result]; ok {
		return message
	}
	return "خطای نامعلوم"
}
