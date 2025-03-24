DELIMITER $$

DROP PROCEDURE commhub_junction.get_members_in_my_workspace;
CREATE PROCEDURE commhub_junction.get_members_in_my_workspace(
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
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status,"","","","";	
	END;
	
-- find out if the user is a member of this workspace and has the visibility		
	SELECT 
		WorkSpaceMmLkp.member_id
	INTO
		requesters_id
	FROM 
		commhub_junction.workspace_members_lkp AS WorkSpaceMmLkp 
	INNER JOIN commhub_junction.end_user AS EnUsr ON 
		WorkSpaceMmLkp.member_id = EnUsr.end_user_id
	INNER JOIN commhub_junction.workspace AS WorkSpaceMain ON 
		WorkSpaceMmLkp.workspace_id = WorkSpaceMain.workspace_id 
	WHERE 
		WorkSpaceMain.workspace_token = in_workspace_token 
	AND
		EnUsr.end_user_token = in_end_user_token
	AND 
		WorkSpaceMmLkp.active = true;
		
-- if nothing found then user isn't an active member of the group - bug out	   
	IF found_rows() = 0 THEN
		LEAVE proc_label;
	END IF;
	
-- all good 				
	SELECT 
		EnUsr.end_user_token, 
		PrmLvl.workspace_permission, 
		EnUsrPrf.first_name, 
		EnUsrPrf.last_name,
		EnUsrPrf.phone
	FROM 
		commhub_junction.workspace_members_lkp AS WorkSpaceMmLkp 
	INNER JOIN commhub_junction.end_user AS EnUsr ON 
		WorkSpaceMmLkp.member_id = EnUsr.end_user_id 
	INNER JOIN commhub_junction.end_user_profile AS EnUsrPrf ON 
		WorkSpaceMmLkp.member_id = EnUsrPrf.end_user_id
	INNER JOIN commhub_junction.workspace AS WorkSpaceMain ON 
		WorkSpaceMmLkp.workspace_id = WorkSpaceMain.workspace_id 
	INNER JOIN commhub_junction.workspace_permission AS PrmLvl ON 
		WorkSpaceMmLkp.workspace_permission_id = PrmLvl.workspace_permission_id 
	WHERE 
		WorkSpaceMain.workspace_token = in_workspace_token 
	AND 
		WorkSpaceMmLkp.active = true;
  
END$$

DELIMITER ;