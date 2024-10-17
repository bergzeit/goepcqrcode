package goepcqrcode_test

import (
	"fmt"
	"log"
	"os"

	sepa "github.com/bergzeit/goepcqrcode"
)

func ExampleGetRawText() {
	conf := sepa.NewConfig(sepa.VERSION_1, sepa.CHAR_SET_UTF8, sepa.IDENTIFICATION_CODE_SCT, "YOURBIC", "YOURNAME", "YOURIBAN", "EUR")

	transf := sepa.NewTransfer("123.45", sepa.PURPOSE_GDDS, "Payment for invoice", "", "")

	rawText, err := sepa.GetRawText(conf, transf)
	if err != nil {
		log.Fatalf("error generating raw text: %v", err)
	}

	// now you got a raw text version of the protocol, not a qr code
	// this is useful e.g. to send over wire and later render a qr code from it etc.
	fmt.Print(rawText)
	// Output:
	// BCD
	// 001
	// 1
	// SCT
	// YOURBIC
	// YOURNAME
	// YOURIBAN
	// EUR123.45
	// GDDS
	// Payment for invoice
}

func ExampleGetQRCode() {
	conf := sepa.NewConfig(sepa.VERSION_1, sepa.CHAR_SET_UTF8, sepa.IDENTIFICATION_CODE_SCT, "YOURBIC", "YOURNAME", "YOURIBAN", "EUR")

	transf := sepa.NewTransfer("123.45", sepa.PURPOSE_GDDS, "Payment for invoice", "", "")

	qrCode, err := sepa.GetQRCode(conf, transf)
	if err != nil {
		log.Fatalf("error generating QR code: %v", err)
	}

	file, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(qrCode)
	if err != nil {
		log.Fatalf("error writing QR code to file: %v", err)
	}

	fmt.Println("QR code saved to qrcode.png")
	// Output:
	// QR code saved to qrcode.png
}
