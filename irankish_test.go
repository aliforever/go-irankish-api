package irankish

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	merchantID = "" // <- Replace Yours
	sha1Key    = "" // <- Replace Yours
)

func TestIranKish_NewToken(t *testing.T) {
	type fields struct {
		merchantID string
	}
	type args struct {
		invoiceID   string
		amount      int64
		callbackUrl string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *token
	}{
		{
			name: "Test1",
			fields: fields{
				merchantID: merchantID,
			},
			args: args{
				invoiceID:   "0123456789abcdef",
				amount:      10000,
				callbackUrl: "http://getback.com",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.fields.merchantID, "")
			if got, err := i.NewToken(tt.args.invoiceID, tt.args.amount, tt.args.callbackUrl).Make(); err != nil {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			} else {
				fmt.Println(fmt.Sprintf("%+v", got))
				fmt.Println(i.SimpleFromRedirectingToGateway(got.Token))
			}
		})
	}
}

func TestIranKish_CallbackHandler(t *testing.T) {
	listener := make(chan *CallbackData)
	go func() {
		for data := range listener {
			fmt.Println(data.Amount, "111")
			data.WriteResponse([]byte("received"))
		}
	}()
	type fields struct {
		merchantID string
	}
	type args struct {
		data chan<- *CallbackData
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantHandler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Test1",
			fields: fields{
				merchantID: merchantID,
			},
			args:        args{data: listener},
			wantHandler: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.fields.merchantID, "")
			http.HandleFunc("/verify", i.CallbackHandler(tt.args.data))
		})
	}

	http.ListenAndServe(":8001", nil)
}

func TestIranKish_VerifyPayment(t *testing.T) {
	type fields struct {
		merchantID string
		sha1Key    string
	}
	type args struct {
		token           string
		referenceNumber int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VerifyPaymentResult
		wantErr bool
	}{
		{
			name: "Test1",
			fields: fields{
				merchantID: merchantID,
				sha1Key:    sha1Key,
			},
			args: args{
				token:           "7B5810F12FCAC64794E0CD5CAA3CAF107564",
				referenceNumber: 0,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IranKish{
				merchantID: tt.fields.merchantID,
				sha1Key:    tt.fields.sha1Key,
			}
			got, err := i.VerifyPayment(tt.args.token, tt.args.referenceNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyPayment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
