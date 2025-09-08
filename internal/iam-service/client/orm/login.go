package orm

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/orm/sqlopt"
	iam_util "github.com/UnicomAI/wanwu/internal/iam-service/pkg/util"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

// 字符集定义
const (
	digits          = "0123456789"
	lowercase       = "abcdefghijklmnopqrstuvwxyz"
	uppercase       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special         = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
	EmailCodeExpire = 5 * time.Minute // 验证码过期时间
	CodeLength      = 6               // 验证码长度
	PasswordLength  = 12              // 密码长度
)

func (c *Client) Login(ctx context.Context, name, password, language string) (*UserInfo, *Permission, *errs.Status) {
	var userInfo *UserInfo
	var permission *Permission

	return userInfo, permission, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// user
		user := &model.User{}
		if err := sqlopt.WithName(name).Apply(tx).First(user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return toErrStatus("iam_login_invalid_name_pwd")
			}
			return toErrStatus("iam_login_err", err.Error())
		}
		if user.Password != util.SHA256(password) {
			return toErrStatus("iam_login_invalid_name_pwd")
		}
		// check status
		if !user.Status {
			return toErrStatus("iam_login_user_disable")
		}
		// org tree
		orgTree, err := getOrgTree(tx)
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		// user info
		userInfo, err = toUserInfoTx(tx, user, orgTree, true)
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		// select org
		orgs, err := selectOrgs(tx, user.ID, orgTree)
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		var orgID uint32
		if len(orgs) == 1 {
			orgID = orgs[0].ID
		} else if len(orgs) > 1 {
			var hasRole bool
			for _, org := range userInfo.Orgs {
				if org.Org.ID == orgs[0].ID {
					if len(org.Roles) > 0 {
						hasRole = true
						orgID = orgs[0].ID
					}
					break
				}
			}
			if !hasRole {
				orgID = orgs[1].ID
			}
		}
		if orgID != 0 {
			permission, err = getUserPermission(tx, user.ID, orgID)
		}
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		// last_login_at & last_exec_at
		nowTS := time.Now().UnixMilli()
		update := map[string]interface{}{
			"last_login_at": nowTS,
			"last_exec_at":  nowTS,
		}
		if language != "" {
			update["language"] = language
			userInfo.Language = language
		}
		if err := tx.Model(user).Updates(update).Error; err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		return nil
	})
}

func (c *Client) SendEmailCode(ctx context.Context, email string) *errs.Status {
	// 1. 生成6位数字验证码
	code := generateNumericCode(CodeLength)

	// 2. 生成包含特殊字符、字母和数字的密码
	password := generateComplexPassword(PasswordLength)
	item := []redis.HashItem{{K: "code", V: code}, {K: "password", V: password}}

	// 3. 存入Redis
	err := redis.IAM().HSetWithExpire(ctx, email, item, EmailCodeExpire)
	if err != nil {
		return toErrStatus("todo", err.Error())
	}

	// 4. 发送邮件（这里需要调用您的邮件发送服务）
	emailContent := fmt.Sprintf(`
尊敬的客户，

您的验证码是：%s
您的临时密码是：%s

验证码有效期为10分钟。

请勿将此信息透露给他人。

感谢您的使用！
	`, code, password)
	// 调用邮件发送服务
	err = iam_util.SendMail(email, "验证码和临时密码", emailContent)
	if err != nil {
		return toErrStatus("todo", err.Error())
	}
	return nil
}

// generateNumericCode 生成指定位数的数字验证码
func generateNumericCode(length int) string {
	var code strings.Builder
	code.Grow(length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// 如果随机数生成失败，使用fallback方法
			num = big.NewInt(int64(i % 10))
		}
		code.WriteString(digits[num.Int64() : num.Int64()+1])
	}
	return code.String()
}

// generateComplexPassword 生成包含特殊字符、字母和数字的复杂密码
func generateComplexPassword(length int) string {
	if length < 8 {
		length = 12 // 确保密码长度足够
	}
	// 组合所有字符
	allChars := digits + lowercase + uppercase + special
	var password strings.Builder
	password.Grow(length)
	// 确保至少包含每种类型的一个字符
	requiredChars := []string{
		digits,
		lowercase,
		uppercase,
		special,
	}
	// 先添加每个类型的至少一个字符
	for _, chars := range requiredChars {
		char, err := randomChar(chars)
		if err == nil {
			password.WriteString(char)
		}
	}
	// 填充剩余长度
	remaining := length - len(requiredChars)
	for i := 0; i < remaining; i++ {
		char, err := randomChar(allChars)
		if err == nil {
			password.WriteString(char)
		}
	}
	// 打乱密码顺序
	result := []rune(password.String())
	shuffleString(result)
	return string(result)
}

// randomChar 从字符串中随机选择一个字符
func randomChar(chars string) (string, error) {
	if len(chars) == 0 {
		return "", fmt.Errorf("empty character set")
	}
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return string(chars[0]), nil
	}
	return string(chars[idx.Int64()]), nil
}

// shuffleString 打乱字符串顺序
func shuffleString(s []rune) {
	for i := len(s) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		s[i], s[j.Int64()] = s[j.Int64()], s[i]
	}
}
