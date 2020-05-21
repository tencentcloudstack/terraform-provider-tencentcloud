/*
Use this resource to create tcaplus application

~> **NOTE:** tcaplus now only supports the following regions:ap-shanghai,ap-hongkong,na-siliconvalley,ap-singapore,ap-seoul,ap-tokyo,eu-frankfurt

Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
```

Import

tcaplus application can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_application.test 26655801
```

```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudTcaplusApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusApplicationCreate,
		Read:   resourceTencentCloudTcaplusApplicationRead,
		Update: resourceTencentCloudTcaplusApplicationUpdate,
		Delete: resourceTencentCloudTcaplusApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"idl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "Idl type of the tcapplus application.Valid values are " + strings.Join(TCAPLUS_IDL_TYPES, ",") + ".",
			},
			"app_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the tcapplus application. length should between 1 and 30.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id of the tcapplus application.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id of the tcapplus application.",
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 12 || len(value) > 16 {
						errors = append(errors, fmt.Errorf("invalid password, length should between 12 and 16"))
						return
					}
					var match = make(map[string]bool)
					for i := 0; i < len(value); i++ {
						if len(match) >= 2 {
							break
						}
						if value[i] >= '0' && value[i] <= '9' {
							match["number"] = true
							continue
						}
						if value[i] >= 'a' && value[i] <= 'z' {
							match["low"] = true
							continue
						}
						if value[i] >= 'A' && value[i] <= 'Z' {
							match["up"] = true
							continue
						}
					}
					if len(match) < 2 {
						errors = append(errors, fmt.Errorf("invalid password, a-z and 0-9 and A-Z must contain"))
					}
					return
				},
				Description: "Password of the tcapplus application. length should between 12 and 16,a-z and 0-9 and A-Z must contain.",
			},
			"old_password_expire_last": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validateIntegerMin(300),
				Description:  "Old password expected expiration seconds after change password,must >= 300.",
			},

			// Computed values.
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network type of the tcapplus application.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the tcapplus application.",
			},
			"password_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password status of the tcapplus application.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.",
			},
			"api_access_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access id of the tcapplus application.For TcaplusDB SDK connect.",
			},
			"api_access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access ip of the tcapplus application.For TcaplusDB SDK connect.",
			},
			"api_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access port of the tcapplus application.For TcaplusDB SDK connect.",
			},
			"old_password_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.",
			},
		},
	}
}

func resourceTencentCloudTcaplusApplicationCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_tcaplus_application.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		idlType  = d.Get("idl_type").(string)
		appName  = d.Get("app_name").(string)
		vpcId    = d.Get("vpc_id").(string)
		subnetId = d.Get("subnet_id").(string)
		password = d.Get("password").(string)
	)

	var applicationId string
	var inErr, outErr error

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		applicationId, inErr = tcaplusService.CreateApp(ctx, idlType, appName, vpcId, subnetId, password)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(applicationId)
	time.Sleep(3 * time.Second)
	return resourceTencentCloudTcaplusApplicationRead(d, meta)
}

func resourceTencentCloudTcaplusApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_application.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	appInfo, has, err := tcaplusService.DescribeApp(ctx, d.Id())
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			appInfo, has, err = tcaplusService.DescribeApp(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("idl_type", appInfo.IdlType)
	_ = d.Set("app_name", appInfo.ClusterName)
	_ = d.Set("vpc_id", appInfo.VpcId)
	_ = d.Set("subnet_id", appInfo.SubnetId)
	_ = d.Set("network_type", appInfo.NetworkType)
	_ = d.Set("create_time", appInfo.CreatedTime)
	_ = d.Set("password_status", appInfo.PasswordStatus)
	_ = d.Set("api_access_id", appInfo.ApiAccessId)
	_ = d.Set("api_access_ip", appInfo.ApiAccessIp)
	_ = d.Set("api_access_port", appInfo.ApiAccessPort)

	if appInfo.OldPasswordExpireTime == nil || *appInfo.OldPasswordExpireTime == "" {
		_ = d.Set("old_password_expire_time", "-")
	} else {
		_ = d.Set("old_password_expire_time", appInfo.OldPasswordExpireTime)
	}

	return nil
}

func resourceTencentCloudTcaplusApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_application.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	d.Partial(true)

	if d.HasChange("app_name") {
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyAppName(ctx, d.Id(), d.Get("app_name").(string))
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("app_name")
	}

	if d.HasChange("password") {
		oldPwd, newPwd := d.GetChange("password")
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyAppPassword(ctx, d.Id(),
				oldPwd.(string),
				newPwd.(string),
				int64(d.Get("old_password_expire_last").(int)))

			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "FailedOperation.OldPasswordInUse" {
					err = fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.OldPasswordInUse,`password_status` is unmodifiable now, can modify after `old_password_expire_time`")
					return resource.NonRetryableError(err)
				}
			}
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("password")
	}
	if d.HasChange("old_password_expire_last") {
		d.SetPartial("old_password_expire_last")
	}
	d.Partial(false)

	return resourceTencentCloudTcaplusApplicationRead(d, meta)
}

func resourceTencentCloudTcaplusApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_application.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	_, err := tcaplusService.DeleteApp(ctx, d.Id())

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err = tcaplusService.DeleteApp(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeApp(ctx, d.Id())
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeApp(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}

			if has {
				err = fmt.Errorf("delete application fail, application still exist from sdk DescribeApps")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete application fail, application still exist from sdk DescribeApps")
	}
}
