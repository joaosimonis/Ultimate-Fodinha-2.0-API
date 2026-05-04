package ranking

import "errors"

var ErrModeRequired = errors.New("modo do ranking é obrigatório")
var ErrInvalidMode = errors.New("modo do ranking inválido")
var ErrTeamSizeRequired = errors.New("team_size é obrigatório para ranking de truco")
var ErrTeamSizeNotAllowed = errors.New("team_size não é permitido para ranking de fodinha")
var ErrInvalidTeamSize = errors.New("team_size deve ser 2 para dupla ou 3 para trio")
var ErrInvalidLimit = errors.New("limit deve ser maior que zero")
