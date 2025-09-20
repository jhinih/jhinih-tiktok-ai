package manager

import (
	"encoding/json"
	"sync"
	"time"
)

type Conn struct {
	id   string
	ch   chan []byte
	once sync.Once
}

func (c *Conn) Close() {
	c.once.Do(func() {
		close(c.ch)
	})
}
func (c *Conn) Chan() <-chan []byte {
	return c.ch
}

var (
	listeners = make(map[string][]*Conn)
	mu        sync.Mutex
)

// Register 注册一个 uid 的监听通道，返回只读通道和取消函数
func Register(id string) (*Conn, func()) {
	conn := &Conn{
		id: id,
		ch: make(chan []byte, 1),
	}

	mu.Lock()
	listeners[id] = append(listeners[id], conn)
	mu.Unlock()

	// 超时保护
	go func() {
		time.Sleep(30 * time.Second)
		conn.Close()
	}()

	cancel := func() {
		conn.Close()
		mu.Lock()
		// 从列表中移除
		conns := listeners[id]
		for i, c := range conns {
			if c == conn {
				listeners[id] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(listeners[id]) == 0 {
			delete(listeners, id)
		}
		mu.Unlock()
	}

	return conn, cancel
}

// Push 向某个 uid 的所有监听者推送结果，只推送一次就关闭通道
func Push(id string, res interface{}) error {
	mu.Lock()
	conns := listeners[id]
	if len(conns) == 0 {
		mu.Unlock()
		return nil
	}
	// 复制一份，避免解锁后操作
	tmp := make([]*Conn, len(conns))
	copy(tmp, conns)
	delete(listeners, id)
	mu.Unlock()

	b, _ := json.Marshal(res)
	for _, conn := range tmp {
		select {
		case conn.ch <- b:
		default:
		}
		conn.Close()
	}
	return nil
}
