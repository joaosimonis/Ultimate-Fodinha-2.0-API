package auth

import "errors"

var ErrInvalidCredentials = errors.New("credenciais inválidas")
var ErrUserNotFound = errors.New("usuário não encontrado")
var ErrInvalidToken = errors.New("token inválido")
var ErrExpiredToken = errors.New("token expirado")
