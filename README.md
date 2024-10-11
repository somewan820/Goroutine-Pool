# goroutine pool
在100w 次执行，原子增量操作下，使用任务池比直接 goroutine 内存分配节省 7000 倍左右, 内存分配次数减少 2700 倍左右
![image](https://github.com/user-attachments/assets/3ec20812-56cf-470c-af0c-3edf08dfa975)
