package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayAPIKeys() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_api_gateway_api_keys.read")()

	var (
		logId                   = getLogId(contextNil)
		ctx                     = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
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

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		apiKeySet, err = apiGatewayService.DescribeApiKeysStatus(ctx, secretName, accessKeyId)
		if err != nil {
			return retryError(err, InternalError)
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

	d.SetId(strings.Join([]string{secretName, accessKeyId}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
