package main

// Create our app and inject our implementation details.
// Postgres, Mux, Bcryp, Jwt, http...
func main() {

	app := App{
		database:          NewPostGreSQLDatabase(),
		router:            NewMuxRouter(),
		encryptionService: NewBCryptEncryptionService(),
		authService:       NewJWTAuthService(),
	}
	app.run()
}
