package app

import (
	"github.com/mohammadshabab/bookstore_users-api/controllers/ping"
	"github.com/mohammadshabab/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("users/:user_id", users.GetUser)
	router.GET("users/search/:user_id", users.SearchUser)
	router.POST("/users", users.CreateUser)
}
