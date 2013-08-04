package app

import (
	"errors"
)

type appFilter func(_app *WebApp) error

var filterMaps = map[string]appFilter{
	"init":   AppInitFilter,
	"input":  AppInputFilter,
	"url":    AppUrlFilter,
	"output": AppOutputFilter,
	"clean":  AppCleanFilter,
}

func InitFilter(filterName string, _app *WebApp) error {
	if filter, ok := filterMaps[filterName]; ok {
		return filter(_app)
	}
	return errors.New("filter.notfound filtername=" + filterName)
}
