package searching

import (
	"github.com/sirupsen/logrus"
)

func logElasticError(url string, sniff bool, username, password string, maxRetry int, step string, err error) {
	logrus.WithFields(logrus.Fields{
		"url":      url,
		"sniff":    sniff,
		"username": username,
		"password": password,
		"maxRetry": maxRetry,
	}).Infof("ElasticSearch|"+step+" failed eturned with Error: %s", err.Error())

}
