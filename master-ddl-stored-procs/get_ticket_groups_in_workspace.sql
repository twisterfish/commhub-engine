DELIMITER $$

DROP PROCEDURE commhub_junction.get_ticket_groups_in_workspace;
CREATE PROCEDURE commhub_junction.get_ticket_groups_in_workspace(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE requesters_id INT;
		
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error;
		
		SELECT out_status;
	END;
	
-- find out if the user has the authority to do this		
	SELECT 
		WorkSpaceMmLkp.member_id
	INTO
		requesters_id
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
	 
		SET out_status = "invalid_user_permissions_level";
		SELECT out_status;
		
		LEAVE proc_label;
	END IF;

-- this is raw and needs to be secured - puppy chow		
	SELECT
		WorkSpace.workspace_token, 
		IFNULL(TicketGroup.ticket_group_id,''),
		IFNULL(TicketGroup.ticket_group_name,''),
		IFNULL(TicketGroup.ticket_group_description,'')
	FROM 
		commhub_junction.ticket_group AS TicketGroup
	INNER JOIN commhub_junction.workspace AS WorkSpace ON
		WorkSpace.workspace_id = TicketGroup.workspace_id
	WHERE 
		WorkSpace.workspace_token = in_workspace_token;
  
END$$

DELIMITER ;