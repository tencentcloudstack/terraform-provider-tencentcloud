/*
Provides a resource to create a dcdb account

Example Usage

```hcl
resource "tencentcloud_dcdb_account" "account" {
  instance_id = &lt;nil&gt;
  user_name = &lt;nil&gt;
  host = &lt;nil&gt;
  password = &lt;nil&gt;
  read_only = &lt;nil&gt;
  description = &lt;nil&gt;
  max_user_connections = &lt;nil&gt;
}
```

Import

dcdb account can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account.account account_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDcdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbAccountCreate,
		Read:   resourceTencentCloudDcdbAccountRead,
		Update: resourceTencentCloudDcdbAccountUpdate,
		Delete: resourceTencentCloudDcdbAccountDelete,
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
				Description: "Account name.",
			},

			"host": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Db host.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Password.",
			},

			"read_only": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether the account is readonly. 0 means not a readonly account.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description for account.",
			},

			"max_user_connections": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Max user connections.",
			},
		},
	}
}

func resourceTencentCloudDcdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCreateAccountRequest()
		response   = dcdb.NewCreateAccountResponse()
		instanceId string
		userName   string
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

	if v, ok := d.GetOkExists("max_user_connections"); ok {
		request.MaxUserConnections = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dcdb account failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, userName}, FILED_SP))

	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	account, err := service.DescribeDcdbAccountById(ctx, instanceId, userName)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

	if account.MaxUserConnections != nil {
		_ = d.Set("max_user_connections", account.MaxUserConnections)
	}

	return nil
}

func resourceTencentCloudDcdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyAccountDescriptionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	request.InstanceId = &instanceId
	request.UserName = &userName

	immutableArgs := []string{"instance_id", "user_name", "host", "password", "read_only", "description", "max_user_connections"}

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyAccountDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb account failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	if err := service.DeleteDcdbAccountById(ctx, instanceId, userName); err != nil {
		return err
	}

	return nil
}
