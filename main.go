package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
}

type GeocodeResult struct {
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"`
}

type WeatherResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Timezone             string  `json:"timezone"`
	Current              Current `json:"current"`
	CurrentUnits         CurrentUnits `json:"current_units"`
	Hourly               *Hourly `json:"hourly,omitempty"`
	HourlyUnits          *HourlyUnits `json:"hourly_units,omitempty"`
	Daily                *Daily `json:"daily,omitempty"`
	DailyUnits           *DailyUnits `json:"daily_units,omitempty"`
}

type Current struct {
	Time                 string  `json:"time"`
	Temperature2m        float64 `json:"temperature_2m"`
	RelativeHumidity2m   int     `json:"relative_humidity_2m"`
	ApparentTemperature  float64 `json:"apparent_temperature"`
	Precipitation        float64 `json:"precipitation"`
	WeatherCode          int     `json:"weather_code"`
	WindSpeed10m         float64 `json:"wind_speed_10m"`
	WindDirection10m     float64 `json:"wind_direction_10m"`
}

type CurrentUnits struct {
	Time                 string `json:"time"`
	Temperature2m        string `json:"temperature_2m"`
	RelativeHumidity2m   string `json:"relative_humidity_2m"`
	ApparentTemperature  string `json:"apparent_temperature"`
	Precipitation        string `json:"precipitation"`
	WeatherCode          string `json:"weather_code"`
	WindSpeed10m         string `json:"wind_speed_10m"`
	WindDirection10m     string `json:"wind_direction_10m"`
}

type Hourly struct {
	Time            []string  `json:"time"`
	Temperature2m   []float64 `json:"temperature_2m"`
	WeatherCode     []int     `json:"weather_code"`
	Precipitation   []float64 `json:"precipitation"`
}

type HourlyUnits struct {
	Time            string `json:"time"`
	Temperature2m   string `json:"temperature_2m"`
	WeatherCode     string `json:"weather_code"`
	Precipitation   string `json:"precipitation"`
}

type Daily struct {
	Time                []string  `json:"time"`
	Temperature2mMax    []float64 `json:"temperature_2m_max"`
	Temperature2mMin    []float64 `json:"temperature_2m_min"`
	WeatherCode         []int     `json:"weather_code"`
	PrecipitationSum    []float64 `json:"precipitation_sum"`
}

type DailyUnits struct {
	Time                string `json:"time"`
	Temperature2mMax    string `json:"temperature_2m_max"`
	Temperature2mMin    string `json:"temperature_2m_min"`
	WeatherCode         string `json:"weather_code"`
	PrecipitationSum    string `json:"precipitation_sum"`
}

func main() {
	var forecast24h = flag.Bool("24h", false, "Show 24-hour forecast")
	var forecast7d = flag.Bool("7d", false, "Show 7-day forecast")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: gocast [options] <location> [country]")
		fmt.Println("Options:")
		fmt.Println("  -24h    Show 24-hour forecast")
		fmt.Println("  -7d     Show 7-day forecast")
		fmt.Println("Examples:")
		fmt.Println("  gocast London")
		fmt.Println("  gocast London GB")
		fmt.Println("  gocast -24h \"New York\" US")
		fmt.Println("  gocast -7d \"San Francisco\"")
		os.Exit(1)
	}

	args := flag.Args()
	var location, country string
	
	if len(args) >= 2 {
		location = strings.Join(args[:len(args)-1], " ")
		country = args[len(args)-1]
	} else {
		location = strings.Join(args, " ")
	}

	coords, err := geocodeLocation(location, country)
	if err != nil {
		log.Fatal(err)
	}

	weather, err := getWeather(coords.Latitude, coords.Longitude, *forecast24h, *forecast7d)
	if err != nil {
		log.Fatal(err)
	}

	displayWeather(weather, coords, *forecast24h, *forecast7d)
}

func geocodeLocation(location, country string) (*GeocodeResult, error) {
	baseURL := "https://geocoding-api.open-meteo.com/v1/search"
	params := url.Values{}
	params.Add("name", location)
	if country != "" {
		params.Add("count", "10")
	} else {
		params.Add("count", "1")
	}
	params.Add("language", "en")
	
	if country != "" {
		countryCode := getCountryCode(country)
		if countryCode != "" {
			params.Add("country_code", countryCode)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %v", err)
	}
	defer resp.Body.Close()

	var geocode GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocode); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %v", err)
	}

	if len(geocode.Results) == 0 {
		if country != "" {
			return nil, fmt.Errorf("location '%s' not found in %s", location, country)
		}
		return nil, fmt.Errorf("location '%s' not found", location)
	}

	// If country is specified, find the first result matching that country
	if country != "" {
		countryCode := getCountryCode(country)
		for _, result := range geocode.Results {
			if result.CountryCode != "" && strings.EqualFold(result.CountryCode, countryCode) {
				return &result, nil
			}
		}
		return nil, fmt.Errorf("location '%s' not found in %s", location, country)
	}

	return &geocode.Results[0], nil
}

func getCountryCode(country string) string {
	country = strings.ToUpper(strings.TrimSpace(country))
	
	// If already a 2-letter code, return as-is
	if len(country) == 2 {
		return country
	}
	
	// Common country mappings
	countryMap := map[string]string{
		"UK":                  "GB",
		"UNITED KINGDOM":      "GB",
		"GREAT BRITAIN":       "GB",
		"ENGLAND":            "GB",
		"SCOTLAND":           "GB", 
		"WALES":              "GB",
		"NORTHERN IRELAND":   "GB",
		"US":                 "US",
		"USA":                "US",
		"UNITED STATES":      "US",
		"AMERICA":            "US",
		"CANADA":             "CA",
		"FRANCE":             "FR",
		"GERMANY":            "DE",
		"ITALY":              "IT",
		"SPAIN":              "ES",
		"PORTUGAL":           "PT",
		"NETHERLANDS":        "NL",
		"HOLLAND":            "NL",
		"BELGIUM":            "BE",
		"SWITZERLAND":        "CH",
		"AUSTRIA":            "AT",
		"DENMARK":            "DK",
		"SWEDEN":             "SE",
		"NORWAY":             "NO",
		"FINLAND":            "FI",
		"IRELAND":            "IE",
		"AUSTRALIA":          "AU",
		"NEW ZEALAND":        "NZ",
		"JAPAN":              "JP",
		"CHINA":              "CN",
		"INDIA":              "IN",
		"RUSSIA":             "RU",
		"BRAZIL":             "BR",
		"MEXICO":             "MX",
		"SOUTH AFRICA":       "ZA",
	}
	
	if code, exists := countryMap[country]; exists {
		return code
	}
	
	// Return empty string if not found - API will work without country filter
	return ""
}

func getWeather(lat, lon float64, forecast24h, forecast7d bool) (*WeatherResponse, error) {
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.6f", lat))
	params.Add("longitude", fmt.Sprintf("%.6f", lon))
	params.Add("current", "temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,weather_code,wind_speed_10m,wind_direction_10m")
	
	if forecast24h {
		params.Add("hourly", "temperature_2m,weather_code,precipitation")
		params.Add("forecast_hours", "24")
	}
	
	if forecast7d {
		params.Add("daily", "temperature_2m_max,temperature_2m_min,weather_code,precipitation_sum")
		params.Add("forecast_days", "7")
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

func displayWeather(weather *WeatherResponse, location *GeocodeResult, forecast24h, forecast7d bool) {
	ascii := getWeatherASCII(weather.Current.WeatherCode)
	
	fmt.Printf("\n%s\n", ascii)
	locationStr := location.Name
	if location.Admin1 != "" && location.Admin1 != location.Name {
		locationStr += ", " + location.Admin1
	}
	if location.Country != "" {
		locationStr += ", " + location.Country
	}
	fmt.Printf("📍 Location: %s\n", locationStr)
	fmt.Printf("🌡️  Temperature: %.1f%s (feels like %.1f%s)\n", 
		weather.Current.Temperature2m, weather.CurrentUnits.Temperature2m,
		weather.Current.ApparentTemperature, weather.CurrentUnits.ApparentTemperature)
	fmt.Printf("💧 Humidity: %d%s\n", weather.Current.RelativeHumidity2m, weather.CurrentUnits.RelativeHumidity2m)
	fmt.Printf("🌧️  Precipitation: %.1f%s\n", weather.Current.Precipitation, weather.CurrentUnits.Precipitation)
	fmt.Printf("💨 Wind: %.1f%s at %.0f°\n", 
		weather.Current.WindSpeed10m, weather.CurrentUnits.WindSpeed10m, weather.Current.WindDirection10m)
	fmt.Printf("⏰ Updated: %s\n\n", weather.Current.Time)
	
	if forecast24h && weather.Hourly != nil {
		displayHourlyForecast(weather)
	}
	
	if forecast7d && weather.Daily != nil {
		displayDailyForecast(weather)
	}
}

func getWeatherASCII(code int) string {
	switch {
	case code == 0:
		return `
    \   /    
     .-.     
  ‒ (   ) ‒  
     '-'     
    /   \    
   Clear Sky  `
	case code >= 1 && code <= 3:
		return `
   .--.      
.-(    ).    
(___.__)__)  
 Partly Cloudy`
	case code >= 45 && code <= 48:
		return `
_ - _ - _ -   
 _ - _ - _    
_ - _ - _ -   
    Fog      `
	case code >= 51 && code <= 67:
		return `
     .-.     
    (   ).   
   (___(__)  
  ‚ ‚ ‚ ‚    
 ‚ ‚ ‚ ‚     
   Drizzle   `
	case code >= 80 && code <= 82:
		return `
     .-.     
    (   ).   
   (___(__)  
  ‚'‚'‚'‚'   
 ‚'‚'‚'‚'    
  Rain Shower`
	case code >= 95 && code <= 99:
		return `
     .-.     
    (   ).   
   (___(__)  
  ‚'⚡'‚'⚡   
 ‚'‚'‚'‚'    
 Thunderstorm`
	default:
		return `
      .--.   
   .-(    ). 
  (___.__)__)
    Cloudy   `
	}
}

func displayHourlyForecast(weather *WeatherResponse) {
	fmt.Println("🕐 24-Hour Forecast:")
	fmt.Println("────────────────────")
	
	for i := 0; i < len(weather.Hourly.Time) && i < 24; i++ {
		timeStr := weather.Hourly.Time[i]
		temp := weather.Hourly.Temperature2m[i]
		precip := weather.Hourly.Precipitation[i]
		code := weather.Hourly.WeatherCode[i]
		
		parsedTime, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			parsedTime, _ = time.Parse(time.RFC3339, timeStr)
		}
		
		icon := getWeatherIcon(code)
		fmt.Printf("%s %s %.1f%s (%.1fmm)\n", 
			parsedTime.Format("15:04"), icon, temp, weather.HourlyUnits.Temperature2m, precip)
	}
	fmt.Println()
}

func displayDailyForecast(weather *WeatherResponse) {
	fmt.Println("📅 7-Day Forecast:")
	fmt.Println("─────────────────")
	
	for i := 0; i < len(weather.Daily.Time) && i < 7; i++ {
		dateStr := weather.Daily.Time[i]
		maxTemp := weather.Daily.Temperature2mMax[i]
		minTemp := weather.Daily.Temperature2mMin[i]
		precip := weather.Daily.PrecipitationSum[i]
		code := weather.Daily.WeatherCode[i]
		
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			parsedDate, _ = time.Parse(time.RFC3339, dateStr)
		}
		
		icon := getWeatherIcon(code)
		fmt.Printf("%s %s %.1f/%.1f%s (%.1fmm)\n", 
			parsedDate.Format("Mon Jan 2"), icon, maxTemp, minTemp, weather.DailyUnits.Temperature2mMax, precip)
	}
	fmt.Println()
}

func getWeatherIcon(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code >= 1 && code <= 3:
		return "⛅"
	case code >= 45 && code <= 48:
		return "🌫️"
	case code >= 51 && code <= 67:
		return "🌦️"
	case code >= 80 && code <= 82:
		return "🌧️"
	case code >= 95 && code <= 99:
		return "⛈️"
	default:
		return "☁️"
	}
}