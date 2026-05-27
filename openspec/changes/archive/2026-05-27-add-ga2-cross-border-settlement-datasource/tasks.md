## 1. Data Source Implementation

- [x] 1.1 Create data source file `tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement.go` with schema definition and Read function
- [x] 1.2 Register the data source in `tencentcloud/provider.go` data source map

## 2. Provider Documentation

- [x] 2.1 Add data source entry to `tencentcloud/provider.md` under GA2 service section

## 3. Example Documentation

- [x] 3.1 Create data source example file `tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement.md`

## 4. Unit Tests

- [x] 4.1 Create unit test file `tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement_test.go` using gomonkey mock for API calls
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify tests pass
