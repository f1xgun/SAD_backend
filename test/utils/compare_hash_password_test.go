package utils

import (
	"github.com/stretchr/testify/assert"
	"sad/internal/utils"
	"testing"
)

func TestCompareHashPassword(t *testing.T) {
	// Тестовый пароль и его хеш
	password := "testpassword"
	hash := "$2a$14$Pr1H5rLHX00wBmst678HxuxbgSaTyQbMzTXwWt2dIB0PrIuSg6BYe"

	// Проверяем, что функция CompareHashPassword возвращает true для правильного пароля
	assert.True(t, utils.CompareHashPassword(password, hash))

	// Проверяем, что функция CompareHashPassword возвращает false для неправильного пароля
	assert.False(t, utils.CompareHashPassword("wrongpassword", hash))
}
