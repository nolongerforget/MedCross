# MedCross 数据库部署指南

## 概述

MedCross平台使用PostgreSQL关系型数据库来存储用户信息、医疗数据元数据、授权记录和跨链交易信息。本文档提供了数据库部署和配置的详细步骤。

## 环境要求

- PostgreSQL 12.0 或更高版本
- 足够的磁盘空间（建议至少10GB）
- 系统内存建议至少4GB

## 安装PostgreSQL

### Windows

1. 从[PostgreSQL官网](https://www.postgresql.org/download/windows/)下载安装程序
2. 运行安装程序，按照向导完成安装
3. 安装过程中设置管理员密码
4. 默认端口为5432，如需更改请记录新端口号

### Linux (Ubuntu/Debian)

```bash
# 安装PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### macOS

```bash
# 使用Homebrew安装
brew install postgresql

# 启动服务
brew services start postgresql
```

## 创建数据库和用户

1. 登录PostgreSQL

```bash
# Linux/macOS
sudo -u postgres psql

# Windows (使用psql命令行工具)
psql -U postgres
```

2. 创建数据库和用户

```sql
CREATE DATABASE medcross;
CREATE USER medcross WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE medcross TO medcross;
```

3. 连接到新创建的数据库

```sql
\c medcross
```

## 初始化数据库结构

### 方法1：使用SQL脚本初始化

1. 导航到项目的database目录

```bash
cd /path/to/MedCross/database
```

2. 使用psql执行初始化脚本

```bash
psql -U medcross -d medcross -f init.sql
```

### 方法2：使用应用程序自动迁移

1. 确保在`config/config.yaml`中正确配置数据库连接信息

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: medcross
  password: your_secure_password
  dbname: medcross
  sslmode: disable
```

2. 启动后端应用程序，它将自动执行数据库迁移

```bash
cd /path/to/MedCross/backend
go run main.go
```

## 数据库备份与恢复

### 备份数据库

```bash
pg_dump -U medcross -d medcross -F c -f medcross_backup.dump
```

### 恢复数据库

```bash
pg_restore -U medcross -d medcross -c medcross_backup.dump
```

## 数据库维护

### 定期清理

```sql
-- 清理已删除的数据（软删除）
VACUUM FULL;

-- 分析数据库以优化查询性能
ANALYZE;
```

### 监控数据库大小

```sql
-- 查看数据库大小
SELECT pg_size_pretty(pg_database_size('medcross'));

-- 查看表大小
SELECT relname, pg_size_pretty(pg_total_relation_size(relid))
FROM pg_catalog.pg_statio_user_tables
ORDER BY pg_total_relation_size(relid) DESC;
```

## 生产环境配置建议

1. **启用SSL连接**

   修改`config/config.yaml`中的`sslmode`为`require`：

   ```yaml
   database:
     sslmode: require
   ```

2. **配置连接池**

   ```yaml
   database:
     max_idle_conns: 10     # 空闲连接池中连接的最大数量
     max_open_conns: 100    # 打开数据库连接的最大数量
     conn_max_lifetime: 3600 # 连接可复用的最大时间（秒）
   ```

3. **定期备份**

   设置定时任务，每天自动备份数据库：

   ```bash
   # 添加到crontab
   0 2 * * * pg_dump -U medcross -d medcross -F c -f /backup/medcross_$(date +\%Y\%m\%d).dump
   ```

4. **使用PgBouncer连接池管理器**

   对于高并发场景，建议使用PgBouncer管理数据库连接池。

## 故障排除

### 连接问题

- 检查PostgreSQL服务是否运行
- 验证用户名和密码是否正确
- 确认防火墙设置允许连接到数据库端口
- 检查`pg_hba.conf`文件中的访问控制设置

### 性能问题

- 确保为常用查询创建了适当的索引
- 定期运行VACUUM和ANALYZE命令
- 监控查询性能，优化慢查询
- 考虑增加数据库服务器的资源（CPU、内存）

## 安全建议

1. 使用强密码保护数据库账户
2. 限制数据库服务器的网络访问
3. 定期更新PostgreSQL到最新版本
4. 审核数据库访问日志
5. 实施数据加密（特别是敏感医疗数据）

## 参考资源

- [PostgreSQL官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL安全最佳实践](https://www.postgresql.org/docs/current/security.html)
- [GORM文档](https://gorm.io/docs/)