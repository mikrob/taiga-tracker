package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"taiga-tracker/taigaclient"

	"github.com/gin-gonic/gin"
)

var (
	taigaUsername = flag.String("l", "api", "taiga login")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
	taigaProject  = flag.String("t", "Ufancyme", "taiga project")
	taigaURL      = flag.String("d", "https://taiga.botsunit.io", "taiga URL")
	taigaManager  *taigaclient.TaigaManager
)

func demosDatas(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()

	taigaManager.MapStoriesPerUsers("Ready for test")
	taigaManager.MapIssuesPerUsers("Ready for test")
	fmt.Println("DEMOS CURRENT MILESTONE IS", taigaManager.CurrentMileStone)

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "demo.tmpl", gin.H{
		"title":            "Demo Short View",
		"userStories":      taigaManager.StoriesPerUsers,
		"pointList":        taigaManager.PointList,
		"roleList":         taigaManager.RoleList,
		"issues":           taigaManager.IssuesPerUsers,
		"currentMilestone": taigaManager.CurrentMileStone,
		"time":             elapsed,
		"taigaURL":         taigaURL,
	})
}

func overtakingUSDatas(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()
	taigaManager.GetStoriesAndElapsedTime()
	taigaManager.TimeTrackStories()
	elapsed := time.Since(start)
	fmt.Println("OVER CURRENT MILESTONE IS", taigaManager.CurrentMileStone)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "overtake.tmpl", gin.H{
		"title":                 "Overtaking US",
		"userStoriesOvertaking": taigaManager.StoriesTimeTrackedPerUsers,
		"currentMilestone":      taigaManager.CurrentMileStone,
		"time":                  elapsed,
		"taigaURL":              taigaURL,
	})
}

func demosCRDatas(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()
	taigaManager.MapStoriesDonePerUsers()
	taigaManager.MapStoriesRejectedPerUsers()
	taigaManager.MapIssuesDonePerUsers()
	taigaManager.MapIssuesRejectedPerUsers()
	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "cr.tmpl", gin.H{
		"title":               "Demo CR",
		"userStoriesDone":     taigaManager.StoriesDonePerUsers,
		"userStoriesRejected": taigaManager.StoriesRejectedPerUsers,
		"issuesDone":          taigaManager.IssuesDonePerUsers,
		"issuesRejected":      taigaManager.IssuesRejectedPerUsers,
		"currentMilestone":    taigaManager.CurrentMileStone,
		"time":                elapsed,
		"taigaURL":            taigaURL,
	})

}

func wipDatas(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()
	taigaManager.MapStoriesPerUsers("In progress")
	taigaManager.MapIssuesPerUsers("In progress")

	fmt.Println("WIP CURRENT MILESTONE IS", taigaManager.CurrentMileStone)

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "wip.tmpl", gin.H{
		"title":            "Work In Progress Short View",
		"userStories":      taigaManager.StoriesPerUsers,
		"pointList":        taigaManager.PointList,
		"roleList":         taigaManager.RoleList,
		"issues":           taigaManager.IssuesPerUsers,
		"currentMilestone": taigaManager.CurrentMileStone,
		"time":             elapsed,
		"taigaURL":         taigaURL,
	})
}

func indexMenu(c *gin.Context) {
	start := time.Now()
	taigaManager.ListMilestones()
	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	fmt.Println("GET INDEX CURRENT MILESTONE IS", taigaManager.CurrentMileStone)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":            "Taiga Tracker",
		"milestoneList":    taigaManager.MilestoneList,
		"currentMilestone": taigaManager.CurrentMileStone,
		"time":             elapsed,
		"taigaURL":         taigaURL,
	})
}

func postMilestone(c *gin.Context) {
	start := time.Now()
	milestone := c.PostForm("milestone")
	milestone = strings.TrimSpace(milestone)
	taigaManager.CurrentMileStone = milestone
	taigaManager.ListMilestones()
	fmt.Println("POST INDEX CURRENT MILESTONE IS", taigaManager.CurrentMileStone)
	elapsed := time.Since(start)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":            "Taiga Tracker",
		"currentMilestone": taigaManager.CurrentMileStone,
		"time":             elapsed,
		"milestoneList":    taigaManager.MilestoneList,
		"taigaURL":         taigaURL,
	})
}

func main() {
	flag.Parse()
	taigaManager = (&taigaclient.TaigaManager{}).NewTaigaManager(taigaUsername, taigaPassword, taigaProject, taigaURL)

	router := gin.Default()
	//router.HTMLRender = createMyRender()
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/fonts", "./fonts")
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin":  "password",
		"mikrob": "password",
		"benj":   "password",
	}))

	router.LoadHTMLGlob("templates/*")

	authorized.GET("/wip", wipDatas)
	authorized.GET("/demo", demosDatas)
	authorized.GET("/cr", demosCRDatas)
	authorized.GET("/over", overtakingUSDatas)
	authorized.GET("/", indexMenu)
	authorized.POST("/", postMilestone)
	router.Run(":8282")

}
