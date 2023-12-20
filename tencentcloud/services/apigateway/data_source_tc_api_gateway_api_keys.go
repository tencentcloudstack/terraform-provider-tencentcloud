package apigateway

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayAPIKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIKeysRead,

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom key name.",
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Created API key ID, this field is exactly the same as ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of API keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API key ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key status. Values: `on`, `off`.",
						},
						"access_key_secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created API key.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_keys.read")()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService       = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiKeySet               []*apigateway.ApiKey
		secretName, accessKeyId string
		err                     error
	)

	if v, ok := d.GetOk("secret_name"); ok {
		secretName = v.(string)
	}
	if v, ok := d.GetOk("api_key_id"); ok {
		accessKeyId = v.(string)
	}

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		apiKeySet, err = apiGatewayService.DescribeApiKeysStatus(ctx, secretName, accessKeyId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(apiKeySet))
	for _, apiKey := range apiKeySet {
		list = append(list, map[string]interface{}{
			"api_key_id":        apiKey.AccessKeyId,
			"status":            API_GATEWAY_KEY_INT2STRS[*apiKey.Status],
			"access_key_secret": apiKey.AccessKeySecret,
			"modify_time":       apiKey.ModifiedTime,
			"create_time":       apiKey.CreatedTime,
		})
	}

	if err := d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{secretName, accessKeyId}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
