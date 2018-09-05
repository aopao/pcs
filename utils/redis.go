package utils

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/astaxie/beego/logs"
)

var (
	redisPool           *redis.Pool
	REDIS_MAX_POOL_SIZE = 20
)

func init() {
	// 从配置文件获取redis配置并连接
	host := beego.AppConfig.String("redis_host")
	db, _ := beego.AppConfig.Int("redis_db")
	password := beego.AppConfig.String("redis_password")
	redisPool = newPool(host, redis.DialPassword(password), redis.DialDatabase(db))
	logs.Info("Redis start init...")
}

func newPool(addr string, dialOption ... redis.DialOption) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     REDIS_MAX_POOL_SIZE,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, dialOption...)
		},
	}
}

func RedisSet(key, value string, expireTime ... int) bool {
	// 操作redis时调用Do方法，第一个参数传入操作名称（字符串），然后根据不同操作传入key、value、数字等
	// 返回2个参数，第一个为操作标识，成功则为1，失败则为0；第二个为错误信息
	rs := redisPool.Get()
	defer rs.Close()
	_, err := rs.Do("SET", key, value)
	// 若操作失败则返回
	if err != nil {
		logs.Error("未知错误，Key：", key, err.Error())
		return false
	}
	if len(expireTime) > 0 {
		n, _ := rs.Do("EXPIRE", key, expireTime)
		if n == int64(1) {
			logs.Info("设置分布式锁失效时间：", expireTime)
		}
	}
	return true
}

func RedisGet(key string) (string, error) {
	rs := redisPool.Get()
	defer rs.Close()
	value, err := redis.String(rs.Do("GET", key))
	// 若操作失败则返回
	if err != nil {
		logs.Error("未知错误，Key：", key, err.Error())
		return "", err
	}
	return value, nil
}

func RedisDel(key string) bool {
	rs := redisPool.Get()
	defer rs.Close()
	_, err := rs.Do("DEL", key)
	// 若操作失败则返回
	if err != nil {
		logs.Error("删除Key出现未知错误，Key：", key, err.Error())
		return false
	}
	return true
}

func GetDistributedLock(key, value string, lockTime int) bool {
	// 返回2个参数，第一个为操作标识，成功则为1，失败则为0；第二个为错误信息
	rs := redisPool.Get()
	defer rs.Close()
	n, err := rs.Do("SETNX", key, value)
	// 若操作失败则返回
	if err != nil {
		logs.Error("未知错误，Key：", key, err.Error())
		return false
	}
	// 返回的n的类型是int64的，所以得将1或0转换成为int64类型的再比较
	if n == int64(1) {
		// 设置过期时间为固定时间
		n, _ := rs.Do("EXPIRE", key, lockTime)
		if n == int64(1) {
			logs.Info("设置分布式锁失效时间：", lockTime)
		} else {
			logs.Error("设置分布式锁失效时间失败！Key:", key)
		}
		return true
	}
	logs.Error("获取分布式锁失败！", key)
	return false
}
