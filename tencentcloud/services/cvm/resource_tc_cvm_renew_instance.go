package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmRenewInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRenewInstanceCreate,
		Read:   resourceTencentCloudCvmRenewInstanceRead,
		Delete: resourceTencentCloudCvmRenewInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"instance_charge_prepaid": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify the renewal duration of the Subscription instance, whether to set automatic renewal, and other attributes. For yearly and monthly subscription instances, this parameter is required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Subscription period; unit: month; valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60. Note: This field may return null, indicating that no valid value is found.",
						},
						"renew_flag": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Auto renewal flag. Valid values:\n" +
								"- `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically;\n" +
								"- `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically;\n" +
								"- `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically;\n" +
								"Default value: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. Note: This field may return null, indicating that no valid value is found.",
						},
					},
				},
			},

			"renew_portable_data_disk": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
				Description: "Whether to renew the elastic data disk. Valid values:\n" +
					"- `TRUE`: Indicates to renew the subscription instance and renew the attached elastic data disk at the same time\n" +
					"- `FALSE`: Indicates that the subscription instance will be renewed and the elastic data disk attached to it will not be renewed\n" +
					"Default value: TRUE.",
			},
		},
	}
}

func resourceTencentCloudCvmRenewInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cvm.NewRenewInstancesRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceIds = []*string{&instanceId}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := cvm.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok && v.(string) != "" {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	if v, _ := d.GetOk("renew_portable_data_disk"); v != nil {
		request.RenewPortableDataDisk = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RenewInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm renewInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmRenewInstanceRead(d, meta)
}

func resourceTencentCloudCvmRenewInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRenewInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
