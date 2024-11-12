// App.js
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import NavBar from "./components/NavBar";
import Home from "./Paginas/Home";
import Hoteles from './Paginas/Hoteles';
import Login from './Paginas/Login'; // Asegúrate de que esta importación esté correcta
import Registro from './Paginas/IniciarSesion';
import MisReservas from './Paginas/Reservas';
import DetalleHotel from './Paginas/DetalleHotel';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsAuthenticated(!!token);
  }, []);

  const handleLogin = () => {
    setIsAuthenticated(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
  };

  return (
    <div className="App">
      <Router>
        <NavBar isAuthenticated={isAuthenticated} onLogout={handleLogout} />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/hoteles" element={<Hoteles />} />
          <Route path="/login" element={<Login onLogin={handleLogin} />} /> {/* Asegúrate de pasar onLogin aquí */}
          <Route path="/registro" element={<Registro />} />
          <Route 
            path="/mis-reservas" 
            element={isAuthenticated ? <MisReservas /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/detalle-hotel/:id" 
            element={<DetalleHotel isAuthenticated={isAuthenticated} />} 
          />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
