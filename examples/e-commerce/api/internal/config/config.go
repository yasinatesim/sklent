package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv            string
	APIPort           string
	DatabaseURL       string
	JWTSecret         string
	CryptoKey         string
	CookieDomain      string
	CORSAllowedOrigin []string
	ChromaURL         string
	OpenRouterKey     string
	OpenRouterModel   string
	IyzicoBaseURL     string
	FrontendBaseURL   string
	GIBBaseURL        string
	SMTPHost          string
	SMTPPort          string
	MailFrom          string
}

func Load() Config {
	return Config{
		AppEnv:            env("APP_ENV", "development"),
		APIPort:           env("API_PORT", "8100"),
		DatabaseURL:       env("DATABASE_URL", "postgres://vela:vela@localhost:5532/vela?sslmode=disable"),
		JWTSecret:         env("JWT_SECRET", "change-me-dev-only-secret"),
		CryptoKey:         env("CRYPTO_KEY", ""),
		CookieDomain:      env("COOKIE_DOMAIN", ""),
		CORSAllowedOrigin: splitCSV(env("CORS_ALLOWED_ORIGINS", "http://localhost:3100")),
		ChromaURL:         env("CHROMA_URL", ""),
		OpenRouterKey:     env("OPENROUTER_API_KEY", ""),
		OpenRouterModel:   env("OPENROUTER_MODEL", "minimax/minimax-01"),
		IyzicoBaseURL:     env("IYZICO_BASE_URL", "https://sandbox-api.iyzipay.com"),
		FrontendBaseURL:   env("FRONTEND_BASE_URL", "http://localhost:3100"),
		GIBBaseURL:        env("GIB_BASE_URL", "https://earsivportaltest.efatura.gov.tr"),
		SMTPHost:          env("SMTP_HOST", ""),
		SMTPPort:          env("SMTP_PORT", "587"),
		MailFrom:          env("MAIL_FROM", "Vela Commerce <no-reply@vela.test>"),
	}
}

func (c Config) IsProduction() bool { return c.AppEnv == "production" }

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func envInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

var _ = envInt
