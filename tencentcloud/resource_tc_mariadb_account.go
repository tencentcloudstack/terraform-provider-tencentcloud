/*
Provides a resource to create a mariadb account

Example Usage

```hcl
resource "tencentcloud_mariadb_account" "account" {
  instance_id = &lt;nil&gt;
  user_name = &lt;nil&gt;
  host = &lt;nil&gt;
  password = &lt;nil&gt;
  read_only = &lt;nil&gt;
  description = &lt;nil&gt;
}
```

Import

mariadb account can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_account.account account_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudMariadbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbAccountCreate,
		Read:   resourceTencentCloudMariadbAccountRead,
		Update: resourceTencentCloudMariadbAccountUpdate,
		Delete: resourceTencentCloudMariadbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User name.",
			},

			"host": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Host.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Account password.",
			},

			"read_only": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Wether account is read only, 0 means not a read only account.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Account description.",
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
		response   = mariadb.NewCreateAccountResponse()
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

	if v, ok := d.GetOkExists("read_only"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mariadb account failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, userName, host}, FILED_SP))

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

	account, err := service.DescribeMariadbAccountById(ctx, instanceId, userName, host)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if account.InstanceId != nil {
		_ = d.Set("instance_id", account.InstanceId)
	}

	if account.UserName != nil {
		_ = d.Set("user_name", account.UserName)
	}

	if account.Host != nil {
		_ = d.Set("host", account.Host)
	}

	if account.Password != nil {
		_ = d.Set("password", account.Password)
	}

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

	immutableArgs := []string{"instance_id", "user_name", "host", "password", "read_only", "description"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyAccountDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mariadb account failed, reason:%+v", logId, err)
		return err
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
