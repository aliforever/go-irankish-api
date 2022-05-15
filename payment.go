package irankish

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type payment struct {
	merchantID      string
	sha1Key         string
	token           string
	referenceNumber string
}

func (p *payment) xml() (xml string, err error) {
	if p.merchantID == "" {
		err = emptyMerchantID
		return
	}

	if p.token == "" {
		err = emptyToken
		return
	}

	if p.sha1Key == "" {
		err = emptySha1Key
		return
	}

	if p.referenceNumber == "" {
		err = emptyReferenceNumber
		return
	}

	var tags []string
	tags = append(tags, fmt.Sprintf("<ns1:token>%s</ns1:token>", p.token))
	tags = append(tags, fmt.Sprintf("<ns1:merchantId>%s</ns1:merchantId>", p.merchantID))
	tags = append(tags, fmt.Sprintf("<ns1:referenceNumber>%s</ns1:referenceNumber>", p.referenceNumber))
	tags = append(tags, fmt.Sprintf("<ns1:sha1Key>%s</ns1:sha1Key>", p.sha1Key))

	xml = strings.ReplaceAll(verifyPaymentXML, "%tags%", strings.Join(tags, "\n"))
	return
}

func (p *payment) verify() (vpr *VerifyPaymentResult, err error) {
	var xml string
	xml, err = p.xml()
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", verifyPaymentUrl, strings.NewReader(xml))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", `"http://tempuri.org/IVerify/KicccPaymentsVerification"`)

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

	vpr, err = p.parseResult(string(b))
	return
}

func (p *payment) parseResult(response string) (vpr *VerifyPaymentResult, err error) {
	if !strings.Contains(response, "<KicccPaymentsVerificationResult") {
		err = verifyPaymentWrongResponse(response)
		return
	}

	vpr = &VerifyPaymentResult{}
	verificationResultTag := []string{"<KicccPaymentsVerificationResult>", "</KicccPaymentsVerificationResult>"}
	verificationResultStartTagIndex := strings.Index(response, verificationResultTag[0])
	verificationResult := response[verificationResultStartTagIndex+len(verificationResultTag[0]) : strings.Index(response, verificationResultTag[1])]
	vpr.Result = verificationResult
	return
}
