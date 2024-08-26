package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/rent", RentHandler)
	mux.HandleFunc("/dashboard", AuthMiddleware(DashboardHandler))

	fmt.Println("Server is Running on port :8081")
	http.ListenAndServe(":8081", mux)
}
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roleUser := ctx.Value("RoleUser")
	secretkey := ctx.Value("Secret")
	response := fmt.Sprintf(" Role : %d\nSecret key is %v", roleUser, secretkey)
	fmt.Fprintf(w, response)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, "RoleUser", "admin")
		ctx = context.WithValue(ctx, "Secret", "usertoken123")

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RentHandler(w http.ResponseWriter, r *http.Request) {
	//context Timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "userID", "1234")

	done := make(chan bool)
	go rentCar(ctx, done)

	select {
	case <-ctx.Done():
		http.Error(w, "Operation cancelled or timeout", http.StatusGatewayTimeout)
		return
	case <-done:
		fmt.Fprintf(w, "Car rented successfully")
	}

}

func rentCar(ctx context.Context, done chan bool) {
	select {
	case <-time.After(6 * time.Second):
		if ctx.Err() != nil {
			fmt.Println("Opearation cancelled:", ctx.Err())
			return
		}

		userID := ctx.Value("userID").(string)
		fmt.Println("Renting car for user: ", userID)
		done <- true
	case <-ctx.Done():
		fmt.Println("Opearation cancelled:", ctx.Err())
		return
	}
}
