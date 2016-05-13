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
	deviceURL = baseURL + "api/devicelist"
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
// ClientID : Client ID from netatmo app registration at http://dev.netatmo.com/dev/listapps
// ClientSecret : Client app secret
// Username : Your netatmo account username
// Password : Your netatmo account password
// Stations : Contains all Station account
type Client struct {
	oauth        *oauth2.Config
	httpClient   *http.Client
	httpResponse *http.Response
}

// DeviceCollection hold all devices from netatmo account (stations and modules)
// Error : returned error (nil if OK)
// Stations : List of stations
// Modules : List of additionnal modules
type DeviceCollection struct {
	Body struct {
		Stations []*Device `json:"devices"`
		Modules  []*Device
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
// MainDevice : Id of main station (only for module)
// AssociatedModules : Associated modules (only for station)
type Device struct {
	ID                string `json:"_id"`
	StationName       string `json:"station_name"`
	ModuleName        string `json:"module_name"`
	Type              string
	DashboardData     DashboardData `json:"dashboard_data"`
	DataType          []string      `json:"data_type"`
	MainDevice        string        `json:"main_device,omitempty"`
	AssociatedModules []*Device     `json:"-"`
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
	Temperature         float32 `json:"Temperature,omitempty"`
	Humidity            int32   `json:"Humidity,omitempty"`
	CO2                 int32   `json:"CO2,omitempty"`
	Noise               int32   `json:"Noise,omitempty"`
	Pressure            float32 `json:"Pressure,omitempty"`
	AbsolutePressure    float32 `json:"AbsolutePressure,omitempty"`
	Rain                float32 `json:"Rain,omitempty"`
	Rain1Hour           float32 `json:"sum_rain_1,omitempty"`
	Rain1Day            float32 `json:"sum_rain_24,omitempty"`
	WindAngle           float32 `json:"WindAngle,omitempty"`
	WindStrength        float32 `json:"WindStrength,omitempty"`
	GustAngle           float32 `json:"GustAngle,omitempty"`
	GustStrengthfloat32 float32 `json:"GustStrengthfloat32,omitempty"`
	LastMesure          float64 `json:"time_utc"`
}

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

// GetDeviceCollection returns the list of stations owned by the user, and their modules
func (c *Client) GetDeviceCollection() (*DeviceCollection, error) {
	//resp, err := c.doHTTPPostForm(deviceURL, url.Values{"app_type": {"app_station"}})
	resp, err := c.doHTTPGet(deviceURL, url.Values{"app_type": {"app_station"}})
	dc := &DeviceCollection{}

	if err = processHTTPResponse(resp, err, dc); err != nil {
		return nil, err
	}

	// associated each module to its station
	for i, station := range dc.Body.Stations {
		for _, module := range dc.Body.Modules {
			if module.MainDevice == station.ID {
				dc.Body.Stations[i].AssociatedModules = append(dc.Body.Stations[i].AssociatedModules, module)
			}
		}
	}

	return dc, nil
}

// Stations returns the list of stations
func (dc *DeviceCollection) Stations() []*Device {
	return dc.Body.Stations
}

// Modules returns the list of modules associated to this station
// also return station itself in the list
func (s *Device) Modules() []*Device {
	modules := s.AssociatedModules

	modules = append(modules, s)
	return modules
}

// Data returns timestamp and the list of sensor value for this module
func (s *Device) Data() (int, map[string]interface{}) {

	m := make(map[string]interface{})
	for _, datatype := range s.DataType {
		m[datatype] = reflect.Indirect(reflect.ValueOf(s.DashboardData)).FieldByName(datatype).Interface()
	}

	return int(s.DashboardData.LastMesure), m
}
