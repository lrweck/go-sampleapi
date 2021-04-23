CREATE TABLE Accounts (
    AccountID uuid PRIMARY KEY, 
	DocNumber TEXT 
);

CREATE TABLE OperationTypes (
    OperationTypeID SMALLINT PRIMARY KEY, 
	Description TEXT 
);

CREATE TABLE Transactions (
    TransactionID uuid PRIMARY KEY,
    AccountID uuid NOT NULL REFERENCES accounts(AccountID),
    OperationTypeID SMALLINT NOT NULL REFERENCES OperationTypes(OperationTypeID),
    Amount NUMERIC NOT NULL DEFAULT 0,
    EventDate TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON Transactions (AccountID,EventDate);
