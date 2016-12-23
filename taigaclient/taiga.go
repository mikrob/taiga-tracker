package taigaclient

import (
	"fmt"
	"strconv"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

const (
	DefaultMilestone = "0.6"
)

// TaigaManager manage interactions with taiga
type TaigaManager struct {
	taigaClient                *taiga.Client
	TaigaProject               string
	Milestone                  *taiga.Milestone
	StoriesPerUsers            map[string][]taiga.Userstory
	StoriesTimeTrackedPerUsers map[string][]*taiga.Userstory
	IssuesPerUsers             map[string][]taiga.Issue
	StoriesDonePerUsers        []taiga.Userstory
	StoriesRejectedPerUsers    []taiga.Userstory
	IssuesDonePerUsers         []taiga.Issue
	IssuesRejectedPerUsers     []taiga.Issue
	PointList                  map[int]string
	RoleList                   map[string]string
	MilestoneList              []*taiga.Milestone
	CurrentMileStone           string
}

var (
	usStatusMap     map[string]int
	issuesStatusMap map[string]int
	userList        map[int]string
)

//NewTaigaManager make initialization of taiga client lib
func (t *TaigaManager) NewTaigaManager(taigaUsername *string, taigaPassword *string, taigaProject *string, taigaURL *string) *TaigaManager {
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", *taigaURL))
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

//ListMilestones allow to list existing milestone
func (t *TaigaManager) ListMilestones() {
	milestoneList, _, err := t.taigaClient.Milestones.ListMilestones()
	t.CurrentMileStone = DefaultMilestone
	if err != nil {
		fmt.Println("Error while listing milestone ", err.Error())
	}
	t.MilestoneList = milestoneList
}

// GetMilestoneWithDetails return a full milestone detailed
func (t *TaigaManager) GetMilestoneWithDetails() {
	milestoneName := t.CurrentMileStone
	if milestoneName == "" {
		milestoneName = DefaultMilestone
		t.CurrentMileStone = milestoneName
	}
	mileStone, _, err := t.taigaClient.Milestones.GetMilestoneDetails(milestoneName, t.TaigaProject)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving milestone"))
	}
	t.Milestone = &mileStone
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
