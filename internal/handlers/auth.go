package handlers

import (
	"awesomeProject1/ent/user"
	"awesomeProject1/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var secretKey = os.Getenv("SECRET_KEY")

// generateToken generates a JWT token for the given username.
//
// @Summary Generate a JWT token
// @Description Generates a JSON Web Token (JWT) with a 1-hour expiration for the given username.
// @Tags Authentication
// @Param username path string true "Username"
// @Success 200 {string} string "JWT token"
// @Failure 500 {object} map[string]string "Error generating token"
// @Router /generate-token/{username} [get]
func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token expires in 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"securepassword"`
}

// loginHandler handles user login and generates a JWT token.
//
// @Summary Login user and generate JWT token
// @Description Authenticates the user with their credentials and returns a JWT token if successful.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body handlers.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Error generating token"
// @Router /login [post]
func loginHandler(c *gin.Context) {
	var req LoginRequest
	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Query the user from the database
	user, err := db.Client.User.
		Query().
		Where(user.Email(req.Username)). // Use the `email` field to query the user
		Only(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"auth-error": "Invalid username or password"})
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"auth-error": "Invalid username or password"})
		return
	}
	tokenChan := make(chan string)
	erroChan := make(chan error)

	// Generate a JWT token
	go func() {
		token, err := generateToken(req.Username)
		if err != nil {
			erroChan <- err
		}
		tokenChan <- token
	}()

	select {
	case token := <-tokenChan:
		c.JSON(http.StatusOK, gin.H{"token": token})

	case err := <-erroChan:
		c.JSON(http.StatusInternalServerError, gin.H{"auth-error": err.Error()})
		return
	}
}
