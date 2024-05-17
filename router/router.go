package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saratkumar-yb/infinityapi/handlers"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/yba", handlers.InsertYbaHandler)
	router.POST("/ybdb", handlers.InsertYbdbHandler)
	router.POST("/compatibility", handlers.InsertCompatibilityHandler)
	router.POST("/compatibility_list", handlers.GetCompatibleYbdbHandler)
	return router
}
