package service

import "github.com/glueops/waggle/internal/repo"

type PolicyService struct {
	platform *repo.PlatformPolicyRepoGorm
	org      *repo.OrgPolicyRepoGorm
	users    *repo.UserRepoGorm
	passkeys *repo.UserPasskeyRepoGorm
}

func NewPolicyService(
	platform *repo.PlatformPolicyRepoGorm,
	org *repo.OrgPolicyRepoGorm,
	users *repo.UserRepoGorm,
	passkeys *repo.UserPasskeyRepoGorm,
) *PolicyService {
	return &PolicyService{
		platform: platform,
		org:      org,
		users:    users,
		passkeys: passkeys,
	}
}
