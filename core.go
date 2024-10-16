// this pkg renders sepa information as a qr code
// many features are missing, please use with caution and only if you know
// what you are doing

package goepcqrcode

import (
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/bergzeit/goepcqrcode/pkg/asserts"
	"github.com/skip2/go-qrcode"
)

// https://de.wikipedia.org/wiki/EPC-QR-Code

// TODO: add selectable element separator LF, CRLF
// TODO: add version1 or version2 selection
// TODO: add a validator for the different fields

const (
	// RAW_TEMPLATE represents the txt raw instance of the information. Do not reformat this string,
	// it is used as a template and the formatting is intentional
	RAW_TEMPLATE = `{{ .ServiceTag }}
{{ .Version }}
{{ .CharacterSet }}
{{ .IdentificationCode }}
{{ .BIC }}
{{ .Name }}
{{ .IBAN }}
{{ .Amount }}
{{ .Purpose }}
{{ .RemittanceReference }}
{{ .RemittanceText }}
{{ .Information}}`

	// Allowed Service tags
	VERSION_1 string = "001" // Version 001
	VERSION_2 string = "002" // Version 002

	// Allowed character sets
	CHAR_SET_UTF8        string = "1" // UTF-8
	CHAR_SET_ISO_8859_1  string = "2" // Latin-1, Westeuropäisch
	CHAR_SET_ISO_8859_2  string = "3" // Latin-2, Mitteleuropäisch
	CHAR_SET_ISO_8859_4  string = "4" // Latin-4, Nordeuropäisch
	CHAR_SET_ISO_8859_5  string = "5" // Kyrillisch
	CHAR_SET_ISO_8859_7  string = "6" // Griechisch
	CHAR_SET_ISO_8859_10 string = "7" // Latin-6, Nordisch
	CHAR_SET_ISO_8859_15 string = "8" // Latin-9, Westeuropäisch

	// Allowed identification codes
	IDENTIFICATION_CODE_SCT  string = "SCT"  // SEPA Credit Transfer
	IDENTIFICATION_CODE_INST string = "INST" // SEPA Instant Credit Transfer

	// Allowed purposes
	PURPOSE_BENE string = "BENE" // Beneficiary payment
	PURPOSE_DEPT string = "DEPT" // Departmental payment
	PURPOSE_GDDS string = "GDDS" // General purpose
	PURPOSE_MTUP string = "MTUP" // Mobile top-up
	PURPOSE_PENS string = "PENS" // Pension payment
	PURPOSE_SALA string = "SALA" // Salary payment
	PURPOSE_TRAD string = "TRAD" // Trade payment

	// Allowed sizes
	MAX_PAYLOAD_SIZE = 331 // the maximum allowed payload size for a QR code is 331 bytes
)

// Errors
var ErrSepaLength error = errors.New("sepa transfer information too long")

// GetRawText generates the raw text from the config and the transfer given
func GetRawText(conf Config, transf Transfer) (string, error) {
	raw, err := newRaw(conf, transf)
	if err != nil {
		return "", fmt.Errorf("error creating raw text: %w", err)
	}

	r, err := raw.render()
	if err != nil {
		return "", fmt.Errorf("error rendering raw text: %w", err)
	}

	return r, nil
}

// QRCode generates a QR code from the config and the transfer given
func GetQRCode(conf Config, transf Transfer) ([]byte, error) {
	r, err := GetRawText(conf, transf)
	if err != nil {
		return []byte{}, fmt.Errorf("error rendering qr code: %w", err)
	}

	var b []byte
	b, err = qrcode.Encode(r, qrcode.Medium, 256)
	if err != nil {
		return []byte{}, fmt.Errorf("error rendering qr code: %w", err)
	}

	asserts.AssertSize(b, 331)

	return b, nil
}

// Transfer represents a SEPA transfer information for a specific transfer
type Transfer struct {
	amount              string
	purpose             string
	remittanceReference string
	remittanceText      string
	information         string
}

func NewTransfer(amount string, purpose string, remittanceReference string, remittanceText string, information string) Transfer {
	return Transfer{
		amount:              amount,
		purpose:             purpose,
		remittanceReference: remittanceReference,
		remittanceText:      remittanceText,
		information:         information,
	}
}

func (t *Transfer) validate() error {
	amount, err := strconv.ParseFloat(t.amount, 64)
	if err != nil {
		return fmt.Errorf("error invalid amount format: %w", err)
	}

	if !asserts.AssertInBetween[float64](0.01, 999_999_999.99, amount) {
		return fmt.Errorf("error amount %v out of valid range", t.amount)
	}

	return nil
}

// raw represents the raw information for the SEPA transfer which will be rendered to text
type raw struct {
	serviceTag          string // required, options: ["BCD"], Service Tag
	version             string // required, options: ["001", "002"], Version
	characterSet        string // required, options: [1, 2, 3, 4, 5, 6, 7, 8, 9], Character set
	identificationCode  string // required, options: ["SCT", "INST"], Identification code
	bic                 string // required with Version 001, optional with Version 002, BIC of the receiver
	name                string // required, name of the receiver
	iban                string // required, IBAN of the receiver
	amount              string // optional, format: "EUR#.##",amount to be paid
	purpose             string // optional, options: ["BENE", "DEPT", "GDDS", "MTUP", "PENS", "SALA", "TRAD"], format ####, reason for payment
	remittanceReference string // optional, text, 25 char, this or remittanceText, Referenz
	remittanceText      string // optional, text, 140char, this or remittanceReference, Verwendungszweck
	information         string // optional, text, 70 char, additional information
}

// newRaw creates a new raw instance of the SEPA information by combinding a config and transfer information
func newRaw(conf Config, transf Transfer) (raw, error) {
	err := transf.validate()
	if err != nil {
		return raw{}, fmt.Errorf("error validating transfer information: %w", err)
	}

	return raw{
		serviceTag:          "BCD",
		version:             conf.version,
		characterSet:        conf.characterSet,
		identificationCode:  conf.identificationCode,
		bic:                 conf.bic,
		name:                conf.name,
		iban:                conf.iban,
		amount:              fmt.Sprintf("%s%s", conf.currency, transf.amount),
		purpose:             transf.purpose,
		remittanceReference: transf.remittanceReference,
		remittanceText:      transf.remittanceText,
		information:         transf.information,
	}, nil
}

func (r *raw) render() (string, error) {
	t, err := template.New("raw").Parse(RAW_TEMPLATE)
	if err != nil {
		return "", fmt.Errorf("error parsing sepa information template: %w", err)
	}

	var sb strings.Builder
	err = t.Execute(&sb, r)
	if err != nil {
		return "", fmt.Errorf("error rendering sepa information to template: %w", err)
	}

	if sb.Len() > MAX_PAYLOAD_SIZE {
		return "", ErrSepaLength
	}

	return sb.String(), nil
}

// Config represents the configuration for the SEPA information, this is the
// static part that mostly does not change between different transfers
type Config struct {
	version            string
	characterSet       string
	identificationCode string
	bic                string
	name               string
	iban               string
	currency           string
}

// NewConfig creates a new configuration for the SEPA information. This
// contains all the static information that often does not change between different transfers
func NewConfig(version string, characterSet string, identificationCode string, bic string, name string, iban string, currency string) Config {
	return Config{
		version:            version,
		characterSet:       characterSet,
		identificationCode: identificationCode,
		bic:                bic,
		name:               name,
		iban:               iban,
		currency:           currency,
	}
}
