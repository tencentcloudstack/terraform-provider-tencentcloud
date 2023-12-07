package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfRepositoryRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query keywords (search by Repository name).",
			},

			"repository_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Repository type (default Repository: default, private Repository: private).",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A list of Repository information that meets the query criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total Repository.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Repository information list. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "repository Id.",
									},
									"repository_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository Name.",
									},
									"repository_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository type (default Repository: default, private Repository: private).",
									},
									"repository_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository description (default warehouse: default, private warehouse: private).",
									},
									"is_used": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the repository is being used. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreationTime. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository bucket name. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"bucket_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository region. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository Directory. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTsfRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_repository.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_type"); ok {
		paramMap["RepositoryType"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var repository *tsf.RepositoryList
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfRepositoryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		repository = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(repository.Content))
	repositoryListMap := map[string]interface{}{}
	if repository != nil {
		if repository.TotalCount != nil {
			repositoryListMap["total_count"] = repository.TotalCount
		}

		if repository.Content != nil {
			contentList := []interface{}{}
			for _, content := range repository.Content {
				contentMap := map[string]interface{}{}

				if content.RepositoryId != nil {
					contentMap["repository_id"] = content.RepositoryId
				}

				if content.RepositoryName != nil {
					contentMap["repository_name"] = content.RepositoryName
				}

				if content.RepositoryType != nil {
					contentMap["repository_type"] = content.RepositoryType
				}

				if content.RepositoryDesc != nil {
					contentMap["repository_desc"] = content.RepositoryDesc
				}

				if content.IsUsed != nil {
					contentMap["is_used"] = content.IsUsed
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.BucketName != nil {
					contentMap["bucket_name"] = content.BucketName
				}

				if content.BucketRegion != nil {
					contentMap["bucket_region"] = content.BucketRegion
				}

				if content.Directory != nil {
					contentMap["directory"] = content.Directory
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.RepositoryId)
			}

			repositoryListMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{repositoryListMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), repositoryListMap); e != nil {
			return e
		}
	}
	return nil
}
