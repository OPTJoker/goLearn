package util

import (
	"os"
	"path/filepath"
)

// ProjectConfig 项目配置结构
type ProjectConfig struct {
	RootDir   string
	WebDir    string
	DocsDir   string
	ScriptDir string
}

// GetProjectConfig 获取项目配置
func GetProjectConfig() *ProjectConfig {
	rootDir := GetProjectRoot()

	return &ProjectConfig{
		RootDir:   rootDir,
		WebDir:    filepath.Join(rootDir, "web"),
		DocsDir:   filepath.Join(rootDir, "docs"),
		ScriptDir: filepath.Join(rootDir, "script"),
	}
}

// GetProjectRoot 获取项目根目录
func GetProjectRoot() string {
	// 方法1: 检查环境变量 PROJECT_ROOT
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		return projectRoot
	}

	// 方法2: 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	// 如果从src目录运行，则向上一级到项目根目录
	if filepath.Base(wd) == "src" {
		return filepath.Dir(wd)
	}

	return wd
}

// GetWebDir 获取Web目录路径
func GetWebDir() string {
	return filepath.Join(GetProjectRoot(), "web")
}

// GetDocsDir 获取文档目录路径
func GetDocsDir() string {
	return filepath.Join(GetProjectRoot(), "docs")
}

// GetScriptDir 获取脚本目录路径
func GetScriptDir() string {
	return filepath.Join(GetProjectRoot(), "script")
}
