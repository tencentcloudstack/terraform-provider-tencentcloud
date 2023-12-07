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
				Description: "application Id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "application list other attribute.",
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
	ids := ""

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_id"); ok {
		ids = v.(string)
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var attribute *tsf.ApplicationAttribute

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationAttributeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		attribute = result
		return nil
	})
	if err != nil {
		return err
	}

	applicationAttributeMap := map[string]interface{}{}
	if attribute != nil {
		if attribute.InstanceCount != nil {
			applicationAttributeMap["instance_count"] = attribute.InstanceCount
		}

		if attribute.RunInstanceCount != nil {
			applicationAttributeMap["run_instance_count"] = attribute.RunInstanceCount
		}

		if attribute.GroupCount != nil {
			applicationAttributeMap["group_count"] = attribute.GroupCount
		}

		_ = d.Set("result", []interface{}{applicationAttributeMap})
	}

	d.SetId(ids)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), applicationAttributeMap); e != nil {
			return e
		}
	}
	return nil
}
