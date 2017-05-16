package vsphere

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi"
	"fmt"
	"github.com/vmware/govmomi/find"
	"context"
	"github.com/vmware/govmomi/govc/host/esxcli"
)

type datastore struct {
	name string
	storeType string
	datacenter string
	host string
	volumeIqn string
}

type StorageDevice struct {
	Device string
	IQN string
}

func resourceVSphereDatastore() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDatastoreCreate,
		Read:   resourceVSphereDatastoreRead,
		Update: resourceVSphereDatastoreUpdate,
		Delete: resourceVSphereDatastoreDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:	  schema.TypeString,
				Required: true,
			},
			"type": {
				Type: 	  schema.TypeString,
				Required: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volumeIqn": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVSphereDatastoreCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] creating datastore: %#v", d)
	client := meta.(*govmomi.Client)

	ds := datastore{}

	if v, ok := d.GetOk("name"); ok {
		ds.name = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		ds.storeType = v.(string)
	}

	err := createDatastore(client, &ds)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("[%v] %v/%v", ds))
	log.Printf("[INFO] Created datastore: %s", ds.name)

	return resourceVSphereDatastoreRead(d, meta)
}

func createDatastore(client *govmomi.Client, ds *datastore) error {
	return nil
}

func resourceVSphereDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereDatastoreUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereDatastoreDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func rescanAllHBAs(client *govmomi.Client, ds *datastore) error {
	finder := find.NewFinder(client.Client, true)

	dc, err := finder.Datacenter(context.TODO(), ds.datacenter)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	finder = finder.SetDatacenter(dc)

	hs, err := finder.HostSystem(context.TODO(), ds.host)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}

	ss, err := hs.ConfigManager().StorageSystem(context.TODO())
	if err != nil {
		return fmt.Errorf("error %s", err)
	}

	scanErr := ss.RescanAllHba(context.TODO())
	if scanErr != nil {
		return fmt.Errorf("error %s", scanErr)
	}

	return nil
}

func getVMwareVolume(client *govmomi.Client, ds *datastore) error {

	// Get devices

	// Find the right device

	return nil
}

func getiSCSITargets(client *govmomi.Client, ds *datastore) ([]StorageDevice, error) {
	finder := find.NewFinder(client, true)
	dc, err := finder.Datacenter(context.TODO(), ds.datacenter)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}
	finder = finder.SetDatacenter(dc)

	hs, err := finder.HostSystem(context.TODO(), ds.host)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}

	esx, err := esxcli.NewExecutor(hs.Client(), hs)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}

	res, err := esx.Run([]string{"storage", "core", "path", "list"})
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}

	var storageDevices []StorageDevice
}