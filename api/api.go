package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	// Intentar conectar al servidor con reintentos
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", "server:8080")
		if err != nil {
			fmt.Println("Error al conectar al servidor, reintentando en 3 segundos...")
			time.Sleep(3 * time.Second) // Espera 3 segundos antes de reintentar
			continue
		}
		break // Sale del bucle si la conexión fue exitosa
	}
	defer conn.Close()

	// Leer la respuesta del servidor
	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')
	fmt.Print(line) // Bienvenida

	// Ingresar el userID
	fmt.Print("Ingresa tu ID de usuario: ")
	var userID int
	fmt.Scanln(&userID)
	fmt.Fprintln(conn, userID)

	// Leer la respuesta del servidor después de ingresar el ID
	line, _ = reader.ReadString('\n')
	fmt.Print(line) // "Archivos CSV leídos correctamente"

	// Leer y mostrar los géneros disponibles
	line, _ = reader.ReadString('\n')
	fmt.Print(line) // Mostrar la instrucción de elegir un género

	// Leer la lista completa de géneros
	var genres []string
	for {
		line, _ = reader.ReadString('\n')
		if line == "[END_OF_GENRES]\n" { // Fin de la lista de géneros
			break
		}
		fmt.Print(line)                                  // Mostrar géneros
		genres = append(genres, strings.TrimSpace(line)) // Eliminar posibles saltos de línea
	}

	// Seleccionar un género
	var genreIndex int
	fmt.Print("Por favor, selecciona un género por el número: ")
	fmt.Scanln(&genreIndex)

	// Validar la selección
	if genreIndex < 1 || genreIndex > len(genres) {
		fmt.Println("Selección inválida")
		return
	}

	// Enviar la selección al servidor
	fmt.Fprintf(conn, "TASK_GENRE:%s\n", genres[genreIndex-1])

	// Leer y mostrar las películas recomendadas con calificación
	fmt.Println("Películas recomendadas:")
	for {
		line, _ = reader.ReadString('\n')
		if line == "\n" { // Fin de las recomendaciones
			break
		}
		fmt.Print(line) // Mostrar películas con calificación promedio
	}
}
