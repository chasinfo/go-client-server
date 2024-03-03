package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TCotacao struct {
	gorm.Model
	Id         uint64 `gorm:"primaryKey; not null"`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

func DbConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./database/dbCotacao.db"), &gorm.Config{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao conectar ao banco de dados: %v\n", err)
	}

	return db
}

func AutoMigrate() error {
	return DbConnect().AutoMigrate(&TCotacao{})
}
