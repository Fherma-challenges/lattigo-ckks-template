package main

import (
	"app/utils"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"math/cmplx"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func main() {
	ccFile := flag.String("cc", "", "")
	skFile := flag.String("sk", "", "")
	outputFile := flag.String("output", "", "")

	flag.Parse()

	params := utils.Parameters{}
	if err := utils.Deserialize(&params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	sk := rlwe.SecretKey{}
	if err := utils.Deserialize(&sk, *skFile); err != nil {
		log.Fatalf(err.Error())
	}

	out := rlwe.Ciphertext{}
	if err := utils.Deserialize(&out, *outputFile); err != nil {
		log.Fatalf(err.Error())
	}

	dec := rlwe.NewDecryptor(params.Scheme, &sk)
	ecd := hefloat.NewEncoder(params.Scheme)

	have := make([]complex128, out.Slots())
	if err := ecd.Decode(dec.DecryptNew(&out), have); err != nil {
		log.Fatalf("%T.Decode: %s", ecd, err.Error())
	}

	r := rand.New(rand.NewSource(0))
	want := make([]complex128, params.Scheme.MaxSlots())
	for i := range want {
		want[i] = complex(2*r.Float64()-1, 2*r.Float64()-1)
	}

	for i := range want{
		want[i] = cmplx.Conj(want[i])
	}

	fmt.Println(have[:4])
	fmt.Println(want[:4])

	fmt.Println(hefloat.GetPrecisionStats(params.Scheme, ecd, nil, have, want, 0, false).String())
}
