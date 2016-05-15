package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	netatmo "github.com/exzz/netatmo-api-go"
)

// filterd device collection
type DeviceCollection struct {
	NetatmoStations []*netatmo.Device
	Modules         []*netatmo.Device
}

func (d *DeviceCollection) Stations() []*netatmo.Device { return d.NetatmoStations }

func main() {
	app := cli.NewApp()
	app.Name = "atnetgo"
	app.Usage = "Read values from the Netatmo API and write to stdout"
	app.Version = "0.0.3"
	app.Author = "github.com/dhogborg"
	app.Email = "d@hogborg.se"

	app.Action = func(c *cli.Context) {
		d := getDevices(c)
		listPrint(d)
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "pretty",
			Usage: "Pretty print the stations and the modules attached",
			Action: func(c *cli.Context) {
				d := getDevices(c)
				prettyPrint(d)
			},
		},
		cli.Command{
			Name:  "list",
			Usage: "List the modules and the values in a greppable list",
			Action: func(c *cli.Context) {
				d := getDevices(c)
				listPrint(d)
			},
		},
		cli.Command{
			Name:  "json",
			Usage: "Output a machine readable json string",
			Action: func(c *cli.Context) {
				d := getDevices(c)
				jsonPrint(d)
			},
		},
		cli.Command{
			Name:  "influx",
			Usage: "Output InfluxDB line format",
			Action: func(c *cli.Context) {
				d := getDevices(c)
				linePrint(d)
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user,u",
			Usage:  "Netatmo login name",
			EnvVar: "NETATMO_USER",
		},
		cli.StringFlag{
			Name:   "password,p",
			Usage:  "Netatmo password",
			EnvVar: "NETATMO_PASSWORD",
		},
		cli.StringFlag{
			Name:   "station,s",
			Usage:  "A station filter, default to none (print everything)",
			EnvVar: "NETATMO_STATION",
		},
	}

	app.Run(os.Args)
}

func getDevices(ctx *cli.Context) *DeviceCollection {

	config := netatmo.Config{
		ClientID:     NetatmoAppID,
		ClientSecret: NetatmoAppSecret,
		Username:     ctx.GlobalString("user"),
		Password:     ctx.GlobalString("password"),
	}

	n, err := netatmo.NewClient(config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("client error")
		os.Exit(1)
	}

	dc, err := n.GetDeviceCollection()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("unable to fetch device collection")
		os.Exit(1)
	}

	collection := filterDevices(ctx, dc)

	return collection

}

func filterDevices(ctx *cli.Context, dc *netatmo.DeviceCollection) *DeviceCollection {
	collection := &DeviceCollection{
		NetatmoStations: dc.Body.Stations,
		Modules:         dc.Body.Modules,
	}

	if sfilter := ctx.GlobalString("station"); sfilter != "" {
		stations := []*netatmo.Device{}
		for _, station := range collection.Stations() {
			if matchesFilter(station, sfilter) {
				stations = append(stations, station)
			}
		}
		collection.NetatmoStations = stations
	}

	return collection
}

func matchesFilter(device *netatmo.Device, filter string) bool {
	return strings.Index(device.StationName, filter) > -1
}

func listPrint(devices *DeviceCollection) {
	for _, station := range devices.Stations() {
		for _, module := range station.Modules() {
			_, data := module.Data()
			for dataType, value := range data {
				fmt.Printf("%s: %s: %s: %s\n", station.StationName, module.ModuleName, dataType, valueString(value))
			}
		}
	}
}

func jsonPrint(devices *DeviceCollection) {

	block := map[string]interface{}{}

	for _, station := range devices.Stations() {
		sblock := map[string]interface{}{}
		for _, module := range station.Modules() {
			mblock := map[string]string{}
			_, data := module.Data()
			for dataType, value := range data {
				mblock[dataType] = valueString(value)
			}
			sblock[module.ModuleName] = mblock
		}
		block[station.StationName] = sblock
	}

	b, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func prettyPrint(devices *DeviceCollection) {
	for _, station := range devices.Stations() {
		fmt.Printf("Station: %s\n", station.StationName)
		for _, module := range station.Modules() {
			fmt.Printf("\t%s:\n", module.ModuleName)
			_, data := module.Data()
			for dataType, value := range data {
				fmt.Printf("\t\t%s: %s\n", dataType, valueString(value))
			}
		}
	}
}

func linePrint(devices *DeviceCollection) {

	// some values should be represented as integers
	typeSuffix := func(t string) string {
		typemap := map[string]string{
			"co2":      "i",
			"humidity": "i",
			"noise":    "i",
		}
		if s, ok := typemap[strings.ToLower(t)]; ok {
			return s
		}
		return ""
	}

	tags := make([]string, 2)

	for _, station := range devices.Stations() {
		tags[0] = "station=" + strings.ToLower(station.StationName)

		for _, module := range station.Modules() {
			tags[1] = "module=" + strings.ToLower(module.ModuleName)

			_, data := module.Data()
			for dataType, value := range data {
				tagstr := strings.Join(tags, ",")
				tagstr = strings.Replace(tagstr, " ", "_", -1)
				fmt.Printf("%s,%s value=%s%s\n", strings.ToLower(dataType), tagstr, valueString(value), typeSuffix(dataType))
			}
		}
	}
}

func valueString(value interface{}) string {
	switch value := value.(type) {
	case float32, float64:
		return fmt.Sprintf("%0.2f", value)
	case string:
		return value
	case int, int32, int64:
		return fmt.Sprintf("%d", value)
	default:
		return "-"
	}
}
