package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type OssService struct {
	client *connectivity.AliyunClient
}

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	raw, err := s.client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(id)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return response, WrapErrorf(err, NotFoundMsg, AliyunOssGoSdk)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetBucketInfo", AliyunOssGoSdk)
	}

	addDebug("GetBucketInfo", raw)
	response, _ = raw.(oss.GetBucketInfoResult)
	return
}

func (s *OssService) WaitForOssBucket(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeOssBucket(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.BucketInfo.Name != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.BucketInfo.Name, status, ProviderERROR)
		}
	}
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", AliyunOssGoSdk)
		}
		addDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(exist), status, ProviderERROR)
		}
	}
}
