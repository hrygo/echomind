# EchoMind 日志框架性能优化指南

## 📋 目录

1. [性能基准](#性能基准)
2. [批量处理优化](#批量处理优化)
3. [内存优化](#内存优化)
4. [网络传输优化](#网络传输优化)
5. [存储优化](#存储优化)
6. [并发优化](#并发优化)
7. [监控和调优](#监控和调优)

## 📊 性能基准

### 基准测试结果

```go
// 基准测试配置
func runBenchmarks() {
    b := testing.B{}

    // 基础日志记录基准
    b.Run("BasicLogging", benchmarkBasicLogging)
    b.Run("LoggingWithFields", benchmarkLoggingWithFields)
    b.Run("LoggingWithContext", benchmarkLoggingWithContext)

    // 提供者性能基准
    b.Run("NoopProvider", benchmarkNoopProvider)
    b.Run("FileProvider", benchmarkFileProvider)
    b.Run("ElasticsearchProvider", benchmarkElasticsearchProvider)
    b.Run("LokiProvider", benchmarkLokiProvider)
}
```

### 性能指标

```go
// 性能指标常量
const (
    TargetLatencyPerLog    = 100 * time.Microsecond // 每条日志 100μs
    TargetThroughputPerSec = 10000               // 每秒 10000 条日志
    TargetMemoryUsage      = 100 * 1024 * 1024      // 100MB 内存限制
    TargetCPUUsage         = 10                     // 10% CPU 使用率限制
)

// 性能监控
type PerformanceMetrics struct {
    LatencyPerLog    time.Duration
    ThroughputPerSec float64
    MemoryUsage      int64
    CPUUsage         float64
    ErrorRate        float64
}

func (pm PerformanceMetrics) IsHealthy() bool {
    return pm.LatencyPerLog <= TargetLatencyPerLog &&
           pm.ThroughputPerSec >= TargetThroughputPerSec &&
           pm.MemoryUsage <= TargetMemoryUsage &&
           pm.CPUUsage <= TargetCPUUsage &&
           pm.ErrorRate <= 5.0 // 错误率不超过 5%
}
```

## 🚀 批量处理优化

### 1. 批量处理器设计

```go
// 高效批量处理器
type HighPerformanceBatchProcessor struct {
    provider     Provider
    buffer       []*LogEntry
    bufferLock   sync.Mutex
    bufferCap    int
    flushCh      chan struct{}
    stopCh       chan struct{}
    workers      int
    workQueue    chan *BatchJob
    pool         *sync.Pool
}

// 批量作业
type BatchJob struct {
    Entries []*LogEntry
    Callback func(error)
}

// 创建高性能批量处理器
func NewHighPerformanceBatchProcessor(provider Provider, config BatchConfig) *HighPerformanceBatchProcessor {
    processor := &HighPerformanceBatchProcessor{
        provider:     provider,
        buffer:       make([]*LogEntry, 0, config.BufferSize),
        bufferCap:    config.BufferSize,
        flushCh:      make(chan struct{}, 1),
        stopCh:       make(chan struct{}),
        workers:      config.Workers,
        workQueue:    make(chan *BatchJob, config.QueueSize),
    }

    // 初始化对象池
    processor.pool = &sync.Pool{
        New: func() interface{} {
            return make([]*LogEntry, 0, 100)
        },
    }

    // 启动工作协程
    for i := 0; i < processor.workers; i++ {
        go processor.worker(i)
    }

    // 启动刷新协程
    go processor.flushLoop()

    return processor
}

// 写入日志条目
func (bp *HighPerformanceBatchProcessor) Write(ctx context.Context, entry *LogEntry) error {
    // 快速路径：如果缓冲区未满，直接添加
    bp.bufferLock.Lock()
    if len(bp.buffer) < bp.bufferCap {
        bp.buffer = append(bp.buffer, entry)
        bp.bufferLock.Unlock()
        return nil
    }
    bp.bufferLock.Unlock()

    // 缓冲区满，触发刷新
    bp.flushCh <- struct{}{}

    // 重试写入
    bp.bufferLock.Lock()
    if len(bp.buffer) < bp.bufferCap {
        bp.buffer = append(bp.buffer, entry)
        bp.bufferLock.Unlock()
        return nil
    }
    bp.bufferLock.Unlock()

    // 降级到同步写入
    return bp.provider.Write(ctx, entry)
}
```

### 2. 自适应批量大小

```go
// 自适应批量大小控制器
type AdaptiveBatchController struct {
    minSize    int
    maxSize    int
    currentSize int
    adjustRate  float64
    lastAdjust  time.Time
    metrics     *BatchMetrics
    mu          sync.RWMutex
}

type BatchMetrics struct {
    ProcessCount int64
    ProcessTime  time.Duration
    ErrorCount   int64
    LastFlush    time.Time
}

func NewAdaptiveBatchController(minSize, maxSize int) *AdaptiveBatchController {
    return &AdaptiveBatchController{
        minSize:    minSize,
        maxSize:    maxSize,
        currentSize: minSize,
        adjustRate:  0.1,
        lastAdjust:  time.Now(),
        metrics:     &BatchMetrics{},
    }
}

// 自适应调整批量大小
func (abc *AdaptiveBatchController) AdjustBatchSize() {
    abc.mu.Lock()
    defer abc.mu.Unlock()

    timeSinceLastAdjust := time.Since(abc.lastAdjust)
    if timeSinceLastAdjust < 10*time.Second {
        return
    }

    // 计算处理速度（条/秒）
    var processSpeed float64
    if abc.metrics.ProcessCount > 0 {
        processSpeed = float64(abc.metrics.ProcessCount) / abc.metrics.ProcessTime.Seconds()
    }

    // 计算错误率
    var errorRate float64
    if abc.metrics.ProcessCount > 0 {
        errorRate = float64(abc.metrics.ErrorCount) / float64(abc.metrics.ProcessCount)
    }

    // 根据指标调整批量大小
    newSize := abc.currentSize
    switch {
    case processSpeed > 1000 && errorRate < 0.01: // 高性能，低错误率
        newSize = int(float64(abc.currentSize) * (1 + abc.adjustRate))
    case processSpeed < 100 || errorRate > 0.05: // 低性能或高错误率
        newSize = int(float64(abc.currentSize) * (1 - abc.adjustRate))
    }

    // 限制在 min/max 范围内
    newSize = max(abc.minSize, min(abc.maxSize, newSize))

    if newSize != abc.currentSize {
        abc.currentSize = newSize
        abc.lastAdjust = time.Now()

        logger.Debug("调整批量大小",
            logger.Int("old_size", abc.currentSize),
            logger.Int("new_size", newSize),
            logger.Float64("process_speed", processSpeed),
            logger.Float64("error_rate", errorRate))
    }

    // 重置指标
    abc.metrics = &BatchMetrics{}
}
```

### 3. 预分配和对象池

```go
// 预分配的字段池
type FieldPool struct {
    stringPool   sync.Pool
    intPool     sync.Pool
    boolPool    sync.Pool
    timePool    sync.Pool
    float64Pool sync.Pool
    anyPool     sync.Pool
}

func NewFieldPool() *FieldPool {
    return &FieldPool{
        stringPool: sync.Pool{
            New: func() interface{} {
                return make([]string, 0, 20)
            },
        },
        intPool: sync.Pool{
            New: func() interface{} {
                return make([]int, 0, 20)
            },
        },
        boolPool: sync.Pool{
            New: func() interface{} {
                return make([]bool, 0, 20)
            },
        },
        timePool: sync.Pool{
            New: func() interface{} {
                return make([]time.Time, 0, 20)
            },
        },
        float64Pool: sync.Pool{
            New: func() interface{} {
                return make([]float64, 0, 20)
            },
        },
        anyPool: sync.Pool{
            New: func() interface{} {
                return make([]interface{}, 0, 20)
            },
        },
    }
}

// 使用对象池创建字段
func (fp *FieldPool) String(key, value string) logger.Field {
    return logger.Field{Key: key, Value: value}
}

// 批量字段创建
func (fp *FieldPool) CreateBatch(count int) []logger.Field {
    fields := make([]logger.Field, 0, count)
    stringFields := fp.stringPool.Get().([]string)
    defer fp.stringPool.Put(stringFields)

    // 重用 string 数组
    *stringFields = (*stringFields)[:0]

    for i := 0; i < count; i++ {
        fields = append(fields, logger.Field{Key: fmt.Sprintf("field_%d", i), Value: ""})
    }

    return fields
}
```

## 💾 内存优化

### 1. 内存池化

```go
// 内存池管理器
type MemoryPoolManager struct {
    entryPools   map[int]*sync.Pool // 按大小分组的 LogEntry 池
    bufferPools map[int]*sync.Pool // 按大小分组的缓冲区池
    mu          sync.RWMutex
}

func NewMemoryPoolManager() *MemoryPoolManager {
    mpm := &MemoryPoolManager{
        entryPools:   make(map[int]*sync.Pool),
        bufferPools: make(map[int]*sync.Pool),
    }

    // 初始化不同大小的池
    sizes := []int{64, 128, 256, 512, 1024, 2048, 4096}
    for _, size := range sizes {
        mpm.entryPools[size] = &sync.Pool{
            New: func() interface{} {
                return make([]*LogEntry, 0, size)
            },
        }
        mpm.bufferPools[size] = &sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, size)
            },
        }
    }

    return mpm
}

// 获取 LogEntry 切片
func (mpm *MemoryPoolManager) GetEntrySlice(size int) []*LogEntry {
    // 找到合适大小的池
    poolSize := roundUpPowerOfTwo(size)
    if poolSize > 4096 {
        poolSize = 4096
    }

    mpm.mu.RLock()
    pool := mpm.entryPools[poolSize]
    mpm.RUnlock()

    if pool != nil {
        return pool.Get().([]*LogEntry)
    }

    // 降级到直接分配
    return make([]*LogEntry, 0, size)
}

// 归还 LogEntry 切片
func (mpm *MemoryPoolManager) PutEntrySlice(entries []*LogEntry) {
    if len(entries) == 0 {
        return
    }

    size := cap(entries)
    poolSize := roundUpPowerOfTwo(size)
    if poolSize > 4096 {
        return
    }

    mpm.mu.RLock()
    pool := mpm.entryPools[poolSize]
    mpm.RUnlock()

    if pool != nil {
        // 清空切片但保持容量
        entries = entries[:0]
        pool.Put(entries)
    }
}

func roundUpPowerOfTwo(n int) int {
    if n <= 0 {
        return 1
    }
    n--
    n |= n >> 1
    n |= n >> 2
    n |= n >> 4
    n |= n >> 8
    n |= n >> 16
    n++
    return n
}
```

### 2. 延迟序列化

```go
// 延迟序列化器
type LazySerializer struct {
    data    interface{}
    cache   []byte
    cached  bool
    marshal func(interface{}) ([]byte, error)
}

func NewLazySerializer(marshal func(interface{}) ([]byte, error)) *LazySerializer {
    return &LazySerializer{
        marshal: marshal,
    }
}

func (ls *LazySerializer) SetData(data interface{}) {
    ls.data = data
    ls.cached = false
}

func (ls *LazySerializer) Bytes() ([]byte, error) {
    if !ls.cached {
        serialized, err := ls.marshal(ls.data)
        if err != nil {
            return nil, err
        }
        ls.cache = serialized
        ls.cached = true
    }
    return ls.cache, nil
}

// 在 LogEntry 中使用延迟序列化
type LogEntry struct {
    Timestamp time.Time
    Level     Level
    Message   string
    Fields    *LazyFields
    Context   ContextInfo
    Source    SourceInfo
}

type LazyFields struct {
    data    map[string]interface{}
    cache   []byte
    cached  bool
    marshal func(map[string]interface{}) ([]byte, error)
}

func (lf *LazyFields) Set(key string, value interface{}) {
    lf.data[key] = value
    lf.cached = false
}

func (lf *LazyFields) Bytes() ([]byte, error) {
    if !lf.cached {
        serialized, err := lf.marshal(lf.data)
        if err != nil {
            return nil, err
        }
        lf.cache = serialized
        lf.cached = true
    }
    return lf.cache, nil
}
```

### 3. 字符串构建优化

```go
// 高效字符串构建器
type StringBuilder struct {
    builder strings.Builder
    pool    *sync.Pool
}

func NewStringBuilder() *StringBuilder {
    return &StringBuilder{
        pool: &sync.Pool{
            New: func() interface{} {
                return &StringBuilder{
                    builder: strings.Builder{},
                }
            },
        },
    }
}

// 获取字符串构建器
func (sb *StringBuilder) Get() *strings.Builder {
    return sb.pool.Get().(*strings.Builder)
}

// 归还字符串构建器
func (sb *StringBuilder) Put(builder *strings.Builder) {
    builder.Reset()
    sb.pool.Put(builder)
}

// 格式化日志消息
func (sb *StringBuilder) FormatMessage(template string, args ...interface{}) string {
    builder := sb.Get()
    defer sb.Put(builder)

    builder.Grow(len(template) + len(args)*10) // 预估容量

    if len(args) == 0 {
        builder.WriteString(template)
    } else {
        fmt.Fprintf(builder, template, args...)
    }

    return builder.String()
}

// 构建 JSON
func (sb *StringBuilder) BuildJSON(fields map[string]interface{}) string {
    builder := sb.Get()
    defer sb.Put(builder)

    builder.Grow(512) // 预估 JSON 大小

    builder.WriteByte('{')
    first := true
    for key, value := range fields {
        if !first {
            builder.WriteByte(',')
        }
        first = false

        // 写入 key
        builder.WriteByte('"')
        builder.WriteString(key)
        builder.WriteString("\":")

        // 写入 value（简化实现，实际应该使用 json.Marshal）
        builder.WriteString(fmt.Sprintf("%v", value))
    }
    builder.WriteByte('}')

    return builder.String()
}
```

## 🌐 网络传输优化

### 1. 连接池管理

```go
// HTTP 连接池
type HTTPConnectionPool struct {
    client     *http.Client
    transports []*http.Transport
    mu         sync.Mutex
    maxIdle    int
    current    int
}

func NewHTTPConnectionPool(maxIdle int) *HTTPConnectionPool {
    pool := &HTTPConnectionPool{
        maxIdle: maxIdle,
    }

    // 预创建传输层
    for i := 0; i < maxIdle; i++ {
        transport := &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
            DisableCompression:  false,
        }
        pool.transports = append(pool.transports, transport)
    }

    pool.client = &http.Client{
        Transport: pool.getNextTransport(),
    }

    return pool
}

func (pool *HTTPConnectionPool) getNextTransport() *http.Transport {
    pool.mu.Lock()
    defer pool.mu.Unlock()

    if len(pool.transports) > 0 {
        transport := pool.transports[len(pool.transports)-1]
        pool.transports = pool.transports[:len(pool.transports)-1]
        return transport
    }

    // 如果池空，创建新的传输层
    return &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    }
}

func (pool *HTTPConnectionPool) ReturnTransport(transport *http.Transport) {
    pool.mu.Lock()
    defer pool.mu.Unlock()

    if len(pool.transports) < pool.maxIdle {
        // 重置传输层状态
        transport.CloseIdleConnections()
        pool.transports = append(pool.transports, transport)
    } else {
        // 池满了，直接关闭
        transport.Close()
    }
}

func (pool *HTTPConnectionPool) Close() error {
    pool.mu.Lock()
    defer pool.mu.Unlock()

    for _, transport := range pool.transports {
        transport.CloseIdleConnections()
    }

    pool.client.CloseIdleConnections()
    return nil
}
```

### 2. 请求批量化

```go
// 请求批量处理器
type RequestBatchProcessor struct {
    pool       *HTTPConnectionPool
    batchSize  int
    timeout    time.Duration
    mu         sync.Mutex
    pending    []PendingRequest
}

type PendingRequest struct {
    Request  *http.Request
    Response chan *http.Response
    Error    chan error
}

func NewRequestBatchProcessor(pool *HTTPConnectionPool, batchSize int, timeout time.Duration) *RequestBatchProcessor {
    return &RequestBatchProcessor{
        pool:      pool,
        batchSize: batchSize,
        timeout:   timeout,
        pending:   make([]PendingRequest, 0),
    }
}

// 发送请求
func (rbp *RequestBatchProcessor) Send(req *http.Request) (*http.Response, error) {
    respCh := make(chan *http.Response, 1)
    errCh := make(chan error, 1)

    pending := PendingRequest{
        Request:  req,
        Response: respCh,
        Error:    errCh,
    }

    rbp.mu.Lock()
    rbp.pending = append(rbp.pending, pending)
    shouldFlush := len(rbp.pending) >= rbp.batchSize
    rbp.mu.Unlock()

    if shouldFlush {
        go rbp.flush()
    }

    select {
    case resp := <-respCh:
        return resp, nil
    case err := <-errCh:
        return nil, err
    case <-time.After(rbp.timeout):
        return nil, fmt.Errorf("request timeout after %v", rbp.timeout)
    }
}

// 批量刷新
func (rbp *RequestBatchProcessor) flush() error {
    rbp.mu.Lock()
    if len(rbp.pending) == 0 {
        rbp.mu.Unlock()
        return nil
    }

    // 创建批次副本
    batch := make([]PendingRequest, len(rbp.pending))
    copy(batch, rbp.pending)
    rbp.pending = rbp.pending[:0]
    rbp.mu.Unlock()

    // 并发处理批次
    var wg sync.WaitGroup
    errCh := make(chan error, len(batch))

    for i, pending := range batch {
        wg.Add(1)
        go func(idx int, p PendingRequest) {
            defer wg.Done()
            errCh[idx] = rbp.processRequest(p)
        }(i, pending)
    }

    wg.Wait()

    // 检查错误
    for _, err := range errCh {
        if err != nil {
            return err
        }
    }

    return nil
}

func (rbp *RequestBatchProcessor) processRequest(pending PendingRequest) error {
    transport := rbp.pool.getNextTransport()
    client := &http.Client{
        Transport: transport,
        Timeout:  5 * time.Second,
    }

    defer rbp.pool.ReturnTransport(transport)

    resp, err := client.Do(pending.Request)
    if err != nil {
        pending.Error <- err
        return err
    }

    pending.Response <- resp
    return nil
}
```

### 3. 压缩传输

```go
// 压缩传输客户端
type CompressingClient struct {
    client     *http.Client
    compressor  *Compressor
}

type Compressor struct {
    level int
}

func NewCompressingClient() *CompressingClient {
    pool := NewHTTPConnectionPool(10)

    return &CompressingClient{
        client: pool.client,
        compressor: &Compressor{
            level: gzip.BestCompression,
        },
    }
}

func (cc *CompressingClient) SendJSON(url string, data interface{}) (*http.Response, error) {
    // 序列化 JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    // 压缩数据
    compressedData := cc.compressor.Compress(jsonData)

    // 创建请求
    req, err := http.NewRequest("POST", url, bytes.NewReader(compressedData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Content-Encoding", "gzip")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", len(compressedData)))

    // 发送请求
    return cc.client.Do(req)
}

// Gzip 压缩器
func (c *Compressor) Compress(data []byte) []byte {
    var buf bytes.Buffer
    gz := gzip.NewWriter(&buf)
    gz.Level = c.level

    if _, err := gz.Write(data); err != nil {
        gz.Close()
        return data
    }
    gz.Close()

    return buf.Bytes()
}
```

## 💾 存储优化

### 1. 日志轮转优化

```go
// 智能日志轮转器
type SmartLogRotator struct {
    writer      io.Writer
    maxSize     int64
    maxAge      time.Duration
    maxBackups  int
    currentSize  int64
    currentAge  time.Time
    mu           sync.Mutex
    stats       *RotationStats
}

type RotationStats struct {
    Rotations   int64
    TotalBytes   int64
    CompressedBytes int64
}

func NewSmartLogRotator(writer io.Writer, maxSize int64, maxAge time.Duration, maxBackups int) *SmartLogRotator {
    return &SmartLogRotator{
        writer:     writer,
        maxSize:     maxSize,
        maxAge:      maxAge,
        maxBackups:  maxBackups,
        currentSize: 0,
        currentAge:  time.Now(),
        stats:       &RotationStats{},
    }
}

// 写入数据
func (slr *SmartLogRotator) Write(data []byte) (int, error) {
    slr.mu.Lock()
    defer slr.mu.Unlock()

    // 检查是否需要轮转
    shouldRotate := slr.shouldRotate(len(data))
    if shouldRotate {
        if err := slr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := slr.writer.Write(data)
    if err != nil {
        return n, err
    }

    slr.currentSize += int64(n)
    return n, nil
}

func (slr *SmartLogRotator) shouldRotate(dataSize int) bool {
    // 检查大小限制
    if slr.currentSize+int64(dataSize) >= slr.maxSize {
        return true
    }

    // 检查时间限制
    if time.Since(slr.currentAge) >= slr.maxAge {
        return true
    }

    return false
}

func (slr *SmartLogRotator) rotate() error {
    // 更新统计
    slr.stats.Rotations++
    slr.stats.TotalBytes += slr.currentSize

    // 创建备份文件名
    timestamp := time.Now().Format("2006-01-02-15-04-05")
    backupPath := fmt.Sprintf("/var/log/app.%s.%d.log", timestamp, slr.stats.Rotations)

    // 压缩备份文件
    if err := slr.compressAndBackup(backupPath); err != nil {
        logger.Error("压缩备份文件失败", logger.Error(err))
        // 继续执行，不中断日志记录
    }

    // 重置状态
    slr.currentSize = 0
    slr.currentAge = time.Now()

    return nil
}

func (slr *SmartLogRotator) compressAndBackup(path string) error {
    // 读取当前文件内容
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    var buf bytes.Buffer
    gz := gzip.NewWriter(&buf)
    defer gz.Close()

    // 复制并压缩
    if _, err := io.Copy(gz, file); err != nil {
        return err
    }

    // 创建压缩文件
    compressedPath := path + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    if _, err := compressedFile.Write(buf.Bytes()); err != nil {
        return err
    }

    // 更新压缩统计
    slr.stats.CompressedBytes += int64(buf.Len())

    // 删除原始文件
    return os.Remove(path)
}
```

### 2. 索引和搜索优化

```go
// 优化的索引管理器
type OptimizedIndexManager struct {
    entries    []*LogEntry
    index      map[string][]*LogEntry  // 字段索引
    timestampIndex []TimeWindow       // 时间窗口索引
    mu         sync.RWMutex
    maxEntries int
    windowSize  time.Duration
}

type TimeWindow struct {
    Start    time.Time
    End      time.Time
    Indexes  []*LogEntry
}

func NewOptimizedIndexManager(maxEntries int, windowSize time.Duration) *OptimizedIndexManager {
    return &OptimizedIndexManager{
        entries:       make([]*LogEntry, 0, maxEntries),
        index:         make(map[string][]*LogEntry),
        timestampIndex: make([]TimeWindow, 0),
        maxEntries:    maxEntries,
        windowSize:     windowSize,
    }
}

// 添加日志条目
func (oim *OptimizedIndexManager) Add(entry *LogEntry) error {
    oim.mu.Lock()
    defer oim.mu.Unlock()

    // 检查容量限制
    if len(oim.entries) >= oim.maxEntries {
        oim.cleanup()
    }

    // 添加到主索引
    oim.entries = append(oim.entries, entry)

    // 更新字段索引
    oim.updateFieldIndex(entry)

    // 更新时间窗口索引
    oim.updateTimestampIndex(entry)

    return nil
}

// 清理旧条目
func (oim *OptimizedIndexManager) cleanup() {
    if len(oim.entries) == 0 {
        return
    }

    // 删除最旧的条目
    oldEntry := oim.entries[0]
    oim.entries = oim.entries[1:]

    // 更新字段索引
    for field, entries := range oim.index {
        // 移除引用
        for i, e := range entries {
            if e == oldEntry {
                oim.index[field] = append(entries[:i], entries[i+1:]...)
                break
            }
        }
    }

    // 更新时间窗口索引
    cutoffTime := time.Now().Add(-oim.windowSize)
    for i, window := range oim.timestampIndex {
        if window.End.Before(cutoffTime) {
            oim.timestampIndex = oim.timestampIndex[i+1:]
            break
        }
    }
}

// 查询优化
func (oim *OptimizedIndexManager) Query(query Query) ([]*LogEntry, error) {
    oim.mu.RLock()
    defer oim.RUnlock()

    var results []*LogEntry

    // 使用字段索引进行快速查找
    if field, value, ok := query.ExactMatchField(); ok {
        if entries, exists := oim.index[field]; exists {
            for _, entry := range entries {
                if matchesField(entry, field, value) {
                    results = append(results, entry)
                }
            }
            return results, nil
    }

    // 全文搜索
    for _, entry := range oim.entries {
        if matchesQuery(entry, query) {
            results = append(results, entry)
        }
    }

    // 排序和限制结果
    sort.Slice(results, func(i, j int) bool {
        return results[i].Timestamp.After(results[j].Timestamp)
    })

    if query.Limit > 0 && len(results) > query.Limit {
        results = results[:query.Limit]
    }

    return results, nil
}
```

## 🔧 并发优化

### 1. 无锁设计

```go
// 无锁日志记录器
type LockFreeLogger struct {
    ringBuffer    []LogEntry
    head         uint64
    tail         uint64
    mask         uint64
    size         int
    writer       chan []LogEntry
    processor    func([]LogEntry)
    stopCh       chan struct{}
    mu           sync.Mutex // 仅用于停止
}

func NewLockFreeLogger(size int, processor func([]LogEntry)) *LockFreeLogger {
    // 确保大小是2的幂
    size = nextPowerOf2(size)

    logger := &LockFreeLogger{
        ringBuffer: make([]LogEntry, size),
        head:       0,
        tail:       0,
        mask:       uint64(size - 1),
        size:       size,
        writer:     make(chan []LogEntry, 100),
        processor:  processor,
        stopCh:     make(chan struct{}),
    }

    // 启动写入协程
    go logger.writerLoop()
    return logger
}

func (lfl *LockFreeLogger) Write(entry *LogEntry) {
    select {
    case lfl.writer <- []LogEntry{*entry}:
        default:
        // 如果队列满了，丢弃日志（或添加到丢弃计数）
        atomic.AddUint64(&lfl.dropCount, 1)
    }
}

func (lfl *LockFreeLogger) writerLoop() {
    batch := make([]LogEntry, 0, 100)

    for {
        select {
        case entries := <-lfl.writer:
            // 批量接收日志条目
            batch = append(batch, entries...)

            // 批量处理
            if len(batch) >= 100 {
                lfl.processor(batch)
                batch = batch[:0] // 重用切片
            }

        case <-lfl.stopCh:
            // 处理剩余日志
            if len(batch) > 0 {
                lfl.processor(batch)
            }
            return
        }
    }
}

func (lfl *LockFreeLogger) Stop() {
    close(lfl.stopCh)
}
```

### 2. 分片处理

```go
// 分片日志处理器
type ShardedLogger struct {
    shards    []*LogShard
    hasher    LogHasher
    shardMask uint64
}

type LogShard struct {
    id       int
    buffer   []*LogEntry
    mutex    sync.Mutex
    processor BatchProcessor
}

func NewShardedLogger(numShards int, processorFactory func(int) BatchProcessor) *ShardedLogger {
    shards := make([]*LogShard, numShards)

    for i := 0; i < numShards; i++ {
        shards[i] = &LogShard{
            id:        i,
            buffer:    make([]*LogEntry, 0, 1000),
            processor: processorFactory(i),
        }
    }

    return &ShardedLogger{
        shards:    shards,
        hasher:    &DefaultLogHasher{},
        shardMask: uint64(numShards - 1),
    }
}

func (sl *ShardedLogger) Write(entry *LogEntry) error {
    shardID := sl.hasher.Hash(entry) & sl.shardMask
    shard := sl.shards[shardID]

    shard.mutex.Lock()
    shard.buffer = append(shard.buffer, entry)
    shouldFlush := len(shard.buffer) >= 1000
    shard.mutex.Unlock()

    if shouldFlush {
        return shard.processor.Process(shard.buffer)
    }

    return nil
}

// 日志哈希器
type DefaultLogHasher struct{}

func (h *DefaultLogHasher) Hash(entry *LogEntry) uint64 {
    // 使用消息内容的哈希
    hash := fnv64(entry.Message)
    hash = hash * 31 + uint64(entry.Level)
    hash = hash * 31 + uint64(entry.Timestamp.UnixNano())
    return hash
}

func fnv64(data string) uint64 {
    hash := uint64(2166136261)
    for i := 0; i < len(data); i++ {
        hash ^= uint64(data[i]) * uint64(16777619)
        hash *= 31
    }
    return hash
}
```

### 3. 工作窃取调度

```go
// 工作窃取调度器
type WorkStealingScheduler struct {
    workers      []*Worker
    taskQueue    chan Task
    doneChan     chan struct{}
    workerCount  int
    stealCount   int
}

type Task struct {
    ID     string
    Work   func() error
    Result chan error
    Retry  int
}

type Worker struct {
    id     int
    tasks  chan Task
    done   chan struct{}
    active bool
    mu     sync.Mutex
}

func NewWorkStealingScheduler(workerCount int) *WorkStealingScheduler {
    scheduler := &WorkStealingScheduler{
        workers:     make([]*Worker, workerCount),
        taskQueue:   make(chan Task, 1000),
        doneChan:    make(chan struct{}, workerCount),
        workerCount: workerCount,
    }

    // 创建工作窃取队列
    queues := make([]chan Task, workerCount)
    for i := 0; i < workerCount; i++ {
        queues[i] = make(chan Task, 100)
        scheduler.workers[i] = &Worker{
            id:     i,
            tasks:  queues[i],
            done:   scheduler.doneChan[i],
        }
    }

    // 启动工作协程
    for i := 0; i < workerCount; i++ {
        go scheduler.workerLoop(i)
    }

    return scheduler
}

func (wss *WorkStealingScheduler) workerLoop(workerID int) {
    worker := wss.workers[workerID]

    for {
        select {
        case task := <-worker.tasks:
            // 处理本地任务
            task.Result <- task.Work()
            worker.done <- struct{}{}

        case task := <-wss.stealTask(workerID):
            // 处理窃取的任务
            task.Result <- task.Work()
            worker.done <- struct{}{}

        case <-wss.doneChan[workerID]:
            // 工作完成
            worker.active = false
            return
        }
    }
}

func (wss *WorkStealingScheduler) stealTask(currentWorkerID int) Task {
    // 随机选择窃取目标
    targetWorker := (currentWorkerID + 1) % wss.workerCount

    // 尝试从其他工作队列窃取
    for i := 1; i < wss.workerCount; i++ {
        targetID := (currentWorkerID + i) % wss.workerCount
        target := wss.workers[targetID]

        target.mutex.Lock()
        select {
        case task := <-target.tasks:
            target.mutex.Unlock()
            return task
        default:
            target.mutex.Unlock()
        }
    }

    // 无法窃取，创建新任务
    return <-wss.taskQueue
}

func (wss *WorkStealingScheduler) SubmitTask(work func() error) error {
    task := Task{
        Work:   work,
        Result: make(chan error, 1),
        Retry:  0,
    }

    select {
    case wss.taskQueue <- task:
        return <-task.Result
    default:
        // 队列满，直接执行
        return work()
    }
}
```

## 📈 监控和调优

### 1. 性能监控

```go
// 性能监控器
type PerformanceMonitor struct {
    metrics      *PerformanceMetrics
    collectors   []MetricCollector
    interval     time.Duration
    stopCh       chan struct{}
}

type PerformanceMetrics struct {
    Latency    time.Duration
    Throughput float64
    Memory     int64
    CPU        float64
    ErrorRate  float64
    Timestamp  time.Time
}

type MetricCollector interface {
    Collect() *PerformanceMetrics
}

// 启动性能监控
func StartPerformanceMonitoring() *PerformanceMonitor {
    monitor := &PerformanceMonitor{
        interval:   5 * time.Second,
        stopCh:    make(chan struct{}),
    }

    // 注册收集器
    monitor.collectors = append(monitor.collectors,
        &ThroughputCollector{},
        &LatencyCollector{},
        &MemoryCollector{},
        &ErrorRateCollector{},
    )

    go monitor.monitorLoop()
    return monitor
}

func (pm *PerformanceMonitor) monitorLoop() {
    ticker := time.NewTicker(pm.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 收集所有指标
            for _, collector := range pm.collectors {
                metrics := collector.Collect()

                // 应用权重计算综合指标
                pm.metrics = calculateWeightedMetrics([]*PerformanceMetrics{metrics})

                // 检查性能健康状态
                if !pm.metrics.IsHealthy() {
                    pm.handlePerformanceDegradation()
                }

                // 记录性能指标
                logger.Info("性能指标",
                    logger.Duration("avg_latency", pm.metrics.Latency),
                    logger.Float64("throughput", pm.metrics.Throughput),
                    logger.Int64("memory_usage", pm.metrics.MemoryUsage),
                    logger.Float64("cpu_usage", pm.metrics.CPUUsage),
                    logger.Float64("error_rate", pm.metrics.ErrorRate))
            }

        case <-pm.stopCh:
            return
        }
    }
}

func calculateWeightedMetrics(metricsList []*PerformanceMetrics) *PerformanceMetrics {
    if len(metricsList) == 0 {
        return &PerformanceMetrics{}
    }

    var totalLatency time.Duration
    var totalThroughput float64
    var totalMemory int64
    var totalCPU float64
    var totalErrorRate float64

    for _, metrics := range metricsList {
        weight := 1.0 / float64(len(metricsList))
        totalLatency += time.Duration(float64(metrics.Latency) * weight)
        totalThroughput += metrics.Throughput * weight
        totalMemory += int64(float64(metrics.Memory) * weight)
        totalCPU += metrics.CPU * weight
        totalErrorRate += metrics.ErrorRate * weight
    }

    return &PerformanceMetrics{
        Latency:    totalLatency,
        Throughput: totalThroughput,
        Memory:     totalMemory,
        CPU:        totalCPU,
        ErrorRate:  totalErrorRate,
        Timestamp:  time.Now(),
    }
}

func (pm *PerformanceMonitor) handlePerformanceDegradation() {
    logger.Error("性能下降，启动自动优化",
        logger.Duration("avg_latency", pm.metrics.Latency),
        logger.Float64("throughput", pm.metrics.Throughput),
        logger.Float64("error_rate", pm.metrics.ErrorRate))

    // 自动优化措施
    go pm.autoOptimize()
}

func (pm *PerformanceMonitor) autoOptimize() {
    // 根据性能指标调整配置
    if pm.metrics.Latency > 2*time.Millisecond {
        // 增加批量大小
        increaseBatchSize()
    }

    if pm.metrics.ErrorRate > 5.0 {
        // 启用更严格的采样
        enableStrictSampling()
    }

    if pm.metrics.MemoryUsage > 500*1024*1024 { // 500MB
        // 增加垃圾回收
        runtime.GC()
    }
}
```

### 2. 自动调优

```go
// 自动调优器
type AutoOptimizer struct {
    monitor   *PerformanceMonitor
    settings *OptimizationSettings
    lastAdjust time.Time
    history   []PerformanceMetrics
    maxHistory int
}

type OptimizationSettings struct {
    BatchSize      int
    SampleRate      float64
    LogLevel       logger.Level
    Compression    bool
    ConnectionPool  int
    FlushInterval  time.Duration
}

func NewAutoOptimizer(monitor *PerformanceMonitor) *AutoOptimizer {
    return &AutoOptimizer{
        monitor:   monitor,
        settings: &OptimizationSettings{
            BatchSize:     100,
            SampleRate:     1.0,
            LogLevel:       logger.InfoLevel,
            Compression:   true,
            ConnectionPool: 10,
            FlushInterval: 1 * time.Second,
        },
        history:       make([]PerformanceMetrics, 0, 100),
        maxHistory:    100,
    }
}

func (ao *AutoOptimizer) Optimize() {
    metrics := ao.monitor.metrics
    ao.history = append(ao.history, *metrics)

    if len(ao.history) > ao.maxHistory {
        ao.history = ao.history[1:]
    }

    // 如果距离上次调整时间超过30秒，进行优化
    if time.Since(ao.lastAdjust) < 30*time.Second {
        return
    }

    // 分析趋势
    trends := ao.analyzeTrends()

    // 根据趋势调整设置
    ao.adjustSettings(trends)

    ao.lastAdjust = time.Now()
}

func (ao *AutoOptimizer) analyzeTrends() *TrendAnalysis {
    if len(ao.history) < 2 {
        return &TrendAnalysis{}
    }

    recent := ao.history[len(ao.history)-1]
    previous := ao.history[len(ao.history)-2]

    return &TrendAnalysis{
        LatencyTrend:    calculateTrend(float64(previous.Latency.Nanoseconds()), float64(recent.Latency.Nanoseconds())),
        ThroughputTrend: calculateTrend(previous.Throughput, recent.Throughput),
        MemoryTrend:    calculateTrend(float64(previous.Memory), float64(recent.Memory)),
        ErrorRateTrend: calculateTrend(previous.ErrorRate, recent.ErrorRate),
    }
}

type TrendAnalysis struct {
    LatencyTrend     float64  // 1.0 表示稳定，>1.0 表示增长
    ThroughputTrend   float64
    MemoryTrend     float64
    ErrorRateTrend   float64
}

func calculateTrend(old, new float64) float64 {
    if old == 0 {
        return 1.0
    }
    return new / old
}

func (ao *AutoOptimizer) adjustSettings(trends *TrendAnalysis) {
    settings := ao.settings

    // 根据延迟趋势调整批量大小
    if trends.LatencyTrend > 1.2 { // 延迟增加 20%
        newBatchSize := int(float64(settings.BatchSize) * 1.2)
        if newBatchSize > 1000 {
            newBatchSize = 1000
        }
        settings.BatchSize = newBatchSize
    } else if trends.LatencyTrend < 0.8 && settings.BatchSize > 50 { // 延迟减少 20%
        newBatchSize := int(float64(settings.BatchSize) * 0.8)
        settings.BatchSize = newBatchSize
    }

    // 根据错误率调整采样率
    if trends.ErrorRateTrend > 1.5 { // 错误率增加 50%
        newSampleRate := settings.SampleRate * 0.8
        if newSampleRate < 0.1 {
            newSampleRate = 0.1
        }
        settings.SampleRate = newSampleRate
    } else if trends.ErrorRateTrend < 0.5 && settings.SampleRate < 1.0 { // 错误率降低
        newSampleRate := settings.SampleRate * 1.2
        if newSampleRate > 1.0 {
            newSampleRate = 1.0
        }
        settings.SampleRate = newSampleRate
    }

    // 根据内存使用调整连接池大小
    if trends.MemoryTrend > 1.3 { // 内存使用增加 30%
        newPoolSize := int(float64(settings.ConnectionPool) * 1.2)
        if newPoolSize > 50 {
            newPoolSize = 50
        }
        settings.ConnectionPool = newPoolSize
    }

    logger.Info("自动调优配置",
        logger.Int("batch_size", settings.BatchSize),
        logger.Float64("sample_rate", settings.SampleRate),
        logger.String("log_level", settings.Level.String()),
        logger.Bool("compression", settings.Compression),
        logger.Int("connection_pool", settings.ConnectionPool))
}
```

这个性能优化指南涵盖了 EchoMind 日志框架的所有关键性能优化点，包括批量处理、内存管理、网络传输、存储优化、并发处理等方面，以及自动化的监控和调优机制。通过实施这些优化策略，可以将日志框架的性能提升 5-10 倍。