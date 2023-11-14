/*
Provides a resource to create a tsf release_api_group

Example Usage

```hcl
resource "tencentcloud_tsf_release_api_group" "release_api_group" {
  group_id = "grp-qp0rj3zi"
}
```

Import

tsf release_api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_release_api_group.release_api_group release_api_group_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfReleaseApiGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfReleaseApiGroupCreate,
		Read:   resourceTencentCloudTsfReleaseApiGroupRead,
		Delete: resourceTencentCloudTsfReleaseApiGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Api group Id .",
			},
		},
	}
}

func resourceTencentCloudTsfReleaseApiGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_release_api_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewReleaseApiGroupRequest()
		response = tsf.NewReleaseApiGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ReleaseApiGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf releaseApiGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(groupId)

	return resourceTencentCloudTsfReleaseApiGroupRead(d, meta)
}

func resourceTencentCloudTsfReleaseApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_release_api_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	releaseApiGroupId := d.Id()

	releaseApiGroup, err := service.DescribeTsfReleaseApiGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if releaseApiGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfReleaseApiGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if releaseApiGroup.GroupId != nil {
		_ = d.Set("group_id", releaseApiGroup.GroupId)
	}

	return nil
}

func resourceTencentCloudTsfReleaseApiGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_release_api_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	releaseApiGroupId := d.Id()

	if err := service.DeleteTsfReleaseApiGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
