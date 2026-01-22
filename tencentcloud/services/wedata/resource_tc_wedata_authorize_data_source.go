package wedata

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataAuthorizeDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataAuthorizeDataSourceCreate,
		Read:   resourceTencentCloudWedataAuthorizeDataSourceRead,
		Update: resourceTencentCloudWedataAuthorizeDataSourceUpdate,
		Delete: resourceTencentCloudWedataAuthorizeDataSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_source_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data source ID.",
			},

			"auth_project_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of target project ID to be authorized.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"auth_users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of users under the authorized project, format: project_id_user_id.\nWhen authorizing multiple objects, the project ID must be consistent.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudWedataAuthorizeDataSourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_authorize_data_source.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		dataSourceId string
	)

	if v, ok := d.GetOk("data_source_id"); ok {
		dataSourceId = v.(string)
	}

	d.SetId(dataSourceId)

	if v, ok := d.GetOk("auth_project_ids"); ok {
		ProjectIdSet := v.(*schema.Set).List()
		for i := range ProjectIdSet {
			projectId := ProjectIdSet[i].(string)
			request := wedatav20250806.NewAuthorizeDataSourceRequest()
			response := wedatav20250806.NewAuthorizeDataSourceResponse()
			request.DataSourceId = &dataSourceId
			request.AuthProjectId = &projectId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AuthorizeDataSourceWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Data == nil {
					return resource.NonRetryableError(fmt.Errorf("Create wedata authorize data source failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s create wedata authorize data source failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.Data.Status == nil || !*response.Response.Data.Status {
				return fmt.Errorf("Create wedata authorize data source failed, Status is false")
			}
		}
	}

	if v, ok := d.GetOk("auth_users"); ok {
		request := wedatav20250806.NewAuthorizeDataSourceRequest()
		response := wedatav20250806.NewAuthorizeDataSourceResponse()
		authUsersSet := v.(*schema.Set).List()
		for i := range authUsersSet {
			authUser := authUsersSet[i].(string)
			request.AuthUsers = append(request.AuthUsers, helper.String(authUser))
		}

		request.DataSourceId = &dataSourceId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AuthorizeDataSourceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil {
				return resource.NonRetryableError(fmt.Errorf("Create wedata authorize data source failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create wedata authorize data source failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.Data.Status == nil || !*response.Response.Data.Status {
			return fmt.Errorf("Create wedata authorize data source failed, Status is false")
		}
	}

	return resourceTencentCloudWedataAuthorizeDataSourceRead(d, meta)
}

func resourceTencentCloudWedataAuthorizeDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_authorize_data_source.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataSourceId = d.Id()
	)

	respData, err := service.DescribeWedataAuthorizeDataSourceById(ctx, dataSourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_authorize_data_source` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("data_source_id", dataSourceId)

	if respData.AuthProjectIds != nil {
		_ = d.Set("auth_project_ids", respData.AuthProjectIds)
	}

	if respData.AuthUsers != nil {
		_ = d.Set("auth_users", respData.AuthUsers)
	}

	return nil
}

func resourceTencentCloudWedataAuthorizeDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_authorize_data_source.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		dataSourceId = d.Id()
	)

	if d.HasChange("auth_project_ids") {
		oldInterface, newInterface := d.GetChange("auth_project_ids")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()
		if len(remove) > 0 {
			for _, item := range remove {
				request := wedatav20250806.NewRevokeDataSourceAuthorizationRequest()
				response := wedatav20250806.NewRevokeDataSourceAuthorizationResponse()
				request.DataSourceId = &dataSourceId
				request.RevokeProjectId = helper.String(item.(string))
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RevokeDataSourceAuthorizationWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.Data == nil {
						return resource.NonRetryableError(fmt.Errorf("Delete wedata authorize data source failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s delete wedata authorize data source failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.Data.Status == nil || !*response.Response.Data.Status {
					return fmt.Errorf("Delete wedata authorize data source failed, Status is false")
				}
			}
		}

		if len(add) > 0 {
			for _, item := range add {
				request := wedatav20250806.NewAuthorizeDataSourceRequest()
				response := wedatav20250806.NewAuthorizeDataSourceResponse()
				request.DataSourceId = &dataSourceId
				request.AuthProjectId = helper.String(item.(string))
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AuthorizeDataSourceWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.Data == nil {
						return resource.NonRetryableError(fmt.Errorf("Create wedata authorize data source failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s create wedata authorize data source failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.Data.Status == nil || !*response.Response.Data.Status {
					return fmt.Errorf("Create wedata authorize data source failed, Status is false")
				}
			}
		}
	}

	if d.HasChange("auth_users") {
		oldInterface, newInterface := d.GetChange("auth_users")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()
		if len(remove) > 0 {
			for _, item := range remove {
				request := wedatav20250806.NewRevokeDataSourceAuthorizationRequest()
				response := wedatav20250806.NewRevokeDataSourceAuthorizationResponse()
				request.DataSourceId = &dataSourceId
				request.RevokeUser = helper.String(item.(string))
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RevokeDataSourceAuthorizationWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.Data == nil {
						return resource.NonRetryableError(fmt.Errorf("Delete wedata authorize data source failed, Response is nil."))
					}

					response = result
					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s delete wedata authorize data source failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				if response.Response.Data.Status == nil || !*response.Response.Data.Status {
					return fmt.Errorf("Delete wedata authorize data source failed, Status is false")
				}
			}
		}

		if len(add) > 0 {
			request := wedatav20250806.NewAuthorizeDataSourceRequest()
			response := wedatav20250806.NewAuthorizeDataSourceResponse()
			for _, item := range add {
				request.AuthUsers = append(request.AuthUsers, helper.String(item.(string)))
			}

			request.DataSourceId = &dataSourceId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AuthorizeDataSourceWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Data == nil {
					return resource.NonRetryableError(fmt.Errorf("Create wedata authorize data source failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s create wedata authorize data source failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.Data.Status == nil || !*response.Response.Data.Status {
				return fmt.Errorf("Create wedata authorize data source failed, Status is false")
			}
		}
	}

	return resourceTencentCloudWedataAuthorizeDataSourceRead(d, meta)
}

func resourceTencentCloudWedataAuthorizeDataSourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_authorize_data_source.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		dataSourceId = d.Id()
	)

	request := wedatav20250806.NewRevokeDataSourceAuthorizationRequest()
	response := wedatav20250806.NewRevokeDataSourceAuthorizationResponse()
	request.DataSourceId = &dataSourceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RevokeDataSourceAuthorizationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata authorize data source failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata authorize data source failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.Status == nil || !*response.Response.Data.Status {
		return fmt.Errorf("Delete wedata authorize data source failed, Status is false")
	}

	return nil
}
