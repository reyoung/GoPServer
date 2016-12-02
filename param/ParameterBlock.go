package param

import "github.com/reyoung/GoPServer/protocol"
import flatbuffers "github.com/google/flatbuffers/go"
import "sync"

type blockJob struct {
	req      *protocol.ParameterServerRequest
	builder  *flatbuffers.Builder
	lock     *sync.Mutex
	onFinish func(flatbuffers.UOffsetT)
}

type parameterBlock struct {
	buffer map[byte][]float64 // Buffers
	job    chan *blockJob
}

type Parameters struct {
	params map[string]*parameterBlock
}

func New() *Parameters {
	return &Parameters{
		params: make(map[string]*parameterBlock),
	}
}

func (param *Parameters) DoJob(req *protocol.Requests) *flatbuffers.Builder {
	reqs := make([]*protocol.ParameterServerRequest, req.ReqLength())
	mtx := &sync.Mutex{}
	builder := flatbuffers.NewBuilder(0)
	resOffsets := make([]flatbuffers.UOffsetT, req.ReqLength())
	resOffsetsIndex := 0
	var wg sync.WaitGroup
	wg.Add(req.ReqLength())
	onFinish := func(idx flatbuffers.UOffsetT) {
		resOffsets[resOffsetsIndex] = idx
		resOffsetsIndex++
		wg.Done()
	}

	for i := 0; i < req.ReqLength(); i++ {
		reqs[i] = new(protocol.ParameterServerRequest)
		req.Req(reqs[i], i)
		n := string(reqs[i].Name())
		block, ok := param.params[n]
		if !ok {
			block = &parameterBlock{
				buffer: make(map[byte][]float64),
				job:    make(chan *blockJob, 10),
			}
			go block.process() // each parameter block use a goroutine.
			param.params[n] = block
		}
		block.job <- &blockJob{
			req:      reqs[i],
			builder:  builder,
			lock:     mtx,
			onFinish: onFinish,
		}
	}
	wg.Wait()

	protocol.ResponsesStartResVector(builder, resOffsetsIndex)
	for i := 0; i < resOffsetsIndex; i++ {
		builder.PrependUOffsetT(resOffsets[i])
	}
	reses := builder.EndVector(resOffsetsIndex)

	protocol.ResponsesStart(builder)
	protocol.ResponsesAddRes(builder, reses)
	protocol.ResponsesEnd(builder)
	return builder
}

func (block *parameterBlock) process() {

}
