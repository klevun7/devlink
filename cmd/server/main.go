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