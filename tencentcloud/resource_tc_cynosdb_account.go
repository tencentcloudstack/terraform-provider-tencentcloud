/*
Provides a resource to create a cynosdb account

Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id = "xxx"
  accounts {
		account_name = ""
		account_password = ""
		host = ""
		description = ""
		max_user_connections =

  }
}
```

Import

cynosdb account can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account.account account_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCynosdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbAccountCreate,
		Read:   resourceTencentCloudCynosdbAccountRead,
		Update: resourceTencentCloudCynosdbAccountUpdate,
		Delete: resourceTencentCloudCynosdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"accounts": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "New Account List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account name, including alphanumeric _, Start with a letter, end with a letter or number, length 1-16.",
						},
						"account_password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Password, with a length range of 8 to 64 characters.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Main engine.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Describe.",
						},
						"max_user_connections": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum number of user connections cannot be greater than 10240.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewCreateAccountsRequest()
		response  = cynosdb.NewCreateAccountsResponse()
		clusterId string
		account   string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("accounts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			newAccount := cynosdb.NewAccount{}
			if v, ok := dMap["account_name"]; ok {
				newAccount.AccountName = helper.String(v.(string))
			}
			if v, ok := dMap["account_password"]; ok {
				newAccount.AccountPassword = helper.String(v.(string))
			}
			if v, ok := dMap["host"]; ok {
				newAccount.Host = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				newAccount.Description = helper.String(v.(string))
			}
			if v, ok := dMap["max_user_connections"]; ok {
				newAccount.MaxUserConnections = helper.IntInt64(v.(int))
			}
			request.Accounts = append(request.Accounts, &newAccount)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateAccounts(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb account failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, account}, FILED_SP))

	return resourceTencentCloudCynosdbAccountRead(d, meta)
}

func resourceTencentCloudCynosdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	account := idSplit[1]

	account, err := service.DescribeCynosdbAccountById(ctx, clusterId, account)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if account.ClusterId != nil {
		_ = d.Set("cluster_id", account.ClusterId)
	}

	if account.Accounts != nil {
		accountsList := []interface{}{}
		for _, accounts := range account.Accounts {
			accountsMap := map[string]interface{}{}

			if account.Accounts.AccountName != nil {
				accountsMap["account_name"] = account.Accounts.AccountName
			}

			if account.Accounts.AccountPassword != nil {
				accountsMap["account_password"] = account.Accounts.AccountPassword
			}

			if account.Accounts.Host != nil {
				accountsMap["host"] = account.Accounts.Host
			}

			if account.Accounts.Description != nil {
				accountsMap["description"] = account.Accounts.Description
			}

			if account.Accounts.MaxUserConnections != nil {
				accountsMap["max_user_connections"] = account.Accounts.MaxUserConnections
			}

			accountsList = append(accountsList, accountsMap)
		}

		_ = d.Set("accounts", accountsList)

	}

	return nil
}

func resourceTencentCloudCynosdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		resetAccountPasswordRequest  = cynosdb.NewResetAccountPasswordRequest()
		resetAccountPasswordResponse = cynosdb.NewResetAccountPasswordResponse()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	account := idSplit[1]

	request.ClusterId = &clusterId
	request.Account = &account

	immutableArgs := []string{"cluster_id", "accounts"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ResetAccountPassword(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb account failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbAccountRead(d, meta)
}

func resourceTencentCloudCynosdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	account := idSplit[1]

	if err := service.DeleteCynosdbAccountById(ctx, clusterId, account); err != nil {
		return err
	}

	return nil
}
