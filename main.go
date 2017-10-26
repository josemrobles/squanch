package squanch

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

func Fetch (d,h,u,p,db,q string) (json string,err error) {
	c,err := connect(d,h,u,p,db)
	if err != nil {
		return "",err
	} else {
		j,err := getData(c,q) 
		if err != nil {
			return "",err
		} else {
			return j,nil
		}
	}
}

// Connects to the specified database using the driver (d), user (u)
// password (p), host(h), and database (db).
// @TODO - Unit test
func connect(d,h,u,p,db string) (c *sql.DB,err error) {
	if d == "mysql" {
		c, err = sql.Open(d, u+":"+p+"@tcp("+h+":3306)/"+db+"?allowOldPasswords=1")
	} else {
		//c, err = sql.Open(d, "postgres://"+u+":"+p+"@"+h+"/"+db+"?sslmode=verify-full")
		c, err = sql.Open(d, "postgres://"+u+":"+p+"@"+h+"/"+db)
	}
	
	if err != nil {
		return nil,err
	} else {
		//defer c.Close()
		return c,nil
	}

}

// Function used to get a bran fromthe database using the supplied brand name
// Function used to gat data from the connected database source. If connected
// and the query is successful, the output will be JSON and nilfor the err var 
// otherwise the result will be an empty json object and the err details/
// @TODO - Unit test
func getData(db *sql.DB, query string) (string, error) {

	// Execute the query
	records, err := db.Query(query)
	if err != nil {
		return "{}", err
	} else {
		
		
	defer records.Close()

	// Get the column names
	cols, err := records.Columns()
	if err != nil {
		return "{}", err
	}
	count := len(cols)
	td := make([]map[string]interface{}, 0)
	data := make([]interface{}, count)
	vals := make([]interface{}, count)
	
	for records.Next() {
		for i := 0; i < count; i++ {
			vals[i] = &data[i]
		}
		records.Scan(vals...)
		entry := make(map[string]interface{})
		for i, column := range cols {
			var v interface{}
			datum := data[i]
			b, ok := datum.([]byte)
			if ok {
				v = string(b)
			} else {
				v = datum
			}
			entry[column] = v
		}
		td = append(td, entry)
	}
	json, err := json.Marshal(td)
	if err != nil {
		return "{}", err
	}
	return string(json), nil
		
	}





}
