package taigaclient

import (
	"fmt"
	"taiga-gitlab/taiga"
)

// TaigaManager manage interactions with taiga
type TaigaManager struct {
	taigaClient     *taiga.Client
	Milestone       *taiga.Milestone
	StoriesPerUsers map[string][]taiga.Userstory
	IssuesPerUsers  map[string][]taiga.Issue
}

var (
	statusMap map[string]int
)

//NewTaigaManager make initialization of taiga client lib
func (t *TaigaManager) NewTaigaManager(taigaUsername *string, taigaPassword *string) *TaigaManager {
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", "https://taiga.botsunit.io"))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while initializating taiga client"))
	}
	taigaManager := &TaigaManager{taigaClient: taigaClient}
	taigaManager.GetStatusUS()
	return taigaManager
}

// GetMilestoneWithDetails return a full milestone detailed
func (t *TaigaManager) GetMilestoneWithDetails(milestoneName string, projectName string) {
	mileStone, _, err := t.taigaClient.Milestones.GetMilestoneDetails(milestoneName, projectName)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving milestone"))
	}
	t.Milestone = &mileStone
}

//RetrieveUserList make a list of user mapped as id -> fullname
func (t *TaigaManager) RetrieveUserList() (map[int]string, error) {
	var userMap map[int]string
	userMap = make(map[int]string)
	userList, _, err := t.taigaClient.Users.ListUsers()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
		return nil, err
	}
	for _, user := range userList {
		userMap[user.ID] = user.FullName
	}
	return userMap, nil
}

//GetStatusUS retrieve the users stories status kind
func (t *TaigaManager) GetStatusUS() {
	statusList, _, err := t.taigaClient.Issues.ListUserstoryStatuses()
	statusMap = make(map[string]int)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
	}
	for _, status := range statusList {
		statusMap[status.Name] = status.ID
	}
}

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers() {
	userList, errUserList := t.RetrieveUserList()
	if errUserList != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
	} else {
		t.StoriesPerUsers = make(map[string][]taiga.Userstory)
		for _, us := range t.Milestone.UserStoryList {
			if us.Assigne != 0 && us.Status == statusMap["In progress"] {
				t.StoriesPerUsers[userList[us.Assigne]] = append(t.StoriesPerUsers[userList[us.Assigne]], *us)
			}
		}
	}
}
