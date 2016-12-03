package param

import "github.com/reyoung/GoPServer/protocol"
import flatbuffers "github.com/google/flatbuffers/go"
import "sync"
import "errors"

// ErrCannotParsePayload is the error when parameterBlock cannot parsing payload
var ErrCannotParsePayload = errors.New("Error when parsing request payload.")

var ErrParamBlockDimMismatch = errors.New("This Parameter block has been inited with different parameter dim")

var ErrParamBlockBufferInited = errors.New("Current ParamBlock Id has been inited")

type blockJob struct {
	req      *protocol.ParameterServerRequest
	builder  *flatbuffers.Builder
	lock     *sync.Mutex
	onFinish func(flatbuffers.UOffsetT)
}

type parameterBlock struct {
	buffer map[int8][]float64 // Buffers
	dim    uint32
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
				buffer: make(map[int8][]float64),
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

func createParameterBlockDone(name string, builder *flatbuffers.Builder, lock *sync.Mutex) flatbuffers.UOffsetT {
	lock.Lock()
	defer lock.Unlock()
	errOffset := builder.CreateString("ok")
	protocol.CreateResponseStart(builder)
	protocol.CreateResponseAddError(builder, errOffset)
	payloadOffset := protocol.CreateResponseEnd(builder)

	nameOffset := builder.CreateString(name)
	protocol.ParameterServerResponseStart(builder)
	protocol.ParameterServerResponseAddName(builder, nameOffset)
	protocol.ParameterServerResponseAddPayloadType(builder, protocol.ResponsePayLoadCreateResponse)
	protocol.ParameterServerResponseAddPayload(builder, payloadOffset)
	return protocol.ParameterServerResponseEnd(builder)
}

func (block *parameterBlock) process() error {
	table := new(flatbuffers.Table)
	createReq := new(protocol.CreateRequest)
	for {
		job, ok := <-block.job
		if !ok {
			break
		}

		switch job.req.PayloadType() {
		case protocol.RequestPayLoadCreateRequest:
			ok = job.req.Payload(table)
			if !ok {
				return ErrCannotParsePayload
			}

			createReq.Init(table.Bytes, table.Pos)
			if block.dim == 0 {
				block.dim = createReq.Dim()
			} else if block.dim != createReq.Dim() {
				return ErrParamBlockDimMismatch
			}

			_, ok = block.buffer[createReq.Id()]
			if ok {
				return ErrParamBlockBufferInited
			}
			buf := make([]float64, createReq.Size())
			for i := 0; i < len(buf); i++ {
				buf[i] = 0.0
			}
			block.buffer[createReq.Id()] = buf
			// CreateDone, Write Return Value.
			job.onFinish(createParameterBlockDone(string(job.req.Name()), job.builder, job.lock))
		}
	}
	return nil
}
