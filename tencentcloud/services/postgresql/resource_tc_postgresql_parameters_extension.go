package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlParametersReadPreHandleResponse0(ctx context.Context, resp *postgresv20170312.DescribeDBInstanceParametersResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if v, ok := d.GetOk("param_list"); ok {
		paramList := []*postgresv20170312.ParamInfo{}
		for _, item := range v.([]interface{}) {
			paramListMap := item.(map[string]interface{})
			if name, ok := paramListMap["name"].(string); ok && v != "" {
				if resp.Detail != nil {
					details := resp.Detail
					for _, detail := range details {
						if *detail.Name == name {
							paramList = append(paramList, detail)
						}
					}
				}
			}
		}
		resp.Detail = paramList
	}

	return nil
}

func resourceTencentCloudPostgresqlParametersUpdateRequestOnError0(ctx context.Context, req *postgresv20170312.ModifyDBInstanceParametersRequest, e error) *resource.RetryError {
	if e, ok := e.(*errors.TencentCloudSDKError); ok {
		if e.GetCode() == "FailedOperation.FailedOperationError" {
			change, err := resourceTencentCloudParamChange(ctx)
			if err != nil || change {
				return resource.NonRetryableError(err)
			}
			return nil
		}
	}
	if e != nil {
		return resource.RetryableError(e)
	}
	return nil
}

func resourceTencentCloudParamChange(ctx context.Context) (bool, error) {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return false, fmt.Errorf("resource data can not be nil")
	}
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return false, fmt.Errorf("provider meta can not be nil")
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	respData, err := service.DescribePostgresqlParametersById(ctx, d.Id())
	if err != nil {
		return false, err
	}

	if respData == nil || respData.Detail == nil {
		return false, fmt.Errorf("respData or respData.Detail be nil")
	}

	change := false
	for _, v := range respData.Detail {
		for _, param := range d.Get("param_list").([]interface{}) {
			paramMap := param.(map[string]interface{})
			if paramMap["name"].(string) == *v.Name {
				if paramMap["expected_value"].(string) != *v.CurrentValue {
					change = true
					break
				}
			}
		}
	}
	return change, nil
}

func resourceTencentCloudPostgresqlParametersUpdatePreHandleResponse0(ctx context.Context, resp *postgresv20170312.ModifyDBInstanceParametersResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if _, err := (&resource.StateChangeConf{
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourcePostgresqlParametersUpdateStateRefreshFunc_0_0(ctx, d.Id()),
		Target:     []string{"Success"},
		Timeout:    300 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func resourcePostgresqlParametersUpdateStateRefreshFunc_0_0(ctx context.Context, dBInstanceId string) resource.StateRefreshFunc {
	var req *postgresv20170312.DescribeTasksRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = postgresv20170312.NewDescribeTasksRequest()
			req.DBInstanceId = helper.String(dBInstanceId)

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().DescribeTasksWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		for _, task := range resp.Response.TaskSet {
			if *task.TaskType == "ModifyInstanceParams" {
				if *task.Status != "Success" {
					return resp.Response, fmt.Sprintf("%v", *task.Status), nil
				}
			}
		}

		return resp.Response, "Success", nil
	}
}
