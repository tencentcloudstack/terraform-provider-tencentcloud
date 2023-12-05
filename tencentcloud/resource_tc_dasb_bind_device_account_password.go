package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbBindDeviceAccountPassword() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_password.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().BindDeviceAccountPassword(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_password.read")()
	defer inconsistentCheck(d, meta)()

	if v, ok := d.GetOkExists("device_account_id"); ok {
		_ = d.Set("device_account_id", v.(int))
	}

	if v, ok := d.GetOk("password"); ok {
		_ = d.Set("password", v.(string))
	}

	return nil
}

func resourceTencentCloudDasbBindDeviceAccountPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_bind_device_account_password.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceAccountId = d.Id()
	)

	if err := service.DeleteDasbBindDeviceAccountPasswordById(ctx, deviceAccountId); err != nil {
		return err
	}

	return nil
}
