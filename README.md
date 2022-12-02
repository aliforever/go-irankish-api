# IranKish Payment Gateway API for Go
The package is used for payments using Iranian IranKish payment gateway

For older version with SOAP check `soap` branch.

## Install
```go get -u github.com/aliforever/go-irankish-api```

## Usage
Variables:
```go
terminalID := "" // <- Place your terminal ID 
acceptorID := "" // <- Place your acceptor ID
passphrase := "" // <- Place your Passphrase
publicKey := "" // <- Place your Public Key

ik, err := irankish.New(terminalID, acceptorID, passphrase, publicKey)
if err != nil {
    panic(err)
}
```
- MakeToken for normal Purchase:
```go
paymentID := ""
requestID := ""
amount := 0
revertUri := ""

token, err := ik.MakePurchaseToken(paymentID, requestID, amount, revertUri)
if err != nil {
	panic(err)
}

fmt.Println(token)
```

After making token you can get a simple redirecting form as a html to redirect users to the gateway:
```go
htmlForm := token.RedirectForm() // <- Token Received from Make Token Method
```

Also, you can either handle incoming POST requests from the gateway yourself or listen to package's incoming callbacks. Here's an example to listen and verify payments using package's callback
```go
http.HandleFunc("/payment_callback", ik.CallbackHandler) // Registers a http handler for callbacks
go http.ListenAndServe(":8001", nil) // Run your http server

for request := range ik.IncomingCallbacks() {
    input, err := request.ParseUserInput()
    if err != nil {
        request.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
        return
    }
    
    fmt.Println(fmt.Sprintf("%+v", input))
    
    if input.RetrievalReferenceNumber != "" && input.SystemTraceAuditNumber != "" {
        result, err := ik.VerifyPurchase(input.Token, input.RetrievalReferenceNumber, input.SystemTraceAuditNumber)
        if err != nil {
            request.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
            fmt.Println(err)
            continue
        }
	
        fmt.Println(result)
    }
    
    request.WriteResponse(http.StatusOK, []byte("ok"))
}
```
-- Make sure to write the response or it will hang on the goroutine
- VerifyPayment:
```go
token := ""
referenceNumber := ""
auditNumber := ""

result, err := ik.VerifyPurchase(token, referenceNumber, auditNumber)
if err != nil {
    request.WriteResponse(http.StatusBadRequest, []byte(err.Error()))
    panic(err)
}

fmt.Println(result)
```

