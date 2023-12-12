package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbDeviceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbDeviceAccountCreate,
		Read:   resourceTencentCloudDasbDeviceAccountRead,
		Delete: resourceTencentCloudDasbDeviceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Device ID.",
			},
			"account": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Device account.",
			},
		},
	}
}

func resourceTencentCloudDasbDeviceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		request         = dasb.NewCreateDeviceAccountRequest()
		response        = dasb.NewCreateDeviceAccountResponse()
		deviceAccountId string
	)

	if v, ok := d.GetOkExists("device_id"); ok {
		request.DeviceId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("account"); ok {
		request.Account = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateDeviceAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.Id != nil {
			e = fmt.Errorf("dasb DeviceAccount not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb DeviceAccount failed, reason:%+v", logId, err)
		return err
	}

	deviceAccountIdInt := *response.Response.Id
	deviceAccountId = strconv.FormatUint(deviceAccountIdInt, 10)
	d.SetId(deviceAccountId)

	return resourceTencentCloudDasbDeviceAccountRead(d, meta)
}

func resourceTencentCloudDasbDeviceAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deviceAccountId = d.Id()
	)

	DeviceAccount, err := service.DescribeDasbDeviceAccountById(ctx, deviceAccountId)
	if err != nil {
		return err
	}

	if DeviceAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbDeviceAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DeviceAccount.DeviceId != nil {
		_ = d.Set("device_id", DeviceAccount.DeviceId)
	}

	if DeviceAccount.Account != nil {
		_ = d.Set("account", DeviceAccount.Account)
	}

	return nil
}

func resourceTencentCloudDasbDeviceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deviceAccountId = d.Id()
	)

	if err := service.DeleteDasbDeviceAccountById(ctx, deviceAccountId); err != nil {
		return err
	}

	return nil
}
