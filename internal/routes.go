package internal

import (
	"github.com/JohnKucharsky/StoreAPI/internal/features/address"
	"github.com/JohnKucharsky/StoreAPI/internal/features/auth"
	"github.com/JohnKucharsky/StoreAPI/internal/features/operations"
	"github.com/JohnKucharsky/StoreAPI/internal/features/order"
	"github.com/JohnKucharsky/StoreAPI/internal/features/product"
	"github.com/JohnKucharsky/StoreAPI/internal/features/shelf"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Register(r *fiber.App, db *pgxpool.Pool, redis *redis.Client) {
	v1 := r.Group("/api")

	// auth
	authStore := auth.NewAuthStore(db, redis)
	authHandler := auth.New(authStore)

	authR := v1.Group("/auth")
	authR.Post("/sign-up", authHandler.SignUp)
	authR.Post("/login", authHandler.SignIn)
	authR.Get("/logout", authHandler.DeserializeUser, authHandler.LogoutUser)
	authR.Get("/refresh", authHandler.RefreshAccessToken)
	authR.Get("/me", authHandler.DeserializeUser, authHandler.GetMe)
	// end auth

	//address
	addressStore := address.NewAddressStore(db)
	addressHandler := address.New(addressStore)

	addrR := v1.Group("/address")
	addrR.Post("/", authHandler.DeserializeUser, addressHandler.Create)
	addrR.Get("/", authHandler.DeserializeUser, addressHandler.Get)
	addrR.Get("/:id", authHandler.DeserializeUser, addressHandler.GetOne)
	addrR.Put("/:id", authHandler.DeserializeUser, addressHandler.Update)
	addrR.Delete("/:id", authHandler.DeserializeUser, addressHandler.Delete)
	// end address

	//product
	productStore := product.NewProductStore(db)
	productHandler := product.New(productStore)

	productR := v1.Group("/product")
	productR.Post("/", authHandler.DeserializeUser, productHandler.Create)
	productR.Get("/", authHandler.DeserializeUser, productHandler.GetMany)
	productR.Get("/:id", authHandler.DeserializeUser, productHandler.GetOne)
	productR.Put("/:id", authHandler.DeserializeUser, productHandler.Update)
	productR.Delete("/:id", authHandler.DeserializeUser, productHandler.Delete)
	//end product

	//order
	orderStore := order.NewOrderStore(db)
	orderHandler := order.New(orderStore)

	orderR := v1.Group("/order")
	orderR.Post("/", authHandler.DeserializeUser, orderHandler.Create)
	orderR.Get("/", authHandler.DeserializeUser, orderHandler.GetMany)
	orderR.Get("/:id", authHandler.DeserializeUser, orderHandler.GetOne)
	orderR.Put("/:id", authHandler.DeserializeUser, orderHandler.Update)
	orderR.Delete("/:id", authHandler.DeserializeUser, orderHandler.Delete)
	//end order

	//shelf
	shelfStore := shelf.NewShelfStore(db)
	shelfHandler := shelf.New(shelfStore)

	shelfR := v1.Group("/shelf")
	shelfR.Post("/", authHandler.DeserializeUser, shelfHandler.Create)
	shelfR.Get("/", authHandler.DeserializeUser, shelfHandler.GetMany)
	shelfR.Get("/:id", authHandler.DeserializeUser, shelfHandler.GetOne)
	shelfR.Put("/:id", authHandler.DeserializeUser, shelfHandler.Update)
	shelfR.Delete("/:id", authHandler.DeserializeUser, shelfHandler.Delete)
	//end shelf

	//operations
	operationsStore := operations.NewOperationsStore(db)
	operationsHandler := operations.New(operationsStore)

	operationsR := v1.Group("/operations")
	operationsR.Get("/assembly_info", authHandler.DeserializeUser, operationsHandler.GetAssemblyInfo)
	//end shelf
}
