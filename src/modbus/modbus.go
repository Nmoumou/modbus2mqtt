package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	config "modbus2mqtt/src/config"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	modubusraw "github.com/goburrow/modbus"
)

func GenModbusTcpClient(config config.Config) (modubusraw.Client, error) {
	// Modbus TCP
	handler := modubusraw.NewTCPClientHandler(config.Tcpmodbus.Host + ":" + strconv.Itoa(config.Tcpmodbus.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = byte(config.Tcpmodbus.SlaveID)
	// handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	if err := handler.Connect(); err != nil {
		fmt.Println("modbus connect error!")
		return nil, errors.New("modbus connect error!")
	}
	fmt.Println("tcpmodbus connect to " + config.Tcpmodbus.Host + " successful")
	// defer handler.Close()
	// results, err := client.ReadDiscreteInputs(15, 2)
	// results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
	// results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
	client := modubusraw.NewClient(handler)
	return client, nil
}

// 读取Modbus信息并根据配置文件发布到指定的主题中
func ReadModbus(client MQTT.Client, config config.Config, modbusclient modubusraw.Client) {
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		for _, item := range config.Tcpmodbus.Devices {
			if item.Register == "holding" {
				for _, mapper := range item.Maptable {
					results, err := modbusclient.ReadHoldingRegisters(uint16(mapper.StartAddr), uint16(mapper.DataLen))
					if err != nil {
						fmt.Println("read holding register error")
					} else {
						var val string
						var sendmsg string
						switch mapper.Type {
						case "int":
							val = strconv.Itoa(int(binary.BigEndian.Uint16(results)))
							sendmsg = "{\"key\":" + "\"" + mapper.Name + "\"," + "\"val\":" + val + "}"
						}
						client.Publish(item.PubTopic, 1, false, sendmsg)
						fmt.Println(sendmsg)
					}
				}
			}
		}
	}
}
