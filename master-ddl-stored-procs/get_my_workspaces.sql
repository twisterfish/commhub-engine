DELIMITER $$

DROP PROCEDURE commhub_junction.get_my_workspaces;
CREATE PROCEDURE commhub_junction.get_my_workspaces(
	IN in_end_user_token VARCHAR(32)
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

-- pull a list of orgs that this user is active in				
	SELECT 
		WorkSpaceMain.workspace_token, 
		PrmLvl.workspace_permission, 
		IFNULL(WorkSpaceMain.workspace_name,''), 
		IFNULL(WorkSpaceMain.workspace_description,'') 
	FROM 
		commhub_junction.workspace_members_lkp AS WorkSpaceMmLkp 
	INNER JOIN commhub_junction.end_user AS EnUsr ON 
		WorkSpaceMmLkp.member_id = EnUsr.end_user_id 
	INNER JOIN commhub_junction.workspace AS WorkSpaceMain ON 
		WorkSpaceMmLkp.workspace_id = WorkSpaceMain.workspace_id 
	INNER JOIN commhub_junction.workspace_permission AS PrmLvl ON 
		WorkSpaceMmLkp.workspace_permission_id = PrmLvl.workspace_permission_id 
	WHERE 
		EnUsr.end_user_token = in_end_user_token 
	AND 
		WorkSpaceMmLkp.active = true;
  
END$$

DELIMITER ;