package nullplatform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about the Service",
		ReadContext: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A system-wide unique ID for the Service.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Service name.",
			},
		},
	}
}

func dataSourceServiceRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	nullOps := m.(NullOps)

	s, err := nullOps.GetService(d.Get("id").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", s.Name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("id").(string))

	return nil
}
