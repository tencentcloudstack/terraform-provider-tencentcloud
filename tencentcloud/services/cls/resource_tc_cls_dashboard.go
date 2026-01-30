package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsDashboard() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClsDashboardRead,
		Create: resourceTencentCloudClsDashboardCreate,
		Update: resourceTencentCloudClsDashboardUpdate,
		Delete: resourceTencentCloudClsDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dashboard_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Dashboard name, which must be unique within the account.",
			},

			"data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dashboard configuration data in JSON format. If not specified, an empty dashboard will be created.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag key-value pairs. Maximum of 10 tags.",
			},

			"dashboard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dashboard ID (globally unique identifier).",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time.",
			},
		},
	}
}

func resourceTencentCloudClsDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dashboard.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateDashboardRequest()
		response = cls.NewCreateDashboardResponse()
	)

	if v, ok := d.GetOk("dashboard_name"); ok {
		request.DashboardName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data"); ok {
		request.Data = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   helper.String(k),
				Value: helper.String(v),
			})
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateDashboard(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls dashboard failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls dashboard failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.DashboardId == nil {
		return fmt.Errorf("DashboardId is nil.")
	}

	dashboardId := *response.Response.DashboardId
	d.SetId(dashboardId)

	return resourceTencentCloudClsDashboardRead(d, meta)
}

func resourceTencentCloudClsDashboardRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dashboard.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dashboardId = d.Id()
	)

	dashboard, err := service.DescribeClsDashboardById(ctx, dashboardId)
	if err != nil {
		return err
	}

	if dashboard == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cls_dashboard` %s does not exist", logId, dashboardId)
		return nil
	}

	if dashboard.DashboardName != nil {
		_ = d.Set("dashboard_name", dashboard.DashboardName)
	}

	if dashboard.Data != nil {
		_ = d.Set("data", dashboard.Data)
	}

	if dashboard.DashboardId != nil {
		_ = d.Set("dashboard_id", dashboard.DashboardId)
	}

	if dashboard.CreateTime != nil {
		_ = d.Set("create_time", dashboard.CreateTime)
	}

	if dashboard.UpdateTime != nil {
		_ = d.Set("update_time", dashboard.UpdateTime)
	}

	// Handle tags from dashboard response
	if dashboard.Tags != nil {
		tags := make(map[string]string)
		for _, tag := range dashboard.Tags {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudClsDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dashboard.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = cls.NewModifyDashboardRequest()
		dashboardId = d.Id()
	)

	request.DashboardId = helper.String(dashboardId)

	if d.HasChange("dashboard_name") {
		if v, ok := d.GetOk("dashboard_name"); ok {
			request.DashboardName = helper.String(v.(string))
		}
	}

	if d.HasChange("data") {
		if v, ok := d.GetOk("data"); ok {
			request.Data = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
			for k, v := range tags {
				request.Tags = append(request.Tags, &cls.Tag{
					Key:   helper.String(k),
					Value: helper.String(v),
				})
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyDashboard(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cls dashboard failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsDashboardRead(d, meta)
}

func resourceTencentCloudClsDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_dashboard.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = cls.NewDeleteDashboardRequest()
		dashboardId = d.Id()
	)

	request.DashboardId = helper.String(dashboardId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteDashboard(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls dashboard failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
