package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/obterCoordenadas", func(w http.ResponseWriter, r *http.Request) {
		destino := r.URL.Query().Get("destino")
		if destino == "" {
			http.Error(w, "Destino ausente", http.StatusBadRequest)
			return
		}

		locations, err := readCSV("locations.csv")
		if err != nil {
			http.Error(w, "Erro ao ler o CSV", http.StatusInternalServerError)
			return
		}

		var latitude, longitude float64
		encontrado := false
		for _, loc := range locations {
			if strings.TrimSpace(loc.Name) == destino {
				latitude = loc.Latitude
				longitude = loc.Longitude
				encontrado = true
				break
			}
		}

		if !encontrado {
			http.Error(w, "Localização não encontrada", http.StatusNotFound)
			return
		}

		resposta := struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  latitude,
			Longitude: longitude,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resposta)
	})

	http.ListenAndServe(":8080", nil)
}

func readCSV(filename string) ([]Location, error) {
	var locations []Location

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		latitude, _ := strconv.ParseFloat(record[1], 64)
		longitude, _ := strconv.ParseFloat(record[2], 64)

		location := Location{
			Name:      strings.TrimSpace(record[0]),
			Latitude:  latitude,
			Longitude: longitude,
		}

		locations = append(locations, location)
	}

	return locations, nil
}
