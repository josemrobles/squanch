package squanch

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-pg/pg"
	_ "github.com/go-sql-driver/mysql"
)



// Function used to get a bran fromthe database using the supplied brand name
// Function used to gat data from the connected database source. If connected
// and the query is successful, the output will be JSON and nilfor the err var 
// otherwise the result will be an empty json object and the err details/
// @TODO - Unit test
func getData(db *sql.DB, query string) (string, error) {

	// Execute the query
	records, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer records.Close()

	// Get the column names
	cols, err := records.Columns()
	if err != nil {
		return "", err
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
		for i, column := range columns {
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
		return "", err
	}
	return string(json), nil
}
