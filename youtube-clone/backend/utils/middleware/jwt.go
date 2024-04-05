package middlewares

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"

	"github.com/joshua468/youtube-clone/backend/utils/models"
)

const (
	identityKey     = "id"
	realm           = "youtube-clone"
	claimsID        = "id"
	isAdminClaims   = "is_admin"
	claimsExpiry    = "exp"
	claimsCreatedAt = "orig_iat"
)

var (
	// ErrUnexpectedSigningMethod occurs when a token does not conform to the expected signing method
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrInvalidToken indicates JWT token has expired. Can't refresh.
	ErrInvalidToken = errors.New("token is invalid/expired")

	// ErrMissingToken reports missing auth token
	ErrMissingToken = errors.New("auth token is missing from the header")

	// ErrUnauthorized reports unauthorized user
	ErrUnauthorized = errors.New("you are not authorized")

	// ErrInvalidTokenHeaderFormat e.g when client passes the token without the Bearer prefix
	ErrInvalidTokenHeaderFormat = errors.New("invalid token header format, token type must be bearer")
)

type Middleware struct {
	jwt  *JwtConfig
	pKey interface{} // Public key for asymmetric signing algorithm
	sKey interface{} // Secret key for symmetric signing algorithm
}

type JwtConfig struct {
	SigningAlgorithm string // Algorithm used for signing (e.g., HS256, RS256)
	Key              []byte // Secret key for symmetric signing algorithm
	CookieDomain     string // Domain for setting cookies
	SecureCookie     bool   // Whether cookies are secure
	CookieHTTPOnly   bool   // Whether cookies are HTTP-only
}

type Tokens struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  string
	RefreshTokenExpiry string
}

func jwtAccessTokenExpiry(env *models.Env) time.Duration {
	ttl, err := strconv.Atoi(env.JWTAccessTokenExpiry)
	if err != nil {
		return time.Minute * 10000
	}
	return time.Minute * time.Duration(ttl)
}

func jwtRefreshTokenExpiry(env *models.Env) time.Duration {
	ttl, err := strconv.Atoi(env.JWTRefreshTokenExpiry)
	if err != nil {
		return time.Hour * 240000
	}
	return time.Hour * time.Duration(ttl)
}

// CreateToken creates a new user access and refresh tokens
func (m *Middleware) CreateToken(env *models.Env, userID string, isAdmin bool) (*Tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod(m.jwt.SigningAlgorithm), jwt.MapClaims{
		claimsID:        userID,
		claimsExpiry:    time.Now().Add(jwtAccessTokenExpiry(env)).Unix(),
		claimsCreatedAt: time.Now().Unix(),
		isAdminClaims:   isAdmin,
	})

	accessTokenString, err := accessToken.SignedString(m.jwt.Key)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod(m.jwt.SigningAlgorithm), jwt.MapClaims{
		claimsID:        userID,
		claimsExpiry:    time.Now().Add(jwtRefreshTokenExpiry(env)).Unix(),
		claimsCreatedAt: time.Now().Unix(),
		isAdminClaims:   isAdmin,
	})

	refreshTokenString, err := refreshToken.SignedString(m.jwt.Key)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		AccessTokenExpiry:  time.Now().Add(jwtAccessTokenExpiry(env)).String(),
		RefreshTokenExpiry: time.Now().Add(jwtRefreshTokenExpiry(env)).String(),
	}, nil
}

// ValidateRefreshToken validates the refresh token
func (m *Middleware) ValidateRefreshToken(z zerolog.Logger, env *models.Env, token string) (*uuid.UUID, error) {
	tokenGotten, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			z.Error().Msgf("RefreshToken unexpected signing method: (%v)", token.Header["alg"])
			return nil, ErrUnexpectedSigningMethod
		}
		return m.jwt.Key, nil
	})

	//any error may be due to token expiration
	if err != nil {
		z.Err(err).Msg("RefreshToken error")
		return nil, err
	}

	//is token valid?
	if err = tokenGotten.Claims.Valid(); err != nil {
		z.Err(err).Msg("RefreshToken failed :: invalid token")
		return nil, err
	}

	claims, ok := tokenGotten.Claims.(jwt.MapClaims)
	claimsUUID := claims[claimsID].(string)

	if ok && tokenGotten.Valid {
		//convert the interface to uuid.UUID
		parsedUUID, err := uuid.Parse(claimsUUID)
		if err != nil {
			z.Err(err).Msgf("RefreshToken error::(%v)", err)
			return nil, ErrInvalidToken
		}

		return &parsedUUID, nil
	}

	return nil, ErrInvalidToken
}

// ParseToken checks if token is valid and parses it
func (m *Middleware) ParseToken(env *models.Env, tokenStr string) (userID string, isAdmin bool, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwt.Key, nil
	})

	if token == nil {
		m.logger.Error().Str("token", tokenStr).Msg("unable to parse token - token is most likely not valid")
		return userID, isAdmin, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims[claimsID].(string)

		if found, ok := claims[isAdminClaims]; ok {
			isAdmin = found.(bool)
		}

		return userID, isAdmin, nil
	}

	return userID, isAdmin, err
}
