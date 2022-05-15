package irankish

const (
	makeTokenXML = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://tempuri.org/">
    <SOAP-ENV:Body>
        <ns1:MakeToken>
            %tags%
        </ns1:MakeToken>
    </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	verifyPaymentXML = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://tempuri.org/">
    <SOAP-ENV:Body>
        <ns1:KicccPaymentsVerification>
            %tags%
        </ns1:KicccPaymentsVerification>
    </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
)
