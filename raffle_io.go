package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/zyedidia/generic/list"
)

func formatArray(arr []int) string {
	return fmt.Sprintf("(%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]"))
}

func printToFile(rolledEntries map[string]*RolledEntry) {
	file, err := os.Create("raffle_results.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for patron, rolledEntry := range rolledEntries {
		writer.WriteString(fmt.Sprintf("+%d for %s %s \n", rolledEntry.points, patron, formatArray(rolledEntry.rolls)))
	}
}

func importDataFile() *list.List[RaffleEntry] {
	raffleEntries := list.New[RaffleEntry]()

	readFile, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineSplit := strings.Split(line, ": ")
		points, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			log.Fatal(err)
		}
		raffleEntries.PushBack(RaffleEntry{Patron: lineSplit[0], Points: points})
	}

	// writeData(raffleEntries, true)

	return raffleEntries
}

func exportDataFile(raffleData *list.List[RaffleEntry]) {
	file, err := os.Create("data_new.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for entry := raffleData.Front; entry != nil; entry = entry.Next {
		writer.WriteString(fmt.Sprintf("%s: %d\n", entry.Value.Patron, entry.Value.Points))
	}
}

func readData() *list.List[RaffleEntry] {
	raffleData := list.New[RaffleEntry]()

	content, err := fs.ReadFile(os.DirFS("."), "points.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	_, err = r.Read()
	if err == io.EOF {
		return raffleData
	} else if err != nil {
		log.Fatal(err)
	}

	for {
		raffleEntry, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		points, err := strconv.Atoi(raffleEntry[1])
		if err != nil {
			log.Fatal(err)
		}

		raffleData.PushBack(RaffleEntry{Patron: raffleEntry[0], Points: points})
	}
	return raffleData
}

func writeData(raffleData *list.List[RaffleEntry], override bool) {
	var fileName string
	if override {
		fileName = "points.csv"
	} else {
		fileName = "points_new.csv"
	}
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"patron", "points"})

	for entry := raffleData.Front; entry != nil; entry = entry.Next {
		writer.Write([]string{entry.Value.Patron, strconv.Itoa(entry.Value.Points)})
	}
}
