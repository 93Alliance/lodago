# Lodago

![](https://github.com/93Alliance/lodago/blob/master/image/logo.png?raw=true)

It has the same goal as the lodash library, providing rich functions for golang.

## Function

- UUID - Generate the uuid string
- Str2MD5 - Converts the string to md5
- Find - Find a specific element in sclice
- Filter - Filter a slice
- Shuffle - Shuffle a slice
- SelectFields - Output json based on the selected field
- DropFields - Output json based on the drop field
- Struct2Map - Converts the struct to map
- Map2Struct - Converts the map to struct
- Map2JSON - Converts the map to json
- JSON2Map - Converts the json to map

## Demo

- UUID

```
func main() {
	fmt.Println(lodago.UUID())
}
```

result

```
9de5f29c-b744-4b3e-b0c6-b5c1e3705510
```

- Str2MD5

```
func main() {
    str := []byte("test12345")
    fmt.Println(lodago.Str2MD5(str))
}
```

result

```
c06db68e819be6ec3d26c6038d8e8d1f
```

- Find

```
type User struct {
	ID            int
	Name          string
}

func main() {
	users := []User {
		User{1, "张三"},
		User{2, "李四"},
		User{3, "赵四"},
	}

	findID := 2

	user, isFound := lodago.Find(users, func(ele interface{}, index int) bool {
		return ele.(User).ID == findID
	})
	if isFound {
		fmt.Println(user.(User))
	}
}
```

result

```
{2 李四}
```

- Filter

```
type User struct {
	ID   int
	Name string
}

func main() {
	users := []User{
		User{1, "张三"},
		User{2, "李四"},
		User{3, "赵四"},
	}

	id := 2

	filterUsers, isFound := lodago.Filter(users, func(ele interface{}, index int) bool {
		return ele.(User).ID != id
	})
	if isFound {
		fmt.Println(filterUsers.([]User))
	}
}
```

result

```
[{1 张三} {3 赵四}]
```

- Shuffle

```
func main() {
	array := []string{"a", "b", "c", "d", "e", "f"}
	fmt.Println("打乱前：", array)
	lodago.Shuffle(array)
	fmt.Println("打乱后：", array)
}
```

result

```
打乱前： [a b c d e f]
打乱后： [b c f a d e]
```

- SelectFields

```
type SearchResult struct {
	Date     string `json:"date"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	result := SearchResult{
		Date:     "2020-03-30",
		ID:       1,
		Name:     "张三",
		Password: "123456",
	}
	fields := []string{"name", "id"}
	out, _ := lodago.SelectFields(&result, fields...)
	fmt.Println(string(out))
}
```

result

```
{"id":1,"name":"张三"}
```

- DropFields

```
type SearchResult struct {
	Date     string `json:"date"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	result := SearchResult{
		Date:     "2020-03-30",
		ID:       1,
		Name:     "张三",
		Password: "123456",
	}
    fields := []string{"password"}
	out, _ := lodago.DropFields(&result, fields...)
	fmt.Println(string(out))
}
```
result

```
{"date":"2020-03-30","id":1,"name":"张三"}
```

- Struct2Map

```
type SearchResult struct {
	Date     string `json:"date"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	result := SearchResult{
		Date:     "2020-03-30",
		ID:       1,
		Name:     "张三",
		Password: "123456",
	}
	mapValues := lodago.Struct2Map(result)
	fmt.Println(mapValues)
}
```

result

```
map[Date:2020-03-30 ID:1 Name:张三 Password:123456]
```

- Map2Struct

```
type SearchResult struct {
	Date     string `json:"date"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	result := SearchResult{
		Date:     "2020-03-30",
		ID:       1,
		Name:     "张三",
		Password: "123456",
	}
	mapValues := lodago.Struct2Map(result)

	var obj SearchResult
	lodago.Map2Struct(mapValues, &obj)
	fmt.Println(obj)
}
```

result

```
{2020-03-30 1 张三 123456}
```

- Map2JSON

```
func main() {
	mapValues := map[string]interface{}{
		"id":   1,
		"name": "peter",
	}
	jsonStr, err := lodago.Map2JSON(mapValues)
	if err == nil {
		fmt.Println(jsonStr)
	}
}
```

result

```
{"id":1,"name":"peter"}
```

- JSON2Map

```
func main() {
	mapValues := map[string]interface{}{
		"id":   1,
		"name": "peter",
	}
	jsonStr, err := lodago.Map2JSON(mapValues)
	if err == nil {
		mapValues, err := lodago.JSON2Map([]byte(jsonStr))
		if err == nil {
			fmt.Println(mapValues)
		}
	}
}
```

result

```
map[id:1 name:peter]
```