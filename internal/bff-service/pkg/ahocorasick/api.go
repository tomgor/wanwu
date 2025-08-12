package ahocorasick

import (
	"fmt"
	"sync"
)

var (
	_ac *acMgr
)

type acDict struct {
	DictID  string
	Version string
	Reply   string

	Matcher *Matcher
}

type acMgr struct {
	mu    sync.RWMutex
	dicts map[string]*acDict
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
	if _ac != nil {
		return fmt.Errorf("aho-corasick already init")
	}
	_ac = &acMgr{
		dicts: make(map[string]*acDict),
	}
	return nil
}

// CheckDictStatus 检查多个字典的状态
func CheckDictStatus(dicts []DictConfig) ([]DictStatus, error) {
	if _ac == nil {
		return nil, fmt.Errorf("aho-corasick not init")
	}
	if len(dicts) == 0 {
		return nil, nil
	}
	_ac.mu.RLock()
	defer _ac.mu.RUnlock()
	result := make([]DictStatus, 0, len(dicts))
	for _, cfg := range dicts {
		if cfg.Version == "" || cfg.DictID == "" {
			return nil, fmt.Errorf("dict config id or version can not be empty")
		}
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		_, exists := _ac.dicts[dictKey]
		result = append(result, DictStatus{
			DictCfg: cfg,
			Status:  exists,
		})
	}
	return result, nil
}

func BuildDict(dict DictConfig, reply string, words []string) error {
	if _ac == nil {
		return fmt.Errorf("aho-corasick not init")
	}
	if dict.DictID == "" || dict.Version == "" {
		return fmt.Errorf("dict config can not be empty")
	}

	dictKey := getDictKey(dict.DictID, dict.Version)
	matcher := NewStringMatcher(words)
	_ac.mu.Lock()
	defer _ac.mu.Unlock()
	_ac.dicts[dictKey] = &acDict{
		DictID:  dict.DictID,
		Version: dict.Version,
		Matcher: matcher,
		Reply:   reply,
	}
	return nil
}

func ContentMatch(content string, dicts []DictConfig, returnFirstMatch bool) ([]MatchResult, error) {
	if _ac == nil {
		return nil, fmt.Errorf("aho-corasick not init")
	}
	if len(content) == 0 || len(dicts) == 0 {
		return nil, nil
	}
	// 1. 加读锁，仅保护字典读取阶段
	_ac.mu.RLock()
	dictsCopy := make(map[string]*acDict) // 复制字典，避免后续操作依赖锁
	for _, cfg := range dicts {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := _ac.dicts[dictKey]; exists && dict != nil {
			dictsCopy[dictKey] = dict
		}
	}
	_ac.mu.RUnlock() // 字典读取完毕，立即释放锁
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

func ContentContain(content string, dicts []DictConfig) (*DictConfig, error) {
	if _ac == nil {
		return nil, fmt.Errorf("aho-corasick not init")
	}
	if len(content) == 0 || len(dicts) == 0 {
		return nil, nil
	}
	// 1. 加读锁，仅保护字典读取阶段
	_ac.mu.RLock()
	dictsCopy := make(map[string]*acDict) // 复制字典，避免后续操作依赖锁
	for _, cfg := range dicts {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := _ac.dicts[dictKey]; exists && dict != nil {
			dictsCopy[dictKey] = dict
		}
	}
	_ac.mu.RUnlock() // 字典读取完毕，立即释放锁
	// 2. 无锁状态匹配字段
	contentBytes := []byte(content)
	for _, cfg := range dictsCopy {
		dictKey := getDictKey(cfg.DictID, cfg.Version)
		if dict, exists := dictsCopy[dictKey]; exists && dict.Matcher.Contains(contentBytes) {
			return &DictConfig{DictID: dict.DictID, Version: dict.Version}, nil
		}
	}
	return nil, nil
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
