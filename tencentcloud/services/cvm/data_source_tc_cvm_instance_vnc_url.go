package cvm

import (
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmInstanceVncUrl() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_instance_vnc_url.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var response *cvm.DescribeInstanceVncUrlResponse
	request := cvm.NewDescribeInstanceVncUrlRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceId = helper.String(instanceId)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeInstanceVncUrl(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
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
