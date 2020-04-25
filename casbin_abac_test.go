package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
)

type Person struct {
	Role string
	Name string
}
type Gate struct {
	Name string
}
type Env struct {
	Time     time.Time
	Location string
}

func (env *Env) IsSchooltime() bool {
	return env.Time.Hour() >= 8 && env.Time.Hour() <= 18
}

const modelText1 = `
[request_definition]
r = sub, obj, act, env

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub.Role=='Teacher' && r.obj.Name=='School Gate' && r.act in('In','Out') && r.env.Time.Hour >=7 && r.env.Time.Hour <= 18
`

func InitEnv(hour int) *Env {
	env := &Env{}
	env.Time = time.Date(2020, 4, 24, hour, 0, 0, 0, time.Local)
	return env
}

//测试matcher支持逻辑表达式：>=,<=
func test1() {

	p1 := Person{Role: "Student", Name: "Yun"}
	p2 := Person{Role: "Teacher", Name: "Devin"}
	persons := []Person{p1, p2}

	g1 := Gate{Name: "School Gate"}
	g2 := Gate{Name: "Factory Gate"}
	gates := []Gate{g1, g2}

	m := model.Model{}
	m.LoadModelFromText(modelText1)
	e := casbin.NewEnforcer(m, true)

	envs := []*Env{InitEnv(9), InitEnv(23)}
	fmt.Println(e.Enforce(p1, g1, "In", InitEnv(10)))
	for _, env := range envs {
		fmt.Println("\r\nTime:", env.Time.Local())
		for _, p := range persons {
			for _, g := range gates {

				pass := e.Enforce(p, g, "In", env)
				fmt.Println(p.Role, p.Name, "In", g.Name, pass)

				pass = e.Enforce(p, g, "Control", env)
				fmt.Println(p.Role, p.Name, "Control", g.Name, pass)
			}
		}
	}
}

const modelText2 = `
[request_definition]
r = sub, obj, act, env

[policy_definition]
p = sub, obj,act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub.Role=='Teacher' && r.obj.Name=='School Gate' && r.act in('In','Out') && r.env.IsSchooltime()
`

//测试支持定义函数方式1
func test2() {

	p1 := Person{Role: "Student", Name: "Yun"}
	p2 := Person{Role: "Teacher", Name: "Devin"}
	persons := []Person{p1, p2}

	g1 := Gate{Name: "School Gate"}
	g2 := Gate{Name: "Factory Gate"}
	gates := []Gate{g1, g2}

	m := model.Model{}

	m.LoadModelFromText(modelText2)
	e := casbin.NewEnforcer(m)
	envs := []*Env{InitEnv(9), InitEnv(23)}

	for _, env := range envs {
		fmt.Println("\r\nTime:", env.Time.Local())
		for _, p := range persons {
			for _, g := range gates {
				pass := e.Enforce(p, g, "In", env)
				fmt.Println(p.Role, p.Name, "In", g.Name, pass)
				pass = e.Enforce(p, g, "Control", env)
				fmt.Println(p.Role, p.Name, "Control", g.Name, pass)
			}
		}
	}
}

const modelText3 = `
[request_definition]
r = sub, obj, act, env

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act && r.env.CheckIp()
`

type Env1 struct {
	IP string
	OS string
}

func (env *Env1) CheckIp() bool {
	fmt.Println(env.IP, env.OS)
	return true
}
func (env *Env1) TCheck(os string) bool {
	fmt.Println(env.IP, env.OS, os)
	return true
}

func test3() {

	env := &Env1{}
	env.IP = "192.168.1.100"
	env.OS = "linux"
	//env := Env1{IP: "192.168.1.100", OS: "Linux"}
	m := model.Model{}
	m.LoadModelFromText(modelText3)
	e := casbin.NewEnforcer(m, true)

	e.AddPolicy("alice", "data1", "read")
	e.AddPolicy("alice", "data2", "read")
	e.AddPolicy("eve", "data3", "read")

	result := e.Enforce("alice", "data1", "read", env)

	fmt.Println(result)
}

//测试支持自定义函数方式2
const modelText4 = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && my_func(r.obj, p.obj) && r.act == p.act
`

func KeyMatch(key1 string, key2 string) bool {
	i := strings.Index(key2, "*")
	if i == -1 {
		return key1 == key2
	}

	if len(key1) > i {
		return key1[:i] == key2[:i]
	}
	return key1 == key2[:i]
}
func KeyMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(KeyMatch(name1, name2)), nil
}
func test4() {

	m := model.Model{}
	m.LoadModelFromText(modelText4)

	e := casbin.NewEnforcer(m, true)
	e.AddFunction("my_func", KeyMatchFunc)

	e.AddPolicy("alice", "data1", "read")
	e.AddPolicy("alice", "data2", "read")
	e.AddPolicy("eve", "data3", "read")
	// e.AddNamedPolicy("P", "alice", "data", "read")
	// e.AddNamedPolicy("P", "alice", "data1", "read")

	result := e.Enforce("alice", "data1", "read")

	fmt.Println(result)
}

const modelText5 = `
[request_definition]
r = sub, obj, act, env

[policy_definition]
p = sub, obj, act, ip

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act && my_func(r.env.IP, p.ip)
`

//不实现判断逻辑，仅看是否能获取参数
func KeyIPMatch(key1 string, key2 string) bool {
	fmt.Println(key1, key2)
	return true
}
func KeyIPMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(KeyIPMatch(name1, name2)), nil
}

//测试policy_effect 自定义扩展项，在matchers中使用
func test5() {
	env := &Env1{}
	env.IP = "192.168.1.100"
	env.OS = "linux"
	//env := Env1{IP: "192.168.1.100", OS: "Linux"}
	m := model.Model{}
	m.LoadModelFromText(modelText5)
	e := casbin.NewEnforcer(m, true)

	e.AddPolicy("alice", "data1", "read", "1.1.1.1")
	e.AddPolicy("alice", "data2", "read", "2.2.2.2")
	e.AddFunction("my_func", KeyIPMatchFunc)
	result := e.Enforce("alice", "data1", "read", env)

	fmt.Println(result)
}

//不实现判断逻辑，仅看是否能获取参数
const modelText6 = `
[request_definition]
r = sub, obj, act, env

[policy_definition]
p = sub, obj, act, ip, os

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act && r.env.TCheck(p.os)
`

func test6() {
	env := &Env1{}
	env.IP = "192.168.1.100"
	env.OS = "linux"
	//env := Env1{IP: "192.168.1.100", OS: "Linux"}
	m := model.Model{}
	m.LoadModelFromText(modelText6)
	e := casbin.NewEnforcer(m, true)

	e.AddPolicy("alice", "data1", "read", "1.1.1.1", "Centos")
	e.AddPolicy("alice", "data2", "read", "2.2.2.2", "Mac")

	result := e.Enforce("alice", "data1", "read", env)

	fmt.Println(result)
}
func main() {
	//test1()
	test6()
}
