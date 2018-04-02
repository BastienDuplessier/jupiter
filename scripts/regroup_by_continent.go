package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "encoding/csv"
	"strings"
  "strconv"
)

type Job struct {
  profession_id string
  contract_type string
  name string
  office_latitude float64
  office_longitude float64
}


type Profession struct {
  id string
  name string
  category_name string
}

func Check(e error) {
  if e != nil {
    panic(e)
  }
}

func ParseFloat64(str string)float64 {
  res, err := strconv.ParseFloat(str, 64)
  if err == nil {
    return res
  } else {
    return -1
  }
}

func ExtractJobs()[]Job {  
  var jobs []Job

  dat, err := ioutil.ReadFile("data/jobs.csv")
  Check(err)

	r := csv.NewReader(strings.NewReader(string(dat)))
	r.Read() // Discard headers

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
    Check(err)

		fmt.Println(rec)
    
    lat := ParseFloat64(rec[3])
    lon := ParseFloat64(rec[4])
    job := Job{rec[0], rec[1], rec[2], lat, lon}
    jobs = append(jobs, job)
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
