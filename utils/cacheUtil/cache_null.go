package cacheUtil

import "time"

/*
@Author: wangc293
@Time: 2020-03-30 10:02
@Description:
*/

type NullCache struct {

}

func (bm *NullCache) Get(key string) interface{} {
	return nil
}

func (bm *NullCache)GetMulti(keys []string) []interface{} {
	return nil
}

func (bm *NullCache)Put(key string, val interface{}, timeout time.Duration) error {
	return nil
}
func (bm *NullCache)Delete(key string) error {
	return nil
}
func (bm *NullCache)Incr(key string) error {
	return nil
}
func (bm *NullCache)Decr(key string) error {
	return nil
}
func (bm *NullCache)IsExist(key string) bool {
	return false
}
func (bm *NullCache)ClearAll() error{
	return nil
}

func (bm *NullCache)StartAndGC(config string) error {
	return nil
}
