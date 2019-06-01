package main

import (
	"fmt"
	"log"
	"os"
	"samplesort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	env string = "ESSENTIA_EXTRACTOR"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	feats, err := samplesort.SampleSort(
		os.Args[1],
		os.Getenv(env),
		logger,
	)

	usage(err)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "SampleSort"
	p.X.Label.Text = "Features"
	p.Y.Label.Text = "Value"

	var items int = 20

	logger.Println("start plotting...")
	lines := make([]interface{}, items)
	for index, sample := range feats[:items] {
		row := make(plotter.XYs, len(sample))
		for i := range row {
			row[i].X = float64(i)
			row[i].Y = sample[i]
		}
		lines[index] = row
	}
	err = plotutil.AddLinePoints(p, lines...)
	if err != nil {
		panic(err)
	}

	logger.Println("saving to file...")
	if err := p.Save(20*vg.Inch, 20*vg.Inch, "/tmp/a.png"); err != nil {
		panic(err)
	}
}

func usage(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", env, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", samplesort.Version, samplesort.Checksum)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}
