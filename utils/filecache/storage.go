package filecache

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/yzimhao/trading_engine/utils/app"
)

type Storage struct {
	db *bolt.DB
}

var obj *Storage

func NewStorage(filename string, timeout time.Duration) *Storage {
	if obj == nil {
		db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: timeout * time.Second})
		if err != nil {
			app.Logger.Panic(fmt.Sprintf("打开%s失败 %s", filename, err))
		}
		obj = &Storage{
			db: db,
		}
	}
	return obj
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Stats() any {
	return s.db.Stats()
}

func (s *Storage) Set(bucket string, key string, raw []byte) {
	s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), raw)
		return err
	})
}

func (s *Storage) Remove(bucket string, key string) {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}

func (s *Storage) Find(bucket string, key string) [][]byte {
	data := make([][]byte, 0)
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if key == "" {
				data = append(data, v)
			} else {
				if key == string(k) {
					data = append(data, v)
				}
			}
		}

		return nil
	})
	return data
}

func (s *Storage) Get(bucket string, key string) ([]byte, bool) {
	data := []byte{}
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		data = b.Get([]byte(key))
		return nil
	})
	if data == nil {
		return nil, false
	}
	return data, true
}
