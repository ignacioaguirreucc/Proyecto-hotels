import React, { useEffect, useState } from 'react';
import { axiosHotelsInstance } from '../axiosConfig';
import styles from './Reservas.module.css';

const Reservas = () => {
  const [reservations, setReservations] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchReservations = async () => {
      try {
        const token = localStorage.getItem('token');
        const userId = localStorage.getItem('user_id'); // Recupera el user_id guardado

        if (!userId) {
          setError('No se pudo encontrar el ID del usuario. Intenta iniciar sesión de nuevo.');
          return;
        }

        // Solicita las reservas usando el user_id
        const response = await axiosHotelsInstance.get(`/users/${userId}/reservations`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        setReservations(response.data || []);
      } catch (err) {
        setError('Error al cargar reservas. Intenta nuevamente.');
      }
    };

    fetchReservations();
  }, []);

  return (
    <div className={styles.reservationsContainer}>
      <h1 className={styles.heading}>Mis Reservas</h1>
      {error && <p className={styles.error}>{error}</p>}
      {reservations.length === 0 ? (
        <p className={styles.noReservations}>No tienes reservas en este momento.</p>
      ) : (
        <div className={styles.reservationList}>
          {reservations.map((reservation) => (
            <div key={reservation.id} className={styles.reservationCard}>
              <h2 className={styles.hotelName}>{reservation.hotelName}</h2>
              <p className={styles.reservationDate}>Fecha de inicio: {reservation.start_date}</p>
              <p className={styles.reservationDate}>Fecha de fin: {reservation.end_date}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Reservas;
