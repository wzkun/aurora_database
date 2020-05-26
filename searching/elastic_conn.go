package searching

import (
	"context"
	"time"

	"code.aliyun.com/bim_backend/zoogoer/gun/errors"
	"github.com/olivere/elastic"
	"github.com/wzkun/aurora_database/errstring"
	"github.com/wzkun/aurora_database/utils"
)

// MakeElasticConn function
func MakeElasticConn(url string, sniff bool, username, password string, maxRetry int) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, password),
		elastic.SetSniff(sniff),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(maxRetry),
	)

	if err != nil {
		logElasticError(url, sniff, username, password, maxRetry, "NewClient", err)
		return nil, errors.NewClientErr(nil, errstring.ESSetUpConnFail, "", "ElasticSearch.MakeElasticConn", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, _, err = client.Ping(url).Do(ctx)
	if err != nil {
		logElasticError(url, sniff, username, password, maxRetry, "Ping", err)
		return nil, errors.NewClientErr(nil, errstring.ESSetUpConnFail, "", "ElasticSearch.MakeElasticConn", nil)
	}
	return client, nil
}

// DoElasticConn function
func DoElasticConn(url string, sniff bool, username, password string, maxRetry int, interval int) (*elastic.Client, error) {
	val, err := utils.Do(interval, errstring.ESSetUpMaxRetry, func() (interface{}, error) {
		return MakeElasticConn(url, sniff, username, password, maxRetry)
	})
	if err != nil {
		return nil, err
	}
	return val.(*elastic.Client), err
}
