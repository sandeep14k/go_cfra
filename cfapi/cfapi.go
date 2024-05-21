package cfapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	models "gin/models"
)

type CodeforcesClient struct {
	client http.Client
}

func (cfClient *CodeforcesClient) RecentActions(maxCount int) ([]models.RecentAction, error) {
	api := "https://codeforces.com/api/recentActions?maxCount=" + strconv.Itoa(maxCount)

	resp, err := cfClient.client.Get(api)
	if err != nil {
		log.Printf("Error occured while calling cf api: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	document, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error occurred while reading the resp body")
		return nil, err
	}

	wrapper := struct {
		Status string
		Result []models.RecentAction
	}{}

	if err = json.Unmarshal(document, &wrapper); err != nil {
		log.Printf("Error while unmarshalling document from cfapi : %v", err)
	}

	return wrapper.Result, err
}
