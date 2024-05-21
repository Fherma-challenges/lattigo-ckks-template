package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"github.com/tuneinsight/lattigo/v5/ring"
	"github.com/tuneinsight/lattigo/v5/utils"
	"github.com/tuneinsight/lattigo/v5/utils/buffer"
)

type bootstrappingParametersLiteral struct {
	Enable                                      bool
	LogN                                        *int                                `json:",omitempty"` // Default: 16
	LogP                                        []int                               `json:",omitempty"` // Default: 61 * max(1, floor(sqrt(#Qi)))
	Xs                                          ring.DistributionParameters         `json:",omitempty"` // Default: ring.Ternary{H: 192}
	Xe                                          ring.DistributionParameters         `json:",omitempty"` // Default: rlwe.DefaultXe
	LogSlots                                    *int                                `json:",omitempty"` // Default: LogN-1
	CoeffsToSlotsFactorizationDepthAndLogScales [][]int                             `json:",omitempty"` // Default: [][]int{min(4, max(LogSlots, 1)) * 56}
	SlotsToCoeffsFactorizationDepthAndLogScales [][]int                             `json:",omitempty"` // Default: [][]int{min(3, max(LogSlots, 1)) * 39}
	EvalModLogScale                             *int                                `json:",omitempty"` // Default: 60
	EphemeralSecretWeight                       *int                                `json:",omitempty"` // Default: 32
	IterationsParameters                        *bootstrapping.IterationsParameters `json:",omitempty"` // Default: nil (default starting level of 0 and 1 iteration)
	Mod1Type                                    hefloat.Mod1Type                    `json:",omitempty"` // Default: hefloat.CosDiscrete
	LogMessageRatio                             *int                                `json:",omitempty"` // Default: 8
	K                                           *int                                `json:",omitempty"` // Default: 16
	Mod1Degree                                  *int                                `json:",omitempty"` // Default: 30
	DoubleAngle                                 *int                                `json:",omitempty"` // Default: 3
	Mod1InvDegree                               *int                                `json:",omitempty"` // Default: 0
}

func (b bootstrappingParametersLiteral) ToLiteral() bootstrapping.ParametersLiteral{
	btpParamsLit := bootstrapping.ParametersLiteral{}
	btpParamsLit.LogN = b.LogN
	btpParamsLit.LogP = b.LogP
	btpParamsLit.Xs = b.Xs
	btpParamsLit.Xe = b.Xe
	btpParamsLit.LogSlots = b.LogSlots
	btpParamsLit.CoeffsToSlotsFactorizationDepthAndLogScales = b.CoeffsToSlotsFactorizationDepthAndLogScales
	btpParamsLit.SlotsToCoeffsFactorizationDepthAndLogScales = b.SlotsToCoeffsFactorizationDepthAndLogScales
	btpParamsLit.EvalModLogScale = b.EvalModLogScale
	btpParamsLit.EphemeralSecretWeight = b.EphemeralSecretWeight
	btpParamsLit.IterationsParameters = b.IterationsParameters
	btpParamsLit.Mod1Type = b.Mod1Type
	btpParamsLit.LogMessageRatio = b.LogMessageRatio
	btpParamsLit.K = b.K
	btpParamsLit.Mod1Degree = b.Mod1Degree
	btpParamsLit.DoubleAngle = b.DoubleAngle
	btpParamsLit.Mod1InvDegree = b.Mod1InvDegree
	return btpParamsLit
}

func (b bootstrappingParametersLiteral) BinarySize() int {
	data, _ := json.Marshal(b)
	return len(data)+4
}

func (p bootstrappingParametersLiteral) WriteTo(w io.Writer) (n int64, err error) {
	switch w := w.(type) {
	case buffer.Writer:

		var data []byte
		if data, err = json.Marshal(p); err != nil{
			return
		}

		if n, err = buffer.WriteAsUint32[int](w, len(data)); err != nil {
			return n, fmt.Errorf("buffer.WriteAsUint32[int]: %w", err)
		}

		var inc int
		if inc, err = w.Write(data); err != nil {
			return int64(n), fmt.Errorf("io.Write.Write: %w", err)
		}

		n += int64(inc)

		return n, w.Flush()
	default:
		return p.WriteTo(bufio.NewWriter(w))
	}
}

func (p *bootstrappingParametersLiteral) ReadFrom(r io.Reader) (n int64, err error) {

	switch r := r.(type) {
	case buffer.Reader:

		var size int
		if n, err = buffer.ReadAsUint32[int](r, &size); err != nil {
			return int64(n), fmt.Errorf("buffer.ReadAsUint64[int]: %w", err)
		}

		bytes := make([]byte, size)

		var inc int
		if inc, err = r.Read(bytes); err != nil {
			return n + int64(inc), fmt.Errorf("io.Reader.Read: %w", err)
		}

		return n + int64(inc), p.UnmarshalJSON(bytes)

	default:
		return p.ReadFrom(bufio.NewReader(r))
	}
}

func (b *bootstrappingParametersLiteral) UnmarshalJSON(data []byte) (err error) {
	var aux map[string]interface{}

	if err := json.Unmarshal(data, &aux); err != nil {
		panic(err)
	}

	if value, ok := aux["Enable"]; ok {
		b.Enable = value.(bool)
	} else {
		return fmt.Errorf("invalid Bootstrapping parameters: field [Enable] should be [true/false]")
	}

	if !b.Enable {
		return
	}

	if value, ok := aux["LogN"]; ok {
		switch value := value.(type) {
		case float64:
			b.LogN = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [LogN] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["LogP"]; ok {
		switch value := value.(type) {
		case []interface{}:
			b.LogP = make([]int, len(value))
			for i := range value {
				b.LogP[i] = int(value[i].(float64))
			}
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [LogP] type should be []interface{} but is %T", value)
		}
	}

	if value, ok := aux["Xs"]; ok {
		switch value := value.(type) {
		case map[string]interface{}:
			if b.Xs, err = ring.ParametersFromMap(value); err != nil {
				return
			}
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [Xs] type should be map[string]interface{} but is %T", value)
		}
	}

	if value, ok := aux["Xe"]; ok {
		switch value := value.(type) {
		case map[string]interface{}:
			if b.Xe, err = ring.ParametersFromMap(value); err != nil {
				return
			}
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [Xe] type should be map[string]interface{} but is %T", value)
		}
	}

	if value, ok := aux["LogSlots"]; ok {
		switch value := value.(type) {
		case float64:
			b.LogSlots = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [LogSlots] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["CoeffsToSlotsFactorizationDepthAndLogScales"]; ok {
		switch value := value.(type) {
		case []interface{}:
			b.CoeffsToSlotsFactorizationDepthAndLogScales = castInterface2DSlice[int](value)
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [CoeffsToSlotsFactorizationDepthAndLogScales] type should be []interface{} but is %T", value)
		}
	}

	if value, ok := aux["SlotsToCoeffsFactorizationDepthAndLogScales"]; ok {
		switch value := value.(type) {
		case []interface{}:
			b.SlotsToCoeffsFactorizationDepthAndLogScales = castInterface2DSlice[int](value)
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [SlotsToCoeffsFactorizationDepthAndLogScales] type should be []interface{} but is %T", value)
		}
	}

	if value, ok := aux["EvalModLogScale"]; ok {
		switch value := value.(type) {
		case float64:
			b.EvalModLogScale = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [EvalModLogScale] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["EphemeralSecretWeight"]; ok {
		switch value := value.(type) {
		case float64:
			b.EphemeralSecretWeight = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [EphemeralSecretWeight] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["IterationsParameters"]; ok {
		switch value := value.(type) {
		case map[string]interface{}:

			b.IterationsParameters = &bootstrapping.IterationsParameters{}

			if value0, ok := value["BootstrappingPrecision"]; ok {
				switch value0 := value0.(type) {
				case []interface{}:
					b.IterationsParameters.BootstrappingPrecision = castInterface1DSlice[float64](value0)
				default:
					return fmt.Errorf("invalid Bootstrapping parameters: field [IterationsParameters.BootstrappingPrecision] type should be []interface{} but is %T", value0)
				}
			} else {
				return fmt.Errorf("invalid Bootstrapping parameters: field [IterationsParameters.BootstrappingPrecision] cannot be empty if field [IterationsParameters] is not empty")
			}

			if value1, ok := value["ReservedPrimeBitSize"]; ok {
				switch value1 := value1.(type) {
				case float64:
					b.IterationsParameters.ReservedPrimeBitSize = int(value1)
				default:
					return fmt.Errorf("invalid Bootstrapping parameters: field [IterationsParameters.ReservedPrimeBitSize] type should be float64 but is %T", value1)
				}
			} else {
				return fmt.Errorf("invalid Bootstrapping parameters: field [IterationsParameters.ReservedPrimeBitSize] cannot be empty if field [IterationsParameters] is not empty")
			}

		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [EphemeralSecretWeight] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["Mod1Type"]; ok {
		switch value := value.(type) {
		case float64:
			b.Mod1Type = hefloat.Mod1Type(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [Mod1Type] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["LogMessageRatio"]; ok {
		switch value := value.(type) {
		case float64:
			b.LogMessageRatio = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [LogMessageRatio] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["K"]; ok {
		switch value := value.(type) {
		case float64:
			b.K = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [K] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["Mod1Degree"]; ok {
		switch value := value.(type) {
		case float64:
			b.Mod1Degree = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [Mod1Degree] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["DoubleAngle"]; ok {
		switch value := value.(type) {
		case float64:
			b.DoubleAngle = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [DoubleAngle] type should be float64 but is %T", value)
		}
	}

	if value, ok := aux["Mod1InvDegree"]; ok {
		switch value := value.(type) {
		case float64:
			b.Mod1InvDegree = utils.Pointy(int(value))
		default:
			return fmt.Errorf("invalid Bootstrapping parameters: field [Mod1InvDegree] type should be float64 but is %T", value)
		}
	}

	return
}

type number interface {
	int | float64
}

func castInterface2DSlice[T number](in []interface{}) (out [][]T) {
	out = make([][]T, len(in))
	for i := range in {
		out[i] = castInterface1DSlice[T](in[i].([]interface{}))
	}
	return
}

func castInterface1DSlice[T number](in []interface{}) (out []T) {
	out = make([]T, len(in))
	for i := range in {
		if x, ok := in[i].(float64); ok {
			out[i] = T(x)
		}
	}
	return
}
