/*
Provides a resource to create a dasb user_group_members

Example Usage

```hcl
resource "tencentcloud_dasb_user_group_members" "example" {
  user_group_id = 3
  member_id_set = [1, 2, 3]
}
```

Import

dasb user_group_members can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group_members.example 3#1,2,3
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudDasbUserGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbUserGroupMembersCreate,
		Read:   resourceTencentCloudDasbUserGroupMembersRead,
		Delete: resourceTencentCloudDasbUserGroupMembersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "User Group ID.",
			},
			"member_id_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Collection of member user IDs.",
			},
		},
	}
}

func resourceTencentCloudDasbUserGroupMembersCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group_members.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		request        = dasb.NewAddUserGroupMembersRequest()
		userGroupId    string
		memberIdSetStr string
	)

	if v, ok := d.GetOkExists("user_group_id"); ok {
		request.Id = helper.IntUint64(v.(int))
		userGroupIdInt := v.(int)
		userGroupId = strconv.Itoa(userGroupIdInt)
	}

	if v, ok := d.GetOk("member_id_set"); ok {
		memberIdSetSet := v.(*schema.Set).List()
		tmpList := make([]string, 0)
		for i := range memberIdSetSet {
			memberIdSet := memberIdSetSet[i].(int)
			request.MemberIdSet = append(request.MemberIdSet, helper.IntUint64(memberIdSet))
			tmpList = append(tmpList, strconv.Itoa(memberIdSet))
		}

		memberIdSetStr = strings.Join(tmpList, COMMA_SP)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().AddUserGroupMembers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("dasb UserGroupMembers not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb UserGroupMembers failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{userGroupId, memberIdSetStr}, FILED_SP))

	return resourceTencentCloudDasbUserGroupMembersRead(d, meta)
}

func resourceTencentCloudDasbUserGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group_members.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	userGroupId := idSplit[0]

	UserGroupMembers, err := service.DescribeDasbUserGroupMembersById(ctx, userGroupId)
	if err != nil {
		return err
	}

	if UserGroupMembers == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbUserGroupMembers` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("user_group_id", userGroupId)
	_ = d.Set("member_id_set", UserGroupMembers)

	return nil
}

func resourceTencentCloudDasbUserGroupMembersDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group_members.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	userGroupId := idSplit[0]
	memberIdSetStr := idSplit[1]

	if err := service.DeleteDasbUserGroupMembersById(ctx, userGroupId, memberIdSetStr); err != nil {
		return err
	}

	return nil
}
