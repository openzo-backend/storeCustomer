package middlewares

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/pb"
)

type User struct {
	ID         string
	Phone      string
	IsVerified bool
}

type Middleware struct {
	UserServiceClient pb.UserServiceClient
}

type MiddlewareInterface interface {
	JwtMiddleware(c *gin.Context)
}

func NewMiddleware(userServiceClient pb.UserServiceClient) MiddlewareInterface {
	return &Middleware{UserServiceClient: userServiceClient}
}

func VerifyTokenAndGetUser(c pb.UserServiceClient, ctx context.Context, token string) (User, error) {
	// Verify JWT token
	res, err := c.GetUserWithJWT(ctx, &pb.Token{Token: token})

	if err != nil {
		return User{}, err
	}

	var user User
	user.ID = res.GetId()
	user.Phone = res.GetPhone()
	user.IsVerified = res.GetIsVerified()
	fmt.Println(user)
	return user, nil

}

func (m *Middleware) JwtMiddleware(c *gin.Context) {
	//get the token from the header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return

	}
	//validate the token
	user, err := VerifyTokenAndGetUser(m.UserServiceClient, c, token)
	if err != nil {

		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	//set the user in the context
	c.Set("user", user)
	c.Next()
}
