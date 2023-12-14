package web

import (
	"encoding/json"
	"fmt"
	. "github.com/doytowin/goquery/core"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type RestService[C any, E any, Q GoQuery] struct {
	*Service[C, E, Q]
	Prefix string
}

func (s *RestService[C, E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match := s.idRgx.FindStringSubmatch(request.URL.Path)
	if len(match) > 0 {
		id := match[1]
		data, err := s.Get(id)
		if data == nil {
			writeResult(writer, fmt.Errorf("record not found. id: %s", id), nil)
		} else {
			writeResult(writer, err, *data)
		}
		return
	}
	query := s.createQuery()
	queryMap := request.URL.Query()
	resolveQuery(queryMap, query)
	pageList, err := s.Page(query)
	writeResult(writer, err, pageList)
}

func resolveQuery(queryMap url.Values, query any) {
	for name, v := range queryMap {
		field := reflect.ValueOf(query).Elem().FieldByName(name)
		if !field.IsValid() {
			continue
		}
		log.Debug("field.Kind: ", field.Kind())
		if field.Kind() == reflect.Bool {
			v0 := strings.EqualFold(v[0], "TRue")
			field.Set(reflect.ValueOf(v0))
		} else if field.Kind() == reflect.Pointer {
			resolvePointer(field, v)
		}
	}
}

func resolvePointer(field reflect.Value, v []string) {
	log.Debug("field.Type: ", field.Type())
	if field.Type().String() == "*[]int" {
		strArr := strings.Split(v[0], ",")
		var v0 []int
		for _, s := range strArr {
			num, err := strconv.Atoi(s)
			if NoError(err) {
				v0 = append(v0, num)
			}
		}
		field.Set(reflect.ValueOf(&v0))
	} else if field.Type().String() == "*int" {
		v0, err := strconv.Atoi(v[0])
		if NoError(err) {
			field.Set(reflect.ValueOf(&v0))
		}
	} else {
		field.Set(reflect.ValueOf(&v[0]))
	}
}

func writeResult(writer http.ResponseWriter, err error, data any) {
	response := Response{Data: data, Success: NoError(err), Error: ReadError(err)}
	marshal, err := json.Marshal(response)
	if NoError(err) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(marshal)
	}
}