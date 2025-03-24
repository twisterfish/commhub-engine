package properties

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for product inventory data
//////////////////////////////////////////////////////////////////////////////////////////

import (

	//"google.golang.org/appengine"
	"bytes"
	"datastores"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

//////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
//////////////////////////////////////////////////////////////////////////////////////////

func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "CreateRealProperty":
		return CreateRealProperty(payload)
	case "CreateRealPropertyType":
		return CreateRealPropertyType(payload)
	case "GetRealPropertyByID":
		return GetRealPropertyByID(payload)
	case "GetAllRealProperties":
		return GetAllRealProperties(payload)
	default:
		return "{\"token\":\"invalid\",\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Enter a new property in the database
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"CreateRealProperty","real_property_type_id":"1","real_property_name":"The Big Fish Hotel","address1":"111 Richie Way","address2":"unit 3","city":"Miami","province_state":"OK","zip_postal_code":"33318","latitude":"26.108900","longitude":"-80.106735"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func CreateRealProperty(payload string) string {

	var prebuf bytes.Buffer

	type NewProperty struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		RPrdTID string `json:"real_property_type_id"`
		RPrNM   string `json:"real_property_name"`
		RPAdd1  string `json:"address1"`
		RPAdd2  string `json:"address2"`
		RPCity  string `json:"city"`
		RPPrSt  string `json:"province_state"`
		RPPZip  string `json:"zip_postal_code"`
		RPPLat  string `json:"latitude"`
		RPPLong string `json:"longitude"`
	}

	var npro NewProperty
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &npro); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	query := "INSERT INTO commhub_junction.real_property (real_property_type_id, real_property_name, address1, address2, city, province_state, zip_postal_code, latitude, longitude) VALUES(" +
		strings.TrimSpace(npro.RPrdTID) + "," +
		"\"" + strings.TrimSpace(npro.RPrNM) + "\"," +
		"\"" + strings.TrimSpace(npro.RPAdd1) + "\"," +
		"\"" + strings.TrimSpace(npro.RPAdd2) + "\"," +
		"\"" + strings.TrimSpace(npro.RPCity) + "\"," +
		"\"" + strings.TrimSpace(npro.RPPrSt) + "\"," +
		"\"" + strings.TrimSpace(npro.RPPZip) + "\"," +
		strings.TrimSpace(npro.RPPLat) + "," +
		strings.TrimSpace(npro.RPPLong) + ")"

	result, err := db.Exec(query) // if there is an error inserting, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	last_insert_id, err2 := result.LastInsertId() // if there is an error getting last insert ID, let them know
	if err2 != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err2.Error() + "\"}"
	}

	str := strconv.FormatInt(last_insert_id, 10)
	prebuf.WriteString("{\"token\":\"" + npro.Token + "\",\"signal\":\"products\",\"action\":\"CreateRealProperty\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Creating a property type such as plumbing, electrical, dry goods, etc.
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"CreateRealPropertyType","description":"Big Mansion"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func CreateRealPropertyType(payload string) string {

	var prebuf bytes.Buffer

	type NewPropertyType struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		PrdDesc string `json:"description"`
	}

	var npro NewPropertyType
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &npro); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	query := "INSERT INTO commhub_junction.real_property_type ( description ) VALUES(" +
		"\"" + strings.TrimSpace(npro.PrdDesc) + "\"" +
		")"

	result, err := db.Exec(query) // if there is an error inserting, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	last_insert_id, err2 := result.LastInsertId() // if there is an error getting last insert ID, let them know
	if err2 != nil {
		defer db.Close()
		return "{\"token\":\"" + npro.Token + "\",\"signal\":\"error\",\"action\": \"" + err2.Error() + "\"}"
	}

	str := strconv.FormatInt(last_insert_id, 10)
	prebuf.WriteString("{\"token\":\"" + npro.Token + "\",\"signal\":\"products\",\"action\":\"CreateRealPropertyType\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of products by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"GetRealPropertyByID","real_property_id":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetRealPropertyByID(payload string) string {

	type PropertyRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		RPoID  string `json:"real_property_id"`
	}
	// used for parsing request
	var prdR PropertyRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &prdR); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Properties struct {
		RPrdID  string `json:"real_property_id"`
		RPrdTID string `json:"real_property_type_id"`
		RPrNM   string `json:"real_property_name"`
		RPAdd1  string `json:"address1"`
		RPAdd2  string `json:"address2"`
		RPCity  string `json:"city"`
		RPPrSt  string `json:"province_state"`
		RPPZip  string `json:"zip_postal_code"`
		RPCNTRY string `json:"country_code"`
		RPPLat  string `json:"latitude"`
		RPPLong string `json:"longitude"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "SELECT real_property_id, real_property_type_id, IFNULL(real_property_name,''), IFNULL(address1,''), IFNULL(address2,''), IFNULL(city,''), IFNULL(province_state,''), IFNULL(zip_postal_code,''), IFNULL(country_code,''), latitude, longitude FROM commhub_junction.real_property WHERE real_property_id = " + prdR.RPoID
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Properties
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.RPrdID, &tag.RPrdTID, &tag.RPrNM, &tag.RPAdd1, &tag.RPAdd2, &tag.RPCity, &tag.RPPrSt, &tag.RPPZip, &tag.RPCNTRY, &tag.RPPLat, &tag.RPPLong)

			if err != nil {
				results.Close()
				return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"real_property_id\":\"" + tag.RPrdID + "\",")
			rows.WriteString("\"real_property_type_id\":\"" + tag.RPrdTID + "\",")
			rows.WriteString("\"real_property_name\":\"" + tag.RPrNM + "\",")
			rows.WriteString("\"address1\":\"" + tag.RPAdd1 + "\",")
			rows.WriteString("\"address2\":\"" + tag.RPAdd2 + "\",")
			rows.WriteString("\"city\":\"" + tag.RPCity + "\",")
			rows.WriteString("\"province_state\":\"" + tag.RPPrSt + "\",")
			rows.WriteString("\"zip_postal_code\":\"" + tag.RPPZip + "\",")
			rows.WriteString("\"country_code\":\"" + tag.RPCNTRY + "\",")
			rows.WriteString("\"latitude\":\"" + tag.RPPLat + "\",")
			rows.WriteString("\"longitude\":\"" + tag.RPPLong + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"token\":\"" + prdR.Token + "\",\"signal\":\"properties\",\"action\":\"GetRealPropertyByID\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of products by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"GetAllRealProperties"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetAllRealProperties(payload string) string {

	type PropertyRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
	}
	// used for parsing request
	var prdR PropertyRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &prdR); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Properties struct {
		RPrdID  string `json:"real_property_id"`
		RPrdTID string `json:"real_property_type_id"`
		RPrNM   string `json:"real_property_name"`
		RPAdd1  string `json:"address1"`
		RPAdd2  string `json:"address2"`
		RPCity  string `json:"city"`
		RPPrSt  string `json:"province_state"`
		RPPZip  string `json:"zip_postal_code"`
		RPCNTRY string `json:"country_code"`
		RPPLat  string `json:"latitude"`
		RPPLong string `json:"longitude"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "SELECT real_property_id, real_property_type_id, IFNULL(real_property_name,''), IFNULL(address1,''), IFNULL(address2,''), IFNULL(city,''), IFNULL(province_state,''), IFNULL(zip_postal_code,''), IFNULL(country_code,''), latitude, longitude FROM commhub_junction.real_property"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Properties
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.RPrdID, &tag.RPrdTID, &tag.RPrNM, &tag.RPAdd1, &tag.RPAdd2, &tag.RPCity, &tag.RPPrSt, &tag.RPPZip, &tag.RPCNTRY, &tag.RPPLat, &tag.RPPLong)

			if err != nil {
				results.Close()
				return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"real_property_id\":\"" + tag.RPrdID + "\",")
			rows.WriteString("\"real_property_type_id\":\"" + tag.RPrdTID + "\",")
			rows.WriteString("\"real_property_name\":\"" + tag.RPrNM + "\",")
			rows.WriteString("\"address1\":\"" + tag.RPAdd1 + "\",")
			rows.WriteString("\"address2\":\"" + tag.RPAdd2 + "\",")
			rows.WriteString("\"city\":\"" + tag.RPCity + "\",")
			rows.WriteString("\"province_state\":\"" + tag.RPPrSt + "\",")
			rows.WriteString("\"zip_postal_code\":\"" + tag.RPPZip + "\",")
			rows.WriteString("\"country_code\":\"" + tag.RPCNTRY + "\",")
			rows.WriteString("\"latitude\":\"" + tag.RPPLat + "\",")
			rows.WriteString("\"longitude\":\"" + tag.RPPLong + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"token\":\"" + prdR.Token + "\",\"signal\":\"properties\",\"action\":\"GetAllRealProperties\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}
