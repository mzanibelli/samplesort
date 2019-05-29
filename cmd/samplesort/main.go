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

	col, err := samplesort.SampleSort(
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

	feats := col.Features()

	for index, sample := range feats {
		row := make(plotter.XYs, len(sample))
		for i := range row {
			row[i].X = float64(i)
			row[i].Y = sample[i]
		}
		err = plotutil.AddLinePoints(p, string(index), row)
		if err != nil {
			panic(err)
		}
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "/tmp/a.png"); err != nil {
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
