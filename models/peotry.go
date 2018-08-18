package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Peotry struct {
	Id       int64     `orm:"column(p_id);pk" json:"id"`
	SId      *Peotryset     `orm:"column(s_id);rel(fk);" json:"set"`
	// SId      int64     `orm:"column(s_id);" json:"set"`
	UId      *User     `orm:"column(u_id);rel(fk);" json:"user"`
	// UId      int64     `orm:"column(u_id);" json:"userId"`
	PTitle   string    `orm:"column(p_title);size(20);null" json:"title"`
	PTime    time.Time `orm:"column(p_time);type(datetime);null" json:"time"`
	PContent string    `orm:"column(p_content);null" json:"content"`
	PEnd     string    `orm:"column(p_end);null" json:"end"`
	PImages  string    `orm:"column(p_images);null" json:"images"`
}

type Peotry2 struct {
	Id       int64     `json:"id"`
	SId      *Peotryset2     `json:"set"`
	UId      *User2     `json:"user"`
	PTitle   string    `json:"title"`
	PTime    time.Time `json:"time"`
	PContent string    `json:"content"`
	PEnd     string    `json:"end"`
	PImages  string    `json:"images"`
}

func (t *Peotry) TableName() string {
	return "peotry"
}

func init() {
	orm.RegisterModel(new(Peotry))
}

// AddPeotry insert a new Peotry into database and returns
// last inserted Id on success.
func AddPeotry(m *Peotry) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPeotryById retrieves Peotry by Id. Returns error if
// Id doesn't exist
func GetPeotryById(id int64) (v *Peotry, err error) {
	o := orm.NewOrm()
	v = &Peotry{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPeotry retrieves all Peotry matches certain condition. Returns empty list if
// no records exist
func GetAllPeotry(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Peotry))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Peotry
	qs = qs.OrderBy(sortFields...)
	//关联查询后的数据顺序不是按照简单查询的数据顺序
	if _, err = qs.Limit(limit, offset).RelatedSel().All(&l, fields...); err == nil {
		var ll []Peotry2
		for _, v := range l {
			vv := Peotry2 {
				Id: v.Id,
				SId: &Peotryset2{Id: v.SId.Id, SName: v.SId.SName},
				UId: &User2{Id: v.UId.Id, UName: v.UId.UName},
				PTitle: v.PTitle,
				PTime: v.PTime,
				PContent: v.PContent,
				PEnd: v.PEnd,
				PImages: v.PImages,
			}

			imageV, err := GetPeotryimageById(v.Id)
			if err == nil {
				vv.PImages = imageV.IImages
			}
			ll = append(ll, vv)
		}

		if len(fields) == 0 {
			for _, v := range ll {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range ll {
				m := make(map[string]interface{})

				val := reflect.ValueOf(v)
				typ := reflect.TypeOf(&v)	// 反射获取struct中的tag
				for _, fname := range fields {
					typElem := typ.Elem()
					field, ok := typElem.FieldByName(fname)
					if ok {
						m[field.Tag.Get("json")] = val.FieldByName(fname).Interface()
					} else {
						m[fname] = val.FieldByName(fname).Interface()
					}
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePeotry updates Peotry by Id and returns error if
// the record to be updated doesn't exist
func UpdatePeotryById(m *Peotry) (err error) {
	o := orm.NewOrm()
	v := Peotry{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePeotry deletes Peotry by Id and returns error if
// the record to be deleted doesn't exist
func DeletePeotry(id int64) (err error) {
	o := orm.NewOrm()
	v := Peotry{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Peotry{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
