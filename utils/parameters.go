package utils

import (
	"io"
	"fmt"
	"bufio"
	"encoding/json"

	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"github.com/tuneinsight/lattigo/v5/utils/buffer"
)

type Parameters struct {
	Scheme        hefloat.Parameters
	Bootstrapping *bootstrapping.Parameters
	btpLiteral bootstrappingParametersLiteral
}

func (p Parameters) BinarySize() int {
	data, _ := p.MarshalJSON()
	return len(data) + 4
}

func (p Parameters) WriteTo(w io.Writer) (n int64, err error) {
	switch w := w.(type) {
	case buffer.Writer:

		var data []byte
		if data, err = p.MarshalJSON(); err != nil{
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

func (p *Parameters) ReadFrom(r io.Reader) (n int64, err error) {

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

func (p Parameters) MarshalJSON() (data []byte, err error){
	aux := struct{
		Scheme hefloat.Parameters
		Bootstrapping bootstrappingParametersLiteral
	}{
		Scheme: p.Scheme,
		Bootstrapping: p.btpLiteral,
	}

	return json.Marshal(aux)
}

func (p *Parameters) UnmarshalJSON(data []byte) (err error) {

	aux := struct {
		Scheme        hefloat.Parameters
		Bootstrapping bootstrappingParametersLiteral
	}{}

	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p.Scheme = aux.Scheme
	p.btpLiteral = aux.Bootstrapping

	if aux.Bootstrapping.Enable {
		
		var btpParams bootstrapping.Parameters
		if btpParams, err = bootstrapping.NewParametersFromLiteral(p.Scheme, aux.Bootstrapping.ToLiteral()); err != nil {
			return
		}

		p.Bootstrapping = &btpParams 
	}

	return
}
