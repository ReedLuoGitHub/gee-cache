分布式缓存
- 最简单的缓存莫过于存储在内存中的k-v键值对缓存
- 但k-v键值对缓存存在问题
    1. 内存不足
    2. 并发写入冲突（Go语言map并发不安全）
    3. 单机性能低

gee-cache主要实现功能
1. 单机缓存和基于 HTTP 的分布式缓存 
2. 最近最少访问(Least Recently Used, LRU) 缓存策略 
3. 使用 Go 锁机制防止缓存击穿 
4. 使用一致性哈希选择节点，实现负载均衡
5. 使用 protobuf 优化节点间二进制通信
6. …

