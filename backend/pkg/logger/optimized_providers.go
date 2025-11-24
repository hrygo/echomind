package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// OptimizedElasticsearchProvider 优化的 Elasticsearch 提供者
type OptimizedElasticsearchProvider struct {
	*Batches
	url       string
	index     string
	client    *http.Client
	mu        sync.RWMutex
	batchPool sync.Pool // 对象池重用 batch
}

// ElasticsearchBatch Elasticsearch 批量请求
type ElasticsearchBatch struct {
	Body bytes.Buffer
}

// NewOptimizedElasticsearchProvider 创建优化的 Elasticsearch 提供者
func NewOptimizedElasticsearchProvider(settings map[string]interface{}) (*OptimizedElasticsearchProvider, error) {
	url, ok := settings["url"].(string)
	if !ok {
		return nil, fmt.Errorf("elasticsearch provider requires 'url' setting")
	}

	index, ok := settings["index"].(string)
	if !ok {
		index = "echomind-logs"
	}

	// 配置 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	provider := &OptimizedElasticsearchProvider{
		url:    url,
		index:  index,
		client: client,
	}

	// 初始化对象池
	provider.batchPool = sync.Pool{
		New: func() interface{} {
			return &ElasticsearchBatch{
				Body: bytes.Buffer{},
			}
		},
	}

	// 创建批量处理器
	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	provider.Batches = NewBatches(provider, batchSize, 5*time.Second)
	provider.AddProcessor(NewBatchProcessor(provider, batchSize, 5*time.Second))

	return provider, nil
}

// Write 实现提供者接口
func (p *OptimizedElasticsearchProvider) Write(ctx context.Context, entry *LogEntry) error {
	// 使用批量处理器
	return p.Batches.Write(ctx, entry)
}

// WriteBatch 直接写入批量数据
func (p *OptimizedElasticsearchProvider) WriteBatch(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// 从对象池获取 batch
	batch := p.batchPool.Get().(*ElasticsearchBatch)
	defer func() {
		batch.Body.Reset()
		p.batchPool.Put(batch)
	}()

	// 构建批量请求
	if err := p.buildBulkRequest(&batch.Body, entries); err != nil {
		return err
	}

	// 发送请求
	req, err := http.NewRequestWithContext(context.Background(), "POST", p.url+"/_bulk", &batch.Body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("elasticsearch bulk request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// buildBulkRequest 构建 Elasticsearch 批量请求
func (p *OptimizedElasticsearchProvider) buildBulkRequest(buf *bytes.Buffer, entries []*LogEntry) error {
	buf.Grow(len(entries) * 500) // 预分配缓冲区大小

	for i, entry := range entries {
		// 索引操作
		indexOp := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": p.index,
				"_id":    fmt.Sprintf("%d-%d", entry.Timestamp.UnixNano(), i),
			},
		}

		if err := json.NewEncoder(buf).Encode(indexOp); err != nil {
			return err
		}
		buf.WriteByte('\n')

		// 文档内容
		if err := json.NewEncoder(buf).Encode(entry); err != nil {
			return err
		}
		buf.WriteByte('\n')
	}

	return nil
}

// Ping 健康检查
func (p *OptimizedElasticsearchProvider) Ping() error {
	resp, err := p.client.Get(p.url + "/_cluster/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Close 关闭提供者
func (p *OptimizedElasticsearchProvider) Close() error {
	p.Batches.Close()
	p.client.CloseIdleConnections()
	return nil
}

// OptimizedLokiProvider 优化的 Loki 提供者
type OptimizedLokiProvider struct {
	*Batches
	url    string
	labels map[string]string
	client *http.Client
	mu     sync.RWMutex
}

// NewOptimizedLokiProvider 创建优化的 Loki 提供者
func NewOptimizedLokiProvider(settings map[string]interface{}) (*OptimizedLokiProvider, error) {
	url, ok := settings["url"].(string)
	if !ok {
		return nil, fmt.Errorf("loki provider requires 'url' setting")
	}

	// 确保 URL 包含完整的路径
	if !hasSuffix(url, "/push") {
		url = url + "/loki/api/v1/push"
	}

	labels := make(map[string]string)
	if l, ok := settings["labels"].(map[string]interface{}); ok {
		for k, v := range l {
			if str, ok := v.(string); ok {
				labels[k] = str
			}
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	provider := &OptimizedLokiProvider{
		url:    url,
		labels: labels,
		client: client,
	}

	// 创建批量处理器
	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	provider.Batches = NewBatches(provider, batchSize, 3*time.Second)
	provider.AddProcessor(NewBatchProcessor(provider, batchSize, 3*time.Second))

	return provider, nil
}

// Write 实现提供者接口
func (p *OptimizedLokiProvider) Write(ctx context.Context, entry *LogEntry) error {
	// 检查是否应该采样（对于 Loki，只记录 WARN 和 ERROR 级别）
	if entry.Level < WarnLevel {
		return nil
	}

	return p.Batches.Write(ctx, entry)
}

// WriteBatch 直接写入批量数据
func (p *OptimizedLokiProvider) WriteBatch(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// 按 labels 分组
	streams := make(map[string][]LokiStream)
	for _, entry := range entries {
		// 合并全局标签和上下文标签
		streamLabels := make(map[string]string)
		for k, v := range p.labels {
			streamLabels[k] = v
		}

		// 添加动态标签
		streamLabels["level"] = entry.Level.String()
		if entry.Context.OrgID != "" {
			streamLabels["org_id"] = entry.Context.OrgID
		}
		if entry.Context.UserID != "" {
			streamLabels["user_id"] = entry.Context.UserID
		}
		if entry.Context.TraceID != "" {
			streamLabels["trace_id"] = entry.Context.TraceID
		}

		// 创建 labels 的字符串表示
		labelsStr := p.labelsToString(streamLabels)

		// 创建 stream
		stream := streams[labelsStr]
		stream = append(stream, LokiStream{
			Stream: streamLabels,
			Values: [][]interface{}{
				{
					fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
					p.formatLogMessage(entry),
				},
			},
		})
		streams[labelsStr] = stream
	}

	// 构建 payload
	payload := LokiPayload{Streams: make([]LokiStream, 0, len(streams))}
	for _, stream := range streams {
		payload.Streams = append(payload.Streams, stream...)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 发送请求
	req, err := http.NewRequestWithContext(context.Background(), "POST", p.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("loki push failed with status: %d", resp.StatusCode)
	}

	return nil
}

// labelsToString 将 labels map 转换为字符串
func (p *OptimizedLokiProvider) labelsToString(labels map[string]string) string {
	var buf bytes.Buffer
	buf.Grow(len(labels) * 20) // 预分配

	buf.WriteByte('{')
	first := true
	for k, v := range labels {
		if !first {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(k)
		buf.WriteString("\"=\"")
		buf.WriteString(v)
		buf.WriteByte('"')
		first = false
	}
	buf.WriteByte('}')

	return buf.String()
}

// formatLogMessage 格式化日志消息
func (p *OptimizedLokiProvider) formatLogMessage(entry *LogEntry) string {
	fieldsJson, _ := json.Marshal(entry.Fields)
	return fmt.Sprintf("%s %s", entry.Message, string(fieldsJson))
}

// LokiPayload Loki 推送载荷
type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream Loki 流
type LokiStream struct {
	Stream map[string]string   `json:"stream"`
	Values [][]interface{}     `json:"values"`
}

// Ping 健康检查
func (p *OptimizedLokiProvider) Ping() error {
	readyURL := p.url[:len(p.url)-len("/push")] + "/ready"
	resp, err := p.client.Get(readyURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Close 关闭提供者
func (p *OptimizedLokiProvider) Close() error {
	p.Batches.Close()
	p.client.CloseIdleConnections()
	return nil
}

// OptimizedSplunkProvider 优化的 Splunk 提供者
type OptimizedSplunkProvider struct {
	*Batches
	url     string
	token   string
	index   string
	source  string
	client  *http.Client
	mu      sync.RWMutex
}

// NewOptimizedSplunkProvider 创建优化的 Splunk 提供者
func NewOptimizedSplunkProvider(settings map[string]interface{}) (*OptimizedSplunkProvider, error) {
	url, ok := settings["url"].(string)
	if !ok {
		return nil, fmt.Errorf("splunk provider requires 'url' setting")
	}

	token, ok := settings["token"].(string)
	if !ok {
		return nil, fmt.Errorf("splunk provider requires 'token' setting")
	}

	index := "echomind"
	if idx, ok := settings["index"].(string); ok {
		index = idx
	}

	source := "backend"
	if src, ok := settings["source"].(string); ok {
		source = src
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	provider := &OptimizedSplunkProvider{
		url:    url,
		token:  token,
		index:  index,
		source: source,
		client: client,
	}

	// 创建批量处理器
	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	provider.Batches = NewBatches(provider, batchSize, 2*time.Second)
	provider.AddProcessor(NewBatchProcessor(provider, batchSize, 2*time.Second))

	return provider, nil
}

// Write 实现提供者接口
func (p *OptimizedSplunkProvider) Write(ctx context.Context, entry *LogEntry) error {
	// 检查是否应该采样（对于 Splunk，只记录 ERROR 和 FATAL 级别）
	if entry.Level < ErrorLevel {
		return nil
	}

	return p.Batches.Write(ctx, entry)
}

// WriteBatch 直接写入批量数据
func (p *OptimizedSplunkProvider) WriteBatch(entries []*LogEntry) error {
	for _, entry := range entries {
		if err := p.writeSingle(entry); err != nil {
			// 记录错误但继续处理其他条目
			continue
		}
	}
	return nil
}

// writeSingle 写入单个条目
func (p *OptimizedSplunkProvider) writeSingle(entry *LogEntry) error {
	event := map[string]interface{}{
		"time": entry.Timestamp.Unix(),
		"index": p.index,
		"source": p.source,
		"sourcetype": "json",
		"event": map[string]interface{}{
			"level":   entry.Level.String(),
			"message": entry.Message,
			"fields":  entry.Fields,
			"context": entry.Context,
			"source":  entry.Source,
		},
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", p.url+"/services/collector/event", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Splunk "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("splunk event failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Ping 健康检查
func (p *OptimizedSplunkProvider) Ping() error {
	req, err := http.NewRequest("GET", p.url+"/services/collector/health", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Splunk "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Close 关闭提供者
func (p *OptimizedSplunkProvider) Close() error {
	p.Batches.Close()
	p.client.CloseIdleConnections()
	return nil
}

// 辅助函数
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}