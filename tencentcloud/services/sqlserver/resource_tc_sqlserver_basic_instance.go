package sqlserver

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverBasicInstance() *schema.Resource {

	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBasicInstanceCreate,
		Read:   resourceTencentCloudSqlserverBasicInstanceRead,
		Update: resourceTencentCloudSqlserverBasicInstanceUpdate,
		Delete: resourceTencentCLoudSqlserverBasicInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the SQL Server basic instance.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The CPU number of the SQL Server basic instance.",
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
			"machine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD, CLOUD_BSSD}),
				Description:  "The host type of the purchased instance, `CLOUD_PREMIUM` for virtual machine high-performance cloud disk, `CLOUD_SSD` for virtual machine SSD cloud disk, `CLOUD_HSSD` for virtual machine enhanced cloud disk, `CLOUD_BSSD` for virtual machine general purpose SSD cloud disk.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      svcpostgresql.COMMON_PAYTYPE_POSTPAID,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{svcpostgresql.COMMON_PAYTYPE_PREPAID, svcpostgresql.COMMON_PAYTYPE_POSTPAID}),
				Description:  "Pay type of the SQL Server basic instance. For now, only `POSTPAID_BY_HOUR` is valid.",
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
			"engine_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "2008R2",
				Description: "Version of the SQL Server basic database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
				Description:  "Purchase instance period, the default value is 1, which means one month. The value does not exceed 48.",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group bound to the instance.",
			},
			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal, the default is 1 automatic renewal. Only valid when purchasing a prepaid instance.",
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
			"maintenance_week_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "A list of integer indicates weekly maintenance. For example, [1,7] presents do weekly maintenance on every Monday and Sunday.",
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
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Availability zone.",
			},
			"collation": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Chinese_PRC_CI_AS",
				Description: "System character set sorting rule, default: Chinese_PRC_CI_AS.",
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
				Description: "Create time of the SQL Server basic instance.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the SQL Server basic instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of the SQL Server basic instance.",
			},
			"dns_pod_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internet address domain name.",
			},
			"tgw_wan_vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "External port number.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBasicInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.create")()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client           = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		sqlserverService = SqlserverService{client: client}
		tagService       = svctag.NewTagService(client)
		region           = client.Region
		paramMap         = make(map[string]interface{})
		name             = d.Get("name").(string)
		payType          = d.Get("charge_type").(string)
		securityGroups   = make([]string, 0)
		voucherIds       = make([]string, 0)
		weekSet          = make([]int, 0)
	)
	if payType == svcpostgresql.COMMON_PAYTYPE_POSTPAID {
		payType = "POSTPAID"
		paramMap["autoRenew"] = 0
	} else {
		if v, ok := d.GetOk("auto_renew"); ok {
			paramMap["autoRenew"] = v.(int)
		} else {
			paramMap["autoRenew"] = 1
		}
	}
	paramMap["cpu"] = d.Get("cpu").(int)
	paramMap["memory"] = d.Get("memory").(int)
	paramMap["storage"] = d.Get("storage").(int)
	paramMap["subnetId"] = d.Get("subnet_id").(string)
	paramMap["vpcId"] = d.Get("vpc_id").(string)
	paramMap["machineType"] = d.Get("machine_type").(string)
	paramMap["payType"] = payType
	paramMap["engineVersion"] = d.Get("engine_version").(string)
	paramMap["period"] = d.Get("period").(int)
	paramMap["autoVoucher"] = d.Get("auto_voucher").(int)
	paramMap["availabilityZone"] = d.Get("availability_zone").(string)
	paramMap["collation"] = d.Get("collation").(string)

	if v, ok := d.GetOk("project_id"); ok {
		paramMap["projectId"] = v.(int)
	}
	if v, ok := d.GetOk("maintenance_start_time"); ok {
		paramMap["startTime"] = v.(string)
	}
	if v, ok := d.GetOk("maintenance_time_span"); ok {
		paramMap["timeSpan"] = v.(int)
	}
	// weekSet
	if v, ok := d.GetOk("maintenance_week_set"); ok {
		mWeekSet := v.(*schema.Set).List()
		for _, vv := range mWeekSet {
			weekSet = append(weekSet, vv.(int))
		}
		paramMap["weekSet"] = weekSet
	}
	// securityGroups
	if temp, ok := d.GetOk("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
		paramMap["securityGroups"] = securityGroups
	}
	// voucherIds
	if temp, ok := d.GetOk("voucher_ids"); ok {
		voucherId := temp.(*schema.Set).List()
		for _, id := range voucherId {
			voucherIds = append(voucherIds, id.(string))
		}
		paramMap["voucherIds"] = voucherIds
	}

	var instanceId string
	var outErr, inErr error
	outErr = resource.Retry(12*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverBasicInstance(ctx, paramMap, weekSet, voucherIds, securityGroups)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	//set name
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudSqlserverBasicInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBasicInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var outErr, inErr error
	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instance, has, outErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}
	chargeType := instance.PayMode
	_ = d.Set("cpu", instance.Cpu)
	_ = d.Set("subnet_id", instance.UniqSubnetId)
	_ = d.Set("vpc_id", instance.UniqVpcId)
	_ = d.Set("machine_type", instance.Type)
	if int(*chargeType) == 1 {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_PREPAID)
		_ = d.Set("auto_renew", instance.RenewFlag)
	} else {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_POSTPAID)
		_ = d.Set("auto_renew", 0)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("engine_version", instance.Version)

	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("status", instance.Status)
	_ = d.Set("cpu", instance.Cpu)
	_ = d.Set("memory", instance.Memory)
	_ = d.Set("storage", instance.Storage)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	if instance.DnsPodDomain != nil {
		_ = d.Set("dns_pod_domain", instance.DnsPodDomain)
	}

	if instance.TgwWanVPort != nil {
		_ = d.Set("tgw_wan_vport", instance.TgwWanVPort)
	}

	//maintanence
	var weekSet []int
	var startTime string
	var timeSpan int
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		weekSet, startTime, timeSpan, inErr = sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	_ = d.Set("maintenance_week_set", weekSet)
	_ = d.Set("maintenance_start_time", startTime)
	_ = d.Set("maintenance_time_span", timeSpan)

	var securityGroup []string
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		securityGroup, inErr = sqlserverService.DescribeInstanceSecurityGroups(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	_ = d.Set("security_groups", securityGroup)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSqlserverBasicInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region
	payType := d.Get("charge_type").(string)

	immutableArgs := []string{"collation"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	var outErr, inErr error
	instanceId := d.Id()
	d.Partial(true)
	//update name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}
	//upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") ||
		d.HasChange("cpu") || d.HasChange("auto_voucher") {
		voucherIds := make([]string, 0)
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		cpu := d.Get("cpu").(int)
		autoVoucher := d.Get("auto_voucher").(int)
		if temp, ok := d.GetOk("voucher_ids"); ok {
			voucherId := temp.(*schema.Set).List()
			for _, id := range voucherId {
				voucherIds = append(voucherIds, id.(string))
			}
		}
		outErr = resource.Retry(12*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.UpgradeSqlserverBasicInstance(ctx, instanceId, memory, storage, cpu, autoVoucher, voucherIds)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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
			outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.RemoveSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}
		for _, v := range newSet {
			sgId := v.(string)
			outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.AddSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}

	}
	//update project id
	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceMaintenanceSpan(ctx, instanceId, weekSet, startTime, timeSpan)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	if payType == svcpostgresql.COMMON_PAYTYPE_PREPAID {
		if d.HasChange("auto_renew") {
			var renewFlag int
			_, newValue := d.GetChange("auto_renew")
			renewFlag = newValue.(int)
			outErr = resource.Retry(2*tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr = sqlserverService.NewModifyDBInstanceRenewFlag(ctx, instanceId, renewFlag)
				if inErr != nil {
					return tccommon.RetryError(inErr)
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
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}
	d.Partial(false)
	return resourceTencentCloudSqlserverBasicInstanceRead(d, meta)
}

func resourceTencentCLoudSqlserverBasicInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	var has bool
	var instance *sqlserver.DBInstance

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
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
	// PREPAID
	if *instance.PayMode == 1 {
		return fmt.Errorf("PREPAID instances are not allowed to be deleted now, please terminate them on console")
	}
	//terminate sql instance
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server basic instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
