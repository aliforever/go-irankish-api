package irankish

type multiplexParameters struct {
	params []multiplexParameter
}

type multiplexParameter struct {
	IBAN   string `json:"iban"`
	Amount int64  `json:"amount"`
}

func NewMultiplexParameters() multiplexParameters {
	return multiplexParameters{}
}

func (p multiplexParameters) Add(iban string, amount int64) multiplexParameters {
	p.params = append(p.params, multiplexParameter{
		IBAN:   iban,
		Amount: amount,
	})

	return p
}
