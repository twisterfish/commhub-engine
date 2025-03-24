DELIMITER $$
DROP PROCEDURE commhub_junction.get_product_price;

CREATE PROCEDURE commhub_junction.get_product_price(
	IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
    IN in_product_id INT,
    IN in_effective_date INT
    )
proc_label:BEGIN
	DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id INT;
	DECLARE out_status VARCHAR(255);
    DECLARE out_price INT;
    
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		SET out_status = @full_error;
        SET out_price = 0;
		SELECT out_status,out_price;
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
		PrmLvl.workspace_permission_id <= 400
	AND 
		WorkSpaceMmLkp.active = true;
		 
	IF found_rows() = 0 THEN
		SET out_status = "invalid_user_permissions_level";
        SET out_price = 0;
		SELECT out_status,out_price;
		LEAVE proc_label;
	END IF;

    SELECT product_id, price FROM commhub_junction.product_price_history WHERE product_id = in_product_id AND effective_date <= FROM_UNIXTIME(in_effective_date) ORDER BY effective_date DESC LIMIT 1;

END$$
DELIMITER ;