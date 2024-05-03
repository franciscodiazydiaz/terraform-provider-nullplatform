package nullplatform

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceParameterValue() *schema.Resource {
	return &schema.Resource{
		Create: ParameterValueCreate,
		Read:   ParameterValueRead,
		Update: ParameterValueUpdate,
		Delete: ParameterValueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"parameter_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"origin_version": {
				Type:     schema.TypeInt,
				Optional: true,
				//Computed: true,
				Default: 0,
			},
			"nrn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ParameterValueCreate(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)

	// FIXME: This code is duplicated in Scope
	dimensionsMap := d.Get("dimensions").(map[string]any)
	// Convert the dimensions to a map[string]string
	dimensions := make(map[string]string)
	for key, value := range dimensionsMap {
		dimensions[key] = value.(string)
	}

	parameterId := d.Get("parameter_id").(int)

	newParameterValue := &ParameterValue{
		OriginVersion: d.Get("origin_version").(int),
		Nrn:           d.Get("nrn").(string),
		Value:         d.Get("value").(string),
		Dimensions:    dimensions,
	}

	paramValue, err := nullOps.CreateParameterValue(parameterId, newParameterValue)

	if err != nil {
		return err
	}

	paramValueId := generateParameterValueID(paramValue)
	d.SetId(paramValueId)

	return nil
}

func ParameterValueRead(d *schema.ResourceData, m any) error {
	var parameterValue *ParameterValue

	nullOps := m.(NullOps)
	parameterId := strconv.Itoa(d.Get("parameter_id").(int))
	parameterValueId := d.Id()

	param, err := nullOps.GetParameter(parameterId)
	if err != nil {
		// FIXME: Validate if error == 404
		/*if !d.IsNewResource() {
			log.Printf("[WARN] Parameter Value ID %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}*/
		return err
	}

	for _, item := range param.Values {
		// -------- DEBUG
		// Convert struct to JSON
		jsonData, err := json.Marshal(item)
		if err != nil {
			return err
		}
		// Print JSON string
		log.Println(string(jsonData))
		// -------- DEBUG
		log.Println("**********************", parameterValueId, strconv.Itoa(item.Id))

		if parameterValueId == generateParameterValueID(item) {
			parameterValue = item
			break
		}
	}

	if parameterValue == nil {
		log.Printf("[WARN] Cannot fetch Parameter Value ID %s", parameterValueId)
		return nil
	}

	/*if err := d.Set("origin_version", parameterValue.OriginVersion); err != nil {
		return err
	}*/

	if err := d.Set("nrn", parameterValue.Nrn); err != nil {
		return err
	}

	if err := d.Set("value", parameterValue.Value); err != nil {
		return err
	}

	if err := d.Set("dimensions", parameterValue.Dimensions); err != nil {
		return err
	}

	return nil
}

func ParameterValueUpdate(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)

	// FIXME: This code is duplicated in Scope
	dimensionsMap := d.Get("dimensions").(map[string]any)
	// Convert the dimensions to a map[string]string
	dimensions := make(map[string]string)
	for key, value := range dimensionsMap {
		dimensions[key] = value.(string)
	}

	parameterId := d.Get("parameter_id").(int)

	newParameterValue := &ParameterValue{}

	if d.HasChange("origin_version") {
		newParameterValue.OriginVersion = d.Get("origin_version").(int)
	}

	if d.HasChange("value") {
		newParameterValue.Value = d.Get("value").(string)
	}

	// The ID of the Parameter Value will change if other value is updated
	// Instead the NRN and Dimensions are composed to generate an ID
	if !reflect.DeepEqual(*newParameterValue, ParameterValue{}) {
		newParameterValue.Nrn = d.Get("nrn").(string)
		// Update the value means creating a new version of it
		paramValue, err := nullOps.CreateParameterValue(parameterId, newParameterValue)
		if err != nil {
			return err
		}
		// -------- DEBUG
		// Convert struct to JSON
		jsonData, err := json.Marshal(paramValue)
		if err != nil {
			return err
		}
		// Print JSON string
		log.Println("****************", string(jsonData))
		// -------- DEBUG
		//d.Set("new_id", paramValue.Id)
		paramValueId := generateParameterValueID(paramValue)
		d.SetId(paramValueId)
	}

	return nil
}

func ParameterValueDelete(d *schema.ResourceData, m any) error {
	nullOps := m.(NullOps)
	parameterId := strconv.Itoa(d.Get("parameter_id").(int))
	parameterValueId := d.Id()

	// FIXME: Most of this logic is duplicated in `ParameterValueRead`
	var parameterValue *ParameterValue

	param, err := nullOps.GetParameter(parameterId)
	if err != nil {
		// FIXME: Validate if error == 404
		/*if !d.IsNewResource() {
			log.Printf("[WARN] Parameter Value ID %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}*/
		return err
	}

	for _, item := range param.Values {
		// -------- DEBUG
		// Convert struct to JSON
		jsonData, err := json.Marshal(item)
		if err != nil {
			return err
		}
		// Print JSON string
		log.Println(string(jsonData))
		// -------- DEBUG
		log.Println("**********************", parameterValueId, strconv.Itoa(item.Id))

		if parameterValueId == generateParameterValueID(item) {
			parameterValue = item
			break
		}
	}

	if parameterValue == nil {
		log.Printf("[WARN] Cannot fetch Parameter Value ID %s", parameterValueId)
		return nil
	}

	err = nullOps.DeleteParameterValue(parameterId, strconv.Itoa(parameterValue.Id))
	if err != nil {
		// FIXME: Validate if error == 404
		log.Printf("[WARN] Parameter Value ID %s not found, removing from state", parameterValueId)
		d.SetId("")
		return nil
	}

	d.SetId("")

	return nil
}

func generateParameterValueID(value *ParameterValue) string {
	var concatenatedString string

	// Concatenate all key-value pairs from the map
	for key, value := range value.Dimensions {
		concatenatedString += key + ":" + value + ";"
	}

	concatenatedString += value.Nrn + ";"

	// Hash the concatenated string using SHA-256
	hash := sha256.New()
	hash.Write([]byte(concatenatedString))
	hashBytes := hash.Sum(nil)

	// Convert the hash bytes to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
