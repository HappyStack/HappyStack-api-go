package main

func main() {
	app := App{ 
				database: NewPostGreSQLDatabase(),
				router: NewMuxRouter(),
				encryptionService: NewBCryptEncryptionService(),
				authService: NewJWTAuthService(),
			}
	app.run()
}
