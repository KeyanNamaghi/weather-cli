package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		Temp float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forcast struct {
		Forecastday []struct {
			Hour []struct {
				Temp float64 `json:"temp_c"`
				TimeEpoch int64 `json:"time_epoch"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset string `json:"sunset"`
			} `json:"astro"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main () {
	apiKey := os.Getenv("WEATHER_API_KEY")

	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=London&days=1&aqi=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Status code error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, temperature, condition, hours := weather.Location.Name, weather.Current.Temp, weather.Current.Condition.Text, weather.Forcast.Forecastday[0].Hour
	fmt.Printf("\033[33m%s: %.0f°C - %s\n\033[0m", location, temperature, condition)
	
	currentTime := time.Now().Format("2006-01-02")
	sunriseTimeString := weather.Forcast.Forecastday[0].Astro.Sunrise
	sunsetTimeString := weather.Forcast.Forecastday[0].Astro.Sunset
	sunriseCombinedString := currentTime + " " + sunriseTimeString
	sunsetCombinedString := currentTime + " " + sunsetTimeString

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic("Error loading location")
	}

	// Parse the combined string in the desired timezone
	sunriseTime, err := time.ParseInLocation("2006-01-02 03:04 AM", sunriseCombinedString, loc)
	if err != nil {
		panic("Error parsing time")
	}
	sunsetTime, err := time.ParseInLocation("2006-01-02 03:04 PM", sunsetCombinedString, loc)
	if err != nil {
		panic("Error parsing time")
	}

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch , 0)

		if date.After(sunriseTime) && date.Before(sunriseTime.Add(1 * time.Hour)) && date.After(time.Now()){
			fmt.Printf("\033[31m%s: Sunrise\n\033[0m", sunriseTime.Format("15:04"))
		}

		if date.After(sunsetTime) && date.Before(sunsetTime.Add(1 * time.Hour)) && date.After(time.Now()) {
			fmt.Printf("\033[31m%s: Sunset\n\033[0m", sunsetTime.Format("15:04"))
		}

		if date.After(time.Now()) {
			fmt.Printf("%s: %.0f°C - %s\n", date.Format("15:04"), hour.Temp, hour.Condition.Text)
		}
	}
}
