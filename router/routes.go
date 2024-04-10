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
	productStore := store.NewProductStore(db)
	orderStore := store.NewOrderStore(db)
	shelfStore := store.NewShelfStore(db)
	operationsStore := store.NewOperationsStore(db)

	h := handler.NewHandler(
		us,
		addressStore,
		productStore,
		orderStore,
		shelfStore,
		operationsStore,
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

	//product
	product := v1.Group("/product")
	product.Post("/", h.DeserializeUser, h.CreateProduct)
	product.Get("/", h.DeserializeUser, h.GetProducts)
	product.Get("/:id", h.DeserializeUser, h.GetOneProduct)
	product.Put("/:id", h.DeserializeUser, h.UpdateProduct)
	product.Delete("/:id", h.DeserializeUser, h.DeleteProduct)
	//end product

	//order
	order := v1.Group("/order")
	order.Post("/", h.DeserializeUser, h.CreateOrder)
	order.Get("/", h.DeserializeUser, h.GetOrders)
	order.Get("/:id", h.DeserializeUser, h.GetOneOrder)
	order.Put("/:id", h.DeserializeUser, h.UpdateOrder)
	order.Delete("/:id", h.DeserializeUser, h.DeleteOrder)
	//end order

	//shelf
	shelf := v1.Group("/shelf")
	shelf.Post("/", h.DeserializeUser, h.CreateShelf)
	shelf.Get("/", h.DeserializeUser, h.GetShelves)
	shelf.Get("/:id", h.DeserializeUser, h.GetOneShelf)
	shelf.Put("/:id", h.DeserializeUser, h.UpdateShelf)
	shelf.Delete("/:id", h.DeserializeUser, h.DeleteShelf)
	//end shelf

	//operations
	operations := v1.Group("/operations")
	operations.Get("/assembly_info", h.DeserializeUser, h.GetAssemblyInfo)
	//end shelf
}
