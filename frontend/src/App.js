import React, { useEffect, useState } from 'react';
import './App.css';

function App() {
  const [genres, setGenres] = useState([]);
  const [selectedGenre, setSelectedGenre] = useState('');
  const [recommendations, setRecommendations] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Fetch géneros disponibles al cargar la página
  useEffect(() => {
    fetchGenres();
  }, []);

  // Función para obtener géneros desde la API (servidor)
  const fetchGenres = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://server:8080/getGenres');
      if (!response.ok) throw new Error('Error al obtener los géneros');
      const data = await response.json();
      setGenres(data);
    } catch (err) {
      setError('Error al obtener los géneros');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Función para obtener las recomendaciones basadas en el género seleccionado
  const getRecommendations = async () => {
    setLoading(true);
    try {
      const response = await fetch(`http://server:8080/getRecommendations?genre=${selectedGenre}`);
      if (!response.ok) throw new Error('Error al obtener las recomendaciones');
      const data = await response.json();
      setRecommendations(data);
    } catch (err) {
      setError('Error al obtener las recomendaciones');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Manejar el cambio en la selección de género
  const handleGenreChange = (event) => {
    setSelectedGenre(event.target.value);
  };

  // Manejar el envío de la selección de género para obtener las recomendaciones
  const handleGetRecommendations = () => {
    if (!selectedGenre) {
      setError('Por favor selecciona un género');
      return;
    }
    getRecommendations();
  };

  return (
    <div className="App">
      <h1>Recomendador de Películas</h1>

      {loading && <div className="loading-spinner">Cargando...</div>}

      <div>
        <label htmlFor="genre">Selecciona un género: </label>
        <select id="genre" value={selectedGenre} onChange={handleGenreChange}>
          <option value="">Selecciona un género</option>
          {genres.map((genre, index) => (
            <option key={index} value={genre}>{genre}</option>
          ))}
        </select>
      </div>

      <button onClick={handleGetRecommendations} disabled={!selectedGenre || loading}>
        Obtener Recomendaciones
      </button>

      {error && <div className="error">{error}</div>}

      {recommendations.length > 0 && (
        <div>
          <h2>Películas recomendadas:</h2>
          <ul>
            {recommendations.map((movie, index) => (
              <li key={index}>
                <strong>{movie.title}</strong> - Géneros: {movie.genres.join(', ')} - Calificación: {movie.avgRating}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default App;
