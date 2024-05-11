package define

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const CACHE_PREFIX = "EmaileCode"

const V1API_HEADER_TOEKN = "token"

func GetMailPwd() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("load env failed!")
		return "", err
	}
	MailPwd := os.Getenv("MailPwd")

	return MailPwd, err
}
