package main

import (
	"sync"
	"testing"
)

// 实现一个线程安全的单例
type Singleton struct{}

var instance *Singleton
var once sync.Once

func GetSingleton() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
func TestConcurrency13(t *testing.T) {

}
