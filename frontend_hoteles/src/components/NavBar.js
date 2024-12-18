import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styles from './NavBar.module.css';

const NavBar = ({ isAuthenticated, onLogout }) => {
  const navigate = useNavigate(); // Hook para redirigir

  const handleLogout = () => {
    onLogout(); // Llamar la función pasada como prop
    navigate('/login'); // Redirigir a la página de login
    window.location.reload(); // Refrescar la página
  };

  // Recuperar el tipo de usuario desde localStorage
  const userType = localStorage.getItem('tipo'); // 'cliente' o 'administrador'

  return (
    <nav className={styles.navbar}>
      <ul className={styles.navlinks}>
        {isAuthenticated && userType === 'cliente' && (
          <>
            <li><Link to="/home">Hogar</Link></li>
            <li><Link to="/hoteles">Resultados</Link></li>
            <li><Link to="/mis-reservas">Mis Reservas</Link></li>
          </>
        )}
        {isAuthenticated && (
          <li>
            <button className={styles.navlinkButton} onClick={handleLogout}>
              Cerrar Sesión
            </button>
          </li>
        )}
        {!isAuthenticated && (
          <>
            <li><Link to="/login">Iniciar Sesión</Link></li>
            <li><Link to="/registro">Registro</Link></li>
          </>
        )}
      </ul>
    </nav>
  );
};

export default NavBar;
