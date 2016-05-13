# netatmo-api-go
Simple API to access Netatmo weather station data written in Go.

Currently tested only with one weather station, outdoor and indoor modules. Let me know if it works with rain or wind gaude.

## Quickstart

- [Create a new netatmo app](https://dev.netatmo.com/dev/createapp)
- Download module ```go get github.com/exzz/netatmo-api-go```
- Try below example (do not forget to edit auth credentials)

## Example

```go
import (
  "fmt"
  "os"

  "github.com/exzz/netatmo-api-go"
)

func main() {

  n, err := netatmo.NewClient(netatmo.Config{
    ClientID:     "YOUR_APP_ID",
    ClientSecret: "YOUR_APP_SECRET",
    Username:     "YOUR_CREDENTIAL",
    Password:     "YOUR_PASSWORD",
  })
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  dc, err := n.GetDeviceCollection()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, station := range dc.Stations() {
    fmt.Printf("Station : %s\n", station.StationName)

    for _, module := range station.Modules() {
      fmt.Printf("\tModule : %s\n", module.ModuleName)

      ts, data := module.Data()
      for dataType, value := range data {
        fmt.Printf("\t\t%s : %s (%d)\n", dataType, value, ts)
      }
    }
  }
}
```

Output should look like this :
```
Station : Home
  Module : Outside
    Temperature : %!s(float32=20.2) (1440302379)
    Humidity : %!s(int32=86) (1440302379)
  Module : Bedroom 1
    CO2 : %!s(int32=500) (1441981664)
    Humidity : %!s(int32=69) (1441981664)
    Temperature : %!s(float32=21.2) (1441981664)
  Module : Bedroom 2
    Temperature : %!s(float32=21) (1441981632)
    CO2 : %!s(int32=508) (1441981632)
    Humidity : %!s(int32=68) (1441981632)
  Module : Living room
    Temperature : %!s(float32=22.1) (1441981633)
    CO2 : %!s(int32=516) (1441981633)
    Humidity : %!s(int32=67) (1441981633)
  Module : Dining room
    Humidity : %!s(int32=75) (1441982895)
    Noise : %!s(int32=36) (1441982895)
    Pressure : %!s(float32=1015.9) (1441982895)
    Temperature : %!s(float32=21.5) (1441982895)
    CO2 : %!s(int32=582) (1441982895)

```
## Tips
- Only GetDeviceCollection() method actually do an API call and refresh all data at once
- Main station is handle as a module, it means that Modules() method returns list of additional modules and station itself.
