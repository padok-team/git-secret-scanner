package progress

import (
	"context"
	"fmt"
	"math"
	"strings"
)

const defaultViewWidth = 40

type progressBar struct {
	description string
	width       int
}

func newProgressBar(description string) progressBar {
	return progressBar{
		description: description,
		width:       defaultViewWidth,
	}
}

func (p progressBar) showBar(b *strings.Builder, percent int) {
	tw := p.width                                               // total width
	fw := int(math.Round(float64(tw) * float64(percent) / 100)) // filled width

	fw = max(0, min(tw, fw))

	// Solid fill
	b.WriteString(strings.Repeat("█", fw))

	// Empty fill
	n := max(0, tw-fw)
	b.WriteString(strings.Repeat("░", n))
}

func (p progressBar) showPercentage(b *strings.Builder, percent int) {
	b.WriteString(" " + fmt.Sprintf("%d%%", percent))
}

func (p progressBar) View(percent int) string {
	b := strings.Builder{}

	b.WriteString(p.description + " ")
	p.showBar(&b, percent)
	p.showPercentage(&b, percent)

	return b.String()
}

func RunTasksWithProgressBar(ctx context.Context, description string, tasks []*Task, concurrency int) error {
	bar := newProgressBar(description)

	pool := newTaskPool(ctx, concurrency, len(tasks))
	for _, task := range tasks {
		// wrap task to update progress once task is completed
		pool.AddTask(task.Wrap(func(fn TaskFunc) TaskFunc {
			return func(ctx context.Context) error {
				err := fn(ctx)
				fmt.Printf("%s\n\033[1A", bar.View(pool.TasksDoneCount*100/pool.Size()))
				return err
			}
		}))
	}

	// print bar for the first time
	fmt.Printf("%s\n\033[1A", bar.View(0))

	err := pool.Run()

	// clear progress
	fmt.Print("\033[K")

	return err
}
