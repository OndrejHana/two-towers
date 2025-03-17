package server

import "errors"

func Run() error {
	// err := godotenv.Load()
	// if err != nil {
	// 	return fmt.Errorf("Error loading .env file")
	// }
	// key := os.Getenv("CLERK_SECRET_KEY")
	// if key == "" {
	// 	return fmt.Errorf("CLERK_SECRET_KEY not set")
	// }
	// clerk.SetKey(key)
	// mux := http.NewServeMux()
	// ls := lobbystore.NewLobbyStore()
	//
	// RegisterRoutes(mux, &ls)
	//
	// return http.ListenAndServe("localhost:8080", mux)

	return errors.New("Not implemented")
}
