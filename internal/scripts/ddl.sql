
DROP TABLE wallets
;
DROP TABLE transactions
;
DROP TABLE clients_cards
;
DROP TABLE clients
;
DROP TABLE cards
;
DROP TABLE card_types
;
DROP TABLE log_sp
;

-- 
-- TABLE: card_types 
--

CREATE TABLE card_types(
    id_card_type    TINYINT        AUTO_INCREMENT,
    card_type       VARCHAR(30)    NOT NULL,
    PRIMARY KEY (id_card_type)
)ENGINE=INNODB
;

INSERT INTO `card_types` (`id_card_type`, `card_type`) VALUES ('1', 'debit');
INSERT INTO `card_types` (`id_card_type`, `card_type`) VALUES ('2', 'credit');



-- 
-- TABLE: cards 
--

CREATE TABLE cards(
    id_card         BIGINT         AUTO_INCREMENT,
    id_card_type    TINYINT        NOT NULL,
    number          CHAR(16)       NOT NULL,
    cardholder      VARCHAR(60)    NOT NULL,
    cvv             CHAR(3)        NOT NULL,
    expire_date     DATE           NOT NULL,
    PRIMARY KEY (id_card), 
    UNIQUE INDEX ui_number(number),
    INDEX ix_cardholder(cardholder),
    INDEX ix_expire_date(expire_date),
    INDEX Ref51(id_card_type), 
    CONSTRAINT Refcard_types1 FOREIGN KEY (id_card_type)
    REFERENCES card_types(id_card_type)
)ENGINE=INNODB
;



-- 
-- TABLE: clients 
--

CREATE TABLE clients(
    id_client    INT    NOT NULL,
    PRIMARY KEY (id_client)
)ENGINE=INNODB
;



-- 
-- TABLE: clients_cards 
--

CREATE TABLE clients_cards(
    id_card      BIGINT      NOT NULL,
    id_client    INT         NOT NULL,
    date         DATETIME    NOT NULL,
    PRIMARY KEY (id_card, id_client), 
    INDEX Ref36(id_client),
    INDEX Ref27(id_card), 
    CONSTRAINT Refclients6 FOREIGN KEY (id_client)
    REFERENCES clients(id_client),
    CONSTRAINT Refcards7 FOREIGN KEY (id_card)
    REFERENCES cards(id_card)
)ENGINE=INNODB
;



-- 
-- TABLE: transactions 
--

CREATE TABLE transactions(
    id_tx          BIGINT            AUTO_INCREMENT,
    id_card        BIGINT            NOT NULL,
    id_client      INT               NOT NULL,
    description    VARCHAR(60),
    date           DATETIME          NOT NULL,
    value          DECIMAL(10, 2)    NOT NULL,
    fee            DECIMAL(10, 2)    NOT NULL,
    status         CHAR(1)           NOT NULL,
    PRIMARY KEY (id_tx, id_card, id_client), 
    UNIQUE INDEX ui_tx(id_tx),
    INDEX ix_date(date),
    INDEX Ref22(id_card),
    INDEX Ref33(id_client), 
    CONSTRAINT Refcards2 FOREIGN KEY (id_card)
    REFERENCES cards(id_card),
    CONSTRAINT Refclients3 FOREIGN KEY (id_client)
    REFERENCES clients(id_client)
)ENGINE=INNODB
;



-- 
-- TABLE: wallets 
--

CREATE TABLE wallets(
    id_wallet          INT               NOT NULL,
    available_funds    DECIMAL(10, 2)    NOT NULL,
    waiting_funds      DECIMAL(10, 2)    NOT NULL,
    PRIMARY KEY (id_wallet), 
    INDEX Ref310(id_wallet), 
    CONSTRAINT Refclients10 FOREIGN KEY (id_wallet)
    REFERENCES clients(id_client)
)ENGINE=INNODB
;

-- 
-- TABLE: log_sp 
--

CREATE TABLE `log_sp` (
  `id_log` bigint NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `sp` varchar(100) NOT NULL,
  `msg` varchar(2000) NOT NULL,
  `params` json DEFAULT NULL,
  PRIMARY KEY (`id_log`)
) ENGINE=InnoDB

