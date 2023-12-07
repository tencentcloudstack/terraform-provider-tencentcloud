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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
		oldName    string
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
		oldName = v.(string)
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

	d.SetId(strings.Join([]string{instanceId, oldName, newName}, FILED_SP))

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
		oldName    string
		newName    string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]
	oldName = idSplit[1]
	newName = idSplit[2]

	generalClone, err := service.DescribeSqlserverGeneralCloneById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalClone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `sqlserver_general_clone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, v := range generalClone {
		if *v.Name == newName {
			_ = d.Set("instance_id", instanceId)
			_ = d.Set("old_name", oldName)
			_ = d.Set("new_name", v.Name)
		}
	}

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
		flowId     int64
		oldName    string
		newName    string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]

	request.InstanceId = &instanceId
	if d.HasChange("old_name") {
		if v, ok := d.GetOk("old_name"); ok {
			request.OldDBName = helper.String(v.(string))
			oldName = v.(string)
		}
	}

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

	d.SetId(strings.Join([]string{instanceId, oldName, newName}, FILED_SP))
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
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId = idSplit[0]
	dbName = idSplit[2]

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
