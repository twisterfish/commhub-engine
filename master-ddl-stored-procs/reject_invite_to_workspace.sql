DELIMITER $$

DROP PROCEDURE commhub_junction.reject_invite_to_workspace;
CREATE PROCEDURE commhub_junction.reject_invite_to_workspace(
	IN in_end_user_token VARCHAR(32),
	IN in_invite_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE out_target_workspace_id, out_target_workspace_permission_id, out_invitation_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error; -- this is set for debug while in dev
		SELECT out_status;
		ROLLBACK;
		
		-- INSERT INTO commhub_junction.app_sp_errors ( error_description ) VALUES out_status;
		
	END;
	
-- find out if this is a legit invite		
	SELECT 
    	invite_id,
    	target_ws_id,
    	target_ws_perm_id
	INTO 
		out_invitation_id,
		out_target_workspace_id,
		out_target_workspace_permission_id
	FROM 
		commhub_junction.invite_to_workspace
	WHERE 
		invite_token = in_invite_token
	AND
		invite_status_id = 1;

		
-- if nothing found then user doesn't have an invite - bug out	   
	IF found_rows() = 0 THEN
		SET out_status = "invalid not found";
		SELECT out_status;
		LEAVE proc_label;
	END IF;
		   
-- all good - set the invite to rejected status		
	UPDATE commhub_junction.invite_to_workspace SET invite_status_id = 3 
		WHERE invite_id = out_invitation_id;
				
	SET out_status = "success";
	SELECT out_status;
  
END$$

DELIMITER ;