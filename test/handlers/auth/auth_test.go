package auth_test

//import (
//	"github.com/gofiber/fiber/v2"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"net/http"
//	"net/http/httptest"
//	"sad/internal/handlers/auth"
//	authModels "sad/internal/models/auth"
//	errorsModels "sad/internal/models/errors"
//	usersModels "sad/internal/models/users"
//	mocks "sad/test/services/auth"
//	"strings"
//	"testing"
//)
//
//type MockAuthService struct {
//	mock.Mock
//}
//
//func (m *MockAuthService) Login(c *fiber.Ctx, user usersModels.UserCredentials) (string, error) {
//	args := m.Called(c, user)
//	return args.String(0), args.Error(1)
//}
//
//func (m *MockAuthService) Register(c *fiber.Ctx, user authModels.UserRegistrationRequest) (string, error) {
//	args := m.Called(c, user)
//	return args.String(0), args.Error(1)
//}
//
//func TestLogin(t *testing.T) {
//	// Создаем мок для AuthService
//	mockAuthService := new(mocks.AuthService)
//
//	// Создаем экземпляр authHandler с моком AuthService
//	authHandler := auth.NewAuthHandler(mockAuthService)
//
//	// Тестовые кейсы
//	testCases := []struct {
//		name           string
//		body           string
//		mockReturnArgs []interface{}
//		expectedStatus int
//		expectedJSON   string
//	}{
//		{
//			name:           "Success",
//			body:           `{"login": "test", "password": "password"}`,
//			mockReturnArgs: []interface{}{"token", nil},
//			expectedStatus: http.StatusOK,
//			expectedJSON:   `{"token": "token"}`,
//		},
//		{
//			name:           "InvalidRequestBody",
//			body:           `invalid`,
//			mockReturnArgs: []interface{}{"", nil},
//			expectedStatus: http.StatusBadRequest,
//			expectedJSON:   `{"error": "invalid request body"}`,
//		},
//		{
//			name:           "UserNotFound",
//			body:           `{"login": "test", "password": "password"}`,
//			mockReturnArgs: []interface{}{"", errorsModels.ErrUserNotFound},
//			expectedStatus: http.StatusNotFound,
//			expectedJSON:   `{"error": "User with this login does not exist"}`,
//		},
//		// Остальные тестовые кейсы...
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			app := fiber.New()
//
//			// Настраиваем ожидаемое поведение мока AuthService
//			mockAuthService.On("Login", mock.Anything, mock.Anything).Return(tc.mockReturnArgs...)
//
//			// Вызываем метод Login
//			err := authHandler.Login(ctx)
//
//			// Проверяем ошибку
//			assert.NoError(t, err)
//
//			// Проверяем статус ответа
//			assert.Equal(t, tc.expectedStatus, res.Code)
//
//			// Проверяем тело ответа
//			assert.JSONEq(t, tc.expectedJSON, res.Body.String())
//
//			// Проверяем, что мок AuthService был вызван с ожидаемыми аргументами
//			mockAuthService.AssertExpectations(t)
//		})
//	}
//}
//
//func TestRegister(t *testing.T) {
//	// Создаем мок для AuthService
//	mockAuthService := new(mocks.AuthService)
//
//	// Создаем экземпляр authHandler с моком AuthService
//	authHandler := auth.NewAuthHandler(mockAuthService)
//
//	// Тестовые кейсы
//	testCases := []struct {
//		name           string
//		body           string
//		mockReturnArgs []interface{}
//		expectedStatus int
//		expectedJSON   string
//	}{
//		{
//			name:           "Success",
//			body:           `{"login": "test", "password": "password"}`,
//			mockReturnArgs: []interface{}{"uuid", nil},
//			expectedStatus: http.StatusOK,
//			expectedJSON:   `{"uuid": "uuid"}`,
//		},
//		{
//			name:           "InvalidRequestBody",
//			body:           invalid,
//			mockReturnArgs: []interface{}{"", nil},
//			expectedStatus: http.StatusBadRequest,
//			expectedJSON:   `{"error": "invalid request body"}`,
//		},
//		{
//			name:           "UserExists",
//			body:           `{"login": "test", "password": "password"}`,
//			mockReturnArgs: []interface{}{"", errorsModels.ErrUserExists},
//			expectedStatus: http.StatusConflict,
//			expectedJSON:   `{"error": "User with this login already exist"}`,
//		},
//		// Остальные тестовые кейсы...
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// Создаем запрос
//			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tc.body))
//			req.Header.Set("Content-Type", "application/json")
//
//			// Создаем ответ
//			res := httptest.NewRecorder()
//
//			// Создаем контекст Fiber
//			ctx := fiber.New().AcquireCtx(req)
//			defer ctx.ReleaseCtx()
//
//			// Настраиваем ожидаемое поведение мока AuthService
//			mockAuthService.On("Register", mock.Anything, mock.Anything).Return(tc.mockReturnArgs...)
//
//			// Вызываем метод Register
//			err := authHandler.Register(ctx)
//
//			// Проверяем ошибку
//			assert.NoError(t, err)
//
//			// Проверяем статус ответа
//			assert.Equal(t, tc.expectedStatus, res.Code)
//
//			// Проверяем тело ответа
//			assert.JSONEq(t, tc.expectedJSON, res.Body.String())
//
//			// Проверяем, что мок AuthService был вызван с ожидаемыми аргументами
//			mockAuthService.AssertExpectations(t)
//		})
//	}
//}
