package role

import (
	"context"
)

func HasRoleFromContext(ctx context.Context, roleName string) bool {
	isRoleValid := false
	roles, ok := ctx.Value("auth_roles").([]string)
	if !ok {
		return false
	}

	if len(roles) > 0 {
		for _, role := range roles {
			if role == roleName {
				isRoleValid = true
				break
			}
		}
	}

	return isRoleValid
}
