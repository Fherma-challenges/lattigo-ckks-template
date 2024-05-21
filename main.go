package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"

	"app/internal/solution"
	"app/utils"
)

func main() {
	now := time.Now()
	cc := flag.String("cc", "", "")
	evkFile := flag.String("key_eval", "", "")
	inputFile := flag.String("input", "", "")
	outputFile := flag.String("output", "", "")

	flag.Parse()

	params := utils.Parameters{}
	evk := utils.EvaluationKeySet{}
	in := rlwe.Ciphertext{}

	if err := utils.Deserialize(&params, *cc); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Deserialize(&evk, *evkFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Deserialize(&in, *inputFile); err != nil {
		log.Fatalf(err.Error())
	}

	out, err := solution.SolveTestcase(params, evk, &in)
	if err != nil {
		log.Fatalf("solution.SolveTestcase: %s", err.Error())
	}

	utils.Serialize(out, *outputFile)
	fmt.Printf("Done: %s\n", time.Since(now))
}
