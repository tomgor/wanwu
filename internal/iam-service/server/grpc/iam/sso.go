package iam

import (
	"fmt"

	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/config"
	"github.com/UnicomAI/wanwu/internal/iam-service/server/grpc/iam/sso"
)

func FetchUserInfoBySso(platform, token string) (*model.User, *config.SSOConfig, error) {
	for _, ssoConfig := range config.Cfg().SSO {
		if ssoConfig.Name == platform && ssoConfig.Enabled {
			if platform == "xietong" {
				user, err := sso.FetchXieTongUserInfo(token, &ssoConfig)
				if err != nil {
					return nil, nil, err
				}
				return user, &ssoConfig, nil
			}

		}
	}
	return nil, nil, fmt.Errorf("unsupported or disabled platform: %s", platform)
}
