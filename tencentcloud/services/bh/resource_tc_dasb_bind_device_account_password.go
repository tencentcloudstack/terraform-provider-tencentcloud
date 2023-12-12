package bh

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbBindDeviceAccountPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbBindDeviceAccountPasswordCreate,
		Read:   resourceTencentCloudDasbBindDeviceAccountPasswordRead,
		Delete: resourceTencentCloudDasbBindDeviceAccountPasswordDelete,

		Schema: map[string]*schema.Schema{
			"device_account_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Host account ID.",
			},
			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Host account password.",
			},
		},
	}
}

func resourceTencentCloudDasbBindDeviceAccountPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_account_password.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		request         = dasb.NewBindDeviceAccountPasswordRequest()
		deviceAccountId string
	)

	if v, ok := d.GetOkExists("device_account_id"); ok {
		request.Id = helper.IntUint64(v.(int))
		deviceAccountId = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().BindDeviceAccountPassword(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb BindDeviceAccountPassword failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(deviceAccountId)

	return resourceTencentCloudDasbBindDeviceAccountPasswordRead(d, meta)
}

func resourceTencentCloudDasbBindDeviceAccountPasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_account_password.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	if v, ok := d.GetOkExists("device_account_id"); ok {
		_ = d.Set("device_account_id", v.(int))
	}

	if v, ok := d.GetOk("password"); ok {
		_ = d.Set("password", v.(string))
	}

	return nil
}

func resourceTencentCloudDasbBindDeviceAccountPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_account_password.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deviceAccountId = d.Id()
	)

	if err := service.DeleteDasbBindDeviceAccountPasswordById(ctx, deviceAccountId); err != nil {
		return err
	}

	return nil
}
