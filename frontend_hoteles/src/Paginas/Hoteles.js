import React, { useState, useEffect } from 'react';
import { axiosSearchInstance } from '../axiosConfig';
import { Link } from 'react-router-dom';
import styles from './Hoteles.module.css';

const Hoteles = () => {
  const [hotels, setHotels] = useState([]);  // Estado para los hoteles obtenidos del backend
  const [error, setError] = useState(null);  // Estado para manejar errores

  useEffect(() => {
    const fetchHotels = async () => {
      try {
        // Cambia el valor de `q` a `*` en lugar de `*:*`
        const response = await axiosSearchInstance.get('/search', {
          params: { q: '*', offset: 0, limit: 10 }  // Consulta todos los hoteles con un límite
        });

        // Mapea los resultados de Solr a los atributos del hotel en frontend
        const hotelsData = response.data.map((hotel) => ({
          id: hotel.id,
          name: Array.isArray(hotel.name) ? hotel.name[0] : hotel.name,
          rating: Array.isArray(hotel.rating) ? hotel.rating[0] : hotel.rating,
          amenities: Array.isArray(hotel.amenities) ? hotel.amenities.join(", ") : 'No disponible',
          city: Array.isArray(hotel.city) ? hotel.city[0] : 'Ubicación no disponible',
          address: Array.isArray(hotel.address) ? hotel.address[0] : 'Ubicación no disponible',
        }));

        setHotels(hotelsData);  // Guarda los hoteles en el estado
      } catch (err) {
        console.error('Error al cargar los hoteles:', err);
        setError('Hubo un problema al cargar los hoteles.');
      }
    };

    fetchHotels();
  }, []);

  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.resultsContainer}>
      <h1 className={styles.heading}>Resultados de Hoteles</h1>
      <p className={styles.subheading}>Explora las mejores opciones para tu próxima estadía.</p>

      <div className={styles.hotelList}>
        {hotels.map((hotel) => (
          <div key={hotel.id} className={styles.hotelCard}>
            <div className={styles.hotelDetails}>
              <h2 className={styles.hotelName}>{hotel.name}</h2>
              <div className={styles.hotelRating}>Puntuación: {hotel.rating}</div>
              <p><strong>Amenities:</strong> {hotel.amenities}</p>
              <p><strong>Ubicación:</strong> {hotel.city}, {hotel.state}</p>
              <Link to={`/detalle-hotel/${hotel.id}`} className={styles.detailButton}>
                Ver detalles
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Hoteles;
