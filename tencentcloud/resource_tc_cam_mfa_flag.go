/*
Provides a resource to create a cam mfa_flag

Example Usage

```hcl
resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = 20003xxxxxxx
  login_flag {
		phone = 0
		stoken = 1
		wechat = 0

  }
  action_flag {
		phone = 0
		stoken = 1
		wechat = 0

  }
}
```

Import

cam mfa_flag can be imported using the id, e.g.

```
terraform import tencentcloud_cam_mfa_flag.mfa_flag mfa_flag_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCamMfaFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamMfaFlagCreate,
		Read:   resourceTencentCloudCamMfaFlagRead,
		Update: resourceTencentCloudCamMfaFlagUpdate,
		Delete: resourceTencentCloudCamMfaFlagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"op_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Operate uin.",
			},

			"login_flag": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Login flag setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Phone.",
						},
						"stoken": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Soft token.",
						},
						"wechat": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Wechat.",
						},
					},
				},
			},

			"action_flag": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Action flag setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Phone.",
						},
						"stoken": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Soft token.",
						},
						"wechat": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Wechat.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCamMfaFlagCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.create")()
	defer inconsistentCheck(d, meta)()

	var opUin uint64
	if v, ok := d.GetOkExists("op_uin"); ok {
		opUin = v.(uint64)
	}

	d.SetId(helper.Int64ToStr(int64(opUin)))

	return resourceTencentCloudCamMfaFlagUpdate(d, meta)
}

func resourceTencentCloudCamMfaFlagRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	mfaFlagId := d.Id()

	mfaFlag, err := service.DescribeCamMfaFlagById(ctx, opUin)
	if err != nil {
		return err
	}

	if mfaFlag == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamMfaFlag` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mfaFlag.OpUin != nil {
		_ = d.Set("op_uin", mfaFlag.OpUin)
	}

	if mfaFlag.LoginFlag != nil {
		loginFlagMap := map[string]interface{}{}

		if mfaFlag.LoginFlag.Phone != nil {
			loginFlagMap["phone"] = mfaFlag.LoginFlag.Phone
		}

		if mfaFlag.LoginFlag.Stoken != nil {
			loginFlagMap["stoken"] = mfaFlag.LoginFlag.Stoken
		}

		if mfaFlag.LoginFlag.Wechat != nil {
			loginFlagMap["wechat"] = mfaFlag.LoginFlag.Wechat
		}

		_ = d.Set("login_flag", []interface{}{loginFlagMap})
	}

	if mfaFlag.ActionFlag != nil {
		actionFlagMap := map[string]interface{}{}

		if mfaFlag.ActionFlag.Phone != nil {
			actionFlagMap["phone"] = mfaFlag.ActionFlag.Phone
		}

		if mfaFlag.ActionFlag.Stoken != nil {
			actionFlagMap["stoken"] = mfaFlag.ActionFlag.Stoken
		}

		if mfaFlag.ActionFlag.Wechat != nil {
			actionFlagMap["wechat"] = mfaFlag.ActionFlag.Wechat
		}

		_ = d.Set("action_flag", []interface{}{actionFlagMap})
	}

	return nil
}

func resourceTencentCloudCamMfaFlagUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewSetMfaFlagRequest()

	mfaFlagId := d.Id()

	request.OpUin = &opUin

	immutableArgs := []string{"op_uin", "login_flag", "action_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("login_flag") {
		if dMap, ok := helper.InterfacesHeadMap(d, "login_flag"); ok {
			loginActionMfaFlag := cam.LoginActionMfaFlag{}
			if v, ok := dMap["phone"]; ok {
				loginActionMfaFlag.Phone = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["stoken"]; ok {
				loginActionMfaFlag.Stoken = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["wechat"]; ok {
				loginActionMfaFlag.Wechat = helper.IntUint64(v.(int))
			}
			request.LoginFlag = &loginActionMfaFlag
		}
	}

	if d.HasChange("action_flag") {
		if dMap, ok := helper.InterfacesHeadMap(d, "action_flag"); ok {
			loginActionMfaFlag := cam.LoginActionMfaFlag{}
			if v, ok := dMap["phone"]; ok {
				loginActionMfaFlag.Phone = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["stoken"]; ok {
				loginActionMfaFlag.Stoken = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["wechat"]; ok {
				loginActionMfaFlag.Wechat = helper.IntUint64(v.(int))
			}
			request.ActionFlag = &loginActionMfaFlag
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().SetMfaFlag(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam mfaFlag failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamMfaFlagRead(d, meta)
}

func resourceTencentCloudCamMfaFlagDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
