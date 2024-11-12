import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_USERS_API || 'http://localhost:8080', // Configuración para desarrollo local
  headers: {
    'Content-Type': 'application/json',
  },
});

export default axiosInstance;
