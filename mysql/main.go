package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

func main() {
	con, err := sql.Open("mysql", "root:123456@/test")
	if err != nil {
		panic(err)
	}
	con.SetMaxOpenConns(100)
	con.SetConnMaxLifetime(time.Second * 3)
	defer con.Close()
	_, err = con.Exec("truncate um_address")
	if err != nil {
		panic(err)
	}

	group := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		group.Add(1)
		i := i
		go func() {
			defer group.Done()
			begin, err := con.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
			if err != nil {
				panic(err)
			}
			begin.Exec("explain UPDATE um_address SET address_number= 0 WHERE ada = ? AND address_type = ? AND address_number = 1 AND status = 1",
				"8216248", "X"+strconv.Itoa(i))
			exec, err := begin.Exec("insert into um_address (ada, name, address_type, address_number, province, city,district, town, post_code, region_code, channel, district_code, address) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", "8216248", "測試名", "X"+strconv.Itoa(i), 2, "", "台北市", "大同區", "", "890",
				"130", "POS", "108", "PO9JCJQHZn0AH600UmfeuA==")
			if err != nil {
				begin.Rollback()
				panic(err)
			}
			fmt.Println(exec.LastInsertId())
			begin.Commit()
		}()
	}
	group.Wait()
}
