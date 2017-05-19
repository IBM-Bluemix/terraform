package ibmcloud

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	slsession "github.com/softlayer/softlayer-go/session"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/api/account/accountv2"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	bxsession "github.com/IBM-Bluemix/bluemix-go/session"
)

//SoftlayerRestEndpoint rest endpoint of SoftLayer
const SoftlayerRestEndpoint = "https://api.softlayer.com/rest/v3"

//Config stores user provider input
type Config struct {
	//BluemixAPIKey is the Bluemix api key
	BluemixAPIKey string
	//Bluemix region
	Region string
	//Bluemix API timeout
	BluemixTimeout time.Duration

	//Softlayer end point url
	SoftLayerEndpointURL string

	//Softlayer API timeout
	SoftLayerTimeout time.Duration

	// Softlayer User Name
	SoftLayerUserName string

	// Softlayer API Key
	SoftLayerAPIKey string

	//SkipServiceConfig is a set of services whose configuration is to be skipped. Valid values could be bluemix, softlayer etc
	SkipServiceConfig *schema.Set

	//Retry Count for API calls
	//Unexposed in the schema at this point as they are used only during session creation for a few calls
	//When sdk implements it we an expose them for expected behaviour
	//https://github.com/softlayer/softlayer-go/issues/41
	RetryCount int
	//Constant Retry Delay for API calls
	RetryDelay time.Duration
}

//Session stores the information required for communication with the SoftLayer and Bluemix API
type Session struct {
	// SoftLayerSesssion is the the SoftLayer session used to connect to the SoftLayer API
	SoftLayerSession *slsession.Session

	// BluemixSession is the the Bluemix session used to connect to the Bluemix API
	BluemixSession *bxsession.Session
}

// ClientSession  contains  Bluemix/SoftLayer session and clients
type ClientSession interface {
	SoftLayerSession() *slsession.Session
	BluemixSession() *bxsession.Session

	ClusterClient() k8sclusterv1.Clusters
	ClusterWorkerClient() k8sclusterv1.Workers
	ClusterSubnetClient() k8sclusterv1.Subnets
	ClusterWebHooksClient() k8sclusterv1.Webhooks

	CloudFoundryOrgClient() cfv2.Organizations
	CloudFoundryServiceInstanceClient() cfv2.ServiceInstances
	CloudFoundryServicePlanClient() cfv2.ServicePlans
	CloudFoundryServiceKeyClient() cfv2.ServiceKeys
	CloudFoundryServiceOfferingClient() cfv2.ServiceOfferings
	CloudFoundrySpaceClient() cfv2.Spaces
	CloudFoundrySpaceQuotaClient() cfv2.SpaceQuotas

	BluemixAcccountClient() accountv2.Accounts
}

type clientSession struct {
	session *Session

	csClient  k8sclusterv1.Clusters
	csWorker  k8sclusterv1.Workers
	csSubnet  k8sclusterv1.Subnets
	csWebHook k8sclusterv1.Webhooks

	cfOrgClient              cfv2.Organizations
	cfServiceInstanceClient  cfv2.ServiceInstances
	cfSpaceClient            cfv2.Spaces
	cfSpaceQuotaClient       cfv2.SpaceQuotas
	cfServicePlanClient      cfv2.ServicePlans
	cfServiceKeysClient      cfv2.ServiceKeys
	cfServiceOfferingsClient cfv2.ServiceOfferings

	bluemixAccountClient accountv2.Accounts
}

// SoftLayerSession providers SoftLayer Session
func (sess clientSession) SoftLayerSession() *slsession.Session {
	return sess.session.SoftLayerSession
}

// CloudFoundryOrgClient providers Cloud Foundary org APIs
func (sess clientSession) CloudFoundryOrgClient() cfv2.Organizations {
	return sess.cfOrgClient
}

// CloudFoundrySpaceClient providers Cloud Foundary space APIs
func (sess clientSession) CloudFoundrySpaceClient() cfv2.Spaces {
	return sess.cfSpaceClient
}

// CloudFoundrySpaceQuotaClient providers Cloud Foundary space quota APIs
func (sess clientSession) CloudFoundrySpaceQuotaClient() cfv2.SpaceQuotas {
	return sess.cfSpaceQuotaClient
}

// CloudFoundryServiceInstanceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceInstanceClient() cfv2.ServiceInstances {
	return sess.cfServiceInstanceClient
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServicePlanClient() cfv2.ServicePlans {
	return sess.cfServicePlanClient
}

// CloudFoundryServiceKeyClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceKeyClient() cfv2.ServiceKeys {
	return sess.cfServiceKeysClient
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceOfferingClient() cfv2.ServiceOfferings {
	return sess.cfServiceOfferingsClient
}

// BluemixAcccountClient providers Bluemix Account APIs
func (sess clientSession) BluemixAcccountClient() accountv2.Accounts {
	return sess.bluemixAccountClient
}

// ClusterClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterClient() k8sclusterv1.Clusters {
	return sess.csClient
}

// ClusterWorkerClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWorkerClient() k8sclusterv1.Workers {
	return sess.csWorker
}

// ClusterSubnetClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterSubnetClient() k8sclusterv1.Subnets {
	return sess.csSubnet
}

// ClusterWebHooksClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWebHooksClient() k8sclusterv1.Webhooks {
	return sess.csWebHook
}

// BluemixSession to provide the Bluemix Session
func (sess clientSession) BluemixSession() *bxsession.Session {
	return sess.session.BluemixSession
}

// ClientSession configures and returns a fully initialized ClientSession
func (c *Config) ClientSession() (interface{}, error) {

	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	session := clientSession{
		session: sess,
	}

	if sess.BluemixSession == nil {
		log.Println("Skipping Bluemix Clients configuration")
		return session, nil
	}

	cfClient, err := cfv2.New(sess.BluemixSession)

	if err != nil {
		return nil, err
	}

	orgAPI := cfClient.Organizations()
	spaceAPI := cfClient.Spaces()
	serviceInstanceAPI := cfClient.ServiceInstances()
	servicePlanAPI := cfClient.ServicePlans()
	serviceKeysAPI := cfClient.ServiceKeys()
	serviceOfferringAPI := cfClient.ServiceOfferings()

	accClient, err := accountv2.New(sess.BluemixSession)
	if err != nil {
		return nil, err
	}
	accountAPI := accClient.Accounts()

	skipClusterConfig := c.SkipServiceConfig.Contains("cluster")

	if !skipClusterConfig {
		clusterClient, err := k8sclusterv1.New(sess.BluemixSession)
		if err != nil {
			if apiErr, ok := err.(bmxerror.Error); ok {
				if apiErr.Code() == endpoints.ErrCodeServiceEndpoint {
					return nil, fmt.Errorf(`Cluster service doesn't exist in the region %q.\nTo remediate the problem please skip the cluster service configuration by specifying "cluster" in skip_service_configuration in the provider block`, c.Region)

				}
			}
			return nil, err
		}
		clustersAPI := clusterClient.Clusters()
		clusterWorkerAPI := clusterClient.Workers()
		clusterSubnetsAPI := clusterClient.Subnets()
		clusterWebhookAPI := clusterClient.WebHooks()

		session.csClient = clustersAPI
		session.csSubnet = clusterSubnetsAPI
		session.csWorker = clusterWorkerAPI
		session.csWebHook = clusterWebhookAPI

	} else {
		log.Println("Skipping cluster configuration")
	}

	session.cfOrgClient = orgAPI
	session.cfServiceInstanceClient = serviceInstanceAPI
	session.cfServiceKeysClient = serviceKeysAPI
	session.cfServicePlanClient = servicePlanAPI
	session.cfServiceOfferingsClient = serviceOfferringAPI
	session.cfSpaceClient = spaceAPI
	session.bluemixAccountClient = accountAPI

	return session, nil
}

func newSession(c *Config) (*Session, error) {
	ibmcloudSession := &Session{}
	skipBluemix, skipSoftLayer := c.SkipServiceConfig.Contains("bluemix"), c.SkipServiceConfig.Contains("softlayer")

	if !skipSoftLayer {
		log.Println("Configuring SoftLayer Session ")
		if c.SoftLayerUserName == "" || c.SoftLayerAPIKey == "" {
			return nil, errors.New("softlayer_username and softlayer_api_key must be provided. Please see the documentation on how to configure them")
		}
		softlayerSession := &slsession.Session{
			Endpoint: c.SoftLayerEndpointURL,
			Timeout:  c.SoftLayerTimeout,
			UserName: c.SoftLayerUserName,
			APIKey:   c.SoftLayerAPIKey,
			Debug:    os.Getenv("TF_LOG") != "",
		}
		ibmcloudSession.SoftLayerSession = softlayerSession
	}
	if !skipBluemix {
		log.Println("Configuring Bluemix Session")
		if c.BluemixAPIKey == "" {
			return nil, errors.New("bluemix_api_key must be provided. Please see the documentation on how to configure it")
		}
		var sess *bxsession.Session
		bmxConfig := &bluemix.Config{
			BluemixAPIKey: c.BluemixAPIKey,
			Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
		}
		sess, err := bxsession.New(bmxConfig)
		if err != nil {
			return nil, err
		}
		ibmcloudSession.BluemixSession = sess
	}
	return ibmcloudSession, nil
}
