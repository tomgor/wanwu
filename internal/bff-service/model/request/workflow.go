package request

type WorkflowIDReq struct {
	WorkflowID string `json:"workflow_id" validate:"required"`
}

func (r *WorkflowIDReq) Check() error {
	return nil
}

type GetWorkflowListReq struct {
	UserId string `form:"userId" json:"userId" validate:"required" `
	OrgId  string `form:"orgId" json:"orgId" validate:"required" `
}

func (g *GetWorkflowListReq) Check() error {
	return nil
}

type WorkflowImportReq struct {
	FileName string `json:"fileName" validate:"required"`
}

func (w *WorkflowImportReq) Check() error {
	return nil
}
