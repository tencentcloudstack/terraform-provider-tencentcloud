/*
Provides a resource to create a cdb local_binlog_config

Example Usage

```hcl
resource "tencentcloud_cdb_local_binlog_config" "local_binlog_config" {
  instance_id = ""
  save_hours =
  max_usage =
}
```

Import

cdb local_binlog_config can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_local_binlog_config.local_binlog_config local_binlog_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
)

func resourceTencentCloudCdbLocalBinlogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbLocalBinlogConfigCreate,
		Read:   resourceTencentCloudCdbLocalBinlogConfigRead,
		Update: resourceTencentCloudCdbLocalBinlogConfigUpdate,
		Delete: resourceTencentCloudCdbLocalBinlogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.",
			},

			"save_hours": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention period of local binlog. Valid range: 72-168 hours. When there is disaster recovery instance, the valid range will be 120-168 hours.",
			},

			"max_usage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Space utilization of local binlog. Value range: [30,50].",
			},
		},
	}
}

func resourceTencentCloudCdbLocalBinlogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_local_binlog_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbLocalBinlogConfigUpdate(d, meta)
}

func resourceTencentCloudCdbLocalBinlogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_local_binlog_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	localBinlogConfigId := d.Id()

	localBinlogConfig, err := service.DescribeCdbLocalBinlogConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if localBinlogConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbLocalBinlogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if localBinlogConfig.InstanceId != nil {
		_ = d.Set("instance_id", localBinlogConfig.InstanceId)
	}

	if localBinlogConfig.SaveHours != nil {
		_ = d.Set("save_hours", localBinlogConfig.SaveHours)
	}

	if localBinlogConfig.MaxUsage != nil {
		_ = d.Set("max_usage", localBinlogConfig.MaxUsage)
	}

	return nil
}

func resourceTencentCloudCdbLocalBinlogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_local_binlog_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyLocalBinlogConfigRequest()

	localBinlogConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "save_hours", "max_usage"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyLocalBinlogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb localBinlogConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbLocalBinlogConfigRead(d, meta)
}

func resourceTencentCloudCdbLocalBinlogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_local_binlog_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
