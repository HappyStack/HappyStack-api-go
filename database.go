package main

type Database interface {
	create(i item) (item, error)
	read(id int) item
	update(i item) (item, error)
	delete(id int) error
	itemsFor(userId int) []item
	passwordForUserEmail(email string) (string, error)
	close()
}