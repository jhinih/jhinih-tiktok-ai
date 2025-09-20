package mq

//
//import (
//	"Tiktok/global"
//	"Tiktok/initialize"
//	"encoding/json"
//	"github.com/streadway/amqp"
//)
//
//type FileRequest struct {
//	Ch *amqp.Channel
//}
//
//func NewFileRequest(Ch *amqp.Channel) *FileRequest {
//	return &FileRequest{
//		Ch: Ch,
//	}
//}
//
//var Ch *amqp.Channel
//
//func (r *FileRequest) Init() {
//	//r.Ch, _ = global.MQ.Channel()
//	Ch := initialize.GetChannel()
//	Ch.ExchangeDeclare(
//		"upload",
//		"direct",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	Ch.QueueDeclare(
//		"upload_queue",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	Ch.QueueBind(
//		"upload_queue",
//		"upload_key",
//		"upload",
//		false,
//		nil,
//	)
//}
//
//func (r *FileRequest) Publish(task interface{}) error {
//	body, _ := json.Marshal(task)
//	return r.Ch.Publish(
//		"upload",
//		"upload_key",
//		false,
//		false,
//		amqp.Publishing{
//			ContentType: "application/json",
//			Body:        body,
//		},
//	)
//}
