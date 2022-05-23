package provider

import (
	"context"
	"fmt"
	"log"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareLogpullRetention() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceCloudflareLogpullRetentionSchema(),
		CreateContext: resourceCloudflareLogpullRetentionSet,
		ReadContext:   resourceCloudflareLogpullRetentionRead,
		UpdateContext: resourceCloudflareLogpullRetentionSet,
		DeleteContext: resourceCloudflareLogpullRetentionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudflareLogpullRetentionImport,
		},
	}
}

func resourceCloudflareLogpullRetentionSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)
	status := d.Get("enabled").(bool)

	_, err := client.SetLogpullRetentionFlag(ctx, zoneID, status)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting Logpull Retention for zone ID %q: %w", zoneID, err))
	}

	d.SetId(stringChecksum("logpull-retention/" + zoneID))

	return resourceCloudflareLogpullRetentionRead(ctx, d, meta)
}

func resourceCloudflareLogpullRetentionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	logpullConf, err := client.GetLogpullRetentionFlag(ctx, zoneID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Logpull Retention for zone ID %q: %w", zoneID, err))
	}

	d.Set("enabled", logpullConf.Flag)

	return nil
}

func resourceCloudflareLogpullRetentionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	_, err := client.SetLogpullRetentionFlag(ctx, zoneID, false)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting Logpull Retention for zone ID %q: %w", zoneID, err))
	}

	d.SetId("")

	return nil
}

func resourceCloudflareLogpullRetentionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	zoneID := d.Id()

	log.Printf("[DEBUG] Importing Cloudflare Logpull Retention option for zone ID: %s", zoneID)

	d.Set("zone_id", zoneID)
	d.SetId(stringChecksum("logpull-retention/" + zoneID))

	resourceCloudflareLogpullRetentionRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}