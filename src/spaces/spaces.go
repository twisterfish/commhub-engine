package spaces

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for product inventory data
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
)

//////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
//////////////////////////////////////////////////////////////////////////////////////////

func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "CreateRealSpace":
		return CreateRealSpace(payload)
	case "CreateRealSpaceType":
		return CreateRealSpaceType(payload)
	case "GetRealSpaceByID":
		return GetRealSpaceByID(payload)
	case "GetAllRealSpaces":
		return GetAllRealSpaces(payload)
	case "GetAssetsForRealSpace":
		return GetAssetsForRealSpace(request, payload)
	default:
		return "{\"token\":\"invalid\",\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Internal to this package only
//////////////////////////////////////////////////////////////////////////////////////////

func getRealSpaceAssetSignedURL(assetName string, request *events.APIGatewayProxyRequest, ContType string) string {

	// SignBytes is a function for implementing custom signing.
	// Since our application is running on Google App Engine, we can use appengine's internal signing function:

	ctx := appengine.NewContext(request)
	acc, _ := appengine.ServiceAccount(ctx)

	url, err := storage.SignedURL("junction-signals.appspot.com/spaces", assetName,
		&storage.SignedURLOptions{
			GoogleAccessID: acc,
			SignBytes: func(b []byte) ([]byte, error) {
				_, signedBytes, err := appengine.SignBytes(ctx, b)
				return signedBytes, err
			},
			Method:      "GET",
			Expires:     time.Now().Add(24 * time.Hour),
			ContentType: ContType,
		})

	if err != nil {
		return err.Error()
	}

	return url

}

//////////////////////////////////////////////////////////////////////////////////////////
// Enter a new property in the database
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"CreateRealSpace","real_property_type_id":"1","real_space_inventory_status_id":"1","space_length":"11.0","space_width":"10.0","space_height":"12.0","description":"Really cool modern Hipster suite"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func CreateRealSpace(payload string) string {

	var prebuf bytes.Buffer

	type NewSpace struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		RPrdTID string `json:"real_property_id"`
		RPrNM   string `json:"real_space_inventory_status_id"`
		RPLen   string `json:"space_length"`
		RPWid   string `json:"space_width"`
		RPHght  string `json:"space_height"`
		RPDesc  string `json:"description"`
	}

	var npro NewSpace
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

	query := "INSERT INTO commhub_junction.real_space_inventory ( real_property_id, real_space_inventory_type_id, real_space_inventory_type_id, space_length, space_width, space_height, description ) VALUES(" +
		strings.TrimSpace(npro.RPrdTID) + "," +
		strings.TrimSpace(npro.RPrNM) + "," +
		strings.TrimSpace(npro.RPLen) + "," +
		strings.TrimSpace(npro.RPWid) + "," +
		strings.TrimSpace(npro.RPHght) + "," +
		"\"" + strings.TrimSpace(npro.RPDesc) + "\"," +
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
	prebuf.WriteString("{\"token\":\"" + npro.Token + "\",\"signal\":\"products\",\"action\":\"CreateRealSpace\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Creating a space type such as bedroom. conference room, dining room, etc.
/*
{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"CreateRealSpaceType","description":"Big Mansion"}
*/
//////////////////////////////////////////////////////////////////////////////////////////

func CreateRealSpaceType(payload string) string {

	var prebuf bytes.Buffer

	type NewSpaceType struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		PrdDesc string `json:"description"`
	}

	var npro NewSpaceType
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
	prebuf.WriteString("{\"token\":\"" + npro.Token + "\",\"signal\":\"products\",\"action\":\"CreateRealSpaceType\",\"status\":\"success\",\"last_insert_id\":\"" + str + "\"}")

	defer db.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull assets associated with this space
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"GetAssetsForRealSpace","real_property_id":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetAssetsForRealSpace(request *events.APIGatewayProxyRequest, payload string) string {

	type ProductAssetRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		RPoID  string `json:"real_property_id"`
	}
	var paR ProductAssetRequest

	type ProductAssetList struct {
		RPoID  string `json:"real_property_id"`
		AsNM   string `json:"asset_name"`
		AsCnT  string `json:"content_type"`
		AsSZE  string `json:"asset_size_bytes"`
		AsSZH  string `json:"asset_size_height"`
		AsSZW  string `json:"asset_size_width"`
		AsDesc string `json:"asset_description"`
	}
	var paL ProductAssetList

	pByte := []byte(payload)
	if err := json.Unmarshal(pByte, &paR); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + paR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()
	defer db.Close()
	// if there is an error opening the connection, handle it
	if err != nil {
		return "{\"token\":\"" + paR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	item_count = 0
	query := "SELECT real_property_id, asset_name, content_type, asset_size_bytes, asset_size_height, asset_size_width, asset_description FROM commhub_junction.real_property_asset WHERE real_property_id = " + paR.RPoID
	results, err := db.Query(query)

	if err != nil {
		return "{\"token\":\"" + paR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our tal struct
			err = results.Scan(&paL.RPoID, &paL.AsNM, &paL.AsCnT, &paL.AsSZE, &paL.AsSZH, &paL.AsSZW, &paL.AsDesc)

			if err != nil {
				//panic(err.Error()) // proper error handling instead of panic in your app
				defer db.Close()
				defer results.Close()
				return "{\"token\":\"" + paR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tal.ID)
			rows.WriteString("{\"real_property_id\":\"" + paL.RPoID + "\",")
			rows.WriteString("\"asset_name\":\"" + paL.AsNM + "\",")
			rows.WriteString("\"content_type\":\"" + paL.AsCnT + "\",")
			rows.WriteString("\"asset_size_bytes\":\"" + paL.AsSZE + "\",")
			rows.WriteString("\"asset_size_height\":\"" + paL.AsSZH + "\",")
			rows.WriteString("\"asset_size_width\":\"" + paL.AsSZW + "\",")

			rows.WriteString("\"signed_url\":\"" + getRealSpaceAssetSignedURL(paL.AsNM, request, paL.AsCnT) + "\",")

			rows.WriteString("\"asset_description\":\"" + paL.AsDesc + "\"},")
			item_count++
		}
	}

	prebuf.WriteString("{\"token\":\"" + paR.Token + "\",\"signal\":\"products\",\"action\":\"GetAssetsForProduct\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of spaces by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"GetRealSpaceByID","real_property_id":"1"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetRealSpaceByID(payload string) string {

	type SpaceRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		RPoID  string `json:"real_property_id"`
	}
	// used for parsing request
	var prdR SpaceRequest
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

	query := "SELECT real_property_id, real_property_type_id, real_property_name, address1, address2, city, province_state, zip_postal_code, latitude, longitude FROM commhub_junction.real_property WHERE real_property_id = " + prdR.RPoID
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Properties
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.RPrdID, &tag.RPrdTID, &tag.RPrNM, &tag.RPAdd1, &tag.RPAdd2, &tag.RPCity, &tag.RPPrSt, &tag.RPPZip, &tag.RPPLat, &tag.RPPLong)

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
			rows.WriteString("\"latitude\":\"" + tag.RPPLat + "\",")
			rows.WriteString("\"longitude\":\"" + tag.RPPLong + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"token\":\"" + prdR.Token + "\",\"signal\":\"properties\",\"action\":\"GetRealSpaceByID\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of products by IDs
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"properties","action":"GetAllRealSpaces"}
//////////////////////////////////////////////////////////////////////////////////////////

func GetAllRealSpaces(payload string) string {

	type SpaceRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
	}
	// used for parsing request
	var prdR SpaceRequest
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

	query := "SELECT real_property_id, real_property_type_id, real_property_name, address1, address2, city, province_state, zip_postal_code, latitude, longitude FROM commhub_junction.real_property"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"token\":\"" + prdR.Token + "\",\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag Properties
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.RPrdID, &tag.RPrdTID, &tag.RPrNM, &tag.RPAdd1, &tag.RPAdd2, &tag.RPCity, &tag.RPPrSt, &tag.RPPZip, &tag.RPPLat, &tag.RPPLong)

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
			rows.WriteString("\"latitude\":\"" + tag.RPPLat + "\",")
			rows.WriteString("\"longitude\":\"" + tag.RPPLong + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"token\":\"" + prdR.Token + "\",\"signal\":\"properties\",\"action\":\"GetAllRealSpaces\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}
