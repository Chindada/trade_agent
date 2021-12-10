// Package mqhandler package mqhandler
package mqhandler

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

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
		SetPassword(mqConf.Password).
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
		log.Get().Infof("MQTT Broker on %s:%s is connected", mqConf.Host, mqConf.Port)
	}
}

func onLost(mqConf config.MQTT) func(mqtt.Client, error) {
	return func(mqtt.Client, error) {
		log.Get().Panicf("MQTT Broker on %s:%s is disconnected", mqConf.Host, mqConf.Port)
	}
}
