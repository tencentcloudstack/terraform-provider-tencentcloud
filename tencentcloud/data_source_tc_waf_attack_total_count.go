/*
Use this data source to query detailed information of waf attack_total_count

Example Usage

Obtain the specified domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "domain.com"
  query_string = "method:GET"
}
```

Obtain all domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "all"
  query_string = "method:GET"
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafAttackTotalCount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackTotalCountRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query domain name, all domain use all.",
			},
			"query_string": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "",
				Description: "Query conditions.",
			},
			"total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of attacks.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafAttackTotalCountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_attack_total_count.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		service          = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		attackTotalCount *waf.GetAttackTotalCountResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackTotalCountByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		attackTotalCount = result
		return nil
	})

	if err != nil {
		return err
	}

	if attackTotalCount.TotalCount != nil {
		_ = d.Set("total_count", attackTotalCount.TotalCount)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
