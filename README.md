# PARITY CHALLENGE

## Setting the Parameters

The `config.json` file provides a `JSON` definition of the scheme and bootstrapping parameters.
We provide a list of the field that can be specified

### Scheme Parameters

#### Mandatory
- LogN: base two logarithm of the ring degree, example: `16`.
- LogQ: base two logarithm list of prime factors the ciphertext modulus, example: `[60, 50, 50, 50]`.
- LogP: base two logarithm list of prime factors the auxiliary prime modulus, example: `[61, 61]`.
- LogDefaultScale: base two logarithm of the default scaling factor: example `50`.

#### Optional
- Xe: distribution of the error, for example `{"Type": "DiscreteGaussian", "Sigma": 6.6, "Bound": 39.6}`.
- Xs: distribution of the secret, for example `{"Type": "Ternary", "H": 192}` (Gaussian secret are also supported).
- RingType: 
	- `Standard` for `R[X]/(X^{N}+1)`: provides N/2 complex slots.
	- `ConjugateInvariant` for `R[X+X^{-1}]/(X^{2N}+1)`: provides N real slots.

For additional information and other optional parameters see `lattigo/schemes/ckks/params.go`

#### Example

```json
"Scheme":{
	"LogN":16,
	"LogQ": [60, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45],
	"LogP": [56, 55, 55, 55],
	"Xe": {"Type": "DiscreteGaussian", "Sigma": 3.2, "Bound": 19.0},
	"Xs": {"Type": "Ternary", "H": 192},
	"RingType": "ConjugateInvariant",
	"LogDefaultScale":45
}
```

### Scheme Evaluation Keys

The user can ask the relinearization key and which Galois key to generate.
The user can either specify the Galois keys by their acting rotation on an encoded plaintext or directly with the Galois element.
Galois elements can be obtained from a rotation by calling `.GaloisElement(k int)` on the scheme parameters.

#### Example

```json
"EvaluationKeys":{
	"Relinearization:": true,
	"Rotations": [1],
	"GaloisElements":[65535]
}
```

### Bootstrapping Parameters

To enable bootstrapping, set `Bootstrapping.Enable: true`.

By default the bootstrapping works without having to set any optional field and provide the following performance:
 - Depth 4 for CoeffsToSlots
 - Depth 8 for EvalMod
 - Depth 3 for SlotsToCoeffs
for a total depth of 15 and a bit consumption of 821
A precision, for complex values with both real and imaginary parts uniformly distributed in -1, 1 of
 - 27.25 bits for H=192
 - 23.8 bits for H=32768,
And a failure probability of 2^{-138.7} for 2^{15} slots.

There are up to 16 optional fields in the `bootstrapping.ParametersLiteral`, enabling fine customization of the bootstrapping.
For additional information about these optional fields, see `lattigo/he/hefloat/bootstrapping/paramters_literal.go`.

Since the bootstrapping parameters are built as an extension of the scheme parameters, it is up to the user to ensure that the produced bootstrapping parameters are secure. See comments in `bootstrapping.NewParametersFromLiteral`, `bootstrapping.EvaluationKeys.GenEvaluationKeys`) and the examples in `examples/single_party/applications/reals_bootstrapping`.

#### Example

```json
"Bootstrapping":{
	"Enable": true,
	"LogN": 16,
	"LogP": [61, 61, 61, 61],
	"Xs": {"Type": "Ternary", "H": 192},
	"Xe": {"Type": "DiscreteGaussian", "Sigma": 3.2, "Bound": 19.0},
	"LogSlots": 15,
	"CoeffsToSlotsFactorizationDepthAndLogScales": [[56], [56], [56], [56]],
	"SlotsToCoeffsFactorizationDepthAndLogScales": [[39], [39], [39]],
	"EvalModLogScale": 60,
	"EphemeralSecretWeight": 32,
	"IterationsParameters": {
		"BootstrappingPrecision": [27.5, 10.0],
		"ReservedPrimeBitSize": 20
	},
	"Mod1Type": 0,
	"LogMessageRatio": 8,
	"K": 16,
	"Mod1Degree": 30,
	"DoubleAngle": 3,
	"Mod1InvDegree": 0
}
```

## Imputing Your Solution

All you have to do is put your solution in the function `SolveTestcase`, which is located in the file `internal/solution/solution.go`.

## Testing Your Solution Locally

- `$ make test-all` to do an end-to-end test of your solution followed by a clean of the temporary files
- `$ make setup` to generate the keys and input ciphertext
- `$ make solution` to run the solution and verify it (assumes that the keys and input ciphertext have been generated)
- `$ make clean` to clean the temporary files

## Packaging & Submitting Your Solution

Simply create a `.zip` containing the folder `app` and submit it on the website.