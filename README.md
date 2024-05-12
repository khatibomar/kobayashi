# kobayashi
<img src=".github/assets/logo.gif" style="width:250px;height:250px;"/>

kobayashi is Go library to get streaming sites direct links easy as doing : <br> <br>

```go
package main

import (
	"fmt"
	"log"

	"github.com/khatibomar/kobayashi"
)

func main () {
  d := kobayashi.NewDecoder()
  directLink , err := d.Decode("https://mixdrop.co/e/zp03jllqawql1q")
  if err != nil {
    // Handle Error...
  }
  fmt.Println(directLink)
  // Result :
  // https://a-delivery15.mxdcontent.net/v/f96b4514c68b0a69ae38c9726488b35e.mp4?s=pWZlKfQUbPmowXdPl_FGAQ&e=1641662900&_t=1641650029

  directLink , err = d.Decode("https://www.fembed.com/api/source/jl4n4fd42rzy3e5")
  if err != nil {
    // Handle Error...
  }
  fmt.Println(directLink)

  // Result :
  // https://fvs.io/redirector?token=bkZHY3IzTVgyMGR2NHRNWGF2NFZ2cHE0YVhIMXJmV3pwTVNKc1d2aVFmdlRGQnhUZUYxMkRRaER0TXdiQTdIdkRMdjFCT2FuejBRd3M4c0d5dUxhVWQzeUw3RnAxdnBzcTNTcWRnMUhBRVpRS21IKzlmVjZDSklCR1IrekZtQ09oMDlscUc4ZS9YcTRGd1VsTUk4cHNLcTEvSm1jTXlEZEZBPT06Zko2allWUUQ5cVlKb01SbThhU3NuUT09yvQ5
}
```

if we pass unsupported link 

```go
directLink , err = d.Decode("https://uptostream.com/tvt5mjn49529")

// Result : host is not supported, yet...
```

if your host no supported yet , feel free to open an issue and let me know 🧑‍💻
