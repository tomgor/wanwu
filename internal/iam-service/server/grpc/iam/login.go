package iam

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/config"
	"github.com/UnicomAI/wanwu/internal/iam-service/pkg/util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	EmailCode = "code"
	Password  = "password"
)

func (s *Service) GetCaptcha(ctx context.Context, req *iam_service.GetCaptchaReq) (*iam_service.GetCaptchaResp, error) {
	code, b64s, err := util.GenerateCaptcha()
	if err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, toErrStatus("iam_login_captcha", err.Error()))
	}
	if err := s.cli.RefreshCaptcha(ctx, req.Key, code); err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, err)
	}
	return &iam_service.GetCaptchaResp{
		Key: req.Key,
		B64: b64s,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *iam_service.LoginReq) (*iam_service.LoginResp, error) {
	// captcha
	if err := s.cli.CheckCaptcha(ctx, req.Key, req.Code); err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, err)
	}
	// login
	user, permission, err := s.cli.Login(ctx, req.UserName, req.Password, req.Language)
	if err != nil {
		return nil, errStatus(errs.Code_IAMLogin, err)
	}
	return &iam_service.LoginResp{
		User:       toUserInfo(user),
		Permission: toPermission(permission),
	}, nil
}

func (s *Service) RegisterByEmail(ctx context.Context, req *iam_service.RegisterByEmailReq) (*emptypb.Empty, error) {
	code, redisErr := redis.IAM().HGet(ctx, req.Email, EmailCode)
	if redisErr != nil {
		return &emptypb.Empty{}, grpc_util.ErrorStatus(errs.Code_IAMRegister)
	}
	if code.V != req.Code {
		return &emptypb.Empty{}, grpc_util.ErrorStatus(errs.Code_IAMRegister)
	}
	passwd, redisErr := redis.IAM().HGet(ctx, req.Email, Password)
	if redisErr != nil {
		return &emptypb.Empty{}, grpc_util.ErrorStatus(errs.Code_IAMRegister)
	}
	userID, err := s.cli.CreateUser(ctx, &model.User{
		Status:    true,
		CreatorID: config.AdminUserID(),
		Name:      req.UserName,
		Email:     req.Email,
		Password:  passwd.V,
	}, config.TopOrgID(), []uint32{})
	if err != nil {
		return &emptypb.Empty{}, errStatus(errs.Code_IAMUser, err)
	}
	_, err = s.cli.CreateOrg(ctx, &model.Org{
		Status:    true,
		CreatorID: userID,
		ParentID:  config.TopOrgID(),
		Name:      req.UserName + "-org",
	})
	if err != nil {
		return &emptypb.Empty{}, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) SendEmailCode(ctx context.Context, req *iam_service.SendEmailCodeReq) (*emptypb.Empty, error) {
	if err := s.cli.SendEmailCode(ctx, req.Email); err != nil {
		return nil, errStatus(errs.Code_IAMRegister, err)
	}
	return &emptypb.Empty{}, nil
}
