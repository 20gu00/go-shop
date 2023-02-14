package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"app_name" json:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"vesion"`
	Port         int    `mapstructure:"app_port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	MaxHeader    int    `mapstructure:"max_header"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`

	*LogConfig    `mapstructure:"log"`
	*MysqlConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
	*SmsConfig    `mapstructure:"sms"`
	NacosConfig   *NacosConfig  `mapstructure:"nacos"`
	ConsulConfig  *ConsulConfig `mapstructure:"consul"`
	UserRpcConfig *UserRpc      `mapstructure:"user-rpc"`
}

//viper读取nacos配置信息,使用nacos来配置服务的信息
//不用设置json这个,因为是本地读取
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataid"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Group     string `mapstructure:"group"`
}
type UserRpc struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"file_name"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxBackup int    `mapstructure:"max_backup"`
	Compress  bool   `mapstructure:"compress"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type SmsConfig struct {
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"key_secret"`
	Exp    int    `mapstructure:"exp"`
}
type MysqlConfig struct {
	MysqlAddr     string `mapstructure:"mysql_addr"`
	MysqlPort     int    `mapstructure:"mysql_port"`
	UserName      string `mapstructure:"user_name"`
	MysqlPassword string `mapstructure:"mysql_password"`
	DBName        string `mapstructure:"db_name"`
	MaxConnection int    `mapstructure:"max_connection"`
	MaxIdle       int    `mapstructure:"max_idle"`
}

type RedisConfig struct {
	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPort     int    `mapstructure:"redis_port"`
	RedisPassword string `mapstructure:"redis_password"`
	DB            int    `mapstructure:"db"`
	PoolSize      int    `mapstructure:"pool_size"`
	MinIdle       int    `mapstructure:"min_idle"`
}

func ConfRead(confFile string) (err error) {
	if confFile != "" {
		viper.SetConfigFile(confFile)
	} else {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./conf")
	}

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("未找到配置文件", err)
		} else {
			fmt.Println("读取配置文件失败", err)
			return err
		}
	}

	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Println("将配置信息添加进结构体失败", err)
	}

	// 不需要额外加goroutine
	// 主goroutine不退出就不影响goroutine,注意分清楚不是函数退出
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已经修改")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Println("将配置信息添加进结构体失败", err)
		}
	})

	// time.Sleep(5000 *time.Second)

	return nil
}
