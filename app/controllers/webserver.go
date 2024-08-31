package controllers

import (
	"encoding/json"
	"go-trade/app/models"
	"go-trade/config"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

var templates = template.Must(template.ParseGlob("app/views/chart.html"))

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var apiValidPath = regexp.MustCompile("^/api/candle/$")

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !apiValidPath.MatchString(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		productCode = config.Config.ProductCode
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Config.Durations[duration]

	_limit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(_limit)
	if err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	df, _ := models.GetAllCandle(productCode, durationTime, limit)

	json, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func StartWebServer() error {
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(":8080", nil)
}
