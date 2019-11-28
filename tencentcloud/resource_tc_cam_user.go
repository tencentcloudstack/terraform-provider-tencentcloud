/*
Provides a resource to create a CAM user.

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
}
```

Import

CAM user can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user.foo cam-user-test
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
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
				Description: "Name of CAM user.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remark of the CAM user.",
			},
			"use_api": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether to generate a secret key or not.",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicade whether the CAM user can login or not.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validateAsConfigPassword,
				Description:  "The password of the CAM user. The password should be set with 8 characters or more and contains uppercase small letters, numbers, and special characters. Only valid when console_login set true. If not set and the value of console_login is true, a random password is automatically generated.",
			},
			"need_reset_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether the CAM user will reset the password the next time he/her logs in.",
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "86",
				Description: "Phone num of the CAM user.",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country code of the phone num, like '86'.",
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
				Description: "Secret Id of the CAM user.",
			},
			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of the CAM user.",
			},
		},
	}
}

func resourceTencentCloudCamUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.create")()

	logId := getLogId(contextNil)
	request := cam.NewAddUserRequest()
	request.Name = stringToPointer(d.Get("name").(string))
	//optional
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = stringToPointer(v.(string))
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
		request.PhoneNum = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("password"); ok {
		request.Password = stringToPointer(v.(string))
	}

	var response *cam.AddUserResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().AddUser(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
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
		return fmt.Errorf("CAM user id is nil")
	}

	d.SetId(*response.Response.Name)
	d.Set("secret_key", *response.Response.SecretKey)
	d.Set("password", *response.Response.Password)
	d.Set("secret_id", *response.Response.SecretId)
	return resourceTencentCloudCamUserRead(d, meta)
}

func resourceTencentCloudCamUserRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	userId := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.GetUserResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserById(ctx, userId)
		if e != nil {
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

	d.Set("name", userId)
	d.Set("uin", int(*instance.Response.Uin))
	d.Set("uid", int(*instance.Response.Uid))
	d.Set("remark", *instance.Response.Remark)
	d.Set("phone_num", *instance.Response.PhoneNum)
	d.Set("country_code", *instance.Response.CountryCode)
	d.Set("email", *instance.Response.Email)
	if int(*instance.Response.ConsoleLogin) == 0 {
		d.Set("console_login", false)
	} else if int(*instance.Response.ConsoleLogin) == 1 {
		d.Set("console_login", true)
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
		request.Remark = stringToPointer(d.Get("remark").(string))
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
		request.PhoneNum = stringToPointer(d.Get("phone_num").(string))
		updateAttrs = append(updateAttrs, "phone_num")
	}
	if d.HasChange("country_code") {
		request.CountryCode = stringToPointer(d.Get("country_code").(string))
		updateAttrs = append(updateAttrs, "country_code")
	}
	if d.HasChange("email") {
		request.Email = stringToPointer(d.Get("email").(string))
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
