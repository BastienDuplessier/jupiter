package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Job struct {
	profession_id    string
	contract_type    string
	name             string
	office_latitude  float64
	office_longitude float64
}

type Profession struct {
	id            string
	name          string
	category_name string
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseFloat64(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return res
	} else {
		return -1
	}
}

func ReadCsv(filename string) *csv.Reader {
	dat, err := ioutil.ReadFile(filename)
	Check(err)

	r := csv.NewReader(strings.NewReader(string(dat)))
	r.Read() // Discard headers
	return r
}

func RecToJob(rec []string) Job {
	lat := ParseFloat64(rec[3])
	lon := ParseFloat64(rec[4])
	return Job{rec[0], rec[1], rec[2], lat, lon}
}

func ReadCsvLine(reader *csv.Reader) []string {
	rec, err := reader.Read()
	if err == io.EOF {
		var empty []string
		return empty
	}
	Check(err)
	return rec
}

func ExtractJobs() []Job {
	var jobs []Job

	r := ReadCsv("data/jobs.csv")
	r.Read()
	for {
		rec := ReadCsvLine(r)
		if len(rec) == 0 {
			break
		}
		jobs = append(jobs, RecToJob(rec))
	}

	return jobs
}

func main() {
	jobs := ExtractJobs()

	// Iterate through list and print its contents.
	for i, job := range jobs {
		fmt.Println("%q - %q", i, job)
	}
}
