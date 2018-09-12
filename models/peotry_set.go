package models

import (
	"fmt"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

type PeotrySet struct {
	ID    	int  	`gorm:"column:id;primary_key" json:"id"`

	UID   	int64  	`gorm:"column:u_id" json:"-"`
	UUser 	*User  	`gorm:"foreignkey:u_id" json:"user,omitempty"`

	SName 	string 	`gorm:"column(s_name);size(100)" json:"name"`
}

func initSystemPeotrySet() {
	setsJson, err := ioutil.ReadFile("data/sys-peotry-set.json")
	if err != nil {
		fmt.Println("read sys-peotry-set.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(setsJson)
	re.ForEach(func (key, value gjson.Result) bool {
		sId := value.Get("s_id").Int()
		uId	:= value.Get("u_id").Int()
		sName := value.Get("s_name").String()
		savePeotrySet(int(sId), uId, sName)
		return true
	})
}

func savePeotrySet(id int, uId int64, name string) {
	peotrySet := PeotrySet {
		ID:		id,
		UID:	uId,
		SName:	name,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Save(peotrySet).Error
	if err != nil {
		fmt.Println(err)
	}
}

// type Peotryset2 struct {
// 	Id    int64  `json:"id"`
// 	SName string `json:"name"`
// }

// func (t *Peotryset) TableName() string {
// 	return "peotryset"
// }

// func init() {
// 	orm.RegisterModel(new(Peotryset))
// }

// // AddPeotryset insert a new Peotryset into database and returns
// // last inserted Id on success.
// func AddPeotryset(m *Peotryset) (id int64, err error) {
// 	o := orm.NewOrm()
// 	id, err = o.Insert(m)
// 	return
// }

// // GetPeotrysetById retrieves Peotryset by Id. Returns error if
// // Id doesn't exist
// func GetPeotrysetById(id int64) (v *Peotryset, err error) {
// 	o := orm.NewOrm()
// 	v = &Peotryset{Id: id}
// 	if err = o.Read(v); err == nil {
// 		return v, nil
// 	}
// 	return nil, err
// }

// // GetAllPeotryset retrieves all Peotryset matches certain condition. Returns empty list if
// // no records exist
// func GetAllPeotryset(query map[string]string, fields []string, sortby []string, order []string,
// 	offset int64, limit int64) (ml []interface{}, err error) {
// 	o := orm.NewOrm()
// 	qs := o.QueryTable(new(Peotryset))
// 	// query k=v
// 	for k, v := range query {
// 		// rewrite dot-notation to Object__Attribute
// 		k = strings.Replace(k, ".", "__", -1)
// 		if strings.Contains(k, "isnull") {
// 			qs = qs.Filter(k, (v == "true" || v == "1"))
// 		} else {
// 			qs = qs.Filter(k, v)
// 		}
// 	}
// 	// order by:
// 	var sortFields []string
// 	if len(sortby) != 0 {
// 		if len(sortby) == len(order) {
// 			// 1) for each sort field, there is an associated order
// 			for i, v := range sortby {
// 				orderby := ""
// 				if order[i] == "desc" {
// 					orderby = "-" + v
// 				} else if order[i] == "asc" {
// 					orderby = v
// 				} else {
// 					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
// 				}
// 				sortFields = append(sortFields, orderby)
// 			}
// 			qs = qs.OrderBy(sortFields...)
// 		} else if len(sortby) != len(order) && len(order) == 1 {
// 			// 2) there is exactly one order, all the sorted fields will be sorted by this order
// 			for _, v := range sortby {
// 				orderby := ""
// 				if order[0] == "desc" {
// 					orderby = "-" + v
// 				} else if order[0] == "asc" {
// 					orderby = v
// 				} else {
// 					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
// 				}
// 				sortFields = append(sortFields, orderby)
// 			}
// 		} else if len(sortby) != len(order) && len(order) != 1 {
// 			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
// 		}
// 	} else {
// 		if len(order) != 0 {
// 			return nil, errors.New("Error: unused 'order' fields")
// 		}
// 	}

// 	var l []Peotryset
// 	qs = qs.OrderBy(sortFields...)
// 	if _, err = qs.Limit(limit, offset).RelatedSel().All(&l, fields...); err == nil {
// 		if len(fields) == 0 {
// 			for _, v := range l {
// 				ml = append(ml, v)
// 			}
// 		} else {
// 			// trim unused fields
// 			for _, v := range l {
// 				m := make(map[string]interface{})
// 				val := reflect.ValueOf(v)
// 				for _, fname := range fields {
// 					m[fname] = val.FieldByName(fname).Interface()
// 				}
// 				ml = append(ml, m)
// 			}
// 		}
// 		return ml, nil
// 	}
// 	return nil, err
// }

// // UpdatePeotryset updates Peotryset by Id and returns error if
// // the record to be updated doesn't exist
// func UpdatePeotrysetById(m *Peotryset) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotryset{Id: m.Id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Update(m); err == nil {
// 			fmt.Println("Number of records updated in database:", num)
// 		}
// 	}
// 	return
// }

// // DeletePeotryset deletes Peotryset by Id and returns error if
// // the record to be deleted doesn't exist
// func DeletePeotryset(id int64) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotryset{Id: id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Delete(&Peotryset{Id: id}); err == nil {
// 			fmt.Println("Number of records deleted in database:", num)
// 		}
// 	}
// 	return
// }
