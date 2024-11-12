import React, { useEffect, useState } from 'react';
import styles from './Reservas.module.css';

const Reservas = () => {
  const [reservations, setReservations] = useState([]);

  useEffect(() => {
    const fetchReservations = async () => {
      const userId = localStorage.getItem("userId"); // Recupera el ID de usuario almacenado
      if (!userId) {
        console.error("Usuario no autenticado.");
        return;
      }

      try {
        const response = await fetch(`http://localhost:8081/users/${userId}/reservations`);

        if (!response.ok) {  // Comprobar si la respuesta es exitosa
          throw new Error(`Error al obtener las reservas: ${response.statusText}`);
        }

        const data = await response.json();
        setReservations(data);
      } catch (error) {
        console.error("Error en la solicitud de reservas:", error);
      }
    };

    fetchReservations();
  }, []);



  return (
      <div className={styles.reservationsContainer}>
        <h1 className={styles.heading}>Mis Reservas</h1>

        {reservations.length === 0 ? (
            <p className={styles.noReservations}>No tienes reservas en este momento.</p>
        ) : (
            <div className={styles.reservationList}>
              {reservations.map((reservation) => (
                  <div key={reservation.id} className={styles.reservationCard}>
                    <h2 className={styles.hotelName}>{reservation.hotelName}</h2>
                    <p className={styles.reservationDate}>Fecha: {reservation.date}</p>
                    <p className={`${styles.reservationStatus} ${styles[reservation.status.toLowerCase()]}`}>
                      Estado: {reservation.status}
                    </p>
                  </div>
              ))}
            </div>
        )}
      </div>
  );
};

export default Reservas;
