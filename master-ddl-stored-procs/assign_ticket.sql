DELIMITER $$

DROP PROCEDURE commhub_junction.assign_ticket;
CREATE PROCEDURE commhub_junction.assign_ticket(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_ticket_group_id INT,
	IN in_work_ticket_id INT,
	IN in_assignee_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE requesters_id, assignee_id, requesters_workspace_permission_id, target_workspace_id, target_group_id, cur_ticket_org INT;
	DECLARE out_status VARCHAR(255);

	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status;
		
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
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
		
-- requester must be at dispatch level at least	
	IF requesters_workspace_permission_id > 200 THEN	 
		SET out_status = "invalid_requester_is_not_authorized";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
		
-- find out if the assignee is a member of the group		
	SELECT 
		WorkSpaceMmLkp.member_id 
	INTO
		assignee_id
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
		EnUsr.end_user_token = in_assignee_token
	AND 
		WorkSpaceMmLkp.active = true;
		
-- if nothing found then assignee isn't even a member of the group	   
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid_assignee_is_not_a_member";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;

-- find out if the user is trying to go outside the org they have permissions for
	SELECT workspace_id INTO cur_ticket_org FROM commhub_junction.ticket_group WHERE ticket_group_id = in_ticket_group_id;
	
-- if nothing found then ticket isn't even in the same workspace the user has permission for	   
	IF found_rows() = 0 OR cur_ticket_org <> target_workspace_id THEN	 
		SET out_status = "invalid_not_authorized";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;		

-- assign the ticket to another user only if they haven't started any work on it
	UPDATE commhub_junction.work_ticket SET assigned_to_user_id = assignee_id, assigned_by_user_id = requesters_id WHERE work_ticket_id = in_work_ticket_id AND ticket_group_id = in_ticket_group_id AND ( UNIX_TIMESTAMP( time_started ) = 1 );

-- if nothing updated then the ticket has been started already   
	IF row_count() = 0 THEN	 
		SET out_status = "invalid_ticket_has_already_been_started";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
	
	SET out_status = "success";
	SELECT out_status;	
	  
END$$

DELIMITER ;