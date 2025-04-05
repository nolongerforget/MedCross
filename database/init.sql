-- MedCross 数据库初始化脚本

-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    org_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 医疗数据表
CREATE TABLE IF NOT EXISTS medical_data (
    data_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    data_hash VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id),
    timestamp BIGINT NOT NULL,
    is_confidential BOOLEAN DEFAULT false,
    ethereum_tx_id VARCHAR(100),
    fabric_tx_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 医疗数据表索引
CREATE INDEX IF NOT EXISTS idx_medical_data_owner ON medical_data(owner_id);
CREATE INDEX IF NOT EXISTS idx_medical_data_type ON medical_data(data_type);
CREATE INDEX IF NOT EXISTS idx_medical_data_ethereum_tx ON medical_data(ethereum_tx_id);
CREATE INDEX IF NOT EXISTS idx_medical_data_fabric_tx ON medical_data(fabric_tx_id);

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- 数据标签关联表
CREATE TABLE IF NOT EXISTS data_tags (
    data_id UUID REFERENCES medical_data(data_id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (data_id, tag_id)
);

-- 数据标签索引
CREATE INDEX IF NOT EXISTS idx_data_tags_data_id ON data_tags(data_id);
CREATE INDEX IF NOT EXISTS idx_data_tags_tag_id ON data_tags(tag_id);

-- 授权表
CREATE TABLE IF NOT EXISTS authorizations (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id),
    authorized_user_id UUID NOT NULL REFERENCES users(id),
    start_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_authorization UNIQUE(data_id, authorized_user_id)
);

-- 授权表索引
CREATE INDEX IF NOT EXISTS idx_authorizations_data_id ON authorizations(data_id);
CREATE INDEX IF NOT EXISTS idx_authorizations_owner_id ON authorizations(owner_id);
CREATE INDEX IF NOT EXISTS idx_authorizations_authorized_user ON authorizations(authorized_user_id);
CREATE INDEX IF NOT EXISTS idx_authorizations_active ON authorizations(is_active);

-- 访问记录表
CREATE TABLE IF NOT EXISTS access_records (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id) ON DELETE CASCADE,
    accessor_id UUID NOT NULL REFERENCES users(id),
    timestamp BIGINT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    chain_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 访问记录表索引
CREATE INDEX IF NOT EXISTS idx_access_records_data_id ON access_records(data_id);
CREATE INDEX IF NOT EXISTS idx_access_records_accessor ON access_records(accessor_id);
CREATE INDEX IF NOT EXISTS idx_access_records_timestamp ON access_records(timestamp);

-- 跨链引用表
CREATE TABLE IF NOT EXISTS cross_chain_references (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id) ON DELETE CASCADE,
    ethereum_tx_id VARCHAR(100),
    fabric_tx_id VARCHAR(100),
    sync_status VARCHAR(20) NOT NULL,
    last_sync_time BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 跨链引用表索引
CREATE INDEX IF NOT EXISTS idx_cross_chain_data_id ON cross_chain_references(data_id);
CREATE INDEX IF NOT EXISTS idx_cross_chain_ethereum_tx ON cross_chain_references(ethereum_tx_id);
CREATE INDEX IF NOT EXISTS idx_cross_chain_fabric_tx ON cross_chain_references(fabric_tx_id);
CREATE INDEX IF NOT EXISTS idx_cross_chain_sync_status ON cross_chain_references(sync_status);

-- 创建初始管理员用户（密码需要在实际部署时更改）
INSERT INTO users (username, password, email, full_name, org_name, role)
VALUES ('admin', '$2a$10$EqKCjGBKQDGMVLzRr1D.6.9UOLWRZGUYfPZB.wqyeAIKQKlc5I5Vy', 'admin@medcross.org', 'System Administrator', 'MedCross', 'admin')
ON CONFLICT (username) DO NOTHING;