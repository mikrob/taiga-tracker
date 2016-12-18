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
	milestone, err := taigaManager.GetMilestoneWithDetails("0.6", "Ufancyme")
	if err != nil {
		fmt.Println("Error while retrieving milestone:", err.Error())
	}
	fmt.Println("Milestone content : ")
	for _, us := range milestone.UserStoryList {
		fmt.Println(us.Subject)
	}
	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
}
