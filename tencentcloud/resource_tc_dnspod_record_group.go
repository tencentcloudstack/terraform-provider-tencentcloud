/*
Provides a resource to create a dnspod record_group

Example Usage

```hcl
resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "dnspod.cn"
  group_name = "group_name_demo"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

dnspod record_group can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_record_group.record_group record_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDnspodRecordGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodRecordGroupCreate,
		Read:   resourceTencentCloudDnspodRecordGroupRead,
		Update: resourceTencentCloudDnspodRecordGroupUpdate,
		Delete: resourceTencentCloudDnspodRecordGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Record Group Name.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudDnspodRecordGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dnspod.NewCreateRecordGroupRequest()
		response = dnspod.NewCreateRecordGroupResponse()
		groupId uint64
		domain string
	)
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().CreateRecordGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod record_group failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(strings.Join([]string{domain, helper.UInt64ToStr(groupId)}, FILED_SP))
	// d.SetId(helper.UInt64ToStr(groupId))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::dnspod:%s:uin/:domainId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudDnspodRecordGroupRead(d, meta)
}

func resourceTencentCloudDnspodRecordGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_record_group id is broken, id is %s", d.Id())
	}
	domain := idSplit[0]
	groupId := helper.StrToUInt64(idSplit[1])

	recordGroup, err := service.DescribeDnspodRecordGroupById(ctx, domain, groupId)
	if err != nil {
		return err
	}

	if recordGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodRecordGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if recordGroup.GroupType != nil {
		_ = d.Set("group_type", recordGroup.GroupType)
	}

	if recordGroup.GroupName != nil {
		_ = d.Set("group_name", recordGroup.GroupName)
	}

	if recordGroup.GroupId != nil {
		_ = d.Set("group_id", recordGroup.GroupId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "dnspod", "domainId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudDnspodRecordGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dnspod.NewModifyRecordGroupRequest()

	groupId := d.Id()

	// request.GroupId = &groupId
	request.GroupId = helper.StrToUint64Point(groupId)

	immutableArgs := []string{"domain", "group_name", "domain_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("domain") {
		if v, ok := d.GetOk("domain"); ok {
			request.Domain = helper.String(v.(string))
		}
	}

	if d.HasChange("group_name") {
		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("domain_id") {
		if v, ok := d.GetOkExists("domain_id"); ok {
			request.DomainId = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifyRecordGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod record_group failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("dnspod", "domainId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudDnspodRecordGroupRead(d, meta)
}

func resourceTencentCloudDnspodRecordGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}
	groupId := d.Id()

	if err := service.DeleteDnspodRecordGroupById(ctx, helper.StrToUint64Point(groupId)); err != nil {
		return err
	}

	return nil
}
