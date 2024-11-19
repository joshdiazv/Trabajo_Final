// proxy/http_proxy.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/recommend", handleRecommendation)
	fmt.Println("Proxy escuchando en puerto 8081")
	http.ListenAndServe(":8081", nil)
}

// Handler para manejar las recomendaciones
func handleRecommendation(w http.ResponseWriter, r *http.Request) {
	// Leer el género de la solicitud
	genre := r.URL.Query().Get("genre")
	if genre == "" {
		http.Error(w, "El género es obligatorio", http.StatusBadRequest)
		return
	}

	// Hacer una solicitud al servidor para obtener las películas filtradas por género
	serverURL := fmt.Sprintf("http://server:8080/movies?genre=%s", genre)
	resp, err := http.Get(serverURL)
	if err != nil {
		http.Error(w, "Error al obtener las películas del servidor", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta del servidor", http.StatusInternalServerError)
		return
	}

	// Escribir la respuesta del servidor al cliente
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
