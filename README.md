# IranKish Payment Gateway API for Go
The package is used for payments using Iranian IranKish payment gateway

## Install
```go get -u github.com/aliforever/go-irankish-api```

## Usage
Variables:
```go
merchantID := "" // <-- Replace
sha1Key := "" // <-- Replace

ik := irankish.New(merchanID, sha1Key)
```
- MakeToken:
```go
invoiceID := ""
amount := 0
callbackUrl := ""

token := ik.NewToken(invoiceID, amount, callbackUrl)

result, err := token.Make()
if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(result)
```
After making token you can write a simple redirecting form as a html to redirect users to the gateway:
```go
htmlForm := ik.SimpleFromRedirectingToGateway(token) // <- Token Received from Make Token Method
```

Also, you can either handle incoming POST requests from the gateway yourself or use a channel and package's handler:
```go
payments := make(chan *irankish.CallbackData)
go func() {
    for data := range payments {
        fmt.Println(data)
        data.WriteResponse([]byte("payment result"))
    }   
}

http.HandleFunc("/verify", ik.CallbackHandler(callbacks))
http.ListenAndServe(":8001", nil)
```
- VerifyPayment:
```go
token := ""
referenceNumber := ""

result, err := i.VerifyPayment(token, referenceNumber)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(result)
```

