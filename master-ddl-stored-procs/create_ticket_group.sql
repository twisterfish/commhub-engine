DELIMITER $$

DROP PROCEDURE commhub_junction.create_ticket_group;
CREATE PROCEDURE commhub_junction.create_ticket_group(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_group_name VARCHAR(64),
	IN in_group_desc TEXT
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id, new_group_id INT;
	DECLARE out_status VARCHAR(255);
	DECLARE out_new_group_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		SET out_new_group_id = 0;
		
		SELECT out_status, out_new_group_id;
		
		ROLLBACK;
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
		
-- if nothing found then user isn't even a member of the group	   
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid_requester_is_not_a_member";
		SET out_new_group_id = 0;		
		SELECT out_status, out_new_group_id;		
		LEAVE proc_label;
	END IF;
		
-- requester must be at dispatch level at least	
	IF requesters_workspace_permission_id > 200 THEN	 
		SET out_status = "invalid_requester_is_not_authorized";
		SET out_new_group_id = 0;		
		SELECT out_status, out_new_group_id;		
		LEAVE proc_label;
	END IF;

-- create the ticket group
	
    INSERT INTO commhub_junction.ticket_group ( ticket_group_creator_id, workspace_id, ticket_group_name, ticket_group_description ) 
    	VALUES ( requesters_id, target_workspace_id, in_group_name, in_group_desc );
   
	SET new_group_id = LAST_INSERT_ID();
			
	SET out_status = "success";
	SET out_new_group_id = new_group_id;
	
	SELECT out_status, out_new_group_id;
  
END$$

DELIMITER ;