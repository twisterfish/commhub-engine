DELIMITER $$

DROP PROCEDURE commhub_junction.set_ticket_start_time;
CREATE PROCEDURE commhub_junction.set_ticket_start_time(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_ticket_id INT,
	IN in_start_time INT,
	IN in_ts_lat FLOAT(10,6),
	IN in_ts_long FLOAT(10,6)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
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
		Member.member_id
	INTO
		requesters_id
	FROM 
		commhub_junction.work_ticket AS Ticket
	INNER JOIN commhub_junction.ticket_group AS TicketGroup ON
		Ticket.ticket_group_id = TicketGroup.ticket_group_id
	INNER JOIN commhub_junction.workspace AS WorkSpace ON
		WorkSpace.workspace_id = TicketGroup.workspace_id
	INNER JOIN commhub_junction.workspace_members_lkp AS Member ON
		Member.workspace_id = WorkSpace.workspace_id
	INNER JOIN commhub_junction.end_user AS EndUser ON
		EndUser.end_user_id = Member.member_id
	WHERE 
		Ticket.work_ticket_id = in_ticket_id
	AND
		EndUser.end_user_token = in_end_user_token
	AND
		WorkSpace.workspace_token = in_workspace_token
	AND 
		Member.workspace_permission_id <= 400
	AND 
		Member.active = true;
		
-- if nothing found then user isn't authorized
	IF FOUND_ROWS() = 0 THEN	 
		SET out_status = "invalid_not_authorized";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
			   
-- all good so far				
	UPDATE commhub_junction.work_ticket SET work_ticket_status_id = 2, time_started = FROM_UNIXTIME( in_start_time ), ts_latitude = in_ts_lat, ts_longitude = in_ts_long WHERE work_ticket_id = in_ticket_id AND assigned_to_user_id = requesters_id;
	
-- if nothing changed then requester isn't the one working on the ticket and can't set a start time
	IF ROW_COUNT() = 0 THEN	 
		SET out_status = "invalid_not_authorized";
		SELECT out_status;		
		LEAVE proc_label;
	END IF;
	
	
	SET out_status = "success";
	SELECT out_status;
  
END$$

DELIMITER ;