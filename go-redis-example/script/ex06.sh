# !bin/bash
go run main.go Ex06 init # 初始化积分
go run main.go Ex06 rev_order # 输出完整榜单
go run main.go Ex06 order_page 1 # 逆序分页输出，page=1
go run main.go Ex06 order_page 2 # 逆序分页输出，page=2
go run main.go Ex06 get_rank user2 # 获取user2的排名
go run main.go Ex06 get_score user2 # 获取user2的分数
go run main.go Ex06 add_user_score user2 10 # 为user2增加10分
go run main.go Ex06 get_rank user2 # 获取user2的排名
go run main.go Ex06 get_score user2 # 获取user2的分数
