package main

import (
	"app/utils"
	"flag"
	"log"
	"math/rand"
	"os"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func main() {
	ccFile := flag.String("cc", "", "")
	skFile := flag.String("sk", "", "")
	evkFile := flag.String("key_eval", "", "")
	inputFile := flag.String("input", "", "")

	flag.Parse()

	dataJSON, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("os.Open(%s): %s", "config.json", err.Error())
	}

	params := utils.Parameters{}
	if err := params.UnmarshalJSON(dataJSON); err != nil {
		log.Fatalf(err.Error())
	}

	sk := rlwe.NewKeyGenerator(params.Scheme).GenSecretKeyNew()

	ecd := hefloat.NewEncoder(params.Scheme)

	enc := rlwe.NewEncryptor(params.Scheme, sk)

	/* #nosec G404 */
	r := rand.New(rand.NewSource(0))
	values := make([]complex128, params.Scheme.MaxSlots())
	for i := range values {
		values[i] = complex(2*r.Float64()-1, 2*r.Float64()-1)
	}

	pt := hefloat.NewPlaintext(params.Scheme, params.Scheme.MaxLevel())


	if err = ecd.Encode(values, pt); err != nil {
		log.Fatalf(err.Error())
	}

	input, err := enc.EncryptNew(pt)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(sk, *skFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(input, *inputFile); err != nil {
		log.Fatalf(err.Error())
	}

	var evk utils.EvaluationKeySet
	if evk, err = utils.NewEvaluationKeySet(params, sk, dataJSON); err != nil{
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(evk, *evkFile); err != nil {
		log.Fatalf(err.Error())
	}
}
