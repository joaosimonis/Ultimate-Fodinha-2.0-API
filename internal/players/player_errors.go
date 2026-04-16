package players

import "errors"

var ErrNotFound = errors.New("jogador não encontrado") // Erro 404
var ErrNameTaken = errors.New("já existe um jogador com esse nome") // Erro 409