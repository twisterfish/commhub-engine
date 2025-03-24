
DELIMITER $$
DROP PROCEDURE commhub_junction.set_product_count;

CREATE PROCEDURE commhub_junction.set_product_count(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
    IN in_product_id INT,
	IN in_total_quantity_on_hand INT
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id INT;
	DECLARE out_status VARCHAR(255);
    
	
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		SET out_status = @full_error; 
		SELECT out_status;
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
	INNER JOIN commhub_junction.product AS prdct ON 
		WorkSpaceMmLkp.workspace_id = prdct.workspace_id 
	WHERE 
		WorkSpaceMain.workspace_token = in_workspace_token 
	AND
		EnUsr.end_user_token = in_end_user_token
	AND 
		prdct.product_id = in_product_id
	AND 
		PrmLvl.workspace_permission_id <= 300
	AND 
		WorkSpaceMmLkp.active = true;
		 
	IF found_rows() = 0 THEN
		SET out_status = "invalid_user_permissions_level";
		SELECT out_status;
		LEAVE proc_label;
	END IF;

	UPDATE commhub_junction.product_inventory SET total_quantity_on_hand = in_total_quantity_on_hand WHERE product_id = in_product_id;

	IF row_count() = 0 THEN	 
			SET out_status = "invalid"; -- this is set for debug while in dev
			SELECT out_status;
			LEAVE proc_label;
	END IF;

-- Return  "success" message
	SET out_status = "success";
	SELECT out_status;
END$$
DELIMITER ;