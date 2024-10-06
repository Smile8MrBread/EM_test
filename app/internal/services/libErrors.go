package services

import "errors"

var (
	ErrInvalidSquad = errors.New("invalid squad")
	ErrInvalidSong  = errors.New("invalid song")
	ErrInvalidId    = errors.New("invalid id")
	ErrInvalidText  = errors.New("invalid text")

	ErrInvalidChunk = errors.New("invalid chunk")
	ErrInvalidField = errors.New("invalid field")
	ErrInvalidOrder = errors.New("invalid order")
)
