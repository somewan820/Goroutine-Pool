package pool

import (
	"errors"
	"log"
	"time"
)

var ErrInvalidPoolCap = errors.New("invalid pool cap")

const (
	RUNNING = 1
	STOPED  = 0
)

func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		// 初始化任务队列, 队列长度为容量
		chTask: make(chan *Task, capacity),
	}, nil
}

func (p *Pool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("Worker panic: %s\n", r)
				}
			}
			p.checkWorker() // worker 退出时检测是否有可运行的 worker
		}()

		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()

	if p.status == STOPED { // 如果任务池处于关闭状态, 再 put 任务会返回 ErrPoolAlreadyClosed 错误
		return ErrPoolAlreadyClosed
	}

	// run worker
	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}

	// send task
	if p.status == RUNNING {
		p.chTask <- task
	}

	return nil
}

var ErrPoolAlreadyClosed = errors.New("pool already closed")

func (p *Pool) Close() {
	p.setStatus(STOPED) // 设置 status 为已停止

	for len(p.chTask) > 0 { // 阻塞等待所有任务被 worker 消费
		time.Sleep(1e6) // 防止等待任务清空 cpu 负载突然变大, 这里小睡一下
	}

	close(p.chTask) // 关闭任务队列
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Unlock()

	// 当前没有 worker 且有任务存在，运行一个 worker 消费任务
	// 没有任务无需考虑 (当前 Put 不会阻塞，下次 Put 会启动 worker)
	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}
