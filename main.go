package main

import (
	"flag"
	"fmt"
	taiga "taiga-gitlab/taiga"
)

var (
	taigaUsername = flag.String("u", "api", "taiga username")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
)

func main() {
	flag.Parse()
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", "https://taiga.botsunit.io"))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}

	userList, _, err := taigaClient.Users.ListUsers()
	if err != nil {
		fmt.Printf("Error while retrieving users : %s", err.Error())
	}

	for _, user := range userList {
		fmt.Println(user.FullName)
	}

	mileStoneList, _, err := taigaClient.Milestones.ListMilestones()

	if err != nil {
		fmt.Printf("Error while retrieving milestone : %s", err.Error())
	}

	milestoneWaned := "0.6"
	var milestoneStruct taiga.Milestone
	for _, milestone := range mileStoneList {
		fmt.Println(milestone.Name)
		if milestone.Name == milestoneWaned {
			milestoneStruct = *milestone
		}
	}

	userStoryList, _, err := taigaClient.Userstories.ListUserstoriesForMilestone(milestoneStruct)

	if err != nil {
		fmt.Printf("Error while retrieving US : %s", err.Error())
	}

	for _, us := range userStoryList {
		fmt.Println(us.Subject)
	}

}
