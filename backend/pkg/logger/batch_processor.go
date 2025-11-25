package logger

import (
	"context"
	"sync"
	"time"
)

// BatchProcessor 批量处理器
type BatchProcessor struct {
	provider     Provider
	buffer       []*LogEntry
	bufferSize   int
	flushTimeout time.Duration
	mu           sync.Mutex
	flushCh      chan struct{}
	stopCh       chan struct{}
	wg           sync.WaitGroup
	onFlush      func([]*LogEntry) error
}

// NewBatchProcessor 创建新的批量处理器
func NewBatchProcessor(provider Provider, bufferSize int, flushTimeout time.Duration) *BatchProcessor {
	bp := &BatchProcessor{
		provider:     provider,
		buffer:       make([]*LogEntry, 0, bufferSize),
		bufferSize:   bufferSize,
		flushTimeout: flushTimeout,
		flushCh:      make(chan struct{}, 1),
		stopCh:       make(chan struct{}),
	}

	// 启动后台 goroutine 处理批量写入
	bp.wg.Add(1)
	go bp.flushLoop()

	return bp
}

// Write 写入日志条目到缓冲区
func (bp *BatchProcessor) Write(ctx context.Context, entry *LogEntry) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.buffer = append(bp.buffer, entry)

	// 如果缓冲区满了，触发立即刷新
	if len(bp.buffer) >= bp.bufferSize {
		select {
		case bp.flushCh <- struct{}{}:
		default:
			// 已经有刷新请求在队列中
		}
	}

	return nil
}

// Close 关闭批量处理器
func (bp *BatchProcessor) Close() error {
	close(bp.stopCh)
	bp.wg.Wait()

	// 刷新剩余的缓冲区
	return bp.flush()
}

// SetFlushCallback 设置刷新回调
func (bp *BatchProcessor) SetFlushCallback(onFlush func([]*LogEntry) error) {
	bp.onFlush = onFlush
}

// flushLoop 后台刷新循环
func (bp *BatchProcessor) flushLoop() {
	defer bp.wg.Done()

	ticker := time.NewTicker(bp.flushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-bp.flushCh:
			bp.flush()
		case <-ticker.C:
			bp.flush()
		case <-bp.stopCh:
			// 最后一次刷新
			bp.flush()
			return
		}
	}
}

// flush 执行实际的刷新操作
func (bp *BatchProcessor) flush() error {
	bp.mu.Lock()
	if len(bp.buffer) == 0 {
		bp.mu.Unlock()
		return nil
	}

	// 创建缓冲区的副本
	entries := make([]*LogEntry, len(bp.buffer))
	copy(entries, bp.buffer)
	bp.buffer = bp.buffer[:0] // 清空缓冲区
	bp.mu.Unlock()

	// 执行回调（如果有的话）
	if bp.onFlush != nil {
		if err := bp.onFlush(entries); err != nil {
			// 记录错误但不阻止后续处理
			// 在实际应用中，可以使用错误队列重试
		}
	}

	// 写入到提供者
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, entry := range entries {
		if err := bp.provider.Write(ctx, entry); err != nil {
			// 记录错误，继续处理其他条目
			// 可以添加重试逻辑或错误队列
		}
	}

	return nil
}

// GetBufferSize 获取当前缓冲区大小
func (bp *BatchProcessor) GetBufferSize() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.buffer)
}

// ForceFlush 强制立即刷新
func (bp *BatchProcessor) ForceFlush() error {
	select {
	case bp.flushCh <- struct{}{}:
	default:
		// 已经有刷新请求在队列中
	}

	// 等待刷新完成
	time.Sleep(10 * time.Millisecond)
	return nil
}

// AsyncBatchProcessor 异步批量处理器
type AsyncBatchProcessor struct {
	*Batches
	workQueue chan *BatchJob
	workers   int
	wg        sync.WaitGroup
	stopCh    chan struct{}
}

// BatchJob 批量作业
type BatchJob struct {
	Entries  []*LogEntry
	Callback func(error)
}

// NewAsyncBatchProcessor 创建异步批量处理器
func NewAsyncBatchProcessor(provider Provider, bufferSize int, flushTimeout time.Duration, workers int) *AsyncBatchProcessor {
	batches := NewBatches(provider, bufferSize, flushTimeout)

	abp := &AsyncBatchProcessor{
		Batches:   batches,
		workQueue: make(chan *BatchJob, 1000), // 缓冲队列
		workers:   workers,
		stopCh:    make(chan struct{}),
	}

	// 启动工作协程
	for i := 0; i < workers; i++ {
		abp.wg.Add(1)
		go abp.worker()
	}

	return abp
}

// Write 异步写入
func (abp *AsyncBatchProcessor) Write(ctx context.Context, entry *LogEntry) error {
	job := &BatchJob{
		Entries: []*LogEntry{entry},
		Callback: func(err error) {
			if err != nil {
				// 处理错误，可以记录到备用日志或重试队列
			}
		},
	}

	select {
	case abp.workQueue <- job:
		return nil
	default:
		// 队列满了，直接写入
		return abp.Batches.Write(ctx, entry)
	}
}

// Close 关闭异步批量处理器
func (abp *AsyncBatchProcessor) Close() error {
	close(abp.stopCh)
	abp.wg.Wait()
	return abp.Batches.Close()
}

// worker 工作协程
func (abp *AsyncBatchProcessor) worker() {
	defer abp.wg.Done()

	for {
		select {
		case job := <-abp.workQueue:
			// 处理作业
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			var err error
			for _, entry := range job.Entries {
				if e := abp.Batches.provider.Write(ctx, entry); e != nil {
					err = e
				}
			}
			cancel()

			if job.Callback != nil {
				job.Callback(err)
			}

		case <-abp.stopCh:
			// 处理剩余的作业
			for len(abp.workQueue) > 0 {
				job := <-abp.workQueue
				if job.Callback != nil {
					job.Callback(context.Canceled)
				}
			}
			return
		}
	}
}

// Batches 多批量管理器
type Batches struct {
	processors []*BatchProcessor
	provider   Provider
	mu         sync.RWMutex
}

// NewBatches 创建多批量管理器
func NewBatches(provider Provider, bufferSize int, flushTimeout time.Duration) *Batches {
	return &Batches{
		processors: make([]*BatchProcessor, 0),
		provider:   provider,
	}
}

// AddProcessor 添加批量处理器
func (b *Batches) AddProcessor(processor *BatchProcessor) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.processors = append(b.processors, processor)
}

// Write 写入到所有处理器
func (b *Batches) Write(ctx context.Context, entry *LogEntry) error {
	b.mu.RLock()
	processors := make([]*BatchProcessor, len(b.processors))
	copy(processors, b.processors)
	b.mu.RUnlock()

	// 写入到所有处理器，如果有任何失败，返回第一个错误
	var firstErr error
	for _, processor := range processors {
		if err := processor.Write(ctx, entry); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// Close 关闭所有处理器
func (b *Batches) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var lastErr error
	for _, processor := range b.processors {
		if err := processor.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// GetStats 获取统计信息
func (b *Batches) GetStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := map[string]interface{}{
		"processor_count": len(b.processors),
		"processors":      make([]map[string]interface{}, len(b.processors)),
	}

	for i, processor := range b.processors {
		stats["processors"].([]map[string]interface{})[i] = map[string]interface{}{
			"buffer_size": processor.GetBufferSize(),
		}
	}

	return stats
}

// SmartBatchProcessor 智能批量处理器，根据日志量自动调整
type SmartBatchProcessor struct {
	*Batches
	currentWorkers int
	minWorkers     int
	maxWorkers     int
	adjustInterval time.Duration
	lastAdjust     time.Time
	totalProcessed int64
	mu             sync.Mutex
}

// NewSmartBatchProcessor 创建智能批量处理器
func NewSmartBatchProcessor(provider Provider, minWorkers, maxWorkers int) *SmartBatchProcessor {
	sb := &SmartBatchProcessor{
		Batches:        NewBatches(provider, 100, 1*time.Second),
		minWorkers:     minWorkers,
		maxWorkers:     maxWorkers,
		currentWorkers: minWorkers,
		adjustInterval: 30 * time.Second,
		lastAdjust:     time.Now(),
	}

	// 初始化工作协程
	sb.adjustWorkers()

	// 启动自动调整协程
	go sb.autoAdjust()

	return sb
}

// Write 写入并自动调整
func (sb *SmartBatchProcessor) Write(ctx context.Context, entry *LogEntry) error {
	sb.mu.Lock()
	sb.totalProcessed++
	sb.mu.Unlock()

	return sb.Batches.Write(ctx, entry)
}

// autoAdjust 自动调整工作协程数量
func (sb *SmartBatchProcessor) autoAdjust() {
	ticker := time.NewTicker(sb.adjustInterval)
	defer ticker.Stop()

	for range ticker.C {
		sb.adjustWorkers()
	}
}

// adjustWorkers 调整工作协程数量
func (sb *SmartBatchProcessor) adjustWorkers() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	// 简单的启发式算法：根据处理速度调整
	processedSinceLastAdjust := sb.totalProcessed
	timeSinceLastAdjust := time.Since(sb.lastAdjust)

	// 计算每秒处理速度
	processingRate := float64(processedSinceLastAdjust) / timeSinceLastAdjust.Seconds()

	// 根据处理率调整协程数量
	var newWorkers int
	switch {
	case processingRate > 1000: // 高负载
		newWorkers = sb.maxWorkers
	case processingRate > 500: // 中等负载
		newWorkers = sb.minWorkers + (sb.maxWorkers-sb.minWorkers)/2
	case processingRate < 100: // 低负载
		newWorkers = sb.minWorkers
	default:
		newWorkers = sb.currentWorkers
	}

	// 如果需要调整
	if newWorkers != sb.currentWorkers {
		// 这里可以实现动态协程调整
		// 为简化，只是记录决策
		sb.currentWorkers = newWorkers
		sb.lastAdjust = time.Now()
		sb.totalProcessed = 0
	}
}

// GetMetrics 获取性能指标
func (sb *SmartBatchProcessor) GetMetrics() map[string]interface{} {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	metrics := sb.Batches.GetStats()
	metrics["current_workers"] = sb.currentWorkers
	metrics["min_workers"] = sb.minWorkers
	metrics["max_workers"] = sb.maxWorkers
	metrics["total_processed"] = sb.totalProcessed

	return metrics
}
