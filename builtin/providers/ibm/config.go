package ibm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

var (
	errEmptySoftLayerCredentials = errors.New("softlayer_username and softlayer_api_key must be provided. Please see the documentation on how to configure them")
	errEmptyBluemixCredentials   = errors.New("bluemix_api_key must be provided. Please see the documentation on how to configure it")
)

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

// ClientSession ...
type ClientSession interface {
	SoftLayerSession() *slsession.Session
	BluemixSession() (*bxsession.Session, error)
	CSAPI() (k8sclusterv1.ClusterServiceAPI, error)
	CFAPI() (cfv2.CfServiceAPI, error)
	BluemixAcccountAPI() (accountv2.AccountServiceAPI, error)
}

type clientSession struct {
	session *Session

	csConfigErr  error
	csServiceAPI k8sclusterv1.ClusterServiceAPI

	cfConfigErr  error
	cfServiceAPI cfv2.CfServiceAPI

	accountConfigErr     error
	bmxAccountServiceAPI accountv2.AccountServiceAPI
}

// SoftLayerSession providers SoftLayer Session
func (sess clientSession) SoftLayerSession() *slsession.Session {
	return sess.session.SoftLayerSession
}

// CFAPI provides Cloud Foundry APIs ...
func (sess clientSession) CFAPI() (cfv2.CfServiceAPI, error) {
	return sess.cfServiceAPI, sess.cfConfigErr
}

// BluemixAcccountAPI ...
func (sess clientSession) BluemixAcccountAPI() (accountv2.AccountServiceAPI, error) {
	return sess.bmxAccountServiceAPI, sess.accountConfigErr
}

// CSAPI provides cluster APIs ...
func (sess clientSession) CSAPI() (k8sclusterv1.ClusterServiceAPI, error) {
	return sess.csServiceAPI, sess.csConfigErr
}

// BluemixSession to provide the Bluemix Session
func (sess clientSession) BluemixSession() (*bxsession.Session, error) {
	return sess.session.BluemixSession, sess.cfConfigErr
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
		//Can be nil only  if bluemix_api_key is not provided
		log.Println("Skipping Bluemix Clients configuration")
		session.csConfigErr = errEmptyBluemixCredentials
		session.cfConfigErr = errEmptyBluemixCredentials
		session.accountConfigErr = errEmptyBluemixCredentials
		return session, nil
	}

	cfAPI, err := cfv2.New(sess.BluemixSession)
	if err != nil {
		return nil, fmt.Errorf("Error occured while configuring Cloud Foundry API: %q", err)
	}
	session.cfServiceAPI = cfAPI

	accAPI, err := accountv2.New(sess.BluemixSession)

	if err != nil {
		return nil, fmt.Errorf("Error occured while configuring Bluemix Account Service: %q", err)
	}
	session.bmxAccountServiceAPI = accAPI

	clusterAPI, err := k8sclusterv1.New(sess.BluemixSession)
	var noSvcError error
	if err != nil {
		if apiErr, ok := err.(bmxerror.Error); ok {
			if apiErr.Code() == endpoints.ErrCodeServiceEndpoint {
				noSvcError = fmt.Errorf(`IBM Container Service for K8s cluster doesn't exist in the region %q`, c.Region)
				session.csConfigErr = noSvcError
			}
		}
		if noSvcError == nil {
			return nil, fmt.Errorf("Error occured while configuring IBM Container Service for K8s cluster: %q", err)
		}
	}
	session.csServiceAPI = clusterAPI
	return session, nil
}

func newSession(c *Config) (*Session, error) {
	ibmSession := &Session{}

	log.Println("Configuring SoftLayer Session ")
	softlayerSession := &slsession.Session{
		Endpoint: c.SoftLayerEndpointURL,
		Timeout:  c.SoftLayerTimeout,
		UserName: c.SoftLayerUserName,
		APIKey:   c.SoftLayerAPIKey,
		Debug:    os.Getenv("TF_LOG") != "",
	}
	ibmSession.SoftLayerSession = softlayerSession

	if c.BluemixAPIKey != "" {
		log.Println("Configuring Bluemix Session")
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
		ibmSession.BluemixSession = sess
	}

	return ibmSession, nil
}
