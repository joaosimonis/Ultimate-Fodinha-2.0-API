package players

import "errors"

var ErrNotFound = errors.New("jogador não encontrado")
var ErrNameTaken = errors.New("já existe um jogador com esse nome")
