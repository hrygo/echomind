-- Migration: Add WeChat Integration Columns to Users Table
-- Description: Supports WeChat Official Account binding for Phase 7.1
-- Date: 2025-11-26

BEGIN;

-- Add WeChat-related columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(64),
ADD COLUMN IF NOT EXISTS wechat_unionid VARCHAR(64),
ADD COLUMN IF NOT EXISTS wechat_config JSONB DEFAULT '{}';

-- Create unique index on wechat_openid for fast binding lookup
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_openid ON users(wechat_openid) WHERE wechat_openid IS NOT NULL;

-- Create index on wechat_unionid for cross-app identification
CREATE INDEX IF NOT EXISTS idx_users_wechat_unionid ON users(wechat_unionid) WHERE wechat_unionid IS NOT NULL;

COMMIT;
