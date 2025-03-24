package authorize

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for authorization of credentials
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"encoding/json"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(action string, payload string) string {
	switch action {
	case "LogIn":
		return CheckUserCredentials(payload)
	case "CreateUserProfile":
		return CreateUserProfile(payload)
	case "TestCreateUserProfile":
		return TestCreateUserProfile(payload)
	default:
		return "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Enter a new profile in the database
// {"signal":"authorize","action":"TestCreateUserProfile","emid":"jon.public@commhubstuff.com","pwid":"somerandompassword","ein_tax_id":"321","ssn_tax_id":"123","last_name":"Public","middle_name":"Q","first_name":"Jon","address1":"111 Richie Way","address2":"unit 3","city":"Miami","province_state":"OK","zip_postal_code":"33318","phone":"123108900","email":"joe.shmuck@anymail.com","country_code":"US"}
// ////////////////////////////////////////////////////////////////////////////////////////
func TestCreateUserProfile(payload string) string {

	type NewProfile struct {
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EMID    string `json:"emid"`
		PWID    string `json:"pwid"`
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

	type UsrRecord struct {
		UsrStatus string `json:"out_status"`
		UsrID     string `json:"out_new_user_id"`
		UsrUUID   string `json:"out_new_user_uuid"`
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

	query := "CALL commhub_junction.test_create_user_profile(" +
		"\"" + strings.TrimSpace(input.EMID) + "\"," +
		"\"" + strings.TrimSpace(input.PWID) + "\"," +
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
			err = results.Scan(&output.UsrStatus, &output.UsrID, &output.UsrUUID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.UsrStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.UsrStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.UsrStatus + "\",")
			rows.WriteString("\"out_new_user_id\":\"" + output.UsrID + "\",")
			rows.WriteString("\"out_new_user_uuid\":\"" + output.UsrUUID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	defer db.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Enter a new profile in the database
// {"signal":"authorize","action":"CreateUserProfile","emid":"jon.public@commhubstuff.com","pwid":"somerandompassword","ein_tax_id":"321","ssn_tax_id":"123","last_name":"Public","middle_name":"Q","first_name":"Jon","address1":"111 Richie Way","address2":"unit 3","city":"Miami","province_state":"OK","zip_postal_code":"33318","phone":"123108900","email":"joe.shmuck@anymail.com","country_code":"US"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CreateUserProfile(payload string) string {

	type NewProfile struct {
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EMID    string `json:"emid"`
		PWID    string `json:"pwid"`
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

	type UsrRecord struct {
		UsrStatus string `json:"out_status"`
		UsrID     string `json:"out_new_user_id"`
		UsrUUID   string `json:"out_new_user_uuid"`
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

	query := "CALL commhub_junction.create_user_profile(" +
		"\"" + strings.TrimSpace(input.EMID) + "\"," +
		"\"" + strings.TrimSpace(input.PWID) + "\"," +
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
			err = results.Scan(&output.UsrStatus, &output.UsrID, &output.UsrUUID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.UsrStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.UsrStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.UsrStatus + "\",")
			rows.WriteString("\"out_new_user_id\":\"" + output.UsrID + "\",")
			rows.WriteString("\"out_new_user_uuid\":\"" + output.UsrUUID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	defer db.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Checks user credentials and returns the current API token on success
// {"signal":"authorize","action":"LogIn","email":"edward.anderson@commhubstuff.com","pass":"testpass"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CheckUserCredentials(payload string) string {

	type inputData struct {
		Signal string `json:"signal"`
		Action string `json:"action"`
		Email  string `json:"email"`
		Pass   string `json:"pass"`
	}

	type outputData struct {
		Status string `json:"out_status"`
		EUID   string `json:"end_user_id"`
		EUTok  string `json:"end_user_token"`
		ApiTok string `json:"api_token"`
	}

	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.check_user_credentials(\"" + strings.TrimSpace(input.Email) + "\", \"" + strings.TrimSpace(input.Pass) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status, &output.EUID, &output.EUTok, &output.ApiTok)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) == "invalid_user" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\",")
			rows.WriteString("\"end_user_id\":\"" + output.EUID + "\",")
			rows.WriteString("\"end_user_token\":\"" + output.EUTok + "\",")
			rows.WriteString("\"api_token\":\"" + output.ApiTok + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"authorize\",\"action\":\"LogIn\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}
