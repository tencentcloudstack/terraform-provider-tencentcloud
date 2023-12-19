package cam

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamGroupCreate,
		Read:   resourceTencentCloudCamGroupRead,
		Update: resourceTencentCloudCamGroupUpdate,
		Delete: resourceTencentCloudCamGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of CAM group.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CAM group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM group.",
			},
		},
	}
}

func resourceTencentCloudCamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_group.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cam.NewCreateGroupRequest()
	request.GroupName = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	var response *cam.CreateGroupResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().CreateGroup(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "GroupNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.GroupId == nil {
		return fmt.Errorf("CAM group id is nil")
	}
	d.SetId(strconv.Itoa(int(*response.Response.GroupId)))

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	groupId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeGroupById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance == nil || instance.Response == nil || instance.Response.GroupId == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamGroupRead(d, meta)
}

func resourceTencentCloudCamGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	groupId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.GetGroupResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeGroupById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.GroupId == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", *instance.Response.GroupName)
	_ = d.Set("create_time", *instance.Response.CreateTime)
	if instance.Response.Remark != nil {
		_ = d.Set("remark", *instance.Response.Remark)
	}
	return nil
}

func resourceTencentCloudCamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_group.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	groupId := d.Id()
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		return e
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewUpdateGroupRequest()
	request.GroupId = &groupIdInt64
	changeFlag := false

	mutableArgs := []string{"name", "remark"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			changeFlag = true
			break
		}
	}

	if changeFlag {
		request.GroupName = helper.String(d.Get("name").(string))
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateGroup(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM group description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCamGroupRead(d, meta)
}

func resourceTencentCloudCamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	groupId := d.Id()
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		return e
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewDeleteGroupRequest()
	request.GroupId = &groupIdInt64
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().DeleteGroup(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
