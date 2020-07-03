package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "sacha"
	dbname = "sacha"
)

var currentId int

type PostGreSQLDatabase struct {
	sqlDB *sql.DB
}

func NewPostGreSQLDatabase() *PostGreSQLDatabase {

	// TODO add DB password
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	adb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = adb.Ping()
	if err != nil {
		panic(err)
	}
	return &PostGreSQLDatabase{sqlDB: adb}
}

func (hsdb *PostGreSQLDatabase) close() {
	hsdb.sqlDB.Close()
}

func (hsdb *PostGreSQLDatabase) itemsFor(userId int) []item {

	query := `SELECT item_id, user_id, name, dosage, taken_today, serving_size, serving_type FROM items WHERE user_id = $1;`
	rows, err := hsdb.sqlDB.Query(query, userId)
	if err != nil {
		log.Fatal(err)
		return []item{}
	}
	defer rows.Close()

	var dbItems []item
	for rows.Next() {
		var dbItem item
		err = rows.Scan(&dbItem.Id, &dbItem.userId, &dbItem.Name, &dbItem.Dosage, &dbItem.TakenToday, &dbItem.ServingSize, &dbItem.ServingType)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(dbItem)
		dbItems = append(dbItems, dbItem)
	}
	return dbItems
}

func (hsdb *PostGreSQLDatabase) read(id int) (item, error) {

	// Look for item in the database.
	query := `SELECT * FROM items WHERE "item_id"=$1;`

	var itemID int
	var itemUserID int
	var itemName string
	var itemDosage string
	var itemTakenToday bool
	var itemServingSize int
	var itemServingType servingType
	var itemTiming time.Time

	row := hsdb.sqlDB.QueryRow(query, id)
	var item item
	switch err := row.Scan(&itemID, &itemUserID, &itemName, &itemDosage, &itemTakenToday, &itemServingSize, &itemServingType, &itemTiming); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return item, &myError{message: "No rows were returned!"}
	case nil:
		item.Id = itemID
		item.userId = itemUserID
		item.Name = itemName
		item.Dosage = itemDosage
		item.TakenToday = itemTakenToday
		item.ServingSize = itemServingSize
		item.ServingType = itemServingType
		// item.Timing = itemTiming
		return item, nil
	default:
		panic(err)
		return item, err
	}
}

func (hsdb *PostGreSQLDatabase) create(i item) (item, error) {

	// Enforce default
	if i.ServingSize == 0 {
		i.ServingSize = 1
	}
	if i.ServingType == "" {
		i.ServingType = pill
	}

	query := `INSERT INTO items (user_id, name, dosage, taken_today, serving_size, serving_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING item_id;`
	var createdItemId int
	err := hsdb.sqlDB.QueryRow(query, i.userId, i.Name, i.Dosage, i.TakenToday, i.ServingSize, i.ServingType).Scan(&createdItemId)
	if err != nil {
		return item{}, err
	}
	i.Id = createdItemId
	return i, nil
}

func (hsdb *PostGreSQLDatabase) update(i item) (item, error) {

	args := []interface{}{i.Id}
	argIndex := 2

	query := `
	UPDATE items
	SET `
	if i.Name != "" {
		log.Printf("i.name is: %v\n", i.Name)
		query += `name = $` + strconv.Itoa(argIndex)
		args = append(args, i.Name)
		argIndex += 1
	}
	if i.Dosage != "" {
		if argIndex > 2 {
			query += `, `
		}
		query += `dosage = $` + strconv.Itoa(argIndex)
		args = append(args, i.Dosage)
		argIndex += 1
	}
	query += ` WHERE item_id = $1;`
	//Timing      time.Time   `json:"timing"`

	res, err := hsdb.sqlDB.Exec(query, args...)
	// res, err := hsdb.sqlDB.Exec(query, i.Id, i.Name)

	// query := `
	// UPDATE items
	// SET name = $2,
	// dosage = $3,
	// taken_today = $4,
	// serving_size = $5,
	// serving_type = $6
	// WHERE item_id = $1;`
	// //Timing      time.Time   `json:"timing"`

	// res, err := hsdb.sqlDB.Exec(query, i.Id, i.Name, i.Dosage, i.TakenToday, i.ServingSize, i.ServingType)

	if err != nil {
		return item{Name: "error"}, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		fmt.Println(err)
	}

	return i, nil
}

func (hsdb *PostGreSQLDatabase) delete(id int) error {
	query := `DELETE FROM items WHERE "item_id"=$1;`
	_, err := hsdb.sqlDB.Exec(query, id)
	return err
}

func (hsdb *PostGreSQLDatabase) userMatchingEmail(email string) (User, error) {
	// Look for username in the database.
	query := `SELECT * FROM users WHERE "email"=$1;`
	var userId int
	var userEmail string
	var userPassword string
	row := hsdb.sqlDB.QueryRow(query, email)
	var user User
	switch err := row.Scan(&userId, &userEmail, &userPassword); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, &myError{message: "No rows were returned!"}
	case nil:
		user.Id = userId
		user.Username = userEmail
		user.Password = userPassword
		fmt.Println(userId, userEmail, userPassword)
		return user, nil
	default:
		panic(err)
		return user, err
	}
}

type myError struct {
	message string
}

func (e *myError) Error() string {
	return fmt.Sprintf("%v", e.message)
}
