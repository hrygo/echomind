package logger

import (
	"sync"
	"time"
)

// sampler 采样控制器
type sampler struct {
	config SamplingConfig
	counts map[Level]int
	mu     sync.RWMutex
	ticker *time.Ticker
	stopCh chan struct{}
}

// newSampler 创建新的采样器
func newSampler(config SamplingConfig) *sampler {
	s := &sampler{
		config: config,
		counts: make(map[Level]int),
		stopCh: make(chan struct{}),
	}

	// 启动定时重置计数器
	s.ticker = time.NewTicker(time.Second)
	go s.resetLoop()

	return s
}

// shouldLog 判断是否应该记录日志
func (s *sampler) shouldLog(level Level) bool {
	if !s.config.Enabled {
		return true
	}

	// 检查是否在采样级别中
	for _, l := range s.config.Levels {
		if l == level {
			break
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查当前计数
	if s.counts[level] >= s.config.Burst {
		return false
	}

	// 检查速率
	if s.counts[level] >= s.config.Rate {
		return false
	}

	s.counts[level]++
	return true
}

// resetLoop 定时重置计数器
func (s *sampler) resetLoop() {
	for {
		select {
		case <-s.ticker.C:
			s.mu.Lock()
			for level := range s.counts {
				s.counts[level] = 0
			}
			s.mu.Unlock()
		case <-s.stopCh:
			return
		}
	}
}
