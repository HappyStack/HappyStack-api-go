package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sacha"
	password = ""
	dbname   = "sacha"
)

var currentId int

var items Items

type HappyStackDatabase struct {
	sqlDB *sql.DB
}

func NewHappyStackDatabase() *HappyStackDatabase {

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
	return &HappyStackDatabase{sqlDB: adb}
}

func (hsdb *HappyStackDatabase) closeDatabase() {
	hsdb.sqlDB.Close()
}

func (hsdb *HappyStackDatabase) allItems() []item {
	query := `SELECT item_id, name, dosage, taken_today, serving_size, serving_type FROM item`
	rows, err := hsdb.sqlDB.Query(query)
	if err != nil {
		log.Fatal(err)
		return []item{}
	}
	defer rows.Close()

	var dbItems []item
	for rows.Next() {
		var dbItem item
		err = rows.Scan(&dbItem.Id, &dbItem.Name, &dbItem.Dosage, &dbItem.TakenToday, &dbItem.ServingSize, &dbItem.ServingType)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(dbItem)
		dbItems = append(dbItems, dbItem)
	}
	return dbItems
}

func (hsdb *HappyStackDatabase) findItem(id int) item {
	for _, i := range items {
		if i.Id == id {
			return i
		}
	}
	// return empty item if not found
	return item{}
}

func (hsdb *HappyStackDatabase) createItem(i item) (item, error) {

	// Enforce default
	if i.ServingSize == 0 {
		i.ServingSize = 1
	}
	if i.ServingType == "" {
		i.ServingType = pill
	}

	query := `INSERT INTO item (name, dosage, taken_today, serving_size, serving_type) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;`
	var createdItemId int
	err := hsdb.sqlDB.QueryRow(query, i.Name, i.Dosage, i.TakenToday, i.ServingSize, i.ServingType).Scan(&createdItemId)
	if err != nil {
		return item{}, err
	}
	i.Id = createdItemId
	return i, nil
}

func (hsdb *HappyStackDatabase) destroyItem(id int) error {
	for i, item := range items {
		if item.Id == id {
			items = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
