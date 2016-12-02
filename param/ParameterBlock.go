package param

import "github.com/reyoung/GoPServer/protocol"
import flatbuffers "github.com/google/flatbuffers/go"

type ParameterBlock struct {
	buffer map[byte][]float64 // Buffers
}

type Parameters struct {
	params map[string]*ParameterBlock
}

func New() *Parameters {
	return &Parameters{
		params: make(map[string]*ParameterBlock),
	}
}

func (param *Parameters) DoJob(req *protocol.Requests) *flatbuffers.Builder {
}

func (block *ParameterBlock) prealloc(req *protocol.ParameterServerRequest, builder *flatbuffers.Builder) {
}
