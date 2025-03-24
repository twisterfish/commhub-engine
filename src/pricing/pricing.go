package pricing

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for perishable items data
//////////////////////////////////////////////////////////////////////////////////////////

import (
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
	case "SetProductPrice":
		return SetProductPrice(payload)
	case "SetProductCost":
		return SetProductCost(payload)
	case "GetProductPrice":
		return GetProductPrice(payload)
	case "GetProductCost":
		return GetProductCost(payload)
	default:
		return "{\"token\":\"invalid\",\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Setting a price on an item - just insert new price - it will be pulled off the stack
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"pricing","action":"SetProductPrice","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1","price":"104.55","effective_date":"91112222"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func SetProductPrice(payload string) string {

	var prebuf bytes.Buffer

	type NewProductPrice struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EUToken     string `json:"end_user_token"`
		WrkSpcToken string `json:"workspace_token"`
		PrdID       string `json:"product_id"`
		PPrc        string `json:"price"`
		PrEfD       string `json:"effective_date"`
	}

	var input NewProductPrice
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		defer db.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	query := "CALL commhub_junction.set_product_price(\"" + strings.TrimSpace(input.EUToken) + "\",\"" + strings.TrimSpace(input.WrkSpcToken) + "\"," + strings.TrimSpace(input.PrdID) + "," + strings.TrimSpace(input.PPrc) + "," + strings.TrimSpace(input.PrEfD) + ");"

	result, err := db.Exec(query) // if there is an error, handle it
	if err != nil {
		defer db.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	last_insert_id, err2 := result.LastInsertId() // if there is an error getting last insert ID, let them know
	if err2 != nil {
		defer db.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err2.Error() + "\"}"
	}

	str := strconv.FormatInt(last_insert_id, 10)
	prebuf.WriteString("{\"token\":\"" + input.Token + "\",\"signal\":\"products\",\"action\":\"SetProductPrice\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Setting a cost on an item - just insert new price - it will be pulled off the stack
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"pricing","action":"SetProductCost","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1","cost":"103.00","effective_date":"1112222"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func SetProductCost(payload string) string {

	var prebuf bytes.Buffer

	type NewProductCost struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EUToken     string `json:"end_user_token"`
		WrkSpcToken string `json:"workspace_token"`
		PrdID       string `json:"product_id"`
		PrdCost     string `json:"cost"`
		PrEfD       string `json:"effective_date"`
	}

	var input NewProductCost
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		defer db.Close()
		return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	query := "CALL commhub_junction.set_product_cost(\"" + strings.TrimSpace(input.EUToken) + "\",\"" + strings.TrimSpace(input.WrkSpcToken) + "\"," + strings.TrimSpace(input.PrdID) + "," + strings.TrimSpace(input.PrdCost) + "," + strings.TrimSpace(input.PrEfD) + ");"

	results, err := db.Exec(query) // if there is an error inserting, handle it
	if err != nil {
		defer db.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	last_insert_id, err2 := results.LastInsertId() // if there is an error getting last insert ID, let them know
	if err2 != nil {
		defer db.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err2.Error() + "\"}"
	}

	str := strconv.FormatInt(last_insert_id, 10)
	prebuf.WriteString("{\"token\":\"" + input.Token + "\",\"signal\":\"products\",\"action\":\"SetProductCost\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a product's cost by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"pricing","action":"GetProductCost","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1","effective_date":"1112222"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetProductCost(payload string) string {

	type ProductRequest struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EUToken     string `json:"end_user_token"`
		WrkSpcToken string `json:"workspace_token"`
		PrdID       string `json:"product_id"`
		PrEfD       string `json:"effective_date"`
	}
	// used for parsing request
	var input ProductRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Products struct {
		PrID  string `json:"product_id"`
		PrCst string `json:"cost"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.get_product_cost(\"" + strings.TrimSpace(input.EUToken) + "\",\"" + strings.TrimSpace(input.WrkSpcToken) + "\"," + strings.TrimSpace(input.PrdID) + "," + strings.TrimSpace(input.PrEfD) + ");"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Products
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.PrID, &tag.PrCst)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"product_id\":\"" + tag.PrID + "\",")
			rows.WriteString("\"product_cost\":\"" + tag.PrCst + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a product's price by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"pricing","action":"GetProductPrice","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1","effective_date":"1112222"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetProductPrice(payload string) string {

	type ProductRequest struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EUToken     string `json:"end_user_token"`
		WrkSpcToken string `json:"workspace_token"`
		PrdID       string `json:"product_id"`
		PrEfD       string `json:"effective_date"`
	}
	// used for parsing request
	var input ProductRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Products struct {
		PrID   string `json:"product_id"`
		PrPrce string `json:"price"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.get_product_price(\"" + strings.TrimSpace(input.EUToken) + "\",\"" + strings.TrimSpace(input.WrkSpcToken) + "\"," + strings.TrimSpace(input.PrdID) + "," + strings.TrimSpace(input.PrEfD) + ");"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Products
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.PrID, &tag.PrPrce)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"product_id\":\"" + tag.PrID + "\",")
			rows.WriteString("\"price\":\"" + tag.PrPrce + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}
