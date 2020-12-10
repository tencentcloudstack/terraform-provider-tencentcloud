/*
Provides a SQL Server instance resource to create read-only database instances.

Example Usage

```hcl
resource "tencentcloud_mysql_readonly_instance" "foo" {
  name = "tf_sqlserver_instance_ro"
  availability_zone = "ap-guangzhou-4"
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id                   = "` + defaultVpcId + `"
  subnet_id = "` + defaultSubnetId + `"
  memory = 2
  storage = 10
  master_instance_id = tencentcloud_sqlserver_instance.test.id
  readonly_group_type = 1
  force_upgrade = true
}
```

Import

SQL Server readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_readonly_instance.foo mssqlro-3cdq7kx5
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverReadonlyInstance() *schema.Resource {
	readonlyInstanceInfo := map[string]*schema.Schema{
		"master_instance_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Indicates the master instance ID of recovery instances.",
		},
		"readonly_group_type": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateAllowedIntValue([]int{1, 3}),
			Description:  "Type of readonly group. Valid values: `1`, `3`. `1` for one auto-assigned readonly instance per one readonly group, `2` for creating new readonly group, `3` for all exist readonly instances stay in the exist readonly group. For now, only `1` and `3` are supported.",
		},
		"force_upgrade": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Indicate that the master instance upgrade or not. `true` for upgrading the master SQL Server instance to cluster type by force. Default is false. Note: this is not supported with `DUAL`(ha_type), `2017`(engine_version) master SQL Server instance, for it will cause ha_type of the master SQL Server instance change.",
		},
		"readonly_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "ID of the readonly group that this instance belongs to. When `readonly_group_type` set value `3`, it must be set with valid value.",
		},
	}

	basic := TencentSqlServerBasicInfo()
	for k, v := range basic {
		readonlyInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudSqlserverReadonlyInstanceCreate,
		Read:   resourceTencentCloudSqlserverReadonlyInstanceRead,
		Update: resourceTencentCloudSqlserverReadonlyInstanceUpdate,
		Delete: resourceTencentCloudSqlserverReadonlyInstanceDelete,

		Schema: readonlyInstanceInfo,
	}
}

func resourceTencentCloudSqlserverReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	sqlserverService := SqlserverService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	var (
		name              = d.Get("name").(string)
		masterInstanceId  = d.Get("master_instance_id").(string)
		payType           = d.Get("charge_type").(string)
		readonlyGroupType = d.Get("readonly_group_type").(int)
		subnetId          = d.Get("subnet_id").(string)
		vpcId             = d.Get("vpc_id").(string)
		zone              = d.Get("availability_zone").(string)
		storage           = d.Get("storage").(int)
		memory            = d.Get("memory").(int)
		forceUpgrade      = d.Get("force_upgrade").(bool)
		readonlyGroupId   = ""
		securityGroups    = make([]string, 0)
	)

	if payType == COMMON_PAYTYPE_POSTPAID {
		payType = "POSTPAID"
	}

	if v, ok := d.GetOk("readonly_group_id"); ok && readonlyGroupType == 3 {
		readonlyGroupId = v.(string)
	}

	if temp, ok := d.GetOkExists("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
	}

	var instanceId string
	var outErr, inErr error

	outErr = resource.Retry(5*writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverReadonlyInstance(ctx, masterInstanceId, subnetId, vpcId, payType, memory, zone, storage, readonlyGroupType, readonlyGroupId, forceUpgrade, securityGroups)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	//set name
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudSqlserverReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_readonly_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var outErr, inErr error
	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	instance, has, err := tencentSqlServerBasicInfoRead(ctx, d, meta)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("engine_version", instance.Version)
	_ = d.Set("ha_type", SQLSERVER_HA_TYPE_FLAGS[*instance.HAFlag])

	//readonly group ID
	readonlyGroupId, masterInstanceId, outErr := sqlserverService.DescribeReadonlyGroupListByReadonlyInstanceId(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			readonlyGroupId, masterInstanceId, inErr = sqlserverService.DescribeReadonlyGroupListByReadonlyInstanceId(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	_ = d.Set("master_instance_id", masterInstanceId)
	_ = d.Set("readonly_group_id", readonlyGroupId)

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSqlserverReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_readonly_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	//basic info update
	if err := sqlServerAllInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}

	return resourceTencentCloudSqlserverReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var outErr, inErr error
	var has bool

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		return nil
	}

	//terminate sql instance
	outErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)

	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)

	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = sqlserverService.RecycleDBInstance(ctx, instanceId)
	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server readonly instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
