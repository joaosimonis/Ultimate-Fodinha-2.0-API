package matches

import "errors"

var ErrNotFound = errors.New("partida não encontrada")
var ErrRoomCodeTaken = errors.New("já existe uma partida com esse código de sala")
var ErrInvalidMode = errors.New("modo de partida inválido")
var ErrMaxScoreRequired = errors.New("pontuação máxima é obrigatória")
var ErrInvalidMaxScore = errors.New("pontuação máxima deve ser maior que zero")
var ErrMaxScoreNotAllowed = errors.New("campo max_score inválido para o modo fodinha")
var ErrMaxScoreValueNotAllowed = errors.New("pontuação máxima não permitida para truco")
var ErrLivesRequired = errors.New("vidas são obrigatórias")
var ErrInvalidLives = errors.New("vidas devem ser maior que zero")
var ErrLivesNotAllowed = errors.New("campo lives inválido para o modo truco")
var ErrLivesValueNotAllowed = errors.New("vidas não permitidas para fodinha")
