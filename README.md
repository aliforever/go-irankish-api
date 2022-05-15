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
token := ik.NewToken(tt.args.invoiceID, tt.args.amount, tt.args.callbackUrl)

result, err := token.Make()
if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(result)
```

- VerifyPayment:
```go
result, err := i.VerifyPayment(tt.args.token, tt.args.referenceNumber)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(result)
```