package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	config "modbus2mqtt/src/config"
	modbus "modbus2mqtt/src/modbus"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// 处理订阅到的MQTT消息
func dealMqttMsg(msg chan [2]string, exit chan bool) {
	for {
		select {
		case incoming := <-msg:
			fmt.Printf("Received message on topic: %s\nMessage: %s\n", incoming[0], incoming[1])
		case <-exit:
			return
		default:
			// fmt.Printf("empty\n")
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func main() {
	config := config.GetConfig()
	hostname, _ := os.Hostname()

	server := "tcp://" + config.Mqttinfo.Host + ":" + strconv.Itoa(config.Mqttinfo.Port)
	subtopic := config.Mqttinfo.SubList
	qos := config.Mqttinfo.Qos
	clientid := hostname + strconv.Itoa(time.Now().Second())
	username := config.Mqttinfo.UserName
	password := config.Mqttinfo.PassWord

	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)
	if username != "" {
		connOpts.SetUsername(username)
		if password != "" {
			connOpts.SetPassword(password)
		}
	}
	// tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	// connOpts.SetTLSConfig(tlsConfig)

	//自动重连机制，如网络不稳定可开启
	//connOpts.SetAutoReconnect(true)//启用自动重连功能
	//connOpts.SetMaxReconnectInterval(30)//每30秒尝试重连

	quit := make(chan bool)
	recmsg := make(chan [2]string, 300)

	connOpts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		recmsg <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	go dealMqttMsg(recmsg, quit)

	client := MQTT.NewClient(connOpts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", server)
	}
	for _, item := range subtopic {
		if token := client.Subscribe(item, byte(qos), nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		} else {
			fmt.Printf("Subscribe topic  %s  successful\n", item)
		}
	}

	if config.Tcpmodbus.Enable {
		modbusclient, error := modbus.GenModbusTcpClient(config)
		if error != nil {
			fmt.Println("modbusclient gen error!")
		} else {
			go modbus.ReadModbus(client, config, modbusclient)
		}

	}

	// 安全退出
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			quit <- true
			go func() {
				go func() {
					client.Disconnect(250)
				}()
				time.Sleep(260 * time.Millisecond)
				cleanup <- true
			}()
			<-cleanup
			log.Println("safe quit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
