package invitations

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for invitations to workspaces and membership in general
//////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"datastores"
	"encoding/json"
	"net/http"
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
	case "InviteMemberToWorkspace":
		return InviteMemberToWorkspace(payload, request)
	case "AcceptInviteToWorkspace":
		return AcceptInviteToWorkspace(payload)
	case "RejectInviteToWorkspace":
		return RejectInviteToWorkspace(payload)
	default:
		return "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Invite a user to a workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"invitations","action":"InviteMemberToWorkspace","end_user_token":"2325418da6fb11e9a58342010a8e0121","workspace_token":"eb170ea62a4448b4a609c0521fbb4cf9","workspace_permission_id":"400","target_email":"jonparse@email.com"}
// ////////////////////////////////////////////////////////////////////////////////////////
func InviteMemberToWorkspace(payload string, request *events.APIGatewayProxyRequest) string {

	type inputData struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		TargetEmail string `json:"target_email"`
		WSToken     string `json:"workspace_token"`
		WSPermID    string `json:"workspace_permission_id"`
	}

	type outputData struct {
		Status  string `json:"out_status"`
		InvTokn string `json:"out_invite_token"`
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

	query := "CALL commhub_junction.invite_member_to_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.WSToken) + "\"," + strings.TrimSpace(input.WSPermID) + ", \"" + strings.TrimSpace(input.TargetEmail) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {

		for results.Next() {

			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status, &output.InvTokn)

			if err != nil {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success" {
				results.Close()
				return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}

			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\",")
			rows.WriteString("\"out_invite_token\":\"" + output.InvTokn + "\"},")

			item_count++
		}
	}

	prebuf.WriteString("{\"signal\":\"" + input.Signal + "\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")
	prebuf.WriteString(strings.TrimRight(rows.String(), ",") + "]}") // take off the trailing comma in the set and cap it

	results.Close()

	iscool := SendEmailForInvitation(strings.TrimSpace(input.TargetEmail), strings.TrimSpace(output.InvTokn))
	//iscool := CurlTest()
	// fire off the email confirmation
	if iscool != "success" {
		return iscool
	}

	return prebuf.String()

}

// ////////////////////////////////////////////////////////////////////////////////////////
// SMTP function to fire off invites sourbeer13
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"invitations","action":"SendEmailForInvitation","email":"SendEmailForInvitation","invite_token":"eb170ea62a4448b4a609c0521fbb4cf9"}
// ////////////////////////////////////////////////////////////////////////////////////////
func SendEmailForInvitation(target_email string, invite_token string) string {

	body := strings.NewReader(`{"email":"` + target_email + `","name":"Hi There!","content":"'https://puppychow.commhubapi.com/?tok=` + invite_token + `'","subject":"Your Invite to commhub"}`)

	//body := strings.NewReader(`{"email":"edward.anderson@commhubstuff.com","name":"Lane","content":"como estas?","subject":"Hola Lane!"}`)

	req, err := http.NewRequest("POST", "https://5jy9za0mr2.execute-api.us-east-1.amazonaws.com/dev", body)
	if err != nil {
		return "failed"
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "failed"
	}
	defer resp.Body.Close()

	return "success"

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Accept an invitation to an workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"invitations","action":"AcceptInviteToWorkspace","end_user_token":"2325418da6fb11e9a58342010a8e0121","invite_token":"989c12a1b51a11e9a58342010a8e0121"}
// ////////////////////////////////////////////////////////////////////////////////////////
func AcceptInviteToWorkspace(payload string) string {

	type inputData struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		InvToken    string `json:"invite_token"`
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

	query := "CALL commhub_junction.accept_invite_to_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.InvToken) + "\")"

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
// Reject an invitation to an workspace
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"invitations","action":"RejectInviteToWorkspace","end_user_token":"2325418da6fb11e9a58342010a8e0121","invite_token":"989c12a1b51a11e9a58342010a8e0121"}
// ////////////////////////////////////////////////////////////////////////////////////////
func RejectInviteToWorkspace(payload string) string {

	type inputData struct {
		Token       string `json:"api_token"`
		Signal      string `json:"signal"`
		Action      string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		InvToken    string `json:"invite_token"`
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

	query := "CALL commhub_junction.reject_invite_to_workspace(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.InvToken) + "\")"

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
