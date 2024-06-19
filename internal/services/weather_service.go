package services

import (
	"errors"
	"regexp"
	"weather-check/internal/adapters/api"
	"weather-check/internal/domain"
)

var (
	ERROR_INVALID_ZIPCODE     = errors.New("invalid zipcode")
	ERROR_CANNOT_GET_LOCALE   = errors.New("error on get zipcode")
	ERROR_CANNOT_FIND_ZIPCODE = errors.New("can not find zipcode")
	ERROR_CANNOT_GET_WHEATER  = errors.New("error on get weather")
)

type WeatherServiceInterface interface {
	CheckWeather(zipcode string) (*domain.Result, error)
}

type WeatherService struct {
	Client api.ClientInterface
}

func NewWeatherService(client api.ClientInterface) *WeatherService {
	return &WeatherService{
		Client: client,
	}
}

func (w *WeatherService) isValidZipcode(zipcode string) bool {
	pattern := `^\d{5}-?\d{3}$`
	match, _ := regexp.MatchString(pattern, zipcode)
	return match
}

func (w *WeatherService) CheckWeather(zipcode string) (*domain.Result, error) {
	if !w.isValidZipcode(zipcode) {
		return nil, ERROR_INVALID_ZIPCODE
	}

	locale, err := w.Client.GetLocaleByZipcode(zipcode)
	if err != nil {
		return nil, ERROR_CANNOT_GET_LOCALE
	}

	if locale.Locale == "" {
		return nil, ERROR_CANNOT_FIND_ZIPCODE
	}

	weather, err := w.Client.GetWeatherByLocale(locale.Locale)
	if err != nil {
		return nil, ERROR_CANNOT_GET_WHEATER
	}

	return &domain.Result{
		TempC: weather.TempC,
		TempF: weather.TempC*1.8 + 32,
		TempK: weather.TempC + 273.15,
	}, nil
}
