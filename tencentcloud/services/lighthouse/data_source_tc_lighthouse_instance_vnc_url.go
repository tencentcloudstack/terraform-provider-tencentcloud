package lighthouse

import (
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseInstanceVncUrl() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_instance_vnc_url.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)

	request := lighthouse.NewDescribeInstanceVncUrlRequest()
	response := lighthouse.NewDescribeInstanceVncUrlResponse()
	request.InstanceId = helper.String(instanceId)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().DescribeInstanceVncUrl(request)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"instance_vnc_url": *response.Response.InstanceVncUrl,
		}); e != nil {
			return e
		}
	}
	return nil
}
