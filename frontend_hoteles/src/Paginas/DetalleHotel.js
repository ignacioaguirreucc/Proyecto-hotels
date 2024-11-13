import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { axiosHotelsInstance } from '../axiosConfig';
import styles from './DetalleHotel.module.css';

const DetalleHotel = ({ isAuthenticated }) => {
  const { id } = useParams();
  const [hotel, setHotel] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchHotel = async () => {
      try {
        const response = await axiosHotelsInstance.get(`/hotels/${id}`);
        console.log('Response from backend:', response.data); // Depuración
        setHotel(response.data);
      } catch (error) {
        console.error('Error al cargar el hotel:', error);
      }
    };
    fetchHotel();
  }, [id]);

  const handleReservation = async () => {
    if (!isAuthenticated) {
      alert("Debes iniciar sesión para reservar.");
      navigate('/login');
      return;
    }
    try {
      const token = localStorage.getItem('token');
      const userId = localStorage.getItem('user_id'); // Extrae el user_id de localStorage
  
      if (!userId) {
        console.error('No se encontró el user_id en localStorage');
        alert('Error al hacer la reserva: Usuario no autenticado.');
        return;
      }
  
      await axiosHotelsInstance.post(
        '/reservations',
        { hotel_id: hotel.id, user_id: userId, start_date: "2024-11-25", end_date: "2024-11-30" },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert('Reserva creada exitosamente');
    } catch (error) {
      console.error('Error al hacer la reserva:', error);
      alert('Error al hacer la reserva. Inténtalo nuevamente.');
    }
  };
  
  

  if (!hotel) return <div>Cargando...</div>;

  return (
    <div className={styles.detailContainer}>
      <h1>{hotel.name}</h1>
      <div className={styles.detailContent}>
        <p><strong>Dirección:</strong> {hotel.address}</p>
        <p><strong>Ciudad:</strong> {hotel.city}</p>
        <p><strong>Estado:</strong> {hotel.state}</p>
        <p><strong>Calificación:</strong> {hotel.rating} / 5</p>
        <p><strong>Amenities:</strong> {hotel.amenities && hotel.amenities.join(', ')}</p>
        <button className={styles.bookButton} onClick={handleReservation}>
          Reservar
        </button>
      </div>
    </div>
  );
};

export default DetalleHotel;
