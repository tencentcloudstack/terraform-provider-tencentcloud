/*
Provides a resource to create a dnspod record_group

Example Usage

```hcl
resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "dnspod.cn"
  group_name = "group_demo"
}
```

Import

dnspod record_group can be imported using the domain#groupId, e.g.

```
terraform import tencentcloud_dnspod_record_group.record_group domain#groupId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

			"group_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Group ID.",
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
		groupId  uint64
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	// if v, ok := d.GetOkExists("domain_id"); ok {
	// 	request.DomainId = helper.IntUint64(v.(int))
	// }

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

	_ = d.Set("domain", domain)

	if recordGroup.GroupName != nil {
		_ = d.Set("group_name", recordGroup.GroupName)
	}

	if recordGroup.GroupId != nil {
		_ = d.Set("group_id", recordGroup.GroupId)
	}

	// if recordGroup.GroupType != nil {
	// 	_ = d.Set("group_type", recordGroup.GroupType)
	// }

	return nil
}

func resourceTencentCloudDnspodRecordGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dnspod.NewModifyRecordGroupRequest()
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_record_group id is broken, id is %s", d.Id())
	}
	request.Domain = helper.String(idSplit[0])
	request.GroupId = helper.StrToUint64Point(idSplit[1])

	immutableArgs := []string{"domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
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

	return resourceTencentCloudDnspodRecordGroupRead(d, meta)
}

func resourceTencentCloudDnspodRecordGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record_group.delete")()
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

	if err := service.DeleteDnspodRecordGroupById(ctx, domain, groupId); err != nil {
		return err
	}

	return nil
}
