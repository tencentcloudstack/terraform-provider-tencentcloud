/*
Provides a resource to create a mysql backup_encryption_status

Example Usage

Enable encryption

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_backup_encryption_status" "example" {
  instance_id       = tencentcloud_mysql_instance.example.id
  encryption_status = "on"
}
```

Disable encryption

```hcl
resource "tencentcloud_mysql_backup_encryption_status" "example" {
  instance_id       = tencentcloud_mysql_instance.example.id
  encryption_status = "off"
}
```

Import

mysql backup_encryption_status can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_backup_encryption_status.backup_encryption_status backup_encryption_status_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlBackupEncryptionStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupEncryptionStatusCreate,
		Read:   resourceTencentCloudMysqlBackupEncryptionStatusRead,
		Update: resourceTencentCloudMysqlBackupEncryptionStatusUpdate,
		Delete: resourceTencentCloudMysqlBackupEncryptionStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-XXXX. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"encryption_status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether physical backup encryption is enabled for the instance. Possible values are `on`, `off`.",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupEncryptionStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_encryption_status.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlBackupEncryptionStatusUpdate(d, meta)
}

func resourceTencentCloudMysqlBackupEncryptionStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_encryption_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	backupEncryptionStatus, err := service.DescribeMysqlBackupEncryptionStatusById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backupEncryptionStatus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlBackupEncryptionStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if backupEncryptionStatus.EncryptionStatus != nil {
		_ = d.Set("encryption_status", backupEncryptionStatus.EncryptionStatus)
	}

	return nil
}

func resourceTencentCloudMysqlBackupEncryptionStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_encryption_status.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := mysql.NewModifyBackupEncryptionStatusRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("encryption_status"); ok {
		request.EncryptionStatus = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyBackupEncryptionStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql backupEncryptionStatus failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlBackupEncryptionStatusRead(d, meta)
}

func resourceTencentCloudMysqlBackupEncryptionStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_encryption_status.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
