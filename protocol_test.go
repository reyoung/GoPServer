package main

import (
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/reyoung/GoPServer/protocol"
	"github.com/stretchr/testify/assert"
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
	// builder.StartVector(flatbuffers.SizeUOffsetT, 2, 0)
	builder.PrependUOffsetT(req1)
	builder.PrependUOffsetT(req0)
	reqs := builder.EndVector(2)

	protocol.RequestsStart(builder)
	protocol.RequestsAddReq(builder, reqs)

	builder.Finish(protocol.RequestsEnd(builder))

	req := protocol.GetRootAsRequests(builder.FinishedBytes(), 0)
	assert.Equal(t, req.ReqLength(), 2)
	r := new(protocol.ParameterServerRequest)
	reqPayload := new(protocol.CreateRequest)
	assert.True(t, req.Req(r, 0))
	assert.Equal(t, "w", string(r.Name()[:]))
	assert.Equal(t, byte(protocol.RequestPayLoadCreateRequest), r.PayloadType())
	tmp := new(flatbuffers.Table)
	assert.True(t, r.Payload(tmp))
	reqPayload.Init(tmp.Bytes, tmp.Pos)

	assert.Equal(t, uint32(200*200), reqPayload.Size())
	assert.Equal(t, uint32(200), reqPayload.Dim())
}
