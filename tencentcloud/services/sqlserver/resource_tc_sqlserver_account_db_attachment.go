package sqlserver

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudSqlserverAccountDBAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverAccountDBAttachmentCreate,
		Read:   resourceTencentCloudSqlserverAccountDBAttachmentRead,
		Update: resourceTencentCloudSqlserverAccountDBAttachmentUpdate,
		Delete: resourceTencentCLoudSqlserverAccountDBAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "SQL Server instance ID that the account belongs to.",
			},
			"account_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "SQL Server account name.",
			},
			"db_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "SQL Server DB name.",
			},
			"privilege": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Privilege of the account on DB. Valid values: `ReadOnly`, `ReadWrite`.",
			},
		},
	}
}

func resourceTencentCloudSqlserverAccountDBAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account_db_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		accountName = d.Get("account_name").(string)
		dbName      = d.Get("db_name").(string)
		instanceId  = d.Get("instance_id").(string)
		privilege   = d.Get("privilege").(string)
	)

	var outErr, inErr error

	//check account exists
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, accountName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf(" SQL Server account %s, instance ID %s is not exist", accountName, instanceId))
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	//check db exists
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeDBDetailsById(ctx, instanceId+tccommon.FILED_SP+dbName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf(" SQL Server DB %s, instance ID %s is not exist", dbName, instanceId))
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.ModifyAccountDBAttachment(ctx, instanceId, accountName, dbName, privilege)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + tccommon.FILED_SP + accountName + tccommon.FILED_SP + dbName)

	return resourceTencentCloudSqlserverAccountDBAttachmentRead(d, meta)
}

func resourceTencentCloudSqlserverAccountDBAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account_db_attachment.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 3 {
		return fmt.Errorf("invalid SQL Server account DB attachment %s", id)
	}
	instanceId := idStrs[0]
	accountName := idStrs[1]
	dbName := idStrs[2]

	var outErr, inErr error

	//update privilege
	if d.HasChange("privilege") {
		privilege := d.Get("privilege").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifyAccountDBAttachment(ctx, instanceId, accountName, dbName, privilege)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudSqlserverAccountDBAttachmentRead(d, meta)
}

func resourceTencentCloudSqlserverAccountDBAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account_db_attachment.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 3 {
		return fmt.Errorf("invalid SQL Server account DB attachment ID %s", id)
	}
	instanceId := idStrs[0]
	accountName := idStrs[1]
	dbName := idStrs[2]

	var outErr, inErr error

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	attachment, has, outErr := sqlserverService.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
	if outErr != nil {
		inErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			attachment, has, inErr = sqlserverService.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("account_name", accountName)
	_ = d.Set("db_name", dbName)
	_ = d.Set("privilege", attachment["privilege"])

	return nil
}

func resourceTencentCLoudSqlserverAccountDBAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account_db_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 3 {
		return fmt.Errorf("invalid SQL Server account DB attachment id %s", id)
	}
	instanceId := idStrs[0]
	accountName := idStrs[1]
	dbName := idStrs[2]

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	var has bool
	privilege := "Delete"

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		return nil
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.ModifyAccountDBAttachment(ctx, instanceId, accountName, dbName, privilege)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server account DB attachment %s fail, account still exists from SDK DescribeSqlserverAccountDBAttachmentById", id)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
