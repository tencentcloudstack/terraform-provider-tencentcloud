package lighthouse

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseRenewInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseRenewInstanceCreate,
		Read:   resourceTencentCloudLighthouseRenewInstanceRead,
		Delete: resourceTencentCloudLighthouseRenewInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"instance_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify attributes such as the purchase duration of the Subscription instance and whether to set automatic renewal.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration of purchasing an instance. Unit is month. Valid values are (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60).",
						},
						"renew_flag": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Automatic renewal logo. Values:\n" +
								"- `NOTIFY_AND_AUTO_RENEW`: notify expiration and renew automatically;\n" +
								"- `NOTIFY_AND_MANUAL_RENEW`: notification of expiration does not renew automatically. Users need to renew manually;\n" +
								"- `DISABLE_NOTIFY_AND_AUTO_RENEW`: no automatic renewal and no notification;\n" +
								"Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis after expiration, when the account balance is sufficient.",
						},
					},
				},
			},

			"renew_data_disk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to renew the data disk. Valid values:true: Indicates that the renewal instance also renews the data disk attached to it.false: Indicates that the instance will be renewed and the data disk attached to it will not be renewed at the same time.Default value: true.",
			},

			"auto_voucher": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
				Description: "Whether to automatically deduct vouchers. Valid values:\n" +
					"- true: Automatically deduct vouchers.\n" +
					"-false:Do not automatically deduct vouchers. Default value: false.",
			},
		},
	}
}

func resourceTencentCloudLighthouseRenewInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = lighthouse.NewRenewInstancesRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{&instanceId}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := lighthouse.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	if v, _ := d.GetOk("renew_data_disk"); v != nil {
		request.RenewDataDisk = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().RenewInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse renewInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseRenewInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseRenewInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseRenewInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
