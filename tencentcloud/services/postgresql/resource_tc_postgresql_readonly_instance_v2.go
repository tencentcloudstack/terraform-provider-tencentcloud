package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudPostgresqlReadonlyInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadonlyInstanceV2Create,
		Read:   resourceTencentCloudPostgresqlReadonlyInstanceV2Read,
		Update: resourceTencentCloudPostgresqlReadonlyInstanceV2Update,
		Delete: resourceTencentCloudPostgresqlReadonlyInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone ID, such as ap-guangzhou-3.",
			},

			"master_db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the primary instance to which the read-only replica belongs.",
			},

			"spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specification code, which can be obtained via DescribeClasses.",
			},

			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance storage capacity in GB, the step is 10.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Purchase duration in months. PREPAID supports 1,2,3,4,5,6,7,8,9,10,11,12,24,36; POSTPAID only supports 1.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC ID, such as vpc-xxxxxxxx.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC subnet ID, such as subnet-xxxxxxxx.",
			},

			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance billing mode. Valid values: PREPAID, POSTPAID_BY_HOUR. Default: PREPAID.",
			},

			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to use voucher automatically, 1 for yes, 0 for no. Default: 0.",
			},

			"voucher_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Voucher ID list, currently only one voucher is supported.",
			},

			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Auto renew flag, 0 for manual renew, 1 for auto renew. Default: 0. Only supports PREPAID.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Project ID. Default: 0, means default project.",
			},

			"read_only_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Read-only group ID.",
			},

			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group IDs bound to the instance.",
			},

			"need_support_ipv6": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Whether to support IPv6 access, 1 for yes, 0 for no. Default: 0.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance name, only supports chinese/english/numbers/_/- with length less than 60.",
			},

			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Dedicated cluster ID.",
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable deletion protection, true for enable, false for disable.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},

			// computed
			"db_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID managed by this resource.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DB instance CPU.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DB instance memory.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = postgresql.NewCreateReadOnlyDBInstanceRequest()
		response   = postgresql.NewCreateReadOnlyDBInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntUint64(v.(int))
	}

	// set default period 1.
	request.Period = helper.IntUint64(1)
	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		request.VoucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroupIds))
		for _, item := range securityGroupIds {
			if sgId, ok := item.(string); ok {
				request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(sgId))
			}
		}
	}

	if v, ok := d.GetOkExists("need_support_ipv6"); ok {
		request.NeedSupportIpv6 = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request.DeletionProtection = helper.Bool(v.(bool))
	}

	request.InstanceCount = helper.Uint64(1)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateReadOnlyDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql readonly instance v2 failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql readonly instance v2 failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.DBInstanceIdSet) == 0 {
		if len(response.Response.DealNames) == 0 {
			return fmt.Errorf("TencentCloud SDK returns empty postgresql ID and Deals")
		}

		// get instance id from deal
		orderReq := postgresql.NewDescribeOrdersRequest()
		orderResp := postgresql.NewDescribeOrdersResponse()
		dealId := response.Response.DealNames[0]
		orderReq.DealNames = []*string{dealId}
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeOrders(orderReq)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe orders failed, Response is nil."))
			}

			if len(result.Response.Deals) == 0 {
				return resource.RetryableError(fmt.Errorf("waiting for deal return instance id"))
			}

			orderResp = result
			return nil
		})

		if err != nil {
			return err
		}

		deals := orderResp.Response.Deals
		if len(deals) > 0 && len(deals[0].DBInstanceIdSet) > 0 {
			if deals[0].DBInstanceIdSet[0] != nil {
				instanceId = *deals[0].DBInstanceIdSet[0]
			}
		}
	} else {
		if response.Response.DBInstanceIdSet[0] != nil {
			instanceId = *response.Response.DBInstanceIdSet[0]
		}
	}

	if instanceId == "" {
		return fmt.Errorf("TencentCloud SDK returns empty postgresql ID")
	}

	d.SetId(instanceId)

	// wait
	waitReq := postgresql.NewDescribeDBInstanceAttributeRequest()
	waitReq.DBInstanceId = &instanceId
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeDBInstanceAttribute(waitReq)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if result == nil || result.Response == nil || result.Response.DBInstance == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe db instance attribute failed, Response is nil."))
		}

		instance := result.Response.DBInstance
		if instance.DBInstanceStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe db instance attribute failed, DBInstanceStatus is nil."))
		}

		if *instance.DBInstanceStatus == "running" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Postgresql readonly instance v2 %s is still creating, status is %s.", instanceId, *instance.DBInstanceStatus))
	})

	if err != nil {
		return err
	}

	// handle tags via tag service if provided
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d, meta)
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()
	instance, err := service.DescribePostgresqlReadonlyInstanceV2ById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgresql_readonly_instance_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if instance.DBInstanceId != nil {
		_ = d.Set("db_instance_id", *instance.DBInstanceId)
	}

	if instance.Zone != nil {
		_ = d.Set("zone", *instance.Zone)
	}

	if instance.MasterDBInstanceId != nil {
		_ = d.Set("master_db_instance_id", *instance.MasterDBInstanceId)
	}

	if instance.DBInstanceClass != nil {
		_ = d.Set("spec_code", *instance.DBInstanceClass)
	}

	if instance.DBInstanceStorage != nil {
		_ = d.Set("storage", *instance.DBInstanceStorage)
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", *instance.VpcId)
	}

	if instance.SubnetId != nil {
		_ = d.Set("subnet_id", *instance.SubnetId)
	}

	if instance.PayType != nil {
		if *instance.PayType == "prepaid" {
			_ = d.Set("instance_charge_type", "PREPAID")
		} else if *instance.PayType == "postpaid" {
			_ = d.Set("instance_charge_type", "POSTPAID_BY_HOUR")
		}
	}

	if instance.AutoRenew != nil {
		_ = d.Set("auto_renew_flag", *instance.AutoRenew)
	}

	if instance.ProjectId != nil {
		_ = d.Set("project_id", *instance.ProjectId)
	}

	if instance.SupportIpv6 != nil {
		_ = d.Set("need_support_ipv6", *instance.SupportIpv6)
	}

	if instance.DBInstanceName != nil {
		_ = d.Set("name", *instance.DBInstanceName)
	}

	if instance.DeletionProtection != nil {
		_ = d.Set("deletion_protection", *instance.DeletionProtection)
	}

	if instance.DBInstanceCpu != nil {
		_ = d.Set("cpu", *instance.DBInstanceCpu)
	}

	if instance.DBInstanceMemory != nil {
		_ = d.Set("memory", *instance.DBInstanceMemory)
	}

	if len(instance.DBNodeSet) > 0 {
		for _, node := range instance.DBNodeSet {
			if node != nil && node.DedicatedClusterId != nil {
				_ = d.Set("dedicated_cluster_id", *node.DedicatedClusterId)
				break
			}
		}
	}

	// read only group
	if instance.MasterDBInstanceId != nil {
		readOnlyGroupId, roErr := service.DescribeReadOnlyGroupsById(ctx, *instance.MasterDBInstanceId, instanceId)
		if roErr != nil {
			return roErr
		}

		if readOnlyGroupId != nil {
			_ = d.Set("read_only_group_id", *readOnlyGroupId)
		}
	}

	// security groups
	sg, sgErr := service.DescribeDBInstanceSecurityGroupsById(ctx, instanceId)
	if sgErr != nil {
		return sgErr
	}

	if len(sg) > 0 {
		_ = d.Set("security_group_ids", sg)
	}

	// tags via tag service
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, tagErr := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, instanceId)
	if tagErr != nil {
		return tagErr
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                 = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId              = d.Id()
		outErr, inErr, checkErr error
	)

	// wait sync time
	waitTime := 5 * time.Second

	if d.HasChange("read_only_group_id") {
		var (
			request          = postgresql.NewModifyDBInstanceReadOnlyGroupRequest()
			masterInstanceId string
			roGroupIdOld     string
			roGroupIdNew     string
		)

		masterInstanceId = d.Get("master_db_instance_id").(string)
		old, new := d.GetChange("read_only_group_id")
		if old != nil {
			roGroupIdOld = old.(string)
		}

		if new != nil {
			roGroupIdNew = new.(string)
		}

		request.DBInstanceId = &instanceId
		request.ReadOnlyGroupId = &roGroupIdOld
		request.NewReadOnlyGroupId = &roGroupIdNew
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceReadOnlyGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s operate postgresql ChangeDbInstanceReadOnlyGroupOperation failed, reason:%+v", logId, err)
			return err
		}

		// need wait sync
		time.Sleep(waitTime)

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"ok"}, d.Timeout(schema.TimeoutUpdate), time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(masterInstanceId, roGroupIdNew, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = service.ModifyPostgresqlInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// need wait sync
		time.Sleep(waitTime)

		// check update name done
		timeoutMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		checkErr = service.CheckDBInstanceStatus(ctx, instanceId, timeoutMinutes)
		if checkErr != nil {
			return checkErr
		}

	}

	if d.HasChange("storage") {
		storage := d.Get("storage").(int)
		cpu := d.Get("cpu").(int)
		memory := d.Get("memory").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = service.UpgradePostgresqlInstance(ctx, instanceId, memory, storage, cpu, 0)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// need wait sync
		time.Sleep(waitTime)

		// check update storage and memory done
		timeoutMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		checkErr = service.CheckDBInstanceStatus(ctx, instanceId, timeoutMinutes)
		if checkErr != nil {
			return checkErr
		}
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = service.ModifyPostgresqlInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// need wait sync
		time.Sleep(waitTime)

		// check update project id done
		timeoutMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())
		checkErr = service.CheckDBInstanceStatus(ctx, instanceId, timeoutMinutes)
		if checkErr != nil {
			return checkErr
		}

	}

	if d.HasChange("security_group_ids") {
		ids := d.Get("security_group_ids").([]interface{})
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}

		err := service.ModifyDBInstanceSecurityGroupsById(ctx, instanceId, sgIds)
		if err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") {
		request := postgres.NewModifyDBInstanceDeletionProtectionRequest()
		if v, ok := d.GetOkExists("deletion_protection"); ok {
			request.DeletionProtection = helper.Bool(v.(bool))
		}

		request.DBInstanceId = &instanceId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceDeletionProtection(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify postgresql db instance deletion protection failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d, meta)
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
		service    = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if err := service.IsolatePostgresqlReadonlyInstanceV2(ctx, instanceId); err != nil {
		return err
	}

	// wait for isolating
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		instance, e := service.DescribePostgresqlReadonlyInstanceV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instance == nil {
			return nil
		}

		if instance.DBInstanceStatus != nil && *instance.DBInstanceStatus == POSTGRESQL_STAUTS_ISOLATED {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("waiting for postgresql_readonly_instance_v2 %s isolating", instanceId))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete postgresql_readonly_instance_v2 id=%s fail, reason:%+v", logId, instanceId, err)
		return err
	}

	// delete
	if err := service.DeletePostgresqlInstanceV2(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
