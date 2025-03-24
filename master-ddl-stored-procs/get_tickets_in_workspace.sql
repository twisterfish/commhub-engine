DELIMITER $$

DROP PROCEDURE commhub_junction.get_tickets_in_workspace;
CREATE PROCEDURE commhub_junction.get_tickets_in_workspace(
	IN in_workspace_token VARCHAR(32)
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

-- find out all the tickets that he is assigned and in groups in which he is an active member		
		SELECT 
		WorkSpace.workspace_token,
		Ticket.work_ticket_id,
		Ticket.local_guid,
		Ticket.ticket_group_id,
		Ticket.assigned_to_user_id,
		Ticket.assigned_by_user_id,
		Ticket.creator_user_id,
		Ticket.work_ticket_status_id,
		Ticket.work_ticket_type_id,
		Ticket.real_property_id,
		Ticket.running_time,
		UNIX_TIMESTAMP(Ticket.time_created),
		UNIX_TIMESTAMP(Ticket.time_started),
		Ticket.ts_latitude,
		Ticket.ts_longitude,
		UNIX_TIMESTAMP(Ticket.time_finished),
		Ticket.tf_latitude,
		Ticket.tf_longitude,
		IFNULL(Ticket.title,''),
		IFNULL(Ticket.description,'')
	FROM 
		commhub_junction.work_ticket AS Ticket
	INNER JOIN commhub_junction.ticket_group AS TicketGroup ON
		Ticket.ticket_group_id = TicketGroup.ticket_group_id
	INNER JOIN commhub_junction.workspace AS WorkSpace ON
		TicketGroup.workspace_id = WorkSpace.workspace_id
	WHERE 
		WorkSpace.workspace_token = in_workspace_token;
  
END$$

DELIMITER ;