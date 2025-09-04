// @Author wangxm 8/13/星期三 15:20:00
package assistant

import (
	"context"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/pkg/util"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AssistantCustomToolCreate(ctx context.Context, req *assistant_service.AssistantCustomToolCreateReq) (*empty.Empty, error) {
	assistantId := util.MustU32(req.AssistantId)

	if status := s.cli.CreateAssistantCustom(ctx, assistantId, req.CustomToolId, req.Identity.UserId, req.Identity.OrgId); status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}

	return &empty.Empty{}, nil
}

func (s *Service) AssistantCustomToolDelete(ctx context.Context, req *assistant_service.AssistantCustomToolDeleteReq) (*empty.Empty, error) {
	assistantId := util.MustU32(req.AssistantId)

	if status := s.cli.DeleteAssistantCustom(ctx, assistantId, req.CustomToolId); status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}
	return &empty.Empty{}, nil
}

func (s *Service) AssistantCustomToolDeleteByCustomToolId(ctx context.Context, req *assistant_service.AssistantCustomToolDeleteByCustomToolIdReq) (*empty.Empty, error) {
	if status := s.cli.DeleteAssistantCustomByCustomToolId(ctx, req.CustomToolId); status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}
	return &empty.Empty{}, nil
}

func (s *Service) AssistantCustomToolEnableSwitch(ctx context.Context, req *assistant_service.AssistantCustomToolEnableSwitchReq) (*empty.Empty, error) {
	assistantId := util.MustU32(req.AssistantId)

	existingCustom, status := s.cli.GetAssistantCustom(ctx, assistantId, req.CustomToolId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}

	existingCustom.Enable = req.Enable
	if status := s.cli.UpdateAssistantCustom(ctx, existingCustom); status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}

	return &empty.Empty{}, nil
}
func (s *Service) AssistantCustomToolGetList(ctx context.Context, req *assistant_service.AssistantCustomToolGetListReq) (*assistant_service.AssistantCustomToolList, error) {
	assistantId := util.MustU32(req.AssistantId)

	customList, status := s.cli.GetAssistantCustomList(ctx, assistantId)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantCustomErr, status)
	}

	assistantCustomInfos := make([]*assistant_service.AssistantCustomToolInfo, len(customList))
	for i, custom := range customList {
		assistantCustomInfos[i] = &assistant_service.AssistantCustomToolInfo{
			Id:           custom.ID,
			AssistantId:  custom.AssistantId,
			CustomToolId: custom.CustomId,
			Enable:       custom.Enable,
			UserId:       custom.UserId,
			OrgId:        custom.OrgId,
			CreatedAt:    custom.CreatedAt,
			UpdatedAt:    custom.UpdatedAt,
		}
	}

	return &assistant_service.AssistantCustomToolList{
		AssistantCustomToolInfos: assistantCustomInfos,
	}, nil
}
