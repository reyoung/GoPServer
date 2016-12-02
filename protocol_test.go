package main

import (
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/reyoung/GoPServer/protocol"
)

func TestFlatBuffer(t *testing.T) {
	builder := flatbuffers.NewBuilder(0)
	// Create two parameter block, w b
	nameW := builder.CreateString("w")

	protocol.CreateRequestStart(builder)
	protocol.CreateRequestAddSize(builder, 200*200)
	protocol.CreateRequestAddDim(builder, 200)
	protocol.CreateRequestAddId(builder, 0)
	payloadW := protocol.CreateRequestEnd(builder)

	protocol.ParameterServerRequestStart(builder)
	protocol.ParameterServerRequestAddName(builder, nameW)
	protocol.ParameterServerRequestAddPayloadType(builder, protocol.RequestPayLoadCreateRequest)
	protocol.ParameterServerRequestAddPayload(builder, payloadW)
	req0 := protocol.ParameterServerRequestEnd(builder)

	nameB := builder.CreateString("b")

	protocol.CreateRequestStart(builder)
	protocol.CreateRequestAddSize(builder, 200*1)
	protocol.CreateRequestAddDim(builder, 200)
	protocol.CreateRequestAddId(builder, 0)
	payloadB := protocol.CreateRequestEnd(builder)

	protocol.ParameterServerRequestStart(builder)
	protocol.ParameterServerRequestAddName(builder, nameB)
	protocol.ParameterServerRequestAddPayloadType(builder, protocol.RequestPayLoadCreateRequest)
	protocol.ParameterServerRequestAddPayload(builder, payloadB)
	req1 := protocol.ParameterServerRequestEnd(builder)

	protocol.RequestsStartReqVector(builder, 2)
	builder.PrependUOffsetT(req0)
	builder.PrependUOffsetT(req1)
	reqs := builder.EndVector(2)

	protocol.RequestsStart(builder)
	protocol.RequestsAddReq(builder, reqs)

	builder.Finish(protocol.RequestsEnd(builder))

	req := protocol.GetRootAsRequests(builder.FinishedBytes(), 0)
	t.Error(req.ReqLength())
}
