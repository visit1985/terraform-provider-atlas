package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/visit1985/atlasgo/services/group"
)

func Provider() *schema.Provider {
	descriptions := map[string]string{
		"group_id": "The group id for API operations. You can retrive this\n" +
			"from the 'Project > Settings' section of the Atlas console.",
	}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["group_id"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"atlas_group_whitelist_entry": resourceGroupWhitelistEntry(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	groupId := d.Get("group_id").(string)
	client := group.New(groupId)
	return client, client.Error
}
