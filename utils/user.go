package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Dom-HTG/warp/models"
)

func GetUserProfile(accessToken string, ctx context.Context) (*models.UserProfile, error) {
	apiURL := os.Getenv("API_ADDRESS")
	var profileURL string = fmt.Sprintf("%s/v1/me", apiURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, profileURL, nil)
	if err != nil {
		return nil, err
	}

	var auth string = fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", auth)

	Client := &http.Client{}
	resp, err1 := Client.Do(req)
	if err1 != nil {
		return nil, err1
	}
	defer resp.Body.Close()

	respBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	var ProfileResponse *models.UserProfile
	if err3 := json.Unmarshal(respBytes, &ProfileResponse); err3 != nil {
		return nil, err3
	}

	return ProfileResponse, nil
}
