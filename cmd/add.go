package cmd

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/jamesdorevski/offbeat/client"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add [key1] [keyN...]",
	Short:   "Use your Jira issues to automatically generate worklogs for a given date range. Worklogs are added between 9:00 and 17:00.",
	Example: "offbeat add ABC-123 ABC-456 -s 2021-01-01 -e 2021-01-07",
	Args:    cobra.MinimumNArgs(1),
	Run:     addRun,
}

var start string
var end string
var weekends bool

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&start, "start", "s", "", "Date to start adding worklogs from. Format: YYYY-MM-DD")
	addCmd.Flags().StringVarP(&end, "end", "e", "", "Date to finish adding worklogs. Format: YYYY-MM-DD")
	addCmd.Flags().BoolVarP(&weekends, "weekends", "w", false, "Include weekends.")

	addCmd.MarkFlagRequired("start")
	addCmd.MarkFlagRequired("end")
}

type Worklog struct {
	issueId          string
	end              time.Time
	timeSpentSeconds int
	isBogus          bool
}

func validDate(start string, end string) bool {
	_, err := time.Parse("2006-01-02", start)
	if err != nil {
		return false
	}

	_, err = time.Parse("2006-01-02", end)
	if err != nil {
		return false
	}

	if start > end {
		return false
	}
	return true
}

func getSobOrFirstLog(target time.Time) time.Time {
	sob := time.Date(target.Year(), target.Month(), target.Day(), 9, 0, 0, 0, time.UTC)
	if sob.Before(target) {
		return sob
	}
	return target
}

func addRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("No Jira keys given. Please provide at least one Jira key.")
		panic("No Jira keys given")
	}

	if !validDate(start, end) {
		fmt.Println("Date range inputted is invalid. Please check that the dates are in the correct format and that the start date is before the end date.")
		panic("Invalid date range")
	}

	req := &client.GetWorklogsRequest{
		Start:  start,
		End:    end,
	}

	resp, err := client.GetWorklogs(req)
	if err != nil {
		panic(err)
	}

	// Key = start time, value = end time
	timeMap := make(map[time.Time]*Worklog)
	for _, log := range resp.Results {
		startTime, err := log.TimeStarted()
		if err != nil {
			panic(err)
		}

		endTime, err := log.TimeFinished()
		if err != nil {
			panic(err)
		}

		l := &Worklog{
			end:              endTime,
			timeSpentSeconds: log.TimeSpentSeconds,
			isBogus:          false,
		}
		timeMap[startTime] = l
	}

	// Need to convert keys to ids because the Tempo API only accepts ids
	var ids []string
	for _, key := range args {
		id, err := client.GetIssueId(key)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Sort keys by earliest to latest so we can iterate through them in order
	sorted := make([]time.Time, 0, len(timeMap))
	for k := range timeMap {
		sorted = append(sorted, k)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Before(sorted[j])
	})

	curr := getSobOrFirstLog(sorted[0])

	for i, k := range sorted {
		if !weekends && (curr.Weekday() == time.Saturday || curr.Weekday() == time.Sunday) {
			curr = timeMap[k].end
			continue
		}

		if curr.Equal(k) {
			curr = timeMap[k].end
			continue
		}

		randIssue := ids[r.Intn(len(ids))]

		// check if k is the last entry in slice or for the day 
		// if it is, add one last timeblock from current to 17:00
		if i == len(sorted)-1 || sorted[i+1].Day() != curr.Day() {
			eob := time.Date(curr.Year(), curr.Month(), curr.Day(), 17, 0, 0, 0, time.UTC)
			if curr.Before(eob) {
				l := &Worklog{
					issueId:          randIssue,
					end:              eob,
					timeSpentSeconds: int(eob.Sub(curr).Seconds()),
					isBogus:          true,
				}
				timeMap[curr] = l

				// Move current to next day
				curr = getSobOrFirstLog(sorted[i+1])
				continue
			}
		}

		l := &Worklog{
			issueId:          randIssue,
			end:              k,
			timeSpentSeconds: int(k.Sub(curr).Seconds()),
			isBogus:          true,
		}
		timeMap[curr] = l

		curr = timeMap[k].end
	}

	for k, v := range timeMap {
		if v.isBogus {
			req := &client.CreateWorklogRequest{
				IssueId:          v.issueId,
				StartDate:        k.Format("2006-01-02"),
				StartTime:        k.Format("15:04:05"),
				TimeSpentSeconds: v.timeSpentSeconds,
			}

			err := client.CreateWorklog(req)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Done!")
}
