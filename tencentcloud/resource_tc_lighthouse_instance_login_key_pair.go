/*
Provides a resource to create a lighthouse instance_login_key_pair

Example Usage

```hcl
resource "tencentcloud_lighthouse_instance_login_key_pair" "instance_login_key_pair" {
  instance_ids =
  permit_login = ""
}
```

Import

lighthouse instance_login_key_pair can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_instance_login_key_pair.instance_login_key_pair instance_login_key_pair_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"log"
)

func resourceTencentCloudLighthouseInstanceLoginKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceLoginKeyPairCreate,
		Read:   resourceTencentCloudLighthouseInstanceLoginKeyPairRead,
		Update: resourceTencentCloudLighthouseInstanceLoginKeyPairUpdate,
		Delete: resourceTencentCloudLighthouseInstanceLoginKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list. Each request can contain up to 100 instances at a time.",
			},

			"permit_login": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether to allow login with the default key pair. Valid values are (YES, NO).",
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceLoginKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_login_key_pair.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudLighthouseInstanceLoginKeyPairUpdate(d, meta)
}

func resourceTencentCloudLighthouseInstanceLoginKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_login_key_pair.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceLoginKeyPairId := d.Id()

	instanceLoginKeyPair, err := service.DescribeLighthouseInstanceLoginKeyPairById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceLoginKeyPair == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseInstanceLoginKeyPair` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceLoginKeyPair.InstanceIds != nil {
		_ = d.Set("instance_ids", instanceLoginKeyPair.InstanceIds)
	}

	if instanceLoginKeyPair.PermitLogin != nil {
		_ = d.Set("permit_login", instanceLoginKeyPair.PermitLogin)
	}

	return nil
}

func resourceTencentCloudLighthouseInstanceLoginKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_login_key_pair.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyInstancesLoginKeyPairAttributeRequest()

	instanceLoginKeyPairId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_ids", "permit_login"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesLoginKeyPairAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse instanceLoginKeyPair failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseInstanceLoginKeyPairRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceLoginKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_login_key_pair.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
