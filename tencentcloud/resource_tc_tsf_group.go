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

func resourceTencentCloudTsfGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfGroupCreate,
		Read:   resourceTencentCloudTsfGroupRead,
		Update: resourceTencentCloudTsfGroupUpdate,
		Delete: resourceTencentCloudTsfGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The application ID to which the group belongs.",
			},

			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the namespace to which the group belongs.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group name field, length 1~60, beginning with a letter or underscore, can contain alphanumeric underscore.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"group_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group description.",
			},

			"alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Deployment Group Notes.",
			},

			"group_resource_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deployment Group Resource Type.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateGroupRequest()
		response = tsf.NewCreateGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("application_id"); ok {
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_desc"); ok {
		request.GroupDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf group failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.Result
	d.SetId(groupId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:group/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfGroupRead(d, meta)
}

func resourceTencentCloudTsfGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()

	group, err := service.DescribeTsfGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if group.ApplicationId != nil {
		_ = d.Set("application_id", group.ApplicationId)
	}

	if group.NamespaceId != nil {
		_ = d.Set("namespace_id", group.NamespaceId)
	}

	if group.GroupName != nil {
		_ = d.Set("group_name", group.GroupName)
	}

	if group.ClusterId != nil {
		_ = d.Set("cluster_id", group.ClusterId)
	}

	if group.GroupDesc != nil {
		_ = d.Set("group_desc", group.GroupDesc)
	}

	if group.GroupResourceType != nil {
		_ = d.Set("group_resource_type", group.GroupResourceType)
	}

	if group.Alias != nil {
		_ = d.Set("alias", group.Alias)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "group", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyGroupRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"application_id", "namespace_id", "cluster_id", "group_resource_type", "alias"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("group_name") {
		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("group_desc") {
		if v, ok := d.GetOk("group_desc"); ok {
			request.GroupDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("alias") {
		if v, ok := d.GetOk("alias"); ok {
			request.Alias = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf group failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "group", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfGroupRead(d, meta)
}

func resourceTencentCloudTsfGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	groupId := d.Id()

	if err := service.DeleteTsfGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
