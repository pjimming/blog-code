# !bin/bash
go run main.go Ex05 init # 初始化用户计数值
go run main.go Ex05 get 1556564194374926  # 打印用户(1556564194374926)的所有计数值
go run main.go Ex05 incr_like 1556564194374926 # 点赞数+1
go run main.go Ex05 incr_collect 1556564194374926 # 收藏数+1
go run main.go Ex05 decr_like 1556564194374926 # 点赞数-1
go run main.go Ex05 decr_collect 1556564194374926 # 收藏数-1
go run main.go Ex05 decr_collect 1556564194374926 # 收藏数-1
