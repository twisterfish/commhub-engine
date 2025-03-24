DELIMITER $$ 

 DROP PROCEDURE commhub_junction.update_ticket;

CREATE PROCEDURE commhub_junction.update_ticket(
  IN in_end_user_token VARCHAR(32),
	IN in_workspace_token VARCHAR(32),
	IN in_work_ticket_id INT,
	IN in_work_ticket_type_id INT,
  IN in_work_status_id INT, 
  IN in_ticket_group_id INT,
  IN in_running_time INT, 
  IN in_assigned_to_user_token VARCHAR(32),
  IN in_real_property_id INT,
  IN in_title VARCHAR(255),
  IN in_description TEXT
) 
  proc_label : BEGIN 
  
DECLARE requesters_id, requesters_workspace_permission_id, target_workspace_id, ticket_group_id, assigned_to_user_id INT;

DECLARE out_status VARCHAR(255);

-- set up the exception handling
DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE,
@errno = MYSQL_ERRNO,
@text = MESSAGE_TEXT;

  SET
    @full_error = CONCAT("ERROR ", @errno, " (", @sqlstate, "): ", @text);

  SET
    out_status = @full_error;

  SELECT
    out_status;

  ROLLBACK;

END;

-- find out if the user has the authority to do this		
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
  INNER JOIN commhub_junction.end_user AS EnUsr ON WorkSpaceMmLkp.member_id = EnUsr.end_user_id
  INNER JOIN commhub_junction.workspace_permission AS PrmLvl ON WorkSpaceMmLkp.workspace_permission_id = PrmLvl.workspace_permission_id
  INNER JOIN commhub_junction.workspace AS WorkSpaceMain ON WorkSpaceMmLkp.workspace_id = WorkSpaceMain.workspace_id
WHERE
  WorkSpaceMain.workspace_token = in_workspace_token
  AND EnUsr.end_user_token = in_end_user_token
  AND WorkSpaceMmLkp.active = true;

-- if nothing found then user isn't even a member of the group	   
IF found_rows() = 0 THEN
SET
  out_status = "invalid_requester_is_not_a_member";

SELECT
  out_status;

LEAVE proc_label;

END IF;

-- requester must be at worker level at least	
IF requesters_workspace_permission_id > 400 THEN
SET
  out_status = "invalid_requester_is_not_authorized";

SELECT
  out_status;

LEAVE proc_label;

END IF;

-- check to see if there is a valid user for this token
SELECT end_user_id INTO assigned_to_user_id FROM commhub_junction.end_user WHERE end_user_token = in_assigned_to_user_token;

-- if nothing found then user isn't even a member of the group	   
IF found_rows() = 0 THEN
SET
  out_status = "invalid_user";

SELECT
  out_status;

LEAVE proc_label;

END IF;

-- update the ticket
UPDATE
  commhub_junction.work_ticket
SET
  work_ticket_status_id = in_work_status_id, 
  work_ticket_type_id = in_work_ticket_type_id,
  ticket_group_id = in_ticket_group_id,
  real_property_id = in_real_property_id,
  assigned_to_user_id = assigned_to_id,
  running_time = in_running_time,
  title = in_title,
  description = in_description
WHERE
  work_ticket_id = in_work_ticket_id;

SET
  out_status = "success";

SELECT
  out_status;

END$$

DELIMITER ;