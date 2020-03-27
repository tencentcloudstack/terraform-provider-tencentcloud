/*
Provides a resource to manage CAM user.

Example Usage

```hcl
resource "tencentcloud_cam_user" "foo" {
  name                = "cam-user-test"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  email               = "hello@test.com"
  country_code        = "86"
  force_delete        = true
}
```

Import

CAM user can be imported using the user name, e.g.

```
$ terraform import tencentcloud_cam_user.foo cam-user-test
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserCreate,
		Read:   resourceTencentCloudCamUserRead,
		Update: resourceTencentCloudCamUserUpdate,
		Delete: resourceTencentCloudCamUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the CAM user.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remark of the CAM user.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether to force deletes the CAM user. If set false, the API secret key will be checked and failed when exists; otherwise the user will be deleted directly. Default is false.",
			},
			"use_api": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether to generate the API secret key or not.",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether the CAM user can login to the web console or not.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validateAsConfigPassword,
				Description:  "The password of the CAM user. Password should be at least 8 characters and no more than 32 characters, includes uppercase letters, lowercase letters, numbers and special characters. Only required when `console_login` is true. If not set, a random password will be automatically generated.",
			},
			"need_reset_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether the CAM user need to reset the password when first logins.",
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number of the CAM user.",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "86",
				Description: "Country code of the phone number, for example: '86'.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email of the CAM user.",
			},
			"uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Uin of the CAM User.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Secret key of the CAM user.",
			},
			"secret_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Secret ID of the CAM user.",
			},
			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the CAM user.",
			},
		},
	}
}

func resourceTencentCloudCamUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.create")()

	logId := getLogId(contextNil)

	request := cam.NewAddUserRequest()
	request.Name = helper.String(d.Get("name").(string))
	//optional
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("use_api"); ok {
		apiBool := v.(bool)
		apiInt := uint64(1)
		if !apiBool {
			apiInt = uint64(0)
		}
		request.UseApi = &apiInt
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		loginBool := v.(bool)
		loginInt := uint64(1)
		if !loginBool {
			loginInt = uint64(0)
		}
		request.ConsoleLogin = &loginInt
	}
	if v, ok := d.GetOkExists("need_reset_password"); ok {
		resetBool := v.(bool)
		resetInt := uint64(1)
		if !resetBool {
			resetInt = uint64(0)
		}
		request.NeedResetPassword = &resetInt
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request.PhoneNum = helper.String(v.(string))
	}
	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}
	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	var response *cam.AddUserResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().AddUser(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				if strings.Contains(errCode, "SubUserNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create user failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.Uid == nil {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CAM user][Create] check: CAM user id returns nil")
	}

	d.SetId(*response.Response.Name)
	_ = d.Set("secret_key", *response.Response.SecretKey)
	_ = d.Set("password", *response.Response.Password)
	_ = d.Set("secret_id", *response.Response.SecretId)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), "logId", logId)
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeUserById(ctx, *response.Response.Name)
		if e != nil {
			return retryError(e, "ResourceNotFound")
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for CAM user ready failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(3 * time.Second)
	return resourceTencentCloudCamUserRead(d, meta)
}

func resourceTencentCloudCamUserRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	deleteForce := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
		_ = d.Set("force_delete", deleteForce)
	}

	userId := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.GetUserResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserById(ctx, userId)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return nil
				}
			}
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.Uid == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", userId)
	_ = d.Set("uin", int(*instance.Response.Uin))
	_ = d.Set("uid", int(*instance.Response.Uid))
	_ = d.Set("remark", *instance.Response.Remark)
	_ = d.Set("phone_num", *instance.Response.PhoneNum)
	_ = d.Set("country_code", *instance.Response.CountryCode)
	_ = d.Set("email", *instance.Response.Email)
	if int(*instance.Response.ConsoleLogin) == 0 {
		_ = d.Set("console_login", false)
	} else if int(*instance.Response.ConsoleLogin) == 1 {
		_ = d.Set("console_login", true)
	}

	return nil
}

func resourceTencentCloudCamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.update")()

	logId := getLogId(contextNil)

	userId := d.Id()

	var updateAttrs []string

	request := cam.NewUpdateUserRequest()
	request.Name = &userId
	if d.HasChange("remark") {
		request.Remark = helper.String(d.Get("remark").(string))
		updateAttrs = append(updateAttrs, "remark")
	}

	if d.HasChange("console_login") {
		consoleLogin := d.Get("console_login").(bool)
		consoleLogin64 := uint64(0)
		if consoleLogin {
			consoleLogin64 = uint64(1)
		}
		request.ConsoleLogin = &consoleLogin64
		updateAttrs = append(updateAttrs, "console_login")
	}

	if d.HasChange("need_reset_password") {
		resetBool := d.Get("need_reset_password").(bool)
		resetBool64 := uint64(0)
		if resetBool {
			resetBool64 = uint64(1)
		}
		request.NeedResetPassword = &resetBool64
		updateAttrs = append(updateAttrs, "need_reset_password")
	}

	if d.HasChange("phone_num") {
		request.PhoneNum = helper.String(d.Get("phone_num").(string))
		updateAttrs = append(updateAttrs, "phone_num")
	}
	if d.HasChange("country_code") {
		request.CountryCode = helper.String(d.Get("country_code").(string))
		updateAttrs = append(updateAttrs, "country_code")
	}
	if d.HasChange("email") {
		request.Email = helper.String(d.Get("email").(string))
		updateAttrs = append(updateAttrs, "email")
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateUser(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM user description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return nil
}

func resourceTencentCloudCamUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.delete")()

	logId := getLogId(contextNil)

	userId := d.Id()
	request := cam.NewDeleteUserRequest()
	request.Name = &userId

	//check is force delete or not
	deleteForce := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
	}

	request.Force = helper.BoolToInt64Pointer(deleteForce)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DeleteUser(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM user failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
