package services

import (
	"application/models"
	"fmt"
	"time"
)

type ComplexReportRepository interface {
	Get(user_id int64, apireport *models.ApiReport) (complexreport *models.ApiComplexReport, err error)
}

type ComplexReportService struct {
	*Repository
}

func NewComplexReportService(repository *Repository) *ComplexReportService {
	return &ComplexReportService{Repository: repository}
}

func (complexreportservice *ComplexReportService) Get(user_id int64, apireport *models.ApiReport) (complexreport *models.ApiComplexReport, err error) {
	query := ""
	if len(apireport.Periods) != 0 {
		for i, period := range apireport.Periods {
			subquery := ""
			begin, _ := time.Parse(models.FORMAT_DATETIME, period.Begin)
			if !(begin.Year() == 1 && begin.Month() == 1 && begin.Day() == 1) {
				subquery += "o.created >= '" + begin.Format(models.FORMAT_DATE) + "'"
			}
			end, _ := time.Parse(models.FORMAT_DATETIME, period.End)
			if !(end.Year() == 1 && end.Month() == 1 && end.Day() == 1) {
				if subquery != "" {
					subquery += " and "
				}
				subquery += "o.created <= '" + end.Format(models.FORMAT_DATE) + "'"
			}
			if subquery != "" {
				query += "(" + subquery + ")"
			}
			if i != len(apireport.Periods)-1 {
				if subquery != "" {
					query += " or "
				}
			}
		}
		if query != "" {
			query = " and (" + query + ")"
		}
	}
	if len(apireport.Projects) != 0 {
		query += " and o.project_id in ("
		for i, project := range apireport.Projects {
			query += fmt.Sprintf("%v", project.Project_ID)
			if i != len(apireport.Projects)-1 {
				query += ","
			}
		}
		query += ")"
	}
	if len(apireport.Orders) != 0 {
		query += " and o.id in ("
		for i, order := range apireport.Orders {
			query += fmt.Sprintf("%v", order.Order_ID)
			if i != len(apireport.Orders)-1 {
				query += ","
			}
		}
		query += ")"
	}
	if len(apireport.Facilities) != 0 {
		query += " and o.service_id in ("
		for i, facility := range apireport.Facilities {
			query += fmt.Sprintf("%v", facility.Facility_ID)
			if i != len(apireport.Facilities)-1 {
				query += ","
			}
		}
		query += ")"
	}
	if len(apireport.ComplexStatuses) != 0 {
		query += " and r.complex_status_id in ("
		for i, complexstatus := range apireport.ComplexStatuses {
			query += fmt.Sprintf("%v", complexstatus.ComplexStatus_ID)
			if i != len(apireport.ComplexStatuses)-1 {
				query += ","
			}
		}
		query += ")"
	}
	if len(apireport.Suppliers) != 0 {
		query += " and o.supplier_id in ("
		for i, supplier := range apireport.Suppliers {
			query += fmt.Sprintf("%v", supplier.Supplier_ID)
			if i != len(apireport.Suppliers)-1 {
				query += ","
			}
		}
		query += ")"
	}
	if apireport.Settings.Field != "" && apireport.Settings.Order != "" {
		query += " order by " + apireport.Settings.Field + " " + apireport.Settings.Order
	}
	if apireport.Settings.Count > 0 {
		query += " limit " + fmt.Sprintf("%v", (apireport.Settings.Page-1)*apireport.Settings.Count) + ", " + fmt.Sprintf("%v", apireport.Settings.Count)
	}
	complexreport = new(models.ApiComplexReport)
	complexreport.Orders = *new([]models.ApiOrderReport)
	_, err = complexreportservice.DbContext.Select(&complexreport.Orders,
		"select o.project_id as projectId, p.name as projectName, o.id as orderId, o.name as orderName,"+
			" o.service_id as serviceId, coalesce(s.name, '') as serviceName, o.supplier_id as supplierId, coalesce(u.name, '') as supplierName,"+
			" r.complex_status_id as statusId, r.name as statusName, cast(o.begin_date as char) as dateBegin,"+
			" case o.end_date when '0001-01-01 00:00:00' then case o.begin_date when '0001-01-01 00:00:00' then cast(o.end_date as char)"+
			" else cast(date_add(o.begin_date, interval o.execution_forecast day) as char) end else cast(o.end_date as char) end as dateEnd,"+
			" o.charged_fee as budget from orders o left join projects p on o.project_id = p.id left join services s on o.service_id = s.id"+
			" left join units u on o.supplier_id = u.id left join order_complex_statuses r on o.id = r.order_id"+
			" where o.unit_id = (select unit_id from users where id = ?)"+query,
		user_id)
	if err != nil {
		log.Error("Error during getting complex report object from database %v", err)
		return nil, err
	}
	var total int64 = 0
	var budget float64 = 0
	services := make(map[int64]*models.DoughnutValues)
	complexstatuses := make(map[int64]*models.DoughnutValues)
	suppliers := make(map[int64]*models.DoughnutValues)
	for _, order := range complexreport.Orders {
		servicevalues, ok := services[order.Facility_ID]
		if !ok {
			servicevalues = new(models.DoughnutValues)
		}
		servicevalues.Count++
		servicevalues.Sum += order.Budget
		services[order.Facility_ID] = servicevalues

		complexstatusvalues, ok := complexstatuses[order.ComplexStatus_ID]
		if !ok {
			complexstatusvalues = new(models.DoughnutValues)
		}
		complexstatusvalues.Count++
		complexstatusvalues.Sum += order.Budget
		complexstatuses[order.ComplexStatus_ID] = complexstatusvalues

		suppliervalues, ok := suppliers[order.Supplier_ID]
		if !ok {
			suppliervalues = new(models.DoughnutValues)
		}
		suppliervalues.Count++
		suppliervalues.Sum += order.Budget
		suppliers[order.Supplier_ID] = suppliervalues

		total++
		budget += order.Budget
	}
	complexreport.Budgets = *new([]models.ApiBudgetedReport)
	complexreport.Facilities = *new([]models.ApiDoughnutReport)
	for id, value := range services {
		complexreport.Facilities = append(complexreport.Facilities, *models.NewApiDoughnutReport(id, 0,
			(float64(value.Count)/float64(total))*100, float64(value.Count)))
		if apireport.Budgeted == models.TYPE_BUDGETEDBY_FACILITY_VALUE || apireport.Budgeted == "" {
			var percentage float64 = 100 / float64(total)
			if budget != 0 {
				percentage = (float64(value.Sum) / float64(budget)) * 100
			}
			complexreport.Budgets = append(complexreport.Budgets, *models.NewApiBudgetedReport(models.TYPE_BUDGETEDBY_FACILITY_VALUE,
				*models.NewApiDoughnutReport(id, 0, percentage, models.Round(float64(value.Sum), 0.5, 2))))
		}

	}

	complexreport.ComplexStatuses = *new([]models.ApiDoughnutReport)
	for id, value := range complexstatuses {
		complexreport.ComplexStatuses = append(complexreport.ComplexStatuses, *models.NewApiDoughnutReport(id, 0,
			(float64(value.Count)/float64(total))*100, float64(value.Count)))
		if apireport.Budgeted == models.TYPE_BUDGETEDBY_COMPLEX_STATUS_VALUE || apireport.Budgeted == "" {
			var percentage float64 = 100 / float64(total)
			if budget != 0 {
				percentage = (float64(value.Sum) / float64(budget)) * 100
			}
			complexreport.Budgets = append(complexreport.Budgets, *models.NewApiBudgetedReport(models.TYPE_BUDGETEDBY_COMPLEX_STATUS_VALUE,
				*models.NewApiDoughnutReport(id, 0, percentage, models.Round(float64(value.Sum), 0.5, 2))))
		}
	}

	complexreport.Suppliers = *new([]models.ApiDoughnutReport)
	for id, value := range suppliers {
		complexreport.Suppliers = append(complexreport.Suppliers, *models.NewApiDoughnutReport(id, 0,
			(float64(value.Count)/float64(total))*100, float64(value.Count)))
		if apireport.Budgeted == models.TYPE_BUDGETEDBY_SUPPLIER_VALUE || apireport.Budgeted == "" {
			var percentage float64 = 100 / float64(total)
			if budget != 0 {
				percentage = (float64(value.Sum) / float64(budget)) * 100
			}
			complexreport.Budgets = append(complexreport.Budgets, *models.NewApiBudgetedReport(models.TYPE_BUDGETEDBY_SUPPLIER_VALUE,
				*models.NewApiDoughnutReport(id, 0, percentage, models.Round(float64(value.Sum), 0.5, 2))))
		}
	}

	return complexreport, nil
}
