package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MqttInfo struct {
	Host     string
	Port     int
	UserName string
	PassWord string
	SubList  []string
	Qos      int
}

type MapTable struct {
	StartAddr int
	DataLen   int
	Type      string
	Name      string
}

type Device struct {
	Register string
	PubTopic string
	Maptable []MapTable
}

type TcpModbus struct {
	Enable   bool
	Host     string
	Port     int
	SlaveID  int
	Interval int
	Devices  []Device
}

type Config struct {
	Mqttinfo  MqttInfo
	Tcpmodbus TcpModbus
}

func GetConfig() Config {

	config := Config{}
	config.Mqttinfo.Host = "127.0.0.1"
	config.Mqttinfo.Port = 1883
	config.Mqttinfo.Qos = 0

	//modbus配置
	config.Tcpmodbus.Enable = false
	config.Tcpmodbus.Host = "127.0.0.1"
	config.Tcpmodbus.Port = 502
	config.Tcpmodbus.SlaveID = 1
	config.Tcpmodbus.Interval = 3

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found")
		} else {
			fmt.Println(err.Error())
		}
		return config
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	return config
}
