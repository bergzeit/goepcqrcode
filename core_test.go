package goepcqrcode_test

import (
	"testing"

	"github.com/bergzeit/goepcqrcode/core"
)

// Mock data for testing
var (
	mockConfig = core.Config{
		Version:            core.VERSION_1,
		CharacterSet:       core.CHAR_SET_UTF8,
		IdentificationCode: core.IDENTIFICATION_CODE_SCT,
		BIC:                "YOURBIC",
		Name:               "John Doe",
		IBAN:               "DE89370400440532013000",
		Currency:           "EUR",
	}

	mockTransfer = core.Transfer{
		Amount:                "123.45",
		Purpose:               core.PURPOSE_SALA,
		RemittanceInformation: "Payment for invoice 12345",
	}
)

func TestGetRawText(t *testing.T) {
	tests := []struct {
		name    string
		conf    core.Config
		transf  core.Transfer
		wantErr bool
	}{
		{
			name:    "Valid input",
			conf:    mockConfig,
			transf:  mockTransfer,
			wantErr: false,
		},
		{
			name: "Invalid amount format",
			conf: mockConfig,
			transf: core.Transfer{
				Amount:                "invalid",
				Purpose:               core.PURPOSE_SALA,
				RemittanceInformation: "Payment for invoice 12345",
			},
			wantErr: true,
		},
		{
			name: "Amount out of range",
			conf: mockConfig,
			transf: core.Transfer{
				Amount:                "1000000000.00",
				Purpose:               core.PURPOSE_SALA,
				RemittanceInformation: "Payment for invoice 12345",
			},
			wantErr: true,
		},
		{
			name: "Payload too long",
			conf: mockConfig,
			transf: core.Transfer{
				Amount:                "123.45",
				Purpose:               core.PURPOSE_SALA,
				RemittanceInformation: string(make([]byte, core.MAX_PAYLOAD_SIZE+1)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := core.GetRawText(tt.conf, tt.transf)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
