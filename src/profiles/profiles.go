package profiles

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for perishable items data
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"encoding/json"
	"guids"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "GetUserProfile":
		return GetUserProfile(payload)
	case "UpdateUserProfile":
		return UpdateUserProfile(payload)
	case "GenNewProfileToken":
		return GenNewProfileToken(payload)
	default:
		return "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a profile by token - NEED TO CHOP THIS LATER!!
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"profiles","action":"GetUserProfile","end_user_token":"8a985a9286ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetUserProfile(payload string) string {

	type ProfileRequest struct {
		ApiToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EuToken  string `json:"end_user_token"`
	}
	// used for parsing request
	var input ProfileRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type UserProfile struct {
		EUTok   string `json:"end_user_token"`
		EINTX   string `json:"ein_tax_id"`
		SSNTX   string `json:"ssn_tax_id"`
		LstNm   string `json:"last_name"`
		MidNM   string `json:"middle_name"`
		FTNM    string `json:"first_name"`
		ADD1    string `json:"address1"`
		ADD2    string `json:"address2"`
		City    string `json:"city"`
		PST     string `json:"province_state"`
		ZipPC   string `json:"zip_postal_code"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Country string `json:"country_code"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "SELECT ENDU.end_user_token, IFNULL(EUPRF.ein_tax_id,''), IFNULL(EUPRF.ssn_tax_id,''), IFNULL(EUPRF.last_name,''), IFNULL(EUPRF.middle_name,''), IFNULL(EUPRF.first_name,''), IFNULL(EUPRF.address1,''), IFNULL( EUPRF.address2,''), IFNULL(EUPRF.city,''), IFNULL(EUPRF.province_state,''), IFNULL(EUPRF.zip_postal_code,''), IFNULL(EUPRF.email,''), IFNULL(EUPRF.phone,''), IFNULL(EUPRF.country_code,'') FROM commhub_junction.end_user AS ENDU INNER JOIN commhub_junction.end_user_profile AS EUPRF ON EUPRF.end_user_id = ENDU.end_user_id WHERE ENDU.end_user_token = \"" + input.EuToken + "\""
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag UserProfile
			// for each record, scan the result into our  struct
			err = results.Scan(&tag.EUTok, &tag.EINTX, &tag.SSNTX, &tag.LstNm, &tag.MidNM, &tag.FTNM, &tag.ADD1, &tag.ADD2, &tag.City, &tag.PST, &tag.ZipPC, &tag.Email, &tag.Phone, &tag.Country)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"end_user_token\":\"" + tag.EUTok + "\",")
			rows.WriteString("\"ein_tax_id\":\"" + tag.EINTX + "\",")
			rows.WriteString("\"ssn_tax_id\":\"" + tag.SSNTX + "\",")
			rows.WriteString("\"last_name\":\"" + tag.LstNm + "\",")
			rows.WriteString("\"middle_name\":\"" + tag.MidNM + "\",")
			rows.WriteString("\"first_name\":\"" + tag.FTNM + "\",")
			rows.WriteString("\"address1\":\"" + tag.ADD1 + "\",")
			rows.WriteString("\"address2\":\"" + tag.ADD2 + "\",")
			rows.WriteString("\"city\":\"" + tag.City + "\",")
			rows.WriteString("\"province_state\":\"" + tag.PST + "\",")
			rows.WriteString("\"zip_postal_code\":\"" + tag.ZipPC + "\",")
			rows.WriteString("\"email\":\"" + tag.Email + "\",")
			rows.WriteString("\"phone\":\"" + tag.Phone + "\",")
			rows.WriteString("\"country_code\":\"" + tag.Country + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"profiles\",\"action\":\"GetUserProfile\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")

	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Generate a new profile token for the user for security purposes - they must be authorized to get a regisration token (for now)
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"profiles","action":"GenNewProfileToken","end_user_token":"98aa64c4c65d46c8bd948f7a0e50f278"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GenNewProfileToken(payload string) string {

	var prebuf bytes.Buffer

	type MemberInfo struct {
		Token     string `json:"api_token"`
		Signal    string `json:"signal"`
		Action    string `json:"action"`
		UserToken string `json:"end_user_token"`
	}

	var input MemberInfo
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it
	if err != nil {
		db.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	new_uuid := guids.GetGUID() // generate a new GUID

	query := "UPDATE commhub_junction.end_user SET end_user_token = \"" + strings.TrimSpace(new_uuid) + "\" WHERE end_user_token = \"" + strings.TrimSpace(input.UserToken) + "\""

	res, err2 := db.Exec(query) // if there is an error, handle it

	if err2 != nil {
		db.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {
		count, err3 := res.RowsAffected()

		if err3 != nil {
			db.Close()
			return "{\"signal\":\"error\",\"action\": \"" + err3.Error() + "\"}"
		} else {
			if count == 0 {
				prebuf.WriteString("{\"signal\":\"profiles\",\"action\":\"GenNewProfileToken\",\"status\":\"failed\",\"end_user_token:\"" + input.UserToken + "\"}")
			} else {
				prebuf.WriteString("{\"signal\":\"profiles\",\"action\":\"GenNewProfileToken\",\"status\":\"success\",\"end_user_token:\"" + new_uuid + "\"}")
			}
		}
	}

	defer db.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Enter a new profile in the database
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"profiles","action":"UpdateUserProfile","api_token":"37b176076fc74698be5aed02f74cbf15","emid":"edward.anderson@commhubstuff.com","pwid":"testpass","npwid":"SOMERANDOMSHIT","ein_tax_id":"321","ssn_tax_id":"123","last_name":"Public","middle_name":"Q","first_name":"Edward","address1":"910 Se 17th Street"","address2":"#420","city":"Fort Lauderdale","province_state":"OK","zip_postal_code":"33316","phone":"1232245945","email":"edward.anderson@commhubstuff.com","country_code":"US"}
// ////////////////////////////////////////////////////////////////////////////////////////
func UpdateUserProfile(payload string) string {

	type NewProfile struct {
		ApiToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		EMID     string `json:"emid"`
		PWID     string `json:"pwid"`
		NPWID    string `json:"npwid"`
		EINTX    string `json:"ein_tax_id"`
		SSNTX    string `json:"ssn_tax_id"`
		LstNm    string `json:"last_name"`
		MidNM    string `json:"middle_name"`
		FTNM     string `json:"first_name"`
		ADD1     string `json:"address1"`
		ADD2     string `json:"address2"`
		City     string `json:"city"`
		PST      string `json:"province_state"`
		ZipPC    string `json:"zip_postal_code"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Country  string `json:"country_code"`
	}

	type UsrRecord struct {
		UsrStatus string `json:"out_status"`
	}
	var output UsrRecord

	var input NewProfile
	var rows bytes.Buffer
	var item_count int
	var prebuf bytes.Buffer
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

	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.update_user_profile(" +
		"\"" + strings.TrimSpace(input.EMID) + "\"," +
		"\"" + strings.TrimSpace(input.PWID) + "\"," +
		"\"" + strings.TrimSpace(input.NPWID) + "\"," +
		"\"" + strings.TrimSpace(input.EINTX) + "\"," +
		"\"" + strings.TrimSpace(input.SSNTX) + "\"," +
		"\"" + strings.TrimSpace(input.LstNm) + "\"," +
		"\"" + strings.TrimSpace(input.MidNM) + "\"," +
		"\"" + strings.TrimSpace(input.FTNM) + "\"," +
		"\"" + strings.TrimSpace(input.ADD1) + "\"," +
		"\"" + strings.TrimSpace(input.ADD2) + "\"," +
		"\"" + strings.TrimSpace(input.City) + "\"," +
		"\"" + strings.TrimSpace(input.PST) + "\"," +
		"\"" + strings.TrimSpace(input.ZipPC) + "\"," +
		"\"" + strings.TrimSpace(input.Email) + "\"," +
		"\"" + strings.TrimSpace(input.Phone) + "\"," +
		"\"" + strings.TrimSpace(input.Country) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.UsrStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.UsrStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.UsrStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.UsrStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	defer db.Close()

	return prebuf.String()
}
