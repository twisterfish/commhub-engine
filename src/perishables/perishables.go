package perishables

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for perishable items data
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// {"token": "1.1.0","signal": "perishables","action": "GetAllPerishableTypes"}
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(signal string, action string, payload string) string {
	switch action {
	case "GetAllPerishableTypes":
		return GetAllPerishableTypes(action, payload)
	default:
		return "{\"token\": \"1.1.0\",\"signal\": \"error\",\"action\": \"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
//
// ////////////////////////////////////////////////////////////////////////////////////////
func GetAllPerishableTypes(action string, payload string) string {

	type PerishablesType struct {
		ID   string `json:"perishables_type_id"`
		Desc string `json:"description"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\": \"1.1.0\",\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"

	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	//results, err := db.Query("SELECT perishables_type_id, description, last_updated FROM commhub_junction.perishables_type;")
	results, err := db.Query("SELECT perishables_type_id, description FROM commhub_junction.perishables_type")
	if err != nil {
		defer db.Close()
		return "{\"token\": \"1.1.0\",\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag PerishablesType
			// for each row, scan the result into our struct
			//err = results.Scan(&tag.ID, &tag.Desc, &tag.TStamp)
			err = results.Scan(&tag.ID, &tag.Desc)

			if err != nil {
				defer db.Close()
				return "{\"token\": \"1.1.0\",\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"perishables_type_id\":\"" + tag.ID + "\",")
			rows.WriteString("\"description\":\"" + tag.Desc + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"ver\":\"1.10\",\"signal\":\"perishables\",\"action\":\"GetAllPerishableTypes\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	defer results.Close()
	// json.NewEncoder(retString).Encode(buf.String())

	return prebuf.String()

}
