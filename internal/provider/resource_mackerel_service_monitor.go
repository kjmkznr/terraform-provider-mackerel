package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelServiceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceMonitorCreate,
		Read:   resourceMackerelServiceMonitorRead,
		Update: resourceMackerelServiceMonitorUpdate,
		Delete: resourceMackerelServiceMonitorDelete,
		Exists: resourceMackerelServiceMonitorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"warning": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"critical": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"missing_duration_warning": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"missing_duration_critical": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_check_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      1,
			},
			"is_mute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceMackerelServiceMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                    "service",
		Name:                    d.Get("name").(string),
		Service:                 d.Get("service").(string),
		Duration:                uint64(d.Get("duration").(int)),
		Metric:                  d.Get("metric").(string),
		Operator:                d.Get("operator").(string),
		Warning:                 pfloat64(d.Get("warning").(float64)),
		Critical:                pfloat64(d.Get("critical").(float64)),
		MissingDurationWarning:  uint64(d.Get("missing_duration_warning").(int)),
		MissingDurationCritical: uint64(d.Get("missing_duration_critical").(int)),
		NotificationInterval:    uint64(d.Get("notification_interval").(int)),
		MaxCheckAttempts:        uint64(d.Get("max_check_attempts").(int)),
		IsMute:                  d.Get("is_mute").(bool),
	}

	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "service" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorServiceMetric)
			_ = d.Set("id", mon.ID)
			_ = d.Set("name", mon.Name)
			_ = d.Set("service", mon.Service)
			_ = d.Set("duration", mon.Duration)
			_ = d.Set("metric", mon.Metric)
			_ = d.Set("operator", mon.Operator)
			_ = d.Set("warning", mon.Warning)
			_ = d.Set("critical", mon.Critical)
			_ = d.Set("missing_duration_warning", mon.MissingDurationWarning)
			_ = d.Set("missing_duration_critical", mon.MissingDurationCritical)
			_ = d.Set("notification_interval", mon.NotificationInterval)
			_ = d.Set("max_check_attempts", mon.MaxCheckAttempts)
			_ = d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelServiceMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                    "service",
		Name:                    d.Get("name").(string),
		Service:                 d.Get("service").(string),
		Duration:                uint64(d.Get("duration").(int)),
		Metric:                  d.Get("metric").(string),
		Operator:                d.Get("operator").(string),
		Warning:                 pfloat64(d.Get("warning").(float64)),
		Critical:                pfloat64(d.Get("critical").(float64)),
		MissingDurationWarning:  uint64(d.Get("missing_duration_warning").(int)),
		MissingDurationCritical: uint64(d.Get("missing_duration_critical").(int)),
		NotificationInterval:    uint64(d.Get("notification_interval").(int)),
		MaxCheckAttempts:        uint64(d.Get("max_check_attempts").(int)),
		IsMute:                  d.Get("is_mute").(bool),
	}

	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*mackerel.Client)
	monitors, err := client.FindMonitors()
	if err != nil {
		return false, err
	}

	for _, m := range monitors {
		if m.MonitorType() == "service" && m.MonitorID() == d.Id() {
			return true, nil
		}
	}

	return false, nil
}

func resourceMackerelServiceMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}
