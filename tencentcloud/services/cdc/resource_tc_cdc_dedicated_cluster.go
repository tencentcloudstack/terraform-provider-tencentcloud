package cdc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdcDedicatedCluster() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudCdcDedicatedClusterCreate,
		Read:   ResourceTencentCloudCdcDedicatedClusterRead,
		Update: ResourceTencentCloudCdcDedicatedClusterUpdate,
		Delete: ResourceTencentCloudCdcDedicatedClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster Site ID.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster Name.",
			},
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster Zone.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster Description.",
			},
		},
	}
}

func ResourceTencentCloudCdcDedicatedClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster.create")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		request            = cdc.NewCreateDedicatedClusterRequest()
		response           = cdc.NewCreateDedicatedClusterResponse()
		dedicatedClusterId string
	)

	if v, ok := d.GetOk("site_id"); ok {
		request.SiteId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().CreateDedicatedCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cdc dedicatedcluster failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cdc dedicatedCluster failed, reason:%+v", logId, err)
		return err
	}

	dedicatedClusterId = *response.Response.DedicatedClusterId
	d.SetId(dedicatedClusterId)

	return ResourceTencentCloudCdcDedicatedClusterRead(d, meta)
}

func ResourceTencentCloudCdcDedicatedClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster.read")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dedicatedClusterId = d.Id()
	)

	dedicatedCluster, err := service.DescribeCdcDedicatedClusterById(ctx, dedicatedClusterId)
	if err != nil {
		return err
	}

	if dedicatedCluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdcDedicatedCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dedicatedCluster.SiteId != nil {
		_ = d.Set("site_id", dedicatedCluster.SiteId)
	}

	if dedicatedCluster.Name != nil {
		_ = d.Set("name", dedicatedCluster.Name)
	}

	if dedicatedCluster.Zone != nil {
		_ = d.Set("zone", dedicatedCluster.Zone)
	}

	if dedicatedCluster.Description != nil {
		_ = d.Set("description", dedicatedCluster.Description)
	}

	return nil
}

func ResourceTencentCloudCdcDedicatedClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster.update")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		request            = cdc.NewModifyDedicatedClusterInfoRequest()
		dedicatedClusterId = d.Id()
	)

	request.DedicatedClusterId = &dedicatedClusterId
	if d.HasChange("site_id") {
		if v, ok := d.GetOk("site_id"); ok {
			request.SiteId = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("zone") {
		if v, ok := d.GetOk("zone"); ok {
			request.Zone = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().ModifyDedicatedClusterInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cdc dedicatedCluster failed, reason:%+v", logId, err)
		return err
	}

	return ResourceTencentCloudCdcDedicatedClusterRead(d, meta)
}

func ResourceTencentCloudCdcDedicatedClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster.delete")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dedicatedClusterId = d.Id()
	)

	if err := service.DeleteCdcDedicatedClusterById(ctx, dedicatedClusterId); err != nil {
		return err
	}

	return nil
}
