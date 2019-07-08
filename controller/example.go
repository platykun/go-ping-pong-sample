package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go/service"
	"gopkg.in/olahol/melody.v1"
	"log"
	"strings"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ping/:id", func(c *gin.Context) {
		id := c.Param("id")
		resMessage := "Hello" + id
		c.JSON(200, gin.H{"message": resMessage})
	})
	r.POST("/terminal", func(c *gin.Context) {
		terminal := service.GenerateID()
		c.JSON(200, gin.H{"terminalId": terminal})
	})
	r.POST("/terminal/:id/data", func(c *gin.Context) {

		// request body の値をrequestBodyMapの中に格納する
		buf := make([]byte, 2048)
		n, _ := c.Request.Body.Read(buf)
		if n == 0 {
			return
		}
		var requestBodyMap map[string]interface{}
		json.Unmarshal(buf[0:n], &requestBodyMap)
		//if err1 != nil {
		//	log.Printf(err1.Error())
		//}
		//heartRate := int(requestBodyMap["heartRate"].(float64))
		//strHeartRate := strconv.Itoa(heartRate)
		//heartRateByte := make([]byte, heartRate)
		//heartRateByte, _ := requestBodyMap["heartRate"].(*[]byte)

		id := c.Param("id")
		c.JSON(200, gin.H(requestBodyMap))

		err := m.BroadcastFilter(buf, func(q *melody.Session) bool {
			url, exists := q.Get("url")
			if exists && strings.Contains(url.(string), id) {
				return true
			}
			return false
		})
		if err != nil {
			log.Printf(err.Error())
		}
		return
	})

	r.GET("/terminal/:id/data/ws", func(c *gin.Context) {
		c.Set("id", "aaa")
		value, _ := c.Get("id")
		log.Printf(value.(string))

		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	m.HandleConnect(func(s *melody.Session) {
		url := s.Request.URL
		s.Set("url", url.String())
		log.Printf("websocket connection open. [session: %#v]\n", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close. [session: %#v]\n", s)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
