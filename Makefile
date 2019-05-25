bin/samplesort: cmd/samplesort/main.go essentia.go engine/* collection/* sample/*
	go build -o bin/samplesort samplesort/cmd/samplesort
