package conn

import (
	//"comm"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type tHost struct {
	id       int
	H_host   string
	H_status string
}

type tJudge struct {
	id     int
	H_host string
	Tag    string
	H_date time.Time
}
type tSvc_name struct {
	id    int
	Sname string
	Stag  int
	Sbool int
}

var db *sql.DB

//var log comm.ConsoleLogger

// func InitDB() (err error) {
// 	db, err = sql.Open("mysql", "root:123456@tcp(mysql-master:3306)/auto_deploy")
// 	//db,err = sql.Open("mysql","root:123456@tcp(172.27.95.86:3306)/auto_deploy")

// 	if err != nil {
// 		fmt.Println("mysql connected failed!")
// 		return
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		return
// 	}
// 	return

// }

// InitDB initializes the database connection
func InitDB() (err error) {
	if db == nil {
		db, err = sql.Open("mysql", "root:123456@tcp(10.43.26.206:3306)/auto_deploy")
		if err != nil {
			log.Error("failed to open database: %v", err)
			return
		}

		// Set connection pool settings
		db.SetMaxOpenConns(25)   // 设置最大打开连接数
		db.SetMaxIdleConns(25)   // 设置最大闲置连接数
		db.SetConnMaxLifetime(0) // 连接可重用的最大时间

		// Test the connection
		if err = db.Ping(); err != nil {
			log.Error("failed to connect to database: %v", err)
			return
		}
		log.Debug("database connection established")
	}
	return
}

// 输入host信息
func InsertHostData(H_host, H_hostname, H_core, H_free, H_disk string) interface{} {
	InitDB()

	sqlStr := "insert into host(H_host, H_hostname, H_core, H_free, H_disk, H_status) values(?,?,?,?,?,0);"

	ret, err := db.Exec(sqlStr, H_host, H_hostname, H_core, H_free, H_disk)
	if err != nil {
		fmt.Printf("exec err:", err)
	}
	theId, err := ret.LastInsertId()
	if err != nil {
		fmt.Println("last id failed")

	}

	fmt.Printf("insert success:%v", ret)
	return theId
}

// 删除host信息
func DelHostData(H_id string) int64 {
	InitDB()
	sqlStr := "delete from host where id =?"
	ret, err := db.Exec(sqlStr, H_id)
	if err != nil {
		fmt.Printf("deltel host error %v\n", err)
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf(" get Rowsaffected is err%v\n", err)
	}

	return n

}

type svc_name struct {
	Sid    string
	Sname  string
	S_tag  string
	Sbool  string
	Simage string
	Stitle string
}
type Count struct {
	cNums int
}

// 取值svc_name表内容并组装
func GetDataAll_svc_name() (interface{}, interface{}) {
	InitDB()
	stmt, err := db.Prepare("select id,Sname,S_tag,bool,image,title  from svc_name")
	if err != nil {
		fmt.Printf("prepare sql准备err:%v\n", err)
	}
	defer stmt.Close()
	var ResultNums Count
	err = db.QueryRow("select count(0) from svc_name").Scan(&ResultNums.cNums)
	//获取结果后调用scan方法，会自动释放数据库链接，注意观察
	if err != nil {
		fmt.Printf("获取目标结果行数出错：", err)
	}
	fmt.Printf("获取的结果行数是：%v\n", ResultNums.cNums)
	//获取到结果的行数，然后赋值给切片的长度
	rows, err := stmt.Query()
	fmt.Printf("执行查询语句后返回的结果:%v 结果行数为：\n", rows)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
	}
	m1 := make([]map[string]string, ResultNums.cNums)
	//ResultNums.cNums是上次获取的结果行数，填到这里，作为切片的长度，实际做了两次查询，一次是获取长度sql，另一次是获取数据sql
	counter := 0
	for rows.Next() {
		var Name svc_name
		err := rows.Scan(&Name.Sid, &Name.Sname, &Name.S_tag, &Name.Sbool, &Name.Simage, &Name.Stitle)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		fmt.Printf("查询到数据库的所有值sid:%v name:%v tag:%v bool:%v image:%v title:%v \n", Name.Sid, Name.Sname, Name.S_tag, Name.Sbool, Name.Simage)

		m1[counter] = make(map[string]string, 100)
		m1[counter]["ID"] = Name.Sid
		m1[counter]["name"] = Name.Sname
		m1[counter]["tag"] = Name.S_tag
		m1[counter]["bool"] = Name.Sbool
		m1[counter]["image"] = Name.Simage
		m1[counter]["title"] = Name.Stitle
		counter++
	}
	//m1=append(m1,`"category":{"relation":"关系型","reno":"非关系型"}`)
	fmt.Printf("===========================================map m1:%v\n", m1)
	m2 := make(map[string]string, ResultNums.cNums)
	m2["relation"] = "关系型数据库"
	m2["relationno"] = "非关系型数据库"
	m2["web"] = "web服务"
	m2["register"] = "注册及队列服务"
	m3 := make([]interface{}, 2)
	m3 = append(m3, m1, m2)
	fmt.Printf("得到m3的切片数据为：%v\n", m3)

	return m1, m2

}

// 查询单条数据
func GetTagOne(Rid int) string {
	InitDB()
	var ret_Judge tJudge
	sqlStr := `select Tag from judge where id=?`
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	current_getOne := db.QueryRow(sqlStr, Rid)
	err := current_getOne.Scan(&ret_Judge.Tag)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
	}
	return ret_Judge.Tag
}

//查询数据库根据服务名称，得到配置文件名称及路径

func GetServerConfigName(Sname string) string {
	InitDB()
	type name struct {
		Sname string
	}
	sqlstr := `select params from svc_name where Sname=?`
	//stmt,err := db.Prepare("select Sname from ?")
	//if err != nil {
	//	fmt.Printf("prepare sql准备err:%v\n",err)
	//}
	//defer stmt.Close()
	fmt.Printf("函数传入的参数是%v\n", Sname)
	var Name name
	ret := db.QueryRow(sqlstr, Sname)
	err := ret.Scan(&Name.Sname)
	if err != nil {
		fmt.Printf("sql返回结果报错；%v\n", err)
		return "MysqlErr"
	}

	//GetRow := db.QueryRow(stmt,Sname)
	return Name.Sname
}

func GetDownLoadURL(SvcName string) string {
	InitDB()
	type name struct {
		Sname string
	}
	sqlstr := `select downloadURL from svc_name where Sname=?`
	fmt.Printf("函数传入的参数是%v\n", SvcName)
	var Name name
	ret := db.QueryRow(sqlstr, SvcName)
	err := ret.Scan(&Name.Sname)
	if err != nil {
		fmt.Printf("sql返回结果报错；%v\n", err)
		return "MysqlErr"
	}

	//GetRow := db.QueryRow(stmt,Sname)
	return Name.Sname
}

func GetHostContext(ID string) map[string]string {
	InitDB() // Ensure DB is initialized

	type HostContextStru struct {
		HHost   string
		HStatus string
	}

	sqlstr := `SELECT H_host, H_status FROM host WHERE id = ?`
	var host HostContextStru

	// Execute SQL query
	err := db.QueryRow(sqlstr, ID).Scan(&host.HHost, &host.HStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("no rows found for id: %s", ID)
			return nil
		}
		log.Error("error querying host context: %v", err)
		return nil
	}

	HostContext := map[string]string{
		"H_host":   host.HHost,
		"H_status": host.HStatus,
	}

	// Debug output
	//log.Debug("getHostContext: retrieved data %v", HostContext)
	return HostContext
}

func GetHostContextAll() []map[string]string {
	InitDB()
	type HostContextStru struct {
		H_id       string
		H_host     string
		H_status   string
		H_core     string
		H_free     string
		H_disk     string
		H_hostname string
	}
	sqlstr := `select id,H_host,H_hostname,H_core,H_free,H_disk,H_status from host`
	rows, err := db.Query(sqlstr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()
	// 循环读取结果集中的数据
	//HostContextAll := make([]map[string]string, 100)
	HostContextAll := make([]map[string]string, 0)
	for rows.Next() {
		var H HostContextStru
		err := rows.Scan(&H.H_id, &H.H_host, &H.H_hostname, &H.H_core, &H.H_free, &H.H_disk, &H.H_status)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		HostContext := make(map[string]string)
		HostContext["H_id"] = H.H_id
		HostContext["H_host"] = H.H_host
		HostContext["H_hostname"] = H.H_hostname

		HostContext["H_core"] = H.H_core
		HostContext["H_free"] = H.H_free
		HostContext["H_disk"] = H.H_disk
		HostContext["H_status"] = H.H_status

		HostContextAll = append(HostContextAll, HostContext)
	}
	return HostContextAll

}
