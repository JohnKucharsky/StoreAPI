package router

import (
	"github.com/JohnKucharsky/StoreAPI/handler"
	"github.com/JohnKucharsky/StoreAPI/store"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Register(r *fiber.App, db *pgxpool.Pool, redis *redis.Client) {
	us := store.NewUserStore(db, redis)
	addressStore := store.NewAddressStore(db)

	h := handler.NewHandler(
		us,
		addressStore,
	)

	v1 := r.Group("/api")

	// auth
	auth := v1.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/login", h.SignIn)
	auth.Get("/logout", h.DeserializeUser, h.LogoutUser)
	auth.Get("/refresh", h.RefreshAccessToken)
	auth.Get("/me", h.DeserializeUser, h.GetMe)
	// end auth

	//address
	address := v1.Group("/address")
	address.Post("/", h.DeserializeUser, h.CreateAddress)
	address.Get("/", h.DeserializeUser, h.GetAddresses)
	address.Get("/:id", h.DeserializeUser, h.GetOneAddress)
	address.Put("/:id", h.DeserializeUser, h.UpdateAddress)
	address.Delete("/:id", h.DeserializeUser, h.DeleteAddress)
	// end address
}
