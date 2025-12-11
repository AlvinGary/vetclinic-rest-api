package utils

import (
	"os"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))
