package api

import (
	"errors"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandler struct {
	userStore db.UserStoreInterface
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type AuthParams struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

func NewAuthHandler(userStore db.UserStoreInterface) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var auth AuthParams
	if err := c.BodyParser(&auth); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), auth.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid Credentials")
		}
	}
	if !types.IsValidPassword(user.EncPassword, auth.Password) {
		return fmt.Errorf("auth failed")
	}
	fmt.Println("Auth success")
	token := CreateTokenFromUser(*user)

	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

func CreateTokenFromUser(user types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Failed to sign token with secret")
	}
	return tokenStr

}
