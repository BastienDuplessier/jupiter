package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"io/ioutil"
	"math"
	"os"
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

type Key struct {
	profession string
	continent  string
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

func ExtractProfessions() map[string]Profession {
	var profs map[string]Profession
	profs = make(map[string]Profession)

	r := ReadCsv("data/professions.csv")
	for {
		rec := ReadCsvLine(r)
		if len(rec) == 0 {
			break
		}

		prof := RecToProf(rec)
		profs[prof.id] = prof
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

func FindContinent(lat float64, lon float64, countries []interface{}) string {
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
	if res != nil {
		return res["region"].(string)
	} else {
		return "unknown"
	}
}

func BuildData(result map[Key]int, cont_res map[string]int, prof_res map[string]int) ([][]string, []string) {
	var data [][]string
	var header []string
	var line []string
	var prof_order []string
	var cont_order []string
	var total int

	_ = result

	// Define continent order
	for cont, _ := range cont_res {
		cont_order = append(cont_order, cont)
	}

	// Total of all offers
	total = 0
	for _, val := range cont_res {
		total += val
	}

	// Header, professions names
	header = append(header, "")
	header = append(header, "TOTAL")

	for prof, _ := range prof_res {
		header = append(header, prof)
		prof_order = append(prof_order, prof)
	}

	// Second line, professions totals
	line = []string{}
	line = append(line, "TOTAL")
	line = append(line, strconv.Itoa(total))

	for _, prof := range prof_order {
		line = append(line, strconv.Itoa(prof_res[prof]))
	}
	data = append(data, line)

	// Now, the rest of the lines, for for each continent
	for _, cont := range cont_order {
		line = []string{}
		line = append(line, strings.ToUpper(cont))
		line = append(line, strconv.Itoa(cont_res[cont]))
		for _, prof := range prof_order {
			key := Key{prof, cont}
			line = append(line, strconv.Itoa(result[key]))
		}
		data = append(data, line)
	}

	return data, header
}

func ComputeTotByContinent(result map[Key]int) map[string]int {
	var res map[string]int
	res = make(map[string]int)

	for key, nb := range result {
		res[key.continent] += nb
	}
	return res
}

func ComputeTotByProf(result map[Key]int) map[string]int {
	var res map[string]int
	res = make(map[string]int)

	for key, nb := range result {
		res[key.profession] += nb
	}
	return res
}

func PrintResult(result map[Key]int) {
	total_by_continent := ComputeTotByContinent(result)
	fmt.Println(total_by_continent)
	total_by_prof := ComputeTotByProf(result)
	fmt.Println(total_by_prof)

	data, header := BuildData(result, total_by_continent, total_by_prof)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

func main() {
	jobs := ExtractJobs()
	professions := ExtractProfessions()
	countries := ExtractCountries()

	var result map[Key]int
	result = make(map[Key]int)

	for _, job := range jobs {
		lat := job.office_latitude
		lon := job.office_longitude
		continent := FindContinent(lat, lon, countries)
		profession := professions[job.profession_id].name
		if profession == "" {
			profession = "Unknown"
		}
		key := Key{continent: continent, profession: profession}
		result[key] += 1
	}

	PrintResult(result)
}
