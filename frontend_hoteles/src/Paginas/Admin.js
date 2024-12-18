import React, { useState, useEffect } from 'react';
import { axiosHotelsInstance, axiosSearchInstance } from '../axiosConfig';
import styles from './Admin.module.css';
import axios from 'axios';


const Admin = () => {
  const [hotels, setHotels] = useState([]);
  const [containerCounts, setContainerCounts] = useState({}); // Estado para almacenar las instancias por imagen
  const [formData, setFormData] = useState({
    name: '',
    address: '',
    city: '',
    state: '',
    rating: '',
    amenities: '',
    descripcion: '',
  });
  const [showModal, setShowModal] = useState(false);

  // Obtener la lista de hoteles
  const fetchHotels = async (query = '*') => {
    try {
      // Realiza una solicitud GET a /search para obtener la lista de hoteles
      const response = await axiosSearchInstance.get('/search', {
        params: { q: query, offset: 0, limit: 10 }  // Consulta con el parámetro de búsqueda
      });

      // Mapea los resultados de Solr a los atributos del hotel en frontend
      const hotelsData = response.data.map((hotel) => ({
        id: hotel.id,
        name: Array.isArray(hotel.name) ? hotel.name[0] : hotel.name,
        state: Array.isArray(hotel.state) ? hotel.state[0] : hotel.state,
        rating: Array.isArray(hotel.rating) ? hotel.rating[0] : hotel.rating,
        amenities: Array.isArray(hotel.amenities) ? hotel.amenities.join(", ") : 'No disponible',
        descripcion: Array.isArray(hotel.descripcion) ? hotel.descripcion[0] : 'No disponible',
        city: Array.isArray(hotel.city) ? hotel.city[0] : hotel.city,
        address: Array.isArray(hotel.address) ? hotel.address[0] : hotel.address,
      }));

      setHotels(hotelsData);  // Guarda los hoteles en el estado
    } catch (err) {
      console.error('Error al cargar los hoteles:', err);
      //setError('Hubo un problema al cargar los hoteles.');
    }
  };

  const fetchContainerCounts = async () => {
    try {
      const response = await axios.get('/containers/json');

      console.log('Contenedores desde Docker API:', response.data); // Debug
      const containers = response.data;
  
      const counts = containers.reduce((acc, container) => {
        const image = container.Image;
        acc[image] = (acc[image] || 0) + 1;
        return acc;
      }, {});
  
      console.log('Conteo de contenedores por imagen:', counts); // Debug
      setContainerCounts(counts);
    } catch (err) {
      console.error('Error al obtener contenedores de Docker:', err);
    }
  };
  

  useEffect(() => {
    fetchHotels();
    fetchContainerCounts();
  }, []);

  
// Crear o editar un hotel
// Crear o editar un hotel
const handleSubmit = async (e) => {
  e.preventDefault();

  const dataToSend = {
    name: formData.name,
    address: formData.address,
    city: formData.city,
    state: formData.state,
    amenities: formData.amenities.split(',').map((amenity) => amenity.trim()),
    descripcion: formData.descripcion.split(',').map((desc) => desc.trim()),
    rating: parseFloat(formData.rating),
  };

  try {
    if (formData.id) {
      // PUT request para actualizar un hotel existente usando formData.id
      await axiosHotelsInstance.put(`/hotels/${formData.id}`, dataToSend);
      console.log('Hotel actualizado correctamente');
    } else {
      // POST request para crear un hotel nuevo
      await axiosHotelsInstance.post('/hotels', dataToSend);
      console.log('Hotel creado correctamente');
    }

    await fetchHotels(); // Refrescar la lista de hoteles después de crear o editar
    closeModal(); // Cierra el modal después de guardar
  } catch (error) {
    console.error('Error al guardar hotel:', error);
    alert('Ocurrió un error al guardar el hotel. Inténtalo de nuevo.');
  }
};


  
  
  

  // Eliminar un hotel
const deleteHotel = async (id) => {
  try {
    await axiosHotelsInstance.delete(`/hotels/${id}`);
    // Actualizar la lista de hoteles eliminando el hotel con el id correspondiente
    setHotels((prevHotels) => prevHotels.filter((hotel) => hotel.id !== id));
  } catch (error) {
    console.error('Error al eliminar hotel:', error);
  }
};

  // Manejo del formulario
  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  // Abrir el modal para crear o editar un hotel
  const openModal = (hotel = null) => {
    setFormData(
      hotel
        ? {
            id: hotel.id || '', // Carga el id del hotel
            name: hotel.name || '',
            address: hotel.address || '',
            city: hotel.city || '',
            state: hotel.state || '',
            rating: hotel.rating || '',
            amenities: Array.isArray(hotel.amenities)
              ? hotel.amenities.join(', ')
              : hotel.amenities || '',
            descripcion: Array.isArray(hotel.descripcion)
              ? hotel.descripcion.join(', ')
              : hotel.descripcion || '',
          }
        : {
            id: '', // Vacía el id al crear un nuevo hotel
            name: '',
            address: '',
            city: '',
            state: '',
            rating: '',
            amenities: '',
            descripcion: '',
          }
    );
    setShowModal(true);
  };
  

  // Cerrar el modal
  const closeModal = () => {
    setShowModal(false);
    setFormData({
      id: '', // Incluye el id vacío al cerrar el modal
      name: '',
      address: '',
      city: '',
      state: '',
      rating: '',
      amenities: '',
      descripcion: '',
    });
  };
  

  return (
    <div className={styles.container}>
      <h1>Administrador de Hoteles</h1>
      <button onClick={() => openModal()} className={styles.createButton}>
        Crear Nuevo Hotel
      </button>

      <ul className={styles.hotelList}>
  {hotels.length > 0 ? (
    hotels.map((hotel) => (
      <li key={hotel.id} className={styles.hotelItem}>
        <div>
          <strong>{hotel.name}</strong>
          <p><strong>Dirección:</strong> {hotel.address}</p>
          <p><strong>Ciudad:</strong> {hotel.city}, {hotel.state}</p>
          <p><strong>Calificación:</strong> {hotel.rating}</p>
          <p><strong>Comodidades:</strong> {hotel.amenities}</p>
          <p><strong>Descripción:</strong> {hotel.descripcion}</p>
        </div>
        <div>
          <button onClick={() => openModal(hotel)} className={styles.editButton}>
            Editar
          </button>
          <button onClick={() => deleteHotel(hotel.id)} className={styles.deleteButton}>
            Eliminar
          </button>
        </div>
      </li>
    ))
  ) : (
    <p>No se encontraron hoteles.</p>
  )}
</ul>


<h2>Instancias de Docker</h2>
      <ul>
        {Object.entries(containerCounts).map(([image, count]) => (
          <li key={image}>
            <strong>{image}:</strong> {count} instancia(s)
          </li>
        ))}
      </ul>




  {/* Modal para Crear/Editar */}
      {showModal && (
        <div className={styles.modal}>
          <div className={styles.modalContent}>
          <h2>{formData.id ? 'Editar Hotel' : 'Crear Nuevo Hotel'}</h2>
            <form onSubmit={handleSubmit}>
  <input
    type="text"
    name="name"
    placeholder="Nombre"
    onChange={handleChange}
    value={formData.name}
    required
  />
  <input
    type="text"
    name="address"
    placeholder="Dirección"
    onChange={handleChange}
    value={formData.address}
    required
  />
  <input
    type="text"
    name="city"
    placeholder="Ciudad"
    onChange={handleChange}
    value={formData.city}
    required
  />
  <input
    type="text"
    name="state"
    placeholder="Estado"
    onChange={handleChange}
    value={formData.state}
    required
  />
  <input
    type="number"
    step="0.1"
    name="rating"
    placeholder="Calificación"
    onChange={handleChange}
    value={formData.rating}
    required
  />
  <input
    type="text"
    name="amenities"
    placeholder="Comodidades (separadas por comas)"
    onChange={handleChange}
    value={formData.amenities}
  />
  <textarea
    name="descripcion"
    placeholder="Descripción (separada por comas)"
    onChange={handleChange}
    value={formData.descripcion}
  />
  <div className={styles.modalActions}>
    <button type="submit" className={styles.saveButton}>
      {formData.id ? 'Guardar Cambios' : 'Crear Hotel'}
    </button>
    <button type="button" onClick={closeModal} className={styles.cancelButton}>
      Cancelar
    </button>
  </div>
</form>

          </div>
        </div>
      )}
    </div>
  );
};

export default Admin;
