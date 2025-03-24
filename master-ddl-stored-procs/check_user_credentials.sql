DELIMITER $$

DROP PROCEDURE commhub_junction.check_user_credentials;
CREATE PROCEDURE commhub_junction.check_user_credentials(
	IN in_user_emid VARCHAR(255),
	IN in_user_pwid VARCHAR(128)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(255);
	DECLARE out_user_id INT;
	DECLARE out_user_token, out_api_token VARCHAR(32);
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		SET out_status = @full_error; -- this is set for debug while in dev
		-- SET out_status = "invalid";
		SET out_user_id = 0;
		SET out_user_token = "null";
		SET out_api_token = "null";
		
		SELECT out_status, out_user_id, out_user_token, out_api_token;
	END;
	
-- get the group creator's internal ID
	SELECT end_user_id, end_user_token INTO out_user_id, out_user_token FROM end_user WHERE emid = in_user_emid AND pwid = MD5(in_user_pwid);
	IF found_rows() = 0 THEN
	 
		SET out_status = "invalid_user"; -- this is set for debug while in dev
		SET out_user_id = 0;
		SET out_user_token = "null";
		SET out_api_token = "null";
		SELECT out_status, out_user_id, out_user_token, out_api_token;
		
		LEAVE proc_label;
	END IF;

-- get the API token if they pass validation
	SELECT token INTO out_api_token FROM commhub_junction.session_token WHERE token_id = 1;    			
	SET out_status = "success";	
	SELECT out_status, out_user_id, out_user_token, out_api_token;
  
END$$

DELIMITER ;