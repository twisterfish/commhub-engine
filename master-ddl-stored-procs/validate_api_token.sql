DELIMITER $$

DROP PROCEDURE commhub_junction.validate_api_token;
CREATE PROCEDURE commhub_junction.validate_api_token(
	IN in_api_token VARCHAR(32)
    )
proc_label:BEGIN
	DECLARE out_status VARCHAR(128);
	DECLARE curr_api_token VARCHAR(32);
	DECLARE token_type, out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
        
		-- SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);		
		-- SET out_status = @full_error; -- this is set for debug while in dev
		SET out_status = "error";
		SELECT out_status,0,0,0,0,1;
	END;
	
-- get the current api token, validate against user's input and fail out on mismatch
	SELECT token_id, debug_mode, sample_mode, sample_rate, sample_duration, sample_time_start  INTO token_type, out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start  FROM commhub_junction.session_token WHERE token = in_api_token;
	
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid";
		SELECT out_status,0,0,0,0,1;	
		LEAVE proc_label;
	END IF;
	
	IF out_smpl_mode = 0 THEN	 
        IF out_dbg_mode = 0 THEN
			SET out_status = "1"; -- sample mode off & debug mode off
			SELECT out_status,out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start; 
		ELSEIF out_dbg_mode = 1 THEN	 
			SET out_status = "2"; -- sample mode off & debug mode on
            SELECT out_status,out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start; 
		END IF;
	ELSE	 
		IF out_dbg_mode = 0 THEN
			SET out_status = "3"; -- sample mode on & debug mode off
			SELECT out_status,out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start; 
		ELSEIF out_dbg_mode = 1 THEN	 
			SET out_status = "4"; -- sample mode on & debug mode on
            SELECT out_status,out_dbg_mode, out_smpl_mode, out_smpl_rate, out_smpl_duration, out_smpl_time_start; 
		END IF;
    END IF;    
  
END$$

DELIMITER ;