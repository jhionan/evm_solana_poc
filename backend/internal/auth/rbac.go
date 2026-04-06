package auth

// Permission represents a single capability that can be granted to a role.
type Permission string

const (
	PermStake          Permission = "stake"
	PermUnstake        Permission = "unstake"
	PermClaimRewards   Permission = "claim_rewards"
	PermGetPosition    Permission = "get_position"
	PermListPositions  Permission = "list_positions"
	PermGetTiers       Permission = "get_tiers"
	PermPause          Permission = "pause"
	PermUpdateTier     Permission = "update_tier"
)

// rolePermissions maps each Role to the set of Permissions it holds.
var rolePermissions = map[Role][]Permission{
	RoleUser: {
		PermStake,
		PermUnstake,
		PermClaimRewards,
		PermGetPosition,
		PermListPositions,
		PermGetTiers,
	},
	RoleAdmin: {
		PermStake,
		PermUnstake,
		PermClaimRewards,
		PermGetPosition,
		PermListPositions,
		PermGetTiers,
		PermPause,
		PermUpdateTier,
	},
}

// HasPermission reports whether role holds the given permission.
func HasPermission(role Role, perm Permission) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
