package drm

type AccessPolicyAgent struct {
	policies map[string]map[string][]string
}

func NewAccessPolicyAgent() *AccessPolicyAgent {
	return &AccessPolicyAgent{
		policies: map[string]map[string][]string{
			"admin": {
				"user":    {"create", "read", "update", "delete"},
				"product": {"create", "read", "update", "delete"},
				"order":   {"create", "read", "update", "delete"},
			},
			"user": {
				"user":    {"read", "update"},
				"product": {"read"},
				"order":   {"create", "read"},
			},
			"guest": {
				"product": {"read"},
			},
		},
	}
}

func (a *AccessPolicyAgent) CheckAccess(command *Command) bool {
	rolePermissions, roleExists := a.policies[command.UserRole]
	if !roleExists {
		return false
	}

	entityPermissions, entityExists := rolePermissions[command.Entity]
	if !entityExists {
		return false
	}

	for _, permission := range entityPermissions {
		if permission == command.Action {
			return true
		}
	}

	return false
}