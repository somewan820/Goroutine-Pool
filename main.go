package main

import (
	Pool "Goroutine-Pool/pool"
	"fmt"
	"time"
)

func main() {
	// 创建任务池
	pool, err := Pool.NewPool(10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		// 任务放入池中
		pool.Put(&Pool.Task{
			Handler: func(v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i},
		})
	}

	time.Sleep(1e9) // 等待执行
}
