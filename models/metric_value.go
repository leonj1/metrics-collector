package models

import (
	"fmt"
	"github.com/kataras/go-errors"
	"time"
)

const MetricValueTable = "metric_value"

type MetricValue struct {
	Id 		int64		`json:"id,string,omitempty"`
	MetricName	string		`json:"metric_name,string,omitempty"`
	Value           int64		`json:"value,string,omitempty"`
	CreateDate 	time.Time	`json:"create_date,string,omitempty"`
}

func (p MetricValue) FindByMetricNameBetweenDates(name string, startTime time.Time, endTime time.Time) (*[]MetricValue, error) {
	if name == "" || startTime.IsZero() || endTime.IsZero() {
		return nil, errors.New("Please provide a name, starTime, and endTime")
	}

	sql := fmt.Sprint("select `name`, `value`, `create_date` from metric_value where `name`=? and `create_date` >=? and `create_date<=?")
	rows, err := db.Query(sql, name, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metricValues []MetricValue
	for rows.Next() {
		p := new(MetricValue)
		err := rows.Scan(&p.MetricName, &p.Value, &p.CreateDate)
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
		sql = fmt.Sprintf("INSERT INTO %s (`name`, `value`, `create_date`) VALUES (?,?,?)", MetricValueTable)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `name`=?, `value`, `create_date`=? WHERE `id`=%d", MetricValueTable, p.Id)
	}

	res, err := db.Exec(sql, p.MetricName, p.Value, p.CreateDate)
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

