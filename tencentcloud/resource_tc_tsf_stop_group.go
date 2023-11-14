/*
Provides a resource to create a tsf stop_group

Example Usage

```hcl
resource "tencentcloud_tsf_stop_group" "stop_group" {
  group_id = ""
}
```

Import

tsf stop_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_stop_group.stop_group stop_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"log"
)

func resourceTencentCloudTsfStopGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfStopGroupCreate,
		Read:   resourceTencentCloudTsfStopGroupRead,
		Update: resourceTencentCloudTsfStopGroupUpdate,
		Delete: resourceTencentCloudTsfStopGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "GroupId.",
			},
		},
	}
}

func resourceTencentCloudTsfStopGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_stop_group.create")()
	defer inconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTsfStopGroupUpdate(d, meta)
}

func resourceTencentCloudTsfStopGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_stop_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	stopGroupId := d.Id()

	stopGroup, err := service.DescribeTsfStopGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if stopGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfStopGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if stopGroup.GroupId != nil {
		_ = d.Set("group_id", stopGroup.GroupId)
	}

	return nil
}

func resourceTencentCloudTsfStopGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_stop_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		startGroupRequest  = tsf.NewStartGroupRequest()
		startGroupResponse = tsf.NewStartGroupResponse()
	)

	stopGroupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().StartGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf stopGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfStopGroupRead(d, meta)
}

func resourceTencentCloudTsfStopGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_stop_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
