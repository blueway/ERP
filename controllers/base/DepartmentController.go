package base

import (
	"encoding/json"
	md "goERP/models"
	"strconv"
	"strings"
)

type DepartmentController struct {
	BaseController
}

func (ctl *DepartmentController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}
func (ctl *DepartmentController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/product/department/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if department, err := md.GetDepartmentById(idInt64); err == nil {
			if err := ctl.ParseForm(&department); err == nil {
				if parentId, err := ctl.GetInt64("parent"); err == nil {
					if parent, err := md.GetDepartmentById(parentId); err == nil {
						department.Parent = parent
					}
				}
				if err := md.UpdateDepartmentById(department); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}
func (ctl *DepartmentController) Get() {
	ctl.PageName = "部门管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()
	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction
	ctl.URL = "/product/department/"
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuDepartmentActive"] = "active"
}
func (ctl *DepartmentController) Edit() {
	id := ctl.Ctx.Input.Param(":id")

	departmentInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {

			if department, err := md.GetDepartmentById(idInt64); err == nil {
				ctl.PageAction = department.Name
				departmentInfo["name"] = department.Name
				parent := make(map[string]interface{})
				if department.Parent != nil {
					parent["id"] = department.Parent.Id
					parent["name"] = department.Parent.Name
				}
				departmentInfo["parent"] = parent

			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordId"] = id

	ctl.Data["Category"] = departmentInfo
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_department_form.html"
}

func (ctl *DepartmentController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

//post请求创建产品分类
func (ctl *DepartmentController) PostCreate() {
	department := new(md.Department)
	if err := ctl.ParseForm(department); err == nil {
		if parentID, err := ctl.GetInt64("parent"); err == nil {
			if parent, err := md.GetDepartmentById(parentID); err == nil {
				department.Parent = parent
			}
		}
		if id, err := md.AddDepartment(department); err == nil {
			ctl.Redirect("/product/department/"+strconv.FormatInt(id, 10)+"?action=detail", 302)
		} else {
			ctl.Get()
		}
	} else {
		ctl.Get()
	}
}
func (ctl *DepartmentController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_department_form.html"
}
func (ctl *DepartmentController) Validator() {
	name := ctl.GetString("name")
	recordID, _ := ctl.GetInt64("recordId")
	name = strings.TrimSpace(name)
	result := make(map[string]bool)
	obj, err := md.GetDepartmentByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.Id {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// 获得符合要求的城市数据
func (ctl *DepartmentController) productCategoryList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.Department
	paginator, arrs, err := md.GetAllDepartment(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		// result["recordsFiltered"] = paginator.TotalCount
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			if line.Parent != nil {
				oneLine["parent"] = line.Parent.Name
			} else {
				oneLine["parent"] = "-"
			}

			oneLine["Id"] = line.Id
			oneLine["id"] = line.Id
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}
func (ctl *DepartmentController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.productCategoryList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *DepartmentController) GetList() {
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-department"
	ctl.TplName = "base/base_list_view.html"
}