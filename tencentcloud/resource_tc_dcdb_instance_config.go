/*
Provides a resource to create a dcdb instance_config

Example Usage

```hcl
resource "tencentcloud_dcdb_instance_config" "instance_config" {
  instance_id = ""
  rs_access_strategy =
}
```

Import

dcdb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_instance_config.instance_config instance_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudDcdbInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbInstanceConfigCreate,
		Read:   resourceTencentCloudDcdbInstanceConfigRead,
		Update: resourceTencentCloudDcdbInstanceConfigUpdate,
		Delete: resourceTencentCloudDcdbInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"rs_access_strategy": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "RS nearest access mode, 0-no policy, 1-nearest access.",
			},
		},
	}
}

func resourceTencentCloudDcdbInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_instance_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	d.SetId(instanceId)

	return resourceTencentCloudDcdbInstanceConfigUpdate(d, meta)
}

func resourceTencentCloudDcdbInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_instance_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	detail, err := service.DescribeDcdbDbInstanceDetailById(ctx, instanceId)
	if err != nil {
		return err
	}

	if detail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbInstanceConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if detail.RsAccessStrategy != nil {
		_ = d.Set("rs_access_strategy", detail.RsAccessStrategy)
	}

	return nil
}

func resourceTencentCloudDcdbInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_instance_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		service    = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		rsStrategy int
	)

	instanceId := d.Id()
	if d.HasChange("rs_access_strategy") {
		if v, ok := d.GetOk("rs_access_strategy"); ok {
			rsStrategy = v.(int)
		}
	}

	err := service.SetRealServerAccessStrategy(ctx, instanceId, rsStrategy)

	if err != nil {
		log.Printf("[CRITAL]%s update dcdb instanceConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbInstanceConfigRead(d, meta)
}

func resourceTencentCloudDcdbInstanceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_instance_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
