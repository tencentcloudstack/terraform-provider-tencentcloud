package tcr

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTCRRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the TCR instance that the repository belongs to.",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the namespace that the repository belongs to.",
			},
			"repository_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the TCR repositories to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"repository_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated TCR repositories.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of repository.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the namespace that the repository belongs to.",
						},
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate that the repository is public or not.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
						"brief_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Brief description of the repository.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the repository.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the repository.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRRepositoriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_repositories.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var namespaceName, instanceId, repositoryName string
	instanceId = d.Get("instance_id").(string)
	namespaceName = d.Get("namespace_name").(string)

	if v, ok := d.GetOk("repository_name"); ok {
		repositoryName = v.(string)
	}

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error
	var domain string
	instance, has, outErr := tcrService.DescribeTCRInstanceById(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = tcrService.DescribeTCRInstanceById(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if has {
		domain = *instance.PublicDomain
	}

	repositories, outErr := tcrService.DescribeTCRRepositories(ctx, instanceId, namespaceName, repositoryName)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			repositories, inErr = tcrService.DescribeTCRRepositories(ctx, instanceId, namespaceName, repositoryName)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(repositories))
	repositoryList := make([]map[string]interface{}, 0, len(repositories))
	for _, repository := range repositories {
		name := strings.Replace(*repository.Name, fmt.Sprintf("%s/", *repository.Namespace), "", 1)
		mapping := map[string]interface{}{
			"name":           name,
			"namespace_name": repository.Namespace,
			"is_public":      repository.Public,
			"create_time":    repository.CreationTime,
			"update_time":    repository.UpdateTime,
			"brief_desc":     repository.BriefDescription,
			"description":    repository.Description,
			"url":            fmt.Sprintf("%s/%s", domain, *repository.Name),
		}

		repositoryList = append(repositoryList, mapping)
		ids = append(ids, instanceId+tccommon.FILED_SP+*repository.Namespace+tccommon.FILED_SP+*repository.Name)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("repository_list", repositoryList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR repository list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), repositoryList); e != nil {
			return e
		}
	}

	return nil

}
