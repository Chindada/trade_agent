// Package mqhandler package mqhandler
package mqhandler

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
	"trade_agent/pkg/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func getMQClient(mqConf config.MQTT) (client mqtt.Client, err error) {
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt((10000)))
	if err != nil {
		return client, err
	}
	random := "-" + randomBigInt.String()
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("ssl://%s:%s", mqConf.Host, mqConf.Port)).
		SetClientID(mqConf.ClientID + random).
		SetTLSConfig(genTLSConfig(mqConf)).
		SetUsername(mqConf.User).
		SetPassword(mqConf.Passwd).
		SetOnConnectHandler(onConnect(mqConf)).
		SetConnectionLostHandler(onLost(mqConf))
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return client, err
	}
	return c, err
}

// genTLSConfig genTLSConfig
func genTLSConfig(mqConf config.MQTT) *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(mqConf.CAPath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	cert, err := tls.LoadX509KeyPair(mqConf.CertPath, mqConf.KeyPath)
	if err != nil {
		log.Get().Panic(err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Get().Panic(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
	}
}

func onConnect(mqConf config.MQTT) func(mqtt.Client) {
	return func(mqtt.Client) {
		log.Get().Infof("MQTT Broker on %s:%s Connected", mqConf.Host, mqConf.Port)
	}
}

func onLost(mqConf config.MQTT) func(mqtt.Client, error) {
	return func(mqtt.Client, error) {
		log.Get().Errorf("MQTT Broker on %s:%s Disconnected", mqConf.Host, mqConf.Port)
		mqConf := config.GetMQTTConfig()
		for {
			if utils.CheckPortIsOpen(mqConf.Host, mqConf.Port) {
				break
			}
			time.Sleep(time.Second)
		}

		// get new mq connection
		newClient, err := getMQClient(mqConf)
		if err != nil {
			log.Get().Panic(err)
		}

		// force sinopac mq srv re-connect mqtt broker
		sinoErr := sinopacapi.Get().AskSinpacMQSRVConnectMQ(mqConf)
		if sinoErr != nil {
			log.Get().Panic(sinoErr)
		}

		Get().client.Disconnect(0)
		Get().client = newClient

		Get().lock.Lock()
		// copy old cb map
		oldCBMap := Get().callbackMap
		oldOnceMap := Get().onceMap
		// empty old map
		Get().callbackMap = make(map[string]MQCallback)
		Get().onceMap = make(map[string]bool)
		Get().lock.Unlock()

		for topic, cb := range oldCBMap {
			err := Get().Sub(MQSubBody{
				MQTopic:  MQTopic(topic),
				Once:     oldOnceMap[topic],
				Callback: cb,
			})
			if err != nil {
				log.Get().Panic("MQTT Resubscribe Fail")
			}
		}

		// get target from cache send event to resubscribe target
		targetArr := cache.GetCache().GetTargets()
		eventbus.Get().Pub(eventbus.TopicSubscribeTargets(), targetArr)
		log.Get().Warn("Resubscribe All Done")
	}
}
