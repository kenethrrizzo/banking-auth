package domain

import "strings"

type RolePermissions struct {
	permissions map[string][]string
}

func (m RolePermissions) IsAuthorizedFor(role, routeName string) bool {
	perms := m.permissions[role]
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return true
		}
	}
	return false
}

func GetRolePermissions() RolePermissions {
	return RolePermissions{map[string][]string{
		"admin": {"GetAllCustomers", "GetCustomer", "NewAccount", "NewTransaction"},
		"user":  {"GetCustomer", "NewTransaction"},
	}}
}
