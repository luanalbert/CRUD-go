package servidor

import (
	"crud/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type toodo struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

//CreateToodo cria toodo no banco de dados
func CreateToodo(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var toodo toodo

	if err := json.Unmarshal(reqBody, &toodo); err != nil {
		w.Write([]byte("Falha ao converter JSON em Struct"))
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	//prepare statement - evita SQL inject
	statement, err := db.Prepare("insert into toodo (title, description) values (?,?)")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
	}
	defer statement.Close()

	insert, err := statement.Exec(toodo.Title, toodo.Description)
	if err != nil {
		w.Write([]byte("erro ao executar o statement"))
		return
	}

	ID, err := insert.LastInsertId()
	if err != nil {
		w.Write([]byte("erro ao obter ID"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso ID: %d", ID)))
}

//GetAllToodo busca todos os toodos
func GetAllToodos(w http.ResponseWriter, r *http.Request) {
	//Abrindo conexão com o banco
	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	//QUERY
	lines, err := db.Query("SELECT * FROM toodo")
	if err != nil {
		w.Write([]byte("erro ao buscar os toodos"))
	}
	defer lines.Close()

	var toodos []toodo //slice de toodos

	for lines.Next() {
		//pega cada tooodo
		var toodo toodo
		//cada toodo é scaneando linha a linha do banco para traqnforma em um struct para preencher o slice de toodos
		if err := lines.Scan(&toodo.ID, &toodo.Title, &toodo.Description); err != nil {
			w.Write([]byte("Erro ao escanear os Usuarios"))
			return
		}
		toodos = append(toodos, toodo)
	}
	w.WriteHeader(http.StatusOK)
	//tranformando o slice de toodos em json
	if err := json.NewEncoder(w).Encode(toodos); err != nil {
		w.Write([]byte("Erro ao covrter os toodos para JSON"))
	}
}

//GetOneToodo busca apenas um toodo especifico
func GetOneToodo(w http.ResponseWriter, r *http.Request) {
	reqParans := mux.Vars(r)

	ID, err := strconv.ParseUint(reqParans["id"], 10, 32) // converter para uint
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}
	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}

	line, err := db.Query("SELECT * from toodo WHERE id = ?", ID)
	if err != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}

	var toodo toodo
	if line.Next() {
		if err := line.Scan(&toodo.ID, &toodo.Title, &toodo.Description); err != nil {
			w.Write([]byte("Erro ao escanear usuario"))
			return
		}
	}
	if toodo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Toodo não encontrado"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toodo); err != nil {
		w.Write([]byte("Erro ao covrter o toodo para JSON"))
	}
}

//UpdateToodo altera os dados de um toodo
func UpdateToodo(w http.ResponseWriter, r *http.Request) {
	reqParans := mux.Vars(r)

	ID, err := strconv.ParseUint(reqParans["id"], 10, 32) // converter para uint
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao ler corpo da requisição"))
		return
	}

	var toodo toodo

	if err := json.Unmarshal(reqBody, &toodo); err != nil {
		w.Write([]byte("Erro ao lconverter toodo para Struct"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update toodo set title = ?, description = ? where id = ?")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
	}
	defer statement.Close()

	if _, err := statement.Exec(toodo.Title, toodo.Description, ID); err != nil {
		w.Write([]byte("Erro ao atualizar toodo"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteToodo(w http.ResponseWriter, r *http.Request) {
	reqParans := mux.Vars(r)

	ID, err := strconv.ParseUint(reqParans["id"], 10, 32) // converter para uint
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from toodo where id = ?")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar toodo"))
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
