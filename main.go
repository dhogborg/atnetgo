package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	netatmo "github.com/exzz/netatmo-api-go"
)

func main() {
	app := cli.NewApp()
	app.Name = "atnetgo"
	app.Usage = "Read values from the Netatmo API and write to stdout"
	app.Version = "0.0.1"
	app.Author = "github.com/dhogborg"
	app.Email = "d@hogborg.se"

	app.Action = func(c *cli.Context) {

		fetchData(c)
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "list",
			Usage: "List the stations and the modules attached",
			Action: func(c *cli.Context) {

				list(c)
			},
		},
		cli.Command{
			Name:  "fetch",
			Usage: "Fetch and display data",
			Action: func(c *cli.Context) {
				fetchData(c)
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "appid",
			Usage:  "Netatmo application ID",
			EnvVar: "NETATMO_APP_ID",
		},
		cli.StringFlag{
			Name:   "appsecret",
			Usage:  "Netatmo application Secret",
			EnvVar: "NETATMO_APP_SECRET",
		},
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
			Usage:  "The station name, default to the first one",
			EnvVar: "NETATMO_STATION",
		},
		cli.StringFlag{
			Name:   "module,m",
			Usage:  "Station module name, default to the first one",
			EnvVar: "NETATMO_MODULE",
		},
	}

	app.Run(os.Args)
}

func list(ctx *cli.Context) {

	dc := getDeviceCollection(netatmo.Config{
		ClientID:     ctx.GlobalString("appid"),
		ClientSecret: ctx.GlobalString("appsecret"),
		Username:     ctx.GlobalString("user"),
		Password:     ctx.GlobalString("password"),
	})

	for _, station := range dc.Stations() {

		fmt.Printf("Station: %s\n", station.StationName)

		for _, module := range station.Modules() {

			fmt.Printf("\tModule: %s\n", module.ModuleName)
			printModuleData(module, "\t\t")
		}
	}
}

func fetchData(ctx *cli.Context) {

	dc := getDeviceCollection(netatmo.Config{
		ClientID:     ctx.GlobalString("appid"),
		ClientSecret: ctx.GlobalString("appsecret"),
		Username:     ctx.GlobalString("user"),
		Password:     ctx.GlobalString("password"),
	})

	var station *netatmo.Device
	var module *netatmo.Device

	for _, s := range dc.Stations() {

		if station == nil || s.StationName == ctx.GlobalString("station") {
			station = s
		}

		for _, m := range s.Modules() {

			if module == nil || m.ModuleName == ctx.GlobalString("module") {
				module = m
			}

		}
	}

	fmt.Printf("%s: %s\n", station.StationName, module.ModuleName)

	printModuleData(module, "")
}

func getDeviceCollection(config netatmo.Config) *netatmo.DeviceCollection {

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

	return dc
}

func printModuleData(module *netatmo.Device, prefix string) {

	_, data := module.Data()
	for dataType, value := range data {
		output := ""
		switch value := value.(type) {
		case float32, float64:
			output = fmt.Sprintf("%0.2f", value)
		case string:
			output = value
		case int, int32, int64:
			output = fmt.Sprintf("%d", value)
		default:
			output = "–"
		}
		fmt.Printf("%s%s: %s\n", prefix, dataType, output)
	}
}
