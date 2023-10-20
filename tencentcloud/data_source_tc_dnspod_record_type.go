/*
Use this data source to query detailed information of dnspod record_type

Example Usage

```hcl
data "tencentcloud_dnspod_record_type" "record_type" {
  domain_grade = "DP_FREE"
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodRecordType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordTypeRead,
		Schema: map[string]*schema.Schema{
			"domain_grade": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain level. + Old packages: D_FREE, D_PLUS, D_EXTRA, D_EXPERT, D_ULTRA correspond to free package, personal luxury, enterprise 1, enterprise 2, enterprise 3. + New packages: DP_FREE, DP_PLUS, DP_EXTRA, DP_EXPERT, DP_ULTRA correspond to new free, personal professional, enterprise basic, enterprise standard, enterprise flagship.",
			},

			"type_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Record type list.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDnspodRecordTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnspod_record_type.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain_grade"); ok {
		paramMap["DomainGrade"] = helper.String(v.(string))
	}

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	var typeList []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodRecordTypeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		typeList = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(typeList))
	if typeList != nil {
		_ = d.Set("type_list", typeList)
	}

	// d.SetId(helper.DataResourceIdsHash(ids))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), typeList); e != nil {
			return e
		}
	}
	return nil
}
