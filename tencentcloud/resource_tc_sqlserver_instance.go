/*
Use this resource to create SQL Server instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_instance" "foo" {
  name = "example"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id      = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
  project_id = 123
  memory = 2
  storage = 100
}
```

Import

SQL Server instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_instance.foo mssql-3cdq7kx5
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func TencentSqlServerBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(1, 60),
			Description:  "Name of the SQL Server instance.",
		},
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      COMMON_PAYTYPE_POSTPAID,
			ForceNew:     true,
			ValidateFunc: validateAllowedStringValue(POSTGRESQL_PAYTYPE),
			Description:  "Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "ID of VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "ID of subnet.",
		},
		"storage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.",
		},
		"memory": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Computed:    true,
			Description: "Availability zone.",
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Security group bound to the instance.",
		},
		//Computed values
		"ro_flag": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Readonly flag. `RO` for readonly instance, `MASTER` for master instance,  `` for not readonly instance.",
		},
		"vip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP for private access.",
		},
		"vport": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port for private access.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Create time of the SQL Server instance.",
		},
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.",
		},
	}
}

func resourceTencentCloudSqlserverInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"multi_zones": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to deploy across availability zones.",
		},
		//RO computed values
		"engine_version": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "2008R2",
			Description: "Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enerprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.",
		},
		"ha_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "DUAL",
			Description: "Instance type. Valid value are `DUAL`, `CLUSTER`. Default is `DUAL`.",
		},
		"maintenance_week_set": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "A list of integer indicates weekly maintenance. For example, [2,7] presents do weekly maintenance on every Tuesday and Sunday.",
		},
		"maintenance_start_time": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Start time of the maintenance in one day, format like `HH:mm`.",
		},
		"maintenance_time_span": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The timespan of maintenance in one day, unit is hour.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Project ID, default value is 0.",
		},
	}
	basic := TencentSqlServerBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceCreate,
		Read:   resourceTencentCloudSqlserverInstanceRead,
		Update: resourceTencentCloudSqlserverInstanceUpdate,
		Delete: resourceTencentCLoudSqlserverInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: specialInfo,
	}
}

func resourceTencentCloudSqlserverInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name           = d.Get("name").(string)
		dbVersion      = d.Get("engine_version").(string)
		payType        = d.Get("charge_type").(string)
		projectId      = d.Get("project_id").(int)
		subnetId       = d.Get("subnet_id").(string)
		vpcId          = d.Get("vpc_id").(string)
		zone           = d.Get("availability_zone").(string)
		storage        = d.Get("storage").(int)
		memory         = d.Get("memory").(int)
		weekSet        = make([]int, 0)
		startTime      = d.Get("maintenance_start_time").(string)
		timeSpan       = d.Get("maintenance_time_span").(int)
		multiZones     = d.Get("multi_zones").(bool)
		haType         = d.Get("ha_type").(string)
		securityGroups = make([]string, 0)
	)

	if v, ok := d.GetOk("maintenance_week_set"); ok {
		mWeekSet := v.(*schema.Set).List()
		for _, vv := range mWeekSet {
			weekSet = append(weekSet, vv.(int))
		}
	}
	if payType == COMMON_PAYTYPE_POSTPAID {
		payType = "POSTPAID"
	}
	var instanceId string
	var outErr, inErr error

	if temp, ok := d.GetOkExists("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
	}
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverInstance(ctx, dbVersion, payType, memory, 0, projectId, subnetId, vpcId, zone, storage, weekSet, startTime, timeSpan, multiZones, haType, securityGroups)
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

	return resourceTencentCloudSqlserverInstanceRead(d, meta)
}

func sqlServerAllInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	var outErr, inErr error
	//update name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		d.SetPartial("name")
	}

	//upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") {
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.UpgradeSqlserverInstance(ctx, instanceId, memory, storage)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("memory")
		d.SetPartial("storage")
	}

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		oldSet := os.List()
		newSet := ns.List()

		for _, v := range oldSet {
			sgId := v.(string)
			outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.RemoveSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}
		for _, v := range newSet {
			sgId := v.(string)
			outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.AddSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}

		d.SetPartial("security_groups")
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	d.Partial(true)

	//basic info update
	if err := sqlServerAllInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}

	var outErr, inErr error
	instanceId := d.Id()

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	//update project id
	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("project_id")
	}

	if d.HasChange("maintenance_week_set") || d.HasChange("maintenance_start_time") || d.HasChange("maintenance_time_span") {
		weekSet := make([]int, 0)
		if v, ok := d.GetOk("maintenance_week_set"); ok {
			mWeekSet := v.(*schema.Set).List()
			for _, vv := range mWeekSet {
				weekSet = append(weekSet, vv.(int))
			}
		}
		startTime := d.Get("maintenance_start_time").(string)
		timeSpan := d.Get("maintenance_time_span").(int)

		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceMaintenanceSpan(ctx, instanceId, weekSet, startTime, timeSpan)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("maintenance_week_set")
		d.SetPartial("maintenance_start_time")
		d.SetPartial("maintenance_time_span")
	}

	d.Partial(false)

	return resourceTencentCloudSqlserverInstanceRead(d, meta)
}
func tencentSqlServerBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (instance *sqlserver.DBInstance,
	has bool, errRet error) {

	if d.Id() == "" {
		return
	}
	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var outErr, inErr error
	instance, has, outErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		errRet = outErr
	}

	if !has {
		return
	}
	//basic info
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("vpc_id", instance.UniqVpcId)
	_ = d.Set("subnet_id", instance.UniqSubnetId)
	_ = d.Set("name", instance.Name)
	_ = d.Set("charge_type", instance.PayMode)

	if int(*instance.PayMode) == 1 {
		_ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
	} else {
		_ = d.Set("charge_type", COMMON_PAYTYPE_POSTPAID)
	}

	var securityGroup []string
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		securityGroup, inErr = sqlserverService.DescribeInstanceSecurityGroups(ctx, instanceId)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		errRet = outErr
	}
	//computed
	_ = d.Set("ro_flag", instance.ROFlag)
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("status", instance.Status)
	_ = d.Set("memory", instance.Memory)
	_ = d.Set("storage", instance.Storage)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	_ = d.Set("security_groups", securityGroup)
	return
}

func resourceTencentCloudSqlserverInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance.read")()

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

	//maintanence
	weekSet, startTime, timeSpan, outErr := sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			weekSet, startTime, timeSpan, inErr = sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	_ = d.Set("maintenance_week_set", weekSet)
	_ = d.Set("maintenance_start_time", startTime)
	_ = d.Set("maintenance_time_span", timeSpan)

	return nil
}

func resourceTencentCLoudSqlserverInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
