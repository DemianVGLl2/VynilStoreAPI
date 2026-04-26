package handlers

import (
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Time Out", Artist: "Dave Brubeck", Price: 37.99},
	{ID: "3", Title: "Flying Beagle", Artist: "Himiko Kikuchi", Price: 69.99},
}

var users = map[string]string{
	"admin": "admin123",
	"user1": "pass1",
}

// token → username - one entry per active session
var activeTokens = map[string]string{}

// Mutexes keep concurrent users from corrupting shared data
var albumsMu sync.RWMutex
var tokensMu sync.RWMutex

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateToken() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = tokenChars[rand.Intn(len(tokenChars))]
	}
	return string(b)
}

// getUserFromToken reads the token from the Authorization header and returns the username
func getUserFromToken(c *gin.Context) (string, string, bool) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return "", "", false
	}
	token := strings.TrimPrefix(header, "Bearer ")

	tokensMu.RLock()
	username, ok := activeTokens[token]
	tokensMu.RUnlock()

	return username, token, ok
}

func RequireAuth(c *gin.Context) {
	username, token, ok := getUserFromToken(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}
	// Pass values to the handler
	c.Set("username", username)
	c.Set("token", token)
	c.Next()
}

func Login(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "provide credentials via Basic Auth"})
		return
	}

	storedPassword, exists := users[username]
	if !exists || storedPassword != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := generateToken()

	tokensMu.Lock()
	activeTokens[token] = username
	tokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username + ", welcome to the Store System",
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	username := c.GetString("username")
	token := c.GetString("token")

	tokensMu.Lock()
	delete(activeTokens, token)
	tokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Bye " + username + ", your token has been revoked",
	})
}

func GetAlbums(c *gin.Context) {
	albumsMu.RLock()
	defer albumsMu.RUnlock()

	c.JSON(http.StatusOK, albums)
}

func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	albumsMu.RLock()
	defer albumsMu.RUnlock()

	for _, album := range albums {
		if album.ID == id {
			c.JSON(http.StatusOK, album)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "album '" + id + "' not found"})
}

func CreateAlbum(c *gin.Context) {
	var newAlbum Album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: " + err.Error()})
		return
	}

	if newAlbum.ID == "" || newAlbum.Title == "" || newAlbum.Artist == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, title and artist are required"})
		return
	}

	albumsMu.RLock()
	defer albumsMu.RUnlock()

	for _, album := range albums {
		if album.ID == newAlbum.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "album id '" + newAlbum.ID + "' already exists"})
			return
		}
	}

	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

func Status(c *gin.Context) {
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username + ", the DPIP System is Up and Running",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}
