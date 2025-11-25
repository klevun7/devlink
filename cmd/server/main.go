package main
import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	_ "github.com/mattn/go-sqlite3"
)

var (
	apiSecret string
	sesClient *ses.Client
	db *sql.DB
)

func main() {
	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		log.Fatal("API_SECRET environment variable is required, check file for valid key")
	}
	var err error
	db, err = sql.Open("sqlite3", "./data/devlink.db")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	initDB()
}

func initDB() {}