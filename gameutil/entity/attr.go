// 非线程安全
package entity

type Attr struct {
	owner    *Entity                // 属于哪个实体
	attrdata map[string]interface{} // 属性容器
}

func (a *Attr) Add(k string, v interface{}) {
	a.attrdata[k] = v
}

func (a *Attr) Del(k string) {
	delete(a.attrdata, k)
}

func (a *Attr) Modify(k string, v interface{}) {
	a.attrdata[k] = v
}
