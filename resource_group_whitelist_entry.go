package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/visit1985/atlasgo/services/group"
	"time"
)

func resourceGroupWhitelistEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupWhitelistEntryCreate,
		Read:   resourceGroupWhitelistEntryRead,
		Update: resourceGroupWhitelistEntryUpdate,
		Delete: resourceGroupWhitelistEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"ip_address"},
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_after_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_address": &schema.Schema{
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"cidr_block"},
			},
		},
	}
}

func resourceGroupWhitelistEntryCreate(d *schema.ResourceData, m interface{}) error {
	var id string
	client := m.(*group.Group)
	params := group.WhitelistEntry{}

	if v, ok := d.GetOk("cidr_block"); ok {
		params.CidrBlock = v.(string)
		id = v.(string)
	}

	if v, ok := d.GetOk("comment"); ok {
		params.Comment = v.(string)
	}

	if v, ok := d.GetOk("delete_after_days"); ok {
		t := time.Now().AddDate(0, 0, v.(int)+1).Truncate(24 * time.Hour)
		params.DeleteAfterDate = &t
	}

	if v, ok := d.GetOk("ip_address"); ok {
		params.IpAddress = v.(string)
		id = v.(string)
	}

	if err := client.SetWhitelistEntry(&group.SetWhitelistEntryInput{params}); err != nil {
		return fmt.Errorf("Error creating Whitelist Entry: %s", err)
	}

	d.SetId(id)
	return resourceGroupWhitelistEntryRead(d, m)
}

func resourceGroupWhitelistEntryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*group.Group)
	res, err := client.GetWhitelistEntry(d.Id())

	if err != nil {
		d.SetId("")
		//TODO: test error type and return if != NotFoundError
		return nil
	}

	if res.IpAddress != "" {
		_ = d.Set("ip_address", res.IpAddress)
	} else {
		_ = d.Set("cidr_block", res.CidrBlock)
	}
	_ = d.Set("comment", res.Comment)

	if res.DeleteAfterDate != nil {
		_ = d.Set("delete_after_days", res.DeleteAfterDate.Sub(time.Now()).Hours()/24)
	}

	return nil
}

func resourceGroupWhitelistEntryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*group.Group)
	params := group.WhitelistEntry{}

	if v, ok := d.GetOk("cidr_block"); ok {
		params.CidrBlock = v.(string)
	}

	if v, ok := d.GetOk("comment"); ok {
		params.Comment = v.(string)
	}

	if v, ok := d.GetOk("delete_after_days"); ok {
		t := time.Now().AddDate(0, 0, v.(int)+1).Truncate(24 * time.Hour)
		params.DeleteAfterDate = &t
	}

	if v, ok := d.GetOk("ip_address"); ok {
		params.IpAddress = v.(string)
	}

	if err := client.SetWhitelistEntry(&group.SetWhitelistEntryInput{params}); err != nil {
		return fmt.Errorf("Error creating Whitelist Entry: %s", err)
	}

	return resourceGroupWhitelistEntryRead(d, m)
}

func resourceGroupWhitelistEntryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*group.Group)

	err := client.DeleteWhitelistEntry(d.Id())

	if err != nil {
		//TODO: test error type and return if != NotFoundError
		return nil
	}

	return nil
}
