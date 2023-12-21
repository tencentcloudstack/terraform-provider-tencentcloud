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

func ResourceTencentCloudLighthouseRenewDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseRenewDiskCreate,
		Read:   resourceTencentCloudLighthouseRenewDiskRead,
		Delete: resourceTencentCloudLighthouseRenewDiskDelete,
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "List of disk ID.",
			},

			"renew_disk_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Renew cloud hard disk subscription related parameter settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Renewal period.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renewal falg. Value:NOTIFY_AND_AUTO_RENEW: Notice expires and auto-renews.NOTIFY_AND_MANUAL_RENEW: Notification expires without automatic renewal, users need to manually renew.DISABLE_NOTIFY_AND_AUTO_RENEW: No automatic renewal and no notification.Default: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the disk will be automatically renewed monthly when the account balance is sufficient.",
						},
						"time_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "newly purchased unit. Default: m.",
						},
						"cur_instance_deadline": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Current instance expiration time. Such as 2018-01-01 00:00:00. Specifying this parameter can align the expiration time of the instance attached to the disk. One of this parameter and Period must be specified, and cannot be specified at the same time.",
						},
					},
				},
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically use the voucher. Not used by default.",
			},
		},
	}
}

func resourceTencentCloudLighthouseRenewDiskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewRenewDisksRequest()
		diskId  string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		diskId = v.(string)
		request.DiskIds = []*string{&diskId}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "renew_disk_charge_prepaid"); ok {
		renewDiskChargePrepaid := lighthouse.RenewDiskChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			renewDiskChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			renewDiskChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		if v, ok := dMap["time_unit"]; ok {
			renewDiskChargePrepaid.TimeUnit = helper.String(v.(string))
		}
		if v, ok := dMap["cur_instance_deadline"]; ok {
			renewDiskChargePrepaid.CurInstanceDeadline = helper.String(v.(string))
		}
		request.RenewDiskChargePrepaid = &renewDiskChargePrepaid
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().RenewDisks(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse renewDisks failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(diskId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskLatestOperationRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseRenewDiskRead(d, meta)
}

func resourceTencentCloudLighthouseRenewDiskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseRenewDiskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
