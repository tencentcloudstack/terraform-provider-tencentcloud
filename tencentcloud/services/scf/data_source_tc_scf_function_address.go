package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfFunctionAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionAddressRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cos address of the function.",
			},

			"code_sha256": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SHA256 code of the function.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_function_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		paramMap["Qualifier"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var res *scf.GetFunctionAddressResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionAddress(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		res = result
		return nil
	})
	if err != nil {
		return err
	}

	resMap := make(map[string]interface{})

	if res.Url != nil {
		_ = d.Set("url", res.Url)
		resMap["url"] = res.Url
	}

	if res.CodeSha256 != nil {
		_ = d.Set("code_sha256", res.CodeSha256)
		resMap["code_sha256"] = res.CodeSha256
	}

	d.SetId(*res.Url)
	output, ok := d.GetOk("result_output_file")

	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), resMap); e != nil {
			return e
		}
	}
	return nil
}
