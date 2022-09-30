package main

import (
	c "backendForSharedProject/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	config, err := c.LoadMainConfig("../config/")
	if err != nil {
		fmt.Println("Error loading main config: ", err)
		os.Exit(1)
	}
	psqlConfig, err := c.LoadPSQLConfig(config.Path)
	if err != nil {
		fmt.Println("Error loading PSQL config: ", err)
		os.Exit(2)
	}

	psqlConnStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		psqlConfig.Host, psqlConfig.Port, psqlConfig.User, psqlConfig.Password, psqlConfig.DBName)

	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		fmt.Println("Error in opening DB: ", err)
		os.Exit(3)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error in DB connection: ", err)
		os.Exit(4)
	}

	fmt.Println("DB successfully connected.")

	//insert example without returning id
	sqlStatement1 := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement1, 30, "ivan@petrov.ru", "Ivan", "Petrov")
	if err != nil {
		fmt.Println("Error during execution of SQL query :", err)
		os.Exit(5)
	}

	//insert example with returning of new entry's id
	sqlStatement2 := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int
	//.QueryRow возвращает ровно одно значение строки. Если в запросе нет срок - ошибку ErrNoRows.
	//если строк несколько - только первую, остальные отбрасывает
	err = db.QueryRow(sqlStatement2, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
	if err != nil {
		fmt.Println("Error during execution of SQL query :", err)
		os.Exit(6)
	}

	//update example
	sqlStatement3 := `
	UPDATE users
	SET first_name = $2, last_name = $3
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement3, 1, "NewFirst", "NewLast")
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Rows updated: ", count)

	//multiple row query example
	rows, err := db.Query("SELECT id, first_name FROM users LIMIT $1", 3)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var firstName string
		err = rows.Scan(&id, &firstName)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(id, firstName)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
