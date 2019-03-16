package Db

func AutoCreateTableTest() {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&ArticleComment{})
	db.AutoMigrate(&ArticleLabel{})
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&TopicLordReply{})
	db.AutoMigrate(&TopicLayerReply{})
}

func Test() {
	db.AutoMigrate(&User{})
	user := User{UserName: "haha", Age: 122, PhoneNumber: "213213213"}
	db.Create(&user)
	db.Save(&user)
}
