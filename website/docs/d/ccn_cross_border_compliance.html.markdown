---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_cross_border_compliance"
sidebar_current: "docs-tencentcloud-datasource-ccn_cross_border_compliance"
description: |-
  Use this data source to query detailed information of ccn cross_border_compliance
---

# tencentcloud_ccn_cross_border_compliance

Use this data source to query detailed information of ccn cross_border_compliance

## Example Usage

```hcl
data "tencentcloud_ccn_cross_border_compliance" "cross_border_compliance" {
  service_provider   = "UNICOM"
  compliance_id      = 10002
  email              = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date   = "2021-07-29"
  state              = "APPROVED"
}
```

## Argument Reference

The following arguments are supported:

* `business_address` - (Optional, String) (Fuzzy query) business license address.
* `company` - (Optional, String) (Fuzzy query) Company name.
* `compliance_id` - (Optional, Int) (Exact match) compliance approval form: 'ID'.
* `email` - (Optional, String) (Exact match) email.
* `issuing_authority` - (Optional, String) (Fuzzy query) Issuing authority.
* `legal_person` - (Optional, String) (Fuzzy query) legal representative.
* `manager_address` - (Optional, String) (Fuzzy query) ID card address of the person in charge.
* `manager_id` - (Optional, String) (Exact query) ID number of the person in charge.
* `manager_telephone` - (Optional, String) (Exact match) contact number of the person in charge.
* `manager` - (Optional, String) (Fuzzy query) Person in charge.
* `post_code` - (Optional, Int) (Exact match) post code.
* `result_output_file` - (Optional, String) Used to save results.
* `service_end_date` - (Optional, String) (Exact match) service end date, such as: '2020-07-28'.
* `service_provider` - (Optional, String) (Exact match) service provider, optional value: 'UNICOM'.
* `service_start_date` - (Optional, String) (Exact match) service start date, such as: '2020-07-28'.
* `state` - (Optional, String) (Exact match) status. Pending: PENDING, Passed: APPROVED, Denied: DENY.
* `uniform_social_credit_code` - (Optional, String) (Exact match) Uniform Social Credit Code.


