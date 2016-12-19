package main

import (
	"flag"
	"fmt"
	"taiga_tracker/taigaclient"
	"time"
)

var (
	taigaUsername = flag.String("u", "api", "taiga username")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
)

func main() {
	flag.Parse()
	start := time.Now()
	var taigaManager *taigaclient.TaigaManager
	taigaManager = (&taigaclient.TaigaManager{}).NewTaigaManager(taigaUsername, taigaPassword)
	taigaManager.GetMilestoneWithDetails("0.5", "Ufancyme")
	taigaManager.MapStoriesPerUsers()

	for user, userStories := range taigaManager.StoriesPerUsers {
		fmt.Println("================================================================================================")
		fmt.Println("User : ", user)
		fmt.Println("Stories : ")
		for _, us := range userStories {
			fmt.Println(us.Subject)
		}
		fmt.Println("================================================================================================")
	}

	taigaManager.GetStatusUS()
	// fmt.Println("Milestone content : ")
	// for _, us := range milestone.UserStoryList {
	// 	fmt.Println(us.Subject)
	// }
	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
}
