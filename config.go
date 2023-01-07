package irankish

import "net/url"

const (
	IranKishShaparakUrl = "https://ikc.shaparak.ir"
	TokenUrl            = "/api/v3/tokenization/make"
	RedirectUrl         = "/iuiv3/IPG/Index/"
	ConfirmationUrl     = "/api/v3/confirmation/purchase"
)

var host, _ = url.Parse(IranKishShaparakUrl)
