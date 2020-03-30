/*
Provides a resource to create a CAM group membership.

Example Usage

```hcl
resource "tencentcloud_cam_group_membership" "foo" {
  group_id = tencentcloud_cam_group.foo.id
  user_ids = [tencentcloud_cam_user.foo.id,tencentcloud_cam_user.bar.id]
}
```

Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.foo 12515263
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudCamGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamGroupMembershipCreate,
		Read:   resourceTencentCloudCamGroupMembershipRead,
		Update: resourceTencentCloudCamGroupMembershipUpdate,
		Delete: resourceTencentCloudCamGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of CAM group.",
			},
			"user_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Id set of the CAM group members.",
			},
		},
	}
}

func resourceTencentCloudCamGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_membership.create")()

	logId := getLogId(contextNil)

	groupId := d.Get("group_id").(string)
	members := d.Get("user_ids").(*schema.Set).List()
	err := addUsersToGroup(members, groupId, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create CAM group membership failed, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(groupId)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeGroupMembershipById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		if len(instance) == 0 {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group membership failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(3 * time.Second)
	return resourceTencentCloudCamGroupMembershipRead(d, meta)
}

func resourceTencentCloudCamGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_membership.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupId := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var members []*string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeGroupMembershipById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		members = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	if len(members) == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)
	_ = d.Set("user_ids", members)

	return nil
}

func resourceTencentCloudCamGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_membership.update")()

	logId := getLogId(contextNil)

	groupId := d.Id()

	if d.HasChange("user_ids") {
		o, n := d.GetChange("user_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(remove) > 0 {
			oErr := removeUsersFromGroup(remove, groupId, meta)
			if oErr != nil {
				log.Printf("[CRITAL]%s update CAM group membership failed, reason:%s\n", logId, oErr.Error())
				return oErr
			}
		}
		if len(add) > 0 {
			nErr := addUsersToGroup(add, groupId, meta)
			if nErr != nil {
				log.Printf("[CRITAL]%s update CAM group membership failed, reason:%s\n", logId, nErr.Error())
				return nErr
			}
		}
	}

	return resourceTencentCloudCamGroupMembershipRead(d, meta)
}

func resourceTencentCloudCamGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_membership.delete")()

	logId := getLogId(contextNil)
	groupId := d.Get("group_id").(string)
	members := d.Get("user_ids").(*schema.Set).List()
	err := removeUsersFromGroup(members, groupId, meta)
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}

func getUidFromName(name string, meta interface{}) (uid *uint64, errRet error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserById(ctx, name)
		if e != nil {
			return retryError(e)
		}
		if result == nil || result.Response == nil || result.Response.Uid == nil {
			return nil
		}
		uid = result.Response.Uid
		return nil
	})
	if err != nil {
		errRet = err
	}
	return
}

func addUsersToGroup(members []interface{}, groupId string, meta interface{}) error {
	logId := getLogId(contextNil)

	request := cam.NewAddUserToGroupRequest()
	request.Info = make([]*cam.GroupIdOfUidInfo, 0)
	for _, member := range members {
		var info cam.GroupIdOfUidInfo
		//get uid from name

		uId, e := getUidFromName(member.(string), meta)
		if e != nil {
			return e
		}
		if uId == nil {
			continue
		}
		info.Uid = uId
		groupIdInt, ee := strconv.Atoi(groupId)
		if ee != nil {
			return ee
		}
		groupIdInt64 := uint64(groupIdInt)
		info.GroupId = &groupIdInt64
		request.Info = append(request.Info, &info)
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().AddUserToGroup(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM group membership failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func removeUsersFromGroup(members []interface{}, groupId string, meta interface{}) error {
	logId := getLogId(contextNil)

	request := cam.NewRemoveUserFromGroupRequest()
	request.Info = make([]*cam.GroupIdOfUidInfo, 0)
	for _, member := range members {
		var info cam.GroupIdOfUidInfo
		uId, e := getUidFromName(member.(string), meta)
		if e != nil {
			//notice case when user is deleted, the uin is not found, and the membership is removed in the user module when deleted
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return e
			}
			if ee.Code == "ResourceNotFound.UserNotExist" {
				continue
			} else {
				return e
			}
		}
		if uId == nil {
			continue
		}
		info.Uid = uId
		groupIdInt, eee := strconv.Atoi(groupId)
		if eee != nil {
			return eee
		}
		groupIdInt64 := uint64(groupIdInt)
		info.GroupId = &groupIdInt64
		request.Info = append(request.Info, &info)
	}
	//no exist user need to remove, then return
	if len(request.Info) == 0 {
		return nil
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().RemoveUserFromGroup(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM group membership failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
