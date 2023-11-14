/*
Provides a resource to create a tsf start_container_group

Example Usage

```hcl
resource "tencentcloud_tsf_start_container_group" "start_container_group" {
  group_id = "group-xxxxxxxx"
}
```

Import

tsf start_container_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_start_container_group.start_container_group start_container_group_id
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

func resourceTencentCloudTsfStartContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfStartContainerGroupCreate,
		Read:   resourceTencentCloudTsfStartContainerGroupRead,
		Update: resourceTencentCloudTsfStartContainerGroupUpdate,
		Delete: resourceTencentCloudTsfStartContainerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group Id.",
			},
		},
	}
}

func resourceTencentCloudTsfStartContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_start_container_group.create")()
	defer inconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTsfStartContainerGroupUpdate(d, meta)
}

func resourceTencentCloudTsfStartContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_start_container_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	startContainerGroupId := d.Id()

	startContainerGroup, err := service.DescribeTsfStartContainerGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if startContainerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfStartContainerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if startContainerGroup.GroupId != nil {
		_ = d.Set("group_id", startContainerGroup.GroupId)
	}

	return nil
}

func resourceTencentCloudTsfStartContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_start_container_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		startContainerGroupRequest  = tsf.NewStartContainerGroupRequest()
		startContainerGroupResponse = tsf.NewStartContainerGroupResponse()
	)

	startContainerGroupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().StartContainerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf startContainerGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfStartContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfStartContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_start_container_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
