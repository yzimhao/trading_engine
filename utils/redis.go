package utils

import "github.com/gomodule/redigo/redis"

func ScanRedisKeys(conn redis.Conn, cursor int, pattern string) ([]string, error) {
	const count = 100 // 一次扫描的数量

	var keys []string
	for {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern, "COUNT", count))
		if err != nil {
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &keys)
		if err != nil {
			return nil, err
		}

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
