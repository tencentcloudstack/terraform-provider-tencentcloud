package gaap

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapProxyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapProxyGroupCreate,
		Read:   resourceTencentCloudGaapProxyGroupRead,
		Update: resourceTencentCloudGaapProxyGroupUpdate,
		Delete: resourceTencentCloudGaapProxyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "ID of the project to which the proxy group belongs.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Channel group alias.",
			},

			"real_server_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "real server region, refer to the interface DescribeDestRegions to return the RegionId in the parameter RegionDetail.",
			},

			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP version, can be taken as IPv4 or IPv6 with a default value of IPv4.",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Package type of channel group. Available values: Thunder and Accelerator. Default is Thunder.",
			},
		},
	}
}

func resourceTencentCloudGaapProxyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = gaap.NewCreateProxyGroupRequest()
		response = gaap.NewCreateProxyGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("real_server_region"); ok {
		request.RealServerRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_version"); ok {
		request.IPAddressVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGaapClient().CreateProxyGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gaap proxyGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(groupId)

	return resourceTencentCloudGaapProxyGroupRead(d, meta)
}

func resourceTencentCloudGaapProxyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()

	proxyGroup, err := service.DescribeGaapProxyGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if proxyGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `GaapProxyGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxyGroup.ProjectId != nil {
		_ = d.Set("project_id", proxyGroup.ProjectId)
	}

	if proxyGroup.GroupName != nil {
		_ = d.Set("group_name", proxyGroup.GroupName)
	}

	if proxyGroup.RealServerRegionInfo != nil && proxyGroup.RealServerRegionInfo.RegionId != nil {
		_ = d.Set("real_server_region", proxyGroup.RealServerRegionInfo.RegionId)
	}

	if proxyGroup.IPAddressVersion != nil {
		_ = d.Set("ip_address_version", proxyGroup.IPAddressVersion)
	}

	if proxyGroup.PackageType != nil {
		_ = d.Set("package_type", proxyGroup.PackageType)
	}

	return nil
}

func resourceTencentCloudGaapProxyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := gaap.NewModifyProxyGroupAttributeRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"real_server_region", "ip_address_version", "package_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	isChanged := false
	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntUint64(v.(int))
			isChanged = true
		}
	}

	if d.HasChange("group_name") {
		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
			isChanged = true
		}
	}

	if isChanged {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGaapClient().ModifyProxyGroupAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update gaap proxyGroup failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudGaapProxyGroupRead(d, meta)
}

func resourceTencentCloudGaapProxyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	groupId := d.Id()

	if err := service.DeleteGaapProxyGroupById(ctx, groupId); err != nil {
		return err
	}

	return nil
}
