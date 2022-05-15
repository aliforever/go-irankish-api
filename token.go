package irankish

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type token struct {
	merchantID       string
	amount           int64
	invoiceNumber    string
	callbackUrl      string
	paymentID        string
	description      string
	specialPaymentID string
	extraParam1      string
	extraParam2      string
	extraParam3      string
	extraParam4      string
}

func (t *token) xml() (xml string, err error) {
	if t.merchantID == "" {
		err = emptyMerchantID
		return
	}

	if t.amount == 0 {
		err = emptyAmount
		return
	}

	if t.invoiceNumber == "" {
		err = emptyInvoiceID
		return
	}

	if t.callbackUrl == "" {
		err = emptyCallbackUrl
		return
	}

	var tags []string
	tags = append(tags, fmt.Sprintf("<ns1:amount>%d</ns1:amount>", t.amount))
	tags = append(tags, fmt.Sprintf("<ns1:merchantId>%s</ns1:merchantId>", t.merchantID))
	tags = append(tags, fmt.Sprintf("<ns1:invoiceNo>%s</ns1:invoiceNo>", t.invoiceNumber))
	tags = append(tags, fmt.Sprintf("<ns1:revertURL>%s</ns1:revertURL>", t.callbackUrl))

	if t.paymentID != "" {
		tags = append(tags, fmt.Sprintf("<ns1:paymentId>%s</ns1:paymentId>", t.paymentID))
	}

	if t.specialPaymentID != "" {
		tags = append(tags, fmt.Sprintf("<ns1:specialPaymentId>%s</ns1:specialPaymentId>", t.specialPaymentID))
	}

	if t.description != "" {
		tags = append(tags, fmt.Sprintf("<ns1:description>%s</ns1:description>", t.description))
	}

	xml = strings.ReplaceAll(makeTokenXML, "%tags%", strings.Join(tags, "\n"))
	return
}

func (t *token) SetPaymentID(paymentID string) *token {
	t.paymentID = paymentID
	return t
}

func (t *token) SetDescription(desc string) *token {
	t.description = desc
	return t
}

func (t *token) SetSpecialPaymentID(specialPaymentID string) *token {
	t.specialPaymentID = specialPaymentID
	return t
}

func (t *token) SetExtraParam1(param string) *token {
	t.extraParam1 = param
	return t
}

func (t *token) SetExtraParam2(param string) *token {
	t.extraParam2 = param
	return t
}

func (t *token) SetExtraParam3(param string) *token {
	t.extraParam3 = param
	return t
}

func (t *token) SetExtraParam4(param string) *token {
	t.extraParam4 = param
	return t
}

func (t *token) Make() (mtr *MakeTokenResult, err error) {
	var xml string
	xml, err = t.xml()
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", makeTokenUrl, strings.NewReader(xml))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", `"http://tempuri.org/ITokens/MakeToken"`)

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	mtr, err = t.parseResult(string(b))
	return
}

func (t *token) parseResult(response string) (mtr *MakeTokenResult, err error) {
	if !strings.Contains(response, "<MakeTokenResult ") {
		err = makeTokenWrongResponse(response)
		return
	}

	mtr = &MakeTokenResult{}

	messageRegex := regexp.MustCompile(`<a:message.*?>(.+)</a:message>`)
	message := messageRegex.FindStringSubmatch(response)
	if len(message) > 1 {
		mtr.Message = strings.TrimSpace(message[1])
	}

	resultRegex := regexp.MustCompile(`<a:result>(.+)</a:result>`)
	result := resultRegex.FindStringSubmatch(response)
	if len(result) > 0 {
		if result[1] == "true" {
			mtr.Result = true
		} else {
			mtr.Result = false
		}
	}

	tokenRegex := regexp.MustCompile(`<a:token>(.+)</a:token>`)
	token := tokenRegex.FindStringSubmatch(response)
	if len(token) > 0 {
		mtr.Token = token[1]
	}

	return
}
