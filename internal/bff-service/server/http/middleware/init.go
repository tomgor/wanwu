package middleware

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	mid "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/mid-wrap"
	"github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/route"
	jwt_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/jwt-util"
)

func Init() {
	jwt_util.InitUserJWT(config.Cfg().JWT.SigningKey)

	mid.InitWrapper(Record)

	// --- openapi ---
	mid.NewSub("openapi", "对外提供原子能力", route.PermNone, false, false)

	// --- callback ---
	mid.NewSub("callback", "系统内部调用", route.PermNone, false, false)

	// --- guest ---
	mid.NewSub("guest", "", route.PermNone, false, false)

	// --- common ---
	mid.NewSub("common", "", route.PermNeedEnable, false, false, JWTUser, CheckUserEnable)

	// --- model ---
	mid.NewSub("model", "模型管理", route.PermNeedEnable, false, false, JWTUser, CheckUserEnable)

	// --- permission ---
	mid.NewSub("permission", "组织管理", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// permission.user
	mid.Sub("permission").NewSub("user", "用户", route.PermNeedCheck, true, true)

	// permission.org
	mid.Sub("permission").NewSub("org", "组织", route.PermNeedCheck, true, true)

	// permission.role
	mid.Sub("permission").NewSub("role", "角色", route.PermNeedCheck, true, true)

	// --- exploration ---
	mid.NewSub("exploration", "探索广场", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- workspace ---
	mid.NewSub("workspace", "工作室", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// workspace.appspace
	mid.Sub("workspace").NewSub("appspace", "应用空间", route.PermNeedCheck, true, true)

	// workspace.knowledge
	mid.Sub("workspace").NewSub("knowledgebase", "知识库", route.PermNeedCheck, true, true)
}
