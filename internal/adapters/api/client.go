package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"weather-check/internal/domain"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ClientInterface interface {
	GetWeatherByLocale(locale string) (*domain.WeatherResult, error)
	GetLocaleByZipcode(zipcode string) (*domain.LocaleResult, error)
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetWeatherByLocale(locale string) (*domain.WeatherResult, error) {
	token := os.Getenv("WEATHER_API_KEY")

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", token, netUrl.QueryEscape(locale))

	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	var weatherAPIResponse WeatherAPIResponse
	err = json.NewDecoder(req.Body).Decode(&weatherAPIResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &domain.WeatherResult{
		TempC: weatherAPIResponse.Current.TempC,
	}, nil
}

func (c *Client) GetLocaleByZipcode(zipcode string) (*domain.LocaleResult, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)

	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	var viaCepResponse ViaCepResponse
	err = json.NewDecoder(req.Body).Decode(&viaCepResponse)
	if err != nil {
		return nil, err
	}

	return &domain.LocaleResult{
		Locale: viaCepResponse.Localidade,
	}, nil
}
