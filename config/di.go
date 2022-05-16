package config

import (
	"log"
)

type Application struct {
	ErrLog 	*log.Logger
	InfoLog *log.Logger
}
