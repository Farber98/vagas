DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_card_create`(pIn json)
BEGIN
	DECLARE pIdCard BIGINT;
    DECLARE pExpireDate DATE DEFAULT pIn->>'$.expire_date';
    
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", COALESCE(@errno, ''), " (", COALESCE(@sqlstate, ''), "): ", COALESCE(@text, ''));
		ROLLBACK;
		INSERT INTO log_sp VALUES(0, NOW(), 'pg_card_create', @full_error, pIn);
		RESIGNAL;
	END;
    
    START TRANSACTION;
		INSERT INTO cards VALUES(0, pIn->>'$.card_type', pIn->>'$.number', pIn->>'$.card_holder', pIn->>'$.cvv', pIn->>'$.expire_date');
		SET pIdCard = LAST_INSERT_ID();
	COMMIT;
    
    SELECT JSON_OBJECT(
            'id_card', id_card,
            'card_type', card_type,
            'card_number', number,
            'card_holder', cardholder,
            'expire_date', expire_date
          ) pOut
	FROM cards INNER JOIN card_types USING(id_card_type)
    WHERE id_card = pIdCard;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_card_fetch`(pIn json)
BEGIN
    SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
    SELECT  JSON_OBJECT(
            'id_card', id_card,
            'card_type', card_type,
            'card_number', number,
            'card_holder', cardholder,
            'expire_date', expire_date
          ) pOut
    FROM  cards  INNER JOIN card_types USING(id_card_type);
    SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_client_create`(pIn json)
BEGIN
	DECLARE pError CONDITION FOR SQLSTATE '45000';
    DECLARE pIdClient INT DEFAULT pIn->>'$.id_client';
    
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", COALESCE(@errno, ''), " (", COALESCE(@sqlstate, ''), "): ", COALESCE(@text, ''));
		ROLLBACK;
		INSERT INTO log_sp VALUES(0, NOW(), 'pg_client_create', @full_error, pIn);
		RESIGNAL;
	END;
    
    START TRANSACTION;
        INSERT INTO clients VALUES(pIdClient);
        INSERT INTO wallets VALUES(pIdClient,0,0);
    COMMIT;
	
    SELECT  JSON_OBJECT(
            'id_client', c.id_client,
            'id_wallet', w.id_wallet,
            'available_funds', w.available_funds,
            'waiting_funds', w.waiting_funds
        ) pOut
    FROM    		clients c  
    INNER	JOIN 	wallets w USING(id_client)
    WHERE 		id_client = pIdClient;
    
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_client_fetch`(pIn json)
BEGIN
    SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
    SELECT  JSON_OBJECT(
            'id_client', c.id_client,
            'id_wallet', w.id_wallet,
            'available_funds', w.available_funds,
            'waiting_funds', w.waiting_funds
          ) pOut
    FROM    		clients c  
    INNER	JOIN 	wallets w USING(id_client);
    SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_client_list_transactions`(pIn json)
BEGIN
    SET SESSION GROUP_CONCAT_MAX_LEN = 1024 * 1024 * 1024;
    SET tmp_table_size = 1024 * 1024 * 1024 * 64;
    SET max_heap_table_size = 1024 * 1024 * 1024 * 64;
    SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
    SELECT  COALESCE(JSON_ARRAYAGG(JSON_OBJECT(
            'id_tx', tx.id_tx,
            'card_number', CONCAT(REPEAT("*",12),RIGHT(ca.number,4)),
            'card_type', ct.card_type,
            'card_holder', ca.cardholder,
            'id_client', tx.id_client,
            'description', tx.description,
            'date', tx.date,
            'value', tx.value,
            'fee', tx.fee,
            'status', tx.status
          )), JSON_ARRAY() ) pOut
    FROM    		transactions tx 
    INNER	JOIN 	cards ca USING(id_card)
    INNER	JOIN    card_types ct USING(id_card_type)
    WHERE			tx.id_client = pIn ->> '$.id_client';

    SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_client_register_card`(pIn json)
BEGIN
    INSERT INTO clients_cards VALUES(pIn->>'$.id_card', pIn->>'$.id_client', NOW());
    SELECT  JSON_OBJECT(
            'id_card', id_card,
            'id_client', id_client
          ) pOut
    FROM    clients_cards
    WHERE 	id_card = pIn->>'$.id_card' AND id_client = pIn->>'$.id_client';
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_transactions_create`(pIn json)
BEGIN

	DECLARE pIdCardType TINYINT DEFAULT pIn->>'$.id_card_type';
	DECLARE pValue DECIMAL(10,2) DEFAULT pIn->>'$.value';
    DECLARE pFee DECIMAL(10,2);
    DECLARE pIdClient INT DEFAULT pIn->>'$.id_client';
	DECLARE pIdTx BIGINT;
    
	DECLARE pError CONDITION FOR SQLSTATE '45000';
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @text = MESSAGE_TEXT;
		SET @full_error = CONCAT("ERROR ", COALESCE(@errno, ''), " (", COALESCE(@sqlstate, ''), "): ", COALESCE(@text, ''));
		ROLLBACK;
		INSERT INTO log_sp VALUES(0, NOW(), 'pg_transactions_create', @full_error, pIn);
		RESIGNAL;
	END;

    START TRANSACTION;

	    SET pFee = IF(pIdCardType = 1, pValue * 0.05, pValue * 0.03);

	    INSERT INTO transactions VALUES(0,pIn->>'$.id_card',pIdClient, pIn->>'$.description', 
	    								IF(pIdCardType = 1, NOW(), ADDDATE(NOW(), INTERVAL 30 DAY)),
                                        pValue, pFee, IF(pIdCardType = 1, 'P', 'W'));

        SET pIdTx = LAST_INSERT_ID();

	    SELECT * FROM wallets WHERE id_wallet = pIdClient FOR UPDATE;

        UPDATE wallets 
        SET available_funds = IF(pIdCardType = 1, available_funds + pValue - pFee, available_funds),
	    	waiting_funds = IF(pIdCardType = 2, waiting_funds + pValue - pFee, waiting_funds)
	    WHERE id_wallet = pIdClient;
    COMMIT; 

	SELECT  JSON_OBJECT(
            'id_tx', tx.id_tx,
            'card_number', CONCAT(REPEAT("*",12),RIGHT(ca.number,4)),
            'card_type', ct.card_type,
            'card_holder', ca.cardholder,
            'id_client', tx.id_client,
            'description', tx.description,
            'date', tx.date,
            'value', tx.value,
            'fee', tx.fee,
            'status', tx.status
          ) pOut
    FROM    		transactions tx 
    INNER	JOIN 	cards ca USING(id_card)
    INNER	JOIN    card_types ct USING(id_card_type)
    WHERE			id_tx = pIdTx;

END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`juan`@`localhost` PROCEDURE `pg_card_list_types`()
BEGIN
	SELECT  COALESCE(JSON_ARRAYAGG(
            JSON_OBJECT(
              'id_card_type', id_card_type,
              'card_type', card_type
            )), JSON_ARRAY()
          ) pOut
	FROM    card_types;
END$$
DELIMITER ;
