package account

import (
	//"database/sql"
	"fmt"
	"log"
)

/********************************role configuration begin***********************************/
type RoleConf struct {
	RoleId uint16
	Att    int
	Def    int
}

var roleConfMgr = make(map[uint16]*RoleConf)

func InitRoleConf() {
	rows, err := confDB.Query("SELECT * FROM roleconf ORDER BY roleId ASC;")
	if err != nil {
		log.Panicf("InitRoleConf() failed")
	}

	if rows == nil {
		log.Panicf("roleconf is empty table")
	}
	defer rows.Close()

	//cols, err := rows.Columns()             // Get the column names; remember to check err
	//vals := make([]sql.RawBytes, len(cols)) // Allocate enough values
	//ints := make([]interface{}, len(cols))  // Make a slice of []interface{}
	//fmt.Printf("%v\n %v\n %v\n", cols, vals, ints)

	for rows.Next() {
		role := &RoleConf{}
		err = rows.Scan(&role.RoleId, &role.Att, &role.Def)
		if err != nil {
			log.Panicf("roleconf query error! %v", role)
		}
		roleConfMgr[role.RoleId] = role
		//fmt.Printf("%v\n", role)
	}
	err = rows.Err()
	if err != nil {
		log.Panicf("roleconf loop query error")
	}

	fmt.Println("InitRoleConf() ok")
}

/********************************role configuration end***********************************/
