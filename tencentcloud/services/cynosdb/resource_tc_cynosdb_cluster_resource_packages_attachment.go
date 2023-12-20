package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbClusterResourcePackagesAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterResourcePackagesAttachmentCreate,
		Read:   resourceTencentCloudCynosdbClusterResourcePackagesAttachmentRead,
		Delete: resourceTencentCloudCynosdbClusterResourcePackagesAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"package_ids": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource Package Unique ID.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterResourcePackagesAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = cynosdb.NewBindClusterResourcePackagesRequest()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("package_ids"); ok {
		packageIdsSet := v.(*schema.Set).List()
		for i := range packageIdsSet {
			packageIds := packageIdsSet[i].(string)
			request.PackageIds = append(request.PackageIds, &packageIds)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().BindClusterResourcePackages(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("cynosdb clusterResourcePackagesAttachment not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterResourcePackagesAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbClusterResourcePackagesAttachmentRead(d, meta)
}

func resourceTencentCloudCynosdbClusterResourcePackagesAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	clusterResourcePackagesAttachment, err := service.DescribeCynosdbClusterResourcePackagesAttachmentById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterResourcePackagesAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterResourcePackagesAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterResourcePackagesAttachment.ClusterId != nil {
		_ = d.Set("cluster_id", clusterResourcePackagesAttachment.ClusterId)
	}

	if clusterResourcePackagesAttachment.ResourcePackages != nil {
		tmpList := []interface{}{}
		for _, v := range clusterResourcePackagesAttachment.ResourcePackages {
			if v.PackageId != nil {
				tmpList = append(tmpList, v.PackageId)
			}
		}
		_ = d.Set("package_ids", tmpList)
	}

	return nil
}

func resourceTencentCloudCynosdbClusterResourcePackagesAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId     = d.Id()
		packageIdsSet []*string
	)

	if v, ok := d.GetOk("package_ids"); ok {
		idsSet := v.(*schema.Set).List()
		for i := range idsSet {
			ids := idsSet[i].(string)
			packageIdsSet = append(packageIdsSet, &ids)
		}

		if err := service.DeleteCynosdbClusterResourcePackagesAttachmentById(ctx, clusterId, packageIdsSet); err != nil {
			return err
		}
	}

	return nil
}
