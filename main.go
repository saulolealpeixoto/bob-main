package main

import (
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Devatoria/go-graylog"
)

// Book struct
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func init() {
	books = append(books, Book{ID: "1", Title: "Book 1", Author: "Author 1"})
	books = append(books, Book{ID: "2", Title: "Book 2", Author: "Author 2"})
}

// LogEntry struct for logging
type LogEntry struct {
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	Method    string    `json:"method"`
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

// Get IP address from request
func getIPAddress(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "localhost"
	}
	return ip
}

// String converte a estrutura LogEntry em uma string formatada
func (l LogEntry) String() string {
	return fmt.Sprintf("ID: %s | Endpoint: %s | Method: %s | IP: %s | Timestamp: %s\n", l.ID, l.Endpoint, l.Method, l.IP, l.Timestamp.Format(time.RFC3339))
}

func main() {
	// Inicialize um novo cliente graylog com TCP
	g, err := graylog.NewGraylog(graylog.Endpoint{
		Transport: graylog.TCP,
		Address:   "172.30.0.1",
		Port:      12201,
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// Middleware para registrar os logs no Graylog
	r.Use(func(c *gin.Context) {
		logEntry := LogEntry{
			ID:        uuid.New().String(),
			Endpoint:  c.FullPath(),
			Method:    c.Request.Method,
			IP:        getIPAddress(c),
			Timestamp: time.Now(),
		}
		err := g.Send(graylog.Message{
			Version:      "1.1",
			Host:         "localhost",
			ShortMessage: "Endpoint accessed",
			FullMessage:  logEntry.String(),
			Timestamp:    time.Now().Unix(),
			Level:        1,
			Extra: map[string]string{
				"Endpoint": logEntry.Endpoint,
				"Method":   logEntry.Method,
				"IP":       logEntry.IP,
			},
		})
		if err != nil {
			log.Printf("Erro ao enviar log para o Graylog: %s\n", err.Error())
		}
		c.Next()
	})

	r.GET("/books", getBooks)
	r.GET("/book/:id", getBook)
	r.POST("/books", createBook)
	r.PUT("/book/:id", updateBook)
	r.DELETE("/book/:id", deleteBook)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(r.Run(":8080"))
}

// Get all books
func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

// Get single book
func getBook(c *gin.Context) {
	id := c.Param("id")
	for _, item := range books {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
}

// Create a new book
func createBook(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, book)
	c.JSON(http.StatusCreated, book)
}

// Update a book
func updateBook(c *gin.Context) {
	id := c.Param("id")
	for index, item := range books {
		if item.ID == id {
			if err := c.BindJSON(&books[index]); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			books[index].ID = id
			c.JSON(http.StatusOK, books[index])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
}

// Delete a book
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	for index, item := range books {
		if item.ID == id {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	c.JSON(http.StatusOK, books)
}

// package main

// import (
// 	"log"
// 	"net/http"
// 	"time"
// 	"fmt"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/Devatoria/go-graylog"
// )

// // Book struct
// type Book struct {
// 	ID     string `json:"id"`
// 	Title  string `json:"title"`
// 	Author string `json:"author"`
// }

// var books []Book

// func init() {
// 	books = append(books, Book{ID: "1", Title: "Book 1", Author: "Author 1"})
// 	books = append(books, Book{ID: "2", Title: "Book 2", Author: "Author 2"})
// }

// // LogEntry struct for logging
// // type LogEntry struct {
// // 	ID        string    `json:"id"`
// // 	Endpoint  string    `json:"endpoint"`
// // 	Method    string    `json:"method"`
// // 	IP        string    `json:"ip"`
// // 	Timestamp time.Time `json:"timestamp"`
// // }

// // Get IP address from request
// func getIPAddress(c *gin.Context) string {
// 	ip := c.ClientIP()
// 	if ip == "::1" {
// 		ip = "localhost"
// 	}
// 	return ip
// }

// // LogMiddleware middleware to log requests
// func LogMiddleware(c *gin.Context) {
// 	logEntry := LogEntry{
// 		ID:        uuid.New().String(),
// 		Endpoint:  c.FullPath(),
// 		Method:    c.Request.Method,
// 		IP:        getIPAddress(c),
// 		Timestamp: time.Now(),
// 	}
// 	log.Printf("ID: %s | Endpoint: %s | Method: %s | IP: %s | Timestamp: %s\n", logEntry.ID, logEntry.Endpoint, logEntry.Method, logEntry.IP, logEntry.Timestamp.Format(time.RFC3339))
// 	c.Next()
// }

// // Get all books
// func getBooks(c *gin.Context) {
// 	c.JSON(http.StatusOK, books)
// }

// // Get single book
// func getBook(c *gin.Context) {
// 	id := c.Param("id")
// 	for _, item := range books {
// 		if item.ID == id {
// 			c.JSON(http.StatusOK, item)
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// }

// // Create a new book
// func createBook(c *gin.Context) {
// 	var book Book
// 	if err := c.BindJSON(&book); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	books = append(books, book)
// 	c.JSON(http.StatusCreated, book)
// }

// // Update a book
// func updateBook(c *gin.Context) {
// 	id := c.Param("id")
// 	for index, item := range books {
// 		if item.ID == id {
// 			if err := c.BindJSON(&books[index]); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 				return
// 			}
// 			books[index].ID = id
// 			c.JSON(http.StatusOK, books[index])
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// }

// // Delete a book
// func deleteBook(c *gin.Context) {
// 	id := c.Param("id")
// 	for index, item := range books {
// 		if item.ID == id {
// 			books = append(books[:index], books[index+1:]...)
// 			break
// 		}
// 	}
// 	c.JSON(http.StatusOK, books)
// }

// func main() {
// 	// Initialize a new graylog client with TCP
// 	g, err := graylog.NewGraylog(graylog.Endpoint{
// 		Transport: graylog.TCP,
// 		Address:   "172.30.0.1",
// 		Port:      12201,
// 	})
	
// 	if err != nil {
// 		panic(err)
// 	}

// 	r := gin.Default()

// 	// Middleware para registrar os logs no Graylog
// 	r.Use(func(c *gin.Context) {
// 		logEntry := LogEntry{
// 			ID:        uuid.New().String(),
// 			Endpoint:  c.FullPath(),
// 			Method:    c.Request.Method,
// 			IP:        getIPAddress(c),
// 			Timestamp: time.Now(),
// 		}
// 		err := g.Send(graylog.Message{
// 			Version:      "1.1",
// 			Host:         "localhost",
// 			ShortMessage: "Endpoint accessed",
// 			FullMessage:  logEntry.String(), // Usando o logEntry como mensagem completa
// 			Timestamp:    time.Now().Unix(),
// 			Level:        1,
// 			Extra: map[string]string{
// 				"Endpoint": logEntry.Endpoint,
// 				"Method":   logEntry.Method,
// 				"IP":       logEntry.IP,
// 			},
// 		})
// 		if err != nil {
// 			log.Printf("Error sending log to Graylog: %s\n", err.Error())
// 		}
// 		c.Next()
// 	})

// 	r.GET("/books", getBooks)
// 	r.GET("/book/:id", getBook)
// 	r.POST("/books", createBook)
// 	r.PUT("/book/:id", updateBook)
// 	r.DELETE("/book/:id", deleteBook)

// 	log.Println("Server started on port 8080")
// 	log.Fatal(r.Run(":8080"))
// }

// // Estrutura LogEntry para registrar informações de log
// type LogEntry struct {
// 	ID        string    `json:"id"`
// 	Endpoint  string    `json:"endpoint"`
// 	Method    string    `json:"method"`
// 	IP        string    `json:"ip"`
// 	Timestamp time.Time `json:"timestamp"`
// }

// // String converte a estrutura LogEntry em uma string formatada
// func (l LogEntry) String() string {
// 	return fmt.Sprintf("ID: %s | Endpoint: %s | Method: %s | IP: %s | Timestamp: %s\n", l.ID, l.Endpoint, l.Method, l.IP, l.Timestamp.Format(time.RFC3339))
// }

