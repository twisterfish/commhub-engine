DELIMITER $$

DROP PROCEDURE commhub_junction.remove_member_from_workspace;
CREATE PROCEDURE commhub_junction.remove_member_from_workspace(
	IN in_end_user_token VARCHAR(32),
	IN in_target_member_token  VARCHAR(32),
	IN in_workspace_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE target_member_id INT;
	DECLARE target_member_workspace_permission_id INT;
	DECLARE target_workspace_id INT;
	DECLARE requesters_workspace_permission_id INT;
	DECLARE requesters_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status;
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
		WorkSpaceMmLkp.active = true;
		
-- if nothing found then user doesn't have permissions to do so	   
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid_requester_is_not_a_member";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
		
-- get the target member's internal ID if valid
	SELECT 
		WorkSpaceMmLkp.member_id, 
		PrmLvl.workspace_permission_id
	INTO
		target_member_id,
		target_member_workspace_permission_id
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
		EnUsr.end_user_token = in_target_member_token;

-- if nothing found then the target member isn't in the org	
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid_target_is_not_a_member";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;

-- requester must be at dispatch level and target can't be any higher	
	IF requesters_workspace_permission_id > 200 AND target_member_workspace_permission_id < 200 THEN	 
		SET out_status = "invalid_requester_is_not_authorized";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
	   
-- all good 				
	UPDATE commhub_junction.workspace_members_lkp SET active = false WHERE workspace_id = target_workspace_id AND member_id = target_member_id;
	SET out_status = "success";
	SELECT out_status;
  
END$$

DELIMITER ;