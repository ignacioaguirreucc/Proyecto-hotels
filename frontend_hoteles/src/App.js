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
import Admin from './Paginas/Admin';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsAuthenticated(!!token);
  }, []);

  const handleLogin = (tipo) => {
    setIsAuthenticated(true);
    // Guarda el tipo de usuario en el estado para usarlo en NavBar y otras partes
    localStorage.setItem('user_tipo', tipo);
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
          <Route path="/home" element={<Home />} />

          <Route exact path="/admin" element={<Admin />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;