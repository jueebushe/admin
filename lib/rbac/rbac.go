package rbac

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	m "github.com/osgochina/admin/models/rbacmodels"
	"strings"
)

func Access() {
	var Check = func(ctx *context.Context) {

		ret := CheckAccess(ctx.Request.RequestURI)
		fmt.Println(ret)
	}
	beego.AddFilter("*", "BeforRouter", Check)
}

func CheckAccess(url string) bool {
	fmt.Println(url)
	urllist := strings.Split(url, "?")
	params := strings.Split(strings.ToLower(urllist[0]), "/")
	if len(params) < 3 {
		return false
	}
	str := fmt.Sprintf("%s/%s/%s", params[0], params[1], params[2])
	accesslist, _ := GetAccessList(1)
	_, ok := accesslist[str]
	if ok != false {
		return true
	}
	return false
}

type AccessNode struct {
	Id        int64
	Name      string
	Childrens []*AccessNode
}

func GetAccessList(uid int64) (map[string]bool, error) {
	list, err := m.AccessList(uid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	alist := make([]*AccessNode, 0)
	for _, l := range list {
		if l["Pid"].(int64) == 0 && l["Level"].(int64) == 1 {
			anode := new(AccessNode)
			anode.Id = l["Id"].(int64)
			anode.Name = l["Name"].(string)
			alist = append(alist, anode)
		}
	}
	for _, l := range list {
		if l["Level"].(int64) == 2 {
			for _, an := range alist {
				if an.Id == l["Pid"].(int64) {
					anode := new(AccessNode)
					anode.Id = l["Id"].(int64)
					anode.Name = l["Name"].(string)
					an.Childrens = append(an.Childrens, anode)
				}
			}
		}
	}
	for _, l := range list {
		if l["Level"].(int64) == 3 {
			for _, an := range alist {
				for _, an1 := range an.Childrens {
					if an1.Id == l["Pid"].(int64) {
						anode := new(AccessNode)
						anode.Id = l["Id"].(int64)
						anode.Name = l["Name"].(string)
						an1.Childrens = append(an1.Childrens, anode)
					}
				}

			}
		}
	}
	accesslist := make(map[string]bool)
	for _, v := range alist {
		for _, v1 := range v.Childrens {
			for _, v2 := range v1.Childrens {
				vname := strings.Split(v.Name, "/")
				v1name := strings.Split(v1.Name, "/")
				v2name := strings.Split(v2.Name, "/")
				str := fmt.Sprintf("%s/%s/%s", strings.ToLower(vname[0]), strings.ToLower(v1name[0]), strings.ToLower(v2name[0]))
				accesslist[str] = true
			}
		}
	}
	return accesslist, nil
}
