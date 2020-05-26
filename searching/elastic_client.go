package searching

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"code.aliyun.com/bim_backend/zoogoer/gun/errors"
	sl "code.aliyun.com/bim_backend/zoogoer/gun/singleflight"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"github.com/wzkun/aurora/utils/decode"
	"github.com/wzkun/aurora_database/errstring"
)

// ElasticClient struct
type ElasticClient struct {
	ec      *elastic.Client
	lg      sl.Group
	indices sync.Map
}

// NewElasticClient create ElasticClient
func NewElasticClient(url string, sniff bool, username, password string, maxRetry int) (*ElasticClient, error) {
	db, err := DoElasticConn(url, sniff, username, password, maxRetry, interval)
	o := new(ElasticClient)
	o.ec = db
	return o, err
}

func (m *ElasticClient) parseBulkError(where string, res *elastic.BulkResponse, err error) (string, bool) {
	if err == nil && !res.Errors {
		return "", false
	}

	if err != nil {
		estr := err.Error()
		logrus.Errorf("ElasticClient %s to ES Failed with Fatal Error % s", where, estr)
		return estr, true
	}

	var docs string
	for _, v := range res.Items {
		value, _ := decode.JSON.Marshal(v)
		docs = docs + string(value)
	}
	logrus.Errorf("ElasticClient %s to ES Failed with Validate Error % s", where, docs)
	return docs, true
}

func (m *ElasticClient) key(index, source string) string {
	b := bytes.Buffer{}
	b.WriteString(index)
	b.WriteString(source)
	return b.String()
}

// indexCached Function indexCached a bucket
func (m *ElasticClient) indexCached(index string) bool {
	_, hit := m.indices.Load(index)
	return hit
}

// CheckOrInitIndex Function CheckOrInitIndex a bucket
func (m *ElasticClient) CheckOrInitIndex(index string) error {
	if hit := m.indexCached(index); hit {
		return nil
	}

	_, err := m.lg.Do("CheckOrInitIndex"+index, func() (interface{}, error) {
		exist, err := m.IndexExists(index)
		if err != nil {
			return nil, err
		}
		if !exist {
			if err := m.CreateIndex(index); err != nil {
				return nil, err
			}
			mappingString := fmt.Sprintf(mapping, index)
			if err := m.PutIndexMapping(index, mappingString); err != nil {
				return nil, err
			}
		}

		m.indices.Store(index, index)
		return index, nil
	})

	return err
}

// IndexExists function
func (m *ElasticClient) IndexExists(index string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	indexExists, err := m.ec.IndexExists(index).Do(ctx)
	if err != nil {
		return false, errors.NewClientErr(nil, errstring.ESCheckIndexExistFailed, "IndexExists Failed", "ElasticClient.IndexExists", nil)
	}
	return indexExists, nil
}

// CreateIndex function
func (m *ElasticClient) CreateIndex(index string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	createIndex, err := m.ec.CreateIndex(index).Do(ctx)
	if err != nil {
		return errors.NewClientErr(nil, errstring.ESCreateIndexFailed, err.Error(), "ElasticClient.CreateIndex", nil)
	}
	if !createIndex.Acknowledged {
		return errors.NewClientErr(nil, errstring.ESCreateIndexFailed, "PutMapping Response is Nil", "ElasticClient.CreateIndex", nil)
	}
	return nil
}

// PutIndexMapping Function PutIndexMapping a bucket
func (m *ElasticClient) PutIndexMapping(index string, mappningString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	putresp, err := m.ec.PutMapping().Index(index).Type(index).BodyString(mappningString).Do(ctx)

	if err != nil {
		return errors.NewClientErr(nil, errstring.ESCreateIndexFailed, err.Error(), "ElasticClient.PutIndexMapping", nil)
	}
	if putresp == nil {
		return errors.NewClientErr(nil, errstring.ESCreateIndexFailed, "PutMapping Response is Nil", "ElasticClient.PutIndexMapping", nil)
	}
	if !putresp.Acknowledged {
		return errors.NewClientErr(nil, errstring.ESCreateIndexFailed, "PutMapping Response not Acknowledged", "ElasticClient.PutIndexMapping", nil)
	}
	return nil
}

// Create function
func (m *ElasticClient) Create(im ElasticItem) error {
	return m.CreateMulti(im)
}

// CreateMulti function
func (m *ElasticClient) CreateMulti(ims ...ElasticItem) error {
	if len(ims) <= 0 {
		return nil
	}

	bulk := m.ec.Bulk()
	for _, im := range ims {
		index := im.Parrent()
		if err := m.CheckOrInitIndex(index); err != nil {
			return err
		}

		id := im.Idx()
		data := im.JSONData()
		req := elastic.NewBulkIndexRequest().Index(index).Type(index).Id(id).Doc(data)
		bulk = bulk.Add(req)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := bulk.Do(ctx)
	estr, hasErr := m.parseBulkError("CreateMulti", res, err)
	if !hasErr {
		return nil
	}
	return errors.NewClientErr(nil, errstring.ESCreateFailed, estr, "ElasticClient.CreateMulti", nil)
}

// Update function
func (m *ElasticClient) Update(im ElasticItem) error {
	return m.CreateMulti(im)
}

// UpdateMulti function
func (m *ElasticClient) UpdateMulti(ims ...ElasticItem) error {
	if len(ims) <= 0 {
		return nil
	}

	return m.CreateMulti(ims...)
}

// Get function
func (m *ElasticClient) Get(id, parrent string) (*json.RawMessage, error) {
	// https://github.com/olivere/elastic/blob/a7f6620ddddaaccb6fa041eaa281182fd4b4fcc4/get.go#L213:48
	// https://github.com/olivere/elastic/blob/release-branch.v6/get_test.go#L23:2
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := m.ec.Get().Index(parrent).Type(parrent).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Source, nil
}

// GetMulti function
func (m *ElasticClient) GetMulti(parrent string, ids ...string) ([]*json.RawMessage, error) {
	// https://github.com/olivere/elastic/blob/a7f6620ddddaaccb6fa041eaa281182fd4b4fcc4/mget_test.go#L75:10
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return nil, err
	}

	bulk := m.ec.MultiGet()
	for _, id := range ids {
		req := elastic.NewMultiGetItem().Index(parrent).Type(parrent).Id(id)
		bulk = bulk.Add(req)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	bulk.Do(ctx)
	return nil, nil
}

// Delete function
func (m *ElasticClient) Delete(id, parrent string) error {
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return err
	}

	return m.DeleteMulti(parrent, id)
}

// DeleteMulti function
func (m *ElasticClient) DeleteMulti(parrent string, ids ...string) error {
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return err
	}

	bulk := m.ec.Bulk()
	for _, id := range ids {
		req := elastic.NewBulkDeleteRequest().Index(parrent).Type(parrent).Id(id)
		bulk = bulk.Add(req)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := bulk.Do(ctx)
	estr, hasErr := m.parseBulkError("DeleteMulti", res, err)
	if !hasErr {
		return nil
	}
	return errors.NewClientErr(nil, errstring.ESCreateFailed, estr, "ElasticClient.DeleteMulti", nil)
}

// Flush function
func (m *ElasticClient) Flush(parrent string) error {
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := m.ec.Flush().Index(parrent).Do(ctx)
	return err
}

// DeleteIndex function
func (m *ElasticClient) DeleteIndex(parrent string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if m.indexCached(parrent) {
		m.indices.Delete(parrent)
	}

	_, err := m.ec.DeleteIndex(parrent).Do(ctx)
	return err
}

// Search function
func (m *ElasticClient) Search(parrent, query string) (string, error) {
	if err := m.CheckOrInitIndex(parrent); err != nil {
		return "", err
	}

	key := m.key(parrent, query)

	response, err := m.lg.Do(key, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		result, err := m.ec.Search().Index(parrent).Source(query).Do(ctx)
		if err != nil {
			logrus.Infof("Elastic Search returned with Error: %s", err.Error())
			return nil, err
		}
		return decode.JSON.Marshal(result)
	})

	if err != nil {
		return "{}", nil
	}

	result := response.([]byte)
	return string(result), nil
}
