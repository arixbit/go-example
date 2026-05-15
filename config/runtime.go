package config

import "github.com/joho/godotenv"

// LoadEnv preloads optional dotenv files without overriding existing environment variables.
func LoadEnv(paths ...string) {
	files := append([]string{}, paths...)
	files = append(files, ".env")

	for _, file := range files {
		_ = godotenv.Load(file)
	}
}
