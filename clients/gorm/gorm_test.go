package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"time"

	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestNewClientFromEtcd(t *testing.T) {
	etcdCli, err := etcd.NewClient(5*time.Second, "")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	type args struct {
		etcdCli *etcd.Client
		values  []interface{}
	}
	tests := []struct {
		name string
		args args
		// wantDbCli *Client
		wantErr bool
	}{
		{"should be ok", args{etcdCli, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClientFromEtcd(tt.args.etcdCli, "order", "dev", tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientFromEtcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotDbCli, tt.wantDbCli) {
			// 	t.Errorf("NewClientFromEtcd() = %v, want %v", gotDbCli, tt.wantDbCli)
			// }
		})
	}
	time.Sleep(25 * time.Second)
}

func Test_1(t *testing.T) {
	db, err := gorm.Open("mysql", "DemoCloudUser:123456@lcb@tcp(192.168.1.50:3306)/information_schema?charset=utf8mb4&parseTime=False")
	e1 := db.Exec("CREATE DATABASE IF NOT EXISTS CodeFirst default charset utf8mb4 COLLATE utf8mb4_general_ci;").Error
	fmt.Println(db, err)
	fmt.Println(e1)
}

func Test_2(t *testing.T) {
	db, err := gorm.Open("mssql", "sqlserver://sa:123456@lcb@192.168.1.29:1433") //?database=gorm
	e1 := db.Exec(`
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'CodeFirst')
BEGIN
CREATE DATABASE [CodeFirst]
END
`).Error
	fmt.Println(db, err)
	fmt.Println(e1)
}
