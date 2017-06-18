package models

import (
	"fmt"
	"github.com/kataras/go-errors"
	"time"
)

const MetricValueTable = "metric_value"

type MetricValue struct {
	Id 		int64		`json:"id,string,omitempty"`
	MetricName	string		`json:"metric_name,omitempty"`
	Value           int64		`json:"value,string,omitempty"`
	CreateDate 	time.Time	`json:"create_date,string,omitempty"`
	Host		string		`json:"host,omitempty"`
}

func (p MetricValue) ListHosts() (*[]string, error) {

	sql := fmt.Sprint("select distinct `host` from metric_value")
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []string
	for rows.Next() {
		p := new(string)
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, *p)
	}

	return &hosts, nil
}

func (p MetricValue) FindByMetricNameBetweenDates(host string, name string, startTime time.Time, endTime time.Time) (*[]MetricValue, error) {
	if name == "" || startTime.IsZero() || endTime.IsZero() || host == "" {
		return nil, errors.New("Please provide a name, starTime, and endTime")
	}

	sql := fmt.Sprint("select `name`, `value`, `host`, `create_date` from metric_value where `host`=? and `name`=? and `create_date` >=? and `create_date` <=?")
	rows, err := db.Query(sql, host, name, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metricValues []MetricValue
	for rows.Next() {
		p := new(MetricValue)
		err := rows.Scan(&p.MetricName, &p.Value, &p.Host, &p.CreateDate)
		if err != nil {
			return nil, err
		}
		metricValues = append(metricValues, *p)
	}

	return &metricValues, nil
}

func (p MetricValue) Save() (*MetricValue, error) {
	var sql string
	if p.Id == 0 {
		p.CreateDate = time.Now()
		p.CreateDate.Format(time.RFC3339)
		sql = fmt.Sprintf("INSERT INTO %s (`name`, `value`, `host`, `create_date`) VALUES (?,?,?)", MetricValueTable)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `name`=?, `value`, `host`, `create_date`=? WHERE `id`=%d", MetricValueTable, p.Id)
	}

	res, err := db.Exec(sql, p.MetricName, p.Value, p.Host, p.CreateDate)
	if err != nil {
		return nil, err
	}

	if p.Id == 0 {
		p.Id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

