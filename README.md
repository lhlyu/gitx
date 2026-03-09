# gitx

用 Go 封装的 Git 命令行工具，提供更友好的交互体验 🚀

## 📦 安装

```bash
go install github.com/lhlyu/gitx@latest
```

## 📖 命令列表

| 命令                   | 说明                                      |
|----------------------|-----------------------------------------|
| `gitx info`          | 显示仓库信息（分支、远程地址、工作区状态等）                  |
| `gitx list [depth]`  | 列出指定深度的 Git 项目及其工作区状态（默认 depth=1）       |
| `gitx pull [depth]`  | 拉取最新代码并显示每个项目是否更新（默认 depth=0，表示只拉取当前目录） |
| `gitx undo [depth]`  | 撤销工作区和暂存区的修改并显示结果（默认 depth=0）           |
| `gitx clean [depth]` | 清理仓库并显示清理结果，重置到最新提交状态 ⚠️（默认 depth=0）    |
| `gitx reset <steps>` | 将当前仓库硬重置到前 N 个提交，即执行 `git reset --hard HEAD~N` ⚠️ |

## 🛠️ 开发

### 项目结构

```
gitx/
├── cmd/              # 命令定义和注册
│   ├── root.go
│   ├── info.go
│   ├── list.go
│   ├── pull.go
│   ├── undo.go
│   ├── clean.go
│   └── reset.go
├── internal/         # 内部实现
│   ├── git/         # Git 客户端封装
│   ├── info/        # info 命令实现
│   ├── list/        # list 命令实现
│   ├── pull/        # pull 命令实现
│   ├── undo/        # undo 命令实现
│   ├── clean/       # clean 命令实现
│   └── reset/       # reset 命令实现
└── main.go          # 程序入口
```

### 添加新功能

遵循 feature 划分原则：

1. 在 `cmd/` 目录下创建命令定义文件
2. 在 `internal/{feature}/` 目录下创建具体实现
3. 使用 `run.go` 作为 feature 的入口
