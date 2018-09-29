package main

func main() {
	app := App{ 
				database: NewPostGreSQLDatabase(),
				router: NewMuxRouter(),
				encryptionService: NewBCryptEncryptionService(),
			}
	app.run()
}
