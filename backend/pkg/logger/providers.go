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

// NoopProvider 空操作提供者
type NoopProvider struct{}

func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

func (p *NoopProvider) Write(ctx context.Context, entry *LogEntry) error {
	return nil
}

func (p *NoopProvider) Close() error {
	return nil
}

func (p *NoopProvider) Ping() error {
	return nil
}

// ElasticsearchProvider Elasticsearch 提供者
type ElasticsearchProvider struct {
	client    *http.Client
	index     string
	batchSize int
	buffer    []*LogEntry
	mu        sync.Mutex
	url       string
}

func NewElasticsearchProvider(settings map[string]interface{}) (*ElasticsearchProvider, error) {
	url, ok := settings["url"].(string)
	if !ok {
		return nil, fmt.Errorf("elasticsearch provider requires 'url' setting")
	}

	index, ok := settings["index"].(string)
	if !ok {
		index = "echomind-logs"
	}

	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	return &ElasticsearchProvider{
		client:    &http.Client{Timeout: 5 * time.Second},
		url:       url,
		index:     index,
		batchSize: batchSize,
		buffer:    make([]*LogEntry, 0, batchSize),
	}, nil
}

func (p *ElasticsearchProvider) Write(ctx context.Context, entry *LogEntry) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer = append(p.buffer, entry)

	if len(p.buffer) >= p.batchSize {
		return p.flushBuffer(ctx)
	}

	return nil
}

func (p *ElasticsearchProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.buffer) > 0 {
		return p.flushBuffer(context.Background())
	}
	return nil
}

func (p *ElasticsearchProvider) Ping() error {
	resp, err := p.client.Get(p.url + "/_cluster/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *ElasticsearchProvider) flushBuffer(ctx context.Context) error {
	if len(p.buffer) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for i, entry := range p.buffer {
		// Elasticsearch bulk format
		indexOp := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": p.index,
				"_id":    fmt.Sprintf("%d-%d", entry.Timestamp.UnixNano(), i),
			},
		}
		indexBytes, _ := json.Marshal(indexOp)
		entryBytes, _ := json.Marshal(entry)

		buf.Write(indexBytes)
		buf.WriteByte('\n')
		buf.Write(entryBytes)
		buf.WriteByte('\n')
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.url+"/_bulk", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p.buffer = p.buffer[:0] // 清空缓冲区
	return nil
}

// LokiProvider Grafana Loki 提供者
type LokiProvider struct {
	client    *http.Client
	batchSize int
	buffer    []*LogEntry
	mu        sync.Mutex
	url       string
	labels    map[string]string
}

func NewLokiProvider(settings map[string]interface{}) (*LokiProvider, error) {
	url, ok := settings["url"].(string)
	if !ok {
		return nil, fmt.Errorf("loki provider requires 'url' setting")
	}

	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	labels := make(map[string]string)
	if l, ok := settings["labels"].(map[string]interface{}); ok {
		for k, v := range l {
			if str, ok := v.(string); ok {
				labels[k] = str
			}
		}
	}

	return &LokiProvider{
		client:    &http.Client{Timeout: 5 * time.Second},
		url:       url + "/loki/api/v1/push",
		batchSize: batchSize,
		buffer:    make([]*LogEntry, 0, batchSize),
		labels:    labels,
	}, nil
}

func (p *LokiProvider) Write(ctx context.Context, entry *LogEntry) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer = append(p.buffer, entry)

	if len(p.buffer) >= p.batchSize {
		return p.flushBuffer(ctx)
	}

	return nil
}

func (p *LokiProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.buffer) > 0 {
		return p.flushBuffer(context.Background())
	}
	return nil
}

func (p *LokiProvider) Ping() error {
	readyURL := p.url[:len(p.url)-len("/push")] + "/ready"
	resp, err := p.client.Get(readyURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *LokiProvider) flushBuffer(ctx context.Context) error {
	if len(p.buffer) == 0 {
		return nil
	}

	streams := make([]map[string]interface{}, 0, len(p.buffer))

	for _, entry := range p.buffer {
		// 合并全局标签和上下文标签
		labels := make(map[string]string)
		for k, v := range p.labels {
			labels[k] = v
		}

		if entry.Context.OrgID != "" {
			labels["org_id"] = entry.Context.OrgID
		}
		if entry.Context.UserID != "" {
			labels["user_id"] = entry.Context.UserID
		}
		if entry.Context.TraceID != "" {
			labels["trace_id"] = entry.Context.TraceID
		}
		labels["level"] = entry.Level.String()

		// 将字段转换为 JSON 字符串
		fieldsJson, _ := json.Marshal(entry.Fields)

		stream := map[string]interface{}{
			"stream": labels,
			"values": [][]interface{}{
				{
					fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
					fmt.Sprintf("%s %s", entry.Message, string(fieldsJson)),
				},
			},
		}
		streams = append(streams, stream)
	}

	payload := map[string]interface{}{
		"streams": streams,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p.buffer = p.buffer[:0] // 清空缓冲区
	return nil
}

// SplunkProvider Splunk HTTP Event Collector 提供者
type SplunkProvider struct {
	client    *http.Client
	token     string
	index     string
	source    string
	batchSize int
	buffer    []*LogEntry
	mu        sync.Mutex
	url       string
}

func NewSplunkProvider(settings map[string]interface{}) (*SplunkProvider, error) {
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

	source := "echomind-backend"
	if src, ok := settings["source"].(string); ok {
		source = src
	}

	batchSize := 100
	if bs, ok := settings["batch_size"].(int); ok {
		batchSize = bs
	}

	return &SplunkProvider{
		client:    &http.Client{Timeout: 5 * time.Second},
		url:       url + "/services/collector/event",
		token:     token,
		index:     index,
		source:    source,
		batchSize: batchSize,
		buffer:    make([]*LogEntry, 0, batchSize),
	}, nil
}

func (p *SplunkProvider) Write(ctx context.Context, entry *LogEntry) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer = append(p.buffer, entry)

	if len(p.buffer) >= p.batchSize {
		return p.flushBuffer(ctx)
	}

	return nil
}

func (p *SplunkProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.buffer) > 0 {
		return p.flushBuffer(context.Background())
	}
	return nil
}

func (p *SplunkProvider) Ping() error {
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

func (p *SplunkProvider) flushBuffer(ctx context.Context) error {
	if len(p.buffer) == 0 {
		return nil
	}

	for _, entry := range p.buffer {
		event := map[string]interface{}{
			"time":       entry.Timestamp.Unix(),
			"index":      p.index,
			"source":     p.source,
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

		req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Splunk "+p.token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	p.buffer = p.buffer[:0] // 清空缓冲区
	return nil
}
