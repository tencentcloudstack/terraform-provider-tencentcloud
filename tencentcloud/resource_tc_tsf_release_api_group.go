package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfReleaseApiGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfReleaseApiGroupCreate,
		Read:   resourceTencentCloudTsfReleaseApiGroupRead,
		Delete: resourceTencentCloudTsfReleaseApiGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "api group Id.",
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

	if !*response.Response.Result {
		return fmt.Errorf("[CRITAL]%s create tsf releaseApiGroup failed", logId)
	}
	d.SetId(groupId)

	return resourceTencentCloudTsfReleaseApiGroupRead(d, meta)
}

func resourceTencentCloudTsfReleaseApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_release_api_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()
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

	return nil
}
