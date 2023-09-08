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
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Required:    true,
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
	var opUin int

	if v, ok := d.GetOk("op_uin"); ok {
		opUin = v.(int)
	}
	d.SetId(strconv.Itoa(opUin))
	return resourceTencentCloudCamMfaFlagUpdate(d, meta)
}

func resourceTencentCloudCamMfaFlagRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	loginFlag, actionFlag, err := service.DescribeCamMfaFlagById(ctx)
	if err != nil {
		return err
	}

	if loginFlag == nil && actionFlag == nil {
		log.Printf("[WARN]%s resource `CamMfaFlag` not found, please check if it has been deleted.\n", logId)
		return nil
	}

	if loginFlag != nil {
		loginFlagMap := map[string]interface{}{}

		if loginFlag.Phone != nil {
			loginFlagMap["phone"] = loginFlag.Phone
		}

		if loginFlag.Stoken != nil {
			loginFlagMap["stoken"] = loginFlag.Stoken
		}

		if loginFlag.Wechat != nil {
			loginFlagMap["wechat"] = loginFlag.Wechat
		}

		_ = d.Set("login_flag", []interface{}{loginFlagMap})
	}

	if actionFlag != nil {
		actionFlagMap := map[string]interface{}{}

		if actionFlag.Phone != nil {
			actionFlagMap["phone"] = actionFlag.Phone
		}

		if actionFlag.Stoken != nil {
			actionFlagMap["stoken"] = actionFlag.Stoken
		}

		if actionFlag.Wechat != nil {
			actionFlagMap["wechat"] = actionFlag.Wechat
		}

		_ = d.Set("action_flag", []interface{}{actionFlagMap})
	}

	return nil
}

func resourceTencentCloudCamMfaFlagUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_mfa_flag.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	opUin := d.Id()
	request := cam.NewSetMfaFlagRequest()
	uin, err := strconv.Atoi(opUin)
	if err != nil {
		return err
	}
	request.OpUin = common.Uint64Ptr(uint64(uin))

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

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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
