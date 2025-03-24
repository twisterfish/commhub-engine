package worktickets

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for perishable items data
//////////////////////////////////////////////////////////////////////////////////////////

import (
	//"cloud.google.com/go/storage"
	//"google.golang.org/appengine"
	//"time"
	"bytes"
	"datastores"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "CreateTicket":
		return CreateTicket(payload)
	case "AssignTicket":
		return AssignTicket(payload)
	case "TestAssignTicket":
		return TestAssignTicket(payload)
	case "RejectTicket":
		return RejectTicket(payload)
	case "GetAssetsForTicket":
		return GetAssetsForTicket(request, payload)
	case "GetTicketByGUID":
		return GetTicketByGUID(payload)
	case "GetMyTickets":
		return GetMyTickets(payload)
	case "GetTicketsInWorkspace":
		return GetTicketsInWorkspace(payload)
	case "GetTicketStatusDescriptions":
		return GetTicketStatusDescriptions(payload)
	case "GetTicketTypeDescriptions":
		return GetTicketTypeDescriptions(payload)
	case "SetTicketTitle":
		return SetTicketTitle(payload)
	case "SetTicketDescription":
		return SetTicketDescription(payload)
	case "SetTicketType":
		return SetTicketType(payload)
	case "SetTicketStatus":
		return SetTicketStatus(payload)
	case "SetTicketTimeStarted":
		return SetTicketTimeStarted(payload)
	case "SetTicketTimeFinished":
		return SetTicketTimeFinished(payload)
	case "SetTicketRunningTime":
		return SetTicketRunningTime(payload)
	case "SetTicketToPaused":
		return SetTicketToPaused(payload)
	case "SetTicketToBlocked":
		return SetTicketToBlocked(payload)
	case "SetTicketToClosed":
		return SetTicketToClosed(payload)
	case "OverrideTicketRunningTime":
		return OverrideTicketRunningTime(payload)
	case "SetTicketInventoryUsage":
		return SetTicketInventoryUsage(payload)
	case "GetSignedURLTest":
		return GetSignedURLTest()
	case "SetProductToTicket":
		return SetProductToTicket(payload)
	default:
		return "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Utility methods
// ////////////////////////////////////////////////////////////////////////////////////////
func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Internal secure URL for assets for tickets
// ////////////////////////////////////////////////////////////////////////////////////////
func getTicketAssetSignedURL(assetName string, request *events.APIGatewayProxyRequest, ContType string) string {

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("golang-test-one"),
		Key:    aws.String("DigMQ/PiOvLVHGSBgsLw+LrUt4RBpFYkhdvX+2hm"),
	})

	url, err := req.Presign(24 * time.Hour)

	if err != nil {
		return err.Error()
	} else {
		return url + "/" + assetName
	}

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Internal secure URL for assets for tickets
// ////////////////////////////////////////////////////////////////////////////////////////
func GetSignedURLTest() string {

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)
	// create a sample test file.ext
	asset := "/CloudLayout_1.png"

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("golang-test-one"),
		Key:    aws.String("DigMQ/PiOvLVHGSBgsLw+LrUt4RBpFYkhdvX+2hm"),
	})

	url, err := req.Presign(24 * time.Hour)

	if err != nil {
		return err.Error()
	} else {
		return url + asset
	}

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Creating a new work order ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"CreateTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_type_id":"7","local_guid":"48b1760150c746124e5aed03334cbf20","real_property_id":"1","title":"some title","description":"do some stuff"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CreateTicket(payload string) string {

	type NewWorkTicket struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WoTID   string `json:"work_ticket_type_id"`
		LocGUID string `json:"local_guid"`
		RpID    string `json:"real_property_id"`
		WoTi    string `json:"title"`
		WoDesc  string `json:"description"`
	}

	// used for parsing request
	var input NewWorkTicket
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		GrpStatus string `json:"out_status"`
		NewTkID   string `json:"out_new_ticket_id"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.create_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WoTID) + ", \"" + strings.TrimSpace(input.LocGUID) + "\", " + strings.TrimSpace(input.RpID) + ", \"" + strings.TrimSpace(input.WoTi) + "\", \"" + strings.TrimSpace(input.WoDesc) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.GrpStatus, &output.NewTkID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.GrpStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.GrpStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.GrpStatus + "\",")
			rows.WriteString("\"out_new_ticket_id\":\"" + output.NewTkID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Updating an existing ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"UpdateTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"7","work_ticket_type_id":"1","work_ticket_status_id":"1","ticket_group_id":"1033","running_time":"700","assigned_to_user_token":"48b1760150c746124e5aed03334cbf20","real_property_id":"1","title":"some title","description":"do some stuff"}
// ////////////////////////////////////////////////////////////////////////////////////////
func UpdateTicket(payload string) string {

	type NewWorkTicket struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WoTID  string `json:"work_ticket_id"`
		WoTTID string `json:"work_ticket_type_id"`
		WoTSID string `json:"work_ticket_status_id"`
		TGID   string `json:"ticket_group_id"`
		RTME   string `json:"running_time"`
		ATUTok string `json:"assigned_to_user_token"`
		RpID   string `json:"real_property_id"`
		WoTi   string `json:"title"`
		WoDesc string `json:"description"`
	}

	// used for parsing request
	var input NewWorkTicket
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		GrpStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.update_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WoTID) + "," + strings.TrimSpace(input.WoTTID) + "," + strings.TrimSpace(input.WoTSID) + "," + strings.TrimSpace(input.TGID) + "," + strings.TrimSpace(input.RTME) + ",\"" + strings.TrimSpace(input.ATUTok) + "\"," + strings.TrimSpace(input.RpID) + ",\"" + strings.TrimSpace(input.WoTi) + "\",\"" + strings.TrimSpace(input.WoDesc) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.GrpStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.GrpStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.GrpStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.GrpStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Assigning a work order ticket to someone
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"AssignTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","ticket_group_id":"1015","work_ticket_id":"1114","assignee_token":"8a9861cc86ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func AssignTicket(payload string) string {

	type TicketAssignee struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		TgrpID string `json:"ticket_group_id"`
		WTID   string `json:"work_ticket_id"`
		ASTok  string `json:"assignee_token"`
	}

	// used for parsing request
	var input TicketAssignee
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.assign_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.TgrpID) + ", " + strings.TrimSpace(input.WTID) + ", \"" + strings.TrimSpace(input.ASTok) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Reject a ticket that has been assigned to him - kick it back to the dispatcher who assigned it and put it in the staging group
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"RejectTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114"}
// ////////////////////////////////////////////////////////////////////////////////////////
func RejectTicket(payload string) string {

	type TicketAssignee struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
	}

	// used for parsing request
	var input TicketAssignee
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.reject_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Setting a work order ticket's description
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketDescription","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_description":"new ticket description verbiage"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketDescription(payload string) string {

	type TicketDesc struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
		TktDsc string `json:"ticket_description"`
	}

	// used for parsing request
	var input TicketDesc
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_description( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", \"" + strings.TrimSpace(input.TktDsc) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Setting a work order ticket's title
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketTitle","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_title":"new ticket title verbiage"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketTitle(payload string) string {

	type TicketTitle struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
		TktTtl string `json:"ticket_title"`
	}

	// used for parsing request
	var input TicketTitle
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_title( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", \"" + strings.TrimSpace(input.TktTtl) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Setting a work order ticket's status
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketStatus","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_status":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketStatus(payload string) string {

	type TicketStatus struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktStat string `json:"ticket_status"`
	}

	// used for parsing request
	var input TicketStatus
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_status( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktStat) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Setting a work order ticket's status
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketType","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_type":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketType(payload string) string {

	type TicketType struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktType string `json:"ticket_type"`
	}

	// used for parsing request
	var input TicketType
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_type( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktType) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set the running time on task.  The app controls the value in seconds of how long the task
// has been running.  server has no way of being accurate or even knowing, so we'll trust the app
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketRunningTime","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_running_time":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketRunningTime(payload string) string {

	type TicketRunning struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"ticket_running_time"`
	}

	// used for parsing request
	var input TicketRunning
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_running_time( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set the running time on task.  The app controls the value in seconds of how long the task
// has been running.  server has no way of being accurate or even knowing, so we'll trust the app
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketToPaused","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_running_time":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketToPaused(payload string) string {

	type TicketRunning struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"ticket_running_time"`
	}

	// used for parsing request
	var input TicketRunning
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_to_paused( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set the running time on task.  The app controls the value in seconds of how long the task
// has been running.  server has no way of being accurate or even knowing, so we'll trust the app
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketToBlocked","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_running_time":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketToBlocked(payload string) string {

	type TicketRunning struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"ticket_running_time"`
	}

	// used for parsing request
	var input TicketRunning
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_to_blocked( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set the running time on task.  The app controls the value in seconds of how long the task
// has been running.  server has no way of being accurate or even knowing, so we'll trust the app
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketToClosed","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketToClosed(payload string) string {

	type TicketRunning struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
	}

	// used for parsing request
	var input TicketRunning
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_to_closed( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Allow administrators to override the running time on the ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"OverrideTicketRunningTime","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","ticket_running_time":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func OverrideTicketRunningTime(payload string) string {

	type TicketRunning struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"ticket_running_time"`
	}

	// used for parsing request
	var input TicketRunning
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.override_ticket_running_time( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Starting a work order ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketTimeStarted","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","time_started":"1","ts_latitude":"26.122438","ts_longitude":"-80.137314"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketTimeStarted(payload string) string {

	type TicketTime struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"time_started"`
		TktLat  string `json:"ts_latitude"`
		TktLong string `json:"ts_longitude"`
	}

	// used for parsing request
	var input TicketTime
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_start_time( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + ", " + strings.TrimSpace(input.TktLat) + ", " + strings.TrimSpace(input.TktLong) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Finished a work order ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketTimeFinished","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","time_finished":"1","tf_latitude":"26.122438","tf_longitude":"-80.137314"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketTimeFinished(payload string) string {

	type TicketTime struct {
		Token   string `json:"api_token"`
		Signal  string `json:"signal"`
		Action  string `json:"action"`
		EUTok   string `json:"end_user_token"`
		WSTok   string `json:"workspace_token"`
		WTID    string `json:"work_ticket_id"`
		TktTime string `json:"time_finished"`
		TktLat  string `json:"tf_latitude"`
		TktLong string `json:"tf_longitude"`
	}

	// used for parsing request
	var input TicketTime
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_finished_time( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.TktTime) + ", " + strings.TrimSpace(input.TktLat) + ", " + strings.TrimSpace(input.TktLong) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Setting the inventory items that are used on the work order
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetTicketInventoryUsage","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","work_ticket_id":"1114","product_id":"1","product_id_qty":"26"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketInventoryUsage(payload string) string {

	type NewTicketItem struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
		PrdID  string `json:"product_id"`
		PrdQ   string `json:"product_id_qty"`
	}

	// used for parsing request
	var input NewTicketItem
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_ticket_inventory_usage( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.PrdID) + ", " + strings.TrimSpace(input.PrdQ) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull assets associated with this ticket
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetAssetsForTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","local_guid":"1"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetAssetsForTicket(request *events.APIGatewayProxyRequest, payload string) string {

	type TicketAssetRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTGUID string `json:"local_guid"`
	}
	var input TicketAssetRequest

	type TicketAssetList struct {
		WTGUID string `json:"local_guid"`
		AsNM   string `json:"asset_name"`
		AsCnT  string `json:"content_type"`
		AsSZE  string `json:"asset_size_bytes"`
		AsSZH  string `json:"asset_size_height"`
		AsSZW  string `json:"asset_size_width"`
		AsDesc string `json:"asset_description"`
	}
	var output TicketAssetList

	pByte := []byte(payload)
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()
	defer db.Close()
	// if there is an error opening the connection, handle it
	if err != nil {
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	item_count = 0
	query := "SELECT local_guid, asset_name, content_type, asset_size_bytes, asset_size_height, asset_size_width, IFNULL(asset_description,'') FROM commhub_junction.work_ticket_asset WHERE local_guid = " + input.WTGUID
	results, err := db.Query(query)

	if err != nil {
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our tal struct
			err = results.Scan(&output.WTGUID, &output.AsNM, &output.AsCnT, &output.AsSZE, &output.AsSZH, &output.AsSZW, &output.AsDesc)

			if err != nil {
				//panic(err.Error()) // proper error handling instead of panic in your app
				defer db.Close()
				defer results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tal.ID)
			rows.WriteString("{\"local_guid\":\"" + output.WTGUID + "\",")
			rows.WriteString("\"asset_name\":\"" + output.AsNM + "\",")
			rows.WriteString("\"content_type\":\"" + output.AsCnT + "\",")
			rows.WriteString("\"asset_size_bytes\":\"" + output.AsSZE + "\",")
			rows.WriteString("\"asset_size_height\":\"" + output.AsSZH + "\",")
			rows.WriteString("\"asset_size_width\":\"" + output.AsSZW + "\",")

			//rows.WriteString("\"signed_url\":\"" +  getTicketAssetSignedURL( output.AsNM, request, output.AsCnT  ) + "\"," )

			rows.WriteString("\"signed_url\":\" test url until lambda is up and running\",")
			rows.WriteString("\"asset_description\":\"" + output.AsDesc + "\"},")
			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a workticket by local GUID - you must be authorized to do so or you will get nothing back
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetTicketByGUID","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","local_guid":"48b176076fc74698be5aed03334cbf20"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetTicketByGUID(payload string) string {

	type TicketRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WoGUID string `json:"local_guid"`
	}
	// used for parsing request
	var input TicketRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type WorkTickets struct {
		WSTok       string `json:"workspace_token"`
		WoID        string `json:"work_ticket_id"`
		LoGUID      string `json:"local_guid"`
		TctGrID     string `json:"ticket_group_id"`
		WoSID       string `json:"work_ticket_status_id"`
		WoTID       string `json:"work_ticket_type_id"`
		RpID        string `json:"real_property_id"`
		Rtime       string `json:"running_time"`
		TimeCr      string `json:"time_created"`
		TimeSt      string `json:"time_started"`
		TsLat       string `json:"ts_latitude"`
		TsLong      string `json:"ts_longitude"`
		TimeFin     string `json:"time_finished"`
		TimeFinLat  string `json:"tf_latitude"`
		TimeFinLong string `json:"tf_longitude"`
		WoTi        string `json:"title"`
		WoDesc      string `json:"description"`
	}

	var output WorkTickets

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

	query := "CALL commhub_junction.get_ticket_info(\"" + input.EUTok + "\",\"" + input.WSTok + "\",\"" + input.WoGUID + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			err = results.Scan(&output.WSTok, &output.WoID, &output.LoGUID, &output.TctGrID, &output.WoSID, &output.WoTID, &output.RpID, &output.Rtime, &output.TimeCr, &output.TimeSt, &output.TsLat, &output.TsLong, &output.TimeFin, &output.TimeFinLat, &output.TimeFinLong, &output.WoTi, &output.WoDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(output.ID)
			rows.WriteString("{\"workspace_token\":\"" + output.WSTok + "\",")
			rows.WriteString("\"work_ticket_id\":\"" + output.WoID + "\",")
			rows.WriteString("\"local_guid\":\"" + output.LoGUID + "\",")
			rows.WriteString("\"ticket_group_id\":\"" + output.TctGrID + "\",")
			rows.WriteString("\"work_ticket_status_id\":\"" + output.WoSID + "\",")
			rows.WriteString("\"work_ticket_type_id\":\"" + output.WoTID + "\",")
			rows.WriteString("\"real_property_id\":\"" + output.RpID + "\",")
			rows.WriteString("\"running_time\":\"" + output.Rtime + "\",")
			rows.WriteString("\"time_created\":\"" + output.TimeCr + "\",")
			rows.WriteString("\"time_started\":\"" + output.TimeSt + "\",")
			rows.WriteString("\"ts_latitude\":\"" + output.TsLat + "\",")
			rows.WriteString("\"ts_longitude\":\"" + output.TsLong + "\",")
			rows.WriteString("\"time_finished\":\"" + output.TimeFin + "\",")
			rows.WriteString("\"tf_latitude\":\"" + output.TimeFinLat + "\",")
			rows.WriteString("\"tf_longitude\":\"" + output.TimeFinLong + "\",")
			rows.WriteString("\"title\":\"" + jsonEscape(output.WoTi) + "\",")
			rows.WriteString("\"description\":\"" + jsonEscape(output.WoDesc) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a a workers own tickets that were assigned to him
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetMyTicketInfo","end_user_token":"8a9861cc86ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetMyTickets(payload string) string {

	type TicketRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WoGUID string `json:"local_guid"`
	}
	// used for parsing request
	var input TicketRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type WorkTickets struct {
		WSTok       string `json:"workspace_token"`
		WoID        string `json:"work_ticket_id"`
		LoGUID      string `json:"local_guid"`
		TctGrID     string `json:"ticket_group_id"`
		WoSID       string `json:"work_ticket_status_id"`
		WoTID       string `json:"work_ticket_type_id"`
		RpID        string `json:"real_property_id"`
		Rtime       string `json:"running_time"`
		TimeCr      string `json:"time_created"`
		TimeSt      string `json:"time_started"`
		TsLat       string `json:"ts_latitude"`
		TsLong      string `json:"ts_longitude"`
		TimeFin     string `json:"time_finished"`
		TimeFinLat  string `json:"tf_latitude"`
		TimeFinLong string `json:"tf_longitude"`
		WoTi        string `json:"title"`
		WoDesc      string `json:"description"`
	}

	var output WorkTickets

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

	query := "CALL commhub_junction.get_my_tickets(\"" + input.EUTok + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			err = results.Scan(&output.WSTok, &output.WoID, &output.LoGUID, &output.TctGrID, &output.WoSID, &output.WoTID, &output.RpID, &output.Rtime, &output.TimeCr, &output.TimeSt, &output.TsLat, &output.TsLong, &output.TimeFin, &output.TimeFinLat, &output.TimeFinLong, &output.WoTi, &output.WoDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(output.ID)
			rows.WriteString("{\"workspace_token\":\"" + output.WSTok + "\",")
			rows.WriteString("\"work_ticket_id\":\"" + output.WoID + "\",")
			rows.WriteString("\"local_guid\":\"" + output.LoGUID + "\",")
			rows.WriteString("\"ticket_group_id\":\"" + output.TctGrID + "\",")
			rows.WriteString("\"work_ticket_status_id\":\"" + output.WoSID + "\",")
			rows.WriteString("\"work_ticket_type_id\":\"" + output.WoTID + "\",")
			rows.WriteString("\"real_property_id\":\"" + output.RpID + "\",")
			rows.WriteString("\"running_time\":\"" + output.Rtime + "\",")
			rows.WriteString("\"time_created\":\"" + output.TimeCr + "\",")
			rows.WriteString("\"time_started\":\"" + output.TimeSt + "\",")
			rows.WriteString("\"ts_latitude\":\"" + output.TsLat + "\",")
			rows.WriteString("\"ts_longitude\":\"" + output.TsLong + "\",")
			rows.WriteString("\"time_finished\":\"" + output.TimeFin + "\",")
			rows.WriteString("\"tf_latitude\":\"" + output.TimeFinLat + "\",")
			rows.WriteString("\"tf_longitude\":\"" + output.TimeFinLong + "\",")
			rows.WriteString("\"title\":\"" + jsonEscape(output.WoTi) + "\",")
			rows.WriteString("\"description\":\"" + jsonEscape(output.WoDesc) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of all tickets in a workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetTicketsInWorkspace","workspace_token":"8a9861cc86ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetTicketsInWorkspace(payload string) string {

	type TicketRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		WSTok  string `json:"workspace_token"`
	}
	// used for parsing request
	var input TicketRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type WorkTickets struct {
		WSTok       string `json:"workspace_token"`
		WoID        string `json:"work_ticket_id"`
		LoGUID      string `json:"local_guid"`
		TctGrID     string `json:"ticket_group_id"`
		As2ID       string `json:"assigned_to_user_id"`
		AsByID      string `json:"assigned_by_user_id"`
		CrtrUsrID   string `json:"creator_user_id"`
		WoSID       string `json:"work_ticket_status_id"`
		WoTID       string `json:"work_ticket_type_id"`
		RpID        string `json:"real_property_id"`
		Rtime       string `json:"running_time"`
		TimeCr      string `json:"time_created"`
		TimeSt      string `json:"time_started"`
		TsLat       string `json:"ts_latitude"`
		TsLong      string `json:"ts_longitude"`
		TimeFin     string `json:"time_finished"`
		TimeFinLat  string `json:"tf_latitude"`
		TimeFinLong string `json:"tf_longitude"`
		WoTi        string `json:"title"`
		WoDesc      string `json:"description"`
	}

	var output WorkTickets

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

	query := "CALL commhub_junction.get_tickets_in_workspace(\"" + input.WSTok + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			err = results.Scan(&output.WSTok, &output.WoID, &output.LoGUID, &output.TctGrID, &output.As2ID, &output.AsByID, &output.CrtrUsrID, &output.WoSID, &output.WoTID, &output.RpID, &output.Rtime, &output.TimeCr, &output.TimeSt, &output.TsLat, &output.TsLong, &output.TimeFin, &output.TimeFinLat, &output.TimeFinLong, &output.WoTi, &output.WoDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(output.ID)
			rows.WriteString("{\"workspace_token\":\"" + output.WSTok + "\",")
			rows.WriteString("\"work_ticket_id\":\"" + output.WoID + "\",")
			rows.WriteString("\"local_guid\":\"" + output.LoGUID + "\",")
			rows.WriteString("\"ticket_group_id\":\"" + output.TctGrID + "\",")
			rows.WriteString("\"assigned_to_user_id\":\"" + output.As2ID + "\",")
			rows.WriteString("\"assigned_by_user_id\":\"" + output.AsByID + "\",")
			rows.WriteString("\"creator_user_id\":\"" + output.CrtrUsrID + "\",")
			rows.WriteString("\"work_ticket_status_id\":\"" + output.WoSID + "\",")
			rows.WriteString("\"work_ticket_type_id\":\"" + output.WoTID + "\",")
			rows.WriteString("\"real_property_id\":\"" + output.RpID + "\",")
			rows.WriteString("\"running_time\":\"" + output.Rtime + "\",")
			rows.WriteString("\"time_created\":\"" + output.TimeCr + "\",")
			rows.WriteString("\"time_started\":\"" + output.TimeSt + "\",")
			rows.WriteString("\"ts_latitude\":\"" + output.TsLat + "\",")
			rows.WriteString("\"ts_longitude\":\"" + output.TsLong + "\",")
			rows.WriteString("\"time_finished\":\"" + output.TimeFin + "\",")
			rows.WriteString("\"tf_latitude\":\"" + output.TimeFinLat + "\",")
			rows.WriteString("\"tf_longitude\":\"" + output.TimeFinLong + "\",")
			rows.WriteString("\"title\":\"" + jsonEscape(output.WoTi) + "\",")
			rows.WriteString("\"description\":\"" + jsonEscape(output.WoDesc) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of the ticket status descriptions
// {"api_token": "37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetTicketStatusDescriptions","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetTicketStatusDescriptions(payload string) string {

	type TicketRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
	}
	// used for parsing request
	var input TicketRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type WorkTicketStatus struct {
		WoSID  string `json:"work_ticket_status_id"`
		WoDesc string `json:"description"`
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

	query := "SELECT work_ticket_status_id, description FROM commhub_junction.work_ticket_status"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag WorkTicketStatus
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.WoSID, &tag.WoDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"work_ticket_status_id\":\"" + tag.WoSID + "\",")
			rows.WriteString("\"description\":\"" + tag.WoDesc + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull a list of the ticket type descriptions
// {"api_token": "37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"GetTicketTypeDescriptions","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetTicketTypeDescriptions(payload string) string {

	type TicketRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
	}
	// used for parsing request
	var input TicketRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type WorkTicketType struct {
		WoTID  string `json:"work_ticket_type_id"`
		WoDesc string `json:"description"`
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

	query := "SELECT work_ticket_type_id, description FROM commhub_junction.work_ticket_type"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag WorkTicketType
			// for each record, scan the result into our tag struct
			err = results.Scan(&tag.WoTID, &tag.WoDesc)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"work_ticket_type_id\":\"" + tag.WoTID + "\",")
			rows.WriteString("\"description\":\"" + tag.WoDesc + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// assigning product inventory items to a work order
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"SetProductToTicket","end_user_token":"db2e3ddb21b411ea95e10ece0304bc53","workspace_token":"db2e5cec21b411ea95e10ece0304bc53","work_ticket_id":"1","product_id":"3","product_id_qty":"100"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetProductToTicket(payload string) string {

	type NewTicketItem struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		WTID   string `json:"work_ticket_id"`
		PrdID  string `json:"product_id"`
		PrdQ   string `json:"product_id_qty"`
	}

	// used for parsing request
	var input NewTicketItem
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.set_product_to_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.WTID) + ", " + strings.TrimSpace(input.PrdID) + ", " + strings.TrimSpace(input.PrdQ) + " )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// *******DO NOT DEPLOY IN A PRODUCTION ENVIRONMENT
// Assigning a work order ticket to someone even if the ticket has started
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"worktickets","action":"AssignTicket","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","ticket_group_id":"1015","work_ticket_id":"1114","assignee_token":"8a9861cc86ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func TestAssignTicket(payload string) string {

	type TicketAssignee struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
		TgrpID string `json:"ticket_group_id"`
		WTID   string `json:"work_ticket_id"`
		ASTok  string `json:"assignee_token"`
	}

	// used for parsing request
	var input TicketAssignee
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type TcktRecord struct {
		TckStatus string `json:"out_status"`
	}
	var output TcktRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.test_assign_ticket( \"" + strings.TrimSpace(input.EUTok) + "\", \"" + strings.TrimSpace(input.WSTok) + "\", " + strings.TrimSpace(input.TgrpID) + ", " + strings.TrimSpace(input.WTID) + ", \"" + strings.TrimSpace(input.ASTok) + "\" )"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.TckStatus)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.TckStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.TckStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.TckStatus + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}
