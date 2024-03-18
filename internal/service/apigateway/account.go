// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	awstypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_api_gateway_account")
func ResourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAccountUpdate,
		ReadWithoutTimeout:   resourceAccountRead,
		UpdateWithoutTimeout: resourceAccountUpdate,
		DeleteWithoutTimeout: schema.NoopContext,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"api_key_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudwatch_role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			"features": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"throttle_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"burst_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rate_limit": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	input := &apigateway.UpdateAccountInput{}

	if v, ok := d.GetOk("cloudwatch_role_arn"); ok {
		input.PatchOperations = []awstypes.PatchOperation{{
			Op:    awstypes.OpReplace,
			Path:  aws.String("/cloudwatchRoleArn"),
			Value: aws.String(v.(string)),
		}}
	} else {
		input.PatchOperations = []awstypes.PatchOperation{{
			Op:    awstypes.OpReplace,
			Path:  aws.String("/cloudwatchRoleArn"),
			Value: aws.String(""),
		}}
	}

	_, err := tfresource.RetryWhen(ctx, propagationTimeout,
		func() (interface{}, error) {
			return conn.UpdateAccount(ctx, input)
		},
		func(err error) (bool, error) {
			if errs.IsAErrorMessageContains[*awstypes.BadRequestException](err, "The role ARN does not have required permissions") {
				return true, err
			}

			if errs.IsAErrorMessageContains[*awstypes.BadRequestException](err, "API Gateway could not successfully write to CloudWatch Logs using the ARN specified") {
				return true, err
			}

			return false, err
		},
	)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating API Gateway Account: %s", err)
	}

	if d.IsNewResource() {
		d.SetId("api-gateway-account")
	}

	return append(diags, resourceAccountRead(ctx, d, meta)...)
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	account, err := conn.GetAccount(ctx, &apigateway.GetAccountInput{})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway Account: %s", err)
	}

	d.Set("api_key_version", account.ApiKeyVersion)
	d.Set("cloudwatch_role_arn", account.CloudwatchRoleArn)
	d.Set("features", account.Features)
	if err := d.Set("throttle_settings", flattenThrottleSettings(account.ThrottleSettings)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting throttle_settings: %s", err)
	}

	return diags
}
