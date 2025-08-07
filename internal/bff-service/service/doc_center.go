package service

import (
	"context"
	"fmt"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const (
	docCenterLocalDir     = "configs/microservice/bff-service/static/docs" // bff-service本地缓存目录
	docCenterStaticPrefix = "../../../service/api/v1"                      // ../../..用于抵消前端固定前缀 aibase/docCenter/pages
	docCenterSnippetLen   = 200                                            // 截取文本长度
)

var (
	mdImageRegex          = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)        // 从markdown文本中匹配 ![](xxxxx) 图片引用
	mdParenthesisRefRegex = regexp.MustCompile(`\((.*?)\)`)                // 从markdown引用中匹配 (xxxxx)
	mdLinkRegex           = regexp.MustCompile(`[^!]\[.*?\]\((.*?\.md)\)`) // 从markdown匹配出跳转链接[](xxxxx)
	mdBracketRegex        = regexp.MustCompile(`\[(.*?)\]`)                // 从markdown匹配[]中的文本

	docSearchers = make(map[string]*riot.Engine) // 初始化搜索引擎
	docSearchMu  sync.RWMutex                    // 搜索引擎读写锁
	docMu        sync.Mutex                      // doc_center全局互斥锁

	docVerRegex = regexp.MustCompile(`v[\d\.]+`)

	docCenter *DocCenter
)

type DocCenter struct {
	DocId    uint32     `json:"docId,omitempty"`    // 文档ID，表主键
	Version  string     `json:"version,omitempty"`  // 版本号
	Desc     string     `json:"desc,omitempty"`     // 描述
	Children []*DocMenu `json:"children,omitempty"` // 菜单 对children 直接mashal json 存到menu里 用的时候unmashal
}

type DocMenu struct {
	Name         string     `json:"name,omitempty"`         // 目录/文档名，如“选模型”、“改模型”、“数据管理”等
	RelativePath string     `json:"relativePath,omitempty"` // 相对路径，md zip包中路径
	FilePath     string     `json:"filePath,omitempty"`     // 文件路径
	Children     []*DocMenu `json:"children,omitempty"`     // 子菜单
}

type MdInfo struct {
	relPath  string
	filePath string
}

func InitDocCenter(ctx context.Context) error {
	var err error
	var mdInfos []MdInfo
	version := config.Cfg().CustomInfo.About.Version
	docDir := path.Join(docCenterLocalDir, version)
	if err := filepath.Walk(docDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查不为目录并且是否为markdown文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// 读取文件内容
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("读取文件%v错误", filePath)
			}
			fileName := info.Name()
			// 将markdown文本中图片引用 ![](xxxxx )与链接引用[](xxxxx )里的 xxxxx 处理为前端可访问的地址
			convertByte := convertMarkdown(path.Join(docDir, fileName), fileName, string(content))

			// 存储文件
			if err = os.WriteFile(path.Join(docDir, fileName), []byte(convertByte), os.ModePerm); err != nil {
				return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_download_save", filePath, err.Error())
			}
			mdInfos = append(mdInfos, MdInfo{
				relPath:  fileName,
				filePath: path.Join(docDir, fileName),
			})

		}
		return nil
	}); err != nil {
		return err
	}
	// 初始化搜索引擎
	se, err := newDocSearcher(docCenterLocalDir, version)
	if err != nil {
		return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_search_init_err", err.Error())
	}
	// 以版本为key添加搜索引擎
	addDocSearcher(version, se)

	// 存储当前版本文档目录结构
	var menus []*DocMenu
	for _, mdInfo := range mdInfos {
		addMdFileToDocMenu(&menus, mdInfo.relPath, mdInfo)
	}
	// 文档name排序
	sortDocMenus(&menus)
	docCenter = &DocCenter{
		Children: menus,
	}
	return err
}

func GetDocCenterMenu(ctx *gin.Context) ([]response.DocMenu, error) {
	var ret []response.DocMenu
	for i, menu := range docCenter.Children {
		ret = append(ret, convertDocMenuToResponse(menu, fmt.Sprintf("doc%d", i+1)))
	}
	return ret, nil
}

func SearchDocCenter(ctx *gin.Context, content string) ([]response.DocSearchResp, error) {
	if isDocSearchersEmpty() {
		return nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_search_init_no")
	}
	// 从文档中心docDir
	var docDir string
	if path, ok := findDocCenterFirstFilePath(docCenter.Children); ok {
		docDir = docVerRegex.FindString(path)
	} else {
		return nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_file_version", strconv.Itoa(int(docCenter.DocId)), docCenter.Version)
	}
	// 找到对应searcher
	se, exists := getDocSearcher(docDir)
	if !exists {
		return nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_file_search", docDir)
	}
	results := se.SearchDoc(types.SearchReq{Text: content})
	var searchResps []response.DocSearchResp
	for _, doc := range results.Docs {
		title := strings.TrimSuffix(filepath.Base(doc.DocId), filepath.Ext(filepath.Base(doc.DocId)))
		snippet, err := util.Md2html([]byte(getMarkdownSnippet(doc.Content, content, docCenterSnippetLen)))
		if err != nil {
			log.Errorf("doc center %v md2html error", doc.DocId)
			continue // 跳过当前doc不做处理
		}
		searchResp := response.DocSearchResp{
			Title: title,
			ContentList: []response.DocSearchContent{
				{
					Title:   title,
					Content: snippet,
					Url:     generateDocCenterUrl(url.PathEscape(strings.TrimPrefix(strings.Replace(doc.DocId, docCenterLocalDir, "", 1), "/"))),
				},
			},
		}
		searchResps = append(searchResps, searchResp)
	}
	return searchResps, nil
}

func GetDocCenterMarkdown(ctx *gin.Context, pathName string) (string, error) {
	pathName, err := url.QueryUnescape(pathName)
	if err != nil {
		return "", grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_file_unescape", err.Error())
	}
	// check fileName
	if !strings.HasSuffix(pathName, ".md") {
		return "", grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_file_md", pathName)
	}
	parts := strings.SplitN(pathName, "/", 2)
	if len(parts) == 0 {
		return "", grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_path_err", pathName)
	}
	docDir := parts[0]
	// 加锁方法
	if err := checkAndInitDocCenter(ctx, docCenterLocalDir, docDir); err != nil {
		return "", err
	}
	// read file
	filePath := path.Join(docCenterLocalDir, pathName)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read %v err: %v", pathName, err)
	}
	return string(b), nil
}

func addMdFileToDocMenu(menus *[]*DocMenu, rest string, mdInfo MdInfo) {
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) == 0 {
		return
	}
	var menu *DocMenu
	for _, curr := range *menus {
		if curr.Name == parts[0] {
			menu = curr
			break
		}
	}
	if menu == nil {
		menu = &DocMenu{
			Name: parts[0],
		}
		if len(parts) == 1 { // 非目录，是md文件
			menu.Name = strings.TrimSuffix(menu.Name, ".md")
			menu.RelativePath = mdInfo.relPath
			menu.FilePath = mdInfo.filePath
		}
		*menus = append(*menus, menu)
	}
	if len(parts) > 1 {
		addMdFileToDocMenu(&menu.Children, parts[1], mdInfo)
	}
}

func convertDocMenuToResponse(menu *DocMenu, index string) response.DocMenu {
	var fileName string
	if menu.FilePath != "" {
		parts := strings.SplitN(menu.FilePath, "/", 2)
		if len(parts) > 1 {
			fileName = parts[1]
		}
	}
	ret := response.DocMenu{
		Name:    menu.Name,
		Index:   index,
		Path:    url.PathEscape(fileName), // 前端要求做path转义
		PathRaw: fileName,
	}
	for i, child := range menu.Children {
		ret.Children = append(ret.Children, convertDocMenuToResponse(child, fmt.Sprintf("%s-%d", index, i+1)))
	}
	return ret
}

func getMarkdownSnippet(content, keyword string, snippetLen int) string {
	//string就是只读的采用utf8编码的字节切片(slice) 因此用len函数获取到的长度并不是字符个数，而是字节个数。
	//rune是int32的别名，代表字符的Unicode编码，采用4个字节存储，将string转成rune就意味着任何一个字符都用4个字节来存储其unicode值，
	//这样每次遍历的时候返回的就是unicode值，而不再是字节。
	runes := []rune(content)
	keyRunes := []rune(keyword)
	index := strings.Index(content, keyword)
	if index == -1 {
		if len(runes) < snippetLen {
			return string(runes)
		} else {
			return string(runes[0:snippetLen])
		}
	}
	runeIndex := len([]rune(content[:index]))
	start := runeIndex - snippetLen
	if start < 0 {
		start = 0
	}
	end := runeIndex + len(keyRunes) + snippetLen
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// 递归查找第一个非空的 filePath
func findDocCenterFirstFilePath(nodes []*DocMenu) (string, bool) {
	for _, node := range nodes {
		if node.FilePath != "" {
			return node.FilePath, true
		}
		// 如果当前节点有子节点，递归查找
		if len(node.Children) > 0 {
			if path, found := findDocCenterFirstFilePath(node.Children); found {
				return path, true
			}
		}
	}
	return "", false
}

func checkAndInitDocCenter(ctx *gin.Context, docCenterLocalDir, docDir string) (retErr error) {
	docMu.Lock()
	defer docMu.Unlock()
	defer func() {
		if retErr != nil {
			if err := os.RemoveAll(path.Join(docCenterLocalDir, docDir)); err != nil {
				log.Errorf("clean doc center %v err: %v", docDir, err)
			}
		}
	}()
	if _, err := os.Stat(path.Join(docCenterLocalDir, docDir)); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// 初始化搜索引擎
		se, err := newDocSearcher(docCenterLocalDir, docDir)
		if err != nil {
			return grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_doc_center_search_init_err", err.Error())
		}
		// 以版本为key添加搜索引擎
		addDocSearcher(docDir, se)
	}
	return nil
}

// 将markdown文本中图片引用 ![](xxxxx )与链接引用[](xxxxx )里的 xxxxx 处理为前端可访问的地址
//
//nolint:staticcheck
func convertMarkdown(mdFilePath, objectName, mdContent string) string {
	convertHttp := mdLinkRegex.ReplaceAllStringFunc(mdContent, func(mdLabel string) string {
		for _, httpRelPaths := range mdParenthesisRefRegex.FindAllStringSubmatch(mdLabel, -1) {
			if len(httpRelPaths) <= 1 {
				return mdLabel
			}
			txt := mdBracketRegex.FindString(mdLabel)
			return txt + "(" + url.PathEscape(path.Join(objectName, "../", httpRelPaths[1])) + ")"
		}
		return mdLabel
	})
	return mdImageRegex.ReplaceAllStringFunc(convertHttp, func(imageLabel string) string {
		// imageLabel是匹配到的图片格式，例如 ![](../assets/append.png)
		for _, imageRelPaths := range mdParenthesisRefRegex.FindAllStringSubmatch(imageLabel, -1) {
			// 从imageLabel中继续匹配，例如imageRelPaths[0] 是 ![](../assets/append.png)，imageRelPaths[1]是 ../assets/append.png
			if len(imageRelPaths) <= 1 {
				return imageLabel
			}
			// 重新生成图片引用，将 ../assets/append.png 处理为 service/api/v1/doc-center/assets/append.png
			// 例如mdFilePath是doc-center/v2.1.0/tips/go-append.md
			// 1. "doc-center/v2.1.0/tips/go-append.md" + "../" + "../assets/append.png" => doc-center/v2.1.0/assets/append.png
			// 2. "../../../service/api/v1" + "doc-center/v2.1.0/assets/append.png" => ../../../service/api/v1/doc-center/v2.1.0/assets/append.png
			// 3. 对路径中的非数字字母等做转义，再将 %2F 转回 /
			return "![](" + strings.ReplaceAll(url.PathEscape(path.Join(docCenterStaticPrefix, path.Join(mdFilePath, "../", imageRelPaths[1]))), "%2F", "/") + ")"
		}
		return imageLabel
	})
}

// --- doc-center search engine ---

func newDocSearcher(docCenterLocalDir, docDir string) (*riot.Engine, error) {
	if se, ok := getDocSearcher(docDir); ok {
		return se, nil
	}
	engine := &riot.Engine{}
	engine.Init(types.EngineOpts{
		Using:   3,    // 使用内存索引
		GseDict: "zh", // 指定中文分词字典
	})
	// 批量加载文件索引
	if err := filepath.Walk(path.Join(docCenterLocalDir, docDir), func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查不为目录并且是否为markdown文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// 读取文件内容
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("读取文件%v错误", filePath)
			}
			// 创建索引
			engine.Index(filePath, types.DocData{
				Content: string(content),
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// 刷新索引
	engine.Flush()
	return engine, nil
}

// 获取 searchers 中的指定版本引擎
func getDocSearcher(docDir string) (*riot.Engine, bool) {
	docSearchMu.RLock()
	defer docSearchMu.RUnlock()
	engine, exists := docSearchers[docDir]
	return engine, exists
}

// 添加新的搜索引擎到 searchers
func addDocSearcher(docDir string, engine *riot.Engine) {
	docSearchMu.Lock()
	defer docSearchMu.Unlock()
	docSearchers[docDir] = engine
}

// 检查 searchers 是否已初始化
func isDocSearchersEmpty() bool {
	docSearchMu.RLock()
	defer docSearchMu.RUnlock()
	return len(docSearchers) == 0
}

func generateDocCenterUrl(suffix string) string {
	return "https://" + os.Getenv("SERVER_EXTERNAL_IP") + ":" + os.Getenv("SERVER_EXTERNAL_PORT") + path.Join("/docCenter/pages", suffix)
}

func sortDocMenus(menus *[]*DocMenu) {
	sort.Slice(*menus, func(i, j int) bool {
		return orderDocNum((*menus)[i].Name, (*menus)[j].Name)
	})
	for _, menu := range *menus {
		if len(menu.Children) > 0 {
			sortDocMenus(&menu.Children)
		}
	}
}

// 实现自然排序（数字优先）
func orderDocNum(s1, s2 string) bool {
	numParts1, isNum1 := extractDocNum(s1)
	numParts2, isNum2 := extractDocNum(s2)
	if isNum1 && isNum2 {
		return numParts1 < numParts2
	} else if isNum1 {
		// 如果一个是数字，一个是非数字，数字部分排在前
		return true
	} else if isNum2 {
		return false
	}
	return s1 < s2
}

// extractNum 将字符串按数字和非数字部分分割
func extractDocNum(s string) (int, bool) {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			break
		}
	}
	num, err := strconv.Atoi(result.String())
	if err != nil {
		return 0, false
	}
	return num, true
}
