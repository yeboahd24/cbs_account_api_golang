CREATE TABLE AccountType (
    AccountID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    Name VARCHAR(255),
    Description TEXT,
    StartRange INT,
    EndRange INT
);

 CREATE TABLE ChartOfAccount (
    AccountID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    AccountTypeID UUID REFERENCES AccountType(AccountID),
    AccountNumber INT,
    Name VARCHAR(255)
);


CREATE TABLE AccountBalance (
    BalanceID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    AccountID UUID,
    Balance INT,
    FOREIGN KEY (AccountID) REFERENCES Account(AccountID)
);

CREATE TABLE JournalEntry (
    TransactionID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    AccountToDebitID UUID,
    AccountToCreditID UUID,
    Date TIMESTAMP,                                      
    Amount INT,
    Description TEXT,
    FOREIGN KEY (AccountToDebitID) REFERENCES Account(AccountID),
    FOREIGN KEY (AccountToCreditID) REFERENCES Account(AccountID)
);

CREATE TABLE Account (
    AccountID UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    Name VARCHAR(255),
    AccountNumber INT UNIQUE,
    COAID UUID,
    FOREIGN KEY (COAID) REFERENCES ChartOfAccount(AccountID)
);

