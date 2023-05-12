/*
Use this data source to query detailed information of lighthouse instance_vnc_url

Example Usage

```hcl
data "tencentcloud_lighthouse_instance_vnc_url" "instance_vnc_url" {
  instance_id = "lhins-123456"
}
```
*/
package tencentcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseInstanceVncUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceVncUrlRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
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

func dataSourceTencentCloudLighthouseInstanceVncUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_instance_vnc_url.read")()
	defer inconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)

	request := lighthouse.NewDescribeInstanceVncUrlRequest()
	response := lighthouse.NewDescribeInstanceVncUrlResponse()
	request.InstanceId = helper.String(instanceId)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {

		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().DescribeInstanceVncUrl(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil || response.Response.InstanceVncUrl == nil {
		d.SetId("")
		return fmt.Errorf("Response is nil")

	}

	d.SetId(instanceId)
	_ = d.Set("instance_vnc_url", *response.Response.InstanceVncUrl)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"instance_vnc_url": *response.Response.InstanceVncUrl,
		}); e != nil {
			return e
		}
	}
	return nil
}
