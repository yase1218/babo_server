package pbhandler

import (
	"babo/pb_code"
	"babo/utility/nw"

	"google.golang.org/protobuf/proto"
)

type HandlerFunc func(nw.IClient, proto.Message) (proto.Message, error)

type handleItem struct {
	key     pb_code.MsgId
	handler HandlerFunc
}

type pbHandler struct {
	handlers map[pb_code.MsgId]*handleItem
}

var Handler = &pbHandler{
	handlers: make(map[pb_code.MsgId]*handleItem),
}

func (h *pbHandler) register(msg_id pb_code.MsgId, handler HandlerFunc) {
	h.handlers[msg_id] = &handleItem{
		key:     msg_id,
		handler: handler,
	}
}

func (h *pbHandler) hans_handler(msg_id pb_code.MsgId) bool {
	_, ok := h.handlers[msg_id]
	return ok
}

func (h *pbHandler) get_handler(msg_id pb_code.MsgId) HandlerFunc {
	if item, ok := h.handlers[msg_id]; ok {
		return item.handler
	}

	return nil
}

func Register(msg_id pb_code.MsgId, handler HandlerFunc) {
	Handler.register(msg_id, handler)
}

func HasHandler(msg_id pb_code.MsgId) bool {
	return Handler.hans_handler(msg_id)
}

func GetHandler(msg_id pb_code.MsgId) HandlerFunc {
	return Handler.get_handler(msg_id)
}
