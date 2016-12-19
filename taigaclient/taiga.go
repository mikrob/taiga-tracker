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
	usStatusMap     map[string]int
	issuesStatusMap map[string]int
	userList        map[int]string
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
	taigaManager.GetStatusIssues()
	taigaManager.GetUserList()
	return taigaManager
}

// GetMilestoneWithDetails return a full milestone detailed
func (t *TaigaManager) GetMilestoneWithDetails(milestoneName string, projectName string, ch chan bool) {
	mileStone, _, err := t.taigaClient.Milestones.GetMilestoneDetails(milestoneName, projectName)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving milestone"))
	}
	t.Milestone = &mileStone
	// ch <- true
}

//MapIssuesPerUsers retrieve issue in progress and map them per users
func (t *TaigaManager) MapIssuesPerUsers() {
	t.IssuesPerUsers = make(map[string][]taiga.Issue)
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		if issue.Status == issuesStatusMap["In progress"] {
			t.IssuesPerUsers[userList[issue.Assigne]] = append(t.IssuesPerUsers[userList[issue.Assigne]], *issue)
		}
	}
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

//GetStatusIssues retrieve the users stories status kind
func (t *TaigaManager) GetStatusIssues() {
	statusList, _, err := t.taigaClient.Issues.ListIssueStatuses()
	issuesStatusMap = make(map[string]int)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
	}
	for _, status := range statusList {
		issuesStatusMap[status.Name] = status.ID
	}
}

//GetStatusUS retrieve the users stories status kind
func (t *TaigaManager) GetStatusUS() {
	statusList, _, err := t.taigaClient.Issues.ListUserstoryStatuses()
	usStatusMap = make(map[string]int)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
	}
	for _, status := range statusList {
		usStatusMap[status.Name] = status.ID
	}
}

//GetUserList initialize userlist
func (t *TaigaManager) GetUserList() {
	users, errUserList := t.RetrieveUserList()
	if errUserList != nil {
		fmt.Println(fmt.Errorf("Error while retrieving list of users"))
	} else {
		userList = users
	}
}

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers() {
	t.StoriesPerUsers = make(map[string][]taiga.Userstory)
	for _, us := range t.Milestone.UserStoryList {
		if us.Assigne != 0 && us.Status == usStatusMap["In progress"] {
			t.StoriesPerUsers[userList[us.Assigne]] = append(t.StoriesPerUsers[userList[us.Assigne]], *us)
		}
	}

}
