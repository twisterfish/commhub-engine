DELIMITER $$

DROP PROCEDURE commhub_junction.create_user_profile;
CREATE PROCEDURE commhub_junction.create_user_profile(
	IN in_emid varchar(255),
  	IN in_pwid varchar(128),
  	IN in_ein_tax_id varchar(255),
  	IN in_ssn_tax_id varchar(255),
	IN in_last_name varchar(255),
	IN in_middle_name varchar(255),
	IN in_first_name varchar(255),
	IN in_address1 varchar(255),
	IN in_address2 varchar(255),
	IN in_city varchar(255),
	IN in_province_state varchar(32),
	IN in_zip_postal_code varchar(16),
	IN in_email varchar(255),
	IN in_phone varchar(255),
	IN in_country_code varchar(32)
    )
proc_label:BEGIN
	DECLARE new_end_user_id, new_workspace_id, new_group_id INT;
	DECLARE new_end_user_uuid VARCHAR(32);
	DECLARE new_workspace_uuid VARCHAR(32);
	DECLARE out_status VARCHAR(255);
	DECLARE out_new_user_uuid VARCHAR(32);
	DECLARE out_new_user_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
    GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		SET out_new_user_id = 0;
		SET out_new_user_uuid = "null";
		
		SELECT out_status, out_new_user_id, out_new_user_uuid;
		
		ROLLBACK;
	END;

    
-- create the end user's authorization record
	SET new_end_user_uuid = REPLACE( UUID(),"-","");

START TRANSACTION; 	           
	INSERT INTO commhub_junction.end_user( end_user_token, emid, pwid ) 
		VALUES ( new_end_user_uuid, in_emid, MD5(in_pwid) );

	SET new_end_user_id = LAST_INSERT_ID();

-- create the end user's profile record	
	INSERT INTO commhub_junction.end_user_profile ( end_user_id, ein_tax_id, ssn_tax_id, last_name, middle_name, first_name, address1, address2, city, province_state, zip_postal_code, email, phone, country_code ) 
		VALUES(	new_end_user_id, in_ein_tax_id, in_ssn_tax_id, in_last_name, in_middle_name, in_first_name, in_address1, in_address2, in_city, in_province_state, in_zip_postal_code, in_email, in_phone, in_country_code );

-- *NG* create default workspace
	SET new_workspace_uuid = REPLACE( UUID(),"-","");
	           
	INSERT INTO commhub_junction.workspace( workspace_creator_id, workspace_owner_id, workspace_token, workspace_name, workspace_description ) 
		VALUES ( new_end_user_id, new_end_user_id, new_workspace_uuid, 'My First Workspace', 'Default Initial Workspace' );
	
	SET new_workspace_id = LAST_INSERT_ID();

-- *NG* add the owner into the workspace as an admin	
	INSERT INTO commhub_junction.workspace_members_lkp ( workspace_id, member_id, workspace_permission_id, active ) 
		VALUES ( new_workspace_id, new_end_user_id, 100, true);

-- *NG* create the ticket group	
    INSERT INTO commhub_junction.ticket_group ( ticket_group_creator_id, workspace_id, ticket_group_name, ticket_group_description ) 
    	VALUES ( new_end_user_id, new_workspace_id, 'STAGING', 'Tickets that are not yet assigned to a group or were rejected by the ticket assignee' );
 
	SET new_group_id = LAST_INSERT_ID();

-- *NG* update the workspace with the new ticket staging group ID        
	UPDATE  commhub_junction.workspace 
		SET 
			staging_group_id = new_group_id 
		WHERE 
			workspace_id = new_workspace_id;

	COMMIT;
			
	SET out_status = "success";
	SET out_new_user_id = new_end_user_id;
	SET out_new_user_uuid = new_end_user_uuid;	
	SELECT out_status, out_new_user_id, out_new_user_uuid;
  
END$$

DELIMITER ;