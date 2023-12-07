package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func datasourceTencentCloudKubernetesClusterCommonNames() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudKubernetesClusterCommonNamesRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster ID.",
			},
			"subaccount_uins": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of sub-account. Up to 50 sub-accounts can be passed in at a time.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"role_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Role ID. Up to 50 sub-accounts can be passed in at a time.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save result.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of the CommonName in the certificate of the client corresponding to the sub-account UIN.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subaccount_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User UIN.",
						},
						"common_names": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CommonName in the certificate of the client corresponding to the sub-account.",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudKubernetesClusterCommonNamesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("datasource.tencentcloud_kubernetes_cluster_common_names.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client}

	clusterId := d.Get("cluster_id").(string)
	request := tke.NewDescribeClusterCommonNamesRequest()
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("subaccount_uins"); ok {
		request.SubaccountUins = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("role_ids"); ok {
		request.RoleIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	names, err := service.DescribeClusterCommonNames(ctx, request)

	if err != nil {
		return err
	}

	result := make([]interface{}, 0, len(names))
	cns := make([]string, 0)

	for i := range names {
		cn := names[i]
		result = append(result, map[string]interface{}{
			"subaccount_uin": cn.SubaccountUin,
			"common_names":   cn.CN,
		})
		cns = append(cns, *cn.CN)
	}

	if err := d.Set("list", result); err != nil {
		return err
	}

	d.SetId(clusterId + FILED_SP + helper.DataResourceIdsHash(cns))

	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), result)
	}

	return nil
}
