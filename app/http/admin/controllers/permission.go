package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gocms/app/http/admin/validates"
	"gocms/app/models/admin"
	"gocms/app/models/casbin"
	"gocms/app/models/permission"
	"gocms/app/service"
	"gocms/pkg/config"
	"gocms/pkg/response"
)

type PermissionController struct{}

var permissionService = &service.PermissionService{}

// 权限节点列表
func (that *PermissionController) Index(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	list := permissionService.GetList(cast.ToInt(page), cast.ToInt(pageSize))

	response.SuccessResponse(list).WriteTo(c)
	return
}

// 数据保存
func (that *PermissionController) Store(c *gin.Context) {
	var params map[string]string
	if validates.VidateCreateOrUpdatePermission(c, &params) == false {
		return
	}
	var permissionModel permission.Permission
	_ = mapstructure.Decode(params, &permissionModel)

	response.SuccessResponse().WriteTo(c)
	return
}

// 数据更新
func (that *PermissionController) Save(c *gin.Context) {
	var params map[string]string
	params["id"] = c.Param("id")
	if validates.VidateCreateOrUpdatePermission(c, &params) == false {
		return
	}
	var permissionModel permission.Permission
	_ = mapstructure.Decode(params, &permissionModel)

	response.SuccessResponse().WriteTo(c)
	return
}

// 数据删除
func (that *PermissionController) destory(c *gin.Context) {
	id := c.Param("id")
	config.Db.Delete(permission.Permission{}, "id = "+id)

	response.SuccessResponse().WriteTo(c)
	return
}

// 权限重置
func (that *PermissionController) reset(c *gin.Context) {
	// 是否权限清除并重置
	is_clear := cast.ToBool(c.PostForm("is_clear"))
	db := config.Db
	if is_clear == true {
		db.Delete(casbin.Casbin{}, "id > 0")
		db.Delete(permission.Permission{}, "id > 0")
	}

	// 自动创建权限
	admin.GeneratePermissionNodes()
	response.SuccessResponse().WriteTo(c)
	return
}