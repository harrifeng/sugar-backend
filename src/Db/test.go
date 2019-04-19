package db

import (
	"fmt"
)

func SetNewVerificationCodeTest() {
	err := SetNewVerificationCode("2323", "21321")
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func Test() {
	type Result struct {
		UserName string
		Age      int
	}
	//var result []UserReply
	var count int
	// Raw SQL
	mysqlDb.Raw(`(select id as topic_lord_reply_key,NULL as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_id as father_id from topic_lord_replies where user_id=? and deleted_at is null)
			union all
			(select  NULL as topic_lord_reply_key,id as topic_layer_reply_key,created_at,content,user_id,thumbs_up_count,
			topic_lord_reply_id as father_id from topic_layer_replies where user_id=? and deleted_at is null)
			`, 103, 103, 0, 2).Count(&count)

}
