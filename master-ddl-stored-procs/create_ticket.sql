DELIMITER $$

DROP PROCEDURE commhub_junction.create_ticket;
CREATE PROCEDURE commhub_junction.create_ticket(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_work_ticket_type_id INT,
	IN in_local_guid VARCHAR(32),
	IN in_real_property_id INT,
	IN in_title VARCHAR(255),
	IN in_description TEXT
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id, staging_group_id INT;
	DECLARE out_status VARCHAR(255);
	DECLARE out_new_ticket_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		SET out_new_ticket_id = 0;
		
		SELECT out_status, out_new_ticket_id;
		
		ROLLBACK;
	END;
	
-- find out if the user has the authority to do this		
	SELECT 
		WorkSpaceMmLkp.member_id, 
		PrmLvl.workspace_permission_id,
		WorkSpaceMain.workspace_id,
		WorkSpaceMain.staging_group_id
	INTO
		requesters_id,
		requesters_workspace_permission_id,
		target_workspace_id,
		staging_group_id
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
		SET out_new_ticket_id = 0;		
		SELECT out_status, out_new_ticket_id;		
		LEAVE proc_label;
	END IF;
		
-- requester must be at worker level at least	
	IF requesters_workspace_permission_id > 400 THEN	 
		SET out_status = "invalid_requester_is_not_authorized";
		SET out_new_ticket_id = 0;		
		SELECT out_status, out_new_ticket_id;		
		LEAVE proc_label;
	END IF;

-- create the ticket
	INSERT INTO commhub_junction.work_ticket( work_ticket_status_id, work_ticket_type_id, ticket_group_id, local_guid, real_property_id, creator_user_id, assigned_to_user_id, running_time, time_created, time_started, time_finished, title, description ) 
		VALUES( 1, in_work_ticket_type_id, staging_group_id, in_local_guid, in_real_property_id, requesters_id, requesters_id, 0, now(), from_unixtime(1), from_unixtime(1), in_title, in_description );
  
	SET out_new_ticket_id = LAST_INSERT_ID();
	SET out_status = "success";
	SELECT out_status, out_new_ticket_id;
	
	  
END$$

DELIMITER ;