bin/samplesort: cmd/samplesort/main.go essentia.go fs.go engine/* collection/* sample/* crypto/* analyze/* extractor/* parser/*
	go build -o bin/samplesort samplesort/cmd/samplesort
