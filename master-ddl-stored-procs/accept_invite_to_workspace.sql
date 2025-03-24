DELIMITER $$

DROP PROCEDURE commhub_junction.accept_invite_to_workspace;
CREATE PROCEDURE commhub_junction.accept_invite_to_workspace(
	IN in_end_user_token VARCHAR(32),
	IN in_invite_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE out_target_ws_id, out_target_ws_permission_id, out_invitation_id, out_target_user_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status;
		ROLLBACK;
	END;
	
-- find out if this is a legit invite		
	SELECT 
    	invite_id,
    	target_ws_id,
    	target_ws_perm_id
	INTO 
		out_invitation_id,
		out_target_ws_id,
		out_target_ws_permission_id
	FROM 
		commhub_junction.invite_to_workspace
	WHERE 
		invite_token = in_invite_token
	AND
		invite_status_id = 1;

		
-- if nothing found then user doesn't have an invite - bug out	   
	IF found_rows() = 0 THEN
		SET out_status = "no invitation found";
		SELECT out_status;
		LEAVE proc_label;
	END IF;

-- find out if this is a legit invited user		
	SELECT 
    	end_user_id
	INTO 
		out_target_user_id
	FROM 
		commhub_junction.end_user
	WHERE 
		end_user_token = in_end_user_token;
	
-- if nothing found then user doesn't have an invite - bug out	   
	IF found_rows() = 0 THEN
		SET out_status = "user not found";
		SELECT out_status;
		LEAVE proc_label;
	END IF;
		   
-- all good 		
	UPDATE commhub_junction.invite_to_workspace SET target_user_id = out_target_user_id, invite_status_id = 2 
	WHERE invite_token = in_invite_token;
	
	INSERT INTO commhub_junction.workspace_members_lkp ( workspace_id, member_id, workspace_permission_id, active ) VALUES ( out_target_ws_id, out_target_user_id, out_target_ws_permission_id, true );
			
	SET out_status = "success";
	SELECT out_status;
  
END$$

DELIMITER ;