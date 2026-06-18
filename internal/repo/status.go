package repo

import (
	"strconv"
	"strings"

	"github.com/lhlyu/gitx/internal/git"
)

// Status 是对一个仓库 `git status -b --porcelain` 结果的解析。
type Status struct {
	Branch       string // 分支名；detached 时为 "(detached)"
	Upstream     string // 上游分支，可能为空
	Ahead        int    // 领先上游的提交数
	Behind       int    // 落后上游的提交数
	ChangedFiles int    // 工作区/暂存区改动的文件数
	NoUpstream   bool   // 是否没有配置上游
	Err          error  // 执行 git 失败时非 nil
}

func (s Status) IsClean() bool { return s.ChangedFiles == 0 }

// FirstLine 返回 git 输出中第一行非空内容，用于在失败时给出真实原因。
func FirstLine(out []byte) string {
	for _, line := range strings.Split(string(out), "\n") {
		if t := strings.TrimSpace(line); t != "" {
			return t
		}
	}
	return ""
}

// Inspect 用一次 `git status -b --porcelain` 取得分支与工作区状态。
func Inspect(client *git.Client, path string) Status {
	out, err := client.RunInDir(path, "status", "-b", "--porcelain")
	if err != nil {
		return Status{Err: err}
	}
	return ParseStatus(string(out))
}

// ParseStatus 解析 `git status -b --porcelain` 的输出。
func ParseStatus(out string) Status {
	var s Status
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	for i, line := range lines {
		if i == 0 && strings.HasPrefix(line, "## ") {
			parseHeader(strings.TrimPrefix(line, "## "), &s)
			continue
		}
		if strings.TrimSpace(line) != "" {
			s.ChangedFiles++
		}
	}
	return s
}

// parseHeader 解析分支行：
//
//	main...origin/main [ahead 1, behind 2]
//	main                              （无上游）
//	No commits yet on main
//	HEAD (no branch)
func parseHeader(h string, s *Status) {
	if strings.HasPrefix(h, "HEAD (no branch)") {
		s.Branch = "(detached)"
		s.NoUpstream = true
		return
	}
	if rest, ok := strings.CutPrefix(h, "No commits yet on "); ok {
		s.Branch = strings.TrimSpace(rest)
		s.NoUpstream = true
		return
	}

	// 拆出 "[ahead x, behind y]" 部分
	track := ""
	if idx := strings.Index(h, " ["); idx >= 0 {
		track = h[idx+2:]
		track = strings.TrimSuffix(track, "]")
		h = h[:idx]
	}

	if idx := strings.Index(h, "..."); idx >= 0 {
		s.Branch = h[:idx]
		s.Upstream = h[idx+3:]
	} else {
		s.Branch = strings.TrimSpace(h)
		s.NoUpstream = true
	}

	for _, part := range strings.Split(track, ",") {
		part = strings.TrimSpace(part)
		if n, ok := strings.CutPrefix(part, "ahead "); ok {
			s.Ahead, _ = strconv.Atoi(strings.TrimSpace(n))
		} else if n, ok := strings.CutPrefix(part, "behind "); ok {
			s.Behind, _ = strconv.Atoi(strings.TrimSpace(n))
		}
	}
}
