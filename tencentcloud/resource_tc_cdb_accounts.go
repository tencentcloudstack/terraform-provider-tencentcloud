/*
Provides a resource to create a cdb accounts

Example Usage

```hcl
resource "tencentcloud_cdb_accounts" "accounts" {
  instance_id = ""
  accounts {
		user = ""
		host = ""

  }
  password = ""
  description = ""
  max_user_connections =
}
```

Import

cdb accounts can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_accounts.accounts accounts_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbAccounts() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbAccountsCreate,
		Read:   resourceTencentCloudCdbAccountsRead,
		Update: resourceTencentCloudCdbAccountsUpdate,
		Delete: resourceTencentCloudCdbAccountsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"accounts": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Cloud database account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "New account name.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The domain name of the new account.",
						},
					},
				},
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The password for the new account.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},

			"max_user_connections": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of available connections for a new account, the default value is 10240, and the maximum value that can be set is 10240.",
			},
		},
	}
}

func resourceTencentCloudCdbAccountsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_accounts.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewCreateAccountsRequest()
		response   = cdb.NewCreateAccountsResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("accounts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			account := cdb.Account{}
			if v, ok := dMap["user"]; ok {
				account.User = helper.String(v.(string))
			}
			if v, ok := dMap["host"]; ok {
				account.Host = helper.String(v.(string))
			}
			request.Accounts = append(request.Accounts, &account)
		}
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_user_connections"); ok {
		request.MaxUserConnections = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateAccounts(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb accounts failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbAccountsStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbAccountsRead(d, meta)
}

func resourceTencentCloudCdbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_accounts.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	accountsId := d.Id()

	accounts, err := service.DescribeCdbAccountsById(ctx, instanceId)
	if err != nil {
		return err
	}

	if accounts == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbAccounts` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accounts.InstanceId != nil {
		_ = d.Set("instance_id", accounts.InstanceId)
	}

	if accounts.Accounts != nil {
		accountsList := []interface{}{}
		for _, accounts := range accounts.Accounts {
			accountsMap := map[string]interface{}{}

			if accounts.Accounts.User != nil {
				accountsMap["user"] = accounts.Accounts.User
			}

			if accounts.Accounts.Host != nil {
				accountsMap["host"] = accounts.Accounts.Host
			}

			accountsList = append(accountsList, accountsMap)
		}

		_ = d.Set("accounts", accountsList)

	}

	if accounts.Password != nil {
		_ = d.Set("password", accounts.Password)
	}

	if accounts.Description != nil {
		_ = d.Set("description", accounts.Description)
	}

	if accounts.MaxUserConnections != nil {
		_ = d.Set("max_user_connections", accounts.MaxUserConnections)
	}

	return nil
}

func resourceTencentCloudCdbAccountsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_accounts.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyAccountHostRequest  = cdb.NewModifyAccountHostRequest()
		modifyAccountHostResponse = cdb.NewModifyAccountHostResponse()
	)

	accountsId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "accounts", "password", "description", "max_user_connections"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("accounts") {
		if v, ok := d.GetOk("accounts"); ok {
			for _, item := range v.([]interface{}) {
				account := cdb.Account{}
				if v, ok := dMap["user"]; ok {
					account.User = helper.String(v.(string))
				}
				if v, ok := dMap["host"]; ok {
					account.Host = helper.String(v.(string))
				}
				request.Accounts = append(request.Accounts, &account)
			}
		}
	}

	if d.HasChange("max_user_connections") {
		if v, ok := d.GetOkExists("max_user_connections"); ok {
			request.MaxUserConnections = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyAccountHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb accounts failed, reason:%+v", logId, err)
		return err
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbAccountsStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbAccountsStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbAccountsRead(d, meta)
}

func resourceTencentCloudCdbAccountsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_accounts.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	accountsId := d.Id()

	if err := service.DeleteCdbAccountsById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
