package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Connect abre conexão com o banco de dados
func Connect() (*sql.DB, error) { //2 returnos
	stringConexao := "root:senha@/nomedobanco?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", stringConexao)
	if err != nil {
		return nil, err // da pra ver os 2 retornos aqui
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
