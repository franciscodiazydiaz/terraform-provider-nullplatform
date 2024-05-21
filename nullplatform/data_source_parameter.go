package nullplatform

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceParameter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceParameterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "A system-wide unique ID representing the resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Definition name of the variable.",
			},
			"nrn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NRN of the application to which the parameter belongs to.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Possible values: [`environment`, `file`]",
			},
			"encoding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Possible values: [`plaintext`, `base64`]",
			},
			"variable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the environment variable. Required when `type = environment`.",
			},
			"destination_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path for file. Required when `type = file`.",
			},
			"secret": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "`true` if the value is a secret, `false` otherwise",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "`true` if the value is a secret, `false` otherwise",
			},
			/*
				"values": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
					Description: "List of unique IDs representing the values",
				},
			*/
		},
	}
}

func dataSourceParameterRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	nullOps := m.(NullOps)

	param, err := nullOps.GetParameter(strconv.Itoa(d.Get("id").(int)))

	/*
		parameterList, err := nullOps.GetParameterList(param.Nrn)
		if err != nil {
			return nil, err
		}

		paramRes, paramExists := parameterExists(parameterList, param)
		if paramExists && param.ImportIfCreated {
			log.Printf("[DEBUG] Parameter with Name: %s and Variable: %s already exists, importing ID: %d", paramRes.Name, paramRes.Variable, paramRes.Id)
			return paramRes, nil
		}
	*/

	err = d.Set("name", param.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("nrn", param.Nrn)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("type", param.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("encoding", param.Encoding)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("variable", param.Variable)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("destination_path", param.DestinationPath)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("secret", param.Secret)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("read_only", param.ReadOnly)
	if err != nil {
		return diag.FromErr(err)
	}

	//fmt.Printf("ResourceData: %+v\n", d)

	// We don't have a unique ID for this data resource so we create one using a
	// timestamp format. I've seen people use a hash of the returned API data as
	// a unique key.
	//
	// NOTE:
	// That hashcode helper is no longer available! It has been moved into an
	// internal directory meaning it's not supposed to be consumed.
	//
	// Reference:
	// https://github.com/hashicorp/terraform-plugin-sdk/blob/master/internal/helper/hashcode/hashcode.go
	//
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
