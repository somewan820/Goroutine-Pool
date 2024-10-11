package main

import (
	Pool "Goroutine-Pool/pool"
	"sync"
	"sync/atomic"
	"testing"
)

var wg = sync.WaitGroup{}

var sum int64

func demoTask(v ...interface{}) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

var runTimes = 1000000

// 原生 goroutine
func BenchmarkGoroutineTimeLifeSetTimes(b *testing.B) {
	b.ReportAllocs() // 启用内存分配报告
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go demoTask()
	}
	wg.Wait() // 等待执行完毕
}

// 使用协程池
func BenchmarkPoolTimeLifeSetTimes(b *testing.B) {
	b.ReportAllocs() // 启用内存分配报告
	pool, err := Pool.NewPool(20)
	if err != nil {
		b.Error(err)
	}

	task := &Pool.Task{
		Handler: demoTask,
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		pool.Put(task)
	}

	wg.Wait() // 等待执行完毕
}
