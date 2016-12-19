package main

import (
	"flag"
	"fmt"
	"net/http"
	"taiga_tracker/taigaclient"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	taigaUsername = flag.String("u", "api", "taiga username")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
	taigaProject  = flag.String("t", "Ufancyme", "taiga project")
	taigaManager  *taigaclient.TaigaManager
)

func refreshDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails("0.5", ch)
	// ready := <-ch
	// if ready {
	taigaManager.MapStoriesPerUsers()
	taigaManager.MapIssuesPerUsers()
	// }

	// usList := taigaManager.StoriesPerUsers["Mathieu Artu"]
	// usList[0].Points

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
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
	taigaManager = (&taigaclient.TaigaManager{}).NewTaigaManager(taigaUsername, taigaPassword, taigaProject)

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

	authorized.GET("/wip", refreshDatas)
	router.Run(":8080")

}
