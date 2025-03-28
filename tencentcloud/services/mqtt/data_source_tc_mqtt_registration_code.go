package mqtt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMqttRegistrationCode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMqttRegistrationCodeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			"registration_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Registration code.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMqttRegistrationCodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mqtt_registration_code.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	var respData *mqttv20240516.ApplyRegistrationCodeResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMqttRegistrationCodeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("ApplyRegistrationCode response is nil."))
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.RegistrationCode != nil {
		_ = d.Set("registration_code", respData.RegistrationCode)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
