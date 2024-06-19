package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"weather-check/internal/domain"
)

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) GetWeatherByLocale(locale string) (*domain.WeatherResult, error) {
	args := m.Called(locale)
	return args.Get(0).(*domain.WeatherResult), args.Error(1)
}

func (m *ApiClientMock) GetLocaleByZipcode(zipcode string) (*domain.LocaleResult, error) {
	args := m.Called(zipcode)
	return args.Get(0).(*domain.LocaleResult), args.Error(1)
}

func TestWeatherService_CheckWeather_InvalidZipcode(t *testing.T) {
	client := new(ApiClientMock)
	service := NewWeatherService(client)

	_, err := service.CheckWeather("12345")

	assert.Equal(t, ERROR_INVALID_ZIPCODE, err)

	_, err = service.CheckWeather("abc")

	assert.Equal(t, ERROR_INVALID_ZIPCODE, err)
}

func TestWeatherService_CheckWeather_CannotGetLocale(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{}, ERROR_CANNOT_GET_LOCALE)
	service := NewWeatherService(client)

	_, err := service.CheckWeather("12345678")

	assert.Equal(t, ERROR_CANNOT_GET_LOCALE, err)
}

func TestWeatherService_CheckWeather_CannotFindZipcode(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{}, nil)
	service := NewWeatherService(client)

	_, err := service.CheckWeather("12345678")

	assert.Equal(t, ERROR_CANNOT_FIND_ZIPCODE, err)
}

func TestWeatherService_CheckWeather_CannotGetWeather(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{Locale: "locale"}, nil)
	client.On("GetWeatherByLocale", "locale").Return(&domain.WeatherResult{}, ERROR_CANNOT_GET_WHEATER)
	service := NewWeatherService(client)

	_, err := service.CheckWeather("12345678")

	assert.Equal(t, ERROR_CANNOT_GET_WHEATER, err)
}

func TestWeatherService_CheckWeather_ShouldCorrectConvertCelsiusToKelvin(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{Locale: "locale"}, nil)
	client.On("GetWeatherByLocale", "locale").Return(&domain.WeatherResult{TempC: 10}, nil)
	service := NewWeatherService(client)

	result, _ := service.CheckWeather("12345678")

	assert.Equal(t, 283.15, result.TempK)
}

func TestWeatherService_CheckWeather_ShouldCorrectConvertCelsiusToFahrenheit(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{Locale: "locale"}, nil)
	client.On("GetWeatherByLocale", "locale").Return(&domain.WeatherResult{TempC: 10}, nil)
	service := NewWeatherService(client)

	result, _ := service.CheckWeather("12345678")

	assert.Equal(t, 50.0, result.TempF)
}

func TestWeatherService_CheckWeather_Success(t *testing.T) {
	client := new(ApiClientMock)
	client.On("GetLocaleByZipcode", "12345678").Return(&domain.LocaleResult{Locale: "locale"}, nil)
	client.On("GetWeatherByLocale", "locale").Return(&domain.WeatherResult{TempC: 10}, nil)
	service := NewWeatherService(client)

	result, _ := service.CheckWeather("12345678")

	assert.Equal(t, 10.0, result.TempC)
	assert.Equal(t, 50.0, result.TempF)
	assert.Equal(t, 283.15, result.TempK)
}
