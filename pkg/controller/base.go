package controller

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
)

type Controller struct {
	DB     *store.DBSession
	Config *models.Config
	Logger *graylogger.GrayLogger
}

var ErrMissingJobId = errors.New("no job id provided")

const (
	EqualitySymbol         = "="
	InequalitySymbol       = "!"
	LessThanSymbol         = "[lt]"
	LessThanEqualSymbol    = "[lte]"
	GreaterThanSymbol      = "[gt]"
	GreaterThanEqualSymbol = "[gte]"
	ContainsSymbol         = "[contains]"
)

func cleanValue(value any) any {
	if s, ok := value.(string); ok {
		if res, err := time.Parse(time.RFC3339, s); err == nil {
			return res
		}
	}
	return value
}

func validateParam(param string, columns []string) error {
	if !slices.Contains(columns, param) {
		return fmt.Errorf("%s is not a valid query parameter", param)
	}
	return nil
}

func (base *Controller) getFilteredQuery(query map[string][]string, tableName string, validColumns []string) (*gocqlx.Queryx, error) {
	var values []any

	q := qb.Select(tableName)
	for k, v := range query {
		val := cleanValue(v[0])
		if strings.Contains(k, LessThanSymbol) {
			param := strings.Split(k, LessThanSymbol)[0]
			if err := validateParam(param, validColumns); err != nil {
				return nil, err
			}
			q.Where(qb.Lt(param))
			values = append(values, val)
			continue
		}
		if strings.Contains(k, LessThanEqualSymbol) {
			param := strings.Split(k, LessThanEqualSymbol)[0]
			if err := validateParam(param, validColumns); err != nil {
				return nil, err
			}
			q.Where(qb.LtOrEq(param))
			values = append(values, val)
			continue
		}
		if strings.Contains(k, GreaterThanSymbol) {
			param := strings.Split(k, GreaterThanSymbol)[0]
			if err := validateParam(param, validColumns); err != nil {
				return nil, err
			}
			q.Where(qb.Gt(param))
			values = append(values, val)
			continue
		}
		if strings.Contains(k, GreaterThanEqualSymbol) {
			param := strings.Split(k, GreaterThanEqualSymbol)[0]
			if err := validateParam(param, validColumns); err != nil {
				return nil, err
			}
			q.Where(qb.GtOrEq(param))
			values = append(values, val)
			continue
		}

		if err := validateParam(k, validColumns); err != nil {
			return nil, err
		}
		q.Where(qb.Eq(k))
		values = append(values, val)
	}

	stmt, names := q.AllowFiltering().ToCql()
	return base.DB.Client.Query(stmt, names).Bind(values...), nil
}
