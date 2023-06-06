/*
Provides a resource to create a dcdb encrypt_attributes_config

~> **NOTE:**  This resource currently only supports the newly created MySQL 8.0.24 version.

Example Usage

```hcl
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

resource "tencentcloud_dcdb_db_instance" "prepaid_instance" {
	instance_name = "test_dcdb_db_post_instance"
	zones = [var.default_az]
	period = 1
	shard_memory = "2"
	shard_storage = "10"
	shard_node_count = "2"
	shard_count = "2"
	vpc_id = local.vpc_id
	subnet_id = local.subnet_id
	db_version_id = "8.0"
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	security_group_ids = [local.sg_id]
}

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
	instance_name = "test_dcdb_db_hourdb_instance"
	zones = [var.default_az]
	shard_memory = "2"
	shard_storage = "10"
	shard_node_count = "2"
	shard_count = "2"
	vpc_id = local.vpc_id
	subnet_id = local.subnet_id
	security_group_id = local.sg_id
	db_version_id = "8.0"
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
}

locals {
	prepaid_dcdb_id = tencentcloud_dcdb_db_instance.prepaid_instance.id
	hourdb_dcdb_id = tencentcloud_dcdb_hourdb_instance.hourdb_instance.id
}

// for postpaid instance
resource "tencentcloud_dcdb_encrypt_attributes_config" "config_hourdb" {
  instance_id = local.hourdb_dcdb_id
  encrypt_enabled = 1
}

// for prepaid instance
resource "tencentcloud_dcdb_encrypt_attributes_config" "config_prepaid" {
  instance_id = local.prepaid_dcdb_id
  encrypt_enabled = 1
}
```

Import

dcdb encrypt_attributes_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_encrypt_attributes_config.encrypt_attributes_config encrypt_attributes_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbEncryptAttributesConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbEncryptAttributesConfigCreate,
		Read:   resourceTencentCloudDcdbEncryptAttributesConfigRead,
		Update: resourceTencentCloudDcdbEncryptAttributesConfigUpdate,
		Delete: resourceTencentCloudDcdbEncryptAttributesConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"encrypt_enabled": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "whether to enable data encryption. Notice: it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.",
			},
		},
	}
}

func resourceTencentCloudDcdbEncryptAttributesConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	d.SetId(instanceId)

	return resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	encryptAttributesConfig, err := service.DescribeDcdbEncryptAttributesConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if encryptAttributesConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbEncryptAttributesConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if encryptAttributesConfig.EncryptStatus != nil {
		_ = d.Set("encrypt_enabled", encryptAttributesConfig.EncryptStatus)
	}

	return nil
}

func resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBEncryptAttributesRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("encrypt_enabled") {
		if v, ok := d.GetOkExists("encrypt_enabled"); ok {
			request.EncryptEnabled = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb encryptAttributesConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbEncryptAttributesConfigRead(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
