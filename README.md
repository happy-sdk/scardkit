# nfcsdk

NFC SDK to read and write NTAG213/215/216 NFC tags 

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
    cmd, err := ntag.Cmd(ntag.CmdGetUID)
    if err != nil {
      return err
    }

    response, err := cmd.Transmit(card)
    if err != nil {
      return err
    }

    logger.Info(cmd.String(),
      slog.String("uid", nfcsdk.FormatByteSlice(response.Payload())),
      response.LogAttr(),
    )
    logger.Info("--- CARD HANDLE COMPLETED ---")
    return nil
  })

  if err := sdk.Run(); err != nil {
    return
  }
}

```

## Development

To set up the development environment:

```bash
git clone git@github.com:happy-sdk/nfcsdk.git
cd nfcsdk
```

To generate code use `go generate .` in package root