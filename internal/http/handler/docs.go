package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	gotemplatedocx "github.com/JJJJJJack/go-template-docx"
	"github.com/s-588/tms/cmd/models"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type ReportData struct {
	Day             string `json:"day"`
	Month           string `json:"month"`
	Year            string `json:"year"`
	OrderID         string `json:"orderID"`
	User            string `json:"user"`
	SourceNode      string `json:"source"`
	DestinationNode string `json:"destination"`
	Weight          string `json:"weight"`
	Price           string `json:"price"`
	PriceWithNDS    string `json:"priceWithNDS"`
	NDS             string `json:"nds"`
}

var russianMonths = []string{
	"", "января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func (h Handler) buildReportData(ctx context.Context, orderID int32) (ReportData, error) {
	order, err := h.DB.GetOrderByID(ctx, orderID)
	if err != nil {
		return ReportData{}, fmt.Errorf("заказ не найден: %w", err)
	}

	data := ReportData{
		OrderID:         fmt.Sprintf("%d", order.OrderID),
		Day:             fmt.Sprintf("%d", order.CreatedAt.Day()),
		Month:           russianMonths[int(order.CreatedAt.Month())],
		Year:            fmt.Sprintf("%d", order.CreatedAt.Year()),
		User:            order.ClientName, // Заказчик в акте
		Weight:          fmt.Sprintf("%d", order.Weight),
		SourceNode:      order.NodeStartName,
		DestinationNode: order.NodeEndName,
	}

	// Расчёт НДС 20%
	price := order.TotalPrice
	nds := price.Mul(decimal.NewFromInt(20)).Div(decimal.NewFromInt(100)).Round(2)
	priceWithNDS := price.Add(nds).Round(2)

	data.Price = price.StringFixed(2)
	data.NDS = nds.StringFixed(2)
	data.PriceWithNDS = priceWithNDS.StringFixed(2)

	return data, nil
}

func generateBytes(templatePath string, data ReportData) ([]byte, error) {
	tmpl, err := gotemplatedocx.NewDocxTemplateFromFilename(templatePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить шаблон %s: %w", templatePath, err)
	}

	if err := tmpl.Apply(data); err != nil {
		return nil, fmt.Errorf("не удалось применить данные: %w", err)
	}

	return tmpl.Bytes(), nil
}

func (h Handler) DownloadContract(w http.ResponseWriter, r *http.Request) {
	orderID, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't get id from request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.buildReportData(r.Context(), orderID)
	if err != nil {
		slog.Error("failed to get order for contract", "orderID", orderID, "error", err)
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	bytes, err := generateBytes("templates/contract-template.docx", data)
	if err != nil {
		slog.Error("Ошибка генерации договора", "error", err)
		http.Error(w, "Ошибка генерации договора: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("договор_%s.docx", data.OrderID)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.PathEscape(filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("can't write response", "error", err)
	}
}

func (h Handler) DownloadAct(w http.ResponseWriter, r *http.Request) {
	orderID, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't get id from request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.buildReportData(r.Context(), orderID)
	if err != nil {
		slog.Error("failed to get order for act", "orderID", orderID, "error", err)
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	bytes, err := generateBytes("templates/act-template.docx", data)
	if err != nil {
		slog.Error("Ошибка генерации акта", "error", err)
		http.Error(w, "Ошибка генерации акта: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("акт_%s.docx", data.OrderID)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.PathEscape(filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("can't write response", "error", err)
	}
}

type ReportPeriod struct {
	Start time.Time
	End   time.Time
}

type TopClient struct {
	Name    string
	Revenue decimal.Decimal
	Count   int
}

type OrderStats struct {
	Period         ReportPeriod
	Orders         []models.Order
	TotalOrders    int
	TotalRevenue   decimal.Decimal
	AvgPrice       decimal.Decimal
	AvgDistance    float64
	AvgWeight      float64
	OrdersByStatus map[models.OrderStatus]int
	OrdersByMonth  map[string]int // "2025-03"
	TopClients     []TopClient
}

func GenerateOrdersReport(stats OrderStats) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("can't close report file", "error", err)
		}
	}()

	// Sheet 1: main orders report with summary and table.
	err := writeSheet1(f, stats)
	if err != nil {
		return nil, err
	}

	// Sheet 2: statistics and charts.
	err = writeSheet2(f, stats)
	if err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSheet1(f *excelize.File, stats OrderStats) error {
	sheet := "Orders Report"
	if err := f.SetSheetName(f.GetSheetName(0), sheet); err != nil {
		slog.Error("can't set sheet name", "name", sheet, "error", err)
		return err
	}

	if err := f.SetCellValue(sheet, "A1", "Orders Report"); err != nil {
		slog.Error("can't set cell value", "cell", "A1", "value", "Orders Report", "error", err)
		return err
	}

	row, err := writeOrdersSummary(stats, f, sheet)
	if err != nil {
		return err
	}

	// Orders table headers.
	err = writeOrdersTable(f, sheet, stats, row)
	if err != nil {
		return err
	}

	if err := f.SetColWidth(sheet, "A", "I", 14); err != nil {
		slog.Error("can't set column width", "sheet", sheet, "error", err)
		return err
	}
	return nil
}

func writeOrdersTable(f *excelize.File, sheet string, stats OrderStats, row int) error {
	row += 2
	headers := []string{"No.", "Date", "Client", "Employee", "Cargo", "Weight, kg", "Distance, km", "Amount, BYN", "Status"}
	for col, h := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, row)
		if err != nil {
			slog.Error("can't convert coordinates to cell name", "col", col+1, "row", row, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			slog.Error("can't set cell value", "cell", cell, "value", h, "error", err)
			return err
		}
	}
	if err := f.SetRowHeight(sheet, row, 20); err != nil {
		slog.Error("can't set row height", "row", row, "error", err)
		return err
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"D9EAD3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	if err != nil {
		slog.Error("can't create header style", "error", err)
		return err
	}
	if err := f.SetCellStyle(sheet, "A"+fmt.Sprint(row), "I"+fmt.Sprint(row), styleHeader); err != nil {
		slog.Error("can't set cell style", "row", row, "error", err)
		return err
	}

	row++
	err = writeOrdersTableData(stats, f, sheet, row)
	return err
}

func writeOrdersTableData(stats OrderStats, f *excelize.File, sheet string, row int) error {
	for i, o := range stats.Orders {
		if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", i+1, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "B"+fmt.Sprint(row), o.CreatedAt.Format("02.01.2006")); err != nil {
			slog.Error("can't set cell value", "cell", "B"+fmt.Sprint(row), "value", o.CreatedAt.Format("02.01.2006"), "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "C"+fmt.Sprint(row), o.ClientName); err != nil {
			slog.Error("can't set cell value", "cell", "C"+fmt.Sprint(row), "value", o.ClientName, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "D"+fmt.Sprint(row), o.EmployeeName); err != nil {
			slog.Error("can't set cell value", "cell", "D"+fmt.Sprint(row), "value", o.EmployeeName, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "E"+fmt.Sprint(row), o.PriceCargoType); err != nil {
			slog.Error("can't set cell value", "cell", "E"+fmt.Sprint(row), "value", o.PriceCargoType, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "F"+fmt.Sprint(row), o.Weight); err != nil {
			slog.Error("can't set cell value", "cell", "F"+fmt.Sprint(row), "value", o.Weight, "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "G"+fmt.Sprint(row), fmt.Sprintf("%.1f", o.Distance)); err != nil {
			slog.Error("can't set cell value", "cell", "G"+fmt.Sprint(row), "value", fmt.Sprintf("%.1f", o.Distance), "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "H"+fmt.Sprint(row), o.TotalPrice.StringFixed(2)); err != nil {
			slog.Error("can't set cell value", "cell", "H"+fmt.Sprint(row), "value", o.TotalPrice.StringFixed(2), "error", err)
			return err
		}
		if err := f.SetCellValue(sheet, "I"+fmt.Sprint(row), string(o.Status)); err != nil {
			slog.Error("can't set cell value", "cell", "I"+fmt.Sprint(row), "value", string(o.Status), "error", err)
			return err
		}
		row++
	}
	return nil
}

func writeOrdersSummary(stats OrderStats, f *excelize.File, sheet string) (row int, err error) {
	var periodStr string
	switch {
	case !stats.Period.Start.IsZero() && !stats.Period.End.IsZero():
		periodStr = fmt.Sprintf("Period: %s — %s",
			stats.Period.Start.Format("02.01.2006"),
			stats.Period.End.Format("02.01.2006"))
	case !stats.Period.Start.IsZero():
		periodStr = fmt.Sprintf("Period: from %s", stats.Period.Start.Format("02.01.2006"))
	case !stats.Period.End.IsZero():
		periodStr = fmt.Sprintf("Period: until %s", stats.Period.End.Format("02.01.2006"))
	default:
		periodStr = "Period: all orders"
	}
	if err := f.SetCellValue(sheet, "A2", periodStr); err != nil {
		slog.Error("can't set cell value", "cell", "A2", "value", periodStr, "error", err)
		return row, err
	}

	row = 4
	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total orders"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", "Total orders", "error", err)
		return row, err
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.TotalOrders); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", stats.TotalOrders, "error", err)
		return row, err
	}
	row++

	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total revenue (excl. VAT)"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", "Total revenue (excl. VAT)", "error", err)
		return row, err
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.TotalRevenue.StringFixed(2)+" BYN"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", stats.TotalRevenue.StringFixed(2)+" BYN", "error", err)
		return row, err
	}
	row++

	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Average check"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", "Average check", "error", err)
		return row, err
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.AvgPrice.StringFixed(2)+" BYN"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", stats.AvgPrice.StringFixed(2)+" BYN", "error", err)
		return row, err
	}
	row++

	if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Average distance"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", "Average distance", "error", err)
		return row, err
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%.1f km", stats.AvgDistance)); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", fmt.Sprintf("%.1f km", stats.AvgDistance), "error", err)
		return row, err
	}
	row++
	return row, nil
}

func writeSheet2(f *excelize.File, stats OrderStats) error {
	sheetStats := "Statistics"
	if _, err := f.NewSheet(sheetStats); err != nil {
		slog.Error("can't create statistics sheet", "name", sheetStats, "error", err)
		return err
	}

	// Orders by month data.
	if err := f.SetCellValue(sheetStats, "A1", "Month"); err != nil {
		slog.Error("can't set cell value", "cell", "A1", "value", "Month", "error", err)
		return err
	}
	if err := f.SetCellValue(sheetStats, "B1", "Order count"); err != nil {
		slog.Error("can't set cell value", "cell", "B1", "value", "Order count", "error", err)
		return err
	}
	row := 2
	for month, cnt := range stats.OrdersByMonth {
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), month); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", month, "error", err)
			return err
		}
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), cnt); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", cnt, "error", err)
			return err
		}
		row++
	}
	monthLastRow := row - 1

	categories := fmt.Sprintf(`%s!$A$2:$A$%d`, sheetStats, monthLastRow)
	values := fmt.Sprintf(`%s!$B$2:$B$%d`, sheetStats, monthLastRow)

	// Line chart: orders by month.
	err := writeOrdersLineChart(f, sheetStats, categories, values)
	if err != nil {
		return err
	}

	// Column chart: orders by month.
	err = writeOrdersColumnChart(f, sheetStats, categories, values)
	if err != nil {
		return err
	}

	// Pie chart: Orders by status data.
	row, err = writeOrdersPieCart(f, sheetStats, row, stats)
	if err != nil {
		return err
	}

	// Top clients by revenue data and bar chart.
	err = writeOrdersBarChart(f, sheetStats, row, stats)
	if err != nil {
		return err
	}

	if err := f.SetColWidth(sheetStats, "A", "C", 25); err != nil {
		slog.Error("can't set column width", "sheet", sheetStats, "error", err)
		return err
	}
	return nil
}

func writeOrdersBarChart(f *excelize.File, sheetStats string, row int, stats OrderStats) error {
	if err := f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row+2), "Client"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row+2), "value", "Client", "error", err)
		return err
	}
	if err := f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row+2), "Revenue (BYN)"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row+2), "value", "Revenue (BYN)", "error", err)
		return err
	}
	if err := f.SetCellValue(sheetStats, fmt.Sprintf("C%d", row+2), "Order count"); err != nil {
		slog.Error("can't set cell value", "cell", fmt.Sprintf("C%d", row+2), "value", "Order count", "error", err)
		return err
	}
	row += 3
	startTop := row

	for _, tc := range stats.TopClients {
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), tc.Name); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", tc.Name, "error", err)
			return err
		}
		revFloat := tc.Revenue.InexactFloat64()
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), revFloat); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", revFloat, "error", err)
			return err
		}
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("C%d", row), tc.Count); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("C%d", row), "value", tc.Count, "error", err)
			return err
		}
		row++
	}

	if len(stats.TopClients) > 0 {
		categories := fmt.Sprintf(`%s!$A$%d:$A$%d`, sheetStats, startTop, row-1)
		values := fmt.Sprintf(`%s!$B$%d:$B$%d`, sheetStats, startTop, row-1)

		if err := f.AddChart(sheetStats, fmt.Sprintf("J%d", startTop-1), &excelize.Chart{
			Type: excelize.Bar,
			Series: []excelize.ChartSeries{
				{
					Name:       "Revenue",
					Categories: categories,
					Values:     values,
				},
			},
			Title:  []excelize.RichTextRun{{Text: "Top clients by revenue"}},
			Legend: excelize.ChartLegend{Position: "bottom"},
		}); err != nil {
			return err
		}
	}
	return nil
}

func writeOrdersPieCart(f *excelize.File, sheetStats string, row int, stats OrderStats) (int, error) {
	if err := f.SetCellValue(sheetStats, "A"+fmt.Sprint(row+2), "Status"); err != nil {
		slog.Error("can't set cell value", "cell", "A"+fmt.Sprint(row+2), "value", "Status", "error", err)
		return row, err
	}
	if err := f.SetCellValue(sheetStats, "B"+fmt.Sprint(row+2), "Count"); err != nil {
		slog.Error("can't set cell value", "cell", "B"+fmt.Sprint(row+2), "value", "Count", "error", err)
		return row, err
	}
	row += 3
	startPie := row
	for status, cnt := range stats.OrdersByStatus {
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), string(status)); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("A%d", row), "value", string(status), "error", err)
			return row, err
		}
		if err := f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), cnt); err != nil {
			slog.Error("can't set cell value", "cell", fmt.Sprintf("B%d", row), "value", cnt, "error", err)
			return row, err
		}
		row++
	}

	categories := fmt.Sprintf(`%s!$A$%d:$A$%d`, sheetStats, startPie, row-1)
	values := fmt.Sprintf(`%s!$B$%d:$B$%d`, sheetStats, startPie, row-1)

	if err := f.AddChart(sheetStats, "D"+fmt.Sprint(startPie), &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{
			{
				Name:       "Statuses",
				Categories: categories,
				Values:     values,
			},
		},
		Title:  []excelize.RichTextRun{{Text: "Distribution by status"}},
		Legend: excelize.ChartLegend{Position: "right"},
	}); err != nil {
		return row, err
	}
	return row, nil
}

func writeOrdersColumnChart(f *excelize.File, sheetStats string, categories string, values string) error {
	if err := f.AddChart(sheetStats, "J2", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Orders",
				Categories: categories,
				Values:     values,
			},
		},
		Title:  []excelize.RichTextRun{{Text: "Orders by month (columns)"}},
		Legend: excelize.ChartLegend{Position: "bottom"},
	}); err != nil {
		return err
	}
	return nil
}

func writeOrdersLineChart(f *excelize.File, sheetStats string, categories string, values string) error {
	if err := f.AddChart(sheetStats, "D2", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       "Orders",
				Categories: categories,
				Values:     values,
				Line: excelize.ChartLine{
					Fill: excelize.Fill{Color: []string{"2F5597"}},
				},
			},
		},
		Title: []excelize.RichTextRun{{Text: "Orders dynamics by month"}},
		Legend: excelize.ChartLegend{
			Position:      "bottom",
			ShowLegendKey: false,
			Font:          &excelize.Font{},
		},
		PlotArea: excelize.ChartPlotArea{ShowLeaderLines: true},
	}); err != nil {
		return err
	}
	return nil
}
