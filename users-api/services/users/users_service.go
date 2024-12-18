package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	dao "users-api/dao/users"
	domain "users-api/domain/users"
	//"github.com/golang-jwt/jwt/v5" Ensure this import is present
)

type Repository interface {
	GetAll() ([]dao.User, error)
	GetByID(id int64) (dao.User, error)
	GetByUsername(username string) (dao.User, error)
	Create(user dao.User) (int64, error)
	Update(user dao.User) error
	Delete(id int64) error
}

type Tokenizer interface {
	GenerateToken(username string, userID int64, userType string) (string, error) // Actualiza la firma
}

type Service struct {
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
	tokenizer           Tokenizer
}

func NewService(mainRepository, cacheRepository, memcachedRepository Repository, tokenizer Tokenizer) Service {
	return Service{
		mainRepository:      mainRepository,
		cacheRepository:     cacheRepository,
		memcachedRepository: memcachedRepository,
		tokenizer:           tokenizer,
	}
}

func (service Service) GetAll() ([]domain.User, error) {
	users, err := service.mainRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	result := make([]domain.User, 0)
	for _, user := range users {
		result = append(result, domain.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		})
	}

	return result, nil
}

func (service Service) GetByID(id int64) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByID(id)
	if err == nil {
		return service.convertUser(user), nil
	}

	// Check in memcached
	user, err = service.memcachedRepository.GetByID(id)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}

	// Check in main repository
	user, err = service.mainRepository.GetByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) GetByUsername(username string) (domain.User, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetByUsername(username)
	if err == nil {
		return service.convertUser(user), nil
	}

	// Check memcached
	user, err = service.memcachedRepository.GetByUsername(username)
	if err == nil {
		if _, err := service.cacheRepository.Create(user); err != nil {
			return domain.User{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.convertUser(user), nil
	}

	// Check main repository
	user, err = service.mainRepository.GetByUsername(username)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.Create(user); err != nil {
		return domain.User{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.convertUser(user), nil
}

func (service Service) Create(user domain.User) (int64, error) {
	// Validación de entrada
	if user.Username == "" || user.Password == "" {
		return 0, fmt.Errorf("invalid input: username and password cannot be empty")
	}

	// Hash the password
	passwordHash := Hash(user.Password)

	newUser := dao.User{
		Username: user.Username,
		Password: passwordHash,
	}

	// Create in main repository
	id, err := service.mainRepository.Create(newUser)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	// Add to cache and memcached
	newUser.ID = id
	if _, err := service.cacheRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error caching new user: %w", err)
	}
	if _, err := service.memcachedRepository.Create(newUser); err != nil {
		return 0, fmt.Errorf("error saving new user in memcached: %w", err)
	}

	return id, nil
}

func (service Service) Update(user domain.User) error {
	// Validaciones de entrada
	if user.Username == "" {
		return fmt.Errorf("invalid input: username cannot be empty")
	}
	if user.ID == 0 {
		return fmt.Errorf("invalid input: user ID cannot be zero")
	}

	// Hash the password if provided
	var passwordHash string
	if user.Password != "" {
		passwordHash = Hash(user.Password)
	} else {
		existingUser, err := service.mainRepository.GetByID(user.ID)
		if err != nil {
			return fmt.Errorf("error retrieving existing user: %w", err)
		}
		passwordHash = existingUser.Password
	}

	// Update in main repository
	err := service.mainRepository.Update(dao.User{
		ID:       user.ID,
		Username: user.Username,
		Password: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Update in cache and memcached
	updatedUser := dao.User{
		ID:       user.ID,
		Username: user.Username,
		Password: passwordHash,
	}
	if err := service.cacheRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in cache: %w", err)
	}
	if err := service.memcachedRepository.Update(updatedUser); err != nil {
		return fmt.Errorf("error updating user in memcached: %w", err)
	}

	return nil
}

func (service Service) Delete(id int64) error {
	// Intenta borrar en el repositorio principal
	if err := service.mainRepository.Delete(id); err != nil {
		return fmt.Errorf("error deleting user: %w", err) // Detiene la ejecución en caso de error
	}

	// Inicializa los errores secundarios
	var errs []string

	// Intenta borrar en el caché
	if err := service.cacheRepository.Delete(id); err != nil {
		errs = append(errs, fmt.Sprintf("error deleting user from cache: %s", err.Error()))
	}

	// Intenta borrar en Memcached
	if err := service.memcachedRepository.Delete(id); err != nil {
		errs = append(errs, fmt.Sprintf("error deleting user from memcached: %s", err.Error()))
	}

	// Combina los errores secundarios si existen
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}

	return nil
}

func (service Service) Login(username string, password string) (domain.LoginResponse, error) {
	// Hash del password ingresado
	passwordHash := Hash(password)

	var user dao.User
	var err error

	// Busca en caché primero
	user, err = service.cacheRepository.GetByUsername(username)
	if err != nil {
		// Si no está en caché, busca en Memcached
		user, err = service.memcachedRepository.GetByUsername(username)
		if err != nil {
			// Si tampoco está en Memcached, busca en el repositorio principal
			user, err = service.mainRepository.GetByUsername(username)
			if err != nil {
				return domain.LoginResponse{}, fmt.Errorf("error getting user by username from main repository: %w", err)
			}

			// Guarda en caché y Memcached
			if _, err := service.cacheRepository.Create(user); err != nil {
				// Manejo del error de guardado en caché (opcional)
			}
			if _, err := service.memcachedRepository.Create(user); err != nil {
				// Manejo del error de guardado en Memcached (opcional)
			}
		} else {
			// Guarda en caché si se encontró en Memcached
			if _, err := service.cacheRepository.Create(user); err != nil {
				// Manejo del error de guardado en caché (opcional)
			}
		}
	}

	// Compara contraseñas
	if user.Password != passwordHash {
		return domain.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// Establece el tipo de usuario por defecto si no está definido
	if user.Tipo == "" {
		user.Tipo = "cliente"
	}

	// Genera el token
	token, err := service.tokenizer.GenerateToken(user.Username, user.ID, user.Tipo)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("error generating token: %w", err)
	}

	// Retorna la respuesta de login
	return domain.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Token:    token,
		Tipo:     user.Tipo,
	}, nil
}

func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (service Service) convertUser(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
}
