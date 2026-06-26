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
| `gitx branch`        | 查看当前仓库的本地分支和远程跟踪分支（基于本地 refs，不主动联网）    |
| `gitx log [n]`       | 查看当前仓库最近 N 条提交记录（默认 n=5，n 必须为正整数）       |
| `gitx search <keyword> [--all]` | 快速查找关键字出现在哪些文件和行号里（依赖 `rg`，按字面量匹配；默认关键字至少 2 个字符、最多显示 20 条，`--all` 解除限制） |
| `gitx list [depth]`  | 列出指定深度的 Git 项目及其工作区状态（默认 depth=1）       |
| `gitx status [depth]`| 查看分支、领先/落后远程的提交数及工作区状态（默认 depth=0）      |
| `gitx pull [depth]`  | 拉取最新代码并显示每个项目是否更新（默认 depth=0，表示只拉取当前目录） |
| `gitx undo [depth]`  | 撤销工作区和暂存区的修改并显示结果（默认 depth=0）           |
| `gitx clean [depth]` | 清理仓库并显示清理结果，重置到最新提交状态 ⚠️（默认 depth=0）    |
| `gitx reset <steps>` | 将当前仓库硬重置到前 N 个提交，即执行 `git reset --hard HEAD~N` ⚠️ |

> 多仓库操作（`list`/`status`/`pull`/`clean`/`undo` 在 depth ≥ 1 时）会按 CPU 核数并发执行，
> 执行期间在终端显示进度（`[3/20]`），重定向或管道时自动静默。
> `status` 的领先/落后基于上次 `fetch` 的远程信息，不会主动联网。
> 扫描到的嵌套仓库以「相对扫描根目录的路径」显示（如 `nest/c`），顶层仓库则显示目录名。

## 🛠️ 开发

### 项目结构

```
gitx/
├── cmd/              # 命令定义和注册
│   ├── root.go
│   ├── info.go
│   ├── branch.go
│   ├── log.go
│   ├── search.go
│   ├── list.go
│   ├── pull.go
│   ├── undo.go
│   ├── clean.go
│   └── reset.go
├── internal/         # 内部实现
│   ├── git/         # Git 客户端封装
│   ├── repo/        # 仓库扫描、并发执行与状态解析（公共）
│   ├── term/        # 终端显示宽度与列对齐工具（处理中英文混排）
│   ├── info/        # info 命令实现
│   ├── branch/      # branch 命令实现
│   ├── log/         # log 命令实现
│   ├── search/      # search 命令实现
│   ├── list/        # list 命令实现
│   ├── status/      # status 命令实现
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
