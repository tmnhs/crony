package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/middlerware"
	"github.com/tmnhs/crony/admin/internal/model/resp"
)

func RegisterRouters(r *gin.Engine) {
	r.Use(middlerware.Cors())

	configRoute(r)

	configNoRoute(r)
}

func configRoute(r *gin.Engine) {

	hello := r.Group("/ping")
	{
		hello.GET("", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
		hello.POST("", func(c *gin.Context) {
			type Hello struct {
				Name string `json:"name" form:"name"`
			}
			var h Hello
			var err error
			err = c.ShouldBindJSON(&h)
			if err != nil {
				c.JSON(resp.ERROR, err.Error())
			}
			c.JSON(200, "hello,"+h.Name)
		})
	}

	base := r.Group("")
	{
		base.POST("register", defaultUserRouter.Register)
		base.POST("login", defaultUserRouter.Login)
	}

	stat := r.Group("/statis")
	stat.Use(middlerware.JWTAuth())
	{
		stat.GET("today", defaultStatRouter.GetTodayStatistics)
		stat.GET("week", defaultStatRouter.GetWeekStatistics)
		stat.GET("system", defaultStatRouter.GetSystemInfo)

	}

	job := r.Group("/job")
	job.Use(middlerware.JWTAuth())
	{
		job.POST("add", defaultJobRouter.CreateOrUpdate)
		job.POST("del", defaultJobRouter.Delete)
		job.GET("find", defaultJobRouter.FindById)
		job.POST("search", defaultJobRouter.Search)
		job.POST("log", defaultJobRouter.SearchLog)
		job.POST("once", defaultJobRouter.Once)
		//job.POST("kill", defaultJobRouter.Kill)
	}

	user := r.Group("/user")
	user.Use(middlerware.JWTAuth())
	{
		user.POST("del", defaultUserRouter.Delete)
		user.POST("update", defaultUserRouter.Update)
		user.POST("change_pw", defaultUserRouter.ChangePassword)
		user.GET("find", defaultUserRouter.FindById)
		user.POST("search", defaultUserRouter.Search)
	}
	node := r.Group("/node")
	node.Use(middlerware.JWTAuth())
	{
		node.POST("search", defaultNodeRouter.Search)
		node.POST("del", defaultNodeRouter.Delete)
	}
	script := r.Group("/script")
	script.Use(middlerware.JWTAuth())
	{
		script.POST("add", defaultScriptRouter.CreateOrUpdate)
		script.POST("del", defaultScriptRouter.Delete)
		script.GET("find", defaultScriptRouter.FindById)
		script.POST("search", defaultScriptRouter.Search)
	}
}

func configNoRoute(r *gin.Engine) {
	r.LoadHTMLGlob("./dist/*.html") // npm?????????dist?????????
	r.StaticFile("favicon.ico", "./dist/favicon.ico")
	r.Static("/css", "./dist/css")         // dist?????????????????????
	r.Static("/fonts", "./dist/fonts")     // dist?????????????????????
	r.Static("/js", "./dist/js")           // dist?????????????????????
	r.Static("/img", "./dist/img")         // dist?????????????????????
	r.StaticFile("/", "./dist/index.html") // ????????????????????????
}
