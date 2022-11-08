/*
Provides a resource to create a mariadb hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "hour_db_instance" {
  zones = ["ap-guangzhou-7","ap-guangzhou-7"]
  node_count = 2
  memory = 2
  storage = 10
  vpc_id = "vpc-68vi2d3h"
  subnet_id = "subnet-ob6clqwk"
  db_version_id = "8.0"
  instance_name = "db-test-2"
  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params  {
      param = "lower_case_table_names"
      value = "1"
  }
  init_params  {
      param = "sync_mode"
      value = "2"
  }
  init_params  {
      param = "innodb_page_size"
      value = "16384"
  }
}

```
Import

mariadb hour_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_hour_db_instance.hour_db_instance hourDbInstance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbHourDbInstanceRead,
		Create: resourceTencentCloudMariadbHourDbInstanceCreate,
		Update: resourceTencentCloudMariadbHourDbInstanceUpdate,
		Delete: resourceTencentCloudMariadbHourDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
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

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
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

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:mariadb-hour-instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbHourDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbHourDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	hourDbInstance, err := service.DescribeMariadbDbInstance(ctx, instanceId)

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

	if hourDbInstance.UniqueVpcId != nil {
		_ = d.Set("vpc_id", hourDbInstance.UniqueVpcId)
	}

	if hourDbInstance.UniqueSubnetId != nil {
		_ = d.Set("subnet_id", hourDbInstance.UniqueSubnetId)
	}

	if hourDbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", hourDbInstance.DbVersionId)
	}

	if hourDbInstance.InstanceName != nil {
		_ = d.Set("instance_name", hourDbInstance.InstanceName)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "mariadb-hour-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMariadbHourDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_hour_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mariadb.NewModifyDBInstanceNameRequest()

	instanceId := d.Id()

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
