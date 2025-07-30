package ahocorasick

import (
	"fmt"
	"sync"
)

var (
	_acs *acsConfig
)

type dict struct {
	DictID  string
	Version string
	Matcher *Matcher
	Reply   string
}

type acsConfig struct {
	mu    sync.RWMutex
	dicts map[string]*dict
}

type MatchResult struct {
	DictCfg   DictConfig
	WordIndex int
	Word      string
	Reply     string
}

type DictConfig struct {
	DictID  string
	Version string
}

type DictStatus struct {
	DictCfg DictConfig
	Status  bool // 词表是否存在
}

// 初始化
func Init() error {
	if _acs != nil {
		return fmt.Errorf("acs client already init")
	}
	_acs = &acsConfig{
		dicts: make(map[string]*dict),
	}
	return nil
}

// CheckDictStatus 检查多个字典的状态
func CheckDictStatus(dicts []DictConfig) ([]DictStatus, error) {
	if _acs == nil {
		return nil, fmt.Errorf("aho-corasick not initialized")
	}
	if len(dicts) == 0 {
		return nil, fmt.Errorf("dict config cannot be empty")
	}
	_acs.mu.RLock()
	defer _acs.mu.RUnlock()
	result := make([]DictStatus, 0, len(dicts))
	for _, cfg := range dicts {
		if cfg.Version == "" || cfg.DictID == "" {
			return nil, fmt.Errorf("dict config can not be empty")
		}
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		dict, exists := _acs.dicts[dictKey]
		result = append(result, DictStatus{
			DictCfg: cfg,
			Status:  exists && dict.Version == cfg.Version,
		})
	}
	return result, nil
}

func BuildDict(dc DictConfig, reply string, words []string) error {
	if _acs == nil {
		return fmt.Errorf("aho-corasick service not initialized")
	}
	if dc.DictID == "" || dc.Version == "" {
		return fmt.Errorf("dict config can not be empty")
	}

	dictKey := getDictKey(dc.DictID, dc.Version)
	matcher := NewStringMatcher(words)
	_acs.mu.Lock()
	defer _acs.mu.Unlock()
	_acs.dicts[dictKey] = &dict{
		DictID:  dc.DictID,
		Version: dc.Version,
		Matcher: matcher,
		Reply:   reply,
	}
	return nil
}

func ContentMatch(content string, dicts []DictConfig, returnFirstMatch bool) ([]MatchResult, error) {
	if _acs == nil || len(content) == 0 || len(dicts) == 0 {
		return nil, fmt.Errorf("invalid parameters: accessor is nil or content/dicts are empty")
	}
	// 1. 加读锁，仅保护字典读取阶段
	_acs.mu.RLock()
	dictsCopy := make(map[string]*dict) // 复制字典，避免后续操作依赖锁
	for _, cfg := range dicts {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := _acs.dicts[dictKey]; exists && dict != nil {
			dictsCopy[dictKey] = dict
		}
	}
	_acs.mu.RUnlock() // 字典读取完毕，立即释放锁
	// 2. 无锁状态下进行匹配
	results := make([]MatchResult, 0)
	contentBytes := []byte(content)
	for _, cfg := range dicts {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := dictsCopy[dictKey]; exists {
			matches := dict.Matcher.MatchThreadSafe(contentBytes)
			if len(matches) > 0 {
				for _, idx := range matches {
					matchResult := MatchResult{
						DictCfg: DictConfig{
							DictID:  dict.DictID,
							Version: dict.Version,
						},
						WordIndex: idx,
						Word:      dict.Matcher.getOriginalWord(idx),
						Reply:     dict.Reply,
					}
					results = append(results, matchResult)
					if returnFirstMatch {
						return results, nil
					}
				}
			}
		}
	}
	return results, nil
}

func ContentContain(content string, dicts []DictConfig) (bool, error) {
	if _acs == nil || len(content) == 0 || len(dicts) == 0 {
		return false, fmt.Errorf("invalid parameters: accessor is nil or content/dictIDs are empty")
	}
	// 1. 加读锁，仅保护字典读取阶段
	_acs.mu.RLock()
	dictsCopy := make(map[string]*dict) // 复制字典，避免后续操作依赖锁
	for _, cfg := range dicts {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := _acs.dicts[dictKey]; exists && dict != nil {
			dictsCopy[dictKey] = dict
		}
	}
	_acs.mu.RUnlock() // 字典读取完毕，立即释放锁
	// 2. 无锁状态匹配字段
	contentBytes := []byte(content)
	for _, cfg := range dictsCopy {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := dictsCopy[dictKey]; exists && dict.Matcher.Contains(contentBytes) {
			return true, nil
		}
	}
	return false, nil
}

// =--- internal ---
func (m *Matcher) getOriginalWord(index int) string {
	for _, node := range m.trie {
		if node.output && node.index == index {
			return string(node.b)
		}
	}
	return ""
}

// getDictKey 生成字典的唯一键
func getDictKey(dictID, version string) string {
	return fmt.Sprintf("%s-%s", dictID, version)
}
