package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKeyPairsRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_name", "project_id"},
				Description:   "ID of the key pair to be queried.",
			},

			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_id"},
				Description:   "Name of the key pair to be queried. Support regular expression search, only `^` and `$` are supported.",
			},

			"key_pair_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of key pair. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the key pair.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the key pair.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the key pair.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID of the key pair.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public key of the key pair.",
						},
					},
				},
			},

			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"key_id"},
				Description:   "Project ID of the key pair to be queried.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_key_pairs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var respData *cvm.DescribeKeyPairsResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKeyPairsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudKeyPairsReadPostHandleResponse0(ctx, paramMap, respData); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudKeyPairsReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
