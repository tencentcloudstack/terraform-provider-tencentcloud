package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbAccount() *schema.Resource {
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

			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name, including alphanumeric _, Start with a letter, end with a letter or number, length 1-16.",
			},
			"account_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password, with a length range of 8 to 64 characters.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "main engine.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "describe.",
			},
			"max_user_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of user connections cannot be greater than 10240.",
			},
		},
	}
}

func resourceTencentCloudCynosdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = cynosdb.NewCreateAccountsRequest()
		clusterId   string
		accountName string
		host        string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	newAccount := cynosdb.NewAccount{}
	if v, ok := d.GetOk("account_name"); ok {
		accountName = v.(string)
		newAccount.AccountName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("account_password"); ok {
		newAccount.AccountPassword = helper.String(v.(string))
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
		newAccount.Host = helper.String(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		newAccount.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOk("max_user_connections"); ok {
		newAccount.MaxUserConnections = helper.IntInt64(v.(int))
	}
	request.Accounts = append(request.Accounts, &newAccount)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CreateAccounts(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + tccommon.FILED_SP + accountName + tccommon.FILED_SP + host)

	return resourceTencentCloudCynosdbAccountRead(d, meta)
}

func resourceTencentCloudCynosdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	accountName := idSplit[1]
	host := idSplit[2]

	account, err := service.DescribeCynosdbAccountById(ctx, clusterId, accountName, host)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if account.AccountName != nil {
		_ = d.Set("account_name", account.AccountName)
	}

	// if account.AccountPassword != nil {
	// 	_ = d.Set("account_password", account.AccountPassword)
	// }

	if account.Host != nil {
		_ = d.Set("host", account.Host)
	}

	if account.Description != nil {
		_ = d.Set("description", account.Description)
	}

	if account.MaxUserConnections != nil {
		_ = d.Set("max_user_connections", account.MaxUserConnections)
	}

	return nil
}

func resourceTencentCloudCynosdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	accountName := idSplit[1]
	host := idSplit[2]

	immutableArgs := []string{"cluster_id", "accounts"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("account_password") {
		request := cynosdb.NewResetAccountPasswordRequest()
		request.ClusterId = &clusterId
		request.AccountName = &accountName
		request.Host = &host
		if v, ok := d.GetOk("account_password"); ok {
			request.AccountPassword = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ResetAccountPassword(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb account failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("description") {
		request := cynosdb.NewModifyAccountDescriptionRequest()
		request.ClusterId = &clusterId
		request.AccountName = &accountName
		request.Host = &host
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyAccountDescription(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb account failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("host") {
		request := cynosdb.NewModifyAccountHostRequest()
		request.ClusterId = &clusterId
		request.Account = &cynosdb.InputAccount{
			AccountName: &accountName,
			Host:        &host,
		}

		var newHost string
		if v, ok := d.GetOk("host"); ok {
			newHost = v.(string)
			request.NewHost = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyAccountHost(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb account failed, reason:%+v", logId, err)
			return err
		}

		d.SetId(clusterId + tccommon.FILED_SP + accountName + tccommon.FILED_SP + newHost)
	}

	if d.HasChange("max_user_connections") {
		request := cynosdb.NewModifyAccountParamsRequest()
		request.ClusterId = &clusterId
		request.Account = &cynosdb.InputAccount{
			AccountName: &accountName,
			Host:        &host,
		}
		if v, ok := d.GetOk("max_user_connections"); ok {
			request.AccountParams = []*cynosdb.AccountParam{
				{
					ParamName:  helper.String("max_user_connections"),
					ParamValue: helper.String(fmt.Sprint(v)),
				},
			}

		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyAccountParams(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb account failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCynosdbAccountRead(d, meta)
}

func resourceTencentCloudCynosdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	accountName := idSplit[1]
	host := idSplit[2]

	if err := service.DeleteCynosdbAccountById(ctx, clusterId, accountName, host); err != nil {
		return err
	}

	return nil
}
