package dcdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbAccountRead,
		Create: resourceTencentCloudDcdbAccountCreate,
		Update: resourceTencentCloudDcdbAccountUpdate,
		Delete: resourceTencentCloudDcdbAccountDelete,
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
				Description: "account name.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db host.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "password.",
				Sensitive:   true,
			},

			"read_only": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "whether the account is readonly. 0 means not a readonly account.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description for account.",
			},

			"max_user_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "max user connections.",
			},
		},
	}
}

func resourceTencentCloudDcdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = dcdb.NewCreateAccountRequest()
		response   *dcdb.CreateAccountResponse
		instanceId string
		userName   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
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

	if v, ok := d.GetOk("read_only"); ok {
		request.ReadOnly = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {

		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_user_connections"); ok {
		request.MaxUserConnections = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().CreateAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb account failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId

	d.SetId(instanceId + tccommon.FILED_SP + userName)
	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	account, err := service.DescribeDcdbAccount(ctx, instanceId, userName)

	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		return fmt.Errorf("resource `account` %s does not exist", d.Id())
	}

	if account.InstanceId != nil {
		_ = d.Set("instance_id", account.InstanceId)
	}

	if account.Users[0] != nil {
		log.Printf("[DEBUG]tencentcloud_dcdb_account.read Users:%v", account.Users[0])
		if account.Users[0].UserName != nil {
			_ = d.Set("user_name", account.Users[0].UserName)
		}

		if account.Users[0].Host != nil {
			_ = d.Set("host", account.Users[0].Host)
		}

		if account.Users[0].ReadOnly != nil {
			_ = d.Set("read_only", account.Users[0].ReadOnly)
		}

		if account.Users[0].Description != nil {
			_ = d.Set("description", account.Users[0].Description)
		}
	}

	return nil
}

func resourceTencentCloudDcdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	// ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := dcdb.NewModifyAccountDescriptionRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	request.InstanceId = &instanceId
	request.UserName = &userName

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("user_name") {
		return fmt.Errorf("`user_name` do not support change now.")
	}

	if d.HasChange("password") {
		// return fmt.Errorf("`password` do not support change now.")
		if v, ok := d.GetOk("password"); ok {
			request := dcdb.NewResetAccountPasswordRequest()
			request.InstanceId = &instanceId
			request.UserName = &userName
			if v, ok := d.GetOk("host"); ok {
				request.Host = helper.String(v.(string))
			}
			request.Password = helper.String(v.(string))

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ResetAccountPassword(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate dcdb resetAccountPasswordOperation failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("read_only") {
		return fmt.Errorf("`read_only` do not support change now.")
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
		if v, ok := d.GetOk("host"); ok {
			request.Host = helper.String(v.(string))
		}
	}

	if d.HasChange("max_user_connections") {
		return fmt.Errorf("`max_user_connections` do not support change now.")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyAccountDescription(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb account failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	var host string
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	} else {
		return fmt.Errorf("host is broken, %s", host)
	}

	instanceId := idSplit[0]
	userName := idSplit[1]

	if err := service.DeleteDcdbAccountById(ctx, instanceId, userName, host); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"deleted"}, tccommon.ReadRetryTimeout, time.Second, service.DcdbAccountRefreshFunc(instanceId, userName, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
