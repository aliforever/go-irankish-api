package irankish

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aliforever/encryptionbox"
	"io"
	"net/http"
)

type IranKish struct {
	terminalID string
	acceptorID string
	passphrase string
	publicKey  *rsa.PublicKey

	callbacks chan IncomingRequest
}

func New(terminalID, acceptorID, passphrase, publicKey string) (*IranKish, error) {
	pKey, err := encryptionbox.EncryptionBox{}.RSA.PublicKeyFromPKIXPEMBytes([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	return &IranKish{terminalID: terminalID, acceptorID: acceptorID, passphrase: passphrase, publicKey: pKey, callbacks: make(chan IncomingRequest)}, nil
}

func (i *IranKish) IncomingCallbacks() chan IncomingRequest {
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

func (i *IranKish) MakePurchaseToken(paymentID, requestID string, amount int64, revertUri string, params ...AdditionalParameter) (*MakeTokenResult, error) {
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

	resp, err := http.Post(TokenUrl, "application/json", bytes.NewReader(j))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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

	resp, err := http.Post(ConfirmationUrl, "application/json", bytes.NewReader(j))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
