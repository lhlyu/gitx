package repo

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// maxWorkers 限制并发执行的 git 进程数上限，避免在多核机器上一次拉起过多进程。
const maxWorkers = 16

// Target 表示一个待处理的 Git 仓库。
type Target struct {
	// Name 是相对于扫描根目录的展示名（顶层仓库即目录名）。
	Name string
	// Path 是仓库的绝对路径。
	Path string
}

// IsGitRepo 判断 dir 是否为 Git 仓库（存在 .git）。
func IsGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// Scan 从 root 出发扫描 Git 仓库。
//
// depth <= 0 时只检查 root 本身；depth >= 1 时向下扫描对应层数的子目录，
// 一旦某目录是 Git 仓库就停止深入（不扫描仓库内部的子目录）。
func Scan(root string, depth int) []Target {
	if depth <= 0 {
		if IsGitRepo(root) {
			return []Target{{Name: filepath.Base(root), Path: root}}
		}
		return nil
	}

	var targets []Target
	scan(root, root, depth, 0, &targets)
	return targets
}

func scan(root, dir string, maxDepth, curDepth int, out *[]Target) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		if IsGitRepo(path) {
			name, relErr := filepath.Rel(root, path)
			if relErr != nil {
				name = entry.Name()
			}
			*out = append(*out, Target{Name: name, Path: path})
		} else if curDepth < maxDepth-1 {
			scan(root, path, maxDepth, curDepth+1, out)
		}
	}
}

// workers 根据任务数与 CPU 核数计算合适的并发度。
func workers(n int) int {
	w := runtime.NumCPU()
	if w > maxWorkers {
		w = maxWorkers
	}
	if w > n {
		w = n
	}
	if w < 1 {
		w = 1
	}
	return w
}

// Process 以并发方式对每个 target 执行 fn，并按 targets 的原始顺序返回结果。
func Process[T any](targets []Target, fn func(Target) T) []T {
	return process(targets, nil, fn)
}

// ProcessWithProgress 与 Process 相同，但在每个任务完成时刷新一行进度（仅终端环境）。
// label 形如 "拉取中"，输出 "拉取中 [3/20]"。
func ProcessWithProgress[T any](targets []Target, label string, fn func(Target) T) []T {
	p := newProgress(label, len(targets))
	results := process(targets, p.inc, fn)
	p.finish()
	return results
}

func process[T any](targets []Target, onDone func(), fn func(Target) T) []T {
	results := make([]T, len(targets))
	if len(targets) == 0 {
		return results
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	for range workers(len(targets)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = fn(targets[i])
				if onDone != nil {
					onDone()
				}
			}
		}()
	}

	for i := range targets {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return results
}
