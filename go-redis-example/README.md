# Redis应用案例

## Example

### 01-基于string Incr的签到案例
```shell
go run main.go Ex01 1165894833417101
```

### 02-基于SETNX的分布式锁
```shell
go run main.go Ex02
```

### 03-基于Incr、Decr的简单限流器
```shell
go run main.go Ex03
```

### 04-基于List的消息队列
```shell
go run main.go Ex04
```

### 05-基于Hash的计数器
```shell
go run main.go Ex05 init # 初始化用户计数值
go run main.go Ex05 get 1556564194374926  # 打印用户(1556564194374926)的所有计数值
go run main.go Ex05 incr_like 1556564194374926 # 点赞数+1
go run main.go Ex05 incr_collect 1556564194374926 # 收藏数+1
go run main.go Ex05 decr_like 1556564194374926 # 点赞数-1
go run main.go Ex05 decr_collect 1556564194374926 # 收藏数-1
```

### 06-基于Zset的排行榜
```shell
go run main.go Ex06 init # 初始化积分
go run main.go Ex06 rev_order # 输出完整榜单
go run main.go Ex06 order_page 1 # 逆序分页输出，page=1
go run main.go Ex06 get_rank user2 # 获取user2的排名
go run main.go Ex06 get_score user2 # 获取user2的分数
go run main.go Ex06 add_user_score user2 10 # 为user2增加10分
```

### 07-基于PubSub的消息订阅
```shell
go run main.go Ex07
```