import axios from 'axios';

// Axios instance for USERS_API
const axiosUsersInstance = axios.create({
  baseURL: process.env.REACT_APP_USERS_API || 'http://localhost:8080', // Configuración para desarrollo local
  headers: {
    'Content-Type': 'application/json',
  },
});

// Axios instance for HOTELS_API
const axiosHotelsInstance = axios.create({
  baseURL: process.env.REACT_APP_HOTELS_API || 'http://localhost:8081', // Configuración para desarrollo local
  headers: {
    'Content-Type': 'application/json',
  },
});

export { axiosUsersInstance, axiosHotelsInstance };