DELIMITER $$

DROP PROCEDURE commhub_junction.invite_member_to_workspace;
CREATE PROCEDURE commhub_junction.invite_member_to_workspace(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_workspace_permission_id INT,
	IN in_target_email VARCHAR(255)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE invite_token VARCHAR(32);
	DECLARE target_workspace_id INT;
	DECLARE requesters_workspace_permission_id INT;
	DECLARE requesters_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status, "error";
	END;
	
-- find out if the user has the authority to do this		
	SELECT 
		WorkSpaceMmLkp.member_id, 
		PrmLvl.workspace_permission_id,
		WorkSpaceMain.workspace_id
	INTO
		requesters_id,
		requesters_workspace_permission_id,
		target_workspace_id
	FROM 
		commhub_junction.workspace_members_lkp AS WorkSpaceMmLkp 
	INNER JOIN commhub_junction.end_user AS EnUsr ON 
		WorkSpaceMmLkp.member_id = EnUsr.end_user_id
	INNER JOIN commhub_junction.workspace_permission AS PrmLvl ON 
		WorkSpaceMmLkp.workspace_permission_id = PrmLvl.workspace_permission_id 
	INNER JOIN commhub_junction.workspace AS WorkSpaceMain ON 
		WorkSpaceMmLkp.workspace_id = WorkSpaceMain.workspace_id 
	WHERE 
		WorkSpaceMain.workspace_token = in_workspace_token 
	AND
		EnUsr.end_user_token = in_end_user_token
	AND 
		PrmLvl.workspace_permission_id <= 200
	AND 
		WorkSpaceMmLkp.active = true;
		
-- if nothing found then user doesn't have permissions to do so	   
	IF found_rows() = 0 THEN
	 
		SET out_status = "invalid_not_authorized";
		SELECT out_status, "invalid";
		
		LEAVE proc_label;
	END IF;

-- requester cannot grant permissions that are at a greater level than their own	
	IF requesters_workspace_permission_id > in_workspace_permission_id THEN
	 
		SET out_status = "invalid_cannot_grant";
		SELECT out_status, "invalid";
		
		LEAVE proc_label;
	END IF;
	
	SET invite_token = REPLACE( UUID(),"-","");
	   
-- all good 
		
	INSERT INTO commhub_junction.invite_to_workspace ( invite_token, end_user_id, target_ws_id, target_ws_perm_id, invite_status_id, target_email ) VALUES ( invite_token, requesters_id, target_workspace_id, in_workspace_permission_id, 1, in_target_email );
	SET out_status = "success";
	SELECT out_status, invite_token;
  
END$$

DELIMITER ;