package users_test

import (
	"errors"
	"testing"
	dao "users-api/dao/users"
	domain "users-api/domain/users"
	"users-api/internal/tokenizers"
	repositories "users-api/repositories/users"
	service "users-api/services/users"

	"github.com/stretchr/testify/assert"
)

var (
	mainRepo      = repositories.NewMock()
	cacheRepo     = repositories.NewMock()
	memcachedRepo = repositories.NewMock()
	tokenizer     = tokenizers.NewMock()
	usersService  = service.NewService(mainRepo, cacheRepo, memcachedRepo, tokenizer)
)

func TestService(t *testing.T) {
	t.Run("GetAll - Success", func(t *testing.T) {
		mockUsers := []dao.User{
			{ID: 1, Username: "user1", Password: "password1"},
			{ID: 2, Username: "user2", Password: "password2"},
		}
		mainRepo.On("GetAll").Return(mockUsers, nil).Once()

		result, err := usersService.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "user1", result[0].Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Error", func(t *testing.T) {
		mainRepo.On("GetAll").Return(nil, errors.New("db error")).Once()

		result, err := usersService.GetAll()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "error getting all users: db error", err.Error()) // Assert expectations    // Assert expectations    // Assert expectations

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Success from Cache", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Not Found in Cache, Found in Memcached", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Not Found in Cache or Memcached, Found in Main Repo", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error in Main Repo", func(t *testing.T) {
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("db error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by ID: db error", err.Error())
		assert.Equal(t, domain.User{}, result)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error caching user after retrieval", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}

		// Simula que el usuario no está en cache ni en memcached.
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		// Simula un error al guardar en el cache.
		cacheRepo.On("Create", mockUser).Return(int64(0), errors.New("cache save error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error caching user after main retrieval: cache save error", err.Error())
		assert.Equal(t, domain.User{}, result)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Success", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}
		mainRepo.On("Create", newUser).Return(int64(1), nil).Once()
		newUser.ID = 1
		cacheRepo.On("Create", newUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", newUser).Return(int64(1), nil).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Error", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}
		mainRepo.On("Create", newUser).Return(int64(0), errors.New("db error")).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "error creating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Error saving in Memcached", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}
		mainRepo.On("Create", newUser).Return(int64(1), nil).Once()

		newUser.ID = 1
		cacheRepo.On("Create", newUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", newUser).Return(int64(0), errors.New("memcached save error")).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.Error(t, err)
		assert.Equal(t, "error saving new user in memcached: memcached save error", err.Error())
		assert.Equal(t, int64(0), id)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Update - Success", func(t *testing.T) {
		updateUser := dao.User{ID: 1, Username: "updateduser", Password: service.Hash("newpassword")}
		mainRepo.On("Update", updateUser).Return(nil).Once()
		cacheRepo.On("Update", updateUser).Return(nil).Once()
		memcachedRepo.On("Update", updateUser).Return(nil).Once()

		userToUpdate := domain.User{ID: 1, Username: "updateduser", Password: "newpassword"}
		err := usersService.Update(userToUpdate)

		assert.NoError(t, err)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Update - Error", func(t *testing.T) {
		updateUser := dao.User{ID: 1, Username: "updateduser", Password: service.Hash("newpassword")}
		mainRepo.On("Update", updateUser).Return(errors.New("db error")).Once()

		userToUpdate := domain.User{ID: 1, Username: "updateduser", Password: "newpassword"}
		err := usersService.Update(userToUpdate)

		assert.Error(t, err)
		assert.Equal(t, "error updating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})
	t.Run("Update - Error updating in cache", func(t *testing.T) {
		updateUser := dao.User{ID: 1, Username: "updateduser", Password: service.Hash("newpassword")}
		mainRepo.On("Update", updateUser).Return(nil).Once()
		cacheRepo.On("Update", updateUser).Return(errors.New("cache update error")).Once()

		userToUpdate := domain.User{ID: 1, Username: "updateduser", Password: "newpassword"}
		err := usersService.Update(userToUpdate)

		assert.Error(t, err)
		assert.Equal(t, "error updating user in cache: cache update error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(nil).Once()
		cacheRepo.On("Delete", int64(1)).Return(nil).Once()
		memcachedRepo.On("Delete", int64(1)).Return(nil).Once()

		err := usersService.Delete(1)

		assert.NoError(t, err)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Delete - Error", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(errors.New("db error")).Once()

		err := usersService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "error deleting user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Success", func(t *testing.T) {
		username := "user1"
		password := "password"
		hashedPassword := service.Hash(password)

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword, Tipo: "cliente"}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()
		tokenizer.On("GenerateToken", username, int64(1), "cliente").Return("token", nil).Once() // Añade el tercer argumento

		response, err := usersService.Login(username, password)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.UserID)
		assert.Equal(t, "token", response.Token)

		cacheRepo.AssertExpectations(t)
		tokenizer.AssertExpectations(t)
	})

	t.Run("Login - Invalid Credentials", func(t *testing.T) {
		username := "user1"
		password := "wrongpassword"
		hashedPassword := service.Hash("password")

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - User Not Found", func(t *testing.T) {
		username := "user1"
		password := "password"

		cacheRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by username from main repository: not found", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})
	/*
		t.Run("Login - Token Generation Error", func(t *testing.T) {
			username := "user1"
			password := "password"
			hashedPassword := service.Hash(password)

			mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
			cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()
			tokenizer.On("GenerateToken", username, int64(1)).Return("", errors.New("token error")).Once()

			response, err := usersService.Login(username, password)

			assert.Error(t, err)
			assert.Equal(t, "error generating token: token error", err.Error())
			assert.Equal(t, domain.LoginResponse{}, response)

			mainRepo.AssertExpectations(t)
			cacheRepo.AssertExpectations(t)
			memcachedRepo.AssertExpectations(t)
		})*/
	t.Run("Login - Error getting user in all repositories", func(t *testing.T) {
		username := "user1"
		password := "password"

		cacheRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		memcachedRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("not found")).Once()
		mainRepo.On("GetByUsername", username).Return(dao.User{}, errors.New("db error")).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by username from main repository: db error", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error saving in Memcached after main repo retrieval", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}

		// Simula que no está en caché ni en memcached.
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		// Simula éxito al guardar en caché, pero error en memcached.
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", mockUser).Return(int64(0), errors.New("memcached save error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error saving user in memcached: memcached save error", err.Error())
		assert.Equal(t, domain.User{}, result)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Create - Empty Username and Password", func(t *testing.T) {
		user := domain.User{Username: "", Password: ""}

		id, err := usersService.Create(user)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "invalid input: username and password cannot be empty", err.Error())
	})

	t.Run("Update - Empty Username", func(t *testing.T) {
		user := domain.User{ID: 1, Username: "", Password: "password"}

		err := usersService.Update(user)

		assert.Error(t, err)
		assert.Equal(t, "invalid input: username cannot be empty", err.Error())
	})

	t.Run("Delete - Error deleting from cache", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(nil).Once()
		cacheRepo.On("Delete", int64(1)).Return(errors.New("cache delete error")).Once()
		memcachedRepo.On("Delete", int64(1)).Return(nil).Once() // Configura la llamada de Memcached

		err := usersService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "error deleting user from cache: cache delete error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	/*t.Run("Login - Token Generation Error", func(t *testing.T) {
		username := "user1"
		password := "password"
		hashedPassword := service.Hash(password)

		mockUser := dao.User{ID: 1, Username: username, Password: hashedPassword}
		cacheRepo.On("GetByUsername", username).Return(mockUser, nil).Once()
		tokenizer.On("GenerateToken", username, int64(1)).Return("", errors.New("token error")).Once()

		response, err := usersService.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "error generating token: token error", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})*/

	t.Run("GetByUsername - Error in Cache and Memcached", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}

		// Simula errores en ambos repositorios secundarios
		cacheRepo.On("GetByUsername", "user1").Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByUsername", "user1").Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByUsername", "user1").Return(mockUser, nil).Once()

		// Permite las llamadas esperadas a Create
		cacheRepo.On("Create", mockUser).Return(int64(1), nil).Times(1)     // Cuando se guarda tras obtener de mainRepository
		memcachedRepo.On("Create", mockUser).Return(int64(1), nil).Times(1) // Cuando se guarda tras obtener de mainRepository

		// Ejecuta el método de servicio
		result, err := usersService.GetByUsername("user1")

		// Validaciones
		assert.NoError(t, err)
		assert.Equal(t, "user1", result.Username)

		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
		mainRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error in Cache and Memcached Save", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}

		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		// Simula error al guardar en caché
		cacheRepo.On("Create", mockUser).Return(int64(0), errors.New("cache save error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error caching user after main retrieval: cache save error", err.Error())
		assert.Equal(t, domain.User{}, result)

		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
		mainRepo.AssertExpectations(t)
	})

	t.Run("Create - Error Saving in Memcached", func(t *testing.T) {
		newUser := dao.User{Username: "newuser", Password: service.Hash("password")}

		mainRepo.On("Create", newUser).Return(int64(1), nil).Once()
		newUser.ID = 1
		cacheRepo.On("Create", newUser).Return(int64(1), nil).Once()
		memcachedRepo.On("Create", newUser).Return(int64(0), errors.New("memcached save error")).Once()

		id, err := usersService.Create(domain.User{Username: "newuser", Password: "password"})

		assert.Error(t, err)
		assert.Equal(t, "error saving new user in memcached: memcached save error", err.Error())
		assert.Equal(t, int64(0), id)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})
	t.Run("Delete - Error Deleting from Memcached", func(t *testing.T) {
		mainRepo.On("Delete", int64(1)).Return(nil).Once()
		cacheRepo.On("Delete", int64(1)).Return(nil).Once()
		memcachedRepo.On("Delete", int64(1)).Return(errors.New("memcached delete error")).Once()

		err := usersService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "error deleting user from memcached: memcached delete error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error saving in Cache after Retrieval", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: "password1"}

		// Simula que no se encuentra en cache ni memcached
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByID", int64(1)).Return(mockUser, nil).Once()

		// Simula error al guardar en caché
		cacheRepo.On("Create", mockUser).Return(int64(0), errors.New("cache save error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error caching user after main retrieval: cache save error", err.Error())
		assert.Equal(t, domain.User{}, result)

		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
		mainRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Error retrieving from all sources", func(t *testing.T) {
		// Simula errores en todos los repositorios
		cacheRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("memcached miss")).Once()
		mainRepo.On("GetByID", int64(1)).Return(dao.User{}, errors.New("db error")).Once()

		result, err := usersService.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by ID: db error", err.Error())
		assert.Equal(t, domain.User{}, result)

		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
		mainRepo.AssertExpectations(t)
	})

	t.Run("Update - Invalid User ID", func(t *testing.T) {
		user := domain.User{ID: 0, Username: "user1", Password: "password"}

		err := usersService.Update(user)

		assert.Error(t, err)
		assert.Equal(t, "invalid input: user ID cannot be zero", err.Error())
	})

	t.Run("Update - Error Updating Main Repository", func(t *testing.T) {
		mockUser := dao.User{ID: 1, Username: "user1", Password: service.Hash("password")}

		mainRepo.On("Update", mockUser).Return(errors.New("db error")).Once()

		err := usersService.Update(domain.User{ID: 1, Username: "user1", Password: "password"})

		assert.Error(t, err)
		assert.Equal(t, "error updating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
	})
	t.Run("Delete - Error in Main Repository", func(t *testing.T) {
		mainRepo = repositories.NewMock()
		cacheRepo = repositories.NewMock()
		memcachedRepo = repositories.NewMock()
		tokenizer = tokenizers.NewMock()
		usersService = service.NewService(mainRepo, cacheRepo, memcachedRepo, tokenizer)

		mainRepo.On("Delete", int64(1)).Return(errors.New("db error")).Once()

		err := usersService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "error deleting user: db error", err.Error())

		cacheRepo.AssertNotCalled(t, "Delete", int64(1))
		memcachedRepo.AssertNotCalled(t, "Delete", int64(1))

		mainRepo.AssertExpectations(t)
	})

	t.Run("Delete - Error in Cache and Memcached", func(t *testing.T) {
		// Configura los mocks para manejar múltiples llamadas
		mainRepo.On("Delete", int64(1)).Return(nil).Once()
		cacheRepo.On("Delete", int64(1)).Return(errors.New("cache delete error")).Once()
		memcachedRepo.On("Delete", int64(1)).Return(errors.New("memcached delete error")).Once()

		// Ejecuta el método
		err := usersService.Delete(1)

		// Valida el resultado
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error deleting user from cache: cache delete error")
		assert.Contains(t, err.Error(), "error deleting user from memcached: memcached delete error")

		// Verifica las expectativas
		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Invalid Password", func(t *testing.T) {
		hashedPassword := service.Hash("password")
		mockUser := dao.User{ID: 1, Username: "user1", Password: hashedPassword}

		// Configura los mocks
		cacheRepo.On("GetByUsername", "user1").Return(mockUser, nil).Once()

		// Ejecuta el método
		response, err := usersService.Login("user1", "wrongpassword")

		// Valida el resultado
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		// Verifica las expectativas
		cacheRepo.AssertExpectations(t)
	})

	t.Run("Login - Token Generation Error", func(t *testing.T) {
		// Resetea mocks
		mainRepo = repositories.NewMock()
		cacheRepo = repositories.NewMock()
		memcachedRepo = repositories.NewMock()
		tokenizer = tokenizers.NewMock()
		usersService = service.NewService(mainRepo, cacheRepo, memcachedRepo, tokenizer)

		// Configurar mocks
		hashedPassword := service.Hash("password")
		mockUser := dao.User{ID: 1, Username: "user1", Password: hashedPassword}

		cacheRepo.On("GetByUsername", "user1").Return(mockUser, nil).Once()
		tokenizer.On("GenerateToken", "user1", int64(1), "cliente").Return("", errors.New("token error")).Once() // Agregar el tercer argumento

		// Ejecutar
		response, err := usersService.Login("user1", "password")

		// Validar
		assert.Error(t, err)
		assert.Equal(t, "error generating token: token error", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		// Verificar expectativas
		cacheRepo.AssertExpectations(t)
		tokenizer.AssertExpectations(t)
	})

	t.Run("Login - Save to Cache and Memcached Fails", func(t *testing.T) {
		// Configurar mocks
		hashedPassword := service.Hash("password")
		mockUser := dao.User{ID: 1, Username: "user1", Password: hashedPassword}

		cacheRepo.On("GetByUsername", "user1").Return(dao.User{}, errors.New("cache miss")).Once()
		memcachedRepo.On("GetByUsername", "user1").Return(mockUser, nil).Once()
		cacheRepo.On("Create", mockUser).Return(int64(0), errors.New("cache save error")).Once()
		tokenizer.On("GenerateToken", "user1", int64(1), "cliente").Return("token", nil).Once()

		// Ejecutar
		response, err := usersService.Login("user1", "password")

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, "user1", response.Username)
		assert.Equal(t, "token", response.Token)

		// Verificar expectativas
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
		tokenizer.AssertExpectations(t)
	})

}
