# atnetgo

atnetgo is a implementation of a client using the [netatmo-api-go](https://github.com/exzz/netatmo-api-go) lib.

It's function is to provide readout of netatmo station values. Using the command line options the output can be filterd to provide only the values you are interested in. 

Using `grep` and `awk` a single value can be extracted from the output.

## Setup 
1. Create a Netatmo app ID: https://dev.netatmo.com/dev/createapp, use the AppID and AppSecret, perferably as environment variables (see **Options**).
2. Use your netatmo credentials, perferably as environment variables to avoid storing passwords in your .bash_history
3. Run atnetgo with the `list` command to see what's on your account.
4. Use `--station` and `--module` to specify a single module to print

## Examples
#### Print all stations and all modules
```
$ atnetgo list
Station: Office
    Module: Outdoor
        Temperature: 9.40
        Humidity: 89
    Module: Indoor
        Humidity: 41
        Noise: 41
        Pressure: 998.50
        Temperature: 22.90
        CO2: 874
Station: Home
    Module: Rain
        Rain: 0.00
    Module: Outdoor
        Temperature: 8.60
        Humidity: 91
    Module: Indoor
        Pressure: 1000.10
        Temperature: 21.60
        CO2: 697
        Humidity: 48
        Noise: 49
```

#### Printing a single module
```
$ atnetgo --station 'Home' --module 'Indoor'
```
```
Pressure: 1000.10
Temperature: 21.60
CO2: 697
Humidity: 48
Noise: 49
```

#### Extracting a single value
```
$ atnetgo --station 'Home' --module 'Indoor' | grep 'Temperature' | awk '{print $2}'
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
   0.0.1

AUTHOR:
  github.com/dhogborg - <d@hogborg.se>

COMMANDS:
   list     List the stations and the modules attached
   fetch    Fetch and display data
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --appid          Netatmo application ID [$NETATMO_APP_ID]
   --appsecret      Netatmo application Secret [$NETATMO_APP_SECRET]
   --user, -u       Netatmo login name [$NETATMO_USER]
   --password, -p   Netatmo password [$NETATMO_PASSWORD]
   --station, -s    The station name, default to the first one [$NETATMO_STATION]
   --module, -m     Station module name, default to the first one [$NETATMO_MODULE]
   --help, -h       show help
   --version, -v    print the version
   ```