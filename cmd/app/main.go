package main

import (
	"awesomeProject1/ent"
	"context"
	//"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jung-kurt/gofpdf" // PDF generation
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"log"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
	"net/http"
	"os"
	"sync"
	"time"
)

// Configurations
const (
	dbFilePath     = "inventory.db"
	reportEndpoint = "/generate-report"
	loginEndpoint  = "/login"
)

var secretKey = os.Getenv("SECRET_KEY")

var client *ent.Client // Global client for demonstration purposes
var ctx context.Context

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateToken Generate JWT token (for demonstration purposes)
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token expires in 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// AuthMiddleware Middleware to validate JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handle login requests
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate credentials (simple hardcoded validation for demonstration)
	if user.Username != "user" || user.Password != "****" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Fetch inventory data concurrently
func fetchInventoryData(resultChan chan map[string]interface{}, errChan chan error) {
	headers := []string{"item_id", "quantity", "price"}

	// Query all inventory items using Ent
	items, err := client.Inventory.Query().All(ctx)
	if err != nil {
		errChan <- err
		return
	}

	// Build the result map
	result := make(map[string]interface{})
	for _, item := range items {
		result[fmt.Sprintf("%d", item.ItemID)] = map[string]interface{}{
			"quantity": item.Quantity,
			"price":    item.Price,
		}
	}

	// Send headers and data to the result channel
	resultChan <- map[string]interface{}{
		"headers": headers,
		"data":    result,
	}
}

// Generate PDF report concurrently
func generatePDFReport(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan map[string]interface{})
	errChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fetchInventoryData(resultChan, errChan)
	}()

	// Wait for goroutine to finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	select {
	case result := <-resultChan:
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, "Inventory Report")
		pdf.Ln(12)

		pdf.SetFont("Arial", "", 12)
		pdf.Cell(30, 10, "Item ID")
		pdf.Cell(30, 10, "Quantity")
		pdf.Cell(30, 10, "Price")
		pdf.Ln(10)

		for itemID, data := range result["data"].(map[string]interface{}) {
			dataMap := data.(map[string]interface{})
			pdf.Cell(30, 10, itemID)
			pdf.Cell(30, 10, fmt.Sprintf("%d", dataMap["quantity"]))
			pdf.Cell(30, 10, fmt.Sprintf("%.2f", dataMap["price"]))
			pdf.Ln(10)
		}

		err := pdf.Output(w)
		if err != nil {
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			return
		}
	case err := <-errChan:
		http.Error(w, "Error fetching inventory data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Inventory represents the database table
type Inventory struct {
	ItemID   int `gorm:"primaryKey"`
	Quantity int
	Price    float64
}

func main() {
	var err error
	// Ensure the secret key is set
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable is not set")
	}

	////// Connect to SQLite database
	//db, err = sql.Open("sqlite", dbFilePath)
	//if err != nil {
	//	log.Fatalf("Failed to connect to database: %v", err)
	//}
	//defer db.Close()

	client, err := ent.Open("sqlite3", "file:test.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Insert a new inventory item
	newItem, err := client.Inventory.
		Create().
		SetItemID(1).
		SetQuantity(100).
		SetPrice(9.99).
		Save(context.Background())
	if err != nil {
		log.Fatalf("failed creating inventory item: %v", err)
	}
	log.Printf("New Inventory Item: %+v\n", newItem)

	// Initialize database schema
	//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS go run -mod=mod entgo.io/ent/cmd/ent new User (
	//	item_id TEXT PRIMARY KEY,
	//	quantity INTEGER,
	//	price REAL
	//)`)

	//if err != nil {
	//	log.Fatalf("Failed to initialize database schema: %v", err)
	//}

	// Setup HTTP server
	http.HandleFunc(loginEndpoint, loginHandler)
	http.Handle(reportEndpoint, AuthMiddleware(http.HandlerFunc(generatePDFReport)))
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
