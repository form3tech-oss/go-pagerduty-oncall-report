package report

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type csvReport struct {
	currency string
	outPath  string
}

func NewCsvReport(currency string, outPath string) Writer {
	return &csvReport{
		currency: strings.TrimSpace(currency),
		outPath:  outPath,
	}
}

func (r *csvReport) GenerateReport(data *PrintableData) (string, error) {

	fmt.Println(separator)
	fmt.Println(fmt.Sprintf("| Generating report(s) from '%s' to '%s'", data.Start.Format("Mon Jan _2 15:04:05 2006"), data.End.Add(time.Second*-1).Format("Mon Jan _2 15:04:05 2006")))
	fmt.Println(separator)

	header := []string{"User", "Email",
		"Weekday Hours", "Weekday Days", "Weekend Hours", "Weekend Days", "Bank Holiday Hours", "Bank Holiday Days",
		"Total Weekday Amount (" + r.currency + ")", "Total Weekend Amount (" + r.currency + ")",
		"Total Bank Holiday Amount (" + r.currency + ")", "Total  Amount (" + r.currency + ")"}

	for _, scheduleData := range data.SchedulesData {
		err := r.writeSingleRotation(scheduleData, data, header)
		if err != nil {
			log.Println("Error creating report for rotation: ", scheduleData.Name, " ID: ", scheduleData.ID, err)
			return "", err
		}
	}

	filename := fmt.Sprintf("%s/pagerduty_oncall_report.%d-%d-Summary.csv", r.outPath, data.Start.Month(), data.Start.Year())
	_ = os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating report file: ", filename, err)
		return "", err
	}
	defer file.Close()
	w := csv.NewWriter(file)

	if err := w.Write(header); err != nil {
		log.Println("error writing record to csv:", err)
		return "", err

	}

	sort.Slice(data.UsersSchedulesSummary, func(i, j int) bool {
		return strings.Compare(data.UsersSchedulesSummary[i].Name, data.UsersSchedulesSummary[j].Name) < 1
	})

	for _, userData := range data.UsersSchedulesSummary {
		err := writeUser(userData, w)
		if err != nil {
			log.Println("error writing user record to csv: ", filename, " user: ", userData.Name, " err: ", err)
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal("Error flushing writr", err)
		return "", err
	}
	return fmt.Sprintf("Report successfully generated: file://%s", filename), nil
}

func (r *csvReport) writeSingleRotation(scheduleData *ScheduleData, data *PrintableData, header []string) error {
	fmt.Println(separator)
	fmt.Println(fmt.Sprintf("| Writing Schedule: '%s' (%s)", scheduleData.Name, scheduleData.ID))
	fmt.Println(fmt.Sprintf("| Time Range: %s to %s", scheduleData.StartDate.Format(time.RFC822), scheduleData.EndDate.Format(time.RFC822)))
	fmt.Println(separator)
	noSpaceName := strings.Replace(scheduleData.Name, " ", "_", -1)

	filename := fmt.Sprintf("%s/pagerduty_oncall_report.%d-%d-%s-%s.csv", r.outPath, data.Start.Month(), data.Start.Year(), noSpaceName, scheduleData.ID)
	_ = os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating report file: ", filename, err)
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)

	if err := w.Write(header); err != nil {
		log.Println("error writing record to csv: ", filename, " err: ", err)
		return err

	}
	sort.Slice(scheduleData.RotaUsers, func(i, j int) bool {
		return strings.Compare(scheduleData.RotaUsers[i].Name, scheduleData.RotaUsers[j].Name) < 1
	})

	for _, userData := range scheduleData.RotaUsers {
		err := writeUser(userData, w)
		if err != nil {
			log.Println("error writing user record to csv: ", filename, " user: ", userData.Name, " err: ", err)
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal("Error flushing writr", err)
		return err
	}
	log.Println(fmt.Sprintf("Report successfully generated: file://%s", filename))
	return nil
}

func writeUser(userData *ScheduleUser, w *csv.Writer) error {
	dat := []string{userData.Name, userData.EmailAddress,
		fmt.Sprintf("%v", userData.NumWorkHours),
		fmt.Sprintf("%.1f", userData.NumWorkDays),
		fmt.Sprintf("%v", userData.NumWeekendHours),
		fmt.Sprintf("%.1f", userData.NumWeekendDays),
		fmt.Sprintf("%v", userData.NumBankHolidaysHours),
		fmt.Sprintf("%.1f", userData.NumBankHolidaysDays),
		fmt.Sprintf("%v", userData.TotalAmountWorkHours),
		fmt.Sprintf("%v", userData.TotalAmountWeekendHours),
		fmt.Sprintf("%v", userData.TotalAmountBankHolidaysHours),
		fmt.Sprintf("%.2f", userData.TotalAmount)}
	if err := w.Write(dat); err != nil {
		log.Println("error writing record to csv:", err)
		return err
	}
	return nil
}
