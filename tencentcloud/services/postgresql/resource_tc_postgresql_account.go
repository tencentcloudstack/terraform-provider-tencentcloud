package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlAccountCreate,
		Read:   resourceTencentCloudPostgresqlAccountRead,
		Update: resourceTencentCloudPostgresqlAccountUpdate,
		Delete: resourceTencentCloudPostgresqlAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of postgres-4wdeb0zv.",
			},
			"user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance username, which can contain 1-16 letters, digits, and underscore (_); can&amp;amp;#39;t be postgres; can&amp;amp;#39;t start with numbers, pg_, and tencentdb_.",
			},
			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Sensitive:   true,
				Description: "Password, which can contain 8-32 letters, digits, and symbols (()`~!@#$%^&amp;amp;amp;*-+=_|{}[]:;&amp;amp;#39;&amp;amp;lt;&amp;amp;gt;,.?/); can&amp;amp;#39;t start with slash /.",
			},
			"type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"normal", "tencentDBSuper"}),
				Description:  "The type of user. Valid values: 1. normal: regular user; 2. tencentDBSuper: user with the pg_tencentdb_superuser role.",
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks correspond to user `UserName`, which can contain 0-60 letters, digits, symbols (-_), and Chinese characters.",
			},
			"lock_status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "whether lock account. true: locked; false: unlock.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = postgresql.NewCreateAccountRequest()
		dBInstanceId string
		userName     string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgres account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dBInstanceId, userName}, tccommon.FILED_SP))

	// lock
	if v, ok := d.GetOkExists("lock_status"); ok {
		lockStatus := v.(bool)
		if lockStatus {
			lockRequest := postgresql.NewLockAccountRequest()
			lockRequest.DBInstanceId = &dBInstanceId
			lockRequest.UserName = &userName
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().LockAccount(lockRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, lockRequest.GetAction(), lockRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s lock postgres account failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudPostgresqlAccountRead(d, meta)
}

func resourceTencentCloudPostgresqlAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dBInstanceId := idSplit[0]
	userName := idSplit[1]

	account, err := service.DescribePostgresqlAccountById(ctx, dBInstanceId, userName)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if account.DBInstanceId != nil {
		_ = d.Set("db_instance_id", account.DBInstanceId)
	}

	if account.UserName != nil {
		_ = d.Set("user_name", account.UserName)
	}

	if account.UserType != nil {
		_ = d.Set("type", account.UserType)
	}

	if account.Remark != nil {
		_ = d.Set("remark", account.Remark)
	}

	if account.Status != nil {
		if *account.Status == 5 {
			_ = d.Set("lock_status", true)
		} else {
			_ = d.Set("lock_status", false)
		}
	}

	return nil
}

func resourceTencentCloudPostgresqlAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		pwdRequest    = postgresql.NewResetAccountPasswordRequest()
		remarkRequest = postgresql.NewModifyAccountRemarkRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dBInstanceId := idSplit[0]
	userName := idSplit[1]

	if d.HasChange("password") {
		pwdRequest.DBInstanceId = &dBInstanceId
		pwdRequest.UserName = &userName
		pwdRequest.Password = helper.String(d.Get("password").(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ResetAccountPassword(pwdRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, pwdRequest.GetAction(), pwdRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update postgres account password failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("remark") {
		remarkRequest.DBInstanceId = &dBInstanceId
		remarkRequest.UserName = &userName
		remarkRequest.Remark = helper.String(d.Get("remark").(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyAccountRemark(remarkRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, remarkRequest.GetAction(), remarkRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update postgres account remark failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("lock_status") {
		if v, ok := d.GetOkExists("lock_status"); ok {
			lockStatus := v.(bool)
			if lockStatus {
				lockRequest := postgresql.NewLockAccountRequest()
				lockRequest.DBInstanceId = &dBInstanceId
				lockRequest.UserName = &userName
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().LockAccount(lockRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, lockRequest.GetAction(), lockRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s lock postgres account failed, reason:%+v", logId, err)
					return err
				}
			} else {
				unlockRequest := postgresql.NewUnlockAccountRequest()
				unlockRequest.DBInstanceId = &dBInstanceId
				unlockRequest.UserName = &userName
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().UnlockAccount(unlockRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, unlockRequest.GetAction(), unlockRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s unlock postgres account failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudPostgresqlAccountRead(d, meta)
}

func resourceTencentCloudPostgresqlAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dBInstanceId := idSplit[0]
	userName := idSplit[1]

	if err := service.DeletePostgresqlAccountById(ctx, dBInstanceId, userName); err != nil {
		return err
	}

	return nil
}
