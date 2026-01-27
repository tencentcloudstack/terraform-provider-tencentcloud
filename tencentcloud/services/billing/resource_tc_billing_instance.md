Provides a resource to create a Billing instance

Example Usage

```hcl
resource "tencentcloud_billing_instance" "example" {
  product_code     = "p_cloudfirewall"
  sub_product_code = "sp_cloudfirewall_svv1"
  region_code      = "ap-guangzhou"
  zone_code        = "ap-guangzhou-6"
  pay_mode         = "PrePay"
  parameter = jsonencode({
    "goodsNum" : 1,
    "pid" : 1002147,
    "productCode" : "p_cloudfirewall",
    "subProductCode" : "sp_cloudfirewall_svv1",
    "sv_cloudfirewall_basic_aeps" : 1,
    "sv_cloudfirewall_basic_eeps" : 0,
    "sv_cloudfirewall_basic_ipsmonth" : 0,
    "sv_cloudfirewall_basic_mon" : 0,
    "sv_cloudfirewall_basic_ueps" : 0,
    "sv_cloudfirewall_extended_ates" : 0,
    "sv_cloudfirewall_extended_clasps" : 1,
    "sv_cloudfirewall_extended_clsesps" : 0,
    "sv_cloudfirewall_extended_ex" : 0,
    "sv_cloudfirewall_extended_ibtesps" : 0,
    "sv_cloudfirewall_extended_nats" : 0,
    "sv_cloudfirewall_extended_ndr" : 0,
    "sv_cloudfirewall_extended_pcs" : 0,
    "sv_cloudfirewall_extended_spt" : 0,
    "sv_cloudfirewall_extended_sra" : 0,
    "sv_cloudfirewall_extended_srb" : 0,
    "sv_cloudfirewall_extended_sub" : 0,
    "sv_cloudfirewall_extended_subs" : 0,
    "sv_cloudfirewall_extended_vpcbges" : 0,
    "timeSpan" : 1,
    "timeUnit" : "m"
  })
  project_id  = 0
  period      = 1
  period_unit = "m"
  renew_flag  = "NOTIFY_AND_MANUAL_RENEW"
}
```
