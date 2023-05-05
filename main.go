package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var sum int

type Bigrams struct {
	count int
	prob  float64
}

// just read the file and convert it into bigram structure
func ReadFile(fileName string) map[string]Bigrams {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("couldn't read file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	birgrams := make(map[string]Bigrams)

	for scanner.Scan() {
		bigrammedWord := fmt.Sprintf("^%s$", scanner.Text())

		for i := 0; bigrammedWord[i] != '$'; i++ {
			sum++
			bigram := string(bigrammedWord[i]) + string(bigrammedWord[i+1])
			b := birgrams[bigram]
			b.count++
			birgrams[bigram] = Bigrams{b.count, 0}
		}
	}
	return birgrams
}

func main() {
	bigrams := ReadFile("name.txt")
	probableBigrams(bigrams)
	firstLetter := ""
	fmt.Println("Choose the first letter that the name will start with:")
	fmt.Scan(&firstLetter)
	fmt.Println("______________________________________________________")
	firstLetter = fmt.Sprintf("^%s", firstLetter)

	for i := 0; i < 6; i++ {
		name := nameFactory(firstLetter, bigrams)
		name = name[1:]
		name = name[:len(name)-1]

		for len(name) < 3 || len(name) > 12 {
			name = nameFactory(firstLetter, bigrams)
			name = name[1:]
			name = name[:len(name)-1]
		}

		fmt.Println(name)
	}
	fmt.Println("______________________________________________________")

	fmt.Println("Я записал вероятности всех биграм в отдельный bigram.txt")

	file, err := os.Create("bigram.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	text := "|Биграмма|Вероятность в %\n"
	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range bigrams {
		text = fmt.Sprintf("  %s    |", i)
		text2 := fmt.Sprintf("%f", bigrams[i].prob*100) + "%"

		_, err = file.WriteString(fmt.Sprintf("%s%s\n", text, text2))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func nameFactory(f string, big map[string]Bigrams) string {
	nextL := f[len(f)-1]

	if nextL == '$' {
		return f
	}

	var contAbleBigrams []string
	var probAbleBigrams []float64
	for i := range big {
		if i[0] == nextL {
			contAbleBigrams = append(contAbleBigrams, i)
			b := big[i]
			probAbleBigrams = append(probAbleBigrams, b.prob)
		}
	}

	var newP float64

	for _, j := range probAbleBigrams {
		newP += j
	}

	newP = float64(1) / newP

	probi2 := make([]float64, len(probAbleBigrams))

	for i := 0; i < len(probAbleBigrams); i++ {
		probi2[i] = newP * probAbleBigrams[i]
	}

	chancesProb := make([]float64, len(probAbleBigrams))
	chancesProb[0] = probi2[0]

	for i := 1; i < len(probi2); i++ {
		chancesProb[i] = chancesProb[i-1] + probi2[i]
	}
	randFloat := rand.Float64()

	var index int

	for i := 0; i < len(chancesProb); i++ {
		if randFloat < chancesProb[i] {
			index = i
			break
		}
	}

	f += string(contAbleBigrams[index][1])

	return nameFactory(f, big)
}

// add a probability for bigrams to appear
func probableBigrams(bigrams map[string]Bigrams) {
	for i := range bigrams {
		b := bigrams[i]
		b.prob = float64(b.count) / float64(sum)
		bigrams[i] = Bigrams{b.count, b.prob}

	}
}
