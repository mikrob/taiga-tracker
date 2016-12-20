package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"gitlab.botsunit.com/infra/taiga-tracker/taigaclient"

	"github.com/gin-gonic/gin"
)

var (
	taigaUsername  = flag.String("l", "api", "taiga login")
	taigaPassword  = flag.String("p", "botsunit8075", "taiga password")
	taigaProject   = flag.String("t", "Ufancyme", "taiga project")
	taigaURL       = flag.String("d", "https://taiga.botsunit.io", "taiga URL")
	taigaMilestone = flag.String("m", "0.5", "taiga Milestone to watch")
	taigaManager   *taigaclient.TaigaManager
)

func demosDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails(*taigaMilestone, ch)
	//ready := <-ch
	//if ready {
	taigaManager.MapStoriesPerUsers("Ready for test")
	taigaManager.MapIssuesPerUsers("Ready for test")
	//	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "demo.tmpl", gin.H{
		"title":       "Demo Short View",
		"userStories": taigaManager.StoriesPerUsers,
		"pointList":   taigaManager.PointList,
		"roleList":    taigaManager.RoleList,
		"issues":      taigaManager.IssuesPerUsers,
		"time":        elapsed,
	})
}

func overtakingUSDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails(*taigaMilestone, ch)

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "cr.tmpl", gin.H{
		"title":                 "Demo CR",
		"userStoriesOvertaking": taigaManager.StoriesDonePerUsers,
		"time":                  elapsed,
	})
}

func demosCRDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails(*taigaMilestone, ch)
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
		"time":                elapsed,
	})

}

func wipDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails(*taigaMilestone, ch)
	// ready := <-ch
	// if ready {
	taigaManager.MapStoriesPerUsers("In progress")
	taigaManager.MapIssuesPerUsers("In progress")
	// }

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "wip.tmpl", gin.H{
		"title":       "Work In Progress Short View",
		"userStories": taigaManager.StoriesPerUsers,
		"pointList":   taigaManager.PointList,
		"roleList":    taigaManager.RoleList,
		"issues":      taigaManager.IssuesPerUsers,
		"time":        elapsed,
	})
}

func main() {
	flag.Parse()
	taigaManager = (&taigaclient.TaigaManager{}).NewTaigaManager(taigaUsername, taigaPassword, taigaProject, taigaURL)

	// for id, name := range taigaManager.RoleList {
	// 	fmt.Println(fmt.Sprintf("Role ID : %d, Name : %s", id, name))
	// }
	//
	// for id, name := range taigaManager.PointList {
	// 	fmt.Println(fmt.Sprintf("Point ID : %d, Name : %s", id, name))
	// }
	router := gin.Default()
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
	router.Run(":8080")

}
