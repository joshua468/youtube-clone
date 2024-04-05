package models

import "os"

type Env struct {
	DBPassword            string
	DBName                string
	DBUsername            string
	DBHost                string
	DBPort                string
	JWTAccessTokenExpiry  string
	JWTRefreshTokenExpiry string
	JWTSigningSecret      string
	PORT                  string
}

func NewEnv() *Env {
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB_NAME")
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	jwtAccessTokenExpiry := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY")
	jwtRefreshTokenExpiry := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY")
	jwtSigningSecret := os.Getenv("JWT_SIGNING_SECRET")
	port := os.Getenv("PORT")

	return &Env{
		DBPassword:            dbPass,
		DBName:                dbName,
		DBUsername:            dbUsername,
		DBHost:                dbHost,
		DBPort:                dbPort,
		JWTAccessTokenExpiry:  jwtAccessTokenExpiry,
		JWTRefreshTokenExpiry: jwtRefreshTokenExpiry,
		JWTSigningSecret:      jwtSigningSecret,
		PORT:                  port,
	}
}
