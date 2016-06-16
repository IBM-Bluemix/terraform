package softlayer

import (
	"log"

	slclient "github.com/TheWeatherCompany/softlayer-go/client"
	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
)

type Config struct {
	Username string
	ApiKey   string
}

type Client struct {
	accountService      softlayer.SoftLayer_Account_Service
	virtualGuestService softlayer.SoftLayer_Virtual_Guest_Service
	sshKeyService       softlayer.SoftLayer_Security_Ssh_Key_Service
	productOrderService softlayer.SoftLayer_Product_Order_Service
	billingItemService  softlayer.SoftLayer_Billing_Item_Service
}

func (c *Config) Client() (*Client, error) {
	slc := slclient.NewSoftLayerClient(c.Username, c.ApiKey)

	accountService, err := slc.GetSoftLayer_Account_Service()
	if err != nil {
		return nil, err
	}

	virtualGuestService, err := slc.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
		return nil, err
	}

	sshKeyService, err := slc.GetSoftLayer_Security_Ssh_Key_Service()
	if err != nil {
		return nil, err
	}

	productOrderService, err := slc.GetSoftLayer_Product_Order_Service()
	if err != nil {
		return nil, err
	}

	billingItemService, err := slc.GetSoftLayer_Billing_Item_Service()
	if err != nil {
		return nil, err
	}

	client := &Client{
		accountService:      accountService,
		virtualGuestService: virtualGuestService,
		sshKeyService:       sshKeyService,
		productOrderService: productOrderService,
		billingItemService:  billingItemService,
	}

	log.Println("[INFO] Created SoftLayer client")

	return client, nil
}
