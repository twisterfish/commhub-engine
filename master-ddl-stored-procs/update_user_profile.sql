DELIMITER $$

DROP PROCEDURE commhub_junction.update_user_profile;
CREATE PROCEDURE commhub_junction.update_user_profile(
	IN in_emid varchar(255),
  IN in_pwid varchar(128),
	IN in_npwid varchar(128),
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
	DECLARE out_status VARCHAR(255);
	DECLARE test_user_id INT;
	
-- set up the exception handling
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
        GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);
		
		SET out_status = @full_error; -- this is set for debug while in dev
		
		SELECT out_status;
		
		ROLLBACK;
	END;
	
-- get the user's internal ID 
	SELECT end_user_id INTO test_user_id FROM end_user WHERE emid = in_emid AND pwid = MD5(in_pwid);
	
	IF found_rows() = 0 THEN	 
		SET out_status = "invalid_user"; -- this is set for debug while in dev
		SELECT out_status;
		
		LEAVE proc_label;
	END IF;
	
	IF in_npwid IS NULL OR CHAR_LENGTH(in_pwid) < 2 THEN 
-- create the end user's profile record	
		UPDATE commhub_junction.end_user_profile
		SET 
			ein_tax_id = in_ein_tax_id, 
			ssn_tax_id = in_ssn_tax_id, 
			last_name = in_last_name, 
			middle_name = in_middle_name, 
			first_name = in_first_name, 
			address1 = in_address1, 
			address2 = in_address2, 
			city = in_city, 
			province_state = in_province_state, 
			zip_postal_code = in_zip_postal_code, 
			email = in_email, 
			phone = in_phone, 
			country_code = in_country_code 
		WHERE 
			end_user_id = test_user_id;

			IF row_count() = 0 THEN	 
			SET out_status = "no rows found no pass"; -- this is set for debug while in dev
			SELECT out_status;
			LEAVE proc_label;
		END IF;

	ELSE
		START TRANSACTION;
		UPDATE commhub_junction.end_user SET pwid = MD5(in_npwid) WHERE end_user_id = test_user_id; 
		UPDATE commhub_junction.end_user_profile
		SET 
			ein_tax_id = in_ein_tax_id, 
			ssn_tax_id = in_ssn_tax_id, 
			last_name = in_last_name, 
			middle_name = in_middle_name, 
			first_name = in_first_name, 
			address1 = in_address1, 
			address2 = in_address2, 
			city = in_city, 
			province_state = in_province_state, 
			zip_postal_code = in_zip_postal_code, 
			email = in_email, 
			phone = in_phone, 
			country_code = in_country_code 
		WHERE 
			end_user_id = test_user_id;
		COMMIT;

		IF row_count() = 0 THEN	 
			SET out_status = "no rows found with pass"; -- this is set for debug while in dev
			SELECT out_status;
			LEAVE proc_label;
		END IF;

	END IF;

	SET out_status = "success";
	
	SELECT out_status;
  
END$$

DELIMITER ;