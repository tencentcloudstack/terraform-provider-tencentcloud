Use this data source to query bill detail (L3 detail bill) of billing.

Example Usage

```hcl
data "tencentcloud_billing_bill_detail" "example_by_month" {
  month            = "2024-01"
  need_record_num  = 1
  pay_mode         = "postPay"
}
```

```hcl
data "tencentcloud_billing_bill_detail" "example_by_time_range" {
  begin_time      = "2024-01-01 00:00:00"
  end_time        = "2024-01-31 23:59:59"
  need_record_num = 1
  resource_id     = "ins-xxxxxxxx"
}
```
