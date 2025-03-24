DELIMITER $$

DROP PROCEDURE commhub_junction.set_workspace_description;
CREATE PROCEDURE commhub_junction.set_workspace_description(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_new_description TEXT
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE requesters_id INT;
	DECLARE target_workspace_id, requesters_workspace_permission_id INT;
		
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
		PrmLvl.workspace_permission_id <= 200
	AND 
		WorkSpaceMmLkp.active = true;
		
-- if nothing found then user doesn't have permissions to do so	   
	IF found_rows() = 0 THEN
	 
		SET out_status = "invalid_user_permissions_level";
		SELECT out_status;
		
		LEAVE proc_label;
	END IF;

			   
-- all good 				
	UPDATE commhub_junction.workspace SET workspace_description = in_new_description WHERE workspace_id = target_workspace_id;
	SET out_status = "success";
	SELECT out_status;
  
END$$

DELIMITER ;