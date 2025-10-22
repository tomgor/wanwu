package service

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"

	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/gin-gonic/gin"
)

type AesKeyIv struct {
	Key []byte // AES 密钥
	IV  []byte // AES 初始化向量 (Initialization Vector)
}

type SsoUserInfo struct {
	RealName  string `json:"real_name"`
	MobileNum string `json:"mobile_num"`
	Username  string `json:"username"`
}

func SimpleSSO(ctx *gin.Context, login *request.SimpleSSO, language string) (*response.Login, error) {
	// 检测平台是否在配置中
	ssoConfig, err := GetSimpleSSOConfigByPlatform(login.Platform)
	if err != nil {
		return nil, err
	}

	// 通过私钥解密出aes key 和 iv
	aesKeyIv, err := RsaDecryptAesKeyIV(ssoConfig.PrivateKey, login.Key)
	if err != nil {
		return nil, fmt.Errorf("RSA解密AES Key和IV失败: %v", err)
	}

	//payload 通过.分隔成json字符串和sign签名字符串
	parts := strings.SplitN(login.Payload, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的Payload格式")
	}
	jsonStr, sig := parts[0], parts[1]

	// 使用解密得到的 AES Key 和 IV 解密 Payload
	userJson, err := AesDecrypt(jsonStr, aesKeyIv)
	if err != nil {
		return nil, fmt.Errorf("AES解密Payload失败: %v", err)
	}

	// 使用RSA公钥验证签名
	isValid, err := VerifySignature(userJson, sig, ssoConfig.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("验证签名失败: %v", err)
	}
	if !isValid {
		return nil, fmt.Errorf("签名验证失败")
	}

	var userInfo SsoUserInfo
	if err := json.Unmarshal([]byte(userJson), &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	userRsp, err := iam.GetUserIDByOrgAndName(ctx.Request.Context(), &iam_service.GetUserIDByOrgAndNameReq{
		OrgId: ssoConfig.FixedOrgID,
		Name:  userInfo.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("获取用户ID失败: %v", err)
	}

	var userId string
	if userRsp.UserId == 0 {
		// 用户不存在，准备创建新用户
		// 尝试创建用户
		resp, err := iam.CreateUser(ctx.Request.Context(), &iam_service.CreateUserReq{
			UserName:  userInfo.Username,
			Phone:     userInfo.MobileNum,
			NickName:  userInfo.RealName,
			OrgId:     ssoConfig.FixedOrgID,
			RoleIds:   []string{ssoConfig.FixedRoleId},
			Remark:    "SimpleSSO自动创建用户",
			Password:  ssoConfig.InitPassword,
			CreatorId: "1",
		})

		if err != nil {
			return nil, err
		}

		// 改个密码，否则会提示用户修改
		_, err = iam.UpdateUserPassword(ctx.Request.Context(), &iam_service.UpdateUserPasswordReq{
			UserId:      resp.Id,
			OldPassword: ssoConfig.InitPassword,
			NewPassword: ssoConfig.InitPassword,
		})

		if err != nil {
			return nil, fmt.Errorf("初始化用户密码失败: %v", err)
		}

		userId = resp.Id
	} else {
		userId = strconv.FormatUint(uint64(userRsp.UserId), 10)
	}

	user, err := iam.GetUserInfo(ctx.Request.Context(), &iam_service.GetUserInfoReq{
		UserId: userId,
		OrgId:  ssoConfig.FixedOrgID,
	})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	userPermission, err := iam.GetUserPermission(ctx.Request.Context(), &iam_service.GetUserPermissionReq{
		UserId: userId,
		OrgId:  ssoConfig.FixedOrgID,
	})

	if err != nil {
		return nil, fmt.Errorf("获取用户权限失败: %v", err)
	}

	// orgs
	orgs, err := iam.GetOrgSelect(ctx.Request.Context(), &iam_service.GetOrgSelectReq{UserId: user.UserId})
	if err != nil {
		return nil, err
	}
	// jwt token
	claims, token, err := jwt_util.GenerateToken(
		user.UserId,
		jwt_util.UserTokenTimeout,
	)
	if err != nil {
		return nil, err
	}
	ctx.Set(gin_util.CLAIMS, &claims)
	// resp
	return &response.Login{
		UID:              user.UserId,
		Username:         user.UserName,
		Nickname:         user.NickName,
		Token:            token,
		ExpiresAt:        claims.StandardClaims.ExpiresAt * 1000, // 超时事件戳毫秒
		ExpireIn:         strconv.FormatInt(jwt_util.UserTokenTimeout, 10),
		Orgs:             toOrgIDNames(ctx, orgs.Selects, user.UserId == config.SystemAdminUserID),
		OrgPermission:    toOrgPermission(ctx, userPermission),
		Language:         getLanguageByCode(user.Language),
		IsUpdatePassword: userPermission.LastUpdatePasswordAt != 0, //是否已修改密码，单点登陆不用修改密码
	}, nil
}

func GetSimpleSSOConfigByPlatform(platform string) (*config.SimpleSSOConfig, error) {
	// 确保平台标识符不为空
	if platform == "" {
		return nil, errors.New("platform 标识符不能为空")
	}

	// 遍历全局配置中的 SimpleSSO 列表
	for _, ssoConfig := range config.Cfg().SimpleSSO {
		// 忽略大小写或精确匹配取决于实际需求，这里使用精确匹配
		if ssoConfig.Platform == platform {
			// 找到匹配的配置，返回其地址
			// 注意：这里返回的是切片元素的副本的地址，
			// 如果需要修改配置并同步到全局变量，需要对 Cfg().SimpleSSO 列表进行基于索引的引用。
			// 但对于配置读取，返回副本的地址通常是安全的做法。
			return &ssoConfig, nil
		}
	}

	// 遍历结束后未找到，返回错误
	return nil, fmt.Errorf("未找到平台 '%s' 对应的 SimpleSSO 配置", platform)
}

// RsaDecryptAesKeyIV 使用 Base64 编码的 RSA 私钥解密 Base64 编码的密文，
// 密文应包含 AES Key 和 IV 的组合数据。
//
// 注意: 密文的结构 (Key + IV) 需要与加密方约定一致。
// 这里假设密文内容是：AES_KEY (N bytes) + IV (16 bytes)
func RsaDecryptAesKeyIV(base64PrivateKey string, base64Ciphertext string) (*AesKeyIv, error) {
	// 1. Base64 解码私钥
	// keyBytes, err := base64.StdEncoding.DecodeString(base64PrivateKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("base64解码私钥失败: %w", err)
	// }

	// // 2. 解析 PEM 格式的私钥
	// block, _ := pem.Decode(keyBytes)
	// if block == nil {
	// 	// 尝试去掉可能的 PEM 头尾，直接解析 PKCS#8 或 PKCS#1 格式
	// 	// Go 的 x509.ParsePKCS8PrivateKey 和 x509.ParsePKCS1PrivateKey
	// 	// 通常要求输入的是 DER 格式 (即 PEM 解码后的 Bytes)。
	// 	// 如果 base64PrivateKey 已经包含了 PEM 块，block 就不会是 nil。
	// 	// 如果用户传入的是裸露的 Base64 编码的 DER 格式，则需要进一步处理。
	// 	// 实际上，更安全和标准的方式是要求用户传入带有 BEGIN/END 标记的 PEM 字符串。

	// 	// 尝试直接解析为 DER 格式
	// 	privKey, parseErr := parsePrivateKey(keyBytes)
	// 	if parseErr != nil {
	// 		return nil, fmt.Errorf("私钥格式解析失败，确保它是合法的 PEM/DER 格式: %w", parseErr)
	// 	}
	// 	privKey = privKey
	// } else {
	// 	// 检查 PEM 类型并解析私钥
	// 	privKey, parseErr := parsePrivateKey(block.Bytes)
	// 	if parseErr != nil {
	// 		return nil, fmt.Errorf("解析PEM私钥块失败: %w", parseErr)
	// 	}
	// }

	// ⚠️ 假设 parsePrivateKey 成功，这里需要重构以正确获取私钥
	var privateKey *rsa.PrivateKey

	// 为了简洁和演示核心逻辑，我们假设私钥解析成功并赋值给 privateKey
	// 实际调用时，需要将上面的私钥解析逻辑正确整合进来。
	// 这里使用一个占位符函数来简化上面的复杂分支
	privateKey, err := parseRsaPrivateKey(base64PrivateKey)
	if err != nil {
		return nil, err
	}

	// 3. Base64 解码密文
	ciphertext, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("base64解码密文失败: %w", err)
	}

	// 4. RSA 解密
	// 推荐使用 RSA-OAEP 模式进行密钥解密，因为 PKCS#1 v1.5 已知存在安全问题。
	// 这里使用 SHA256 作为散列函数，并使用空标签（label）。
	// decryptedBytes, err := rsa.DecryptOAEP(
	// 	sha256.New(),
	// 	rand.Reader,
	// 	privateKey,
	// 	ciphertext,
	// 	nil, // label
	// )
	decryptedBytes, err := rsa.DecryptPKCS1v15(nil, privateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("RSA解密失败 (密文或填充错误): %w", err)
	}

	// 5. 解析 AES Key 和 IV
	// 假设 AES Key 长度是 32 字节 (AES-256)，IV 长度是 16 字节 (AES 块大小)
	const aesKeyLen = 32
	const ivLen = 16

	if len(decryptedBytes) < aesKeyLen+ivLen {
		return nil, fmt.Errorf("解密后的数据长度不足，期望至少 %d 字节，实际 %d 字节", aesKeyLen+ivLen, len(decryptedBytes))
	}

	// 约定 Key 在前，IV 在后
	aesKey := decryptedBytes[:aesKeyLen]
	aesIV := decryptedBytes[aesKeyLen : aesKeyLen+ivLen]

	return &AesKeyIv{
		Key: aesKey,
		IV:  aesIV,
	}, nil
}

// 辅助函数：解析 Base64 字符串中的 RSA 私钥
func parseRsaPrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}

	// 尝试解析 PEM 格式
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		keyBytes = block.Bytes
	} else {
		// 如果没有 PEM 块，尝试清理 Base64 字符串中的空格和换行符
		keyBytes = []byte(strings.TrimSpace(string(keyBytes)))
		// 再次尝试 Base64 解码，以防原始输入是裸露的 Base64 且带有空格
		// 简化处理，只针对标准 PEM/DER 格式

		// 如果 block 为 nil，但用户传入的是裸露的 Base64(DER)，则 keyBytes 已经是 DER
	}

	// 尝试 PKCS#8 (最常见)
	if priv, err := x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		if rsaPriv, ok := priv.(*rsa.PrivateKey); ok {
			return rsaPriv, nil
		}
		return nil, errors.New("解析到的私钥不是 RSA 类型")
	}

	// 尝试 PKCS#1
	if priv, err := x509.ParsePKCS1PrivateKey(keyBytes); err == nil {
		return priv, nil
	}

	return nil, errors.New("无法解析 RSA 私钥，请检查 Base64 编码和 PEM/DER 格式")
}

const (
	AESKeySize = 32 // 默认使用 AES-256，即 32 字节密钥
	AESIVSize  = 16 // IV 长度固定为 16 字节
)

// pkcs7UnPadding 移除 PKCS#7 填充。
// 在 Go 的 crypto/cipher 库中，CBC 模式需要手动处理填充。
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("密文数据为空，无法进行去填充操作")
	}

	// 填充长度由最后一个字节的值决定
	paddingLen := int(data[length-1])

	// 检查填充长度是否合法（必须大于 0 且小于或等于数据长度）
	if paddingLen <= 0 || paddingLen > length {
		return nil, errors.New("填充长度非法，可能密文被篡改或密钥/IV错误")
	}

	// 检查填充字节是否一致
	for i := 0; i < paddingLen; i++ {
		if data[length-1-i] != byte(paddingLen) {
			return nil, errors.New("填充字节不一致，去填充失败")
		}
	}

	// 返回去除填充后的数据
	return data[:length-paddingLen], nil
}

// AesDecrypt 使用 AES-256 CBC 模式解密 Base64 编码的数据。
//
// 参数:
//
//	encryptedData: Base64编码的加密数据
//	aesKey: AES密钥（必须是 32 字节）
//	iv: 初始化向量（必须是 16 字节）
//
// 返回值:
//
//	解密后的字符串，或错误信息
func AesDecrypt(encryptedData string, aesKeyIv *AesKeyIv) (string, error) {
	// 1. 输入参数校验 (与 Java 版本保持一致)
	if encryptedData == "" {
		return "", errors.New("加密数据不能为空")
	}
	if len(aesKeyIv.Key) != AESKeySize {
		return "", fmt.Errorf("AES密钥长度必须为 %d 字节，实际为 %d 字节", AESKeySize, len(aesKeyIv.Key))
	}
	if len(aesKeyIv.IV) != AESIVSize {
		return "", fmt.Errorf("IV长度必须为 %d 字节，实际为 %d 字节", AESIVSize, len(aesKeyIv.IV))
	}

	// 2. Base64 解码加密数据
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	// 检查密文长度是否是 AES 块大小的倍数 (AES 块大小固定为 16)
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", errors.New("加密数据长度不是 AES 块大小的倍数，密文可能不完整或格式错误")
	}

	// 3. 创建 AES 密码块
	block, err := aes.NewCipher(aesKeyIv.Key)
	if err != nil {
		return "", fmt.Errorf("创建AES密码块失败: %w", err)
	}

	// 4. 创建 CBC 解密器
	blockMode := cipher.NewCBCDecrypter(block, aesKeyIv.IV)

	// 5. 执行解密操作 (数据长度不变)
	decryptedBytes := make([]byte, len(encryptedBytes))
	blockMode.CryptBlocks(decryptedBytes, encryptedBytes)

	// 6. 移除 PKCS#7 填充
	unpaddedBytes, err := pkcs7UnPadding(decryptedBytes)
	if err != nil {
		return "", fmt.Errorf("去填充失败: %w", err)
	}

	return string(unpaddedBytes), nil
}

// parseRsaPublicKey 解析 Base64 编码的 PEM/DER 格式的 RSA 公钥字符串。
func parseRsaPublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	// Base64解码公钥字符串
	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		// 尝试直接使用原始字符串（可能公钥字符串本身就是 PEM 格式）
		keyBytes = []byte(publicKeyStr)
	}

	// 解析 PEM 块 (如果存在)
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		keyBytes = block.Bytes
	}

	// 1. 尝试 PKIX (X.509) 格式 (最常见，用于公钥证书)
	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
		return nil, errors.New("解析到的公钥不是 RSA 类型")
	}

	// 2. 尝试 PKCS#1 格式
	rsaPub, err := x509.ParsePKCS1PublicKey(keyBytes)
	if err == nil {
		return rsaPub, nil
	}

	return nil, fmt.Errorf("无法解析 RSA 公钥: %v", err)
}

// VerifySignature 验证 RSA 签名（使用 SHA256 散列和 PSS 填充）。
//
// 参数:
//
//	data: 待验证的原始数据字符串。
//	signatureBase64: Base64 编码的签名数据。
//	publicKeyStr: Base64 编码的 PEM/DER 格式的 RSA 公钥字符串。
//
// 返回值:
//
//	bool: 验证结果 (true 为通过，false 为失败)。
//	error: 验证过程中发生的错误（如格式错误、解析失败）。
func VerifySignature(data string, signatureBase64 string, publicKeyStr string) (bool, error) {
	// 1. 输入参数校验 (与 Java 版本逻辑一致)
	if data == "" {
		return false, errors.New("待验证数据不能为空")
	}
	if signatureBase64 == "" {
		return false, errors.New("签名不能为空")
	}
	if publicKeyStr == "" {
		return false, errors.New("公钥不能为空")
	}

	// 2. 解析公钥
	publicKey, err := parseRsaPublicKey(publicKeyStr)
	if err != nil {
		return false, fmt.Errorf("解析公钥失败: %w", err)
	}

	// 3. Base64 解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, fmt.Errorf("Base64解码签名失败: %w", err)
	}

	// // 4. 计算数据的哈希值 (SHA256)
	// // Go 的 RSA 签名验证 API 要求传入数据的哈希值
	// hashed := sha256.Sum256([]byte(data))

	// // 5. 验证签名 (使用 PSS 模式)
	// // 如果加密方使用 PKCS#1 v1.5 (Java 默认的 "SHA256withRSA")，则应使用 rsa.VerifyPKCS1v15
	// // 这里推荐使用更安全的 PSS 模式
	// err = rsa.VerifyPKCS1v15(
	// 	publicKey,
	// 	crypto.SHA256,
	// 	hashed[:], // 传入哈希切片
	// 	signatureBytes,
	// )

	hash := sha256.New()
	hash.Write([]byte(data))
	hashedData := hash.Sum(nil)

	// 4. 验证签名 (使用PKCS#1 v1.5填充)
	// rsa.VerifyPKCS1v15 在验证成功时返回 nil，失败时返回一个 error
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedData, signatureBytes)

	if err != nil {
		// rsa.VerifyPSS/VerifyPKCS1v15 验证失败时返回非 nil 错误
		// 签名失败是正常的业务流程结果，我们只返回 false，将错误记录在日志中（或在调用方处理）
		// fmt.Printf("签名验证失败: %v\n", err) // 可以在此处进行调试日志记录
		return false, fmt.Errorf("签名验证失败: %w", err)
	}

	// 签名验证成功
	return true, nil
}
