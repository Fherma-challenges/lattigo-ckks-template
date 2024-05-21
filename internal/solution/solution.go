package solution

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"app/utils"
)

func SolveTestcase(
	params utils.Parameters,
	evk utils.EvaluationKeySet,
	in *rlwe.Ciphertext,
) (out *rlwe.Ciphertext, err error) {

	paramsScheme := params.Scheme
	paramsBootstrapping := params.Bootstrapping

	// Instantiate the evaluator with the evaluation keys
	eval := hefloat.NewEvaluator(paramsScheme, evk.Scheme)

	if err = eval.Conjugate(in, in); err != nil{
		return
	}

	// Instantiate the bootstrapping Evaluator (if enabled)
	var btp *bootstrapping.Evaluator
	if paramsBootstrapping != nil{
		if btp, err = bootstrapping.NewEvaluator(*paramsBootstrapping, evk.Bootstrapping); err != nil{
			return
		}
	}

	if btp != nil{
		// bootstrapping.Evaluator is compliant to the interface he.Bootstrapper[rlwe.Ciphertext] (/he/bootstrapper.go)
		// see /he/hefloat/bootstrapping/bootstrapping for individual methods of the bootstrapping evaluator
		// see examples/single_party/applications/reals_bootstrapping for bootstrapping examples
		if out, err = btp.Bootstrap(in); err != nil{
			return
		}
	}else{
		out = in
	}

	// Put your solution here
	return out, nil
}
