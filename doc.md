# Account Type

## Sample Request:

```bash
curl -X POST \
  http://localhost:8000/accounttype \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer e4f2A5gWIJrmsZ1wWxch3VBnef7pnlvWZfiLadc3i-c' \
  -d '{
    "name": "Assets",
    "description": "Main Assets",
    "start_range": 1000,
    "end_range": 1999
  }'
```

## Response:

```json
{ "status_code": 200, "message": "Account type created successfully" }
```

# CharOfAccount

## Sample Request:

```bash
curl -X POST \
  http://localhost:8000/coa \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer e4f2A5gWIJrmsZ1wWxch3VBnef7pnlvWZfiLadc3i-c' \
  -d '{
    "account_type_id": "df94e5a3-9a2a-496a-b177-23b5305f6e5b",
    "account_number": 2100,
    "name": "Loan to members"
  }'
```

## Response:

```json
{ "status_code": 200, "message": "Chart of account created successfully" }
```

# Account

## Sample Request:

```bash
curl -X POST \
  http://localhost:8000/account \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer b5H0IaFg5rafZ9qqI9UuidtnrnAbR875Jj7gtuHzUZ0' \
  -d '{
    "coa_id": "bfc0618b-4ad3-46b5-b0a8-b88321770494",
    "name": "Loan to members"
  }'
```

## Response:

```json
{ "status_code": 200, "message": "Account created successfully" }
```

# Accounts

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/account \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer b5H0IaFg5rafZ9qqI9UuidtnrnAbR875Jj7gtuHzUZ0' \
```

## Response:

```json
[
  {
    "id": "0e9807d6-d127-436a-b5ee-331081de2242",
    "name": "Loan to Groups",
    "account_number": 2101
  },
  {
    "id": "03c38358-0c37-4295-bb7b-1a03be7db025",
    "name": "Loan to Members",
    "account_number": 2102
  }
]
```

# Chart Of Accounts

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/coa \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9'
```

## Response:

```json
[
  {
    "id": "bfc0618b-4ad3-46b5-b0a8-b88321770494",
    "name": "Loan to members",
    "account_number": 2100
  }
]
```

# Account Types

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/accounttype \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer b5H0IaFg5rafZ9qqI9UuidtnrnAbR875Jj7gtuHzUZ0'
```

## Response

```json
[
  {
    "id": "df94e5a3-9a2a-496a-b177-23b5305b79",,
    "name": "Assets",
    "description": "Main Assets",
    "start_range": 1000,
    "end_range": 1999
  },
  {
    "id": "df94e5a3-9a2a-496a-b177-23b5305f6e5b",
    "name": "Liability",
    "description": "Main Liability",
    "start_range": 2000,
    "end_range": 2999
  }
]
```

# Journal Entry

## Sample Request:

```bash
curl -X POST \
  http://localhost:8000/journalentry \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer b5H0IaFg5rafZ9qqI9UuidtnrnAbR875Jj7gtuHzUZ0' \
  -d '{
    "credit_account": 2101,
    "debit_account": 2102,
    "amount": 1000,
    "description": "Payment"
  }'

```

## Response:

```json
{ "status_code": 200, "message": "Journal entry created successfully" }
```

# Journal Entries

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/journalentry \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer b5H0IaFg5rafZ9qqI9UuidtnrnAbR875Jj7gtuHzUZ0'
```

## Response:

```json
[
  {
    "transactionid": "0dcf5466-81d2-4d63-8676-b7ae5ade87e4",
    "accounttodebitid": 2102,
    "accounttocreditid": 2101,
    "date": "0001-01-01T00:00:00Z",
    "amount": 1000,
    "description": "Payment"
  },
  {
    "transactionid": "21a1c28e-ec46-4208-bff1-f4fd6a644cb1",
    "accounttodebitid": 2102,
    "accounttocreditid": 2101,
    "date": "0001-01-01T00:00:00Z",
    "amount": 1000,
    "description": "Payment"
  },
  {
    "transactionid": "c45780e2-5d5f-43d0-98b6-759c744423f5",
    "accounttodebitid": 2102,
    "accounttocreditid": 2101,
    "date": "2024-06-07T15:47:40.894838Z",
    "amount": 2000,
    "description": "Payment"
  }
]
```

# Profit and Loss

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/profitandloss \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer PtHy78sGh-pat6dKu6OYqXaBcrJtvSvCFaUH8Pj9E-Y'
```

## Response:

```json
{
  "status_code": 200,
  "message": "Profit and loss data retrieved successfully",
  "data": {
    "expenses": [
      {
        "account_name": "Expenses",
        "balance": 8000
      }
    ],
    "incomes": [],
    "total_expenses": 8000,
    "total_incomes": 0
  }
}
```

# Balance Sheet

## Sample Request:

```bash
curl -X GET \
  http://localhost:8000/balancesheet \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9'
```

## Response:

```json
{
  "status_code": 200,
  "message": "Balance sheet data retrieved successfully",
  "data": {
    "assets": [],
    "equity": [
      {
        "account_name": "Loan Payment Interest",
        "account_number": 3101,
        "balance": 3000
      }
    ],
    "liabilities": [
      {
        "account_name": "Expenses",
        "account_number": 2104,
        "balance": -2000
      },
      {
        "account_name": "Church Loans",
        "account_number": 2106,
        "balance": 0
      }
    ]
  }
}
```
