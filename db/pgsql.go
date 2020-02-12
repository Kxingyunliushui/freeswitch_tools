package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wonderivan/logger"
)

var Pgdb *sql.DB


type Fs_info_st struct {
	Phone      string
	Start_time string
	End_time   string
	Talk_time  string
	Status	string
}

func PgsqlOpen(host, user, password, dbname string, port int) {
	var err error
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	Pgdb, err = sql.Open("postgres", pgInfo)
	//port是数据库的端口号，默认是5432，如果改了，这里一定要自定义；
	//user就是你数据库的登录帐号;
	//dbname就是你在数据库里面建立的数据库的名字;
	//sslmode就是安全验证模式;

	//还可以是这种方式打开
	//db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	checkErr(err)
}

func PgsqlClose() {
	Pgdb.Close()
}

func Pgsql_fs_info_insert(fs_info Fs_info_st) error {

	stmt, err := Pgdb.Prepare("INSERT INTO subject_identity.fs_info(phone,start_time,end_time,talk_time,status) values($1,$2,$3,$4,$5)")
	if err != nil {
		logger.Error("prepare: sql Error")
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(fs_info.Phone, fs_info.Start_time, fs_info.End_time, fs_info.Talk_time, fs_info.Status)
	if nil != err {
		logger.Error("fs_info insert failed", err)
		return err
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
