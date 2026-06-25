package dbdc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewDbdcService(client *connectivity.TencentCloudClient) DbdcService {
	return DbdcService{client: client}
}

type DbdcService struct {
	client *connectivity.TencentCloudClient
}

func (me *DbdcService) DescribeDBCustomClustersByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.DBCustomCluster, totalCount int64, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = dbdcv20201029.NewDescribeDBCustomClustersRequest()
		response = dbdcv20201029.NewDescribeDBCustomClustersResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterIds" {
			request.ClusterIds = v.([]*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*dbdcv20201029.Filter)
		}
		if k == "Tags" {
			request.Tags = v.([]*dbdcv20201029.Tag)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseDbdcV20201029Client().DescribeDBCustomClusters(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.ClusterSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe dbdc_db_custom_clusters failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.ClusterSet...)
		if totalCount == 0 && response.Response.TotalCount != nil {
			totalCount = *response.Response.TotalCount
		}

		if len(response.Response.ClusterSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbdcService) DescribeDBCustomImagesByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.DBCustomImage, totalCount int64, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = dbdcv20201029.NewDescribeDBCustomImagesRequest()
		response = dbdcv20201029.NewDescribeDBCustomImagesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseDbdcV20201029Client().DescribeDBCustomImages(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.ImageSet == nil {
				log.Printf("[DATASOURCE] read empty, skip SetId")
				return resource.NonRetryableError(fmt.Errorf("Describe dbdc_db_custom_images failed, Response is nil or empty."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.ImageSet...)
		if totalCount == 0 && response.Response.TotalCount != nil {
			totalCount = *response.Response.TotalCount
		}

		if len(response.Response.ImageSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
