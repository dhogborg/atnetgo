# atnetgo

atnetgo is a implementation of a client using the [netatmo-api-go](https://github.com/exzz/netatmo-api-go) lib.

It's function is to provide readout of netatmo station values. Using the command line options the output can be filterd to provide only the values you are interested in. 

Using `grep` and `awk` a single value can be extracted from the output.

## Running
Binary release is available from the Releases page: https://github.com/dhogborg/atnetgo/releases/

Currently binary release is provided for the following platforms:
* OS X (x64)
* Linux (x64)
* Linux (arm5)
* Linux (arm7)
* Windows (x64)
* Windows (386)

Using the binary you don't need to provide your own Netatmo AppID and Client secret.

## Setup 
1. Create a Netatmo app ID: https://dev.netatmo.com/dev/createapp, put the ID and the Client Secret in the secrets.go file (follow instructions in secrets.example.go).
2. Use your netatmo credentials, perferably as environment variables to avoid storing passwords in your .bash_history
3. Run atnetgo with the `list` command to see what's on your account.
4. Use `--station` and `--module` to specify a single module to print

## Examples
#### Print all stations and all modules
```
$ atnetgo pretty
Station: Office
	Outdoor:
		Humidity: 81
		Temperature: 10.40
	Indoor:
		Pressure: 992.40
		Temperature: 21.40
		CO2: 435
		Humidity: 39
		Noise: 40
Station: Home
	Rain:
		Rain: 0.00
	Outdoor:
		Temperature: 9.00
		Humidity: 83
	Indoor:
		Temperature: 22.00
		CO2: 1057
		Humidity: 49
		Noise: 39
		Pressure: 993.80
```

```
$ atnetgo list
Office: Outdoor: Temperature: 10.40
Office: Outdoor: Humidity: 81
Office: Indoor: Temperature: 21.40
Office: Indoor: CO2: 435
Office: Indoor: Humidity: 39
Office: Indoor: Noise: 40
Office: Indoor: Pressure: 992.40
Home: Rain: Rain: 0.00
Home: Outdoor: Temperature: 9.00
Home: Outdoor: Humidity: 83
Home: Indoor: Temperature: 22.00
Home: Indoor: CO2: 1057
Home: Indoor: Humidity: 49
Home: Indoor: Noise: 39
Home: Indoor: Pressure: 993.80
```

```
$ atnetgo json
{
	"Office": {
		"Indoor": {
			"CO2": "431",
			"Humidity": "39",
			"Noise": "37",
			"Pressure": "992.40",
			"Temperature": "21.40"
		},
		"Outdoor": {
			"Humidity": "81",
			"Temperature": "10.50"
		}
	},
	"Home": {
		"Indoor": {
			"CO2": "1171",
			"Humidity": "49",
			"Noise": "37",
			"Pressure": "994.00",
			"Temperature": "22.00"
		},
		"Outdoor": {
			"Humidity": "83",
			"Temperature": "9.00"
		},
		"Rain": {
			"Rain": "0.00"
		}
	}
}
```
The example above is prettyfied for clarity. Actual output is plain json without newlines or tabs.

#### Extracting a single value
```
$ atnetgo list | grep 'Temperature' | awk '{print $4}'
```
```
21.60
```


## Options
```
NAME:
   atnetgo - Read values from the Netatmo API and write to stdout

USAGE:
   atnetgo [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
  github.com/dhogborg - <d@hogborg.se>

COMMANDS:
   pretty	Pretty print the stations and the modules attached
   list		List the modules and the values in a greppable list
   json		Output a machine readable json string
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --user, -u 		Netatmo login name [$NETATMO_USER]
   --password, -p 	Netatmo password [$NETATMO_PASSWORD]
   --station, -s 	The station name, default to the first one [$NETATMO_STATION]
   --module, -m 	Station module name, default to the first one [$NETATMO_MODULE]
   --help, -h		show help
   --version, -v	print the version
   
   ```