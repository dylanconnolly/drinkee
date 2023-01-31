package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	SongCount int32  `json:"song_count"`
}

var albums = []album{
	{"1", "Like a Bird", "Nelly Furtado", 9},
	{"2", "Get Back", "Ludacris", 12},
	{"3", "Hi Hi Momma", "Tupac", 11},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func CreateRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/albums", getAlbums)
	return r
}
