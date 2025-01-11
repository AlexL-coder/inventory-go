package handlers

import (
	"awesomeProject1/ent"
	"awesomeProject1/internal/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

func createUser(c *gin.Context) {
	var req CreateUserRequest
	// Bind the incoming JSON payload to the request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	resultChan := make(chan *ent.User)
	errorChan := make(chan error)

	go func() {
		// Create the user in the database
		user, err := db.Client.User.
			Create().
			SetName(req.Name).
			SetEmail(req.Email).
			SetAge(req.Age).
			SetPwd(string(hashedPwd)).
			Save(c.Request.Context())
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- user
	}()

	select {
	case user := <-resultChan:
		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	case err := <-errorChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
	}
}

func listUsers(c *gin.Context) {

	resultChan := make(chan []*ent.User)
	errorChan := make(chan error)

	go func() {
		users, err := db.Client.User.Query().All(c.Request.Context())
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- users

	}()

	select {
	case users := <-resultChan:
		c.JSON(http.StatusOK, users)
	case err := <-errorChan:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
}

//// Fetch inventory data concurrently
//func fetchInventoryData(resultChan chan map[string]interface{}, errChan chan error) {
//	headers := []string{"item_id", "quantity", "price"}
//
//	// Query all inventory items using Ent
//	items, err := client.Inventory.Query().All(ctx)
//	if err != nil {
//		errChan <- err
//		return
//	}
//
//	// Build the result map
//	result := make(map[string]interface{})
//	for _, item := range items {
//		result[fmt.Sprintf("%d", item.ItemID)] = map[string]interface{}{
//			"quantity": item.Quantity,
//			"price":    item.Price,
//		}
//	}
//
//	// Send headers and data to the result channel
//	resultChan <- map[string]interface{}{
//		"headers": headers,
//		"data":    result,
//	}
//}

//// Generate PDF report concurrently
//func generatePDFReport(w http.ResponseWriter, r *http.Request) {
//	resultChan := make(chan map[string]interface{})
//	errChan := make(chan error)
//	var wg sync.WaitGroup
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		fetchInventoryData(resultChan, errChan)
//	}()
//
//	// Wait for goroutine to finish
//	go func() {
//		wg.Wait()
//		close(resultChan)
//		close(errChan)
//	}()
//
//	select {
//	case result := <-resultChan:
//		pdf := gofpdf.New("P", "mm", "A4", "")
//		pdf.AddPage()
//		pdf.SetFont("Arial", "B", 16)
//		pdf.Cell(40, 10, "Inventory Report")
//		pdf.Ln(12)
//
//		pdf.SetFont("Arial", "", 12)
//		pdf.Cell(30, 10, "Item ID")
//		pdf.Cell(30, 10, "Quantity")
//		pdf.Cell(30, 10, "Price")
//		pdf.Ln(10)
//
//		for itemID, data := range result["data"].(map[string]interface{}) {
//			dataMap := data.(map[string]interface{})
//			pdf.Cell(30, 10, itemID)
//			pdf.Cell(30, 10, fmt.Sprintf("%d", dataMap["quantity"]))
//			pdf.Cell(30, 10, fmt.Sprintf("%.2f", dataMap["price"]))
//			pdf.Ln(10)
//		}
//
//		err := pdf.Output(w)
//		if err != nil {
//			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
//			return
//		}
//	case err := <-errChan:
//		http.Error(w, "Error fetching inventory data: "+err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
