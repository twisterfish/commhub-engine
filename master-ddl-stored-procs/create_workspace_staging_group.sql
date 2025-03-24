DELIMITER $$

DROP PROCEDURE commhub_junction.create_workspace_staging_group;
CREATE PROCEDURE commhub_junction.create_workspace_staging_group(
	IN in_creator_token VARCHAR(32),
	IN in_workspace_name VARCHAR(64),
	IN in_workspace_desc TEXT
    )
proc_label:BEGIN
	DECLARE creator_id, new_workspace_id, new_group_id, out_new_group_id INT;
	DECLARE new_workspace_uuid VARCHAR(32);
	DECLARE out_status VARCHAR(255);
	DECLARE out_new_workspace_uuid VARCHAR(32);
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		SET out_new_workspace_uuid = "null";
		SET out_new_group_id = 0;
		
		SELECT out_status, out_new_workspace_uuid, out_new_group_id;
		
		ROLLBACK;
	END;
	
-- get the group creator's internal ID
	SELECT end_user_id INTO creator_id FROM end_user WHERE end_user_token = in_creator_token;
	IF found_rows() = 0 THEN 
		SET out_status = "invalid_user"; -- this is set for debug while in dev
		SET out_new_workspace_uuid = "null";
		SET out_new_group_id = 0;
		
		SELECT out_status, out_new_workspace_uuid, out_new_group_id;
		LEAVE proc_label;
	END IF;
    
-- create the workspace
	SET new_workspace_uuid = REPLACE( UUID(),"-","");

	START TRANSACTION;
	           
	INSERT INTO commhub_junction.workspace( workspace_creator_id, workspace_owner_id, workspace_token, workspace_name, workspace_description ) 
		VALUES ( creator_id, creator_id, new_workspace_uuid, in_workspace_name, in_workspace_desc );
	
	SET new_workspace_id = LAST_INSERT_ID();

-- add the owner into the workspace as an admin	
	INSERT INTO commhub_junction.workspace_members_lkp ( workspace_id, member_id, workspace_permission_id, active ) 
		VALUES ( new_workspace_id, creator_id, 100, true);

-- create the ticket group	
    INSERT INTO commhub_junction.ticket_group ( ticket_group_creator_id, workspace_id, ticket_group_name, ticket_group_description ) 
    	VALUES ( creator_id, new_workspace_id, 'STAGING', 'Tickets that are not yet assigned to a group or were rejected by the ticket assignee' );
   
	SET new_group_id = LAST_INSERT_ID();

-- update the workspace with the new ticket staging group ID        
	UPDATE  commhub_junction.workspace 
		SET 
			staging_group_id = new_group_id 
		WHERE 
			workspace_id = new_workspace_id;
	COMMIT;

	SET out_status = "success";
	SET out_new_workspace_uuid = new_workspace_uuid;
	SET out_new_group_id = new_group_id;
	
	SELECT out_status, out_new_workspace_uuid, out_new_group_id;
  
END$$

DELIMITER ;