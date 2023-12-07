package tencentcloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodSubApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSubApplicationCreate,
		Read:   resourceTencentCloudVodSubApplicationRead,
		Update: resourceTencentCloudVodSubApplicationUpdate,
		Delete: resourceTencentCloudVodSubApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 40),
				Description:  "Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Sub appliaction status.",
				ValidateFunc: validateAllowedStringValue(VOD_SUB_APPLICATION_STATUS),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sub application description.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the sub application was created.",
			},
		},
	}
}

func resourceTencentCloudVodSubApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.create")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewCreateSubAppIdRequest()
		subAppId   *uint64
		subAppName *string
	)

	if v, ok := d.GetOk("name"); ok {
		subAppName = helper.String(v.(string))
		request.Name = subAppName
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateSubAppId(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		subAppId = response.Response.SubAppId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*subAppName + FILED_SP + helper.UInt64ToStr(*subAppId))

	if v, ok := d.GetOk("status"); ok {
		statusResquest := vod.NewModifySubAppIdStatusRequest()
		statusResquest.SubAppId = subAppId
		statusResquest.Status = helper.String(v.(string))

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(statusResquest.GetAction())
			_, err := meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySubAppIdStatus(statusResquest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentCloudVodSubApplicationRead(d, meta)
}

func resourceTencentCloudVodSubApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.read")()
	defer inconsistentCheck(d, meta)()

	var (
		//logId   = getLogId(contextNil)
		//ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		client  = meta.(*TencentCloudClient).apiV3Conn
		request = vod.NewDescribeSubAppIdsRequest()
		appInfo = vod.SubAppIdInfo{}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppName := idSplit[0]
	subAppId := idSplit[1]

	request.Name = &subAppName

	response, err := client.UseVodClient().DescribeSubAppIds(request)
	if err != nil {
		return err
	}
	infoSet := response.Response.SubAppIdInfoSet
	if len(infoSet) == 0 {
		d.SetId("")
		return nil
	}

	for _, info := range infoSet {
		if helper.UInt64ToStr(helper.PUint64(info.SubAppId)) == subAppId {
			appInfo = *info
			break
		}
	}
	_ = d.Set("name", appInfo.Name)
	_ = d.Set("description", appInfo.Description)
	// there may hide a bug, appInfo do not return status. So use the user input
	//_ = d.Set("status", appInfo.Status)
	_ = d.Set("status", helper.String(d.Get("status").(string)))
	_ = d.Set("create_time", appInfo.CreateTime)
	return nil
}

func resourceTencentCloudVodSubApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.update")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewModifySubAppIdInfoRequest()
		changeFlag = false
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[1]

	if d.HasChange("name") {
		changeFlag = true
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}
	if d.HasChange("description") {
		changeFlag = true
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}
	request.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))

	if changeFlag {
		var err error
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySubAppIdInfo(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if d.HasChange("status") {
		var statusRequest = vod.NewModifySubAppIdStatusRequest()
		if v, ok := d.GetOk("status"); ok {
			statusRequest.Status = helper.String(v.(string))
		}
		statusRequest.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))
		var err error
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(statusRequest.GetAction())
			_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySubAppIdStatus(statusRequest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
				return retryError(err)
			}
			return nil
		})
	}
	return resourceTencentCloudVodSubApplicationRead(d, meta)
}

func resourceTencentCloudVodSubApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.delete")()
	logId := getLogId(contextNil)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[1]
	var statusRequest = vod.NewModifySubAppIdStatusRequest()
	// first turn off
	statusRequest.Status = helper.String("Off")
	statusRequest.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(statusRequest.GetAction())
		if _, err := meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySubAppIdStatus(statusRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	// then destroy
	statusRequest.Status = helper.String("Destroyed")
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(statusRequest.GetAction())
		if _, err := meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySubAppIdStatus(statusRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
