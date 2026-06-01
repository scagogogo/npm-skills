package models

// Script 表示 NPM 包的脚本命令定义
//
// 使用 map[string]string 类型以支持 npm 包中定义的任意脚本命令，
// 例如 "build"、"lint"、"dev" 等非标准脚本。
// 最常用的脚本命令包括:
//   - "test": 测试脚本命令
//   - "start": 启动项目脚本命令
//   - "build": 构建脚本命令
//   - "lint": 代码检查脚本命令
type Script map[string]string
