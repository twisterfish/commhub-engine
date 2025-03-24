package workspaces

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for workspaces
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

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// move to permission labels and tags
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {

	switch action {
	case "CreateWorkspace":
		return CreateWorkspace(payload)
	case "CreateWorkspaceOnDevice":
		return CreateWorkspaceOnDevice(payload)
	case "GetWorkspaceByOwner":
		return GetWorkspaceByOwner(payload)
	case "AddMemberToWorkspace":
		return AddMemberToWorkspace(payload)
	case "RemoveMemberFromWorkspace":
		return RemoveMemberFromWorkspace(payload)
	case "GetMembersInMyWorkspace":
		return GetMembersInMyWorkspace(payload)
	case "GetMyWorkspaces":
		return GetMyWorkspaces(payload)
	case "GetTicketGroupsInWorkspace":
		return GetTicketGroupsInWorkspace(payload)
	case "SetWorkspaceName": // use token
		return SetWorkspaceName(payload)
	case "SetWorkspaceDescription": // use token
		return SetWorkspaceDescription(payload)
	case "CreateTicketGroup":
		return CreateTicketGroup(payload)
	case "SetTicketGroupName":
		return SetTicketGroupName(payload)
	case "SetTicketGroupDescription":
		return SetTicketGroupDescription(payload)
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
// Pull the ticket groups in a workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"GetTicketGroupsInWorkspace","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetTicketGroupsInWorkspace(payload string) string {

	type WorkspaceRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		ENUTok string `json:"end_user_token"`
		WSTok  string `json:"workspace_token"`
	}

	type WorkspacesMem struct {
		WorkspaceToken string `json:"workspace_token"`
		TicketGrpID    string `json:"ticket_group_id"`
		TicketGrpNM    string `json:"ticket_group_name"`
		TicketGrpDESC  string `json:"ticket_group_description"`
	}

	var input WorkspaceRequest
	var output WorkspacesMem
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.get_ticket_groups_in_workspace(\"" + strings.TrimSpace(input.ENUTok) + "\",\"" + strings.TrimSpace(input.WSTok) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			err = results.Scan(&output.WorkspaceToken, &output.TicketGrpID, &output.TicketGrpNM, &output.TicketGrpDESC)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"workspace_token\":\"" + output.WorkspaceToken + "\",")
			rows.WriteString("\"ticket_group_id\":\"" + output.TicketGrpID + "\",")
			rows.WriteString("\"ticket_group_name\":\"" + jsonEscape(output.TicketGrpNM) + "\",")
			rows.WriteString("\"ticket_group_description\":\"" + jsonEscape(output.TicketGrpDESC) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Enter a new workspace in the data
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"CreateWorkspace","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_name":"COOL COMPANY","workspace_description":"Based out of Fort Lauderdale"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CreateWorkspace(payload string) string {

	type NewWorkspace struct {
		Token         string `json:"api_token"`
		Signal        string `json:"signal"`
		Action        string `json:"action"`
		WorkspaceCRTR string `json:"end_user_token"`
		WorkspaceNM   string `json:"workspace_name"`
		WorkspaceDESC string `json:"workspace_description"`
	}

	// used for parsing request
	var input NewWorkspace
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type WorkspaceRecord struct {
		WorkspaceStatus string `json:"out_status"`
		WorkspaceUUID   string `json:"out_new_ws_uuid"`
		GrpID           string `json:"out_new_group_id"`
	}
	var output WorkspaceRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.create_workspace_staging_group(\"" + strings.TrimSpace(input.WorkspaceCRTR) + "\", \"" + strings.TrimSpace(input.WorkspaceNM) + "\", \"" + strings.TrimSpace(input.WorkspaceDESC) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.WorkspaceStatus, &output.WorkspaceUUID, &output.GrpID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.WorkspaceStatus) == "invalid_user" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.WorkspaceStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.WorkspaceStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.WorkspaceStatus + "\",")
			rows.WriteString("\"out_new_workspace_uuid\":\"" + output.WorkspaceUUID + "\",")
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Create a new workspace with a device generated GUID
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"CreateWorkspaceOnDevice","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"YC98oij345trghui45tergf9huie","workspace_name":"COOL COMPANY","workspace_description":"Based out of Fort Lauderdale"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CreateWorkspaceOnDevice(payload string) string {

	type NewWorkspace struct {
		Token         string `json:"api_token"`
		Signal        string `json:"signal"`
		Action        string `json:"action"`
		WorkspaceCRTR string `json:"end_user_token"`
		WorkspaceTKN  string `json:"workspace_token"`
		WorkspaceNM   string `json:"workspace_name"`
		WorkspaceDESC string `json:"workspace_description"`
	}

	// used for parsing request
	var input NewWorkspace
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type WorkspaceRecord struct {
		WorkspaceStatus string `json:"out_status"`
		WorkspaceUUID   string `json:"out_new_workspace_uuid"`
		GrpID           string `json:"out_new_group_id"`
	}
	var output WorkspaceRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.create_workspace_on_device(\"" + strings.TrimSpace(input.WorkspaceCRTR) + "\", \"" + strings.TrimSpace(input.WorkspaceTKN) + "\", \"" + strings.TrimSpace(input.WorkspaceNM) + "\", \"" + strings.TrimSpace(input.WorkspaceDESC) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.WorkspaceStatus, &output.WorkspaceUUID, &output.GrpID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.WorkspaceStatus) == "invalid_user" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.WorkspaceStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.WorkspaceStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.WorkspaceStatus + "\",")
			rows.WriteString("\"out_new_workspace_uuid\":\"" + output.WorkspaceUUID + "\",")
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull an workspace by owner's token
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"GetWorkspaceByOwner","workspace_owner_token":"8a985cfe86ce11e9bc42526af7764f64"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetWorkspaceByOwner(payload string) string {

	type WorkspaceRequest struct {
		Token    string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		OwnToken string `json:"workspace_owner_token"`
	}
	// used for parsing request
	var input WorkspaceRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type WorkspaceRecord struct {
		WorkspaceToken string `json:"workspace_token"`
		WorkspaceNM    string `json:"workspace_name"`
		WorkspaceDESC  string `json:"workspace_description"`
	}

	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	db, err := datastores.OpenRDS()
	defer db.Close()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	item_count = 0

	query := "SELECT WSP.workspace_token, IFNULL(WSP.workspace_name,''), IFNULL(WSP.workspace_description,'') FROM commhub_junction.workspace AS WSP INNER JOIN commhub_junction.end_user AS EndU ON WSP.workspace_owner_id = EndU.end_user_id WHERE EndU.end_user_token = \"" + strings.TrimSpace(input.OwnToken) + "\""
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag WorkspaceRecord
			// for each record, scan the result into our  struct
			err = results.Scan(&tag.WorkspaceToken, &tag.WorkspaceNM, &tag.WorkspaceDESC)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"workspace_token\":\"" + tag.WorkspaceToken + "\",")
			rows.WriteString("\"workspace_name\":\"" + jsonEscape(tag.WorkspaceNM) + "\",")
			rows.WriteString("\"workspace_description\":\"" + jsonEscape(tag.WorkspaceDESC) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Add a user to an workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"AddMemberToWorkspace","end_user_token":"8a985a9286ce11e9bc42526af7764f64","new_member_token":"8a9861cc86ce11e9bc42526af7764f64","workspace_token":"17f0a2866e7e408d9ca3810dc801e046","workspace_permission_id":"200"}
// ////////////////////////////////////////////////////////////////////////////////////////
func AddMemberToWorkspace(payload string) string {

	type inputData struct {
		Token           string `json:"api_token"`
		Signal          string `json:"signal"`
		Action          string `json:"action"`
		EndUsrToken     string `json:"end_user_token"`
		NewMemToken     string `json:"new_member_token"`
		WorkspaceToken  string `json:"workspace_token"`
		WorkspacePermID string `json:"workspace_permission_id"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}

	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.add_member_to_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.NewMemToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\", \"" + strings.TrimSpace(input.WorkspacePermID) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Remove a member from an workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"RemoveMemberFromWorkspace","end_user_token":"8a985a9286ce11e9bc42526af7764f64","target_member_token":"8a9861cc86ce11e9bc42526af7764f64","workspace_token":"17f0a2866e7e408d9ca3810dc801e046"}
// ////////////////////////////////////////////////////////////////////////////////////////
func RemoveMemberFromWorkspace(payload string) string {

	type inputData struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		TrgMemToken    string `json:"target_member_token"`
		WorkspaceToken string `json:"workspace_token"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.remove_member_from_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.TrgMemToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull members in an workspace by workspace token
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"GetMembersInMyWorkspace","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"17f0a2866e7e408d9ca3810dc801e046"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetMembersInMyWorkspace(payload string) string {

	type WorkspaceRequest struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		WorkspaceToken string `json:"workspace_token"`
	}
	// used for parsing request
	var input WorkspaceRequest
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	type WorkspaceMemXREF struct {
		WorkspaceMEM string `json:"end_user_token"`
		WorkspaceLvl string `json:"workspace_permission"`
		MemFnm       string `json:"first_name"`
		MemLnm       string `json:"last_name"`
		PhoNm        string `json:"phone"`
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

	query := "CALL commhub_junction.get_members_in_my_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			var tag WorkspaceMemXREF
			// for each record, scan the result into our  struct
			err = results.Scan(&tag.WorkspaceMEM, &tag.WorkspaceLvl, &tag.MemFnm, &tag.MemLnm, &tag.PhoNm)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"end_user_token\":\"" + tag.WorkspaceMEM + "\",")
			rows.WriteString("\"workspace_permission\":\"" + tag.WorkspaceLvl + "\",")
			rows.WriteString("\"first_name\":\"" + tag.MemFnm + "\",")
			rows.WriteString("\"last_name\":\"" + tag.MemLnm + "\",")
			rows.WriteString("\"phone_number\":\"" + tag.PhoNm + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Pull workspaces I'm in
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"GetMyWorkspaces","end_user_token":"eb170ea62a4448b4a609c0521fbb4cf9"}
// ////////////////////////////////////////////////////////////////////////////////////////
func GetMyWorkspaces(payload string) string {

	type WorkspaceRequest struct {
		Token  string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EUTok  string `json:"end_user_token"`
	}

	type WorkspacesMem struct {
		WorkspaceToken string `json:"workspace_token"`
		WorkspaceLvl   string `json:"workspace_permission"`
		WorkspaceNM    string `json:"workspace_name"`
		WorkspaceDESC  string `json:"workspace_description"`
	}

	var input WorkspaceRequest
	var output WorkspacesMem
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int
	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.get_my_workspaces(\"" + strings.TrimSpace(input.EUTok) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {
			err = results.Scan(&output.WorkspaceToken, &output.WorkspaceLvl, &output.WorkspaceNM, &output.WorkspaceDESC)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"workspace_token\":\"" + output.WorkspaceToken + "\",")
			rows.WriteString("\"workspace_permission\":\"" + output.WorkspaceLvl + "\",")
			rows.WriteString("\"workspace_name\":\"" + jsonEscape(output.WorkspaceNM) + "\",")
			rows.WriteString("\"workspace_description\":\"" + jsonEscape(output.WorkspaceDESC) + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set or update a workspace name - only workspace owners can do this
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"SetWorkspaceName","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","workspace_name":"COOL COMPANY"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetWorkspaceName(payload string) string {

	type inputData struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		WorkspaceToken string `json:"workspace_token"`
		WorkspaceName  string `json:"workspace_name"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.set_workspace_name(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\", \"" + strings.TrimSpace(input.WorkspaceName) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set or update an workspace description - only the workspace owner can do this.
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"SetWorkspaceDescription","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","workspace_description":"COOL COMPANY"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetWorkspaceDescription(payload string) string {

	type inputData struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		WorkspaceToken string `json:"workspace_token"`
		WorkspaceDesc  string `json:"workspace_description"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.set_workspace_description(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\", \"" + strings.TrimSpace(input.WorkspaceDesc) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Enter a new ticket group in the data
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"CreateTicketGroup","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_name":"B OCEAN","ticket_group_description":"Based out of Fort Lauderdale"}
// ////////////////////////////////////////////////////////////////////////////////////////
func CreateTicketGroup(payload string) string {

	type NewGroup struct {
		Token        string `json:"api_token"`
		Signal       string `json:"signal"`
		Action       string `json:"action"`
		EuTok        string `json:"end_user_token"`
		WorkspaceTok string `json:"workspace_token"`
		GrpNM        string `json:"ticket_group_name"`
		GrpDESC      string `json:"ticket_group_description"`
	}

	// used for parsing request
	var input NewGroup
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

	pByte := []byte(payload)

	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	type GrpRecord struct {
		GrpStatus string `json:"out_status"`
		GrpID     string `json:"out_new_group_id"`
	}
	var output GrpRecord

	db, err := datastores.OpenRDS()

	// if there is an error opening the connection, handle it - don't close on bad connection?
	if err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	item_count = 0

	query := "CALL commhub_junction.create_ticket_group(\"" + strings.TrimSpace(input.EuTok) + "\", \"" + strings.TrimSpace(input.WorkspaceTok) + "\", \"" + strings.TrimSpace(input.GrpNM) + "\", \"" + strings.TrimSpace(input.GrpDESC) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.GrpStatus, &output.GrpID)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.GrpStatus) == "invalid_user" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.GrpStatus) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.GrpStatus) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.GrpStatus + "\",")
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set or update a ticket_group name in the data
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"SetTicketGroupName","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_id":"1029","ticket_group_name":"Water Pipes"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketGroupName(payload string) string {

	type inputData struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		WorkspaceToken string `json:"workspace_token"`
		GrpID          string `json:"ticket_group_id"`
		GrpName        string `json:"ticket_group_name"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.set_ticket_group_name(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\", \"" + strings.TrimSpace(input.GrpID) + "\", \"" + strings.TrimSpace(input.GrpName) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Set or update a ticket_group name in the data
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"workspaces","action":"SetTicketGroupDescription","end_user_token":"8a985a9286ce11e9bc42526af7764f64","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_id":"1029","ticket_group_description":"B Ocean Hotel"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SetTicketGroupDescription(payload string) string {

	type inputData struct {
		Token          string `json:"api_token"`
		Signal         string `json:"signal"`
		Action         string `json:"action"`
		EndUsrToken    string `json:"end_user_token"`
		WorkspaceToken string `json:"workspace_token"`
		GrpID          string `json:"ticket_group_id"`
		GrpDesc        string `json:"ticket_group_description"`
	}

	type outputData struct {
		Status string `json:"out_status"`
	}
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
	var item_count int

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

	query := "CALL commhub_junction.set_ticket_group_description(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WorkspaceToken) + "\", \"" + strings.TrimSpace(input.GrpID) + "\", \"" + strings.TrimSpace(input.GrpDesc) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
