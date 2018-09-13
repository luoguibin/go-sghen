package models

type PeotryImage struct {
	ID    		int64  	`orm:"column(id);pk" json:"id,omitempty"`
	IImages   	string  `orm:"column(i_images);" json:"images"`
	ICount 		int 	`orm:"column(i_count);" json:"count"`
}

func SavePeotryImage(id int64, images string, count int) error {
	peotryImage := &PeotryImage {
		ID:			id,
		IImages:	images,
		ICount:		count,
	}

	err := dbOrmDefault.Model(&PeotryImage{}).Save(peotryImage).Error
	return err
}

// func (t *Peotryimage) TableName() string {
// 	return "peotryimage"
// }

// func init() {
// 	orm.RegisterModel(new(Peotryimage))
// }

// // AddPeotryimage insert a new Peotryimage into database and returns
// // last inserted Id on success.
// func AddPeotryimage(m *Peotryimage) (id int64, err error) {
// 	o := orm.NewOrm()
// 	id, err = o.Insert(m)
// 	return
// }

// // GetPeotryimageById retrieves Peotryimage by Id. Returns error if
// // Id doesn't exist
// func GetPeotryimageById(id int64) (v *Peotryimage, err error) {
// 	o := orm.NewOrm()
// 	v = &Peotryimage{Id: id}
// 	if err = o.Read(v); err == nil {
// 		return v, nil
// 	}
// 	return nil, err
// }

// // GetAllPeotryimage retrieves all Peotryimage matches certain condition. Returns empty list if
// // no records exist
// func GetAllPeotryimage(query map[string]string, fields []string, sortby []string, order []string,
// 	offset int64, limit int64) (ml []interface{}, err error) {
// 	o := orm.NewOrm()
// 	qs := o.QueryTable(new(Peotryimage))
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

// 	var l []Peotryimage
// 	qs = qs.OrderBy(sortFields...)
// 	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
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

// // UpdatePeotryimage updates Peotryimage by Id and returns error if
// // the record to be updated doesn't exist
// func UpdatePeotryimageById(m *Peotryimage) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotryimage{Id: m.Id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Update(m); err == nil {
// 			fmt.Println("Number of records updated in database:", num)
// 		}
// 	}
// 	return
// }

// // DeletePeotryimage deletes Peotryimage by Id and returns error if
// // the record to be deleted doesn't exist
// func DeletePeotryimage(id int64) (err error) {
// 	o := orm.NewOrm()
// 	v := Peotryimage{Id: id}
// 	// ascertain id exists in the database
// 	if err = o.Read(&v); err == nil {
// 		var num int64
// 		if num, err = o.Delete(&Peotryimage{Id: id}); err == nil {
// 			fmt.Println("Number of records deleted in database:", num)
// 		}
// 	}
// 	return
// }

// // baseStr: dataUrl format
// func SavePeotryimage(baseStr string, rename string) (format string, err error) {
// 	if len(baseStr) == 0 {
// 		return "", errors.New("空数据")
// 	}

// 	baseIndex := strings.Index(baseStr, "base64")
// 	if baseIndex < 15 {
// 		return "", errors.New("数据错误")
// 	}

// 	format = baseStr[11:baseIndex - 1]

// 	data, err := base64.StdEncoding.DecodeString(baseStr[baseIndex + 7:])
// 	if err != nil {
// 		return "", err
// 	}

// 	isExist, _ := utils.PathExists(MConfig.ImageSavePath)
// 	if !isExist {
// 		isMade := utils.MkdirAll(MConfig.ImageSavePath)
// 		if !isMade {
// 			return "", errors.New("创建文件夹失败")
// 		}
// 	}

// 	err2 := ioutil.WriteFile(MConfig.ImageSavePath + rename + "." + format, data, 0666) 
// 	if err2 != nil {
// 		return "", err2
// 	} 
	
// 	return format, nil
// }
