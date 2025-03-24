package organizations

//////////////////////////////////////////////////////////////////////////////////////////
// This package handles all signal requests for organizations
// THIS PACKAGE IS DEPRECATED!!!!
//////////////////////////////////////////////////////////////////////////////////////////

import (
    _"github.com/go-sql-driver/mysql"
    "encoding/json"
    "bytes"
    "strconv"
    "strings"
    "net/http"
    "datastores"
    "github.com/aws/aws-lambda-go/events"
)

//////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// move to permission labels and tags
//////////////////////////////////////////////////////////////////////////////////////////
func DoAction( signal string, action string, payload string, *events.APIGatewayProxyRequest )( string ) {
	
	switch action {
		case "CreateOrganization":
			return CreateOrganization( payload )
		case "CreateOrganizationOnDevice":
			return CreateOrganizationOnDevice( payload )
		case "GetOrganizationByOwner":
			return GetOrganizationByOwner( payload )
		case "AddMemberToOrg":
			return AddMemberToOrg( payload )
		case "RemoveMemberFromOrg":
			return RemoveMemberFromOrg( payload )
		case "GetMembersInMyOrganization":
			return GetMembersInMyOrganization( payload )
		case "GetMyOrganizations":
			return GetMyOrganizations( payload )
		case "SetOrgName":// use token
			return SetOrgName( payload )
		case "SetOrgDescription":// use token
			return SetOrgDescription( payload )
		case "CreateTicketGroup":
			return CreateTicketGroup( payload )
		case "SetTicketGroupName":
			return SetTicketGroupName( payload )
		case "SetTicketGroupDescription":
			return SetTicketGroupDescription( payload )
		default:
    		return  "{\"signal\":\"error\",\"action\":\"Your action request is invalid.\"}" 
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Enter a new organization in the data
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"CreateOrganization","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_name":"COOL COMPANY","organization_description":"Based out of Fort Lauderdale"}
//////////////////////////////////////////////////////////////////////////////////////////
func CreateOrganization( payload string )( string ) {
  
	type NewOrg struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		OrgCRTR string `json:"end_user_token"`
		OrgNM string `json:"organization_name"`
		OrgDESC string `json:"organization_description"`
	}

	// used for parsing request
	var input NewOrg
	var rows bytes.Buffer
    var prebuf bytes.Buffer
    var item_count int
    	
	pByte :=  []byte ( payload )
	
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
    	return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}
	
	type OrgRecord struct{
		OrgStatus string `json:"out_status"`
		OrgUUID string `json:"out_new_org_uuid"`
		GrpID string `json:"out_new_group_id"`
	}
	var output OrgRecord

    db, err := datastores.OpenRDS()
    
    // if there is an error opening the connection, handle it - don't close on bad connection?
    if err != nil {
        return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"        
    }
        
    // defer the close till after the main function has finished executing
    defer db.Close()
    item_count = 0
    
    query := "CALL commhub_junction.create_org_staging_group(\"" + strings.TrimSpace(input.OrgCRTR) + "\", \"" + strings.TrimSpace(input.OrgNM) + "\", \"" + strings.TrimSpace(input.OrgDESC) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.OrgStatus, &output.OrgUUID, &output.GrpID )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.OrgStatus) == "invalid_user"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.OrgStatus) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.OrgStatus) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.OrgStatus + "\"," )
			rows.WriteString("\"out_new_org_uuid\":\"" + output.OrgUUID + "\"," )
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Create a new organization with a device generated GUID
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"CreateOrganizationOnDevice","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"YC98oij345trghui45tergf9huie","organization_name":"COOL COMPANY","organization_description":"Based out of Fort Lauderdale"}
//////////////////////////////////////////////////////////////////////////////////////////
func CreateOrganizationOnDevice( payload string )( string ) {
  
	type NewOrg struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		OrgCRTR string `json:"end_user_token"`
		OrgTKN string `json:"organization_token"`
		OrgNM string `json:"organization_name"`
		OrgDESC string `json:"organization_description"`
	}

	// used for parsing request
	var input NewOrg
	var rows bytes.Buffer
    var prebuf bytes.Buffer
    var item_count int
    	
	pByte :=  []byte ( payload )
	
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
    	return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}
	
	type OrgRecord struct{
		OrgStatus string `json:"out_status"`
		OrgUUID string `json:"out_new_org_uuid"`
		GrpID string `json:"out_new_group_id"`
	}
	var output OrgRecord

    db, err := datastores.OpenRDS()
    
    // if there is an error opening the connection, handle it - don't close on bad connection?
    if err != nil {
        return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"        
    }
        
    // defer the close till after the main function has finished executing
    defer db.Close()
    item_count = 0
    
    query := "CALL commhub_junction.create_org_on_device(\"" + strings.TrimSpace(input.OrgCRTR) + "\", \"" + strings.TrimSpace(input.OrgTKN) + "\", \"" + strings.TrimSpace(input.OrgNM) + "\", \"" + strings.TrimSpace(input.OrgDESC) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.OrgStatus, &output.OrgUUID, &output.GrpID )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.OrgStatus) == "invalid_user"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.OrgStatus) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.OrgStatus) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.OrgStatus + "\"," )
			rows.WriteString("\"out_new_org_uuid\":\"" + output.OrgUUID + "\"," )
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull an organization by owner's token
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"GetOrganizationByOwner","org_owner_token":"8a985cfe86ce11e9bc42526af7764f64"}
//////////////////////////////////////////////////////////////////////////////////////////
func GetOrganizationByOwner( payload string )( string ) {

	type OrgRequest struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`	
		OwnToken string `json:"org_owner_token"`
	}
	// used for parsing request
	var input OrgRequest	
	pByte :=  []byte ( payload )
	
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
    	return "{\"token\":\"" + input.Token + "\",\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}
	
	type OrgRecord struct{
		OrgToken string `json:"organization_token"`
		OrgNM string `json:"organization_name"`
		OrgDESC string `json:"organization_description"`
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
    
    query := "SELECT ORG.organization_token, IFNULL(ORG.organization_name,''), IFNULL(ORG.organization_description,'') FROM commhub_junction.organization AS ORG INNER JOIN commhub_junction.end_user AS EndU ON ORG.organization_owner_id = EndU.end_user_id WHERE EndU.end_user_token = \"" + strings.TrimSpace(input.OwnToken) + "\""
	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			var tag OrgRecord
			// for each record, scan the result into our  struct
			err = results.Scan( &tag.OrgToken, &tag.OrgNM, &tag.OrgDESC )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"organization_token\":\"" + tag.OrgToken + "\"," )
			rows.WriteString("\"organization_name\":\"" + tag.OrgNM + "\"," )
			rows.WriteString("\"organization_description\":\"" + tag.OrgDESC + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
    		
}

//////////////////////////////////////////////////////////////////////////////////////////
// Add a user to an organization 
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"AddMemberToOrg","end_user_token":"8a985a9286ce11e9bc42526af7764f64","new_member_token":"8a9861cc86ce11e9bc42526af7764f64","organization_token":"17f0a2866e7e408d9ca3810dc801e046","org_permission_id":"200"}
//////////////////////////////////////////////////////////////////////////////////////////
func AddMemberToOrg( payload string )( string ) {
  
	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		NewMemToken string `json:"new_member_token"`
		OrgToken string `json:"organization_token"`		
		OrgPermID string `json:"org_permission_id"`
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}
	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.add_member_to_org(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.NewMemToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\", \"" + strings.TrimSpace(input.OrgPermID) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Remove a member from an organization 
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"RemoveMemberFromOrg","end_user_token":"8a985a9286ce11e9bc42526af7764f64","target_member_token":"8a9861cc86ce11e9bc42526af7764f64","organization_token":"17f0a2866e7e408d9ca3810dc801e046"}
//////////////////////////////////////////////////////////////////////////////////////////
func RemoveMemberFromOrg( payload string )( string ) {

	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		TrgMemToken string `json:"target_member_token"`
		OrgToken string `json:"organization_token"`		
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.remove_member_from_org(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.TrgMemToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull members in an organization by org token
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"GetMembersInMyOrganization","end_user_token":"2325418da6fb11e9a58342010a8e0121","organization_token":"17f0a2866e7e408d9ca3810dc801e046"}
//////////////////////////////////////////////////////////////////////////////////////////
func GetMembersInMyOrganization( payload string )( string ) {

	type OrgRequest struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`	
		OrgToken string `json:"organization_token"`
	}
	// used for parsing request
	var input OrgRequest	
	pByte :=  []byte ( payload )
	
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
    	return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}
	
	type OrgMemXREF struct{
		OrgMEM string `json:"end_user_token"`
		OrgLvl string `json:"org_permission"`
		MemFnm string `json:"first_name"`
		MemLnm string `json:"last_name"`
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
    
    query := "CALL commhub_junction.get_members_in_my_org(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\")"

	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			var tag OrgMemXREF
			// for each record, scan the result into our  struct
			err = results.Scan( &tag.OrgMEM, &tag.OrgLvl, &tag.MemFnm, &tag.MemLnm )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"end_user_token\":\"" + tag.OrgMEM + "\",")
			rows.WriteString("\"org_permission\":\"" + tag.OrgLvl + "\",")
			rows.WriteString("\"first_name\":\"" + tag.MemFnm + "\",")
			rows.WriteString("\"last_name\":\"" + tag.MemLnm + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
    		
}

//////////////////////////////////////////////////////////////////////////////////////////
// Pull organizations I'm in
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"GetMyOrganizations","end_user_token":"eb170ea62a4448b4a609c0521fbb4cf9"}
//////////////////////////////////////////////////////////////////////////////////////////
func GetMyOrganizations( payload string )( string ) {

	type OrgRequest struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`	
		EUTok string `json:"end_user_token"`
	}

	type OrgsMem struct{
		OrgToken string `json:"organization_token"`
		OrgLvl string `json:"org_permission"`
		OrgNM string `json:"organization_name"`
		OrgDESC string `json:"organization_description"`
	}
	
	var input OrgRequest
	var output OrgsMem
    var rows bytes.Buffer
    var prebuf bytes.Buffer
    var item_count int
    pByte :=  []byte ( payload )	
	
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
    
	query := "CALL commhub_junction.get_my_organizations(\"" + strings.TrimSpace(input.EUTok) + "\")"
	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			err = results.Scan( &output.OrgToken, &output.OrgLvl, &output.OrgNM, &output.OrgDESC )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			}

			rows.WriteString("{\"organization_token\":\"" + output.OrgToken + "\",")
			rows.WriteString("\"org_permission\":\"" + output.OrgLvl + "\",")
			rows.WriteString("\"organization_name\":\"" + output.OrgNM + "\",") 
			rows.WriteString("\"organization_description\":\"" + output.OrgDESC + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
    		
}

//////////////////////////////////////////////////////////////////////////////////////////
// Set or update a organization name - only org owners can do this
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"SetOrgName","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"eb170ea62a4448b4a609c0521fbb4cf9","organization_name":"COOL COMPANY"}
//////////////////////////////////////////////////////////////////////////////////////////
func SetOrgName( payload string )( string ) {

	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		OrgToken string `json:"organization_token"`
		OrgName string `json:"organization_name"`		
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.set_organization_name(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\", \"" + strings.TrimSpace(input.OrgName) + "\")"   
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Set or update an organization description - only the org owner can do this.
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"SetOrgDescription","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"eb170ea62a4448b4a609c0521fbb4cf9","organization_description":"COOL COMPANY"}
//////////////////////////////////////////////////////////////////////////////////////////
func SetOrgDescription( payload string )( string ) {

	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		OrgToken string `json:"organization_token"`
		OrgDesc string `json:"organization_description"`		
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.set_organization_description(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\", \"" + strings.TrimSpace(input.OrgDesc) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}	

//////////////////////////////////////////////////////////////////////////////////////////
// Enter a new ticket group in the data 
//{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"CreateTicketGroup","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_name":"B OCEAN","ticket_group_description":"Based out of Fort Lauderdale"}
//////////////////////////////////////////////////////////////////////////////////////////
func CreateTicketGroup( payload string )( string ) {

	type NewGroup struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EuTok string `json:"end_user_token"`
		OrgTok string `json:"organization_token"`
		GrpNM string `json:"ticket_group_name"`
		GrpDESC string `json:"ticket_group_description"`
	}

	// used for parsing request
	var input NewGroup
	var rows bytes.Buffer
    var prebuf bytes.Buffer
    var item_count int
    	
	pByte :=  []byte ( payload )
	
	if err := json.Unmarshal(pByte, &input); err != nil {
		// Malformed JSON - kick it back
    	return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
	}
	
	type GrpRecord struct{
		GrpStatus string `json:"out_status"`
		GrpID string `json:"out_new_group_id"`
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
    
    query := "CALL commhub_junction.create_ticket_group(\"" + strings.TrimSpace(input.EuTok) + "\", \"" + strings.TrimSpace(input.OrgTok) + "\", \"" + strings.TrimSpace(input.GrpNM) + "\", \"" + strings.TrimSpace(input.GrpDESC) + "\")"
    
	results, err := db.Query(query)

	if err != nil {
		results.Close()
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.GrpStatus, &output.GrpID )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.GrpStatus) == "invalid_user"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"invalid_user\"}"
			} else if strings.TrimSpace(output.GrpStatus) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.GrpStatus) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.GrpStatus + "\"," )
			rows.WriteString("\"out_new_group_id\":\"" + output.GrpID + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Set or update a ticket_group name in the data 
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"SetTicketGroupName","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_id":"1029","ticket_group_name":"Water Pipes"}
//////////////////////////////////////////////////////////////////////////////////////////
func SetTicketGroupName( payload string )( string ) {
	
	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		OrgToken string `json:"organization_token"`
		GrpID string `json:"ticket_group_id"`
		GrpName string `json:"ticket_group_name"`		
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.set_ticket_group_name(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\", \"" + strings.TrimSpace(input.GrpID) + "\", \"" + strings.TrimSpace(input.GrpName) + "\")"   
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

//////////////////////////////////////////////////////////////////////////////////////////
// Set or update a ticket_group name in the data 
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"organizations","action":"SetTicketGroupDescription","end_user_token":"8a985a9286ce11e9bc42526af7764f64","organization_token":"eb170ea62a4448b4a609c0521fbb4cf9","ticket_group_id":"1029","ticket_group_description":"B Ocean Hotel"}
//////////////////////////////////////////////////////////////////////////////////////////
func SetTicketGroupDescription( payload string )( string ) {
	
	type inputData struct{
		Token string `json:"api_token"`
		Signal string `json:"signal"`
		Action string `json:"action"`
		EndUsrToken string `json:"end_user_token"`
		OrgToken string `json:"organization_token"`
		GrpID string `json:"ticket_group_id"`
		GrpDesc string `json:"ticket_group_description"`		
	}
	
	type outputData struct{
		Status string `json:"out_status"`
	}	
	var input inputData
	var output outputData
	var rows bytes.Buffer
	var prebuf bytes.Buffer
    var item_count int

	pByte :=  []byte ( payload )

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
		
    query := "CALL commhub_junction.set_ticket_group_description(\"" + strings.TrimSpace(input.EndUsrToken) + "\", \"" + strings.TrimSpace(input.OrgToken) + "\", \"" + strings.TrimSpace(input.GrpID) + "\", \"" + strings.TrimSpace(input.GrpDesc) + "\")"   
	results, err := db.Query(query)

	if err != nil {
		results.Close()                                
        return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
	
		for results.Next() {
			
			// for each record, scan the result into our  struct
			err = results.Scan( &output.Status )
		
			if err != nil {
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}"
			} else if strings.TrimSpace(output.Status) != "success"{
				results.Close()		
        		return "{\"signal\":\"error\",\"action\":\"" + strings.TrimSpace(output.Status) + "\"}"
			}
			
			//log.Printf(tag.ID)
			rows.WriteString("{\"out_status\":\"" + output.Status + "\"},") 
		
			item_count++
		}
	}
	
	prebuf.WriteString("{\"signal\":\""+ input.Signal +"\",\"action\":\"" + input.Action + "\",\"item_count\":\"" + strconv.Itoa(item_count) + "\"," + "\"item_list\":[")		
	prebuf.WriteString ( strings.TrimRight(rows.String(),",") + "]}" )// take off the trailing comma in the set and cap it
   
    results.Close()
    
    return prebuf.String()
}

////////////////////////////////////////////////////////////////////////////////////////// 