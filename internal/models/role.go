package models

// Role represents a user role in the system
type Role string

// Role constants
const (
	RoleSuperUser Role = "superuser"
	RoleStore     Role = "store"
	RoleDelivery  Role = "delivery"
	RoleClient    Role = "client"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleSuperUser, RoleStore, RoleDelivery, RoleClient:
		return true
	}
	return false
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}
