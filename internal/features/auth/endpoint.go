package auth

import (
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	Service interface {
		SignUp(ctx *fiber.Ctx) error
		SignIn(ctx *fiber.Ctx) error
		RefreshAccessToken(ctx *fiber.Ctx) error
		DeserializeUser(ctx *fiber.Ctx) error
		GetMe(ctx *fiber.Ctx) error
		LogoutUser(ctx *fiber.Ctx) error
	}

	service struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &service{repository: store}
}

func (h *service) SignUp(c *fiber.Ctx) error {
	var input domain.SignUpInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	err := input.HashPassword()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	res, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(res))
}

func (h *service) SignIn(c *fiber.Ctx) error {
	var input domain.SignInInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	signedRes, err := h.repository.GetOne(c.Context(), strings.ToLower(input.Email), "")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	ok, err := input.CheckPassword(signedRes.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(shared.ErrorRes("passwords don't match"))
	}

	accessToken, err := h.repository.SetAccessToken(c.Context(), signedRes.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	refreshToken, err := h.repository.SetRefreshToken(c.Context(), signedRes.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	var accessTokenMaxAgeString = os.Getenv("ACCESS_TOKEN_MAXAGE")
	var refreshTokenMaxAgeString = os.Getenv("REFRESH_TOKEN_MAXAGE")
	var accessTokenMaxAge, _ = strconv.Atoi(accessTokenMaxAgeString)
	var refreshTokenMaxAge, _ = strconv.Atoi(refreshTokenMaxAgeString)

	c.Cookie(
		&fiber.Cookie{
			Name:     "access_token",
			Value:    *accessToken,
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "refresh_token",
			Value:    *refreshToken,
			Path:     "/",
			MaxAge:   refreshTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: false,
		},
	)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"access_token": accessToken,
			"data":         signedRes,
		},
	)
}

func (h *service) RefreshAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(
			shared.ErrorRes("no refresh token in cookies"),
		)
	}

	userID, err := h.repository.GetByRefreshTokenRedis(c.Context(), refreshToken)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(shared.ErrorRes(err.Error()))
	}

	user, err := h.repository.GetOne(c.Context(), "", userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	accessToken, err := h.repository.SetAccessToken(c.Context(), user.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	var accessTokenMaxAgeString = os.Getenv("ACCESS_TOKEN_MAXAGE")
	var accessTokenMaxAge, _ = strconv.Atoi(accessTokenMaxAgeString)

	c.Cookie(
		&fiber.Cookie{
			Name:     "access_token",
			Value:    *accessToken,
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: false,
		},
	)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"access_token": accessToken,
		},
	)
}

func (h *service) DeserializeUser(c *fiber.Ctx) error {
	var accessToken string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		accessToken = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		accessToken = c.Cookies("access_token")
	}

	if accessToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(
			shared.ErrorRes("No access token"),
		)
	}

	userID, tokenUUID, err := h.repository.GetByAccessTokenRedis(c.Context(), accessToken)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(shared.ErrorRes(err.Error()))
	}

	res, err := h.repository.GetOne(c.Context(), "", userID)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(shared.ErrorRes(err.Error()))
	}

	c.Locals("user", res)
	c.Locals("access_token_uuid", tokenUUID)

	return c.Next()
}

func (h *service) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(user))
}

func (h *service) LogoutUser(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared.ErrorRes("No refresh token in the cookies"))
	}
	accessToken, ok := c.Locals("access_token_uuid").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(shared.ErrorRes("Access token is not a string"))
	}
	if accessToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(shared.ErrorRes("Access token an empty string"))
	}

	err := h.repository.DeleteTokensRedis(c.Context(), refreshToken, accessToken)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	now := time.Now()

	c.Cookie(
		&fiber.Cookie{
			Name:    "access_token",
			Value:   "",
			Expires: now,
		},
	)
	c.Cookie(
		&fiber.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: now,
		},
	)
	c.Cookie(
		&fiber.Cookie{
			Name:    "logged_in",
			Value:   "",
			Expires: now,
		},
	)

	return c.SendStatus(http.StatusOK)
}
