# modbus2mqtt
读取modbus数据并通过mqtt发布



版本持续更新中，当前仅支持tcpmodbus方式



配置文件使用yaml格式，具体含义如下

```yaml
mqttinfo:

 host: 127.0.0.1 #MQTT服务器地址
 
 port: 1883   #MQTT端口号

 username: aaaa  #用户名

 password: 1234   #密码

 sublist:  #MQTT订阅列表

   - test

   - test2

 qos: 0 #qos

tcpmodbus:

 enable: True  #是否启用modbuds tcp采集数据

 host: 127.0.0.1 #modbus slave IP地址

 port: 502  #端口号

 slaveid: 1  #slave ID

 interval: 4 #读取频率 4秒/次

 devices:  #采集设备配置

    - register: holding  #要读取的寄存器 holding或 coil

      pubtopic: dev1pub  #采集的数据要发布的MQTT主题

      maptable:  #数据含义及地址映射表

       - startaddr: 0  #起始地址

          datalen: 1   #数据长度

          type: int   #数据类型

          name: temperature  #数据含义

       - startaddr: 3

          datalen: 1

          type: int

          name: humidity
```



如果成功采集到数据，程序会通过以下JSON格式发布到对应的**dev1pub**主题中：

```json
{"key":"temperature","val":1}

{"key":"humidity","val":4} 
```

