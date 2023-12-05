package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTCRTokens() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRTokensRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the instance that the token belongs to.",
			},
			"token_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the TCR token to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"token_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated TCR tokens.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of TCR token.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate that the token is enabled or not.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the token.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRTokensRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_tokens.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var tokenId, instanceId string
	instanceId = d.Get("instance_id").(string)
	if v, ok := d.GetOk("token_id"); ok {
		tokenId = v.(string)
	}

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error
	tokens, outErr := tcrService.DescribeTCRTokens(ctx, instanceId, tokenId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			tokens, inErr = tcrService.DescribeTCRTokens(ctx, instanceId, tokenId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(tokens))
	tokenList := make([]map[string]interface{}, 0, len(tokens))
	for i := range tokens {
		token := tokens[i]
		mapping := map[string]interface{}{
			"token_id":    token.Id,
			"enable":      token.Enabled,
			"create_time": token.CreatedAt,
			"description": token.Desc,
		}

		tokenList = append(tokenList, mapping)
		ids = append(ids, instanceId+FILED_SP+*token.Id)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("token_list", tokenList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR token list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tokenList); e != nil {
			return e
		}
	}

	return nil

}
