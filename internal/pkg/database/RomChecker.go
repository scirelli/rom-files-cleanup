package database

import (
	"github.com/scirelli/rom-files-cleanup/internal/pkg/model"
)

type RomChecker interface {
	IsValid(hash string) bool
	LookUp(hash string) model.Game
}
