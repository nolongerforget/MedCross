# MedCross 数据库设计文档

## 概述

MedCross平台使用PostgreSQL关系型数据库存储用户信息、医疗数据元数据、授权记录和跨链交易信息。实际的医疗数据内容通过区块链（以太坊和Hyperledger Fabric）进行存储和共享。

## 数据库配置

在`backend/config/config.yaml`中添加以下数据库配置：

```yaml
# 数据库配置
database:
  driver: postgres
  host: localhost
  port: 5432
  user: medcross
  password: your_secure_password
  dbname: medcross
  sslmode: disable
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

## 表结构设计

### 1. users 表

存储系统用户信息，包括医疗机构人员、研究人员等。

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    org_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### 2. medical_data 表

存储医疗数据的元数据信息，实际数据内容存储在区块链上。

```sql
CREATE TABLE medical_data (
    data_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

-- 索引
CREATE INDEX idx_medical_data_owner ON medical_data(owner_id);
CREATE INDEX idx_medical_data_type ON medical_data(data_type);
CREATE INDEX idx_medical_data_ethereum_tx ON medical_data(ethereum_tx_id);
CREATE INDEX idx_medical_data_fabric_tx ON medical_data(fabric_tx_id);
```

### 3. data_tags 表

存储医疗数据的标签，与medical_data表是多对多关系。

```sql
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE data_tags (
    data_id UUID REFERENCES medical_data(data_id),
    tag_id INTEGER REFERENCES tags(id),
    PRIMARY KEY (data_id, tag_id)
);

-- 索引
CREATE INDEX idx_data_tags_data_id ON data_tags(data_id);
CREATE INDEX idx_data_tags_tag_id ON data_tags(tag_id);
```

### 4. authorizations 表

存储医疗数据的授权信息，记录谁可以访问哪些数据。

```sql
CREATE TABLE authorizations (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id),
    owner_id UUID NOT NULL REFERENCES users(id),
    authorized_user_id UUID NOT NULL REFERENCES users(id),
    start_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_authorization UNIQUE(data_id, authorized_user_id)
);

-- 索引
CREATE INDEX idx_authorizations_data_id ON authorizations(data_id);
CREATE INDEX idx_authorizations_owner_id ON authorizations(owner_id);
CREATE INDEX idx_authorizations_authorized_user ON authorizations(authorized_user_id);
CREATE INDEX idx_authorizations_active ON authorizations(is_active);
```

### 5. access_records 表

记录医疗数据的访问历史。

```sql
CREATE TABLE access_records (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id),
    accessor_id UUID NOT NULL REFERENCES users(id),
    timestamp BIGINT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    chain_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_access_records_data_id ON access_records(data_id);
CREATE INDEX idx_access_records_accessor ON access_records(accessor_id);
CREATE INDEX idx_access_records_timestamp ON access_records(timestamp);
```

### 6. cross_chain_references 表

维护不同区块链上同一数据的引用关系。

```sql
CREATE TABLE cross_chain_references (
    id SERIAL PRIMARY KEY,
    data_id UUID NOT NULL REFERENCES medical_data(data_id),
    ethereum_tx_id VARCHAR(100),
    fabric_tx_id VARCHAR(100),
    sync_status VARCHAR(20) NOT NULL,
    last_sync_time BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_cross_chain_data_id ON cross_chain_references(data_id);
CREATE INDEX idx_cross_chain_ethereum_tx ON cross_chain_references(ethereum_tx_id);
CREATE INDEX idx_cross_chain_fabric_tx ON cross_chain_references(fabric_tx_id);
CREATE INDEX idx_cross_chain_sync_status ON cross_chain_references(sync_status);
```

## 数据库迁移

### 初始化数据库

1. 创建数据库和用户：

```sql
CREATE DATABASE medcross;
CREATE USER medcross WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE medcross TO medcross;
```

2. 连接到数据库并执行上述建表语句：

```bash
psql -U medcross -d medcross -f init.sql
```

## 数据库连接配置

在`backend/config/config.go`中添加数据库配置结构：

```go
// Config 存储应用程序配置
type Config struct {
    // 服务器配置
    Server struct {
        Port string `mapstructure:"port"`
    }

    // 数据库配置
    Database struct {
        Driver          string `mapstructure:"driver"`
        Host            string `mapstructure:"host"`
        Port            string `mapstructure:"port"`
        User            string `mapstructure:"user"`
        Password        string `mapstructure:"password"`
        DBName          string `mapstructure:"dbname"`
        SSLMode         string `mapstructure:"sslmode"`
        MaxIdleConns    int    `mapstructure:"max_idle_conns"`
        MaxOpenConns    int    `mapstructure:"max_open_conns"`
        ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
    }

    // 以太坊配置
    Ethereum struct {
        RPCURL          string `mapstructure:"rpc_url"`
        ContractAddress string `mapstructure:"contract_address"`
        PrivateKey      string `mapstructure:"private_key"`
    }

    // Fabric配置
    Fabric struct {
        ConfigPath    string `mapstructure:"config_path"`
        ChannelID     string `mapstructure:"channel_id"`
        ChaincodeName string `mapstructure:"chaincode_name"`
        UserName      string `mapstructure:"user_name"`
        OrgName       string `mapstructure:"org_name"`
    }

    // JWT配置
    JWT struct {
        Secret    string `mapstructure:"secret"`
        ExpiresIn int    `mapstructure:"expires_in"`
    }
}
```

在`setDefaults()`函数中添加数据库默认配置：

```go
// 数据库默认配置
viper.SetDefault("database.driver", "postgres")
viper.SetDefault("database.host", "localhost")
viper.SetDefault("database.port", "5432")
viper.SetDefault("database.user", "medcross")
viper.SetDefault("database.password", "")
viper.SetDefault("database.dbname", "medcross")
viper.SetDefault("database.sslmode", "disable")
viper.SetDefault("database.max_idle_conns", 10)
viper.SetDefault("database.max_open_conns", 100)
viper.SetDefault("database.conn_max_lifetime", 3600)
```

## 数据库连接管理

创建`backend/database/db.go`文件用于管理数据库连接：

```go
package database

import (
    "fmt"
    "log"
    "time"

    "github.com/your-org/medcross/backend/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
    dbConfig := config.AppConfig.Database
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    }

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        log.Fatalf("无法连接到数据库: %v", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("获取数据库连接失败: %v", err)
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
    sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

    log.Println("数据库连接成功")
}

// CloseDB 关闭数据库连接
func CloseDB() {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err != nil {
            log.Printf("获取数据库连接失败: %v\n", err)
            return
        }
        sqlDB.Close()
        log.Println("数据库连接已关闭")
    }
}
```

## 数据库模型与GORM集成

修改现有模型以支持GORM：

```go
// 在 models/user.go 中
package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"type:varchar(100);not null" json:"password,omitempty"`
    Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    FullName  string         `gorm:"type:varchar(100);not null" json:"fullName"`
    OrgName   string         `gorm:"type:varchar(100);not null" json:"orgName"`
    Role      string         `gorm:"type:varchar(20);not null" json:"role"`
    CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 自动迁移

在`main.go`中添加自动迁移代码：

```go
// 在 import 部分添加
"github.com/your-org/medcross/backend/database"
"github.com/your-org/medcross/backend/models"

// 在 main() 函数中添加
func main() {
    // 初始化配置
    config.InitConfig()
    
    // 初始化数据库
    database.InitDB()
    defer database.CloseDB()
    
    // 自动迁移数据库模型
    database.DB.AutoMigrate(
        &models.User{},
        &models.MedicalData{},
        &models.Authorization{},
        &models.AccessRecord{},
        &models.CrossChainReference{},
    )
    
    // 其他初始化代码...
}
```

## 数据库备份与恢复

### 备份

```bash
pg_dump -U medcross -d medcross -F c -f medcross_backup.dump
```

### 恢复

```bash
pg_restore -U medcross -d medcross -c medcross_backup.dump
```

## 注意事项

1. 生产环境中，请确保使用强密码并启用SSL连接
2. 定期备份数据库
3. 考虑使用数据库连接池管理工具，如PgBouncer，以优化连接性能
4. 敏感数据（如用户密码）应当在存储前进行加密或哈希处理