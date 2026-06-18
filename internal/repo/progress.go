package repo

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

// progress 在任务执行期间向 stderr 刷新一行 "label [done/total]"。
// 非终端环境（重定向、管道、CI）下完全静默，避免污染输出。
type progress struct {
	label string
	total int
	done  atomic.Int64
	mu    sync.Mutex
	tty   bool
}

func newProgress(label string, total int) *progress {
	return &progress{
		label: label,
		total: total,
		tty:   isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()),
	}
}

func (p *progress) inc() {
	if !p.tty {
		return
	}
	n := p.done.Add(1)
	p.mu.Lock()
	fmt.Fprintf(os.Stderr, "\r\033[K%s [%d/%d]", p.label, n, p.total)
	p.mu.Unlock()
}

// finish 清除进度行，让后续正式输出从干净的一行开始。
func (p *progress) finish() {
	if !p.tty {
		return
	}
	fmt.Fprint(os.Stderr, "\r\033[K")
}
