# nfcsdk

NFC SDK to read and write NTAG213/215/216 NFC tags 

[![Coverage Status](https://coveralls.io/repos/github/happy-sdk/nfcsdk/badge.svg?branch=main)](https://coveralls.io/github/happy-sdk/nfcsdk?branch=main)

## Example

```go
package main

import (
  "log/slog"
  "os"

  "github.com/happy-sdk/nfcsdk"
)

func main() {
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
    ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
      if a.Key == slog.TimeKey && len(groups) == 0 {
        return slog.Attr{}
      }
      return a
    },
  }))

  // ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
  sdk, err := nfcsdk.New(nil, logger)
  if err != nil {
    os.Exit(1)
  }

  // Default behavior: Automatically selects the first available reader 
  // if no SelectReader callback is set.
  sdk.SelectReader(func(rs []nfcsdk.Reader) ([]nfcsdk.Reader, error) {
    logger.Info("select first reader")
    rs[0].Use = true
    return rs, nil
  })

  // Reads card uid when card is present
  sdk.OnCardPresent(func(card nfcsdk.Card) error {
    logger.Info("--- HANDLE CARD PRESENT ---")
    cmdGetUID := ntag.NewGetUIDCmd()

    response, err := card.Transmit(cmdGetUID)
    if err != nil {
      return err
    }

    logger.Info(cmdGetUID.String(), slog.String("uid", tag.UID), uidr.LogAttr())
    logger.Info("--- CARD HANDLE COMPLETED ---")
    return nil
  })

  if err := sdk.Run(); err != nil {
    return
  }
}
// example out:
// level=DEBUG msg="nfc: scard context established"
// level=DEBUG msg="nfc: found reader" nfc.reader.id=1 nfc.reader.name="<your-reader>"
// level=INFO msg="select first reader"
// level=INFO msg="nfc: started" nfc.time="<current-time>"
// level=DEBUG msg="nfc: no card present, waiting..."
// level=DEBUG msg="nfc: card is present"
// level=DEBUG msg="nfc: card connected" nfc.protocols="T1, Any"
// level=INFO msg="nfc: connected card status" nfc.state="Present, Powered, Negotiable" nfc.protocol="T1, Any" nfc.reader="<your-reader>" nfc.atr=<reader-atr>
// level=DEBUG msg="--- HANDLE CARD PRESENT ---"
// level=DEBUG msg="GET_UID [FF:CA:00:00:00]" uid=01:02:03:04:05:06:07 status=success sw1=90 sw2=00 payload.len=7
// level=DEBUG msg="--- CARD HANDLE COMPLETED ---"
// level=INFO msg="nfc: card disconnected"
// level=DEBUG msg="nfc: no card present, waiting..."

```

## Development

To set up the development environment:

```bash
git clone git@github.com:happy-sdk/nfcsdk.git
cd nfcsdk
```

To generate code use `go generate .` in package root