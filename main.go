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
	taigaUsername = flag.String("l", "login", "taiga login")
	taigaPassword = flag.String("p", "password", "taiga password")
	taigaProject  = flag.String("t", "project", "taiga project")
	taigaURL      = flag.String("d", "https://taiga.example.io", "taiga URL")
	taigaManager  *taigaclient.TaigaManager
)

func demosDatas(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()

	taigaManager.MapStoriesDemos()
	taigaManager.MapIssuesDemos()

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run \n", elapsed)
	c.HTML(http.StatusOK, "demo.tmpl", gin.H{
		"title":            "Demo Short View",
		"userStories":      taigaManager.StoriesDemos,
		"pointList":        taigaManager.PointList,
		"roleList":         taigaManager.RoleList,
		"issues":           taigaManager.IssuesDemos,
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
	taigaManager.MapStoriesWipPerUsers()
	taigaManager.MapIssuesWipPerUsers()

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
	elapsed := time.Since(start)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":            "Taiga Tracker",
		"currentMilestone": taigaManager.CurrentMileStone,
		"time":             elapsed,
		"milestoneList":    taigaManager.MilestoneList,
		"taigaURL":         taigaURL,
	})
}

//synchronizeStories set US to to test if all task are ready to test
// or set US to closed if all task are done
func synchronizeStories(c *gin.Context) {
	start := time.Now()
	taigaManager.GetMilestoneWithDetails()
	usListSync, err := taigaManager.SynchronizeMilestone()
	if err != nil {
		fmt.Println("Error while syncinc :", err.Error())
	}

	elapsed := time.Since(start)
	c.HTML(http.StatusOK, "tasksSync.tmpl", gin.H{
		"title":            "Task Synchronizer",
		"usListSync":       usListSync,
		"time":             elapsed,
		"taigaURL":         taigaURL,
		"currentMilestone": taigaManager.CurrentMileStone,
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
	authorized.POST("/synchronizeStories", synchronizeStories)
	router.Run(":8282")

}
