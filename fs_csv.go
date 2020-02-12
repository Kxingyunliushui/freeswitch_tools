package main

//./fs_csv -path ./ test.csv
import (
	"encoding/csv"
	"flag"
	"fmt"
	"freeswitch_tools/db"
	. "github.com/0x19/goesl"
	"github.com/wonderivan/logger"
	"os"
)

var (
	fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
	//fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	//password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

type Fs_status struct {
	Status_EN  string
	Status_CN  string
	Status_ID    int
}

var fs_status = []Fs_status {
	{"NO_USER_RESPONSE", "直接挂断", 0},
	{"NO_ANSWER", "无人接听", 1},
	{"USER_BUSY", "用户忙", 2},
	{"NORMAL_CLEARING", "正常接听", 3},
}

var user_phone = make([]string, 0, 10)
var fs_info = make([]db.Fs_info_st, 0, 10)
var fs_status_map = make(map[string]Fs_status)

//NO_USER_RESPONSE  直接挂断
//NO_ANSWER  无人接听
//USER_BUSY  用户忙

func read_phone_csv(path, filename string) {
	//准备读取文件
	fileName := path + "/" +filename
	fs, err := os.Open(fileName)
	if err != nil {
		logger.Error("can not open the file, err is %+v", err)
	}
	defer fs.Close()

	//r := csv.NewReader(fs)
	////针对大文件，一行一行的读取文件
	//for {
	//	row, err := r.Read()
	//	if err != nil && err != io.EOF {
	//		logger.Error("can not read, err is %+v", err)
	//	}
	//	if err == io.EOF {
	//		break
	//	}
	//	fmt.Println(row)
	//}

	//针对小文件，也可以一次性读取所有的文件
	//注意，r要重新赋值，因为readall是读取剩下的
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		logger.Error("can not readall, err is %+v", err)
	}
	for _, row := range content {
		user_phone = append(user_phone,row[0])
		fmt.Println(row)
	}
	fmt.Println(user_phone)
}

func read_master_csv() { //读取通话结果
	//准备读取文件
	fileName := "/home/lizhangming/work/freeswitch-install/var/log/freeswitch/cdr-csv/Master.csv"
	fs, err := os.Open(fileName)
	if err != nil {
		logger.Error("can not open the file, err is %+v", err)
	}
	defer fs.Close()

	//r := csv.NewReader(fs)
	////针对大文件，一行一行的读取文件
	//for {
	//	row, err := r.Read()
	//	if err != nil && err != io.EOF {
	//		logger.Error("can not read, err is %+v", err)
	//	}
	//	if err == io.EOF {
	//		break
	//	}
	//	fmt.Println(row)
	//}

	//针对小文件，也可以一次性读取所有的文件
	//注意，r要重新赋值，因为readall是读取剩下的
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		logger.Error("can not readall, err is %+v", err)
	}
	for _, row := range content {
		fmt.Println(row)
		var fs_tmp db.Fs_info_st

		if row[0] == "Outbound Call" {
			fs_tmp.Phone = row[1]
		}
		if row[12] == "" {
			if row[1] != "" && row[1] != "0" {
				fs_tmp.Phone = row[1]
			} else if row[2] != "" && row[2] != "0"{
				fs_tmp.Phone = row[2]
			}
		}
		if fs_tmp.Phone != "" {
			fs_tmp.Start_time = row[5]
			fs_tmp.End_time = row[6]
			fs_tmp.Talk_time = row[8]
			fs_tmp.Status = fs_status_map[row[9]].Status_CN
			fs_info = append(fs_info,fs_tmp)
			fmt.Println(fs_tmp)
		}
	}
	//写入数据库
	for _, fs :=range fs_info {
		db.Pgsql_fs_info_insert(fs)
	}

}

const (
	pgname = "postgres"
	pgpw   = "contrail123"
	pghost = "49.235.175.29"
	dbport = 5432
	dbname = "campus3"
	ftable = "user_data"
)

func originate_phone(client *Client)  {
	for _, phone := range user_phone {
		cli := "originate user/" + phone + " &"
		client.BgApi(cli) // 非阻塞 拨打1004 和 1002 电话
	}
}

func init()  {
	db.PgsqlOpen(pghost, pgname, pgpw, dbname, dbport)
	for _, status := range fs_status {
		fs_status_map[status.Status_EN] = status
	}
}

func flaginit(path, filename *string) {
	//命令行是：test.exe -u root -p root123 -h localhost -port 8080
	//var path, filename string
	flag.StringVar(path, "path", "./", "path")
	flag.StringVar(filename, "filename", "test.csv", "file name")

	flag.Parse() //解析注册的flag，必须

	return
}

func main() {
	// Boost it as much as it can go ...
	// We don't need this since Go 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())
	var path, filename string
	flaginit(&path, &filename)
	client, err := NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!
	go client.Handle()

	client.Send("events json ALL")

	//client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))
	read_phone_csv(path, filename)
	originate_phone(client)
	//client.BgApi("originate user/1002 &") // 非阻塞 拨打1004 和 1002 电话
	//client.BgApi("originate user/1004 &")
	//read_master_csv()
	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				Error("Error while reading Freeswitch message: %s", err)
			}

			break
		}

		fmt.Println("Got new message: %s", msg)
	}
	defer db.PgsqlClose()
}

//func main()  {
//	read_master_csv()
//}