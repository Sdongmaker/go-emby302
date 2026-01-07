package path

import (
	"fmt"
)

// OpenlistPathRes 路径转换结果（已废弃，保留用于兼容性）
type OpenlistPathRes struct {
	Success bool
	Path    string
	Range   func() ([]string, error)
}

// Emby2Openlist Emby 资源路径转 Openlist 资源路径（已废弃）
// 此函数已不再使用，保留用于兼容性
func Emby2Openlist(embyPath string) OpenlistPathRes {
	return OpenlistPathRes{
		Success: false,
		Path:    embyPath,
		Range: func() ([]string, error) {
			return nil, fmt.Errorf("path service is deprecated")
		},
	}
}

// SplitFromSecondSlash 找到给定字符串 str 中第二个 '/' 字符的位置（已废弃）
func SplitFromSecondSlash(str string) (string, error) {
	return "", fmt.Errorf("path service is deprecated")
}
