/*
Use this data source to query detailed information of TCR repositories.

Example Usage

```hcl
data "tencentcloud_tcr_repositories" "name" {
  name       = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTCRRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the TCR instance that the repository belongs to.",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the namespace that the repository belongs to.",
			},
			"repository_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the TCR repositories to query.",
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
						//Computed values
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
	defer logElapsed("data_source.tencentcloud_tcr_repositories.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var namespaceName, instanceId, repositoryName string
	instanceId = d.Get("instance_id").(string)
	namespaceName = d.Get("namespace_name").(string)

	if v, ok := d.GetOk("repository_name"); ok {
		repositoryName = v.(string)
	}

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error
	var domain string
	instance, has, outErr := tcrService.DescribeTCRInstanceById(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = tcrService.DescribeTCRInstanceById(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
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
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			repositories, inErr = tcrService.DescribeTCRRepositories(ctx, instanceId, namespaceName, repositoryName)
			if inErr != nil {
				return retryError(inErr)
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
		ids = append(ids, instanceId+FILED_SP+*repository.Namespace+FILED_SP+*repository.Name)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("repository_list", repositoryList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR repository list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), repositoryList); e != nil {
			return e
		}
	}

	return nil

}
