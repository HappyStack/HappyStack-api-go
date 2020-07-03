package main

type Database interface {
	create(i item) (item, error)
	read(id int) (item, error)
	update(i item) (item, error)
	delete(id int) error
	itemsFor(userId int) []item
	userMatchingEmail(email string) (User, error)
	close()
}
