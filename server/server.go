package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Movie struct {
	MovieID int
	Title   string
	Genres  []string // Lista de géneros
}

type Rating struct {
	UserID  int
	MovieID int
	Rating  float64
}

type Recommendation struct {
	MovieID   int
	Title     string
	Genres    []string
	AvgRating float64
	Count     int
}

var movies = make(map[int]Movie)
var ratings = make([]Rating, 0)

// Mapa para almacenar las recomendaciones de diferentes clientes y combinarlas
var genreRecommendations = make(map[string]map[int]*Recommendation)
var mutex sync.Mutex // Mutex para proteger el acceso concurrente al mapa

func main() {
	// Cargar películas desde el archivo
	loadMovies("movies.csv")

	// Configurar el servidor HTTP para escuchar en el puerto 8080
	http.HandleFunc("/movies", getMoviesByGenreHandler)

	// Iniciar el servidor
	fmt.Println("Servidor escuchando en puerto 8080")
	http.ListenAndServe(":8080", nil)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	switch {
	case strings.HasPrefix(message, "GET_GENRES"):
		// Paso 1: Enviar instrucción al cliente
		fmt.Fprintln(conn, "Por favor, elige un género del listado siguiente:")

		// Paso 2: Enviar la lista de géneros
		genres := getTopGenres()
		for _, genre := range genres {
			fmt.Fprintln(conn, genre)
		}
		fmt.Fprintln(conn, "[END_OF_GENRES]") // Indicar el final de la lista de géneros

	case strings.HasPrefix(message, "TASK_GENRE:"):
		genre := strings.TrimPrefix(message, "TASK_GENRE:")
		go processGenreTask(conn, genre) // Procesa la tarea para el género en paralelo

	default:
		fmt.Fprintln(conn, "Comando no reconocido")
	}
}

// Endpoint para obtener películas por género
func getMoviesByGenreHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el género desde la consulta (query)
	genre := r.URL.Query().Get("genre")
	if genre == "" {
		http.Error(w, "El género es obligatorio", http.StatusBadRequest)
		return
	}

	// Filtrar las películas por género
	recommendedMovies := getMoviesByGenre(genre)

	// Enviar la lista de películas como respuesta en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendedMovies)
}

func processGenreTask(conn net.Conn, genre string) {
	// Paso 1: Obtener las películas recomendadas por género
	recommendedMovies := getMoviesByGenre(genre)

	// Paso 2: Enviar las recomendaciones de películas al cliente
	fmt.Fprintln(conn, "Películas recomendadas:")
	for _, movie := range recommendedMovies {
		fmt.Fprintf(conn, "Título: %s, Calificación Promedio: %.2f\n", movie.Title, calculateAverageRating(movie.MovieID))
	}
}

func combineRecommendations(genre string, recommendations []Movie) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := genreRecommendations[genre]; !exists {
		genreRecommendations[genre] = make(map[int]*Recommendation)
	}

	for _, movie := range recommendations {
		avgRating := calculateAverageRating(movie.MovieID)
		if rec, exists := genreRecommendations[genre][movie.MovieID]; exists {
			// Si ya existe una recomendación, actualizar el promedio de calificación
			rec.AvgRating = (rec.AvgRating*float64(rec.Count) + avgRating) / float64(rec.Count+1)
			rec.Count++
		} else {
			// Si es una nueva recomendación, agregarla al mapa
			genreRecommendations[genre][movie.MovieID] = &Recommendation{
				MovieID:   movie.MovieID,
				Title:     movie.Title,
				Genres:    movie.Genres,
				AvgRating: avgRating,
				Count:     1,
			}
		}
	}
}

func displayCombinedRecommendations(writer *bufio.Writer, genre string) {
	mutex.Lock()
	defer mutex.Unlock()

	if recs, exists := genreRecommendations[genre]; exists {
		fmt.Fprintln(writer, "Películas recomendadas combinadas para el género:", genre)
		count := 0
		for _, rec := range recs {
			if count >= 5 { // Limitar a las 5 primeras
				break
			}
			fmt.Fprintf(writer, "%d. Título: %s, Géneros: %s, Calificación Promedio Combinada: %.2f\n",
				count+1, rec.Title, strings.Join(rec.Genres, ", "), rec.AvgRating)
			count++
		}
	} else {
		fmt.Fprintln(writer, "No se encontraron recomendaciones para el género seleccionado.")
	}
	writer.Flush()
}

// Función para calcular la calificación promedio de una película
func calculateAverageRating(movieID int) float64 {
	var totalRating float64
	var count int
	for _, rating := range ratings {
		if rating.MovieID == movieID {
			totalRating += rating.Rating
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return totalRating / float64(count)
}

func getTopGenres() []string {
	genreCount := make(map[string]int)
	for _, movie := range movies {
		for _, genre := range movie.Genres {
			genreCount[genre]++
		}
	}

	var genreList []string
	for genre := range genreCount {
		genreList = append(genreList, genre)
	}

	sort.Slice(genreList, func(i, j int) bool {
		return genreCount[genreList[i]] > genreCount[genreList[j]]
	})

	if len(genreList) > 15 {
		genreList = genreList[:15]
	}
	return genreList
}

func getMoviesByGenre(preferredGenre string) []Movie {
	var recommendedMovies []Movie
	for _, movie := range movies {
		for _, genre := range movie.Genres {
			if strings.Contains(strings.ToLower(genre), strings.ToLower(preferredGenre)) {
				recommendedMovies = append(recommendedMovies, movie)
				break
			}
		}
	}

	if len(recommendedMovies) > 5 {
		recommendedMovies = recommendedMovies[:5]
	}
	return recommendedMovies
}

// Cargar las películas desde un archivo CSV
func loadMovies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error al abrir archivo de películas:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error al leer archivo de películas:", err)
		return
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		movieID, _ := strconv.Atoi(record[0])
		genres := strings.Split(record[2], "|")

		movies[movieID] = Movie{
			MovieID: movieID,
			Title:   record[1],
			Genres:  genres,
		}
	}
}

// Cargar las calificaciones desde un archivo CSV
func loadRatings(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error al abrir archivo de calificaciones:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error al leer archivo de calificaciones:", err)
		return
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		userID, _ := strconv.Atoi(record[0])
		movieID, _ := strconv.Atoi(record[1])
		rating, _ := strconv.ParseFloat(record[2], 64)

		ratings = append(ratings, Rating{
			UserID:  userID,
			MovieID: movieID,
			Rating:  rating,
		})
	}
}
