package batch

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/zerjioang/rooftop/v2/analysis"
)

// CSVRow represents a row in the CSV file.
type CSVRow struct {
	BlockTimestamp string
	Address        string
	Bytecode       string
}

func ProcessCSV(path string, limit uint, minsize int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReaderSize(file, 64*1024)) // 64 KB buffer
	reader.ReuseRecord = false                                  // Reuse memory for each record
	reader.FieldsPerRecord = 3

	var wg sync.WaitGroup
	workerCount := 1
	rows := make(chan []string, workerCount*4) // Buffered to reduce contention

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		go func() {
			log.Println("creating csv processor goroutine")
			for record := range rows {
				// Avoid allocations by using direct assignment
				row := CSVRow{
					BlockTimestamp: record[0],
					Address:        record[1],
					Bytecode:       record[2],
				}
				processRow(&row)
				wg.Done()
			}
		}()
	}

	// Read and dispatch rows
	// Skip header
	log.Println("reading header")
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	processed := uint(0)
	log.Println("reading rows")
	for {
		fmt.Println(processed + 1)
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading record: %w", err)
		}
		if len(record[2]) < minsize {
			continue
		}

		wg.Add(1)
		rows <- record
		if limit != 0 && processed >= limit {
			break
		}
		processed++
	}

	// Wait and cleanup
	wg.Wait()
	close(rows)
	return nil
}

// processRow is your row-processing function (runs in goroutines).
func processRow(row *CSVRow) {
	// Replace this with your actual processing logic.
	cli := analysis.NewCLI()
	log.Printf("scanning contract [%s] = %s\n", row.Address, row.Bytecode)
	name := fmt.Sprintf("%s_%s",
		padUint32(len(row.Bytecode)),
		row.Address,
	)
	if err := cli.Run(context.Background(), name, row.Bytecode, time.Now()); err != nil {
		log.Fatal(err)
	}
}

func padUint32(n int) string {
	return fmt.Sprintf("%010d", n)
}
