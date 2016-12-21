package taigaclient

import (
	"fmt"
	"strconv"
	"sync"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
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

//getElapsedTimeAttributeId find attribute id for elapsed times
func (t *TaigaManager) getElapsedTimeAttributeID() int {
	customAttributes, _, err := t.taigaClient.Userstories.GetUserStoryCustomAttributes()
	result := 1
	if err != nil {
		fmt.Println("Error while retrieving custom attribute list")
	}
	for _, ca := range customAttributes {
		if ca.Name == "Elapsed Time" {
			return ca.ID
		}
	}
	return result
}

//GetAttributeValue retrieve attribute value in a goproc with waitgroup to allow parralel call
func (t *TaigaManager) GetAttributeValue(us *taiga.Userstory, attributeID int, wg *sync.WaitGroup, idx int) {
	defer wg.Done()
	attributesValues, _, err := t.taigaClient.Userstories.GetUserStoryCustomAttributeValue(us.ID)
	if err != nil {
		fmt.Println("Error while retrieving custom user story attributes value ", err.Error())
	}
	if len(attributesValues.Values) != 0 {
		elapsedTime := attributesValues.Values[strconv.Itoa(attributeID)]
		elapsedTimeFloat, err := strconv.ParseFloat(elapsedTime, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Cannot convert %v to float64", elapsedTime))
			fmt.Println(err.Error())
		} else {
			us.ElapsedTime = elapsedTimeFloat
		}
	}
}

//GetStoriesAndElapsedTime iterate over stories and detect the one that are overtaking
func (t *TaigaManager) GetStoriesAndElapsedTime() {
	elapsedTimeAttributeID := t.getElapsedTimeAttributeID()
	t.StoriesTimeTrackedPerUsers = make(map[string][]*taiga.Userstory)
	var wg sync.WaitGroup
	for idx, us := range t.Milestone.UserStoryList {
		if us.Assigne != 0 && us.Status == usStatusMap["Done"] {
			t.StoriesTimeTrackedPerUsers[userList[us.Assigne]] = append(t.StoriesTimeTrackedPerUsers[userList[us.Assigne]], us)
			wg.Add(1)
			go t.GetAttributeValue(us, elapsedTimeAttributeID, &wg, idx)
		}
	}
	wg.Wait()
}

//TimeTrackStories timetrack the stories
func (t *TaigaManager) TimeTrackStories() {
	points, _, err := t.taigaClient.Points.ListPoints(&taiga.ListPointsOptions{})
	if err != nil {
		fmt.Println("Error while retrieving points", err.Error())
	}
	pointListFloat := make(map[int]float64)
	for _, point := range points {
		pointListFloat[point.ID] = point.Value
	}
	//roleList := t.RoleList
	for _, usList := range t.StoriesTimeTrackedPerUsers {
		for _, us := range usList {
			fmt.Println("ElapsedTime ", us.ElapsedTime)
			usTotalPoints := 0.0
			for _, pointID := range us.Points {
				usTotalPoints += pointListFloat[pointID]
			}
			fmt.Println("US TOTAL POINTS :", usTotalPoints)
			us.TotalPoint = usTotalPoints
			if us.ElapsedTime == 0 || usTotalPoints == 0 {
				us.Color = "card-panel deep-orange lighten-2"
			} else if us.ElapsedTime > usTotalPoints {
				us.Overtaking = true
				us.Color = "card-panel red darken-4"
			} else if us.ElapsedTime < usTotalPoints {
				us.Undertaking = true
				us.Color = "card-panel teal lighten-1"
			} else if (usTotalPoints == us.ElapsedTime) && usTotalPoints != 0 {
				us.RightTime = true
				us.Color = "card-panel blue lighten-2"
			}
		}
	}
}
