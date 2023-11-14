/*
Provides a resource to create a cam user_account_group

Example Usage

```hcl
resource "tencentcloud_cam_user_account_group" "user_account_group" {
  info {
		group_id = &lt;nil&gt;
		uid = &lt;nil&gt;
		uin = &lt;nil&gt;

  }
}
```

Import

cam user_account_group can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_account_group.user_account_group user_account_group_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCamUserAccountGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserAccountGroupCreate,
		Read:   resourceTencentCloudCamUserAccountGroupRead,
		Delete: resourceTencentCloudCamUserAccountGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The association relationship between the sub-user and the group ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Group Id.",
						},
						"uid": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sub-user uid. at least one of Uid and Uin must be filled in.",
						},
						"uin": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sub-user uin. at least one of Uid and Uin must be filled in.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCamUserAccountGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_account_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewAddUserToGroupRequest()
		response = cam.NewAddUserToGroupResponse()
		subUin   int
		groupId  int
	)
	if v, ok := d.GetOk("info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			groupIdOfUidInfo := cam.GroupIdOfUidInfo{}
			if v, ok := dMap["group_id"]; ok {
				groupIdOfUidInfo.GroupId = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["uid"]; ok {
				groupIdOfUidInfo.Uid = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["uin"]; ok {
				groupIdOfUidInfo.Uin = helper.IntUint64(v.(int))
			}
			request.Info = append(request.Info, &groupIdOfUidInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().AddUserToGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam UserAccountGroup failed, reason:%+v", logId, err)
		return err
	}

	subUin = *response.Response.SubUin
	d.SetId(strings.Join([]string{helper.Int64ToStr(int64(subUin)), helper.Int64ToStr(int64(groupId))}, FILED_SP))

	return resourceTencentCloudCamUserAccountGroupRead(d, meta)
}

func resourceTencentCloudCamUserAccountGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_account_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	subUin := idSplit[0]
	groupId := idSplit[1]

	UserAccountGroup, err := service.DescribeCamUserAccountGroupById(ctx, subUin, groupId)
	if err != nil {
		return err
	}

	if UserAccountGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamUserAccountGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if UserAccountGroup.Info != nil {
		infoList := []interface{}{}
		for _, info := range UserAccountGroup.Info {
			infoMap := map[string]interface{}{}

			if UserAccountGroup.Info.GroupId != nil {
				infoMap["group_id"] = UserAccountGroup.Info.GroupId
			}

			if UserAccountGroup.Info.Uid != nil {
				infoMap["uid"] = UserAccountGroup.Info.Uid
			}

			if UserAccountGroup.Info.Uin != nil {
				infoMap["uin"] = UserAccountGroup.Info.Uin
			}

			infoList = append(infoList, infoMap)
		}

		_ = d.Set("info", infoList)

	}

	return nil
}

func resourceTencentCloudCamUserAccountGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_account_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	subUin := idSplit[0]
	groupId := idSplit[1]

	if err := service.DeleteCamUserAccountGroupById(ctx, subUin, groupId); err != nil {
		return err
	}

	return nil
}
