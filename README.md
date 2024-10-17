# Go EPC QR code

This is a small go library that supports you creating EPC standard compliant SEPA transfer QR codes. 

## Compatibility

Currently the pkg is compliant to EPC069-12 version 3.0 from 13th September 2022. [More Information](https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2022-09/EPC069-12%20v3.0%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20an%20SCT_0.pdf)

## Installation

To install the package, use the following command:

```sh
go get github.com/bergzeit/goepcqrcode
```

## Usage

The library supports you generating QR codes or constructing plain text versions of EPC protocol sepa transfers.

### GetRawText

The GetRawText function generates the raw text from the given configuration and transfer details.

``` go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/epc-qrcode/core"
)

func main() {
    conf := core.Config{
        Version:           "001",
        CharacterSet:      "UTF-8",
        IdentificationCode: "SCT",
        BIC:               "YOURBIC",
        Name:              "John Doe",
        IBAN:              "DE89370400440532013000",
        Currency:          "EUR",
    }

    transf := core.Transfer{
        Amount:   "123.45",
        Purpose:  "Invoice 12345",
        RemittanceInformation: "Payment for invoice 12345",
    }

    rawText, err := core.GetRawText(conf, transf)
    if err != nil {
        log.Fatalf("Error generating raw text: %v", err)
    }

    // now you got a raw text version of the protocol, not a qr code
    // this is useful e.g. to send over wire and later render a qr code from it etc.
    fmt.Println("Raw Text:", rawText)
}
```

### GetQRCode

The GetQRCode function generates a QR code image from the given configuration and transfer details. You will receive a slice of bytes which
you can then render (store, send to browser as base64 etc), using your favorite method.

``` go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/yourusername/epc-qrcode/core"
)

func main() {
    conf := core.Config{
        Version:           "001",
        CharacterSet:      "UTF-8",
        IdentificationCode: "SCT",
        BIC:               "YOURBIC",
        Name:              "John Doe",
        IBAN:              "DE89370400440532013000",
        Currency:          "EUR",
    }

    transf := core.Transfer{
        Amount:   "123.45",
        Purpose:  "Invoice 12345",
        RemittanceInformation: "Payment for invoice 12345",
    }

    qrCode, err := core.GetQRCode(conf, transf)
    if err != nil {
        log.Fatalf("Error generating QR code: %v", err)
    }

    file, err := os.Create("qrcode.png")
    if err != nil {
        log.Fatalf("Error creating file: %v", err)
    }
    defer file.Close()

    _, err = file.Write(qrCode)
    if err != nil {
        log.Fatalf("Error writing QR code to file: %v", err)
    }

    fmt.Println("QR code saved to qrcode.png")
```

## Issues

When encountering issues please open a github issue describing the issue you have encountered and provide as much information as you can.

## Contribution

If you want to contribute to the project please create a fork and open a pull request to the project.
