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
- Concat - Concat multi slice
- After - The function does not actually execute until it has been called n times
- Remove - Remove a element in slice by index
- LowerFirst - Lower first char
- UpperFirst - Upper first char
- CamelCase - Converts string to camelCase, `group_id -> groupId`
- SnakeCase - Converts string to SnakeCase, `GroupId -> group_id`

## Demo

- **UUID**

```
func main() {
	fmt.Println(lodago.UUID())
}
```

result

```
9de5f29c-b744-4b3e-b0c6-b5c1e3705510
```

- **Str2MD5**

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

- **Find**

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

- **Filter**

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

- **Shuffle**

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

- **SelectFields**

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

- **DropFields**

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

- **Struct2Map**

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

- **Map2Struct**

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

- **Map2JSON**

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

- **JSON2Map**

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

- **Concat**

```
a1 := []int{1, 2, 3}
a2 := []int{4, 5}
a3 := []int{6}
result := lodago.Concat(a1, a2, a3)
fmt.Println(result)         // interface{} type
fmt.Println(result.([]int)) // original type
```
result

```
[1 2 3 4 5 6]
[1 2 3 4 5 6]
```

or

```
type User struct {
	ID   int
	Name string
}
func main() {
	a1 := []User{User{1, "u1"}, User{2, "u2"}, User{3, "u3"}}
	a2 := []User{User{4, "u4"}, User{5, "u5"}}
	a3 := []User{User{6, "u6"}}
	result := lodago.Concat(a1, a2, a3)
	fmt.Println(result)          // interface{} type
	fmt.Println(result.([]User)) // original type
}
```
result
```
[{1 u1} {2 u2} {3 u3} {4 u4} {5 u5} {6 u6}]
[{1 u1} {2 u2} {3 u3} {4 u4} {5 u5} {6 u6}]
```

- **Fill**

```
func main() {
	var slice1 []int = make([]int, 10)
	var slice2 []int = make([]int, 10)
	var slice3 []int = make([]int, 10)
	fmt.Println("slice1 origin: ", slice1)
	lodago.Fill(slice1, 6, 0, 5)
	fmt.Println("slice1 fill:   ", slice1)

	fmt.Println("slice2 origin: ", slice2)
	lodago.Fill(slice2, 7)
	fmt.Println("slice2 fill:   ", slice2)

	fmt.Println("slice3 origin: ", slice3)
	lodago.Fill(slice3, 8, 3)
	fmt.Println("slice3 fill:   ", slice3)
}
```
result

```
slice1 origin:  [0 0 0 0 0 0 0 0 0 0]
slice1 fill:    [6 6 6 6 6 0 0 0 0 0]
slice2 origin:  [0 0 0 0 0 0 0 0 0 0]
slice2 fill:    [7 7 7 7 7 7 7 7 7 7]
slice3 origin:  [0 0 0 0 0 0 0 0 0 0]
slice3 fill:    [0 0 0 8 8 8 8 8 8 8]
```

- **After**

```
func main() {
	cf := lodago.After(3, func(idx int) {
		fmt.Println("idx is ", idx)
	}).(func(int))

	cf(1)
	cf(2)
	cf(3)
	cf(4)
}
```
result

```
idx is  3
idx is  4
```

or

```
func main() {
	cf := lodago.After(3, func(idx int) (string, error) {
		fmt.Println("idx is ", idx)
		return "ok", errors.New("error")
	}).(func(int) (string, error))

	fmt.Println(cf(1))
	fmt.Println(cf(2))
	fmt.Println(cf(3))
	fmt.Println(cf(4))
}
```

result

```
 <nil>
 <nil>
idx is  3
ok error
idx is  4
ok error
```

- **Remove**

```
// User 测试
type User struct {
	ID   int
	Name string
}

func main() {
	a := []string{"A", "B", "C", "D", "E"}
	lodago.Remove(&a, 1)
	fmt.Println(a)
	b := []int{1, 2, 3, 4, 5}
	lodago.Remove(&b, 4)
	fmt.Println(b)
	c := []User{
		User{1, "张三"},
		User{2, "李四"},
		User{3, "小芳"},
		User{4, "二狗子"},
	}
	lodago.Remove(&c, 2)
	fmt.Println(c)
}
```

result

```
[A C D E]
[1 2 3 4]
[{1 张三} {2 李四} {4 二狗子}]
```