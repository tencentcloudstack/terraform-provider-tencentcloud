/*
Use this data source to query detailed information of tsf application_attribute

Example Usage

```hcl
data "tencentcloud_tsf_application_attribute" "application_attribute" {
  application_id = "application-a24x29xv"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfApplicationAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationAttributeRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application Id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Application list other attribute.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of instances.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"run_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of running instances.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"group_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of deployment groups under the application.Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTsfApplicationAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_application_attribute.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.ApplicationAttribute

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationAttributeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		applicationAttributeMap := map[string]interface{}{}

		if result.InstanceCount != nil {
			applicationAttributeMap["instance_count"] = result.InstanceCount
		}

		if result.RunInstanceCount != nil {
			applicationAttributeMap["run_instance_count"] = result.RunInstanceCount
		}

		if result.GroupCount != nil {
			applicationAttributeMap["group_count"] = result.GroupCount
		}

		ids = append(ids, *result.ApplicationId)
		_ = d.Set("result", applicationAttributeMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), applicationAttributeMap); e != nil {
			return e
		}
	}
	return nil
}
