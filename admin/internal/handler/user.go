package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/middlerware"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"time"
)

type UserRouter struct {
}

var defaultUserRouter = new(UserRouter)

func (u *UserRouter) Login(c *gin.Context) {
	var req request.ReqUserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[user_login] request parameter error", c)
		return
	}
	user, err := service.DefaultUserService.Login(req.UserName, req.Password)
	if err != nil || user == nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_login] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[user_login] username or password is incorrect", c)
		return
	}
	j := middlerware.NewJWT() // 唯一签名
	claims := j.CreateClaims(middlerware.BaseClaims{
		ID:       user.ID,
		UserName: user.UserName,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		logger.GetLogger().Debug(fmt.Sprintf("获取token失败 :%s", err.Error()))
		resp.FailWithMessage(resp.ErrorTokenGenerate, "获取token失败", c)
		return
	}
	resp.OkWithDetailed(resp.RspLogin{
		User:  user,
		Token: token,
	}, "login success", c)
}

func (u *UserRouter) Register(c *gin.Context) {
	var req request.ReqUserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_register] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[user_register] request parameter error", c)
		return
	}
	user, err := service.DefaultUserService.FindByUserName(req.UserName)
	if err != nil || user == nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_register] db find by username:%s error:%v", req.UserName, err))
		resp.FailWithMessage(resp.ErrorUserNameExist, "[user_register] the user name has already been used", c)
		return
	}
	if req.Role == 0 {
		req.Role = models.RoleNormal
	}
	userModel := &models.User{
		UserName: req.UserName,
		Password: utils.MD5(req.Password),
		Role:     req.Role,
		Email:    req.Email,
		Created:  time.Now().Unix(),
	}
	insertId, err := userModel.Insert()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_register] db insert error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[user_register] db insert error", c)
		return
	}
	userModel.ID = insertId
	resp.OkWithDetailed(userModel, "register success", c)
}

func (u *UserRouter) Update(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[update_user] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[update_user] request parameter error", c)
		return
	}
	req.Updated = time.Now().Unix()
	err := req.Update()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[update_user] db update error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[update_user] db update error", c)
		return
	}
	resp.OkWithMessage("update success", c)
}

func (u *UserRouter) Delete(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_user] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_user] request parameter error", c)
		return
	}
	userModel := models.User{ID: req.ID}
	err := userModel.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_user] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[delete_user] db error", c)
		return
	}
	resp.OkWithMessage("delete success", c)
}

func (u *UserRouter) ChangePassword(c *gin.Context) {
	var req request.ReqChangePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[change_password] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[change_password] request parameter error", c)
		return
	}
	err := service.DefaultUserService.ChangePassword(middlerware.GetUserInfo(c).ID, req.Password, req.NewPassword)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[change_password] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[change_password] db error", c)
		return
	}
	resp.OkWithMessage("update success", c)
}

func (u *UserRouter) FindById(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_user] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[find_user] request parameter error", c)
		return
	}
	userModel := models.User{ID: req.ID}
	err := userModel.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_user] db update error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[find_user] db update error", c)
		return
	}
	resp.OkWithDetailed(userModel, "find success", c)
}

func (u *UserRouter) Search(c *gin.Context) {
	var req request.ReqUserSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_user] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_user] request parameter error", c)
		return
	}
	req.Check()
	users, total, err := service.DefaultUserService.Search(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_user] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[search_user] db error", c)
		return
	}
	resp.OkWithDetailed(resp.PageResult{
		List:     users,
		Total:    total,
		PageSize: req.PageSize,
		Page:     req.Page,
	}, "search success", c)
}

func (u *UserRouter) JoinGroup(c *gin.Context) {
	var req models.UserGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_join_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[user_join_group] request parameter error", c)
		return
	}
	_, err := req.Insert()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_join_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[user_join_group] db error", c)
		return
	}
	resp.OkWithMessage("join success", c)
}

func (u *UserRouter) KickGroup(c *gin.Context) {
	var req models.UserGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_kick_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[user_kick_group] request parameter error", c)
		return
	}
	err := req.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_kick_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[user_kick_group] db error", c)
		return
	}
	resp.OkWithMessage("kick success", c)
}

func (u *UserRouter) GetByGroupId(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_get_by_group] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[user_get_by_group] request parameter error", c)
		return
	}
	users, err := service.DefaultUserService.FindByGroupId(req.ID)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[user_get_by_group] db error:%v", err))
		resp.FailWithMessage(resp.ERROR, "[user_get_by_group] db error", c)
		return
	}
	resp.OkWithDetailed(users, "get success", c)
}
