// NavBar.js
import React from 'react';
import { Link } from 'react-router-dom';
import styles from './NavBar.module.css';

const NavBar = ({ isAuthenticated, onLogout }) => {
  return (
    <nav className={styles.navbar}>
      <ul className={styles.navlinks}>
        <li><Link to="/home">Hogar</Link></li>
        <li><Link to="/hoteles">Resultados</Link></li>
        {isAuthenticated ? (
          <>
            <li><Link to="/mis-reservas">Mis Reservas</Link></li>
            <li>
              <button className={styles.navlinkButton} onClick={onLogout}>
                Cerrar Sesión
              </button>
            </li>
          </>
        ) : (
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