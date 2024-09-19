package memdb

import (
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func TestMemDB_Write(t *testing.T) {
	type User struct {
		Id   int    `gorm:"column:id;primaryKey"`
		Name string `gorm:"column:name;type:varchar(128)"`
	}

	type Msg struct {
		Id        int    `gorm:"column:id;primaryKey"`
		CreatedAt int64  `gorm:"column:created_at"`
		Content   string `gorm:"column:content;type:text"`
		Status    int    `gorm:"column:status;default:0"`
		UserId    int    `gorm:"column:user_id;default:0"`
		GroupId   int    `gorm:"column:group_id;default:0"`
	}

	bs := make([]byte, 0, 1024)
	dbbs := Create("demo.db", bs)
	conn := gormlite.Open("file:/demo.db?vfs=memdb")

	var (
		err error
		db  *gorm.DB
	)

	if db, err = gorm.Open(conn); err != nil {
		log.Fatalf("gorm open err: %v", err)
	}
	db = db.Debug()

	if err = db.AutoMigrate(&User{}, &Msg{}); err != nil {
		log.Fatalf("auto migrate err: %v", err)
	}

	if err = db.Create(&User{Name: "alice"}).Error; err != nil {
		log.Fatalf("create user error: %v", err)
	}

	if err = db.Create(&User{Name: "bob"}).Error; err != nil {
		log.Fatalf("create user error: %v", err)
	}

	if err = db.Create(&([]*Msg{
		{
			CreatedAt: time.Now().UnixMilli(),
			Content:   "this is first msg content",
			Status:    1,
			UserId:    2,
			GroupId:   3,
		},
		{
			CreatedAt: time.Now().UnixMilli(),
			Content:   "this is second msg content",
			Status:    11111,
			UserId:    22222,
			GroupId:   33333,
		},
		{
			CreatedAt: time.Now().UnixMilli(),
			Content:   "&)$&W()$*&Q)(#E&($&&$&$&$&$&$",
			Status:    -1,
			UserId:    -2,
			GroupId:   -3,
		},
	})).Error; err != nil {
		log.Fatalf("create msgs error: %v", err)
	}

	list := make([]*User, 0)
	if err = db.Model(&User{}).Limit(10).Find(&list).Error; err != nil {
		log.Fatalf("find user list err: %v", err)
	}

	for _, v := range list {
		log.Printf("|| user || [%05d] name = %10s", v.Id, v.Name)
	}

	list2 := make([]*Msg, 0)
	if err = db.Model(&Msg{}).Limit(10).Find(&list2).Error; err != nil {
		log.Fatalf("find msg list err: %v", err)
	}
	for _, v := range list2 {
		log.Printf("|| msg  || [%05d] content = %s", v.Id, v.Content)
	}

	r := dbbs.Dump()
	dbs, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("read dbs err: %v", err)
	}

	os.WriteFile("dump.db", dbs, 0644)
}
