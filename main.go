package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = "v1.0.1"

type BodyfileEntry struct {
	Zero           string
	Name           string
	Inode          string
	Mode           string
	UID            string
	GID            string
	Size           string
	ATime          string
	MTime          string
	CTime          string
	CrTime         string
	ATimeReadable  string
	MTimeReadable  string
	CTimeReadable  string
	CrTimeReadable string
}

func main() {
	// Command line flags
	inputFile := flag.String("i", "", "Input bodyfile path (required)")
	outputFile := flag.String("o", "", "Output CSV file path (required)")
	epochOnly := flag.Bool("e", false, "Keep timestamps in epoch format only (default is human-readable)")
	showVersion := flag.Bool("v", false, "Show version and exit")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *showVersion {
		fmt.Printf("bdf2csv %s\n", version)
		return
	}

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Error: Both input and output file paths are required")
		showHelp()
		os.Exit(1)
	}

	err := convertBodyfileToCSV(*inputFile, *outputFile, *epochOnly)
	if err != nil {
		log.Fatalf("Error converting bodyfile: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", *inputFile, *outputFile)
}

func showHelp() {
	fmt.Println("Linux Bodyfile to CSV Converter")
	fmt.Println("Usage:")
	fmt.Println("  bdf2csv -i <bodyfile> -o <csvfile> [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -i string         Input bodyfile path (required)")
	fmt.Println("  -o string         Output CSV file path (required)")
	fmt.Println("  -e                Keep timestamps in epoch format only (default is human-readable)")
	fmt.Println("  -v                Show version and exit")
	fmt.Println("  -h                Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  bdf2csv -i bodyfile.txt -o bodyfile.csv")
	fmt.Println("  bdf2csv -i bodyfile.txt -o bodyfile.csv -e")
}

func convertBodyfileToCSV(inputPath, outputPath string, epochOnly bool) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create CSV writer
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Write CSV header
	header := []string{
		"0", "Name", "Inode", "Mode", "UID", "GID", "Size",
		"ATime", "MTime", "CTime", "CrTime",
	}

	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Read and process bodyfile line by line
	scanner := bufio.NewScanner(inputFile)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		entry, err := parseBodyfileLine(line)
		if err != nil {
			log.Printf("Warning: Failed to parse line %d: %v", lineNumber, err)
			continue
		}

		// Convert to CSV record
		record := []string{
			entry.Zero, entry.Name, entry.Inode, entry.Mode,
			entry.UID, entry.GID, entry.Size,
		}

		if epochOnly {
			// Use epoch timestamps
			record = append(record, entry.ATime, entry.MTime, entry.CTime, entry.CrTime)
		} else {
			// Default: use human-readable timestamps only
			record = append(record, entry.ATimeReadable, entry.MTimeReadable,
				entry.CTimeReadable, entry.CrTimeReadable)
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record at line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	return nil
}

func parseBodyfileLine(line string) (*BodyfileEntry, error) {
	// Split by pipe character
	fields := strings.Split(line, "|")

	// Bodyfile format can have 10 or 11 fields (CrTime may not be available on all filesystems)
	if len(fields) < 10 || len(fields) > 11 {
		return nil, fmt.Errorf("invalid bodyfile format: expected 10 or 11 fields, got %d", len(fields))
	}

	entry := &BodyfileEntry{
		Zero:  fields[0],
		Name:  fields[1],
		Inode: fields[2],
		Mode:  fields[3],
		UID:   fields[4],
		GID:   fields[5],
		Size:  fields[6],
		ATime: fields[7],
		MTime: fields[8],
		CTime: fields[9],
	}

	// Handle CrTime - may not be present on all filesystems
	if len(fields) == 11 {
		entry.CrTime = fields[10]
	} else {
		entry.CrTime = "0" // Default to 0 if not available
	}

	// Convert timestamps to human-readable format
	entry.ATimeReadable = convertTimestamp(fields[7])
	entry.MTimeReadable = convertTimestamp(fields[8])
	entry.CTimeReadable = convertTimestamp(fields[9])
	entry.CrTimeReadable = convertTimestamp(entry.CrTime)

	return entry, nil
}

func convertTimestamp(epochStr string) string {
	if epochStr == "" || epochStr == "0" {
		return "N/A"
	}

	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return "Invalid"
	}

	// Convert to time and format as RFC3339 (ISO 8601)
	t := time.Unix(epoch, 0).UTC()
	return t.Format("2006-01-02 15:04:05 UTC")
}
