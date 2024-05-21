package worker

import (
	"fmt"
	"gin/cfapi"
	"gin/models"
	"gin/store"
	"log"
	"sync"
	"time"
)

func PerformWork(m *store.MongoStore, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		m.OpenConnectionWithMongoDB()
		obj := new(cfapi.CodeforcesClient)
		RecentActions, err := obj.RecentActions(100)
		if err != nil {
			fmt.Println("error occurred")
			return
		}

		maxTimeStamp, err := m.GetMaxTimeStamp()
		if err != nil {
			log.Printf("Error while getting maxTimeStamp: %v", maxTimeStamp)
		}

		log.Printf("Got maxTimeStamp successfully")

		var NewData []models.RecentAction

		for i := 0; i < len(RecentActions); i++ {
			if RecentActions[i].TimeSeconds > int64(maxTimeStamp) {
				NewData = append(NewData, RecentActions[i])
			}
		}

		log.Printf("RecentActions stored in NewData successfully ")

		err = m.StoreRecentActionsInTheDatabase(NewData)
		if err != nil {
			log.Printf("Error occurred while storing data : %v", err)
			return
		}

		var temp []models.RecentAction

		NewData = temp

		log.Printf("The worker will sleep for 5 min now.")
		time.Sleep(5 * time.Minute)
	}
}
