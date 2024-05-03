package nullplatform

import (
	"context"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceParameter() *schema.Resource {
	return &schema.Resource{
		Create: ParameterCreate,
		Read:   ParameterRead,
		Update: ParameterUpdate,
		Delete: ParameterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nrn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "environment", //  Possible values: [environment, file]
				Optional: true,
				ForceNew: true,
			},
			"encoding": {
				Type:     schema.TypeString,
				Default:  "plaintext", //  Possible values: [plaintext, base64]
				Optional: true,
				ForceNew: true,
			},
			"variable": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_path": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"secret": {
				Type:     schema.TypeBool,
				Default:  "false",
				Optional: true,
				ForceNew: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Default:  "false",
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func ParameterCreate(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)

	newParameter := &Parameter{
		Name:            d.Get("name").(string),
		Nrn:             d.Get("nrn").(string),
		Type:            d.Get("type").(string),
		Encoding:        d.Get("encoding").(string),
		Variable:        d.Get("variable").(string),
		DestinationPath: d.Get("destination_path").(string),
		Secret:          d.Get("secret").(bool),
		ReadOnly:        d.Get("read_only").(bool),
	}

	param, err := nullOps.CreateParameter(newParameter)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(param.Id))

	return ParameterRead(d, m)
}

func ParameterRead(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)
	parameterId := d.Id()

	param, err := nullOps.GetParameter(parameterId)
	if err != nil {
		// FIXME: Validate if error == 404
		if !d.IsNewResource() {
			log.Printf("[WARN] Parameter ID %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("name", param.Name); err != nil {
		return err
	}

	if err := d.Set("nrn", param.Nrn); err != nil {
		return err
	}

	if err := d.Set("type", param.Type); err != nil {
		return err
	}

	if err := d.Set("encoding", param.Encoding); err != nil {
		return err
	}

	if err := d.Set("variable", param.Variable); err != nil {
		return err
	}

	if err := d.Set("destination_path", param.DestinationPath); err != nil {
		return err
	}

	if err := d.Set("secret", param.Secret); err != nil {
		return err
	}

	if err := d.Set("read_only", param.ReadOnly); err != nil {
		return err
	}

	return nil
}

func ParameterUpdate(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)

	parameterId := d.Id()

	param := &Parameter{}

	if d.HasChange("name") {
		param.Name = d.Get("name").(string)
	}

	if d.HasChange("nrn") {
		param.Nrn = d.Get("nrn").(string)
	}

	if d.HasChange("type") {
		param.Type = d.Get("type").(string)
	}

	if d.HasChange("encoding") {
		param.Encoding = d.Get("encoding").(string)
	}

	if d.HasChange("variable") {
		param.Variable = d.Get("variable").(string)
	}

	if d.HasChange("destination_path") {
		param.DestinationPath = d.Get("destination_path").(string)
	}

	if d.HasChange("secret") {
		param.Secret = d.Get("secret").(bool)
	}

	if d.HasChange("read_only") {
		param.ReadOnly = d.Get("read_only").(bool)
	}

	if !reflect.DeepEqual(*param, Parameter{}) {
		err := nullOps.PatchParameter(parameterId, param)
		if err != nil {
			return err
		}
	}

	return nil //ParameterRead(d, m)
}

func ParameterDelete(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)
	parameterId := d.Id()

	err := nullOps.DeleteParameter(parameterId)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
