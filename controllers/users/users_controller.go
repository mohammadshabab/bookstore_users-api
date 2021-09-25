package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohammadshabab/bookstore_users-api/domain/users"
	"github.com/mohammadshabab/bookstore_users-api/services"
	"github.com/mohammadshabab/bookstore_users-api/utils/errors"
)

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("User id should be a number")
		c.JSON(err.Status, err)
		return
	}
	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}
func CreateUser(c *gin.Context) {
	var user users.User
	// //	doing same work in 1 line below
	// fmt.Println(user)
	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// 	fmt.Println(err.Error()) // Handle Json Error
	// 	return
	// }

	if err := c.ShouldBindJSON(&user); err != nil {
		// restError := errors.RestErr{
		// 	Message: "invalid json body",
		// 	Status:  http.StatusBadRequest,
		// 	Error:   "bad_request",
		// }
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr) // here we are returning error as it is because we are not sure what service is returning err about
		return
	}
	//fmt.Println(string(bytes))
	c.JSON(http.StatusCreated, result)
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me first!")
}
