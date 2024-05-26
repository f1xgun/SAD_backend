package utils

import (
	"github.com/stretchr/testify/assert"
	"sad/internal/utils"
	"testing"
)

func TestGenerateHashPassword(t *testing.T) {
	// Тестовый пароль
	password := "testpassword"

	// Вызываем функцию GenerateHashPassword для генерации хеша пароля
	hash, err := utils.GenerateHashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Проверяем, что сгенерированный хеш соответствует паролю
	assert.True(t, utils.CompareHashPassword(password, hash))
}
