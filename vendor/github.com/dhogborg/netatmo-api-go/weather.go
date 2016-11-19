package netatmo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"golang.org/x/oauth2"
)

const (
	// DefaultBaseURL is netatmo api url
	baseURL = "https://api.netatmo.net/"
	// DefaultAuthURL is netatmo auth url
	authURL = baseURL + "oauth2/token"
	// DefaultDeviceURL is netatmo device url
	deviceURL = baseURL + "/api/getstationsdata"
)

// Config is used to specify credential to Netatmo API
// ClientID : Client ID from netatmo app registration at http://dev.netatmo.com/dev/listapps
// ClientSecret : Client app secret
// Username : Your netatmo account username
// Password : Your netatmo account password
type Config struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

// Client use to make request to Netatmo API
type Client struct {
	oauth        *oauth2.Config
	httpClient   *http.Client
	httpResponse *http.Response
	Dc           *DeviceCollection
}

// DeviceCollection hold all devices from netatmo account
type DeviceCollection struct {
	Body struct {
		Devices []*Device `json:"devices"`
	}
}

// Device is a station or a module
// ID : Mac address
// StationName : Station name (only for station)
// ModuleName : Module name
// Type : Module type :
//  "NAMain" : for the base station
//  "NAModule1" : for the outdoor module
//  "NAModule4" : for the additionnal indoor module
//  "NAModule3" : for the rain gauge module
//  "NAModule2" : for the wind gauge module
// DashboardData : Data collection from device sensors
// DataType : List of available datas
// LinkedModules : Associated modules (only for station)
type Device struct {
	ID            string `json:"_id"`
	StationName   string `json:"station_name"`
	ModuleName    string `json:"module_name"`
	Type          string
	DashboardData DashboardData `json:"dashboard_data"`
	DataType      []string      `json:"data_type"`
	LinkedModules []*Device     `json:"modules"`
}

// DashboardData is used to store sensor values
// Temperature : Last temperature measure @ LastMesure (in °C)
// Humidity : Last humidity measured @ LastMesure (in %)
// CO2 : Last Co2 measured @ time_utc (in ppm)
// Noise : Last noise measured @ LastMesure (in db)
// Pressure : Last Sea level pressure measured @ LastMesure (in mb)
// AbsolutePressure : Real measured pressure @ LastMesure (in mb)
// Rain : Last rain measured (in mm)
// Rain1Hour : Amount of rain in last hour
// Rain1Day : Amount of rain today
// WindAngle : Current 5 min average wind direction @ LastMesure (in °)
// WindStrength : Current 5 min average wind speed @ LastMesure (in km/h)
// GustAngle : Direction of the last 5 min highest gust wind @ LastMesure (in °)
// GustStrength : Speed of the last 5 min highest gust wind @ LastMesure (in km/h)
// LastMessage : Contains timestamp of last data received
type DashboardData struct {
	Temperature      float32 `json:"Temperature,omitempty"`
	Humidity         int32   `json:"Humidity,omitempty"`
	CO2              int32   `json:"CO2,omitempty"`
	Noise            int32   `json:"Noise,omitempty"`
	Pressure         float32 `json:"Pressure,omitempty"`
	AbsolutePressure float32 `json:"AbsolutePressure,omitempty"`
	Rain             float32 `json:"Rain,omitempty"`
	Rain1Hour        float32 `json:"sum_rain_1,omitempty"`
	Rain1Day         float32 `json:"sum_rain_24,omitempty"`
	WindAngle        float32 `json:"WindAngle,omitempty"`
	WindStrength     float32 `json:"WindStrength,omitempty"`
	GustAngle        float32 `json:"GustAngle,omitempty"`
	GustStrength     float32 `json:"GustStrength,omitempty"`
	LastMeasure      float64 `json:"time_utc"`
}

var NAModuleMap = map[string][]string{
	"NAMain": []string{
		NAMainTemperature,
		NAMainHumidity,
		NAMainCO2,
		NAMainNoise,
		NAMainPressure,
		NAMainAbsolutePressure,
	},
	"NAModule1": []string{
		NAModule1Temperature,
		NAModule1Humidity,
	},
	"NAModule2": []string{
		NAModule2WindAngle,
		NAModule2WindStrength,
		NAModule2GustAngle,
		NAModule2GustStrength,
	},
	"NAModule3": []string{
		NAModule3Rain,
		NAModule3Rain1Hour,
		NAModule3Rain1Day,
	},
	"NAModule4": []string{
		NAModule4Temperature,
		NAModule4Humidity,
		NAModule4CO2,
	},
}

// Main module
const (
	NAMainTemperature      = "Temperature"
	NAMainHumidity         = "Humidity"
	NAMainCO2              = "CO2"
	NAMainNoise            = "Noise"
	NAMainPressure         = "Pressure"
	NAMainAbsolutePressure = "AbsolutePressure"
)

// Outdoor module
const (
	NAModule1Temperature = "Temperature"
	NAModule1Humidity    = "Humidity"
)

// Wind module
const (
	NAModule2WindAngle    = "WindAngle"
	NAModule2WindStrength = "WindStrength"
	NAModule2GustAngle    = "GustAngle"
	NAModule2GustStrength = "GustStrength"
)

// Rain module
const (
	NAModule3Rain      = "Rain"
	NAModule3Rain1Hour = "Rain1Hour"
	NAModule3Rain1Day  = "Rain1Day"
)

// Aux Indoor module
const (
	NAModule4Temperature = "Temperature"
	NAModule4Humidity    = "Humidity"
	NAModule4CO2         = "CO2"
)

// NewClient create a handle authentication to Netamo API
func NewClient(config Config) (*Client, error) {
	oauth := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"read_station"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  baseURL,
			TokenURL: authURL,
		},
	}

	token, err := oauth.PasswordCredentialsToken(oauth2.NoContext, config.Username, config.Password)

	return &Client{
		oauth:      oauth,
		httpClient: oauth.Client(oauth2.NoContext, token),
		Dc:         &DeviceCollection{},
	}, err
}

// do a url encoded HTTP POST request
func (c *Client) doHTTPPostForm(url string, data url.Values) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	//req.ContentLength = int64(reader.Len())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doHTTP(req)
}

// send http GET request
func (c *Client) doHTTPGet(url string, data url.Values) (*http.Response, error) {
	if data != nil {
		url = url + "?" + data.Encode()
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.doHTTP(req)
}

// do a generic HTTP request
func (c *Client) doHTTP(req *http.Request) (*http.Response, error) {

	// debug
	//debug, _ := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%s\n\n", debug)

	var err error
	c.httpResponse, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return c.httpResponse, nil
}

// process HTTP response
// Unmarshall received data into holder struct
func processHTTPResponse(resp *http.Response, err error, holder interface{}) error {
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// debug
	//debug, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("%s\n\n", debug)

	// check http return code
	if resp.StatusCode != 200 {
		//bytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Bad HTTP return code %d", resp.StatusCode)
	}

	// Unmarshall response into given struct
	if err = json.NewDecoder(resp.Body).Decode(holder); err != nil {
		return err
	}

	return nil
}

// GetStations returns the list of stations owned by the user, and their modules
func (c *Client) Read() (*DeviceCollection, error) {
	resp, err := c.doHTTPGet(deviceURL, url.Values{"app_type": {"app_station"}})
	//dc := &DeviceCollection{}

	if err = processHTTPResponse(resp, err, c.Dc); err != nil {
		return nil, err
	}

	return c.Dc, nil
}

// Devices returns the list of devices
func (dc *DeviceCollection) Devices() []*Device {
	return dc.Body.Devices
}

// Stations is an alias of Devices
func (dc *DeviceCollection) Stations() []*Device {
	return dc.Devices()
}

// Modules returns associated device module
func (d *Device) Modules() []*Device {
	modules := d.LinkedModules
	modules = append(modules, d)

	return modules
}

// Data returns timestamp and the list of sensor value for this module
func (d *Device) Data() (int, map[string]interface{}) {

	m := make(map[string]interface{})
	if moduleMap, ok := NAModuleMap[d.Type]; ok {
		for _, datatype := range moduleMap {
			m[datatype] = reflect.Indirect(reflect.ValueOf(d.DashboardData)).FieldByName(datatype).Interface()
		}
	}

	return int(d.DashboardData.LastMeasure), m
}
