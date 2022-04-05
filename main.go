package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// Структура новости
type new struct {
	ID     string  `json:"id" db:"id"`
	Title  string  `json:"title" db:"title"`
	Text string  `json:"text" db:"text"`
	Date  string `json:"time" db:"time"`
}

var schema = `
CREATE TABLE news (
    id text,
    title text,
    text text,
	time text
)`

//Слайс новостей
var news = []new{}

func main() {
	//Коннект к БД
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 password=123321 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	//Создание БД - используется при первом запуске
	//db.MustExec(schema)

	//Чтение списка из БД
	db.Select(&news, "SELECT * FROM news ORDER BY id ASC")
	db.Close()
	
	router := gin.Default()
	
	router.GET("/news", getNews)
	router.GET("/news/id/:id", getNewByID)
	router.GET("/news/title/:title", getNewByTitle)
	router.POST("/news", postNews)

	router.Run("localhost:8080")
}

// Список всех новостей
func getNews(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, news)
}

// Добавление записи через POST
func postNews(c *gin.Context) {
	//Открытие БД
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 password=123321 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	//Новая запись
	var newNew new

	// Перевод из JSON
	if err := c.BindJSON(&newNew); err != nil {
		c.IndentedJSON(http.StatusBadRequest, newNew)
		return
	}
	// Добавляем запись в слайс
	news = append(news, newNew)
	// И в БД
	_, err = db.Exec("INSERT INTO news (id, title, text, time) VALUES ($1, $2, $3, $4)", newNew.ID, newNew.Title, newNew.Text, newNew.Date)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, newNew)
		return
	}
	db.Close()
	c.IndentedJSON(http.StatusCreated, newNew)
}

// Получение записи по ID
func getNewByID(c *gin.Context) {
	id := c.Param("id")

	// Цикл поиска по ID
	for _, a := range news {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// Получение записи по заголовку
func getNewByTitle(c *gin.Context) {
	title := c.Param("title")

	// Цикл поиска
	for _, a := range news {
		if a.Title == title {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
