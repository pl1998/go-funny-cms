package pkg1

import (
	"fmt"
	"gocms/app/service"
)

func init() {

}

func Echo() {
	//roleService := new(service.RoleService)
	//role := role2.RoleModel{
	//	Name:        "Surest",
	//	Description: "Surest",
	//}
	//roleService.UpdateOrCreateById(role)

	permissionService := new(service.PermissionService)
	permissionService.AddPermissionForUser("/api/auth/me", "GET", "Surest")
	permissionService.AddRoleForUser("Surest", "chenf")
	fmt.Println(permissionService.HasPermissionForUser("chenf", "GET", "/api/auth/me"))

}
