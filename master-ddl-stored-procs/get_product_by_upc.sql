DELIMITER $$
DROP PROCEDURE commhub_junction.get_product_by_upc;

CREATE PROCEDURE commhub_junction.get_product_by_upc(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
    IN in_upc VARCHAR(64)
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id INT;
	DECLARE out_upc, out_sku, out_description VARCHAR(255);
	DECLARE out_product_id, out_vendor_id, out_product_unit_sold_id INT;
   
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		-- SET out_status = @full_error; 
		-- Declare & return Dummy output variables to avoid "result.Scan()" method errors in GoLang
		SET out_product_id = 0;
        SET out_upc = "NULL";
		SET out_sku = "NULL";
        SET out_description = "NULL";
		SET out_vendor_id = 0;
        SET out_product_unit_sold_id = 0;
		SELECT out_product_id,out_vendor_id, out_product_unit_sold_id, out_upc, out_sku, out_description;
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
		prdct.upc = in_upc
	AND 
		PrmLvl.workspace_permission_id <= 400
	AND 
		WorkSpaceMmLkp.active = true;
	-- if permission check fails (i.e: return no rows) retun invalid permission level and dummt values
	 IF found_rows() = 0 THEN
		-- SET out_status = "invalid_user_permissions_level";
		-- Declare & return Dummy output variables to avoid "result.Scan()" method errors in GoLang
        SET out_product_id = 0;
        SET out_upc = "NULL";
		SET out_sku = "NULL";
        SET out_description = "NULL";
		SET out_vendor_id = 0;
        SET out_product_unit_sold_id = 0;
		SELECT out_product_id,out_vendor_id, out_product_unit_sold_id, out_upc, out_sku, out_description;
		LEAVE proc_label;
	END IF;

    SELECT product_id, vendor_id, product_unit_sold_id, upc, sku, description FROM commhub_junction.product WHERE upc = in_upc AND workspace_id = target_workspace_id;

END$$
DELIMITER ;