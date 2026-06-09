## 1. Service Layer

- [x] 1.1 Append `DescribeClsNoticeContentsByFilter()` — wraps `DescribeNoticeContents` with Offset/Limit pagination retry loop (Limit=100)

## 2. Data Source Implementation

- [x] 2.1 Create `data_source_tc_cls_notice_contents.go` with schema and read handler
- [x] 2.2 Optional `filters` param; computed `notice_content_list` with nested `notice_contents`
- [x] 2.3 Read handler: build filters, call service, flatten results

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cls_notice_contents` in `provider.go` DataSourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `data_source_tc_cls_notice_contents.md`
- [x] 4.2 Create `data_source_tc_cls_notice_contents_test.go`

