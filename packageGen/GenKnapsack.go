package packageGen

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	// "bufio"
    "encoding/csv"
	// "log"
	// "io"
	"strconv"
	"gonum.org/v1/gonum/stat"
)
type Data struct {
	Seed int64 `json:"seed"`
	Score float64 `json:"score"`
	Payload [] Entry `json:"payload"`
}

type Entry struct {
	TextID int `json:"textid"`
	WordCount  int `json:"# words"`
	Genre string `json:"genre"`
	Year int `json:"year"`
}

type SelectionObj struct {
	selectionArray[] int
	fitnessScore float64
}

type SelectionStatObj struct {
	genresWC map[string] int
	years map[int] int
	yearsWC map[int] int
	genres map[string] int
}

func Gen(seed int64) (* Data){
	fmt.Println(seed)
	rand.Seed(seed)
	domain, domain_len := select_domain(1830, 1839)
	generation := 3
	inds_count := 1024
	inds := make([]*SelectionObj, 0)

	for i := 0; i < generation; i++ {
		if i  % 10 == 0 {
			fmt.Println("Generation: ", i, " ", len(inds), " individuals")
		}
		if i == 0 {
			fmt.Println("New Spec")
			for j := 0; j < inds_count; j++{
				inds = append(inds, new_random_selection(domain_len, seed))
			}
		}
		for j := range inds {
			calc_fitness(domain, inds[j], false)
		}
		inds = sort_inds(inds)
		if i % 10 == 0 {
			fmt.Println(inds[0].fitnessScore, inds[63].fitnessScore)
		}
		inds = breed(inds, domain, seed)
		
	}
	for j := range inds {
		calc_fitness(domain, inds[j], false)
	}
	inds = sort_inds(inds)
	calc_fitness(domain, inds[0], true)
	selected := print_values(domain, inds)
	//main_test()
	data := new(Data)
	data.Seed = seed
	data.Score = inds[0].fitnessScore
	data.Payload = selected
	return data
}

func print_values(domain []Entry, inds []*SelectionObj)([] Entry){
	var selected [] Entry
	for i := 0; i < len(inds[0].selectionArray); i++ {
		if inds[0].selectionArray[i] == 1 {
			selected = append(selected, domain[i])
		}
	}
	fmt.Println(len(selected))
	return selected
}

func breed(inds []*SelectionObj, domain []Entry, seed int64)([]*SelectionObj) {
	halfway := len(inds[0].selectionArray) / 2
	offsprings := make([]*SelectionObj, 0)
	// rand.Seed(seed)
	for i := 0; i < (len(inds) / 2) - 0; i = i + 2 {
		male := inds[i].selectionArray
		female := inds[rand.Intn(len(inds))].selectionArray
		
		offspring1 := new(SelectionObj)
		offspring2 := new(SelectionObj)

		offspring1.selectionArray = append(male[halfway:], female[:halfway]...)
		offspring2.selectionArray = append(female[halfway:], male[:halfway]...)

		offsprings = append(offsprings, offspring1, offspring2, inds[i], mass_mutate(inds[i], domain, seed))

	}
	return offsprings
}

func mass_mutate( obj *SelectionObj, domain []Entry, seed int64)(*SelectionObj) {
	// rand.Seed(seed)
	mutationChance := .05
	remutate_max := 3
	stats := calc_fitness(domain, obj, false)
	for key, value := range stats.genresWC {
		words_needed := 2000000 - value
		for i := 0; i < len(obj.selectionArray); i++ {
			if words_needed > 0 && obj.selectionArray[i] == 0 && domain[i].Genre == key && rand.Float64() < mutationChance{
				obj.selectionArray[i] = 1
				words_needed = words_needed - domain[i].WordCount
			} else if words_needed < 0 && obj.selectionArray[i] == 1 && domain[i].Genre == key && rand.Float64() < mutationChance {
				obj.selectionArray[i] = 0
				words_needed = words_needed + domain[i].WordCount
			}
			if words_needed < 1000 && words_needed > 0 {
				break
			}
			if remutate_max > 0 && i == len(obj.selectionArray) - 1 && (words_needed > 1000 || words_needed < 0) {
				remutate_max--
				i = 0
			}
		}
	}
	return obj
}

func sort_inds(inds []*SelectionObj)([]*SelectionObj) {
	for i := 0; i < len(inds) -1; i++ {
		for j := 0; j < len(inds) - i - 1; j++ {
			if inds[i].fitnessScore < inds[j].fitnessScore {
				temp := new(SelectionObj)
				temp = inds[i]
				inds[i] = inds[j]
				inds[j] = temp
			}
		}
	}
	return inds
}

func str_to_int(s string)(int){
	i,_:= strconv.Atoi(s)
	return i
}
func select_domain(year_lower_bound int, year_upper_bound int) ([]Entry, int) {
	fmt.Println("opening file")
	// file, _ := os.Open("test.csv")
	file, _ := os.Open("./packageGen/sources_coha_for_algo.csv")
	var entires [] Entry
	var domain [] Entry
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
    if err != nil {
        panic(err)
    }
	for _, line := range lines {
		entires = append(entires, Entry {
			TextID: str_to_int(line[0]),
			WordCount: str_to_int(line[1]),
			Genre: line[2],
			Year: str_to_int(line[3]),
		})
	}

	for k := range entires {
		if in_bounds(entires[k].Year, year_lower_bound, year_upper_bound) && entires[k].Genre != "NEWS"{
			domain = append(domain, entires[k])
		}
	}
	fmt.Println(domain[0])
	return domain, len(domain)
}

func in_bounds(value int, lower_bound int, upper_bound int) (bool){
	if value >= lower_bound && value <= upper_bound {
		return true
	} else {
		return false
	}
}

func new_random_selection(len int, seed int64)(*SelectionObj) {
	// rand.Seed(seed)
	obj := new(SelectionObj)
	obj.selectionArray = make([]int, len)
	for i := range obj.selectionArray {
		if rand.Float64() < .20{
			obj.selectionArray[i] = 1
		}	
	}
	return obj
}

func print_selectionObj(obj *SelectionObj){
	fmt.Println(obj.selectionArray)
	fmt.Println(obj.fitnessScore)
}

func calc_fitness(domain []Entry, obj *SelectionObj, print bool)(*SelectionStatObj){
	stat := new(SelectionStatObj)
	stat.years = make(map[int] int)
	stat.yearsWC = make(map[int] int)
	stat.genresWC = make(map[string] int)
	stat.genres = make(map[string] int)
	for i := range obj.selectionArray {
		if obj.selectionArray[i] == 1 {
			stat.years[domain[i].Year]++
			stat.yearsWC[domain[i].Year] = stat.yearsWC[domain[i].Year] + domain[i].WordCount
			stat.genresWC[domain[i].Genre] = stat.genresWC[domain[i].Genre] + domain[i].WordCount
			stat.genres[domain[i].Genre]++
		}
	}
	if print {
		fmt.Println(stat.years)
		fmt.Println(stat.genresWC)
		fmt.Println(stat.genres)
		fmt.Println(stat.yearsWC)
	}
	
	obj.fitnessScore = fitnessScore(stat)
	return stat
}
func fitnessScore(stats *SelectionStatObj)(float64){

	yearsVector := make([]float64, 0, len(stats.years))
	for _, value := range stats.years {
		yearsVector= append(yearsVector, float64(value))
	}
	yearsScore := stat.StdDev(yearsVector,nil)

	yearsWCVector := make([]float64, 0, len(stats.yearsWC))
	for _, value := range stats.yearsWC {
		yearsWCVector = append(yearsWCVector, float64(value))
	}
	yearsWCScore := stat.StdDev(yearsWCVector,nil)

	genresVector := make([]float64, 0, len(stats.genres))
	for _, value := range stats.genres {
		genresVector= append(genresVector, float64(value))
	}
	genresScore :=stat.StdDev(genresVector,nil)

	genresWCVector := make([]float64, 0, len(stats.genresWC))
	var genresWCSum int
	for _, value := range stats.genresWC {
		genresWCVector= append(genresWCVector, float64(value))
		genresWCSum = genresWCSum + value
	}
	genresWCScore := stat.StdDev(genresWCVector,nil) + math.Abs(float64(6000000 - genresWCSum))
	
	return 1 + yearsScore + genresScore + genresWCScore + yearsWCScore
}

func main_test() {
	test1 := new(SelectionStatObj)
	test2 := new(SelectionStatObj)
	test1.years = map[int] int {
		1820: 10,
		1821: 10,
		1822: 10,
		1823: 10,
		1824: 10,
		1825: 10,
		1826: 10,
		1827: 10,
		1828: 10,
		1829: 10,
	}
	test2.years = map[int] int {
		1820: 50,
		1821: 0,
		1822: 0,
		1823: 0,
		1824: 0,
		1825: 0,
		1826: 0,
		1827: 0,
		1828: 0,
		1829: 0,
	}
	test1.genres = map[string] int {
		"NF" : 33,
		"FIC" : 33,
		"MAG" : 33,
	}
	test2.genres = map[string] int {
		"NF" : 100,
		"FIC" : 0,
		"MAG" : 0,
	}
	test1.genresWC = map[string] int {
		"NF" : 2000001,
		"FIC" : 2000001,
		"MAG" : 2000001,
	}
	test2.genresWC = map[string] int {
		"NF" : 6000000,
		"FIC" : 0,
		"MAG" : 0,
	}
	select_1830 := new(SelectionStatObj)
	select_1830.years = map[int] int {
		1820: 26,
		1821: 47,
		1822: 63,
		1823: 37,
		1824: 32,
		1825: 41,
		1826: 18,
		1827: 24,
		1828: 29,
		1829: 35,
	}
	select_1830.genres = map[string] int {
		"NF" : 43,
		"FIC" : 49,
		"MAG" : 260,
	}
	select_1830.genresWC = map[string] int {
		"NF" : 2002602,
		"FIC" : 2000906,
		"MAG" : 2000098,
	}
	fmt.Println(fitnessScore(test1))
	fmt.Println(fitnessScore(test2))
	fmt.Println(fitnessScore(select_1830))
}