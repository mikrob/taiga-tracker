package taigaclient

import (
	"fmt"
	"strconv"
	"sync"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

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
		if us.Status == usStatusMap["Done"] {
			assign := ""
			if us.Assigne == 0 || userList[us.Assigne] == "" {
				assign = NotAssigned
			} else {
				assign = userList[us.Assigne]
			}
			us.AssignedUser = assign
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
			usTotalPoints := 0.0
			for _, pointID := range us.Points {
				usTotalPoints += pointListFloat[pointID]
			}
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
