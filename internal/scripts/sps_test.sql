-- pg_card_list_types
CALL `pagarme`.`pg_card_list_types`();

-- pg_client_create
CALL `pagarme`.`pg_client_create`('{"id_client": 1}');

-- pg_client_fetch
CALL `pagarme`.`pg_client_fetch`('{"id_client": 1}');

-- pg_card_create

CALL `pagarme`.`pg_card_create`('{
    "id_card_type":1,
    "card_number": "1111222233334444",
    "card_holder": "Debit Holder",
    "cvv": "123",
    "expire_date": "2022-11-10"
}');

CALL `pagarme`.`pg_card_create`('{
    "id_card_type": 2,
    "card_number": "6666777788889999",
    "card_holder": "Credit Holder",
    "cvv": "123",
    "expire_date": "2022-11-10"
}');


-- pg_card_fetch_by_id
CALL `pagarme`.`pg_card_fetch_by_id`('{"id_card":1}');

-- pg_card_fetch_by_number
CALL `pagarme`.`pg_card_fetch_by_number`('{"card_number": "6666777788889999"}');

-- pg_client_register_card
CALL `pagarme`.`pg_client_register_card`('{"id_card":1, "id_client":1}');

-- pg_transactions_create


-- pg_client_list_transactions
    
/* 
    Creates client and wallet.
    Creates card and registers card for client.
    Creates tx in paid status (D+0) with 3% fee.
    Returns created tx. 
*/

CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 1,
    "id_card": 1,
    "value": 100,
    "description": "Smartband XYZ 1.0"
  }');

/* 
    Creates tx in paid status (D+0) with 3% fee.
    Returns created tx.
*/

CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 1,
    "id_card": 1,
    "value": 200,
    "description": "Smartband XYZ 2.0"
  }');


/* 
    Creates card and registers card for client.
    Creates tx in waiting status (D+30) with 5% fee.
    Returns created tx.
*/
CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 2,
    "id_card": 2,
    "value": 300,
    "description": "Smartband XYZ 3.0"
  }');

/*  
    Creates client and wallet.
    Registers card for client.
    Creates tx in waiting status (D+30) with 5% fee.
    Returns created tx.
*/
CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 3,
    "id_card": 2,
    "value": 400,
    "description": "Smartband XYZ 4.0",
  }');
