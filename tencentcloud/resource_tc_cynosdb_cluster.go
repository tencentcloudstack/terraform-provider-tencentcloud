/*
Provide a resource to create a CynosDB cluster.

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = "ap-guangzhou-4"
  vpc_id                       = "vpc-h70b6b49"
  subnet_id                    = "subnet-q6fhy1mi"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2

  tags = {
    test = "test"
  }

  force_delete = false

  rw_group_sg = [
    "sg-ibyjkl6r",
  ]
  ro_group_sg = [
    "sg-ibyjkl6r",
  ]
}
```

Import

CynosDB cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_cluster.foo cynosdbmysql-dzj5l8gz
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCynosdbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterCreate,
		Read:   resourceTencentCloudCynosdbClusterRead,
		Update: resourceTencentCloudCynosdbClusterUpdate,
		Delete: resourceTencentCloudCynosdbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: TencentCynosdbClusterBaseInfo(),
	}
}

func resourceTencentCloudCynosdbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster.create")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		client         = meta.(*TencentCloudClient).apiV3Conn
		cynosdbService = CynosdbService{client: client}
		tagService     = TagService{client: client}
		region         = client.Region

		request = cynosdb.NewCreateClustersRequest()
	)

	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.Zone = helper.String(d.Get("available_zone").(string))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	request.Port = helper.IntInt64(d.Get("port").(int))
	request.DbType = helper.String(d.Get("db_type").(string))
	request.DbVersion = helper.String(d.Get("db_version").(string))
	request.StorageLimit = helper.IntInt64(d.Get("storage_limit").(int))
	request.ClusterName = helper.String(d.Get("cluster_name").(string))
	request.AdminPassword = helper.String(d.Get("password").(string))
	request.RollbackStrategy = helper.String("noneRollback")

	// instance info
	request.Cpu = helper.IntInt64(d.Get("instance_cpu_core").(int))
	request.Memory = helper.IntInt64(d.Get("instance_memory_size").(int))

	var chargeType int64 = 0
	if v, ok := d.GetOk("charge_type"); ok {
		if v == CYNOSDB_CHARGE_TYPE_PREPAID {
			chargeType = 1
			if vv, ok := d.GetOk("prepaid_period"); ok {
				request.TimeSpan = helper.IntInt64(vv.(int))
				request.TimeUnit = helper.String("m")
			} else {
				return fmt.Errorf("prepaid period can not be empty when charge type is %s", CYNOSDB_CHARGE_TYPE_PREPAID)
			}
			if vv, ok := d.GetOk("auto_renew_flag"); ok {
				request.AutoRenewFlag = helper.IntInt64(vv.(int))
			}
		}
	}
	request.PayMode = &chargeType

	request.InstanceCount = helper.Int64(1)

	var response *cynosdb.CreateClustersResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateClusters(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && len(response.Response.ClusterIds) != 1 {
		return fmt.Errorf("cynosdb cluster id count isn't 1")
	}
	d.SetId(*response.Response.ClusterIds[0])
	id := d.Id()

	_, has, err := cynosdbService.DescribeClusterById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s creating cynosdb cluster failed, cluster doesn't exist", logId)
	}

	// set tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("cynosdb", "cluster", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// set maintenance info
	cluster, _, err := cynosdbService.DescribeClusterById(ctx, id)
	if err != nil {
		return err
	}
	var rwInstanceId string
	for _, v := range cluster.InstanceSet {
		_, instance, has, err := cynosdbService.DescribeInstanceById(ctx, *v.InstanceId)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		if *instance.InstanceType == CYNOSDB_INSTANCE_RW_TYPE {
			rwInstanceId = *instance.InstanceId
			break
		}
	}
	// set maintenance info
	var weekdays []interface{}
	if v, ok := d.GetOk("instance_maintain_weekdays"); ok {
		weekdays = v.(*schema.Set).List()
	} else {
		weekdays = []interface{}{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	reqWeekdays := make([]*string, 0, len(weekdays))
	for _, v := range weekdays {
		reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
	}
	startTime := int64(d.Get("instance_maintain_start_time").(int))
	duration := int64(d.Get("instance_maintain_duration").(int))
	err = cynosdbService.ModifyMaintainPeriodConfig(ctx, rwInstanceId, startTime, duration, reqWeekdays)
	if err != nil {
		return err
	}

	// set sg
	insGrps, err := cynosdbService.DescribeClusterInstanceGrps(ctx, id)
	if err != nil {
		return err
	}
	var rwGroupId string
	var roGroupId string
	for _, insGrp := range insGrps.Response.InstanceGrpInfoList {
		if *insGrp.Type == CYNOSDB_INSGRP_HA {
			rwGroupId = *insGrp.InstanceGrpId
		} else if *insGrp.Type == CYNOSDB_INSGRP_RO {
			roGroupId = *insGrp.InstanceGrpId
		}
	}
	if v, ok := d.GetOk("rw_group_sg"); ok {
		vv := v.([]interface{})
		vvv := make([]*string, 0, len(vv))
		for _, item := range vv {
			vvv = append(vvv, helper.String(item.(string)))
		}
		if err = cynosdbService.ModifyInsGrpSecurityGroups(ctx, rwGroupId, d.Get("available_zone").(string), vvv); err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("ro_group_sg"); ok {
		vv := v.([]interface{})
		vvv := make([]*string, 0, len(vv))
		for _, item := range vv {
			vvv = append(vvv, helper.String(item.(string)))
		}
		if err = cynosdbService.ModifyInsGrpSecurityGroups(ctx, roGroupId, d.Get("available_zone").(string), vvv); err != nil {
			return err
		}
	}

	return resourceTencentCloudCynosdbClusterRead(d, meta)
}

func resourceTencentCloudCynosdbClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	cynosdbService := CynosdbService{client: client}
	cluster, has, err := cynosdbService.DescribeClusterById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", cluster.ProjectID)
	_ = d.Set("available_zone", cluster.Zone)
	_ = d.Set("vpc_id", cluster.VpcId)
	_ = d.Set("subnet_id", cluster.SubnetId)
	_ = d.Set("port", cluster.Vport)
	_ = d.Set("db_type", cluster.DbType)
	_ = d.Set("db_version", cluster.DbVersion)
	_ = d.Set("cluster_name", cluster.ClusterName)
	_ = d.Set("charge_type", CYNOSDB_CHARGE_TYPE[*cluster.PayMode])
	_ = d.Set("charset", cluster.Charset)
	_ = d.Set("cluster_status", cluster.Status)
	_ = d.Set("create_time", cluster.CreateTime)
	_ = d.Set("storage_used", *cluster.UsedStorage/1000/1000)

	//tag
	tagService := &TagService{client: client}
	tags, err := tagService.DescribeResourceTags(ctx, "cynosdb", "cluster", client.Region, id)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	for _, v := range cluster.InstanceSet {
		_, instance, has, err := cynosdbService.DescribeInstanceById(ctx, *v.InstanceId)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		if *instance.InstanceType == CYNOSDB_INSTANCE_RW_TYPE {
			maintain, err := cynosdbService.DescribeMaintainPeriod(ctx, *v.InstanceId)
			if err != nil {
				return err
			}
			_ = d.Set("instance_cpu_core", instance.Cpu)
			_ = d.Set("instance_memory_size", instance.Memory)
			_ = d.Set("instance_id", instance.InstanceId)
			_ = d.Set("instance_name", instance.InstanceName)
			_ = d.Set("instance_status", instance.Status)
			_ = d.Set("instance_storage_size", instance.Storage)
			_ = d.Set("instance_maintain_weekdays", maintain.Response.MaintainWeekDays)
			_ = d.Set("instance_maintain_start_time", maintain.Response.MaintainStartTime)
			_ = d.Set("instance_maintain_duration", maintain.Response.MaintainDuration)
			break
		}
	}

	// instance group infos
	insGrps, err := cynosdbService.DescribeClusterInstanceGrps(ctx, id)
	if err != nil {
		return err
	}
	var rwGroupId string
	rwGroupIns := make([]map[string]interface{}, 0)
	rwGroupAddr := make([]map[string]interface{}, 0)
	var roGroupId string
	roGroupIns := make([]map[string]interface{}, 0)
	roGroupAddr := make([]map[string]interface{}, 0)
	for _, insGrp := range insGrps.Response.InstanceGrpInfoList {
		if *insGrp.Type == CYNOSDB_INSGRP_HA {
			rwGroupId = *insGrp.InstanceGrpId
			_ = d.Set("rw_group_id", rwGroupId)
			for _, rwIns := range insGrp.InstanceSet {
				rwGroupIns = append(rwGroupIns, map[string]interface{}{
					"instance_id":   *rwIns.InstanceId,
					"instance_name": *rwIns.InstanceName,
				})
			}
			rwGroupAddr = append(rwGroupAddr, map[string]interface{}{
				"ip":   *insGrp.Vip,
				"port": *insGrp.Vport,
			})
		} else if *insGrp.Type == CYNOSDB_INSGRP_RO {
			roGroupId = *insGrp.InstanceGrpId
			_ = d.Set("ro_group_id", roGroupId)
			for _, roIns := range insGrp.InstanceSet {
				roGroupIns = append(roGroupIns, map[string]interface{}{
					"instance_id":   *roIns.InstanceId,
					"instance_name": *roIns.InstanceName,
				})
			}
			roGroupAddr = append(roGroupAddr, map[string]interface{}{
				"ip":   *insGrp.Vip,
				"port": *insGrp.Vport,
			})
		}
	}
	_ = d.Set("rw_group_instances", rwGroupIns)
	_ = d.Set("rw_group_addr", rwGroupAddr)
	_ = d.Set("ro_group_instances", roGroupIns)
	_ = d.Set("ro_group_addr", roGroupAddr)

	// sg infos
	if rwGroupId != "" {
		sgs, err := cynosdbService.DescribeInsGrpSecurityGroups(ctx, rwGroupId)
		if err != nil {
			return err
		}
		if sgs != nil {
			sgIds := make([]*string, 0, len(sgs.Response.Groups))
			for _, item := range sgs.Response.Groups {
				sgIds = append(sgIds, item.SecurityGroupId)
			}
			_ = d.Set("rw_group_sg", sgIds)
		}
	}
	if roGroupId != "" {
		sgs, err := cynosdbService.DescribeInsGrpSecurityGroups(ctx, roGroupId)
		if err != nil {
			return err
		}
		if sgs != nil {
			sgIds := make([]*string, 0, len(sgs.Response.Groups))
			for _, item := range sgs.Response.Groups {
				sgIds = append(sgIds, item.SecurityGroupId)
			}
			_ = d.Set("ro_group_sg", sgIds)
		}
	}

	return nil
}

func resourceTencentCloudCynosdbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster.update")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId      = d.Id()
		instanceId     = d.Get("instance_id").(string)
		client         = meta.(*TencentCloudClient).apiV3Conn
		cynosdbService = CynosdbService{client: client}
		tagService     = TagService{client: client}
		region         = client.Region
	)

	// check unsupported field modification
	unsupported := []string{
		"project_id", "vpc_id", "subnet_id", "port", "storage_limit", "cluster_name",
		"password", "prepaid_period", "auto_renew_flag",
	}
	for _, v := range unsupported {
		if d.HasChange(v) {
			return fmt.Errorf("[CRITAL] field %s is not allowed to be modified", v)
		}
	}

	d.Partial(true)

	if d.HasChange("instance_cpu_core") || d.HasChange("instance_memory_size") {
		cpu := int64(d.Get("instance_cpu_core").(int))
		memory := int64(d.Get("instance_memory_size").(int))
		err := cynosdbService.UpgradeInstance(ctx, instanceId, cpu, memory)
		if err != nil {
			return err
		}

		errUpdate := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, infos, has, e := cynosdbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s updating cynosdb cluster instance failed, instance doesn't exist", logId))
			}

			cpuReal := *infos.Cpu
			memReal := *infos.Memory
			if cpu != cpuReal || memory != memReal {
				return resource.RetryableError(fmt.Errorf("[CRITAL] updating cynosdb instance, current cpu and memory values: %d, %d, waiting for them becoming new value: %d, %d", cpuReal, memReal, cpu, memory))
			}
			return nil
		})
		if errUpdate != nil {
			return errUpdate
		}

		d.SetPartial("instance_cpu_core")
		d.SetPartial("instance_memory_size")
	}

	if d.HasChange("instance_maintain_weekdays") || d.HasChange("instance_maintain_start_time") || d.HasChange("instance_maintain_duration") {
		weekdays := d.Get("instance_maintain_weekdays").(*schema.Set).List()
		reqWeekdays := make([]*string, 0, len(weekdays))
		for _, v := range weekdays {
			reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
		}
		startTime := int64(d.Get("instance_maintain_start_time").(int))
		duration := int64(d.Get("instance_maintain_duration").(int))
		err := cynosdbService.ModifyMaintainPeriodConfig(ctx, instanceId, startTime, duration, reqWeekdays)
		if err != nil {
			return err
		}

		d.SetPartial("instance_maintain_weekdays")
		d.SetPartial("instance_maintain_start_time")
		d.SetPartial("instance_maintain_duration")
	}

	// update tags
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("cynosdb", "cluster", region, clusterId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	// update sg
	if d.HasChange("rw_group_sg") {
		v := d.Get("rw_group_sg").([]interface{})
		vv := make([]*string, 0, len(v))
		for _, item := range v {
			vv = append(vv, helper.String(item.(string)))
		}
		if err := cynosdbService.ModifyInsGrpSecurityGroups(ctx, d.Get("rw_group_id").(string), d.Get("available_zone").(string), vv); err != nil {
			return err
		}
		d.SetPartial("rw_group_sg")
	}
	if d.HasChange("ro_group_sg") {
		v := d.Get("ro_group_sg").([]interface{})
		vv := make([]*string, 0, len(v))
		for _, item := range v {
			vv = append(vv, helper.String(item.(string)))
		}
		if err := cynosdbService.ModifyInsGrpSecurityGroups(ctx, d.Get("ro_group_id").(string), d.Get("available_zone").(string), vv); err != nil {
			return err
		}
		d.SetPartial("ro_group_sg")
	}

	d.Partial(false)

	return resourceTencentCloudCynosdbClusterRead(d, meta)
}

func resourceTencentCloudCynosdbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clusterID := d.Id()
	cynosdbService := CynosdbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	forceDelete := d.Get("force_delete").(bool)
	var err error
	if err = cynosdbService.IsolateCluster(ctx, clusterID); err != nil {
		return err
	}

	if forceDelete {
		if err = cynosdbService.OfflineCluster(ctx, clusterID); err != nil {
			return err
		}
	}

	return nil
}
