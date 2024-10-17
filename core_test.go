package goepcqrcode_test

import (
	"testing"

	sepa "github.com/bergzeit/goepcqrcode"
)

// Mock data for testing
var (
	mockConfig   = sepa.NewConfig(sepa.VERSION_1, sepa.CHAR_SET_UTF8, sepa.IDENTIFICATION_CODE_SCT, "YOURBIC", "John Doe", "DE89370400440532013000", "EUR")
	mockTransfer = sepa.NewTransfer("123.45", sepa.PURPOSE_SALA, "Payment for invoice 12345", "", "")
)

func TestGetRawText(t *testing.T) {
	tests := []struct {
		name    string
		conf    sepa.Config
		transf  sepa.Transfer
		wantErr bool
	}{
		{
			name:    "Valid input",
			conf:    mockConfig,
			transf:  mockTransfer,
			wantErr: false,
		},
		{
			name:    "Invalid amount format",
			conf:    mockConfig,
			transf:  sepa.NewTransfer("invalid", sepa.PURPOSE_SALA, "Payment for invoice 12345", "", ""),
			wantErr: true,
		},
		{
			name:    "Amount out of range",
			conf:    mockConfig,
			transf:  sepa.NewTransfer("1_000_000_000.00", sepa.PURPOSE_GDDS, "Payment for invoice 12345", "", ""),
			wantErr: true,
		},
		{
			name:    "Payload too long",
			conf:    mockConfig,
			transf:  sepa.NewTransfer("123.45", sepa.PURPOSE_GDDS, string(make([]byte, sepa.MAX_PAYLOAD_SIZE+1)), "", ""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sepa.GetRawText(tt.conf, tt.transf)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
