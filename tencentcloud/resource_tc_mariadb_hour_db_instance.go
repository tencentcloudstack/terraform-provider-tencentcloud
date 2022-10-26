/*
Provides a resource to create a mariadb hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "hour_db_instance" {
  zones = ""
  node_count = ""
  memory = ""
  storage = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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
				Description: "vpc id for this instance.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id for this instance, it&amp;#39;s required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "db engine version for this instance, default to 10.1.9.",
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
			if err, ok := e.(*errors.TencentCloudSDKError); ok {
				if err.Code == "FailedOperation.PayFailed" {
					return &resource.RetryError{e, false}
				}
			}
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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

	hourDbInstances, err := service.DescribeMariadbHourDbInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if hourDbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `hourDbInstance` %s does not exist", instanceId)
	}

	if len(hourDbInstances.Instances) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `hourDbInstance` %s does not exist", instanceId)
	}
	hourDbInstance := hourDbInstances.Instances[0]

	if hourDbInstance.Zone != nil {
		_ = d.Set("zones", hourDbInstance.Zone)
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
