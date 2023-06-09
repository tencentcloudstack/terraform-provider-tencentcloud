/*
Provides a resource to create a mariadb hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "10.0"
  instance_name = "db-test-del"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  vip           = "10.0.0.197"
  zones         = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]
  tags          = {
    createdBy   = "terraform"
  }
}

```
Import

mariadb hour_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_hour_db_instance.hour_db_instance tdsql-kjqih9nn
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbHourDbInstanceCreate,
		Read:   resourceTencentCloudMariadbHourDbInstanceRead,
		Update: resourceTencentCloudMariadbHourDbInstanceUpdate,
		Delete: resourceTencentCloudMariadbHourDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "available zone of instance.",
			},

			"node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "number of node for instance.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance memory.",
			},

			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance disk storage.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "project id.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vpc id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id, it&amp;#39;s required when vpcId is set.",
			},

			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vip.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "db engine version, default to 10.1.9.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of this instance.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudMariadbHourDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewCreateHourDBInstanceRequest()
		response   *mariadb.CreateHourDBInstanceResponse
		instanceId string
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOk("node_count"); ok {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateHourDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb hourDbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	initParams := []*mariadb.DBParamValue{
		{
			Param: helper.String("character_set_server"),
			Value: helper.String("utf8mb4"),
		},
		{
			Param: helper.String("lower_case_table_names"),
			Value: helper.String("1"),
		},
		{
			Param: helper.String("sync_mode"),
			Value: helper.String("2"),
		},
		{
			Param: helper.String("innodb_page_size"),
			Value: helper.String("16384"),
		},
	}

	initRet, err := service.InitDbInstance(ctx, instanceId, initParams)
	if err != nil {
		return err
	}

	if !initRet {
		return fmt.Errorf("db instance init failed")
	}

	// set Tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:mariadb-hour-instance/%s", region, instanceId)
		if err = tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbHourDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbHourDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	hourDbInstance, err := service.DescribeMariadbDbInstanceDetail(ctx, instanceId)
	if err != nil {
		return err
	}

	if hourDbInstance == nil {
		d.SetId("")
		return fmt.Errorf("resource `hourDbInstance` %s does not exist", instanceId)
	}

	if hourDbInstance.Zone != nil {
		var zones []*string
		zones = append(zones, hourDbInstance.Zone)
		if hourDbInstance.SlaveZones != nil {
			zones = append(zones, hourDbInstance.SlaveZones...)
		}
		_ = d.Set("zones", zones)
	}

	if hourDbInstance.NodeCount != nil {
		_ = d.Set("node_count", hourDbInstance.NodeCount)
	}

	if hourDbInstance.Memory != nil {
		_ = d.Set("memory", hourDbInstance.Memory)
	}

	if hourDbInstance.Storage != nil {
		_ = d.Set("storage", hourDbInstance.Storage)
	}

	if hourDbInstance.ProjectId != nil {
		_ = d.Set("project_id", hourDbInstance.ProjectId)
	}

	if hourDbInstance.Vip != nil {
		_ = d.Set("vip", hourDbInstance.Vip)
	}

	if hourDbInstance.VpcId != nil {
		_ = d.Set("vpc_id", hourDbInstance.VpcId)
	}

	if hourDbInstance.SubnetId != nil {
		_ = d.Set("subnet_id", hourDbInstance.SubnetId)
	}

	if hourDbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", hourDbInstance.DbVersionId)
	}

	if hourDbInstance.InstanceName != nil {
		_ = d.Set("instance_name", hourDbInstance.InstanceName)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "mariadb-hour-instance", tcClient.Region, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMariadbHourDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewModifyDBInstanceNameRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	if d.HasChange("zones") {
		return fmt.Errorf("`zones` do not support change now.")
	}

	if d.HasChange("node_count") {
		return fmt.Errorf("`node_count` do not support change now.")
	}

	if d.HasChange("memory") {
		return fmt.Errorf("`memory` do not support change now.")
	}

	if d.HasChange("storage") {
		return fmt.Errorf("`storage` do not support change now.")
	}

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}

	if d.HasChange("db_version_id") {
		return fmt.Errorf("`db_version_id` do not support change now.")
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstanceName(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create mariadb hourDbInstance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("mariadb", "mariadb-hour-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			projectId := int64(v.(int))
			MPRequest := mariadb.NewModifyDBInstancesProjectRequest()
			MPRequest.InstanceIds = common.StringPtrs([]string{instanceId})
			MPRequest.ProjectId = &projectId

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstancesProject(MPRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb modifyInstanceProject failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("vip") {
		if v, ok := d.GetOk("vip"); ok {
			Vip := v.(string)
			var VipFlowId int64
			VipRequest := mariadb.NewModifyInstanceNetworkRequest()
			VipRequest.InstanceId = &instanceId
			VipRequest.Vip = &Vip
			if v, ok := d.GetOk("vpc_id"); ok {
				VipRequest.VpcId = helper.String(v.(string))
			}

			if v, ok := d.GetOk("subnet_id"); ok {
				VipRequest.SubnetId = helper.String(v.(string))
			}

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyInstanceNetwork(VipRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				VipFlowId = *result.Response.FlowId
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb network failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if VipFlowId != NONE_FLOW_TASK {
				err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
					result, e := service.DescribeFlowById(ctx, VipFlowId)
					if e != nil {
						return retryError(e)
					}

					if *result.Status == MARIADB_TASK_SUCCESS {
						return nil
					} else if *result.Status == MARIADB_TASK_RUNNING {
						return resource.RetryableError(fmt.Errorf("operate mariadb network status is running"))
					} else if *result.Status == MARIADB_TASK_FAIL {
						return resource.NonRetryableError(fmt.Errorf("operate mariadb network status is fail"))
					} else {
						e = fmt.Errorf("operate mariadb network status illegal")
						return resource.NonRetryableError(e)
					}
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb network task failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudMariadbHourDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbHourDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteMariadbHourDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
