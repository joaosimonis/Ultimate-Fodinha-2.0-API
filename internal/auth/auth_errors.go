package auth

import "errors"

var ErrInvalidCredentials = errors.New("credenciais inválidas")
var ErrUserNotFound = errors.New("usuário não encontrado")
