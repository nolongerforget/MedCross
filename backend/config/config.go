package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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
		ConfigPath     string `mapstructure:"config_path"`
		ChannelID      string `mapstructure:"channel_id"`
		ChaincodeName  string `mapstructure:"chaincode_name"`
		UserName       string `mapstructure:"user_name"`
		OrgName        string `mapstructure:"org_name"`
	}

	// JWT配置
	JWT struct {
		Secret    string `mapstructure:"secret"`
		ExpiresIn int    `mapstructure:"expires_in"`
	}
}

// AppConfig 全局配置实例
var AppConfig Config

// InitConfig 初始化配置
func InitConfig() {
	// 设置默认配置
	setDefaults()

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，创建默认配置文件
			log.Println("配置文件不存在，创建默认配置文件")
			createDefaultConfigFile()
		} else {
			// 其他错误
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	// 将配置绑定到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("无法解析配置: %v", err)
	}

	log.Println("配置加载成功")
}

// 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", "8080")

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

	// 以太坊默认配置
	viper.SetDefault("ethereum.rpc_url", "http://localhost:8545")
	viper.SetDefault("ethereum.contract_address", "")
	viper.SetDefault("ethereum.private_key", "")

	// Fabric默认配置
	viper.SetDefault("fabric.config_path", "./fabric-config")
	viper.SetDefault("fabric.channel_id", "mychannel")
	viper.SetDefault("fabric.chaincode_name", "medicaldata")
	viper.SetDefault("fabric.user_name", "Admin")
	viper.SetDefault("fabric.org_name", "Org1")

	// JWT默认配置
	viper.SetDefault("jwt.secret", "medcross_secret_key")
	viper.SetDefault("jwt.expires_in", 86400) // 24小时
}

// 创建默认配置文件
func createDefaultConfigFile() {
	// 确保config目录存在
	configDir := "./config"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}

	// 写入配置文件
	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		log.Fatalf("无法创建配置文件: %v", err)
	}

	log.Printf("已创建默认配置文件: %s", configPath)
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return &AppConfig
}