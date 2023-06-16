package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pkg/errors"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TagService struct {
	client *connectivity.TencentCloudClient
}

func (me *TagService) ModifyTags(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
	request := tag.NewModifyResourceTagsRequest()
	request.Resource = &resourceName
	if len(replaceTags) > 0 {
		request.ReplaceTags = make([]*tag.Tag, 0, len(replaceTags))
		for k, v := range replaceTags {
			key := k
			value := v
			replaceTag := &tag.Tag{
				TagKey:   &key,
				TagValue: &value,
			}
			request.ReplaceTags = append(request.ReplaceTags, replaceTag)
		}
	}
	if len(deleteKeys) > 0 {
		request.DeleteTags = make([]*tag.TagKeyObject, 0, len(deleteKeys))
		for _, v := range deleteKeys {
			key := v
			deleteKey := &tag.TagKeyObject{
				TagKey: &key,
			}
			request.DeleteTags = append(request.DeleteTags, deleteKey)
		}
	}

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseTagClient().ModifyResourceTags(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *TagService) DescribeResourceTags(ctx context.Context, serviceType, resourceType, region, resourceId string) (tags map[string]string, err error) {
	request := tag.NewDescribeResourceTagsByResourceIdsRequest()
	request.ServiceType = &serviceType
	request.ResourcePrefix = &resourceType
	request.ResourceRegion = &region
	request.ResourceIds = []*string{&resourceId}
	request.Limit = helper.IntUint64(DESCRIBE_TAGS_LIMIT)

	var offset uint64
	request.Offset = &offset

	// for run loop at least once
	count := DESCRIBE_TAGS_LIMIT
	for count == DESCRIBE_TAGS_LIMIT {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseTagClient().DescribeResourceTagsByResourceIds(request)
			if err != nil {
				count = 0

				return retryError(errors.WithStack(err))
			}

			allTags := response.Response.Tags
			count = len(allTags)

			for _, t := range allTags {
				if *t.ResourceId != resourceId {
					continue
				}
				if tags == nil {
					tags = make(map[string]string)
				}

				tags[*t.TagKey] = *t.TagValue
			}

			return nil
		}); err != nil {
			return nil, err
		}

		offset += uint64(count)
	}

	return
}

func diffTags(oldTags, newTags map[string]interface{}) (replaceTags map[string]string, deleteTags []string) {
	replaceTags = make(map[string]string)
	deleteTags = make([]string, 0)
	for k, v := range newTags {
		_, ok := oldTags[k]
		if !ok || oldTags[k].(string) != v.(string) {
			replaceTags[k] = v.(string)
		}
	}
	for k := range oldTags {
		_, ok := newTags[k]
		if !ok {
			deleteTags = append(deleteTags, k)
		}
	}
	return
}

func (me *TagService) DescribeProjectById(ctx context.Context, projectId uint64) (project *tag.Project, disable *uint64, errRet error) {
	logId := getLogId(ctx)

	request := tag.NewDescribeProjectsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	// query enable project
	request.AllList = helper.Uint64(0)

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*tag.Project, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		instances = append(instances, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}
		offset += limit
	}

	for _, instance := range instances {
		if *instance.ProjectId == projectId {
			project = instance
			disable = helper.Uint64(0)
			break
		}
	}

	if project != nil {
		return
	}

	// query all project
	offset = 0
	limit = 20

	request.AllList = helper.Uint64(1)
	instances = make([]*tag.Project, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		instances = append(instances, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}
		offset += limit
	}

	for _, instance := range instances {
		if *instance.ProjectId == projectId {
			project = instance
			disable = helper.Uint64(1)
			break
		}
	}

	return
}

func (me *TagService) DisableProjectById(ctx context.Context, projectId uint64) (errRet error) {
	logId := getLogId(ctx)

	request := tag.NewUpdateProjectRequest()
	request.ProjectId = &projectId
	request.Disable = helper.Int64(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().UpdateProject(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TagService) DescribeProjects(ctx context.Context, param map[string]interface{}) (project []*tag.Project, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tag.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AllList" {
			request.AllList = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		project = append(project, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
