package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Partition: int32(0),
		Key:       sarama.StringEncoder("key"),
	}

	var value string
	var msgType string
	for {
		_, err := fmt.Scanf("%s", &value)
		if err != nil {
			break
		}
		fmt.Scanf("%s", &msgType)
		fmt.Println("msgType = ", msgType, ",value = ", value)
		msg.Topic = msgType
		msg.Timestamp = time.Now()
		msg.Value = sarama.ByteEncoder(value)
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
	}
}
