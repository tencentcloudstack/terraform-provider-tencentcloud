package tke

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudKubernetesClusterAdminRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterAdminRoleRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},
			"request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request ID returned by the API.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKubernetesClusterAdminRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_cluster_admin_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clusterId := d.Get("cluster_id").(string)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var requestId string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.AcquireClusterAdminRole(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		requestId = result
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(clusterId)
	_ = d.Set("request_id", requestId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		data := map[string]interface{}{
			"cluster_id": clusterId,
			"request_id": requestId,
		}
		if e := tccommon.WriteToFile(output.(string), data); e != nil {
			return e
		}
	}

	return nil
}
