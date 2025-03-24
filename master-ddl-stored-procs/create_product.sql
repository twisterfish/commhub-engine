DELIMITER $$
DROP PROCEDURE commhub_junction.create_product;

CREATE PROCEDURE commhub_junction.create_product(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_vendor_id INT,
	IN in_product_unit_sold_id INT,
	IN in_upc VARCHAR(64),
	IN in_sku VARCHAR(64),
	IN in_description TEXT
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id INT;
	DECLARE out_status VARCHAR(255);
	DECLARE out_new_product_id INT;
    
	
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; 
		SET out_new_product_id = 0;
		
		SELECT out_status, out_new_product_id;
		
		ROLLBACK;
	END;
	
	
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
		 
	IF found_rows() = 0 THEN
		SET out_status = "invalid_user_permissions_level";
		SELECT out_status;
		LEAVE proc_label;
	END IF;

	INSERT INTO commhub_junction.product (workspace_id, vendor_id, product_unit_sold_id, upc, sku, description)
		VALUES( target_workspace_id, in_vendor_id, in_product_unit_sold_id, in_upc, in_sku, in_description);
  
-- Return Last insterted ID and "success" message
	SET out_new_product_id = LAST_INSERT_ID();
	SET out_status = "success";
	SELECT out_status, out_new_product_id;
END$$
DELIMITER ;