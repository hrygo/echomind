/**
 * 搜索增强功能相关的 TypeScript 类型定义
 * 对应后端 API: /api/v1/search
 */

// 搜索结果项
export interface SearchResult {
  email_id: string;
  subject: string;
  sender: string;
  snippet: string;
  date: string;
  score: number;
  context_id?: string;
}

// 聚类类型
export type ClusterType = 'sender' | 'time' | 'topic';

// 搜索聚类
export interface SearchCluster {
  id: string;
  type: ClusterType;
  label: string;
  count: number;
  results: SearchResult[];
}

// 搜索结果摘要
export interface SearchResultsSummary {
  natural_summary: string;      // 自然语言总结
  key_topics: string[];          // 关键主题
  important_people: string[];    // 重要联系人
}

// 搜索 API 响应
export interface SearchResponse {
  query: string;
  results: SearchResult[];
  count: number;
  // 可选的增强功能字段
  clusters?: SearchCluster[];
  cluster_type?: ClusterType;
  summary?: SearchResultsSummary;
}

// 搜索选项
export interface SearchOptions {
  query: string;
  limit?: number;
  context_id?: string;
  // 增强功能开关
  enable_clustering?: boolean;
  cluster_type?: ClusterType;
  enable_summary?: boolean;
  // 过滤器
  sender?: string;
  start_date?: string;
  end_date?: string;
}

// 视图模式
export type SearchViewMode = 'all' | 'clustered';

// Copilot 模式
export type CopilotMode = 'idle' | 'search' | 'chat';
