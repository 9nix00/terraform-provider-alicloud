package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamLoginProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamLoginProfileCreate,
		Read:   resourceAlicloudRamLoginProfileRead,
		Update: resourceAlicloudRamLoginProfileUpdate,
		Delete: resourceAlicloudRamLoginProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"password_reset_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mfa_bind_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudRamLoginProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramSercvice := RamService{client}
	request := ram.CreateCreateLoginProfileRequest()
	request.UserName = d.Get("user_name").(string)
	request.Password = d.Get("password").(string)
	if v, ok := d.GetOk("password_reset_required"); ok {
		request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("mfa_bind_required"); ok {
		request.MFABindRequired = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateLoginProfile(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_login_profile", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId(request.UserName)
	err = ramSercvice.WaitForRamLoginProfile(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamLoginProfileRead(d, meta)
}

func resourceAlicloudRamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateLoginProfileRequest()
	request.Password = d.Get("password").(string)
	request.UserName = d.Id()

	if d.HasChange("password_reset_required") {
		request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(d.Get("password_reset_required").(bool)))
	}

	if d.HasChange("mfa_bind_required") {
		request.MFABindRequired = requests.Boolean(strconv.FormatBool(d.Get("mfa_bind_required").(bool)))
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.UpdateLoginProfile(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return resourceAlicloudRamLoginProfileRead(d, meta)
}

func resourceAlicloudRamLoginProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamLoginProfile(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	profile := object.LoginProfile
	d.Set("user_name", profile.UserName)
	d.Set("mfa_bind_required", profile.MFABindRequired)
	d.Set("password_reset_required", profile.PasswordResetRequired)
	return nil
}

func resourceAlicloudRamLoginProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateDeleteLoginProfileRequest()
	request.UserName = d.Id()

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DeleteLoginProfile(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return WrapError(ramService.WaitForRamLoginProfile(d.Id(), Deleted, DefaultTimeout))

}
