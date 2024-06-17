package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"regexp"
)

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

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherCheckResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherCheckHandler struct{}

func NewWeatherCheckHandler() *WeatherCheckHandler {
	return &WeatherCheckHandler{}
}

func (h *WeatherCheckHandler) isValidCEP(cep string) bool {
	pattern := `^\d{5}-?\d{3}$`
	match, _ := regexp.MatchString(pattern, cep)
	return match
}

func (h *WeatherCheckHandler) getWeather(locale string) (*WeatherAPIResponse, error) {
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

	return &weatherAPIResponse, nil
}

func (h *WeatherCheckHandler) getCep(cep string) (*ViaCepResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

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

	return &viaCepResponse, nil
}

func (h *WeatherCheckHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if !h.isValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	viaCepResponse, err := h.getCep(cep)
	if err != nil {
		log.Println(err)
		http.Error(w, "error on get zipcode", http.StatusInternalServerError)
	}

	if viaCepResponse.Localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weatherAPIResponse, err := h.getWeather(viaCepResponse.Localidade)
	if err != nil {
		log.Println(err)
		http.Error(w, "error on get weather", http.StatusInternalServerError)
		return
	}

	weatherCheckResponse := WeatherCheckResponse{
		TempC: weatherAPIResponse.Current.TempC,
		TempF: weatherAPIResponse.Current.TempF,
		TempK: weatherAPIResponse.Current.TempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(weatherCheckResponse)

	if err != nil {
		http.Error(w, "error on encode response", http.StatusInternalServerError)
		return
	}
}
