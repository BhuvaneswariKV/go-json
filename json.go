package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type JSON struct {
	jmap map[string]interface{}
}

/////////////////////////////////////////////////////
////////////////////// Create ///////////////////////
/////////////////////////////////////////////////////

func New() (newObject JSON) {
	m := make(map[string]interface{})
	newObject.jmap = m
	return newObject
}

func NewFrom(m map[string]interface{}) (newObject JSON) {
	newObject.jmap = m
	return newObject
}

/////////////////////////////////////////////////////
/////////////////////// Parse ///////////////////////
/////////////////////////////////////////////////////

func ParseString(data string) JSON {
	return Parse([]byte(data))
}

func Parse(data []byte) (parsed JSON) {
	//Currently it supports top level map not array
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return
	}
	switch f.(type) {
	case []interface{}:
		log.Println("Found Array, Not Supported")
		return
	default:
		m := f.(map[string]interface{})
		parsed.jmap = m
		return parsed
	}
	return
}

func ParseReadCloser(rc io.ReadCloser) JSON {
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(-255)
	}
	return Parse(data)
}

/////////////////////////////////////////////////////
////////////////////// Availble /////////////////////
/////////////////////////////////////////////////////

func (jobj *JSON) HasKey(k string) bool {
	if _, ok := jobj.jmap[k]; ok {
		return true
	}
	return false
}

/////////////////////////////////////////////////////
//////////////////////GetKeyList/////////////////////
/////////////////////////////////////////////////////

func (jobj *JSON) GetKeyList() []string {
	keys := make([]string, 0)
	for k, _ := range jobj.jmap {
		keys = append(keys, k)
	}
	return keys
}

/////////////////////////////////////////////////////
////////////////////// SetValue /////////////////////
/////////////////////////////////////////////////////

func (jobj *JSON) Put(k string, v interface{}) {
	jobj.jmap[k] = v
}

/////////////////////////////////////////////////////
////////////////////// GetValue /////////////////////
/////////////////////////////////////////////////////

func (jobj *JSON) Get(k string) interface{} {
	return jobj.jmap[k]
}

func (jobj *JSON) GetString(k string) string {
	data := jobj.jmap[k]
	if str, ok := data.(string); ok {
		return str
	}
	return ""
}

func (j *JSON) GetJSON(k string) JSON {
	f := j.jmap[k]
	switch f.(type) {
	case interface{}:
		m := f.(map[string]interface{})
		return NewFrom(m)
	}
	return New()
}

func (j *JSON) GetJSONArray(k string) []JSON {
	f := j.jmap[k]
	jarr := make([]JSON, 0)
	switch val := f.(type) {
	case []interface{}:
		for _, u := range val {
			switch u.(type) {
			case map[string]interface{}:
				m := u.(map[string]interface{})
				jarr = append(jarr, NewFrom(m))
			}
		}
	}
	return jarr
}

func (jobj *JSON) GetAsStringArray(k string) []string {
	v := jobj.jmap[k]
	str := make([]string, 0)
	switch val := v.(type) {
	case []float64: //float array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%f", u))
		}
	case []int: //int array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%d", u))
		}
	case []int64: //int array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%d", u))
		}
	case []string: //String array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%s", u))
		}
	case []bool: //String array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%t`", u))
		}
	case []interface{}: //String array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%s", u))
		}
	case []JSON: //JSON array
		for _, u := range val {
			str = append(str, fmt.Sprintf("%s", u))
		}
	default:
		log.Println(k, "is of a type I don't know how to handle ", reflect.TypeOf(v))
	}
	return str
}

/////////////////////////////////////////////////////
////////////////////// ToString /////////////////////
/////////////////////////////////////////////////////

func ToJSONString(v interface{}) string {
	return string(ToJSONByte(v))
}

func ToJSONByte(v interface{}) []byte {
	buff, _ := json.Marshal(v)
	return buff
}

func (jobj *JSON) ToString() string {
	m := jobj.jmap
	str := "{"
	for k, v := range m {
		switch val := v.(type) {
		case string:
			str += fmt.Sprintf("\"%s\":\"%s\",", k, val)
		case bool:
			str += fmt.Sprintf("\"%s\":\"%t\",", k, val)
		case int:
			str += fmt.Sprintf("\"%s\":%d,", k, val)
		case int64:
			str += fmt.Sprintf("\"%s\":%d,", k, val)
		case float64:
			str += fmt.Sprintf("\"%s\":%f,", k, val)
		case JSON:
			str += fmt.Sprintf("\"%s\":%s,", k, val.ToString())
		case []float64: //float array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("%f,", u)
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case []int: //int array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("%d,", u)
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case []int64: //int array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("%d,", u)
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case []string: //String array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("\"%s\",", u)
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case []interface{}: //String array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("\"%s\",", u)
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case []JSON: //JSON array
			str += fmt.Sprintf("\"%s\":[", k)
			for _, u := range val {
				str += fmt.Sprintf("%s,", u.ToString())
			}
			str = prune(str, ",") //prune extra comma if needed
			str += "],"
		case map[string]interface{}: //JSON
			j := NewFrom(val)
			str += fmt.Sprintf("\"%s\":%s", k, j.ToString())
		default:
			log.Println(k, "is of a type I don't know how to handle ", reflect.TypeOf(v))
		}
	}
	str = prune(str, ",") //prune extra comma if needed
	return str + "}"
}

func prune(str string, splitter string) string {
	if strings.LastIndex(str, splitter)+1 == len(str) {
		return str[0 : len(str)-1]
	}
	return str
}
