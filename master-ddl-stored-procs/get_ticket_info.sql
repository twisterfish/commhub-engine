DELIMITER $$

DROP PROCEDURE commhub_junction.get_ticket_info;
CREATE PROCEDURE commhub_junction.get_ticket_info(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_ticket_guid  VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
		
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status;
	END;

-- find out if the user has the authority to do this and get the ticket info if he does and nothing if not		
	SELECT
		WorkSpace.workspace_token, 
		Ticket.work_ticket_id,
		Ticket.local_guid,
		Ticket.ticket_group_id,
		TicketStatus.description,
		Ticket.work_ticket_type_id,
		Ticket.real_property_id,
		Ticket.running_time,
		UNIX_TIMESTAMP(Ticket.time_created),
		Ticket.time_started,
		Ticket.ts_latitude,
		Ticket.ts_longitude,
		UNIX_TIMESTAMP(Ticket.time_finished),
		Ticket.tf_latitude,
		Ticket.tf_longitude,
		IFNULL(Ticket.title,''),
		IFNULL(Ticket.description,'')
	FROM 
		commhub_junction.work_ticket AS Ticket
	INNER JOIN commhub_junction.work_ticket_status AS TicketStatus ON
		Ticket.work_ticket_status_id = TicketStatus.work_ticket_status_id
	INNER JOIN commhub_junction.ticket_group AS TicketGroup ON
		Ticket.ticket_group_id = TicketGroup.ticket_group_id
	INNER JOIN commhub_junction.workspace AS WorkSpace ON
		WorkSpace.workspace_id = TicketGroup.workspace_id
	INNER JOIN commhub_junction.workspace_members_lkp AS Member ON
		Member.workspace_id = WorkSpace.workspace_id
	INNER JOIN commhub_junction.end_user AS EndUser ON
		EndUser.end_user_id = Member.member_id
	WHERE 
		Ticket.local_guid = in_ticket_guid
	AND
		EndUser.end_user_token = in_end_user_token
	AND
		WorkSpace.workspace_token = in_workspace_token
	AND 
		Member.workspace_permission_id <= 400
	AND 
		Member.active = true;
  
END$$

DELIMITER ;