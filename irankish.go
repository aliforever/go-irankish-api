package irankish

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aliforever/encryptionbox"
)

type IranKish struct {
	terminalID string
	acceptorID string
	passphrase string
	publicKey  *rsa.PublicKey
	logger     Logger
	callbacks  chan IncomingRequest
	host       *url.URL
}

func New(terminalID, acceptorID, passphrase, publicKey string, logger Logger) (*IranKish, error) {
	pKey, err := encryptionbox.EncryptionBox{}.RSA.PublicKeyFromPKIXPEMBytes([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	return &IranKish{
		terminalID: terminalID,
		acceptorID: acceptorID,
		passphrase: passphrase,
		publicKey:  pKey,
		callbacks:  make(chan IncomingRequest),
		logger:     logger,
		host:       host}, nil
}

func NewWithProxyHost(
	terminalID,
	acceptorID,
	passphrase,
	publicKey string,
	proxyAddress string,
	logger Logger,
) (*IranKish, error) {

	pKey, err := encryptionbox.EncryptionBox{}.RSA.PublicKeyFromPKIXPEMBytes([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	proxyUrl, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, err
	}

	return &IranKish{
		terminalID: terminalID,
		acceptorID: acceptorID,
		passphrase: passphrase,
		publicKey:  pKey,
		callbacks:  make(chan IncomingRequest),
		logger:     logger,
		host:       proxyUrl}, nil
}

func (i *IranKish) IncomingCallbacks() <-chan IncomingRequest {
	return i.callbacks
}

func (i *IranKish) CallbackHandler(wr http.ResponseWriter, r *http.Request) {
	ir := IncomingRequest{
		Request:        r,
		responseWriter: wr,
		Done:           make(chan bool),
	}

	i.callbacks <- ir

	<-ir.Done
}

func (i *IranKish) MakePurchaseToken(
	paymentID,
	requestID string,
	amount int64,
	revertUri string,
	params ...AdditionalParameter,
) (*MakeTokenResult, error) {

	token := newMakeToken().
		SetPaymentID(paymentID).
		SetRequestID(requestID).
		SetTerminalID(i.terminalID).
		SetAcceptorID(i.acceptorID).
		SetRequestTimestampNow().
		SetRevertUri(revertUri).
		SetTransactionTypePurchase().
		SetAmount(amount)

	if len(params) > 0 {
		token.SetAdditionalParameters(params...)
	}

	iv, data, err := i.createAuthenticationEnvelopeHex(token.Amount)
	if err != nil {
		return nil, err
	}

	token.SetTerminalID(i.terminalID)

	body := map[string]interface{}{
		"request": token.Build(),
		"authenticationEnvelope": map[string]interface{}{
			"iv":   iv,
			"data": data,
		},
	}

	j, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", i.host.String()+TokenUrl, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if i.logger != nil {
		fmt.Println(string(result))
		go i.logger.Println(string(result))
	}

	var r *Response

	err = json.Unmarshal(result, &r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, fmt.Errorf("%s - %s", r.ResponseCode, r.Description)
	}

	var makeTokenResult *MakeTokenResult

	err = json.Unmarshal(r.Result, &makeTokenResult)
	if err != nil {
		return nil, err
	}

	return makeTokenResult, nil
}

func (i *IranKish) VerifyPurchase(token, referenceNumber, auditNumber string) (*TransactionResult, error) {
	body := map[string]interface{}{
		"terminalId":               i.terminalID,
		"retrievalReferenceNumber": referenceNumber,
		"systemTraceAuditNumber":   auditNumber,
		"tokenIdentity":            token,
	}

	j, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", i.host.String()+ConfirmationUrl, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if i.logger != nil {
		go i.logger.Println(string(result))
	}

	var r *Response
	err = json.Unmarshal(result, &r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, fmt.Errorf("%s - %s", r.ResponseCode, r.Description)
	}

	var makeTokenResult *TransactionResult
	err = json.Unmarshal(r.Result, &makeTokenResult)
	if err != nil {
		return nil, err
	}

	return makeTokenResult, nil
}

func (i *IranKish) SingleInquiryByReferenceNumber(referenceNumber string) (*InquiryResult, error) {
	return i.singleInquiry(&referenceNumber, nil, nil)
}

func (i *IranKish) SingleInquiryByToken(token string) (*InquiryResult, error) {
	return i.singleInquiry(nil, &token, nil)
}

func (i *IranKish) singleInquiry(referenceNumber *string, token *string, requestID *string) (*InquiryResult, error) {
	payload := newInquiry().
		SetPassPhrase(i.passphrase).
		SetTerminalID(i.terminalID)

	if referenceNumber != nil {
		payload.SetByReferenceNumber(*referenceNumber)
	} else if token != nil {
		payload.SetByPaymentToken(*token)
	} else if requestID != nil {
		payload.SetByRequestID(*requestID)
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", i.host.String()+InquiryUrl, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if i.logger != nil {
		go i.logger.Println(string(result))
	}

	var r *Response
	err = json.Unmarshal(result, &r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, fmt.Errorf("%s - %s", r.ResponseCode, r.Description)
	}

	var inquiryResult *InquiryResult
	err = json.Unmarshal(r.Result, &inquiryResult)
	if err != nil {
		return nil, err
	}

	return inquiryResult, nil
}

func (i *IranKish) createAuthenticationEnvelopeHex(amount int64) (iv, data string, err error) {
	str := fmt.Sprintf("%s%s%012d00", i.terminalID, i.passphrase, amount)

	binStr, err := hex.DecodeString(str)
	if err != nil {
		return "", "", err
	}

	key, ivBytes, encrypted, err := encryptionbox.EncryptionBox{}.AES.EncryptCbcPkcs5PaddingIvKey16Bytes(binStr)
	if err != nil {
		return "", "", err
	}

	hashWriter := sha256.New()
	hashWriter.Write(encrypted)

	finalData := make([]byte, 48)
	copy(finalData[:16], key)
	copy(finalData[16:], hashWriter.Sum(nil))

	result, err := encryptionbox.EncryptionBox{}.RSA.PublicKeyEncryptPKCS1v15(i.publicKey, finalData)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(ivBytes), hex.EncodeToString(result), nil
}
