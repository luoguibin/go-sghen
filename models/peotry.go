package models

import (
	"SghenApi/helper"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

type Peotry struct {
	ID       	int64     		`gorm:"column:id;primary_key;" json:"id"`

	UID			int64			`gorm:"column:u_id" json:"-"`
	UUser      	*User     		`gorm:"foreignkey:u_id;" json:"user"`

	SID			int				`gorm:"column:s_id" json:"-"`
	SSet     	*PeotrySet  	`gorm:"foreignkey:s_id" json:"set"`

	PTitle   	string    		`gorm:"column:p_title;type:varchar(20)" json:"title"`
	PTime    	time.Time 		`gorm:"column:p_time" json:"time"`
	PContent 	string    		`gorm:"column:p_content;type:mediumtext" json:"content"`
	PEnd     	string    		`gorm:"column:p_end" json:"end"`

	PImage		*PeotryImage	`gorm:"foreignkey:id" json:"image,omitempty"`
}

func initSystemPeotry() {
	peotriesJson, err := ioutil.ReadFile("data/sys-peotry.json")
	if err != nil {
		fmt.Println("read sys-peotry.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(peotriesJson)
	re.ForEach(func(key, value gjson.Result) bool {
		uId := value.Get("u_id").Int()
		sId := value.Get("s_id").Int()
		pTitle := value.Get("p_title").String()
		pTime := value.Get("p_time").String()
		pContent := value.Get("p_content").String()
		pEnd := value.Get("p_end").String()
		pImages := value.Get("p_images").String()
		SavePeotry(uId, int(sId), pTitle, pTime, pContent, pEnd, pImages)
		return true
	})
}

func SavePeotry(userId int64, setId int, title string, pTime string, content string, end string, images string) {
	curTime := time.Now().UnixNano() / 1e3
	peotry := Peotry{
		ID:				curTime,
		UID:			userId,
		SID:			setId,
		PTitle:			title,
		PTime:			helper.StrToTimeStamp(pTime),	
		PContent:		content,
		PEnd:			end,
	}

	err := dbOrmDefault.Model(&Peotry{}).Save(peotry).Error
	if err != nil {
		fmt.Println(err)
	} else {
		res := gjson.Parse(images)
		imgs := res.Array()
		l := len(imgs)
		if l > 0 {
			SavePeotryImage(curTime, images, l)
		}
	}
}

func QueryPeotry(id int64, setId int, page int, limit int, content string) ([]Peotry, error, int, int ,int, int) {
	list := make([]Peotry, 0)
	totalPage := 0
	count := 0
	currentPage := page
	pageIsEnd := 0

	if limit == 0 {
		limit = 10
	}
	fmt.Println(id, setId, limit, page, content)

	db := dbOrmDefault.Model(&Peotry{})

	if id > 0 {
		peotry := Peotry{
			ID:	id,
		}
		err := db.Preload("UUser").Preload("SSet").Preload("PImage").Find(&peotry).Error
		if err == nil {
			peotry.UUser.UToken = ""
			peotry.SSet.UUser = nil
			list = append(list, peotry)
		} else {
			return nil, err, 0, 0, 0, 0
		}
	} else {
		if setId > 0 {
			query := &Peotry{
				SID:	setId,
			}
			db = db.Where(query)
		}
		if len(content) > 1 {
			db = db.Where("p_content LIKE ?", "%" + content + "%")
		}
		db.Count(&count)
		db = db.Preload("UUser").Preload("SSet").Preload("PImage")
		err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error
	
		if err == nil {
			totalPage, pageIsEnd = helper.PageTotal(limit, page, int64(count))
		} else {
			return nil, err, 0, 0, 0, 0
		}
	}
	return list, nil, count, totalPage, currentPage, pageIsEnd
}

// type Peotry2 struct {
// 	Id       int64     `json:"id"`
// 	SId      *Peotryset2     `json:"set"`
// 	UId      *User2     `json:"user"`
// 	PTitle   string    `json:"title"`
// 	PTime    time.Time `json:"time"`
// 	PContent string    `json:"content"`
// 	PEnd     string    `json:"end"`
// 	IId  *Peotryimage    `json:"img"`
// }

// func (t *Peotry) TableName() string {
// 	return "peotry"
// }

// func init() {
// 	orm.RegisterModel(new(Peotry))
// }

// func initSystemPeotry() {
	
// }

// // AddPeotry insert a new Peotry into database and returns
// // last inserted Id on success.
// func AddPeotry(m *Peotry) (id int64, err error) {
// 	o := orm.NewOrm()
// 	id, err = o.Insert(m)
// 	return
// }

// // GetPeotryById retrieves Peotry by Id. Returns error if
// // Id doesn't exist
// func GetPeotryById(id int64) (v *Peotry, err error) {
// 	o := orm.NewOrm()
// 	v = &Peotry{Id: id}
// 	if err = o.Read(v); err == nil {
// 		return v, nil
// 	}
// 	return nil, err
// }

// // GetAllPeotry retrieves all Peotry matches certain condition. Returns empty list if
// // no records exist
// func GetAllPeotry(query map[string]string, fields []string, sortby []string, order []string,
// 	offset int64, limit int64) (ml []interface{}, err error) {
// 	o := orm.NewOrm()
// 	qs := o.QueryTable(new(Peotry))
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

// 	var l []Peotry
// 	qs = qs.OrderBy(sortFields...)
// 	//关联查询后的数据顺序不是按照简单查询的数据顺序
// 	if _, err = qs.Limit(limit, offset).RelatedSel().All(&l, fields...); err == nil {
// 		var ll []Peotry2
// 		for _, v := range l {
// 			vv := Peotry2 {
// 				Id: v.Id,
// 				SId: &Peotryset2{Id: v.SId.Id, SName: v.SId.SName},
// 				UId: &User2{Id: v.UId.Id, UName: v.UId.UName},
// 				PTitle: v.PTitle,
// 				PTime: v.PTime,
// 				PContent: v.PContent,
// 				PEnd: v.PEnd,
// 				IId: v.IId,
// 			}

// 			// imageV, err := GetPeotryimageById(v.Id)
// 			// if err == nil {
// 			// 	vv.PImages = imageV.IImages
// 			// }
// 			ll = append(ll, vv)
// 		}

// 		if len(fields) == 0 {
// 			for _, v := range ll {
// 				ml = append(ml, v)
// 			}
// 		} else {
// 			// trim unused fields
// 			for _, v := range ll {
// 				m := make(map[string]interface{})

// 				val := reflect.ValueOf(v)
// 				typ := reflect.TypeOf(&v)	// 反射获取struct中的tag
// 				for _, fname := range fields {
// 					typElem := typ.Elem()
// 					field, ok := typElem.FieldByName(fname)
// 					if ok {
// 						m[field.Tag.Get("json")] = val.FieldByName(fname).Interface()
// 					} else {
// 						m[fname] = val.FieldByName(fname).Interface()
// 					}
// 				}
// 				ml = append(ml, m)
// 			}
// 		}
// 		return ml, nil
// 	}
// 	return nil, err
// }

// // UpdatePeotry updates Peotry by Id and returns error if
// // the record to be updated doesn't exist
// func UpdatePeotryById(m *Peotry) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotry{Id: m.Id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Update(m); err == nil {
// 			fmt.Println("Number of records updated in database:", num)
// 		}
// 	}
// 	return
// }

// // DeletePeotry deletes Peotry by Id and returns error if
// // the record to be deleted doesn't exist
// func DeletePeotry(id int64) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotry{Id: id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Delete(&Peotry{Id: id}); err == nil {
// 			fmt.Println("Number of records deleted in database:", num)
// 		}
// 	}
// 	return
// }
