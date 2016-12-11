package mackerel

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelExternalMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelExternalMonitorCreate,
		Read:   resourceMackerelExternalMonitorRead,
		Update: resourceMackerelExternalMonitorUpdate,
		Delete: resourceMackerelExternalMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"response_time_warning": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"response_time_critical": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"response_time_duration": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				//ValidateFunc: validateDurationTime 1 minute - 10 minute
			},
			"contains_string": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_check_attempts": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  1,
				//ValidateFunc: validateDurationTime 1 - 10
			},
			"certification_expiration_warning": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"certification_expiration_critical": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_mute": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceMackerelExternalMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := getMackerelExternalMonitorInput(d)
	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelExternalMonitorRead(d, meta)
}

func resourceMackerelExternalMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "external" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorExternalHTTP)
			d.Set("id", mon.MonitorID())
			d.Set("name", mon.MonitorName())
			d.Set("url", mon.URL)
			d.Set("service", mon.Service)
			d.Set("notification_interval", mon.NotificationInterval)
			d.Set("response_time_duration", mon.ResponseTimeDuration)
			d.Set("response_time_warning", mon.ResponseTimeWarning)
			d.Set("response_time_critical", mon.ResponseTimeCritical)
			d.Set("contains_string", mon.ContainsString)
			d.Set("max_check_attempts", mon.MaxCheckAttempts)
			d.Set("certification_expiration_warning", mon.CertificationExpirationWarning)
			d.Set("certification_expiration_critical", mon.CertificationExpirationCritical)
			d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelExternalMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := getMackerelExternalMonitorInput(d)
	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelExternalMonitorRead(d, meta)
}

func resourceMackerelExternalMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}

func getMackerelExternalMonitorInput(d *schema.ResourceData) *mackerel.MonitorExternalHTTP {
	input := &mackerel.MonitorExternalHTTP{
		Type: "external",
		Name: d.Get("name").(string),
		URL:  d.Get("url").(string),
	}

	if v, ok := d.GetOk("service"); ok {
		input.Service = v.(string)
	}
	if v, ok := d.GetOk("notification_interval"); ok {
		input.NotificationInterval = uint64(v.(int))
	}
	if v, ok := d.GetOk("response_time_duration"); ok {
		input.ResponseTimeDuration = v.(float64)
	}
	if v, ok := d.GetOk("response_time_warning"); ok {
		input.ResponseTimeWarning = v.(float64)
	}
	if v, ok := d.GetOk("response_time_critical"); ok {
		input.ResponseTimeCritical = v.(float64)
	}
	if v, ok := d.GetOk("contains_string"); ok {
		input.ContainsString = v.(string)
	}
	if v, ok := d.GetOk("max_check_attempts"); ok {
		input.MaxCheckAttempts = v.(float64)
	}
	if v, ok := d.GetOk("certification_expiration_warning"); ok {
		input.CertificationExpirationWarning = uint64(v.(int))
	}
	if v, ok := d.GetOk("certification_expiration_critical"); ok {
		input.CertificationExpirationCritical = uint64(v.(int))
	}
	if v, ok := d.GetOk("is_mute"); ok {
		input.IsMute = v.(bool)
	}

	return input
}
