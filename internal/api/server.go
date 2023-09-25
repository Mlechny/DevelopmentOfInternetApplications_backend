package api

import (
	"github.com/gin-gonic/gin"
	"lab1/internal/api/ds"
	"log"
	"net/http"
	"strconv"
)

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	pipe := ds.GetPipeline()

	r.LoadHTMLGlob("templates/html/*")

	r.GET("/", func(c *gin.Context) {
		filter := c.Query("filter")
		field := c.Query("field")
		filteredCodes := filterCodes(pipe, filter, field)

		c.HTML(http.StatusOK, "all_codes.tmpl", gin.H{
			"title":  "Список услуг",
			"codes":  filteredCodes,
			"filter": filter,
			"field":  field,
		})
	})

	r.GET("/code/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 || id >= len(pipe) {
			c.String(http.StatusNotFound, "Страница не найдена")
			return
		}

		code := pipe[id]

		c.HTML(http.StatusOK, "code.tmpl", gin.H{
			"title": "Подробная информация",
			"code":  code,
		})
	})

	r.Static("/images", "./resources")
	r.Static("/styles", "./templates/css")

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")

}

func filterCodes(codes []ds.Code, filter string, field string) []ds.Code {
	if filter == "" {
		return codes
	}

	var filtered []ds.Code
	for _, code := range codes {
		if field == "Name" && contains(code.Name, filter) {
			filtered = append(filtered, code)
		} else if field == "Language" && contains(code.Subject, filter) {
			filtered = append(filtered, code)
		}
	}
	return filtered
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
