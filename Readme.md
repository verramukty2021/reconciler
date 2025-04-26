## Step to Test
1. Generate csv dummy input
2. Run reconciler program

All command should be run from project root dir /reconciler

## Generate csv dummy input
#### Generate csv file using default config
```
go run cmd/csv_generator/main.go
```
#### Generate csv file using input config
```
go run cmd/csv_generator/main.go <totalTrx> <totalBankA> <totalBankB> <totalBankB> <startDate> <endDate>
```
The command will generate csv files on these folders: csv/input_bank_trx, csv/input_system_trx
Increase <totalTrx> i.e. 1000 to test how the reconciler performs on large number of transactions.
Configure <totalBankA> <totalBankB> <totalBankB> to split bank statement data into multiple banks.
Configure <startDate> <endDate> to set transaction time range

#### Case 1: All match case
```
go run cmd/csv_generator/main.go 10 3 3 4 2025-01-01 2025-01-31
```
**Reconciliation Result**
All transactionn will be matched

#### Case 2: System transaction is missing on bank statement
```
go run cmd/csv_generator/main.go 10 3 3 2 2025-01-01 2025-01-31
```

**Reconciliation Result**
2 system transactions will be listed on "System transaction details missing in bank statement(s)"

#### Case 3: Bank Data (bank_a & bank_b) is missing on system transaction
```
go run cmd/csv_generator/main.go 10 3 3 4 2025-01-01 2025-01-31
```
After the above command, delete manually from bank_a & bank_b csv file to create the discrepancy.
**Reconciliation Result**
The deleted bank data will be listed on "Bank statement details missing in system transactions (grouped by bank)".

## Run reconciler program
#### Reconcile using default connfig
```
go run cmd/reconciler/main.go
```
Default config: 2025-01-01 2025-01-31

#### Reconcile using input config
```
go run cmd/reconciler/main.go <startDate> <endDate>
```

#### Example
```
go run cmd/reconciler/main.go 2025-01-01 2025-01-31
```
Configure <startDate> <endDate> to filter transaction data that will be processed during reconciliation.

## Algorithm Overview
1. Read all system transaction & bank statement data, and put it into map.
    * Both map will have the same key format that allows them to be matched.
    * Key format: CREDIT_AbsoluteAmount, Example: CREDIT_100.00, DEBIT_50.00
    * Map Illustration:
        * bankTrxMap[CREDIT_100.00][{bankTrxDetails_1}, {bankTrxDetails_2}, ...]
        * systemTrxMap[CREDIT_100.00][{systemTrxDetails_1}, {systemTrxDetails_2}, ...]
    * Space complexity is O(n) for both map
    * Time complexity is O(n) on this process
2. Iterate system transaction map, search the bank statement with the same key.
    * During iteration, delete matched transaction from both map.
    * At the end of the iteration, we got the unmatched transaction to be displayed on result.
    * Time complexity is O(n) for this iteration
3. Calculate all the output metrics, and print result.

##### Assumptions
* Each system transaction has (1:1) relationship with bank statement


##### Algorithm Complexity
* Time Complexity: O(n)
* Space Complexity: O(n)
