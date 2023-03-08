package routes

import (
	"basic-user-auth/database"
	"basic-user-auth/models"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Auth(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", register)
	auth.Post("/login", login)
	auth.Get("/confidential", confidential)
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	password := string(hash)
	user.Password = password
	users := database.GetDBCollection("users")
	res, err := users.InsertOne(c.Context(), user)
	if err != nil {
		return err
	}
	c.JSON(res)
	return nil
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.New("Error while parsing the body").Error(),
		})
	}
	resDb := new(models.User)
	users := database.GetDBCollection("users")
	err := users.FindOne(c.Context(), bson.M{"email": user.Email}).Decode(resDb)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.New("Unauthorized").Error(),
		})
	}
	match := bcrypt.CompareHashAndPassword([]byte(resDb.Password), []byte(user.Password))
	if match != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.New("Wrong password").Error(),
		})
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iss":   os.Getenv("JWT_ISSUER"),
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.New("Error while signing the token").Error(),
		})
	}

	c.JSON(fiber.Map{
		"token": token,
		"user":  resDb,
	})
	return nil
}

type Claim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func verifyJWT(c *fiber.Ctx) error {
	btoken := c.Request().Header.Peek("Authorization")[7:]
	if btoken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	fmt.Println(string(btoken))
	token, err := jwt.ParseWithClaims(string(btoken), &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	_, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Next()
}

func confidential(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "boom",
	})
}
