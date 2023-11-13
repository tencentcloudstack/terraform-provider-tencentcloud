/*
Use this data source to query detailed information of cvm instance_vnc_url

Example Usage

```hcl
data "tencentcloud_cvm_instance_vnc_url" "instance_vnc_url" {
  instance_id = "ins-r9hr2upy"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmInstanceVncUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmInstanceVncUrlRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID. To obtain the instance IDs, you can call `DescribeInstances` and look for `InstanceId` in the response.",
			},

			"instance_vnc_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance VNC URL.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmInstanceVncUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_instance_vnc_url.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmInstanceVncUrlByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceVncUrl = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceVncUrl))
	if instanceVncUrl != nil {
		_ = d.Set("instance_vnc_url", instanceVncUrl)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
