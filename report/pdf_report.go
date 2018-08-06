package report

import (
	"fmt"
	"log"
	"time"

	"sort"

	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/mitchellh/go-homedir"
)

const (
	matrixRowFormat = "%-23s %8v %8v %11v %10v %10v %13v %9v"
)

type pdfReport struct {
	currency string
}

func NewPDFReport(currency string) Writer {
	return &pdfReport{
		currency: currency,
	}
}

func (r *pdfReport) GenerateReport(data *PrintableData) (string, error) {

	log.Println("Generating pdf report...")
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetTopMargin(30)
	pdf.SetHeaderFunc(func() {
		//pdf.Image(example.ImageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.CellFormat(0, 10,
			fmt.Sprintf("PagerDuty oncall report(s) from %s to %s ", data.Start.Format("02/01/2006"), data.End.Add(time.Second*-1).Format("02/01/2006")),
			"R", 0, "R", false, 0, "")
		pdf.Ln(20)
	})

	pdf.AddPage()

	for _, scheduleData := range data.SchedulesData {

		pdf.SetFont("Arial", "B", 13)
		pdf.CellFormat(0, 5,
			fmt.Sprintf("  Schedule: '%s' (%s)", scheduleData.Name, scheduleData.ID),
			"L", 0, "L", false, 0, "")
		pdf.Ln(8)

		pdf.SetFont("Courier", "B", 9)
		pdf.CellFormat(0, 5,
			fmt.Sprintf(matrixRowFormat, "USER", "WEEKDAY", "WEEKEND", "B. HOLIDAY", "WEEKDAY", "WEEKEND", "B. HOLIDAY", "TOTAL"),
			"", 0, "L", false, 0, "")
		pdf.Ln(3)
		pdf.CellFormat(0, 5,
			fmt.Sprintf(matrixRowFormat, "", "HOURS", "HOURS", "HOURS", "AMOUNT", "AMOUNT", "AMOUNT", "AMOUNT"),
			"", 0, "L", false, 0, "")
		pdf.Ln(3)
		pdf.CellFormat(0, 5,
			fmt.Sprintf(matrixRowFormat, "", "DAYS", "DAYS", "DAYS", "", "", "", ""),
			"B", 0, "L", false, 0, "")
		pdf.Ln(5)

		pdf.SetFont("Courier", "", 9)

		sort.Slice(scheduleData.RotaUsers, func(i, j int) bool {
			return strings.Compare(scheduleData.RotaUsers[i].Name, scheduleData.RotaUsers[j].Name) < 1
		})

		for _, userData := range scheduleData.RotaUsers {
			pdf.CellFormat(0, 5,
				fmt.Sprintf(matrixRowFormat, tr(userData.Name),
					fmt.Sprintf("%v h", userData.NumWorkHours),
					fmt.Sprintf("%v h", userData.NumWeekendHours),
					fmt.Sprintf("%v h", userData.NumBankHolidaysHours),
					tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountWorkHours)),
					tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountWeekendHours)),
					tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountBankHolidaysHours)),
					tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmount))),
				"", 0, "L", false, 0, "")
			pdf.Ln(3)
			pdf.CellFormat(0, 5,
				fmt.Sprintf(matrixRowFormat, "_______________________",
					fmt.Sprintf("%.1f d", userData.NumWorkDays),
					fmt.Sprintf("%.1f d", userData.NumWeekendDays),
					fmt.Sprintf("%.1f d", userData.NumBankHolidaysDays),
					"__________", "__________", "_____________", "_________"),
				"", 0, "L", false, 0, "")
			pdf.Ln(5)
		}

		pdf.Ln(10)
	}

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 5, "  Users summary",
		"L", 0, "L", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Courier", "B", 9)
	pdf.CellFormat(0, 5,
		fmt.Sprintf(matrixRowFormat, "USER", "WEEKDAY", "WEEKEND", "B. HOLIDAY", "WEEKDAY", "WEEKEND", "B. HOLIDAY", "TOTAL"),
		"", 0, "L", false, 0, "")
	pdf.Ln(3)
	pdf.CellFormat(0, 5,
		fmt.Sprintf(matrixRowFormat, "", "HOURS", "HOURS", "HOURS", "AMOUNT", "AMOUNT", "AMOUNT", "AMOUNT"),
		"", 0, "L", false, 0, "")
	pdf.Ln(3)
	pdf.CellFormat(0, 5,
		fmt.Sprintf(matrixRowFormat, "", "DAYS", "DAYS", "DAYS", "", "", "", ""),
		"B", 0, "L", false, 0, "")
	pdf.Ln(5)

	sort.Slice(data.UsersSchedulesSummary, func(i, j int) bool {
		return strings.Compare(data.UsersSchedulesSummary[i].Name, data.UsersSchedulesSummary[j].Name) < 1
	})

	pdf.SetFont("Courier", "", 9)
	for _, userData := range data.UsersSchedulesSummary {
		pdf.CellFormat(0, 5,
			fmt.Sprintf(matrixRowFormat, tr(userData.Name),
				fmt.Sprintf("%v h", userData.NumWorkHours),
				fmt.Sprintf("%v h", userData.NumWeekendHours),
				fmt.Sprintf("%v h", userData.NumBankHolidaysHours),
				tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountWorkHours)),
				tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountWeekendHours)),
				tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmountBankHolidaysHours)),
				tr(fmt.Sprintf("%s%v", r.currency, userData.TotalAmount))),
			"", 0, "L", false, 0, "")
		pdf.Ln(3)
		pdf.CellFormat(0, 5,
			fmt.Sprintf(matrixRowFormat, "_______________________",
				fmt.Sprintf("%.1f d", userData.NumWorkDays),
				fmt.Sprintf("%.1f d", userData.NumWeekendDays),
				fmt.Sprintf("%.1f d", userData.NumBankHolidaysDays),
				"__________", "__________", "_____________", "_________"),
			"", 0, "L", false, 0, "")
		pdf.Ln(5)
	}

	dir, _ := homedir.Dir()
	filename := fmt.Sprintf("%s/pagerduty_oncall_report.%d-%d.pdf", dir, data.Start.Month(), data.Start.Year())
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Report successfully generated: file://%s", filename), nil
}
