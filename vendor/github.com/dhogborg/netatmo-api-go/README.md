# netatmo-api-go
Simple API to access Netatmo weather station data written in Go.

Currently tested only with one weather station, outdoor and indoor modules and rain gauge. Let me know if it works with wind gauge.

## Quickstart

- [Create a new netatmo app](https://dev.netatmo.com/dev/createapp)
- Download module ```go get github.com/exzz/netatmo-api-go```
- Try below example (do not forget to edit auth credentials)

## Example

```
go
package main

import (
	"fmt"
	"os"

	netatmo "github.com/exzz/netatmo-api-go"
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

	dc, err := n.Read()
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
        Module : Chambre Enfant
                Temperature : %!s(float32=18.4) (1479127223)
                CO2 : %!s(int32=567) (1479127223)
                Humidity : %!s(int32=65) (1479127223)
        Module : Chambre
                Temperature : %!s(float32=18.1) (1479127230)
                CO2 : %!s(int32=494) (1479127230)
                Humidity : %!s(int32=65) (1479127230)
        Module : Salon
                Temperature : %!s(float32=18.3) (1479127217)
                CO2 : %!s(int32=434) (1479127217)
                Humidity : %!s(int32=63) (1479127217)
        Module : Exterieur
                Temperature : %!s(float32=11.9) (1479127243)
                Humidity : %!s(int32=81) (1479127243)
        Module : Pluie
                Rain : %!s(float32=0) (1479127249)
        Module : Salle à manger
                Temperature : %!s(float32=17.8) (1479127255)
                CO2 : %!s(int32=473) (1479127255)
                Humidity : %!s(int32=68) (1479127255)
                Noise : %!s(int32=36) (1479127255)
                Pressure : %!s(float32=1033.3) (1479127255)

```
## Tips
- Only Read() method actually do an API call and refresh all data at once
- Main station is handle as a module, it means that Modules() method returns list of additional modules and station itself.
