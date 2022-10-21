-- pg_card_list_types

-- pg_client_create

-- pg_client_fetch

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
    "card_number": "1111222233334444",
    "card_holder": "Credit Holder",
    "cvv": "123,
    "expire_date": "2022-11-10"
}');

-- pg_card_fetch

-- pg_client_register_card

-- pg_transactions_create

-- pg_client_list_transactions

-- pg_client_fetch

-- pg_transactions_create
    
/* 
    Creates client and wallet.
    Creates card and registers card for client.
    Creates tx in paid status (D+0) with 3% fee.
    Returns created tx. 
*/

CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 1234,
    "value": 100,
    "description": "Smartband XYZ 1.0",
    "id_card_type": 1,
    "card_number": "1111222233334444",
    "card_holder": "The Peps",
    "expire_date": "2022-11-10",
    "cvv": 123
  }');

/* 
    Creates tx in paid status (D+0) with 3% fee.
    Returns created tx.
*/

CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 1234,
    "value": 200,
    "description": "Smartband XYZ 2.0",
    "id_card_type": 1,
    "card_number": "1111222233334444",
    "card_holder": "The Peps",
    "expire_date": "2022-11-10",
    "cvv": 123
  }');


/* 
    Creates card and registers card for client.
    Creates tx in waiting status (D+30) with 5% fee.
    Returns created tx.
*/
CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 1234,
    "value": 300,
    "description": "Smartband XYZ 3.0",
    "id_card_type": 2,
    "card_number": "9999888877776666",
    "card_holder": "The Peps",
    "expire_date": "2022-11-10",
    "cvv": 123
  }');

/*  
    Creates client and wallet.
    Registers card for client.
    Creates tx in waiting status (D+30) with 5% fee.
    Returns created tx.
*/
CALL `pagarme`.`pg_transactions_create`('{
    "id_client": 9876,
    "value": 400,
    "description": "Smartband XYZ 4.0",
    "id_card_type": 2,
    "card_number": "9999888877776666",
    "card_holder": "The Peps",
    "expire_date": "2022-11-10",
    "cvv": 123
  }');
