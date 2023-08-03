/*
Use this resource to create SQL Server instance

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 4
  storage           = 100
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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentSqlServerBasicInfo(isROInstance bool) map[string]*schema.Schema {
	basicSchema := map[string]*schema.Schema{
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
			ValidateFunc: validateAllowedStringValue([]string{COMMON_PAYTYPE_PREPAID, COMMON_PAYTYPE_POSTPAID}),
			Description:  "Pay type of the SQL Server instance. Available values `PREPAID`, `POSTPAID_BY_HOUR`.",
		},
		"period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateIntegerInRange(1, 48),
			Description:  "Purchase instance period in month. The value does not exceed 48.",
		},
		"auto_voucher": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Whether to use the voucher automatically; 1 for yes, 0 for no, the default is 0.",
		},
		"voucher_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "An array of voucher IDs, currently only one can be used for a single order.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID of VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
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
			Description: "Readonly flag. `RO` (read-only instance), `MASTER` (primary instance with read-only instances). If it is left empty, it refers to an instance which is not read-only and has no RO group.",
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
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The tags of the SQL Server.",
		},
		"wait_switch": {
			Type:        schema.TypeInt,
			Optional:    true,
			Deprecated:  "It has been deprecated from version 1.81.2.",
			Description: "The way to execute the allocation. Supported values include: 0 - execute immediately, 1 - execute in maintenance window.",
		},
	}

	if !isROInstance {
		basicSchema["auto_renew"] = &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal (Default). Only valid when purchasing a prepaid instance.",
		}
	}

	return basicSchema
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
			Description: "Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.",
		},
		"ha_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "DUAL",
			Description: "Instance type. `DUAL` (dual-server high availability), `CLUSTER` (cluster). Default is `DUAL`.",
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
	basic := TencentSqlServerBasicInfo(false)
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceCreate,
		Read:   resourceTencentCloudSqlserverInstanceRead,
		Update: resourceTencentCloudSqlserverInstanceUpdate,
		Delete: resourceTencentCLoudSqlserverInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"auto_voucher": 0,
			}),
		},
		Schema: specialInfo,
	}
}

func resourceTencentCloudSqlserverInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	sqlserverService := SqlserverService{client: client}
	tagService := TagService{client: client}
	region := client.Region
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

	var instanceId string
	var outErr, inErr error

	if temp, ok := d.GetOkExists("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
	}

	request := sqlserver.NewCreateDBInstancesRequest()
	request.DBVersion = &dbVersion
	request.Memory = helper.IntInt64(memory)
	request.Storage = helper.IntInt64(storage)
	request.SubnetId = &subnetId
	request.VpcId = &vpcId
	request.HAType = &haType
	request.MultiZones = &multiZones

	if payType == COMMON_PAYTYPE_POSTPAID {
		request.InstanceChargeType = helper.String("POSTPAID")
	}
	if payType == COMMON_PAYTYPE_PREPAID {
		request.InstanceChargeType = helper.String("PREPAID")
		if v, ok := d.Get("auto_renew").(int); ok {
			request.AutoRenewFlag = helper.IntInt64(v)
		}

		if v, ok := d.Get("period").(int); ok {
			request.Period = helper.IntInt64(v)
		}
	}

	if v, ok := d.Get("auto_voucher").(int); ok {
		request.AutoVoucher = helper.IntInt64(v)
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIds := v.(*schema.Set).List()
		request.VoucherIds = helper.InterfacesStringsPoint(voucherIds)
	}

	if projectId != 0 {
		request.ProjectId = helper.IntInt64(projectId)
	}

	if len(weekSet) > 0 {
		request.Weekly = make([]*int64, 0)
		for _, i := range weekSet {
			request.Weekly = append(request.Weekly, helper.IntInt64(i))
		}
	}
	if startTime != "" {
		request.StartTime = &startTime
	}
	if timeSpan != 0 {
		request.Span = helper.IntInt64(timeSpan)
	}

	request.SecurityGroupList = make([]*string, 0, len(securityGroups))
	for _, v := range securityGroups {
		request.SecurityGroupList = append(request.SecurityGroupList, &v)
	}

	request.GoodsNum = helper.IntInt64(1)
	request.Zone = &zone

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverInstance(ctx, request)
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
	return resourceTencentCloudSqlserverInstanceRead(d, meta)
}

func sqlServerAllInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).apiV3Conn
	sqlserverService := SqlserverService{client: client}
	tagService := TagService{client: client}
	region := client.Region
	instanceId := d.Id()

	var outErr, inErr error

	//upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("auto_voucher") || d.HasChange("voucher_ids") {
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		autoVoucher := d.Get("auto_voucher").(int)
		voucherIds := d.Get("voucher_ids").(*schema.Set).List()
		outErr = sqlserverService.UpgradeSqlserverInstance(ctx, instanceId, memory, storage, autoVoucher, helper.InterfacesStringsPoint(voucherIds))

		if outErr != nil {
			return outErr
		}

	}

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
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return nil
}

func sqlServerAllInstanceNetUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewModifyDBInstanceNetworkRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
		instanceId  = d.Id()
	)

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		vpcId := d.Get("vpc_id").(string)
		subnetId := d.Get("subnet_id").(string)
		request.InstanceId = &instanceId
		request.NewVpcId = &vpcId
		request.NewSubnetId = &subnetId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceNetwork(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil {
				e = fmt.Errorf("sqlserver configInstanceNetwork not exists")
				return resource.NonRetryableError(e)
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
			return err
		}

		flowRequest.FlowId = &flowId
		err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DescribeFlowStatus(flowRequest)
			if e != nil {
				return retryError(e)
			}

			if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
				return nil
			} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
				return resource.RetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is running"))
			} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
				return resource.NonRetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is fail"))
			} else {
				e = fmt.Errorf("sqlserver configInstanceNetwork status illegal")
				return resource.NonRetryableError(e)
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s create sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
			return err
		}
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

	//update network
	if err := sqlServerAllInstanceNetUpdate(d, meta); err != nil {
		return err
	}

	var outErr, inErr error
	instanceId := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	sqlserverService := SqlserverService{client: client}
	tagService := TagService{client: client}
	region := client.Region
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

	}
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

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
	_ = d.Set("charge_type", helper.Int64ToStr(*instance.PayMode))

	if int(*instance.PayMode) == 1 {
		_ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
		if _, ok := d.GetOk("auto_renew"); ok {
			_ = d.Set("auto_renew", instance.RenewFlag)
		}
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

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

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

	return nil
}
