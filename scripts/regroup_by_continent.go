package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
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

func RecToProf(rec []string) Profession {
	return Profession{rec[0], rec[1], rec[2]}
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
	for {
		rec := ReadCsvLine(r)
		if len(rec) == 0 {
			break
		}
		jobs = append(jobs, RecToJob(rec))
	}

	return jobs
}

func ExtractProfessions() []Profession {
	var profs []Profession

	r := ReadCsv("data/professions.csv")
	for {
		rec := ReadCsvLine(r)
		if len(rec) == 0 {
			break
		}

		profs = append(profs, RecToProf(rec))
	}

	return profs
}

func ExtractCountries() []interface{} {
	var countries []interface{}

	files, err := ioutil.ReadDir("data/countries")
	Check(err)

	for _, file := range files {
		byt, err := ioutil.ReadFile("data/countries/" + file.Name())
		Check(err)
		var dat map[string]interface{}

		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}
		countries = append(countries, dat)
	}

	return countries
}

func FindContinent(lat float64, lon float64, countries []interface{}) map[string]interface{} {
	var res map[string]interface{}
	dist := -1.0
	for _, country := range countries {
		country := country.(map[string]interface{})
		mn_lon := country["MinLongitude"].(float64)
		mx_lon := country["MaxLongitude"].(float64)
		mn_lat := country["MinLatitude"].(float64)
		mx_lat := country["MaxLatitude"].(float64)

		lon_in := mn_lon <= lon && mx_lon >= lon
		lat_in := mn_lat <= lat && mx_lat >= lat

		if lon_in && lat_in {
			clat := country["Latitude"].(float64)
			clon := country["Longitude"].(float64)
			cdist := math.Abs((clat - lat) + (clon - lon))
			if dist == -1.0 || dist > cdist {
				dist = cdist
				res = country
			}
		}
	}
	return res
}

func main() {
	jobs := ExtractJobs()
	professions := ExtractProfessions()
	countries := ExtractCountries()

	// Iterate through list and print its contents.
	for i, job := range jobs {
		fmt.Println(i, job)
	}

	// Iterate through list and print its contents.
	for i, prof := range professions {
		fmt.Println(i, prof)
	}

	// Iterate through list and print its contents.
	for i, country := range countries {
		fmt.Println(i, country)
	}

	for _, job := range jobs {
		lat := job.office_latitude
		lon := job.office_longitude
		fmt.Println(lat, lon)

		country := FindContinent(lat, lon, countries)
		fmt.Println(country["region"])
		break
	}

}
