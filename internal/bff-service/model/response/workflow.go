package response

type CozeWorkflowListResp struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data *CozeWorkflowListData `json:"data,omitempty"`
}

type CozeWorkflowListData struct {
	Workflows []*CozeWorkflowListDataWorkflow `json:"workflow_list"`
}

type CozeWorkflowListDataWorkflow struct {
	WorkflowId string `json:"workflow_id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type CozeWorkflowCreateResp struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *CozeWorkflowCreateData `json:"data,omitempty"`
}

type CozeWorkflowCreateData struct {
	WorkflowID string `json:"workflow_id"`
}

type CozeWorkflowDeleteResp struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *CozeWorkflowDeleteData `json:"data,omitempty"`
}

type CozeWorkflowDeleteData struct {
	Status int64 `json:"status"`
}

func (d *CozeWorkflowDeleteData) GetStatus() int64 {
	if d == nil {
		return 0
	}
	return d.Status
}
