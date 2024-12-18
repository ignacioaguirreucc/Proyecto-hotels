import axios from 'axios';

// Axios instance for USERS_API (ahora conectado a NGINX)
const axiosUsersInstance = axios.create({
  baseURL: process.env.REACT_APP_USERS_API || 'http://localhost:8085', // Configuración para desarrollo local
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
axiosHotelsInstance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

const axiosSearchInstance = axios.create({
  baseURL: process.env.REACT_APP_SEARCH_API || 'http://localhost:8082', // Configuración para desarrollo local
  headers: {
    'Content-Type': 'application/json',
  },
});

export { axiosUsersInstance, axiosHotelsInstance ,axiosSearchInstance};