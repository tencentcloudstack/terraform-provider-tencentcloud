/*
Use this data source to query detailed information of lighthouse instance_traffic_package

Example Usage

```hcl
data "tencentcloud_lighthouse_instance_traffic_package" "instance_traffic_package" {
  instance_ids =
  offset = 0
  limit = 20
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseInstanceTrafficPackage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceTrafficPackageRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseInstanceTrafficPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_instance_traffic_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceReturnableSet []*lighthouse.InstanceReturnable

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseInstanceTrafficPackageByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceReturnableSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceReturnableSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceReturnableSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
