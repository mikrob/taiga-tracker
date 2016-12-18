package taigaclient

import (
	"fmt"
	"taiga-gitlab/taiga"
)

// TaigaManager manage interactions with taiga
type TaigaManager struct {
	taigaClient *taiga.Client
	Milestone   *taiga.Milestone
}

//NewTaigaManager make initialization of taiga client lib
func (t *TaigaManager) NewTaigaManager(taigaUsername *string, taigaPassword *string) *TaigaManager {
	taigaClient := taiga.NewClient(nil, *taigaUsername, *taigaPassword)
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", "https://taiga.botsunit.io"))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}
	return &TaigaManager{taigaClient: taigaClient}
}

// GetMilestoneWithDetails return a full milestone detailed
func (t *TaigaManager) GetMilestoneWithDetails(milestoneName string, projectName string) (*taiga.Milestone, error) {
	milestone, _, err := t.taigaClient.Milestones.GetMilestoneDetails(milestoneName, projectName)
	return &milestone, err
}

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers() {

}
