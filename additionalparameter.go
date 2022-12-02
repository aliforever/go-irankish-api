package irankish

type additionalParameters struct {
	params map[string]string
}

func NewAdditionalParameters() additionalParameters {
	return additionalParameters{params: map[string]string{}}
}

func (p additionalParameters) Add(key, value string) additionalParameters {
	p.params[key] = value

	return p
}
