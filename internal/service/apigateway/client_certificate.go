// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	awstypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_api_gateway_client_certificate", name="Client Certificate")
// @Tags(identifierAttribute="arn")
func ResourceClientCertificate() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClientCertificateCreate,
		ReadWithoutTimeout:   resourceClientCertificateRead,
		UpdateWithoutTimeout: resourceClientCertificateUpdate,
		DeleteWithoutTimeout: resourceClientCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pem_encoded_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClientCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	input := &apigateway.GenerateClientCertificateInput{
		Tags: getTagsIn(ctx),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	output, err := conn.GenerateClientCertificate(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating API Gateway Client Certificate: %s", err)
	}

	d.SetId(aws.ToString(output.ClientCertificateId))

	return append(diags, resourceClientCertificateRead(ctx, d, meta)...)
}

func resourceClientCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	cert, err := FindClientCertificateByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] API Gateway Client Certificate (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway Client Certificate (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "apigateway",
		Region:    meta.(*conns.AWSClient).Region,
		Resource:  fmt.Sprintf("/clientcertificates/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("created_date", cert.CreatedDate.String())
	d.Set("description", cert.Description)
	d.Set("expiration_date", cert.ExpirationDate.String())
	d.Set("pem_encoded_certificate", cert.PemEncodedCertificate)

	setTagsOut(ctx, cert.Tags)

	return diags
}

func resourceClientCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &apigateway.UpdateClientCertificateInput{
			ClientCertificateId: aws.String(d.Id()),
			PatchOperations: []awstypes.PatchOperation{
				{
					Op:    awstypes.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String(d.Get("description").(string)),
				},
			},
		}

		_, err := conn.UpdateClientCertificate(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating API Gateway Client Certificate (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceClientCertificateRead(ctx, d, meta)...)
}

func resourceClientCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayClient(ctx)

	log.Printf("[DEBUG] Deleting API Gateway Client Certificate: %s", d.Id())
	_, err := conn.DeleteClientCertificate(ctx, &apigateway.DeleteClientCertificateInput{
		ClientCertificateId: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting API Gateway Client Certificate (%s): %s", d.Id(), err)
	}

	return diags
}

func FindClientCertificateByID(ctx context.Context, conn *apigateway.Client, id string) (*apigateway.GetClientCertificateOutput, error) {
	input := &apigateway.GetClientCertificateInput{
		ClientCertificateId: aws.String(id),
	}

	output, err := conn.GetClientCertificate(ctx, input)

	if errs.IsA[*awstypes.NotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
