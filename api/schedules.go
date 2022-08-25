package api

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type Schedule struct {
	ID            string
	Name          string
	TimeZone      string
	FinalSchedule ScheduleLayer
}

type ScheduleLayer struct {
	RenderedScheduleEntries []RenderedScheduleEntry
}

type RenderedScheduleEntry struct {
	Start string
	End   string
	User  User
}

type ScheduleInfo struct {
	ID            string
	Name          string
	Location      *time.Location
	Start         time.Time
	End           time.Time
	FinalSchedule ScheduleLayer
}

func (p *PagerDutyClient) ListSchedules() ([]*Schedule, error) {
	var opts pagerduty.ListSchedulesOptions
	var scheduleList []*Schedule

	more := true
	for more {
		listSchedulesResponse, err := p.ApiClient.ListSchedules(opts)
		if err != nil {
			return nil, err
		}
		for _, schedule := range listSchedulesResponse.Schedules {
			scheduleList = append(scheduleList, convertSchedule(&schedule))
		}
		more = listSchedulesResponse.More
		opts.Offset = listSchedulesResponse.Limit
	}

	return scheduleList, nil
}

func (p *PagerDutyClient) GetSchedule(scheduleID, startDate, endDate string) (*Schedule, error) {
	var opts pagerduty.GetScheduleOptions
	opts.Since = startDate
	opts.Until = endDate
	scheduleResponse, err := p.ApiClient.GetSchedule(scheduleID, opts)
	if err != nil {
		return nil, err
	}

	return convertSchedule(scheduleResponse), nil
}

func convertSchedule(schedule *pagerduty.Schedule) *Schedule {
	return &Schedule{
		ID:            schedule.ID,
		Name:          schedule.Name,
		TimeZone:      schedule.TimeZone,
		FinalSchedule: convertScheduleLayer(schedule.FinalSchedule),
	}
}

func convertScheduleLayer(layer pagerduty.ScheduleLayer) ScheduleLayer {
	return ScheduleLayer{
		RenderedScheduleEntries: convertRenderedScheduleEntry(layer.RenderedScheduleEntries),
	}
}

func convertRenderedScheduleEntry(entries []pagerduty.RenderedScheduleEntry) []RenderedScheduleEntry {
	var entryList []RenderedScheduleEntry
	for _, entry := range entries {
		entryList = append(entryList, RenderedScheduleEntry{
			Start: entry.Start,
			End:   entry.End,
			User: User{
				ID:      entry.User.ID,
				Summary: entry.User.Summary,
			},
		})
	}

	return entryList
}
