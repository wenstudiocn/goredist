package e

import (
	"errors"
)

var (
	ErrKeyNotExists        = errors.New("key of map is not exists")
	ErrBlank               = errors.New("blank slice or array")
	ErrNotFound            = errors.New("target not exists")
	ErrUnknown             = errors.New("unknown error")
	ErrOutOfRange          = errors.New("index out of range")
	ErrParameters          = errors.New("bad parameters")
	ErrLength              = errors.New("invalid length")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrTimeout             = errors.New("timeout")
	ErrExpired             = errors.New("expired")
	ErrTooFrequently       = errors.New("operation too frequently")
	ErrOccupied            = errors.New("resource has been occupied")
	ErrDupOper             = errors.New("duplicate operations")
	ErrUnbind              = errors.New("have not bind")
	ErrFormat              = errors.New("invalid format")
	ErrNoReaction          = errors.New("target not reaction")
	ErrVersion             = errors.New("version mismatch")
	ErrOnline              = errors.New("target online")
	ErrOffline             = errors.New("target offline")
	ErrAuth                = errors.New("Auth failed")
	ErrFailed              = errors.New("operation failed")
	ErrPending             = errors.New("pending")
	ErrMaintaining         = errors.New("service is maintaining")
)
