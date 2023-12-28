package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvironmentVariables struct {
	PortNumber  int
	DatabaseUrl string
	SecretKey   string
}

var Environment *EnvironmentVariables

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Unable to load environment file")
	}

	portNumber, err := strconv.Atoi(os.Getenv("PORT_NUMBER"))
	dbUrl := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("SECRET_KEY")
	if err != nil {
		log.Fatalln("Unable to load port number for server")
	}

	Environment = &EnvironmentVariables{
		PortNumber:  portNumber,
		DatabaseUrl: dbUrl,
		SecretKey:   secretKey,
	}
}
