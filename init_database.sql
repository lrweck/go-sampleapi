CREATE TABLE Accounts (
    AccountID uuid PRIMARY KEY, 
	DocNumber TEXT,
    AvailableCreditLimit NUMERIC DEFAULT 200 
);

CREATE TABLE OperationTypes (
    OperationTypeID SMALLINT PRIMARY KEY, 
	Description TEXT 
);
insert into OperationTypes
values 
(1,'COMPRA A VISTA'),
(2,'COMPRA PARCELADA'),
(3,'SAQUE'),
(4,'PAGAMENTO');

CREATE TABLE Transactions (
    TransactionID uuid PRIMARY KEY,
    AccountID uuid NOT NULL REFERENCES accounts(AccountID),
    OperationTypeID SMALLINT NOT NULL REFERENCES OperationTypes(OperationTypeID),
    Amount NUMERIC NOT NULL DEFAULT 0,
    EventDate TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON Transactions (AccountID,EventDate);
