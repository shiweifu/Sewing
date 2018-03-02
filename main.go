package main

import (
  "strings"
    "encoding/json"
    "fmt"
)


type JSONValue = map[string]interface{}

var subModelList = []string{}

var jsonStr = `
{
    "from_user_id": "8",
    "to_user_id": "4",
    "status": "1",
    "create_time": "1517477417",
    "type": "0",
    "user": {
        "user_id": "4",
        "unique_id": "",
        "name": "minnerhuang",
        "avatar_url": "https://static.chat.jianda.com/2018-01-30/5a70295b7de71.jpg"
    }
}
`

func formatClassName(name string) string {
    return strings.Title(formatPropertyName(name))
}

func formatPropertyName(name string) string  {
    arr := strings.Split(name, "_")
    if len(arr) == 1 {
        return strings.ToLower(name)
    }

    result := ""

    for idx, elem := range arr {
        if idx == 0 {
            result += strings.ToLower(elem)
            continue
        }

        result += strings.Title(elem)
    }
    return result
}

var jsonModelFormat = `
class %s: Mappable {
%s

  required init(_ map: [String: Any]) {
%s
  }
} `

func jsonDict2Class(jsonValue JSONValue, clsName string, subModel bool) {
    var propertyList = []string{}
    var mapList = []string{}

    for k, v := range jsonValue {
        switch v.(type) {
        case string:
            var property = formatPropertyName(k)
            var line = fmt.Sprintln("    var", property, ":",  "String? = nil")
            var initLine = fmt.Sprintf("    %s <- map.property(\"%s\")\n", property, k)
            propertyList = append(propertyList, line)
            mapList = append(mapList, initLine)
        case int:
            var property = formatPropertyName(k)
            var line = fmt.Sprintln("    var", property, ":",  "Int? = nil")
            var initLine = fmt.Sprintf("    %s <- map.property(\"%s\")\n", property, k)
            propertyList = append(propertyList, line)
            mapList = append(mapList, initLine)

        case float64:
            var property = formatPropertyName(k)
            var line = fmt.Sprintln("    var", property, ":",  "Float? = nil")
            var initLine = fmt.Sprintf("    %s <- map.property(\"%s\")\n", property, k)
            propertyList = append(propertyList, line)
            mapList = append(mapList, initLine)
        case interface{}:
            var property = formatPropertyName(k)
            var clsName = formatClassName(k)
            var line = fmt.Sprintln("var", property, ":",  clsName, "? = nil")
            var initLine = fmt.Sprintf("    %s <- map.property(\"%s\")\n", property, k)
            propertyList = append(propertyList, line)
            mapList = append(mapList, initLine)

            if newV, ok := v.(JSONValue); ok {
                fmt.Println("进入子model里面")
                jsonDict2Class(newV, clsName, true)
            }
        default:
        }
    }

    var result = fmt.Sprintf(jsonModelFormat, clsName, strings.Join(propertyList, ""), strings.Join(mapList, ""))
    if !subModel {
        fmt.Println(result)
    } else {
        subModelList = append(subModelList, result)
    }
}

func main() {
    jsonValue := JSONValue{}
    err := json.Unmarshal([]byte(jsonStr), &jsonValue)
    if err != nil {
        fmt.Println(err)
        return
    }
  jsonDict2Class(jsonValue, "UserObj", false)

    for _, v := range subModelList {
        fmt.Println(v)
        fmt.Println()
    }

}

