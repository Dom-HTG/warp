package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Dom-HTG/warp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetAccessToken(authCode string) (*models.AccessTokenPayload, error) {
	exURL := os.Getenv("EXCHANGE_URL")
	redirectURI := os.Getenv("REDIRECT_URI")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	formData := map[string]string{
		"grant_type":   "authorization_code",
		"code":         authCode,
		"redirect_uri": redirectURI,
	}

	//form object.
	form := url.Values{}

	for k, v := range formData {
		form.Set(k, v)
	}

	encodedForm := form.Encode()
	body := strings.NewReader(encodedForm)

	//Make http request to url.
	req, err := http.NewRequest(http.MethodPost, exURL, body)
	if err != nil {
		return nil, err
	}

	//Set headers.
	req.Header.Set("Content-Type", " application/x-www-form-urlencoded")
	credentials := fmt.Sprintf("%s:%s", clientID, clientSecret)
	encString := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encString))

	//send request.
	client := &http.Client{}
	response, err1 := client.Do(req)
	if err1 != nil {
		return nil, err1
	}
	defer response.Body.Close()

	respBytes, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		return nil, err3
	}
	fmt.Print(string(respBytes))

	payload := &models.AccessTokenPayload{}
	// err2 := json.NewDecoder(response.Body).Decode(&payload)
	// if err2 != nil {
	// 	return &models.AccessTokenPayload{}, err2
	// }

	if err4 := json.Unmarshal(respBytes, &payload); err4 != nil {
		return nil, err4
	}

	return payload, nil
}

func GenerateState() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func GetStateDB(db *gorm.DB, id uint) (string, error) {
	var userData models.User
	tx := db.Where("id = ?", id).First(&userData)
	if tx.Error != nil {
		return "", tx.Error
	}
	return userData.StateValue, nil
}

func InitDB() (*gorm.DB, error) {
	//construct postgres URL.
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s", host, user, password, dbname, sslmode)

	//start the database.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err1 := db.AutoMigrate(&models.User{}); err1 != nil {
		return nil, err1
	}
	fmt.Print("mode migration success")
	return db, nil
}
