package router
import "github.com/gin-gonic/gin"

func Application(argv string){
	router := gin.Default()
	router.GET("/getfaces")
	router.Run()
}