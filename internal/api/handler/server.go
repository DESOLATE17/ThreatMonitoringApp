package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"threat-monitoring/internal/models"
)

var cards = []models.Threat{
	{
		Id:          "0",
		Name:        "Фишинг",
		Description: "Попытка обманом заставить людей поделиться конфиденциальной информацией, такой как пароли или данные кредитной карты, выдавая себя за заслуживающую доверия организацию посредством электронной почты, сообщений или веб-сайтов.",
		Image:       "/image/phishing.jpg",
		Count:       1850392,
		Price:       30000,
	},
	{
		Id:          "1",
		Name:        "Вредоносное ПО",
		Description: "Вредоносное программное обеспечение, предназначенное для нарушения работы, повреждения или получения несанкционированного доступа к компьютерным системам, включая вирусы, черви, трояны и программы-вымогатели.",
		Image:       "/image/malware.jpeg",
		Count:       505879385,
		Price:       40000,
	},
	{
		Id:          "2",
		Name:        "DDoS Атака",
		Description: "Перегрузка сети или веб-сайта потоком трафика из нескольких источников, приводящая к его зависанию или сбою.",
		Image:       "/image/ddos.jpeg",
		Count:       384800,
		Price:       70000,
	},
	{
		Id:          "3",
		Name:        "Поиск SQL инъекций",
		Description: "Использование уязвимостей в базе данных веб-сайта для внедрения вредоносных команд SQL, что потенциально позволяет злоумышленникам получить доступ к конфиденциальным данным или манипулировать ими.",
		Image:       "/image/sql_injection.jpeg",
		Price:       35000,
	},
}

func StartServer() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")

	r.GET("/home", GetCardsList)
	r.GET("/card/:id", GetCardById)

	r.Static("/image", "./resources")
	return r
}

func GetCardsList(c *gin.Context) {
	query := c.Query("query")
	resultThreats := make([]models.Threat, 0, len(cards))
	fmt.Println(query)
	if query != "" {
		for _, v := range cards {
			if strings.Contains(v.Name, query) {
				resultThreats = append(resultThreats, v)
			}
		}
	} else {
		resultThreats = cards
	}
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title": "Main website",
		"items": resultThreats,
	})
}

func GetCardById(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.HTML(http.StatusOK, "card.gohtml", gin.H{
		"title": "Main website",
		"card":  cards[id],
	})
}
