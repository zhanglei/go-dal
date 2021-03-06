package main

import (
	"fmt"
	"time"

	"github.com/antlinker/go-dal"

	_ "github.com/antlinker/go-dal/mysql"
)

type Student struct {
	ID       int64
	StuCode  string
	StuName  string
	Sex      int
	Age      int
	Birthday time.Time
	Memo     string
}

func main() {
	dal.RegisterProvider(dal.MYSQL,
		`{
		"datasource":"root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8",
		"maxopen":100,
		"maxidle":50,
		"print":true
	}`)
	// insert()
	// update()
	// delete()
	// transaction()
	list()
	// single()
	// insertManyData()
	// pager()
}

func getAddEntity() dal.TranEntity {
	stud := Student{
		StuCode:  "S001",
		StuName:  "Lyric",
		Sex:      1,
		Age:      25,
		Birthday: time.Now(),
		Memo:     "Message...",
	}
	return dal.NewTranAEntity("student", stud).Entity
}

func insert() {
	entity := getAddEntity()
	result := dal.Exec(entity)
	if err := result.Error; err != nil {
		panic(err)
	}
	fmt.Println("===> Student Insert:", result.Result)
}

func getUpdateEntity() dal.TranEntity {
	stud := map[string]interface{}{
		"StuName": "Elva",
		"Sex":     2,
		"Age":     26,
	}
	cond := dal.NewFieldsKvCondition(map[string]interface{}{"StuCode": "S001"}).Condition
	entity := dal.NewTranUEntity("student", stud, cond).Entity
	return entity
}

func update() {
	entity := getUpdateEntity()
	result := dal.Exec(entity)
	if err := result.Error; err != nil {
		panic(err)
	}
	fmt.Println("===> Student Update:", result.Result)
}

func delete() {
	cond := dal.NewFieldsKvCondition(Student{StuCode: "S001"}).Condition
	entity := dal.NewTranDEntity("student", cond).Entity
	result := dal.Exec(entity)
	if err := result.Error; err != nil {
		panic(err)
	}
	fmt.Println("===> Student Delete:", result.Result)
}

func single() {
	entity := dal.NewQueryEntity("student",
		dal.NewFieldsKvCondition(map[string]interface{}{"StuCode": "S001"}).Condition,
		"StuName", "StuCode", "Memo")().Entity

	data, err := dal.Single(entity)
	if err != nil {
		panic(err)
	}
	fmt.Println("===> Student Single:", data)
}

func list() {
	entity := dal.NewQueryEntity("student",
		dal.NewCondition("order by Id limit 10").Condition,
		"Id", "StuCode", "StuName")().Entity

	var stuData []Student
	err := dal.AssignList(entity, &stuData)
	if err != nil {
		panic(err)
	}
	for i, item := range stuData {
		fmt.Println(i+1, ",Student data:", item.ID, item.StuCode, item.StuName)
	}
}

func transaction() {
	var transEntity []dal.TranEntity
	transEntity = append(transEntity, getAddEntity())
	transEntity = append(transEntity, getUpdateEntity())
	result := dal.ExecTrans(transEntity)
	if err := result.Error; err != nil {
		panic(err)
	}
	fmt.Println("===> Execute Transaction:", result.Result)
}

func insertManyData() {
	var entities []dal.TranEntity
	for i := 0; i < 1000; i++ {
		var stu Student
		stu.StuCode = fmt.Sprintf("S-%d", i)
		stu.StuName = fmt.Sprintf("SName-%d", i)
		stu.Birthday = time.Now()
		entities = append(entities, dal.NewTranAEntity("student", stu).Entity)
	}
	result := dal.ExecTrans(entities)
	if err := result.Error; err != nil {
		panic(err)
	}
	fmt.Println("===> Insert data numbers:", result.Result)
}

func pager() {
	entity := dal.NewQueryPagerEntity("student",
		dal.NewCondition("where StuCode like ? order by ID", "S-%").Condition,
		dal.NewPagerParam(1, 20),
		"StuCode", "StuName", "Birthday").Entity
	result, err := dal.Pager(entity)
	if err != nil {
		panic(err)
	}
	fmt.Println("===> Query total:")
	fmt.Println(result.Total)
	fmt.Println("===> Query rows:")
	fmt.Println(result.Rows)
}
