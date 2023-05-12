/*
Provides a resource to create a sqlserver general_communication

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_clone" "general_clone" {
  instance_id = "Instance ID in the format of mssql-j8kv137v"
  old_name    = "old_db_name"
  new_name    = "new_db_name"
}
```

Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.general_communication general_communication_id
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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverGeneralClone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloneCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloneRead,
		Update: resourceTencentCloudSqlserverGeneralCloneUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"old_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.",
			},
			"new_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.",
			},
			"db_detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Sqlserver db Clone detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collation_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "database collation.",
						},
						"is_auto_cleanup_on": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to automatically clean up after turning on CT 0: No 1: Yes.",
						},
						"is_broker_enabled": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Proxy enabled 0: No 1: Yes.",
						},
						"is_cdc_enabled": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether CDC is enabled/disabled 0: Disabled 1: Enabled.",
						},
						"is_db_chaining_on": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether CT is enabled/disabled 0: Disabled 1: Enabled.",
						},
						"is_encrypted": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to encrypt 0: No 1: Yes.",
						},
						"is_fulltext_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable full text 0: No 1: Yes.",
						},
						"is_mirroring": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether it is a mirror image 0: No 1: Yes.",
						},
						"is_published": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Published or not 0: No 1: Yes.",
						},
						"is_read_committed_snapshot_on": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable snapshot 0: No 1: Yes.",
						},
						"is_subscribed": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subscribed 0: No 1: Yes.",
						},
						"is_trust_worthy_on": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trustworthy 0: No 1: Yes.",
						},
						"mirroring_state": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mirror status.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "db name.",
						},
						"recovery_model_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "recovery mode.",
						},
						"retention_period": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "retention days.",
						},
						"state_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "database status.",
						},
						"user_access_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "user type.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_clone.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = sqlserver.NewCloneDBRequest()
		response   = sqlserver.NewCloneDBResponse()
		instanceId string
		newName    string
		flowId     int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	renameRestore := sqlserver.RenameRestoreDatabase{}
	if v, ok := d.GetOk("old_name"); ok {
		renameRestore.OldName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("new_name"); ok {
		renameRestore.NewName = helper.String(v.(string))
		newName = v.(string)
	}

	request.RenameRestore = append(request.RenameRestore, &renameRestore)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CloneDB(request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "FailedOperation.DBError" {
					e = fmt.Errorf("%s", sdkerr.Message)
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver clone %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		response = result
		flowId = *response.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver clone failed, reason:%+v", logId, err)
		return err
	}

	// wait for sqlserver clone done.
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver clone instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_CLONE_RUNNING {
			return resource.RetryableError(fmt.Errorf("create sqlserver clone task status is running"))
		}

		if *result.Status == SQLSERVER_CLONE_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_CLONE_FAIL {
			return resource.NonRetryableError(fmt.Errorf("create sqlserver clone task status is failed"))
		}

		e = fmt.Errorf("create sqlserver clone task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver clone task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{instanceId, newName}, FILED_SP))
	return resourceTencentCloudSqlserverGeneralCloneRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_clone.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId string
		dbName     string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]
	dbName = idSplit[1]

	generalClone, err := service.DescribeSqlserverGeneralCloneById(ctx, instanceId)
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0)
	for _, v := range generalClone {
		if *v.Name == dbName {
			var infoMap = map[string]interface{}{}
			if v.CollationName != nil {
				infoMap["collation_name"] = v.CollationName
			}
			if v.IsAutoCleanupOn != nil {
				infoMap["is_auto_cleanup_on"] = v.IsAutoCleanupOn
			}
			if v.IsBrokerEnabled != nil {
				infoMap["is_broker_enabled"] = v.IsBrokerEnabled
			}
			if v.IsCdcEnabled != nil {
				infoMap["is_cdc_enabled"] = v.IsCdcEnabled
			}
			if v.IsDbChainingOn != nil {
				infoMap["is_db_chaining_on"] = v.IsDbChainingOn
			}
			if v.IsEncrypted != nil {
				infoMap["is_encrypted"] = v.IsEncrypted
			}
			if v.IsFulltextEnabled != nil {
				infoMap["is_fulltext_enabled"] = v.IsFulltextEnabled
			}
			if v.IsMirroring != nil {
				infoMap["is_mirroring"] = v.IsMirroring
			}
			if v.IsPublished != nil {
				infoMap["is_published"] = v.IsPublished
			}
			if v.IsReadCommittedSnapshotOn != nil {
				infoMap["is_read_committed_snapshot_on"] = v.IsReadCommittedSnapshotOn
			}
			if v.IsSubscribed != nil {
				infoMap["is_subscribed"] = v.IsSubscribed
			}
			if v.IsTrustworthyOn != nil {
				infoMap["is_trust_worthy_on"] = v.IsTrustworthyOn
			}
			if v.MirroringState != nil {
				infoMap["mirroring_state"] = v.MirroringState
			}
			if v.Name != nil {
				infoMap["name"] = v.Name
			}
			if v.RecoveryModelDesc != nil {
				infoMap["recovery_model_desc"] = v.RecoveryModelDesc
			}
			if v.RetentionPeriod != nil {
				infoMap["retention_period"] = v.RetentionPeriod
			}
			if v.StateDesc != nil {
				infoMap["state_desc"] = v.StateDesc
			}
			if v.UserAccessDesc != nil {
				infoMap["user_access_desc"] = v.UserAccessDesc
			}
			list = append(list, infoMap)
			break
		}
	}
	_ = d.Set("db_detail", list)

	return nil
}

func resourceTencentCloudSqlserverGeneralCloneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_clone.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = sqlserver.NewModifyDBNameRequest()
		instanceId string
		dbName     string
		newName    string
		flowId     int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]
	dbName = idSplit[1]

	request.InstanceId = &instanceId
	request.OldDBName = &dbName

	if d.HasChange("new_name") {
		if v, ok := d.GetOk("new_name"); ok {
			request.NewDBName = helper.String(v.(string))
			newName = v.(string)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBName(request)
		if e != nil {

			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver clone %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver general clone failed, reason:%+v", logId, err)
		return err
	}

	// wait for sqlserver clone update done.
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver clone instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_CLONE_RUNNING {
			return resource.RetryableError(fmt.Errorf("update sqlserver clone task status is running"))
		}

		if *result.Status == SQLSERVER_CLONE_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_CLONE_FAIL {
			return resource.NonRetryableError(fmt.Errorf("update sqlserver clone task status is failed"))
		}

		e = fmt.Errorf("update sqlserver clone task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver clone task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{instanceId, newName}, FILED_SP))
	return resourceTencentCloudSqlserverGeneralCloneRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_clone.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId string
		dbName     string
		flowId     int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]
	dbName = idSplit[1]

	result, err := service.DeleteSqlserverGeneralCloneDB(ctx, instanceId, dbName)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("delete sqlserver clone task failed")
	}

	flowId = *result.Response.FlowId
	// wait for sqlserver clone delete done.
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		cloneStatus, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if cloneStatus == nil {
			e = fmt.Errorf("sqlserver clone instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *cloneStatus.Status == SQLSERVER_CLONE_RUNNING {
			return resource.RetryableError(fmt.Errorf("delete sqlserver clone task status is running"))
		}

		if *cloneStatus.Status == SQLSERVER_CLONE_SUCCESS {
			return nil
		}

		if *cloneStatus.Status == SQLSERVER_CLONE_FAIL {
			return resource.NonRetryableError(fmt.Errorf("delete sqlserver clone task status is failed"))
		}

		e = fmt.Errorf("delete sqlserver clone task status is %v, we won't wait for it finish", *cloneStatus.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete sqlserver clone task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
