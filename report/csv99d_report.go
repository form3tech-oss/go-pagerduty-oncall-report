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

type csvOncallReport struct {
	outPath string
}

func NewCsvOncallReport(outPath string) Writer {
	return &csvOncallReport{
		outPath: outPath,
	}
}

func (r *csvOncallReport) GenerateReport(data *PrintableData) (string, error) {
	separator := strings.Repeat("-", 80)
	fmt.Println(separator)
	fmt.Printf("| Generating report(s) from '%s' to '%s'\n", data.Start.Format("Mon Jan _2 15:04:05 2006"), data.End.Add(-time.Second).Format("Mon Jan _2 15:04:05 2006"))
	fmt.Println(separator)

	header := []string{"User", "Weekday Hours", "Weekday Days", "Weekend Hours", "Weekend Days", "Holiday Hours", "Holiday Days", "Manager Approval"}

	for _, scheduleData := range data.SchedulesData {
		err := r.writeSingleRotation(scheduleData, data, header)
		if err != nil {
			log.Printf("Error creating report for rotation: %s, ID: %s, %v\n", scheduleData.Name, scheduleData.ID, err)
			return "", err
		}
	}

	filename := fmt.Sprintf("%s/pagerduty_oncall_report.%d-%d-Summary.csv", r.outPath, data.Start.Month(), data.Start.Year())
	_ = os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating report file: %s, %v\n", filename, err)
		return "", err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	// Write the first header with the start and end dates
	if err := w.Write([]string{"Start date:", data.Start.Format("Mon Jan _2 15:04:05 2006"), "End date:", data.End.Add(-time.Second).Format("Mon Jan _2 15:04:05 2006")}); err != nil {
		log.Printf("Error writing record to csv: %v\n", err)
		return "", err
	}

	if err := w.Write(header); err != nil {
		log.Printf("Error writing record to csv: %v\n", err)
		return "", err
	}

	sort.Slice(data.UsersSchedulesSummary, func(i, j int) bool {
		return strings.Compare(data.UsersSchedulesSummary[i].Name, data.UsersSchedulesSummary[j].Name) < 1
	})

	for _, userData := range data.UsersSchedulesSummary {
		err := writeUserData(userData, w)
		if err != nil {
			log.Printf("Error writing user record to csv: %s, user: %s, err: %v\n", filename, userData.Name, err)
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("Error flushing writer: %v\n", err)
		return "", err
	}
	return fmt.Sprintf("Report successfully generated: file://%s", filename), nil
}

func (r *csvOncallReport) writeSingleRotation(scheduleData *ScheduleData, data *PrintableData, header []string) error {
	separator := strings.Repeat("-", 80)
	fmt.Println(separator)
	fmt.Printf("| Writing Schedule: '%s' (%s)\n", scheduleData.Name, scheduleData.ID)
	fmt.Printf("| Time Range: %s to %s\n", scheduleData.StartDate.Format(time.RFC822), scheduleData.EndDate.Format(time.RFC822))
	fmt.Println(separator)

	noSpaceName := strings.ReplaceAll(scheduleData.Name, " ", "_")
	filename := fmt.Sprintf("%s/pagerduty_oncall_report.%d-%d-%s-%s.csv", r.outPath, data.Start.Month(), data.Start.Year(), noSpaceName, scheduleData.ID)

	_ = os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating report file: %s, %v\n", filename, err)
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	// Write the first header with the start and end dates
	if err := w.Write([]string{"Start date:", data.Start.Format("Mon Jan _2 15:04:05 2006"), "End date:", data.End.Add(-time.Second).Format("Mon Jan _2 15:04:05 2006")}); err != nil {
		log.Printf("Error writing record to csv: %v\n", err)
		return err
	}

	if err := w.Write(header); err != nil {
		log.Printf("Error writing record to csv: %s, err: %v\n", filename, err)
		return err
	}

	sort.Slice(scheduleData.RotaUsers, func(i, j int) bool {
		return strings.Compare(scheduleData.RotaUsers[i].Name, scheduleData.RotaUsers[j].Name) < 1
	})

	for _, userData := range scheduleData.RotaUsers {
		err := writeUserData(userData, w)
		if err != nil {
			log.Printf("Error writing user record to csv: %s, user: %s, err: %v\n", filename, userData.Name, err)
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("Error flushing writer: %v\n", err)
		return err
	}
	log.Printf("Report successfully generated: file://%s\n", filename)
	return nil
}

func writeUserData(userData *ScheduleUser, w *csv.Writer) error {
	data := []string{
		userData.Name,
		fmt.Sprint(userData.NumWorkHours),
		fmt.Sprintf("%.1f", userData.NumWorkDays),
		fmt.Sprint(userData.NumWeekendHours),
		fmt.Sprintf("%.1f", userData.NumWeekendDays),
		fmt.Sprint(userData.NumBankHolidaysHours),
		fmt.Sprintf("%.1f", userData.NumBankHolidaysDays),
	}
	if err := w.Write(data); err != nil {
		log.Printf("error writing record to csv: %v", err)
		return err
	}
	return nil
}
