package header2

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
)

// Profile 用户信息结构体
type Profile struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	DepartmentID string `json:"department_id"`
}

// GetProfile 从request头部获取用户信息
func GetProfile(c *gin.Context) Profile {
	return getProfileFromGIN(c)
}

func getProfileFromGIN(c *gin.Context) Profile {
	userID := c.GetHeader(_userID)
	userName := c.GetHeader(_userName)
	departmentID := c.GetHeader(_departmentID)

	return Profile{
		UserID:       userID,
		UserName:     userName,
		DepartmentID: strings.Split(departmentID, ",")[0],
	}
}

// GetDepartments GetDepartments
func GetDepartments(c *gin.Context) []string {
	departmentID := c.GetHeader(_departmentID)
	return strings.Split(departmentID, ",")
}

// CloneProfile CloneProfile
func CloneProfile(dst *http.Header, src http.Header) {
	dst.Set(_userID, deepCopy(src.Values(_userID)))
	dst.Set(_userName, deepCopy(src.Values(_userName)))
	dst.Add(_departmentID, deepCopy(src.Values(_departmentID)))
}

func deepCopy(src []string) string {
	for _, elem := range src {
		if elem != "" {
			return elem
		}
	}
	return ""
}

const (
	roleName = "Role"
)

// SetRole SetRole
func SetRole(c *gin.Context, role ...string) {
	roles := strings.Join(role, ",")
	c.Request.Header.Set(roleName, roles)
	c.Writer.Header().Set(roleName, roles)
}

// Role role
type Role struct {
	Role []string
}

// IsSuper IsSuper
func (r *Role) IsSuper() bool {
	for _, role := range r.Role {
		if role == "super" {
			return true
		}
	}
	return false
}

// GetRole GetRole
func GetRole(c *gin.Context) *Role {
	roleStr := c.Request.Header.Get(roleName)
	roles := strings.Split(roleStr, ",")
	return &Role{Role: roles}
}
