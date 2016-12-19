package taigaclient

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

// TaigaManager manage interactions with taiga
type TaigaManager struct {
	taigaClient             *taiga.Client
	TaigaProject            string
	Milestone               *taiga.Milestone
	StoriesPerUsers         map[string][]taiga.Userstory
	IssuesPerUsers          map[string][]taiga.Issue
	StoriesDonePerUsers     []taiga.Userstory
	StoriesRejectedPerUsers []taiga.Userstory
	IssuesDonePerUsers      []taiga.Issue
	IssuesRejectedPerUsers  []taiga.Issue
	PointList               map[int]string
	RoleList                map[string]string
}

var (
	usStatusMap     map[string]int
	issuesStatusMap map[string]int
	userList        map[int]string
)

//NewTaigaManager make initialization of taiga client lib
func (t *TaigaManager) NewTaigaManager(taigaUsername *string, taigaPassword *string, taigaProject *string) *TaigaManager {
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", "https://taiga.botsunit.io"))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while initializating taiga client"))
	}
	taigaManager := &TaigaManager{taigaClient: taigaClient, TaigaProject: *taigaProject}
	taigaManager.GetStatusUS()
	taigaManager.GetStatusIssues()
	taigaManager.GetUserList()
	taigaManager.GetPoints()
	taigaManager.GetRoles()
	return taigaManager
}

// GetMilestoneWithDetails return a full milestone detailed
func (t *TaigaManager) GetMilestoneWithDetails(milestoneName string, ch chan bool) {
	mileStone, _, err := t.taigaClient.Milestones.GetMilestoneDetails(milestoneName, t.TaigaProject)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving milestone"))
	}
	t.Milestone = &mileStone
	// ch <- true
}

//GetRoles allow to retrieve role list
func (t *TaigaManager) GetRoles() {
	t.RoleList = make(map[string]string)
	projectID, err := t.taigaClient.Projects.GetProjectID(t.TaigaProject)
	if err != nil {
		fmt.Println("Error while retrieving project ID")
		return
	}
	project, _, err := t.taigaClient.Projects.GetProject(projectID)
	if err != nil {
		fmt.Println("Error while retrieving roles", err.Error())
		return
	}
	for _, role := range project.Roles {
		strRoleID := strconv.Itoa(role.ID)
		t.RoleList[strRoleID] = role.Name
	}
}

//GetPoints Retrieve the points
func (t *TaigaManager) GetPoints() {
	points, _, err := t.taigaClient.Points.ListPoints(&taiga.ListPointsOptions{})
	t.PointList = make(map[int]string)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving points"))
	} else {
		for _, pt := range points {
			t.PointList[pt.ID] = pt.Name
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
		fmt.Println(status.Name)
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
func (t *TaigaManager) MapStoriesPerUsers(status string) {
	t.StoriesPerUsers = make(map[string][]taiga.Userstory)
	for _, us := range t.Milestone.UserStoryList {
		if us.Assigne != 0 && us.Status == usStatusMap[status] {
			t.StoriesPerUsers[userList[us.Assigne]] = append(t.StoriesPerUsers[userList[us.Assigne]], *us)
		}
	}
}

// Sort history entries : filter the one that are status modification only, and take the latest modified today
func getLatestHistoryEntryWithStatusModification(historyEntries []*taiga.HistoryEntry) *taiga.HistoryEntry {
	var historyEntryResult *taiga.HistoryEntry
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, historyEntry := range historyEntries {
		if historyEntry.Diff.Status != nil && len(historyEntry.Diff.Status) == 2 {
			historyModificationYear, historyModificationMonth, historyModificationDay := historyEntry.CreatedAt.Date()
			if historyModificationYear == nowYear && historyModificationMonth == nowMonth && historyModificationDay == nowDay {
				if historyEntryResult == nil {
					historyEntryResult = historyEntry
				} else {
					if historyEntry.CreatedAt.After(historyEntryResult.CreatedAt) {
						historyEntryResult = historyEntry
					}
				}
			}
		}
	}
	return historyEntryResult
}

func (t *TaigaManager) retrieveUserStoryHistory(us taiga.Userstory) (string, string) {
	historyEntries, _, err := t.taigaClient.Userstories.GetUserStoryHistory(us.ID)
	if err != nil {
		fmt.Println("Error while retrieving history", err.Error())
	}
	latestHistoryEntry := getLatestHistoryEntryWithStatusModification(historyEntries)
	fromStatus := latestHistoryEntry.HistoryValueList.Status[0]
	toStatus := latestHistoryEntry.HistoryValueList.Status[1]
	return fromStatus, toStatus
}

func (t *TaigaManager) retrieveIssueHistory(issue taiga.Issue) (string, string) {
	historyEntries, _, err := t.taigaClient.Userstories.GetIssueHistory(issue.ID)
	if err != nil {
		fmt.Println("Error while retrieving history", err.Error())
	}
	latestHistoryEntry := getLatestHistoryEntryWithStatusModification(historyEntries)
	fromStatus := latestHistoryEntry.HistoryValueList.Status[0]
	toStatus := latestHistoryEntry.HistoryValueList.Status[1]
	return fromStatus, toStatus
}

//MapStoriesDonePerUsers make data to storie that have been done today
func (t *TaigaManager) MapStoriesDonePerUsers() {
	var storiesDones []taiga.Userstory
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, us := range t.Milestone.UserStoryList {
		year, month, day := us.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveUserStoryHistory(*us)
			if fromStatus == "Ready for test" && toStatus == "Done" {
				us.AssignedUser = userList[us.Assigne]
				storiesDones = append(storiesDones, *us)
			}
		}
	}
	t.StoriesDonePerUsers = storiesDones
}

//MapStoriesRejectedPerUsers make data to storie that have been done today
func (t *TaigaManager) MapStoriesRejectedPerUsers() {
	var storiesRejected []taiga.Userstory
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, us := range t.Milestone.UserStoryList {
		year, month, day := us.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveUserStoryHistory(*us)
			fmt.Println(fmt.Sprintf("US : %s, fromStatus : %s, toStatus :%s", us.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "In progress" {
				us.AssignedUser = userList[us.Assigne]
				storiesRejected = append(storiesRejected, *us)
			}
		}
	}
	t.StoriesRejectedPerUsers = storiesRejected
}

//MapIssuesDonePerUsers map issue done per users
func (t *TaigaManager) MapIssuesDonePerUsers() {
	var issuesDone []taiga.Issue
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	nowYear, nowMonth, nowDay := time.Now().Date()

	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		year, month, day := issue.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveIssueHistory(*issue)
			fmt.Println(fmt.Sprintf("Issue : %s, FromStatus : %s, toStatus : %s", issue.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "Closed" {
				issue.AssignedUser = userList[issue.Assigne]
				issuesDone = append(issuesDone, *issue)
			}
		}
	}
	t.IssuesDonePerUsers = issuesDone
}

//MapIssuesRejectedPerUsers map issue rejected per users
func (t *TaigaManager) MapIssuesRejectedPerUsers() {
	var issuesRejected []taiga.Issue
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	nowYear, nowMonth, nowDay := time.Now().Date()

	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		year, month, day := issue.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveIssueHistory(*issue)
			fmt.Println(fmt.Sprintf("Issue : %s, FromStatus : %s, toStatus : %s", issue.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "In progress" {
				issue.AssignedUser = userList[issue.Assigne]
				issuesRejected = append(issuesRejected, *issue)
			}
		}
	}
	t.IssuesRejectedPerUsers = issuesRejected
}

//MapIssuesPerUsers retrieve issue in progress and map them per users
func (t *TaigaManager) MapIssuesPerUsers(status string) {
	t.IssuesPerUsers = make(map[string][]taiga.Issue)
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		if issue.Status == issuesStatusMap[status] {
			t.IssuesPerUsers[userList[issue.Assigne]] = append(t.IssuesPerUsers[userList[issue.Assigne]], *issue)
		}
	}
}
