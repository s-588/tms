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
	Day          string `json:"day"`
	Month        string `json:"month"`
	Year         string `json:"year"`
	OrderID      string `json:"orderID"`
	User         string `json:"user"`         
	SourceNode         string `json:"source"`         
	DestinationNode         string `json:"destination"`         
	Weight       string `json:"weight"`       
	Price        string `json:"price"`        
	PriceWithNDS string `json:"priceWithNDS"` 
	NDS          string `json:"nds"`          
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
		OrderID: fmt.Sprintf("%d", order.OrderID),
		Day:     fmt.Sprintf("%d", order.CreatedAt.Day()),
		Month:   russianMonths[int(order.CreatedAt.Month())],
		Year:    fmt.Sprintf("%d", order.CreatedAt.Year()),
		User:    order.ClientName, // Заказчик в акте
		Weight:  fmt.Sprintf("%d", order.Weight),
		SourceNode: order.NodeStartName,
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
	w.Write(bytes)
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
	w.Write(bytes)
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

// GenerateOrdersReport создаёт Excel-файл с отчётом и графиками
func GenerateOrdersReport(stats OrderStats) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// --------------------------------------------------
	// Лист 1 — Основной отчёт (таблица + сводка)
	// --------------------------------------------------
	sheet := "Отчёт по заказам"
	f.SetSheetName(f.GetSheetName(0), sheet)
	// Заголовки сводки
	f.SetCellValue(sheet, "A1", "Отчёт по заказам")

	periodStr := ""
	switch {
	case !stats.Period.Start.IsZero() && !stats.Period.End.IsZero():
		periodStr = fmt.Sprintf("Период: %s — %s",
			stats.Period.Start.Format("02.01.2006"),
			stats.Period.End.Format("02.01.2006"))
	case !stats.Period.Start.IsZero():
		periodStr = fmt.Sprintf("Период: с %s", stats.Period.Start.Format("02.01.2006"))
	case !stats.Period.End.IsZero():
		periodStr = fmt.Sprintf("Период: по %s", stats.Period.End.Format("02.01.2006"))
	default:
		periodStr = "Период: все заказы"
	}
	f.SetCellValue(sheet, "A2", periodStr)

	row := 4
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Всего заказов")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.TotalOrders)
	row++

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Общая выручка (без НДС)")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.TotalRevenue.StringFixed(2)+" BYN")
	row++

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Средний чек")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stats.AvgPrice.StringFixed(2)+" BYN")
	row++

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Среднее расстояние")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%.1f км", stats.AvgDistance))
	row++

	// Таблица заказов
	row += 2
	headers := []string{"№", "Дата", "Клиент", "Сотрудник", "Груз", "Вес, кг", "Расстояние, км", "Сумма, BYN", "Статус"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, row)
		f.SetCellValue(sheet, cell, h)
	}
	f.SetRowHeight(sheet, row, 20)

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"D9EAD3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	f.SetCellStyle(sheet, "A"+fmt.Sprint(row), "I"+fmt.Sprint(row), styleHeader)

	row++
	for i, o := range stats.Orders {
		colA := fmt.Sprintf("A%d", row)
		f.SetCellValue(sheet, colA, i+1)
		f.SetCellValue(sheet, "B"+fmt.Sprint(row), o.CreatedAt.Format("02.01.2006"))
		f.SetCellValue(sheet, "C"+fmt.Sprint(row), o.ClientName)
		f.SetCellValue(sheet, "D"+fmt.Sprint(row), o.EmployeeName)
		f.SetCellValue(sheet, "E"+fmt.Sprint(row), o.PriceCargoType)
		f.SetCellValue(sheet, "F"+fmt.Sprint(row), o.Weight)
		f.SetCellValue(sheet, "G"+fmt.Sprint(row), fmt.Sprintf("%.1f", o.Distance))
		f.SetCellValue(sheet, "H"+fmt.Sprint(row), o.TotalPrice.StringFixed(2))
		f.SetCellValue(sheet, "I"+fmt.Sprint(row), string(o.Status))
		row++
	}

	f.SetColWidth(sheet, "A", "I", 14)

	// --------------------------------------------------
	// Лист 2 — Статистика и графики
	// --------------------------------------------------
	sheetStats := "Статистика"
	f.NewSheet(sheetStats)

	// 1. Заказы по месяцам (данные)
	f.SetCellValue(sheetStats, "A1", "Месяц")
	f.SetCellValue(sheetStats, "B1", "Кол-во заказов")
	row = 2
	for month, cnt := range stats.OrdersByMonth {
		f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), month)
		f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), cnt)
		row++
	}
	monthLastRow := row - 1

	categories := fmt.Sprintf(`%s!$A$2:$A$%d`, sheetStats, monthLastRow)
	values := fmt.Sprintf(`%s!$B$2:$B$%d`, sheetStats, monthLastRow)

	// Линейный график — заказы по месяцам
	if err := f.AddChart(sheetStats, "D2", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       "Заказы",
				Categories: categories,
				Values:     values,
				Line: excelize.ChartLine{
					Fill: excelize.Fill{Color: []string{"2F5597"}},
				},
			},
		},
		Title: []excelize.RichTextRun{{Text: "Динамика заказов по месяцам"}},
		Legend: excelize.ChartLegend{Position: "bottom"},
		PlotArea: excelize.ChartPlotArea{ShowLeaderLines: true},
	}); err != nil {
		return nil, err
	}

	// === НОВЫЙ ГРАФИК 1: Столбчатая диаграмма заказов по месяцам (справа) ===
	if err := f.AddChart(sheetStats, "J2", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Заказы",
				Categories: categories,
				Values:     values,
			},
		},
		Title: []excelize.RichTextRun{{Text: "Заказы по месяцам (столбцы)"}},
		Legend: excelize.ChartLegend{Position: "bottom"},
	}); err != nil {
		return nil, err
	}

	// 2. Заказы по статусам (данные + круговая)
	f.SetCellValue(sheetStats, "A"+fmt.Sprint(row+2), "Статус")
	f.SetCellValue(sheetStats, "B"+fmt.Sprint(row+2), "Количество")
	row += 3
	startPie := row
	for status, cnt := range stats.OrdersByStatus {
		f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), string(status))
		f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), cnt)
		row++
	}

	categories = fmt.Sprintf(`%s!$A$%d:$A$%d`, sheetStats, startPie, row-1)
	values = fmt.Sprintf(`%s!$B$%d:$B$%d`, sheetStats, startPie, row-1)

	if err := f.AddChart(sheetStats, "D"+fmt.Sprint(startPie), &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{
			{
				Name:       "Статусы",
				Categories: categories,
				Values:     values,
			},
		},
		Title: []excelize.RichTextRun{{Text: "Распределение по статусам"}},
		Legend: excelize.ChartLegend{Position: "right"},
	}); err != nil {
		return nil, err
	}

	// === НОВЫЙ ГРАФИК 2: Топ-клиенты по выручке (горизонтальная столбчатая) ===
	f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row+2), "Клиент")
	f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row+2), "Выручка (BYN)")
	f.SetCellValue(sheetStats, fmt.Sprintf("C%d", row+2), "Кол-во заказов")
	row += 3
	startTop := row

	for _, tc := range stats.TopClients {
		f.SetCellValue(sheetStats, fmt.Sprintf("A%d", row), tc.Name)
		revFloat := tc.Revenue.InexactFloat64() // числовое значение для графика
		f.SetCellValue(sheetStats, fmt.Sprintf("B%d", row), revFloat)
		f.SetCellValue(sheetStats, fmt.Sprintf("C%d", row), tc.Count)
		row++
	}

	if len(stats.TopClients) > 0 {
		categories = fmt.Sprintf(`%s!$A$%d:$A$%d`, sheetStats, startTop, row-1)
		values = fmt.Sprintf(`%s!$B$%d:$B$%d`, sheetStats, startTop, row-1)

		if err := f.AddChart(sheetStats, fmt.Sprintf("J%d", startTop-1), &excelize.Chart{
			Type: excelize.Bar, // горизонтальная — удобно для длинных названий клиентов
			Series: []excelize.ChartSeries{
				{
					Name:       "Выручка",
					Categories: categories,
					Values:     values,
				},
			},
			Title: []excelize.RichTextRun{{Text: "Топ клиентов по выручке"}},
			Legend: excelize.ChartLegend{Position: "bottom"},
		}); err != nil {
			return nil, err
		}
	}

	// Авто-ширина на листе статистики
	f.SetColWidth(sheetStats, "A", "C", 25)

	// Сохраняем в память
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}