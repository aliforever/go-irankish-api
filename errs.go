package irankish

import "errors"

var (
	emptyMerchantID            = errors.New("empty_merchant_id")
	emptyAmount                = errors.New("empty_amount")
	emptyInvoiceID             = errors.New("empty_invoice_id")
	emptyCallbackUrl           = errors.New("empty_callback_url")
	makeTokenWrongResponse     = func(response string) error { return errors.New("invalid_response_" + response) }
	verifyPaymentWrongResponse = func(response string) error { return errors.New("invalid_response_" + response) }
	emptyToken                 = errors.New("empty_token")
	emptySha1Key               = errors.New("empty_sha1key")
	emptyReferenceNumber       = errors.New("empty_reference_number")
)

var verifyErrors = map[string]string{
	"-20": "وجود کاراکترهای غیرمجاز در درخواست",
	"-30": "تراکنش قبلا برگشت خورده است",
	"-50": "طول رشته درخواست غیرمجاز است",
	"-51": "خطا در درخواست",
	"-80": "تراکنش مورد نظر یافت نشد",
	"-81": "خطای داخلی بانک",
	"-90": "تراکنش قبلا تایید شده است",
}

var callbackCodes = map[string]string{
	"100": "تراکنش با موفقیت انجام شد",
	"110": "خریدار انصراف داده است",
	"120": "موجودی کافی نیست",
	"130": "اطلاعات کارت اشتباه است",
	"131": "رمز کارت اشتباه است",
	"132": "کارت مسدود شده است",
	"133": "کارت منقضی شده است",
	"140": "زمان مورد نظر به پایان رسیده است",
	"150": "خطای داخلی بانک",
	"160": "خطا در اطلاعات کارت",
	"166": "بانک صادر کننده مجور انجام تراکنش را صادر نکرده است",
	"200": "مبلغ تراکنش بیش از سقف مجاز برای هر تراکنش می باشد",
	"201": "مبلغ تراکنش بیشتر از سقف مجاز در روز می باشد",
	"202": "مبلغ تراکنش بیشتر از سقف مجاز در ماه می باشد",
}
