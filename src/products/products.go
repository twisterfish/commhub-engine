package products

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for product inventory data
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

//////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
//////////////////////////////////////////////////////////////////////////////////////////

func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "CreateProduct":
		return CreateProduct(payload)
	case "SetProductCount":
		return SetProductCount(payload)
	case "GetProductByID":
		return GetProductByID(payload)
	case "GetProductBySKU":
		return GetProductBySKU(payload)
	case "GetProductByUPC":
		return GetProductByUPC(payload)
	default:
		return "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Enter a new product in the database
/*
{"api_token": "37b176076fc74698be5aed02f74cbf15","signal":"products","action":"CreateProduct","end_user_token":"1","workspace_token":"1","vendor_id":"1","product_unit_sold_id":"1","upc":"983459345934345","sku":"234567Y","description":"cool ass widget"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func CreateProduct(payload string) string {

	type NewProduct struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EUTok    string `json:"end_user_token"`
		WSTok    string `json:"workspace_token"`
		VendID   string `json:"vendor_id"`
		PusID    string `json:"product_unit_sold_id"`
		PrdUPC   string `json:"upc"`
		PrdSKU   string `json:"sku"`
		PrdDesc  string `json:"description"`
	}

	// used for parsing request
	var input NewProduct
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type ProdRecord struct {
		PrdStatus string `json:"out_status"`
		NewProdID string `json:"out_new_product_id"`
	}
	var output ProdRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.create_product( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.VendID) + ", \"" + strings.TrimSpace(input.PusID) + "\", " + strings.TrimSpace(input.PrdUPC) + ", \"" + strings.TrimSpace(input.PrdSKU) + "\", \"" + strings.TrimSpace(input.PrdDesc) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.PrdStatus, &output.NewProdID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.PrdStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.PrdStatus) + "\"}"
			}

			rows.WriteString("{\"out_status\":\"" + output.PrdStatus + "\",")
			rows.WriteString("\"out_new_product_id\":\"" + output.NewProdID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Setting a cost on an item - just insert new price - it will be pulled off the stack
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"products","action":"SetProductCount","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1","total_quantity_on_hand":"10"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func SetProductCount(payload string) string {

	type NewProductQty struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EUTok    string `json:"end_user_token"`
		WSTok    string `json:"workspace_token"`
		PrdID    string `json:"product_id"`
		TqOH     string `json:"total_quantity_on_hand"`
	}

	var input NewProductQty
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}
	type ProdRecord struct {
		PrdStatus string `json:"out_status"`
	}
	var output ProdRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_product_count( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\"," + strings.TrimSpace(input.PrdID) + "," + strings.TrimSpace(input.TqOH) + ")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.PrdStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.PrdStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.PrdStatus) + "\"}"
			}

			rows.WriteString("{\"out_status\":\"" + output.PrdStatus + "\"},")
			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of products by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"products","action":"GetProductByID","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","product_id":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetProductByID(payload string) string {

	type ProductRequest struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EUTok    string `json:"end_user_token"`
		WSTok    string `json:"workspace_token"`
		PrdID    string `json:"product_id"`
	}
	// used for parsing request
	var input ProductRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Products struct {
		PrID   string `json:"product_id"`
		VenID  string `json:"vendor_id"`
		PuUSID string `json:"product_unit_sold_id"`
		PUPC   string `json:"upc"`
		PSKU   string `json:"sku"`
		PDesc  string `json:"description"`
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

	query := "CALL commhub_junction.get_product_by_id(\"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\"," + strings.TrimSpace(input.PrdID) + ")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Products
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.PrID, &tag.VenID, &tag.PuUSID, &tag.PUPC, &tag.PSKU, &tag.PDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"product_id\":\"" + tag.PrID + "\",")
			rows.WriteString("\"vendor_id\":\"" + tag.VenID + "\",")
			rows.WriteString("\"product_unit_sold_id\":\"" + tag.PuUSID + "\",")
			rows.WriteString("\"upc\":\"" + tag.PUPC + "\",")
			rows.WriteString("\"sku\":\"" + tag.PSKU + "\",")
			rows.WriteString("\"description\":\"" + tag.PDesc + "\"},")

			item_count++
		}
	}

	//prebuf.WriteString("{\"token\":\"" + input.Token + "\",\"signal\":\"products\",\"action\":\"GetProductByID\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	//prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of products by SKU
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"products","action":"GetProductBySKU","sku":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetProductBySKU(payload string) string {

	type ProductRequest struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EUTok    string `json:"end_user_token"`
		WSTok    string `json:"workspace_token"`
		PrdSKU   string `json:"sku"`
	}
	// used for parsing request
	var input ProductRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Products struct {
		PrID   string `json:"product_id"`
		VenID  string `json:"vendor_id"`
		PuUSID string `json:"product_unit_sold_id"`
		PUPC   string `json:"upc"`
		PSKU   string `json:"sku"`
		PDesc  string `json:"description"`
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
	query := "CALL commhub_junction.get_product_by_sku(\"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", \"" + strings.TrimSpace(input.PrdSKU) + "\")"
	results, err := db.Query(query)
	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
		var tag Products
		for results.Next() {
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.PrID, &tag.VenID, &tag.PuUSID, &tag.PUPC, &tag.PSKU, &tag.PDesc)
			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"product_id\":\"" + tag.PrID + "\",")
			rows.WriteString("\"vendor_id\":\"" + tag.VenID + "\",")
			rows.WriteString("\"product_unit_sold_id\":\"" + tag.PuUSID + "\",")
			rows.WriteString("\"upc\":\"" + tag.PUPC + "\",")
			rows.WriteString("\"sku\":\"" + tag.PSKU + "\",")
			rows.WriteString("\"description\":\"" + tag.PDesc + "\"},")

			item_count++
		}
		if strings.TrimSpace(tag.PrID) == "0" {
			item_count = 0
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of workorders by creator and assignee IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"products","action":"GetProductByUPC","upc":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetProductByUPC(payload string) string {

	type ProductRequest struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EUTok    string `json:"end_user_token"`
		WSTok    string `json:"workspace_token"`
		PrdUPC   string `json:"upc"`
	}
	// used for parsing request
	var input ProductRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type Products struct {
		PrID   string `json:"product_id"`
		VenID  string `json:"vendor_id"`
		PuUSID string `json:"product_unit_sold_id"`
		PUPC   string `json:"upc"`
		PSKU   string `json:"sku"`
		PDesc  string `json:"description"`
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

	query := "CALL commhub_junction.get_product_by_upc(\"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", \"" + strings.TrimSpace(input.PrdUPC) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
		var tag Products
		for results.Next() {

			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.PrID, &tag.VenID, &tag.PuUSID, &tag.PUPC, &tag.PSKU, &tag.PDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"product_id\":\"" + tag.PrID + "\",")
			rows.WriteString("\"vendor_id\":\"" + tag.VenID + "\",")
			rows.WriteString("\"product_unit_sold_id\":\"" + tag.PuUSID + "\",")
			rows.WriteString("\"upc\":\"" + tag.PUPC + "\",")
			rows.WriteString("\"sku\":\"" + tag.PSKU + "\",")
			rows.WriteString("\"description\":\"" + tag.PDesc + "\"},")

			item_count++
		}
		if strings.TrimSpace(tag.PrID) == "0" {
			item_count = 0
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}
