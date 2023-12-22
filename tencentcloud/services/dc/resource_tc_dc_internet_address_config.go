package dc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
)

func ResourceTencentCloudDcInternetAddressConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcInternetAddressConfigCreate,
		Read:   resourceTencentCloudDcInternetAddressConfigRead,
		Update: resourceTencentCloudDcInternetAddressConfigUpdate,
		Delete: resourceTencentCloudDcInternetAddressConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "internet public address id.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "whether enable internet address.",
			},
		},
	}
}

func resourceTencentCloudDcInternetAddressConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_internet_address_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)

	d.SetId(instanceId)

	return resourceTencentCloudDcInternetAddressConfigUpdate(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_internet_address_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	internetAddressConfig, err := service.DescribeDcInternetAddressById(ctx, instanceId)
	if err != nil {
		return err
	}

	if internetAddressConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcInternetAddressConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if internetAddressConfig.InstanceId != nil {
		_ = d.Set("instance_id", internetAddressConfig.InstanceId)
	}

	if *internetAddressConfig.Status == 0 {
		_ = d.Set("enable", true)
	} else {
		_ = d.Set("enable", false)
	}

	return nil
}

func resourceTencentCloudDcInternetAddressConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_internet_address_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		enable         bool
		enableRequest  = dc.NewEnableInternetAddressRequest()
		disableRequest = dc.NewDisableInternetAddressRequest()
	)

	instanceId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.InstanceId = &instanceId

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcClient().EnableInternetAddress(enableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc internetAddressConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.InstanceId = &instanceId

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcClient().DisableInternetAddress(disableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc internetAddressConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDcInternetAddressConfigRead(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_internet_address_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
