package utils

import(
	"io"
	"bufio"
	"encoding/json"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"github.com/tuneinsight/lattigo/v5/utils"
	"github.com/tuneinsight/lattigo/v5/utils/buffer"
)

type EvaluationKeySet struct{
	Scheme *rlwe.MemEvaluationKeySet
	Bootstrapping *bootstrapping.EvaluationKeys
}

func NewEvaluationKeySet(params Parameters, sk *rlwe.SecretKey, JSON []byte) (evk EvaluationKeySet, err error){
	aux := struct{
		EvaluationKeys struct{
			Rotations []int
			GaloisElements []uint64
			Relinearization bool
		}
	}{}

	if err = json.Unmarshal(JSON, &aux); err != nil{
		return
	}

	kgen := rlwe.NewKeyGenerator(params.Scheme)

	var rlk *rlwe.RelinearizationKey
	if aux.EvaluationKeys.Relinearization {
		rlk = kgen.GenRelinearizationKeyNew(sk)
	}

	var gks []*rlwe.GaloisKey


	keys := map[uint64]bool{}
	for _, k := range aux.EvaluationKeys.Rotations{
		keys[params.Scheme.GaloisElement(k)] = true
	}

	for _, galEl := range aux.EvaluationKeys.GaloisElements{
		keys[galEl] = true
	}

	if  galEls := utils.GetSortedKeys(keys); len(galEls) != 0{
		gks = make([]*rlwe.GaloisKey, len(aux.EvaluationKeys.GaloisElements))
		for i, galEl := range aux.EvaluationKeys.GaloisElements{
			gks[i] = kgen.GenGaloisKeyNew(galEl, sk)
		}
	}

	evk.Scheme = rlwe.NewMemEvaluationKeySet(rlk, gks...)

	if params.Bootstrapping != nil{
		if evk.Bootstrapping, _, err = params.Bootstrapping.GenEvaluationKeys(sk); err != nil{
			return
		}
	}

	return
}

func (evk EvaluationKeySet) BinarySize() int {
	return 1 + evk.Scheme.BinarySize() + 8 + evk.Bootstrapping.BinarySize()
}

func (evk EvaluationKeySet) WriteTo(w io.Writer) (n int64, err error) {
	switch w := w.(type) {
	case buffer.Writer:

		var inc int64

		if evk.Scheme != nil{
			if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
				return
			}


			n += inc

			if inc, err = evk.Scheme.WriteTo(w); err != nil{
				return
			}

			n += inc

		}else{
			if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
				return
			}

			n += inc
		}

		if evk.Bootstrapping != nil{

			if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
				return
			}

			n += inc

			if evk.Bootstrapping.EvkN1ToN2 != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkN1ToN2.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.EvkN2ToN1 != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkN2ToN1.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.EvkRealToCmplx != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkRealToCmplx.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.EvkCmplxToReal != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkCmplxToReal.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.EvkDenseToSparse != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkDenseToSparse.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.EvkSparseToDense != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.EvkSparseToDense.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}

			if evk.Bootstrapping.MemEvaluationKeySet != nil {

				if inc, err = buffer.WriteAsUint8[int](w, 1); err != nil {
					return
				}

				n += inc

				if inc, err = evk.Bootstrapping.MemEvaluationKeySet.WriteTo(w); err != nil {
					return
				}

				n += inc

			} else {

				if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
					return
				}

				n += inc
			}
		}else{
			if inc, err = buffer.WriteAsUint8[int](w, 0); err != nil {
				return
			}

			n += inc
		}

		return n, w.Flush()
	default:
		return evk.WriteTo(bufio.NewWriter(w))
	}
}

func (evk *EvaluationKeySet) ReadFrom(r io.Reader) (n int64, err error) {
	switch r := r.(type) {
	case buffer.Reader:

		var inc int64

		var exist int

		if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
			return
		}

		n += inc

		if exist == 1{
			if evk.Scheme == nil{
				evk.Scheme = &rlwe.MemEvaluationKeySet{}
			}

			if inc, err = evk.Scheme.ReadFrom(r); err != nil {
				return
			}

			n += inc
		}

		if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
			return
		}

		if exist == 1{

			if evk.Bootstrapping == nil{
				evk.Bootstrapping = &bootstrapping.EvaluationKeys{}
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkN1ToN2 == nil{
					evk.Bootstrapping.EvkN1ToN2 = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkN1ToN2.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkN2ToN1 == nil{
					evk.Bootstrapping.EvkN2ToN1 = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkN2ToN1.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkRealToCmplx == nil{
					evk.Bootstrapping.EvkRealToCmplx = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkRealToCmplx.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkCmplxToReal == nil{
					evk.Bootstrapping.EvkCmplxToReal = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkCmplxToReal.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkDenseToSparse == nil{
					evk.Bootstrapping.EvkDenseToSparse = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkDenseToSparse.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.EvkSparseToDense == nil{
					evk.Bootstrapping.EvkSparseToDense = &rlwe.EvaluationKey{}
				}

				if inc, err = evk.Bootstrapping.EvkSparseToDense.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}

			if inc, err = buffer.ReadAsUint8[int](r, &exist); err != nil {
				return
			}

			n += inc

			if exist == 1 {

				if evk.Bootstrapping.MemEvaluationKeySet == nil{
					evk.Bootstrapping.MemEvaluationKeySet = &rlwe.MemEvaluationKeySet{}
				}

				if inc, err = evk.Bootstrapping.MemEvaluationKeySet.ReadFrom(r); err != nil {
					return
				}

				n += inc
			}
		}

		return n, nil
	default:
		return evk.ReadFrom(bufio.NewReader(r))
	}
}