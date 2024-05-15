package auth

import (
	"errors"
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

type (
	AuthStoreI interface {
		Create(ctx *fasthttp.RequestCtx, u domain.SignUpInput) (*domain.User, error)
		GetOne(ctx *fasthttp.RequestCtx, email string, id string) (*domain.User, error)
		SetAccessToken(ctx *fasthttp.RequestCtx, uuid uuid.UUID) (*string, error)
		SetRefreshToken(ctx *fasthttp.RequestCtx, uuid uuid.UUID) (*string, error)
		GetByRefreshTokenRedis(ctx *fasthttp.RequestCtx, token string) (string, error)
		GetByAccessTokenRedis(ctx *fasthttp.RequestCtx, token string) (string, string, error)
		DeleteTokensRedis(ctx *fasthttp.RequestCtx, refreshToken, accessToken string) error
	}

	AuthStore struct {
		db    *pgxpool.Pool
		redis *redis.Client
	}
)

func NewAuthStore(db *pgxpool.Pool, redis *redis.Client) *AuthStore {
	return &AuthStore{
		db:    db,
		redis: redis,
	}
}

func (us *AuthStore) Create(ctx *fasthttp.RequestCtx, u domain.SignUpInput) (*domain.User, error) {
	rows, err := us.db.Query(
		ctx, `
        INSERT INTO users (id, name, last_name, middle_name, email, password)
        VALUES (@id, @name, @last_name, @middle_name, @email, @password)
        RETURNING id, name, last_name, middle_name, email, password, created_at, updated_at`,
		pgx.NamedArgs{
			"id":          uuid.Must(uuid.NewV7()),
			"name":        u.Name,
			"last_name":   u.LastName,
			"middle_name": u.MiddleName,
			"email":       u.Email,
			"password":    u.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.User],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us *AuthStore) SetAccessToken(ctx *fasthttp.RequestCtx, id uuid.UUID) (*string, error) {
	now := time.Now()

	var accessTokenExpiresInString = os.Getenv("ACCESS_TOKEN_EXPIRED_IN")
	var accessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenExpiresIn, _ := time.ParseDuration(accessTokenExpiresInString)

	accessTokenDetails, err := shared.CreateToken(
		id.String(),
		accessTokenExpiresIn,
		accessTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	if err = us.redis.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		id.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now),
	).Err(); err != nil {
		return nil, err
	}

	return accessTokenDetails.Token, nil
}

func (us *AuthStore) SetRefreshToken(ctx *fasthttp.RequestCtx, id uuid.UUID) (*string, error) {
	now := time.Now()

	var refreshTokenExpiresInString = os.Getenv("REFRESH_TOKEN_EXPIRED_IN")
	var refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")

	refreshTokenExpiresIn, _ := time.ParseDuration(refreshTokenExpiresInString)

	refreshTokenDetails, err := shared.CreateToken(
		id.String(),
		refreshTokenExpiresIn,
		refreshTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	if err = us.redis.Set(
		ctx,
		refreshTokenDetails.TokenUUID,
		id.String(),
		time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(now),
	).Err(); err != nil {
		return nil, err
	}

	return refreshTokenDetails.Token, nil
}

func (us *AuthStore) GetOne(ctx *fasthttp.RequestCtx, email, id string) (*domain.User, error) {
	var query, param string
	if email == "" {
		query = `select id,name,last_name,middle_name,email,
       password,created_at,updated_at from users where id::text = @param`
		param = id
	} else if id == "" {
		query = `select id,name,last_name,middle_name,email,
       password,created_at,updated_at from users where email = @param`
		param = email
	} else {
		return nil, errors.New("you have to provide either email or id")
	}

	rows, err := us.db.Query(
		ctx, query,
		pgx.NamedArgs{"param": param},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.User],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us *AuthStore) GetByRefreshTokenRedis(ctx *fasthttp.RequestCtx, token string) (string, error) {
	var refreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	tokenClaims, err := shared.ValidateToken(
		token,
		refreshTokenPublicKey,
	)
	if err != nil {
		return "", err
	}

	userID, err := us.redis.Get(
		ctx,
		tokenClaims.TokenUUID,
	).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (us *AuthStore) GetByAccessTokenRedis(ctx *fasthttp.RequestCtx, token string) (
	string,
	string,
	error,
) {
	var accessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	tokenClaims, err := shared.ValidateToken(
		token,
		accessTokenPublicKey,
	)
	if err != nil {
		return "", "", err
	}

	userID, err := us.redis.Get(
		ctx,
		tokenClaims.TokenUUID,
	).Result()
	if errors.Is(err, redis.Nil) {
		return "", "", err
	}

	return userID, tokenClaims.TokenUUID, nil
}

func (us *AuthStore) DeleteTokensRedis(ctx *fasthttp.RequestCtx, refreshToken, accessToken string) error {
	var refreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	tokenClaims, err := shared.ValidateToken(
		refreshToken,
		refreshTokenPublicKey,
	)
	if err != nil {
		return err
	}

	_, err = us.redis.Del(
		ctx,
		tokenClaims.TokenUUID,
		accessToken,
	).Result()
	if err != nil {
		return err
	}

	return nil
}
