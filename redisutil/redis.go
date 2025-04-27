package redisutil

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/vivasoft-ltd/golang-course-utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	utils "github.com/vivasoft-ltd/golang-course-utils/methods"
)

type Redis struct {
	Prefix      string
	RedisClient *redis.Client
}

/*
Connect method takes the redis credentials and prefix as input. It's then
connect to redis instance and return Redis util object otherwise create panic
*/
func Connect(host, port, pass string, db int, prefix string) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       db,
	})
	logger.Info("connecting to redis at ", host, ":", port, "...")
	if _, err := redisClient.Ping().Result(); err != nil {
		logger.Error("failed to connect redis: ", err)
		panic(err)
	}
	logger.Info("redis connection successful...")
	return &Redis{
		RedisClient: redisClient,
		Prefix:      prefix,
	}

}

func (r *Redis) Set(key string, value interface{}, ttl int) error {
	key = r.getKey(key)
	if utils.IsEmpty(key) || utils.IsEmpty(value) {
		return errutil.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.RedisClient.Set(key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (r *Redis) SetString(key string, value string, ttl int) error {
	key = r.getKey(key)
	if utils.IsEmpty(key) || utils.IsEmpty(value) {
		return errutil.ErrEmptyRedisKeyValue
	}

	return r.RedisClient.Set(key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *Redis) SetStruct(key string, value interface{}, ttl time.Duration) error {
	key = r.getKey(key)
	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.RedisClient.Set(key, string(serializedValue), ttl*time.Second).Err()
}

func (r *Redis) Get(key string) (string, error) {
	key = r.getKey(key)
	if utils.IsEmpty(key) {
		return "", errutil.ErrEmptyRedisKeyValue
	}

	return r.RedisClient.Get(key).Result()
}

func (r *Redis) GetInt(key string) (int, error) {
	key = r.getKey(key)
	if utils.IsEmpty(key) {
		return 0, errutil.ErrEmptyRedisKeyValue
	}

	str, err := r.RedisClient.Get(key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}
func (r *Redis) GetStruct(key string, outputStruct interface{}) error {
	key = r.getKey(key)
	if utils.IsEmpty(key) {
		return errutil.ErrEmptyRedisKeyValue
	}

	serializedValue, err := r.RedisClient.Get(key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (r *Redis) HasKey(key string) bool {
	key = r.getKey(key)
	exists, err := r.RedisClient.Exists(key).Result()
	if err != nil {
		return false
	}

	return exists == 1
}
func (r *Redis) Exists(key string) bool {
	return r.HasKey(key)
}

func (r *Redis) IncBy(key string, value int) error {
	key = r.getKey(key)
	return r.RedisClient.IncrBy(key, int64(value)).Err()
}
func (r *Redis) INCR(key string) error {
	key = r.getKey(key)
	return r.RedisClient.Incr(key).Err()
}

func (r *Redis) Del(keys ...string) error {
	newKey := []string{}
	for _, v := range keys {
		v = r.getKey(v)
		newKey = append(newKey, v)
	}
	return r.RedisClient.Del(newKey...).Err()
}

func (r *Redis) DelPattern(pattern string) error {
	pattern = r.getKey(pattern)
	iter := r.RedisClient.Scan(0, pattern, 0).Iterator()

	for iter.Next() {
		err := r.RedisClient.Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) getKey(key string) string {
	return r.Prefix + key
}
