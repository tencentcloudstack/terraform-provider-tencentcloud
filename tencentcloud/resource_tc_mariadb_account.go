/*
Provides a resource to create a mariadb account

Example Usage

```hcl
resource "tencentcloud_mariadb_account" "account" {
	instance_id = "tdsql-4pzs5b67"
	user_name   = "account-test"
	host        = "10.101.202.22"
	password    = "Password123."
	read_only   = 0
	description = "desc"
}

```
Import

mariadb account can be imported using the instance_id#user_name#host, e.g.
```
$ terraform import tencentcloud_mariadb_account.account tdsql-4pzs5b67#account-test#10.101.202.22
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbAccountRead,
		Create: resourceTencentCloudMariadbAccountCreate,
		Update: resourceTencentCloudMariadbAccountUpdate,
		Delete: resourceTencentCloudMariadbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "user name.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "host.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "account password.",
			},

			"read_only": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "wether account is read only, 0 means not a read only account.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "account description.",
			},
		},
	}
}

func resourceTencentCloudMariadbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewCreateAccountRequest()
		instanceId string
		userName   string
		host       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only"); ok {
		request.ReadOnly = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + userName + FILED_SP + host)
	return resourceTencentCloudMariadbAccountRead(d, meta)
}

func resourceTencentCloudMariadbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]
	host := idSplit[2]

	account, err := service.DescribeMariadbAccount(ctx, instanceId, userName, host)

	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		return fmt.Errorf("resource `account` %s does not exist", userName)
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("user_name", userName)
	_ = d.Set("host", host)

	// if account.Password != nil {
	// 	_ = d.Set("password", account.Password)
	// }

	if account.ReadOnly != nil {
		_ = d.Set("read_only", account.ReadOnly)
	}

	if account.Description != nil {
		_ = d.Set("description", account.Description)
	}

	return nil
}

func resourceTencentCloudMariadbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyAccountDescriptionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]
	host := idSplit[2]

	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &host

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("user_name") {
		return fmt.Errorf("`user_name` do not support change now.")
	}

	if d.HasChange("host") {
		return fmt.Errorf("`host` do not support change now.")
	}

	if d.HasChange("read_only") {
		return fmt.Errorf("`read_only` do not support change now.")
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyAccountDescription(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create mariadb account failed, reason:%+v", logId, err)
			return err
		}
	}

	// update pwd
	if d.HasChange("password") {
		PwdRequest := mariadb.NewResetAccountPasswordRequest()
		if v, ok := d.GetOk("password"); ok {
			PwdRequest.Password = helper.String(v.(string))
		}

		if v, ok := d.GetOk("user_name"); ok {
			PwdRequest.UserName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("host"); ok {
			PwdRequest.Host = helper.String(v.(string))
		}

		PwdRequest.InstanceId = &instanceId

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ResetAccountPassword(PwdRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s operate mariadb resetPassword failed, reason:%+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudMariadbAccountRead(d, meta)
}

func resourceTencentCloudMariadbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]
	host := idSplit[2]

	if err := service.DeleteMariadbAccountById(ctx, instanceId, userName, host); err != nil {
		return err
	}

	return nil
}
