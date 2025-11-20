package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 实现一个线程安全的队列
type Queue[T any] struct {
	items []T
	mu    *sync.Mutex
	cond  *sync.Cond
}

func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{}
	q.mu = new(sync.Mutex)
	q.cond = sync.NewCond(q.mu)
	return q
}

func (q *Queue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *Queue[T]) Pop() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait() // 陷入阻塞并释放锁
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func TestConcurrency21(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		q := NewQueue[int]()

		// 测试 Push 和 Pop
		q.Push(1)
		q.Push(2)
		q.Push(3)

		if val := q.Pop(); val != 1 {
			t.Errorf("期望 Pop 得到 1，实际得到 %d", val)
		}
		if val := q.Pop(); val != 2 {
			t.Errorf("期望 Pop 得到 2，实际得到 %d", val)
		}
		if val := q.Pop(); val != 3 {
			t.Errorf("期望 Pop 得到 3，实际得到 %d", val)
		}

		t.Log("基本 Push/Pop 功能测试通过")
	})

	t.Run("空队列阻塞测试", func(t *testing.T) {
		q := NewQueue[int]()

		var wg sync.WaitGroup
		var result int

		// 启动消费者，会在空队列上阻塞
		wg.Add(1)
		go func() {
			defer wg.Done()
			result = q.Pop()
		}()

		// 等待一会确保 Pop 已经在阻塞
		time.Sleep(50 * time.Millisecond)

		// Push 一个值来唤醒阻塞的 Pop
		q.Push(42)

		wg.Wait()

		if result != 42 {
			t.Errorf("期望 Pop 得到 42，实际得到 %d", result)
		}
		t.Log("空队列阻塞测试通过")
	})

	t.Run("多生产者单消费者测试", func(t *testing.T) {
		q := NewQueue[int]()
		producers := 5
		itemsPerProducer := 100
		totalItems := producers * itemsPerProducer

		var wg sync.WaitGroup

		// 启动消费者
		received := make([]int, 0, totalItems)
		var receiveMu sync.Mutex
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < totalItems; i++ {
				val := q.Pop()
				receiveMu.Lock()
				received = append(received, val)
				receiveMu.Unlock()
			}
		}()

		// 启动多个生产者
		for i := 0; i < producers; i++ {
			wg.Add(1)
			go func(producerID int) {
				defer wg.Done()
				for j := 0; j < itemsPerProducer; j++ {
					q.Push(producerID*1000 + j)
				}
			}(i)
		}

		wg.Wait()

		if len(received) != totalItems {
			t.Errorf("期望接收 %d 个元素，实际接收 %d 个", totalItems, len(received))
		}
		t.Logf("多生产者单消费者测试通过，接收了 %d 个元素", len(received))
	})

	t.Run("单生产者多消费者测试", func(t *testing.T) {
		q := NewQueue[int]()
		consumers := 5
		totalItems := 500

		var wg sync.WaitGroup
		var counter int32

		// 启动多个消费者
		for i := 0; i < consumers; i++ {
			wg.Add(1)
			go func(consumerID int) {
				defer wg.Done()
				for {
					val := q.Pop()
					if val == -1 { // 结束信号
						return
					}
					atomic.AddInt32(&counter, 1)
				}
			}(i)
		}

		// 生产者
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < totalItems; i++ {
				q.Push(i)
			}
			// 发送结束信号给所有消费者
			for i := 0; i < consumers; i++ {
				q.Push(-1)
			}
		}()

		wg.Wait()

		if counter != int32(totalItems) {
			t.Errorf("期望消费 %d 个元素，实际消费 %d 个", totalItems, counter)
		}
		t.Logf("单生产者多消费者测试通过，消费了 %d 个元素", counter)
	})

	t.Run("多生产者多消费者测试", func(t *testing.T) {
		q := NewQueue[int]()
		producers := 10
		consumers := 10
		itemsPerProducer := 100
		totalItems := producers * itemsPerProducer

		var wg sync.WaitGroup
		var producedCounter int32
		var consumedCounter int32

		// 启动多个消费者
		for i := 0; i < consumers; i++ {
			wg.Add(1)
			go func(consumerID int) {
				defer wg.Done()
				for {
					val := q.Pop()
					if val == -1 { // 结束信号
						return
					}
					atomic.AddInt32(&consumedCounter, 1)
				}
			}(i)
		}

		// 启动多个生产者
		for i := 0; i < producers; i++ {
			wg.Add(1)
			go func(producerID int) {
				defer wg.Done()
				for j := 0; j < itemsPerProducer; j++ {
					q.Push(producerID*1000 + j)
					atomic.AddInt32(&producedCounter, 1)
				}
			}(i)
		}

		// 等待所有生产者完成
		time.Sleep(100 * time.Millisecond)

		// 发送结束信号给所有消费者
		for i := 0; i < consumers; i++ {
			q.Push(-1)
		}

		wg.Wait()

		if producedCounter != int32(totalItems) {
			t.Errorf("期望生产 %d 个元素，实际生产 %d 个", totalItems, producedCounter)
		}
		if consumedCounter != int32(totalItems) {
			t.Errorf("期望消费 %d 个元素，实际消费 %d 个", totalItems, consumedCounter)
		}
		t.Logf("多生产者多消费者测试通过：生产 %d 个，消费 %d 个", producedCounter, consumedCounter)
	})

	t.Run("FIFO 顺序测试", func(t *testing.T) {
		q := NewQueue[int]()

		// 先 Push 一批数据
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, val := range expected {
			q.Push(val)
		}

		// 验证 Pop 的顺序
		for i, expectedVal := range expected {
			actualVal := q.Pop()
			if actualVal != expectedVal {
				t.Errorf("索引 %d: 期望 %d，实际 %d", i, expectedVal, actualVal)
			}
		}

		t.Log("FIFO 顺序测试通过")
	})

	t.Run("交替 Push 和 Pop 测试", func(t *testing.T) {
		q := NewQueue[int]()

		for i := 0; i < 100; i++ {
			q.Push(i)
			val := q.Pop()
			if val != i {
				t.Errorf("迭代 %d: 期望 %d，实际 %d", i, i, val)
			}
		}

		t.Log("交替 Push/Pop 测试通过")
	})

	t.Run("泛型支持测试 - string", func(t *testing.T) {
		q := NewQueue[string]()

		q.Push("hello")
		q.Push("world")
		q.Push("go")

		if val := q.Pop(); val != "hello" {
			t.Errorf("期望 'hello'，实际 '%s'", val)
		}
		if val := q.Pop(); val != "world" {
			t.Errorf("期望 'world'，实际 '%s'", val)
		}
		if val := q.Pop(); val != "go" {
			t.Errorf("期望 'go'，实际 '%s'", val)
		}

		t.Log("泛型 string 类型测试通过")
	})

	t.Run("泛型支持测试 - struct", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		q := NewQueue[User]()

		users := []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		}

		for _, user := range users {
			q.Push(user)
		}

		for i, expectedUser := range users {
			actualUser := q.Pop()
			if actualUser.ID != expectedUser.ID || actualUser.Name != expectedUser.Name {
				t.Errorf("索引 %d: 期望 %+v，实际 %+v", i, expectedUser, actualUser)
			}
		}

		t.Log("泛型 struct 类型测试通过")
	})

	t.Run("高并发压力测试", func(t *testing.T) {
		q := NewQueue[int]()
		producers := 20
		consumers := 20
		itemsPerProducer := 1000
		totalItems := producers * itemsPerProducer

		var wg sync.WaitGroup
		var consumedCounter int32

		start := time.Now()

		// 启动消费者
		for i := 0; i < consumers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					val := q.Pop()
					if val == -1 {
						return
					}
					atomic.AddInt32(&consumedCounter, 1)
				}
			}()
		}

		// 启动生产者
		for i := 0; i < producers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < itemsPerProducer; j++ {
					q.Push(id*itemsPerProducer + j)
				}
			}(i)
		}

		// 等待所有生产完成
		time.Sleep(200 * time.Millisecond)

		// 发送结束信号
		for i := 0; i < consumers; i++ {
			q.Push(-1)
		}

		wg.Wait()
		duration := time.Since(start)

		if consumedCounter != int32(totalItems) {
			t.Errorf("期望消费 %d 个元素，实际消费 %d 个", totalItems, consumedCounter)
		}
		t.Logf("高并发压力测试通过：%d 个生产者，%d 个消费者，处理 %d 个元素，耗时 %v",
			producers, consumers, consumedCounter, duration)
	})

	t.Run("多个消费者阻塞唤醒测试", func(t *testing.T) {
		q := NewQueue[int]()
		consumers := 5

		var wg sync.WaitGroup
		results := make([]int, consumers)

		// 启动多个消费者，都会在空队列上阻塞
		for i := 0; i < consumers; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				val := q.Pop()
				results[idx] = val
			}(i)
		}

		// 等待所有消费者进入阻塞状态
		time.Sleep(50 * time.Millisecond)

		// 逐个 Push 数据来唤醒消费者
		for i := 0; i < consumers; i++ {
			q.Push(i + 1)
			time.Sleep(10 * time.Millisecond)
		}

		wg.Wait()

		// 验证所有消费者都收到了数据
		receivedMap := make(map[int]bool)
		for _, val := range results {
			if val < 1 || val > consumers {
				t.Errorf("收到无效值: %d", val)
			}
			receivedMap[val] = true
		}

		if len(receivedMap) != consumers {
			t.Errorf("期望收到 %d 个不同的值，实际收到 %d 个", consumers, len(receivedMap))
		}

		t.Logf("多个消费者阻塞唤醒测试通过，%d 个消费者都被正确唤醒", consumers)
	})

	t.Run("快速 Push/Pop 循环测试", func(t *testing.T) {
		q := NewQueue[int]()
		iterations := 10000

		var wg sync.WaitGroup

		// 生产者
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				q.Push(i)
			}
		}()

		// 消费者
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				val := q.Pop()
				if val < 0 || val >= iterations {
					t.Errorf("收到无效值: %d", val)
				}
			}
		}()

		wg.Wait()
		t.Logf("快速 Push/Pop 循环测试通过，处理了 %d 次操作", iterations)
	})

	t.Run("零值类型测试", func(t *testing.T) {
		q := NewQueue[int]()

		// 测试 Push 零值
		q.Push(0)
		q.Push(1)
		q.Push(0)
		q.Push(2)

		values := []int{q.Pop(), q.Pop(), q.Pop(), q.Pop()}
		expected := []int{0, 1, 0, 2}

		for i, val := range values {
			if val != expected[i] {
				t.Errorf("索引 %d: 期望 %d，实际 %d", i, expected[i], val)
			}
		}

		t.Log("零值类型测试通过")
	})
}
