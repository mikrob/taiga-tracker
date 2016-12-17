package main

import (
	"flag"
	"fmt"
	taiga "taiga-gitlab/taiga"
	"taiga_tracker/lib"
)

var (
	taigaUsername = flag.String("u", "api", "taiga username")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
)

// TaigaManager manage interactions with taiga
type TaigaManager struct {
	taigaClient *taiga.Client
	Milestone   *taiga.Milestone
}

//NewTaigaManager make initialization of taiga client lib
func (t *TaigaManager) NewTaigaManager() *TaigaManager {
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", "https://taiga.botsunit.io"))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}
	return &TaigaManager{taigaClient: taigaClient}
}

// StoriesForMileStone return all stories for a given milestone
func (t *TaigaManager) StoriesForMileStone(milestone string) ([]*taiga.Userstory, error) {
	mileStoneList, _, err := t.taigaClient.Milestones.ListMilestones()
	if err != nil {
		fmt.Printf("Error while retrieving milestone : %s", err.Error())
	}
	milestoneFiltered := lib.FilterMilestone(mileStoneList, func(m *taiga.Milestone) bool {
		return m.Name == milestone
	})

	var userStoryList []*taiga.Userstory
	if milestoneFiltered != nil {
		userStoryListResult, _, errUserStoryList := t.taigaClient.Userstories.ListUserstoriesForMilestone(*milestoneFiltered[0])
		if errUserStoryList != nil {
			fmt.Printf("Error while retrieving US : %s", errUserStoryList.Error())
		}
		userStoryList = userStoryListResult
	} else {
		fmt.Printf("Error while retrieving milestone, no milestone matching %s", milestone)
	}

	return userStoryList, err
}

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers() {

}

func main() {
	flag.Parse()
	var taigaManager *TaigaManager
	taigaManager = (&TaigaManager{}).NewTaigaManager()
	userStoryList, err := taigaManager.StoriesForMileStone("0.6")
	if err != nil {
		fmt.Printf("Error with retrieving stories for milestone %s", err.Error())
	} else {
		for _, story := range userStoryList {
			fmt.Println(story.Subject)
		}
	}

}
