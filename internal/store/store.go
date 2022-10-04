package store

import (
	c "backendForSharedProject/config"
	"backendForSharedProject/internal/handlers"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Database struct {
	DB *sql.DB
}

type User struct {
	ID                int
	Login             string
	EncryptedPassword string
}

func ConnectToDB(path string) *sql.DB {
	//загрузка конфига с параметрами БД
	psqlConfig, err := c.LoadPSQLConfig(path)
	if err != nil {
		log.Println("error reading psql config:", err)
		os.Exit(2)
	}

	//форматирование параметров для подключения к БД
	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		psqlConfig.Host, psqlConfig.Port, psqlConfig.User, psqlConfig.Password, psqlConfig.DBName)

	//открытие БД
	db, err := sql.Open("postgres", psqlConnectionString)
	if err != nil {
		log.Println("error opening database:", err)
		os.Exit(3)
	}

	//подключение к БД
	err = db.Ping()
	if err != nil {
		log.Println("error connecting to database:", err)
		os.Exit(4)
	}
	log.Println("DB successfully connected.")

	return db
}

func GetPassword(db *sql.DB, username string) (encryptedPassword string, err error) {
	sqlStatement := `
	SELECT "user".password
	FROM "user"
	WHERE "user".username = $1;`

	if err = db.QueryRow(sqlStatement, username).Scan(&encryptedPassword); err != nil {
		return "", err
	}
	return encryptedPassword, nil
}

func CreateNewUser(db *sql.DB, user *User) (id int, err error) {
	encrPass, err := handlers.HashPassword(user.EncryptedPassword)
	if err != nil {
		err := "Password Encryption  failed"
		fmt.Println(err)
	}

	sqlStatement := `
	INSERT INTO "user" (username, password) 
	VALUES ($1, $2)
	RETURNING user_id;`

	err = db.QueryRow(sqlStatement, user.Login, encrPass).Scan(&id)
	if err != nil {
		return
	}
	return
}
func (db *Database) AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	var a handlers.Authn
	err := handlers.DecodeJSONBody(w, r, &a)
	if err != nil {
		var mr *handlers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "cheking login: %+v\n", a.Login)

	//check if username exists and retrieve its encrypted password
	encrPass, err := GetPassword(db.DB, a.Login)
	if err == sql.ErrNoRows {
		//there's no such username, create new
		u := &User{
			Login:             a.Login,
			EncryptedPassword: a.Password,
		}
		newID, err := CreateNewUser(db.DB, u)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
			return
		}
		fmt.Fprintf(w, "Created new user, ID: %d\n", newID)
		return
	} else if err != nil {
		log.Println("username search failed")
		return
	}

	corresponds := handlers.CheckPasswordHash(encrPass, a.Password)
	if corresponds {
		fmt.Fprintf(w, "successfully logged in as %s\n", a.Login)
		return
	} else {
		gotWrong, _ := handlers.HashPassword(a.Password)
		fmt.Fprintf(w, "wrong password:\n"+
			"expecting: %s\n"+
			"got: %s\n", encrPass, gotWrong)
		return
	}

	// encrypt password

	fmt.Fprintf(w, "hashed pw: %+v\n", encrPass)
}
