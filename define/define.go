package define

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetMailPwd() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("load env failed!")
		return "", err
	}
	MailPwd := os.Getenv("MailPwd")

	return MailPwd, err
}
