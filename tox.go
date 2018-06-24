package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

void callbackFriendRequestWrapperForC(Tox *, uint8_t *, uint8_t *, uint16_t, void*);
void callbackFriendMessageWrapperForC(Tox *, uint32_t, int, uint8_t*, uint32_t, void*);
void callbackFriendNameWrapperForC(Tox *, uint32_t, uint8_t*, uint32_t, void*);
void callbackFriendStatusMessageWrapperForC(Tox *, uint32_t, uint8_t*, uint32_t, void*);
void callbackFriendStatusWrapperForC(Tox *, uint32_t, int, void*);
void callbackFriendConnectionStatusWrapperForC(Tox *, uint32_t, int, void*);
void callbackFriendTypingWrapperForC(Tox *, uint32_t, uint8_t, void*);
void callbackFriendReadReceiptWrapperForC(Tox *, uint32_t, uint32_t, void*);
void callbackFriendLossyPacketWrapperForC(Tox *, uint32_t, uint8_t*, size_t, void*);
void callbackFriendLosslessPacketWrapperForC(Tox *, uint32_t, uint8_t*, size_t, void*);
void callbackSelfConnectionStatusWrapperForC(Tox *, int, void*);
void callbackFileRecvControlWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number,
                                      TOX_FILE_CONTROL control, void *user_data);
void callbackFileRecvWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint32_t kind,
                               uint64_t file_size, uint8_t *filename, size_t filename_length, void *user_data);
void callbackFileRecvChunkWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position,
                                    uint8_t *data, size_t length, void *user_data);
void callbackFileChunkRequestWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position,
                                       size_t length, void *user_data);

// fix nouse compile warning
static inline __attribute__((__unused__)) void fixnousetox() { }

*/
import "C"
import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"

	deadlock "github.com/sasha-s/go-deadlock"
)

// friend callback type
type cb_friend_request_ftype func(this *Tox, pubkey string, message string, userData interface{})
type cb_friend_message_ftype func(this *Tox, friendNumber uint32, message string, userData interface{})
type cb_friend_name_ftype func(this *Tox, friendNumber uint32, newName string, userData interface{})
type cb_friend_status_message_ftype func(this *Tox, friendNumber uint32, newStatus string, userData interface{})
type cb_friend_status_ftype func(this *Tox, friendNumber uint32, status int, userData interface{})
type cb_friend_connection_status_ftype func(this *Tox, friendNumber uint32, status int, userData interface{})
type cb_friend_typing_ftype func(this *Tox, friendNumber uint32, isTyping uint8, userData interface{})
type cb_friend_read_receipt_ftype func(this *Tox, friendNumber uint32, receipt uint32, userData interface{})
type cb_friend_lossy_packet_ftype func(this *Tox, friendNumber uint32, data string, userData interface{})
type cb_friend_lossless_packet_ftype func(this *Tox, friendNumber uint32, data string, userData interface{})

// self callback type
type cb_self_connection_status_ftype func(this *Tox, status int, userData interface{})

// file callback type
type cb_file_recv_control_ftype func(this *Tox, friendNumber uint32, fileNumber uint32,
	control int, userData interface{})
type cb_file_recv_ftype func(this *Tox, friendNumber uint32, fileNumber uint32, kind uint32, fileSize uint64,
	fileName string, userData interface{})
type cb_file_recv_chunk_ftype func(this *Tox, friendNumber uint32, fileNumber uint32, position uint64,
	data []byte, userData interface{})
type cb_file_chunk_request_ftype func(this *Tox, friend_number uint32, file_number uint32, position uint64,
	length int, user_data interface{})

type Tox struct {
	opts       *ToxOptions
	toxcore    *C.Tox // save C.Tox
	threadSafe bool
	mu         deadlock.RWMutex
	// mu sync.RWMutex

	// some callbacks, should be private
	cb_friend_requests           map[unsafe.Pointer]interface{}
	cb_friend_messages           map[unsafe.Pointer]interface{}
	cb_friend_names              map[unsafe.Pointer]interface{}
	cb_friend_status_messages    map[unsafe.Pointer]interface{}
	cb_friend_statuss            map[unsafe.Pointer]interface{}
	cb_friend_connection_statuss map[unsafe.Pointer]interface{}
	cb_friend_typings            map[unsafe.Pointer]interface{}
	cb_friend_read_receipts      map[unsafe.Pointer]interface{}
	cb_friend_lossy_packets      map[unsafe.Pointer]interface{}
	cb_friend_lossless_packets   map[unsafe.Pointer]interface{}
	cb_self_connection_statuss   map[unsafe.Pointer]interface{}

	cb_conference_invites            map[unsafe.Pointer]interface{}
	cb_conference_messages           map[unsafe.Pointer]interface{}
	cb_conference_actions            map[unsafe.Pointer]interface{}
	cb_conference_titles             map[unsafe.Pointer]interface{}
	cb_conference_peer_names         map[unsafe.Pointer]interface{}
	cb_conference_peer_list_changeds map[unsafe.Pointer]interface{}

	cb_file_recv_controls  map[unsafe.Pointer]interface{}
	cb_file_recvs          map[unsafe.Pointer]interface{}
	cb_file_recv_chunks    map[unsafe.Pointer]interface{}
	cb_file_chunk_requests map[unsafe.Pointer]interface{}

	cb_iterate_data              interface{}
	cb_conference_message_setted bool

	hooks  callHookMethods
	cbevts []func() // no need lock
}

var cbUserDatas = newUserData()

//export callbackFriendRequestWrapperForC
func callbackFriendRequestWrapperForC(m *C.Tox, a0 *C.uint8_t, a1 *C.uint8_t, a2 C.uint16_t, a3 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_requests {
		pubkey_b := C.GoBytes(unsafe.Pointer(a0), C.int(PUBLIC_KEY_SIZE))
		pubkey := hex.EncodeToString(pubkey_b)
		pubkey = strings.ToUpper(pubkey)
		message_b := C.GoBytes(unsafe.Pointer(a1), C.int(a2))
		message := string(message_b)
		cbfn := *(*cb_friend_request_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, pubkey, message, ud) })
	}
}

func (this *Tox) CallbackFriendRequest(cbfn cb_friend_request_ftype, userData interface{}) {
	this.CallbackFriendRequestAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendRequestAdd(cbfn cb_friend_request_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_requests[cbfnp]; ok {
		return
	}
	this.cb_friend_requests[cbfnp] = userData

	C.tox_callback_friend_request(this.toxcore, (*C.tox_friend_request_cb)(C.callbackFriendRequestWrapperForC))
}

//export callbackFriendMessageWrapperForC
func callbackFriendMessageWrapperForC(m *C.Tox, a0 C.uint32_t, mtype C.int,
	a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_messages {
		message_ := C.GoStringN((*C.char)(unsafe.Pointer(a1)), (C.int)(a2))
		cbfn := *(*cb_friend_message_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), message_, ud) })
	}
}

func (this *Tox) CallbackFriendMessage(cbfn cb_friend_message_ftype, userData interface{}) {
	this.CallbackFriendMessageAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendMessageAdd(cbfn cb_friend_message_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_messages[cbfnp]; ok {
		return
	}
	this.cb_friend_messages[cbfnp] = userData

	C.tox_callback_friend_message(this.toxcore, (*C.tox_friend_message_cb)(C.callbackFriendMessageWrapperForC))
}

//export callbackFriendNameWrapperForC
func callbackFriendNameWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_names {
		name := C.GoStringN((*C.char)((unsafe.Pointer)(a1)), C.int(a2))
		cbfn := *(*cb_friend_name_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), name, ud) })
	}
}

func (this *Tox) CallbackFriendName(cbfn cb_friend_name_ftype, userData interface{}) {
	this.CallbackFriendNameAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendNameAdd(cbfn cb_friend_name_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_names[cbfnp]; ok {
		return
	}
	this.cb_friend_names[cbfnp] = userData

	C.tox_callback_friend_name(this.toxcore, (*C.tox_friend_name_cb)(C.callbackFriendNameWrapperForC))
}

//export callbackFriendStatusMessageWrapperForC
func callbackFriendStatusMessageWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_status_messages {
		statusText := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(a2))
		cbfn := *(*cb_friend_status_message_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), statusText, ud) })
	}
}

func (this *Tox) CallbackFriendStatusMessage(cbfn cb_friend_status_message_ftype, userData interface{}) {
	this.CallbackFriendStatusMessageAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendStatusMessageAdd(cbfn cb_friend_status_message_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_status_messages[cbfnp]; ok {
		return
	}
	this.cb_friend_status_messages[cbfnp] = userData

	C.tox_callback_friend_status_message(this.toxcore, (*C.tox_friend_status_message_cb)(C.callbackFriendStatusMessageWrapperForC))
}

//export callbackFriendStatusWrapperForC
func callbackFriendStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.int, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_statuss {
		cbfn := *(*cb_friend_status_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), int(a1), ud) })
	}
}

func (this *Tox) CallbackFriendStatus(cbfn cb_friend_status_ftype, userData interface{}) {
	this.CallbackFriendStatusAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendStatusAdd(cbfn cb_friend_status_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_statuss[cbfnp]; ok {
		return
	}
	this.cb_friend_statuss[cbfnp] = userData

	C.tox_callback_friend_status(this.toxcore, (*C.tox_friend_status_cb)(C.callbackFriendStatusWrapperForC))
}

//export callbackFriendConnectionStatusWrapperForC
func callbackFriendConnectionStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.int, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_connection_statuss {
		cbfn := *(*cb_friend_connection_status_ftype)((unsafe.Pointer)(cbfni))
		this.putcbevts(func() { cbfn(this, uint32(a0), int(a1), ud) })
	}
}

func (this *Tox) CallbackFriendConnectionStatus(cbfn cb_friend_connection_status_ftype, userData interface{}) {
	this.CallbackFriendConnectionStatusAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendConnectionStatusAdd(cbfn cb_friend_connection_status_ftype, userData interface{}) {
	cbfnp := unsafe.Pointer(&cbfn)
	if _, ok := this.cb_friend_connection_statuss[cbfnp]; ok {
		return
	}
	this.cb_friend_connection_statuss[cbfnp] = userData

	C.tox_callback_friend_connection_status(this.toxcore, (*C.tox_friend_connection_status_cb)(C.callbackFriendConnectionStatusWrapperForC))
}

//export callbackFriendTypingWrapperForC
func callbackFriendTypingWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint8_t, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_typings {
		cbfn := *(*cb_friend_typing_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), uint8(a1), ud) })
	}
}

func (this *Tox) CallbackFriendTyping(cbfn cb_friend_typing_ftype, userData interface{}) {
	this.CallbackFriendTypingAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendTypingAdd(cbfn cb_friend_typing_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_typings[cbfnp]; ok {
		return
	}
	this.cb_friend_typings[cbfnp] = userData

	C.tox_callback_friend_typing(this.toxcore, (*C.tox_friend_typing_cb)(C.callbackFriendTypingWrapperForC))
}

//export callbackFriendReadReceiptWrapperForC
func callbackFriendReadReceiptWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_read_receipts {
		cbfn := *(*cb_friend_read_receipt_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), uint32(a1), ud) })
	}
}

func (this *Tox) CallbackFriendReadReceipt(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
	this.CallbackFriendReadReceiptAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendReadReceiptAdd(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_read_receipts[cbfnp]; ok {
		return
	}
	this.cb_friend_read_receipts[cbfnp] = userData

	C.tox_callback_friend_read_receipt(this.toxcore, (*C.tox_friend_read_receipt_cb)(C.callbackFriendReadReceiptWrapperForC))
}

//export callbackFriendLossyPacketWrapperForC
func callbackFriendLossyPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_lossy_packets {
		cbfn := *(*cb_friend_lossy_packet_ftype)(cbfni)
		msg := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(len))
		this.putcbevts(func() { cbfn(this, uint32(a0), msg, ud) })
	}
}

func (this *Tox) CallbackFriendLossyPacket(cbfn cb_friend_lossy_packet_ftype, userData interface{}) {
	this.CallbackFriendLossyPacketAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendLossyPacketAdd(cbfn cb_friend_lossy_packet_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_lossy_packets[cbfnp]; ok {
		return
	}
	this.cb_friend_lossy_packets[cbfnp] = userData

	C.tox_callback_friend_lossy_packet(this.toxcore, (*C.tox_friend_lossy_packet_cb)(C.callbackFriendLossyPacketWrapperForC))
}

//export callbackFriendLosslessPacketWrapperForC
func callbackFriendLosslessPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_friend_lossless_packets {
		cbfn := *(*cb_friend_lossless_packet_ftype)(cbfni)
		msg := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(len))
		this.putcbevts(func() { cbfn(this, uint32(a0), msg, ud) })
	}
}

func (this *Tox) CallbackFriendLosslessPacket(cbfn cb_friend_lossless_packet_ftype, userData interface{}) {
	this.CallbackFriendLosslessPacketAdd(cbfn, userData)
}
func (this *Tox) CallbackFriendLosslessPacketAdd(cbfn cb_friend_lossless_packet_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_friend_lossless_packets[cbfnp]; ok {
		return
	}
	this.cb_friend_lossless_packets[cbfnp] = userData

	C.tox_callback_friend_lossless_packet(this.toxcore, (*C.tox_friend_lossless_packet_cb)(C.callbackFriendLosslessPacketWrapperForC))
}

//export callbackSelfConnectionStatusWrapperForC
func callbackSelfConnectionStatusWrapperForC(m *C.Tox, status C.int, a2 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_self_connection_statuss {
		cbfn := *(*cb_self_connection_status_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, int(status), ud) })
	}
}

func (this *Tox) CallbackSelfConnectionStatus(cbfn cb_self_connection_status_ftype, userData interface{}) {
	this.CallbackSelfConnectionStatusAdd(cbfn, userData)
}
func (this *Tox) CallbackSelfConnectionStatusAdd(cbfn cb_self_connection_status_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_self_connection_statuss[cbfnp]; ok {
		return
	}
	this.cb_self_connection_statuss[cbfnp] = userData

	C.tox_callback_self_connection_status(this.toxcore, (*C.tox_self_connection_status_cb)(C.callbackSelfConnectionStatusWrapperForC))
}

//export callbackFileRecvControlWrapperForC
func callbackFileRecvControlWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	control C.TOX_FILE_CONTROL, userData unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_file_recv_controls {
		cbfn := *(*cb_file_recv_control_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(friendNumber), uint32(fileNumber), int(control), ud) })
	}
}

func (this *Tox) CallbackFileRecvControl(cbfn cb_file_recv_control_ftype, userData interface{}) {
	this.CallbackFileRecvControlAdd(cbfn, userData)
}
func (this *Tox) CallbackFileRecvControlAdd(cbfn cb_file_recv_control_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_file_recv_controls[cbfnp]; ok {
		return
	}
	this.cb_file_recv_controls[cbfnp] = userData

	C.tox_callback_file_recv_control(this.toxcore, (*C.tox_file_recv_control_cb)(C.callbackFileRecvControlWrapperForC))
}

//export callbackFileRecvWrapperForC
func callbackFileRecvWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, kind C.uint32_t,
	fileSize C.uint64_t, fileName *C.uint8_t, fileNameLength C.size_t, userData unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_file_recvs {
		cbfn := *(*cb_file_recv_ftype)(cbfni)
		fileName_ := C.GoStringN((*C.char)(unsafe.Pointer(fileName)), C.int(fileNameLength))
		this.putcbevts(func() {
			cbfn(this, uint32(friendNumber), uint32(fileNumber), uint32(kind),
				uint64(fileSize), fileName_, ud)
		})
	}
}

func (this *Tox) CallbackFileRecv(cbfn cb_file_recv_ftype, userData interface{}) {
	this.CallbackFileRecvAdd(cbfn, userData)
}
func (this *Tox) CallbackFileRecvAdd(cbfn cb_file_recv_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_file_recvs[cbfnp]; ok {
		return
	}
	this.cb_file_recvs[cbfnp] = userData

	C.tox_callback_file_recv(this.toxcore, (*C.tox_file_recv_cb)(C.callbackFileRecvWrapperForC))
}

//export callbackFileRecvChunkWrapperForC
func callbackFileRecvChunkWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	position C.uint64_t, data *C.uint8_t, length C.size_t, userData unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_file_recv_chunks {
		cbfn := *(*cb_file_recv_chunk_ftype)(cbfni)
		data_ := C.GoBytes((unsafe.Pointer)(data), C.int(length))
		this.putcbevts(func() { cbfn(this, uint32(friendNumber), uint32(fileNumber), uint64(position), data_, ud) })
	}
}

func (this *Tox) CallbackFileRecvChunk(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
	this.CallbackFileRecvChunkAdd(cbfn, userData)
}
func (this *Tox) CallbackFileRecvChunkAdd(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_file_recv_chunks[cbfnp]; ok {
		return
	}
	this.cb_file_recv_chunks[cbfnp] = userData

	C.tox_callback_file_recv_chunk(this.toxcore, (*C.tox_file_recv_chunk_cb)(C.callbackFileRecvChunkWrapperForC))
}

//export callbackFileChunkRequestWrapperForC
func callbackFileChunkRequestWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	position C.uint64_t, length C.size_t, userData unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_file_chunk_requests {
		cbfn := *(*cb_file_chunk_request_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(friendNumber), uint32(fileNumber), uint64(position), int(length), ud) })
	}
}

func (this *Tox) CallbackFileChunkRequest(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	this.CallbackFileChunkRequestAdd(cbfn, userData)
}
func (this *Tox) CallbackFileChunkRequestAdd(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_file_chunk_requests[cbfnp]; ok {
		return
	}
	this.cb_file_chunk_requests[cbfnp] = userData

	C.tox_callback_file_chunk_request(this.toxcore, (*C.tox_file_chunk_request_cb)(C.callbackFileChunkRequestWrapperForC))
}

func NewTox(opts *ToxOptions) (*Tox, error) {
	var tox = new(Tox)
	if opts != nil {
		tox.opts = opts
	} else {
		tox.opts = NewToxOptions()
	}

	toxopts := tox.opts.toCToxOptions()
	defer C.tox_options_free(toxopts)

	var cerr C.TOX_ERR_NEW

	var toxcore = C.tox_new(toxopts, &cerr)

	switch cerr {
	case C.TOX_ERR_NEW_OK:
		assert(toxcore != nil, "toxcore != nil")

	case C.TOX_ERR_NEW_NULL:
		return nil, fmt.Errorf("one of the arguments to the function was nil when it was not expected")

	case C.TOX_ERR_NEW_MALLOC:
		return nil, fmt.Errorf("the function was unable to allocate enough memory to store the internal structures for the Tox object")

	case C.TOX_ERR_NEW_PORT_ALLOC:
		return nil, fmt.Errorf("the function was unable to bind to a port. This may mean that all ports have already been bound, e.g. by other Tox instances, or it may mean a permission error. You may be able to gather more information from errno")

	case C.TOX_ERR_NEW_PROXY_BAD_TYPE:
		return nil, fmt.Errorf("proxyType was invalid")

	case C.TOX_ERR_NEW_PROXY_BAD_HOST:
		return nil, fmt.Errorf("proxyType was valid but the proxyHost passed had an invalid format or was nil")

	case C.TOX_ERR_NEW_PROXY_BAD_PORT:
		return nil, fmt.Errorf("proxyType was valid, but the proxyPort was invalid")

	case C.TOX_ERR_NEW_PROXY_NOT_FOUND:
		return nil, fmt.Errorf("the proxy address passed could not be resolved")

	case C.TOX_ERR_NEW_LOAD_ENCRYPTED:
		return nil, fmt.Errorf("the byte array to be loaded contained an encrypted save")

	case C.TOX_ERR_NEW_LOAD_BAD_FORMAT:
		return nil, fmt.Errorf("the data format was invalid. This can happen when loading data that was saved by an older version of Tox, or when the data has been corrupted. When loading from badly formatted data, some data may have been loaded, and the rest is discarded. Passing an invalid length parameter also causes this error")

	default:
		return nil, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}

	tox.toxcore = toxcore
	cbUserDatas.set(toxcore, tox)

	tox.cb_friend_requests = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_messages = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_names = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_status_messages = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_statuss = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_connection_statuss = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_typings = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_read_receipts = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_lossy_packets = make(map[unsafe.Pointer]interface{})
	tox.cb_friend_lossless_packets = make(map[unsafe.Pointer]interface{})
	tox.cb_self_connection_statuss = make(map[unsafe.Pointer]interface{})

	tox.cb_conference_invites = make(map[unsafe.Pointer]interface{})
	tox.cb_conference_messages = make(map[unsafe.Pointer]interface{})
	tox.cb_conference_actions = make(map[unsafe.Pointer]interface{})
	tox.cb_conference_titles = make(map[unsafe.Pointer]interface{})
	tox.cb_conference_peer_names = make(map[unsafe.Pointer]interface{})
	tox.cb_conference_peer_list_changeds = make(map[unsafe.Pointer]interface{})

	tox.cb_file_recv_controls = make(map[unsafe.Pointer]interface{})
	tox.cb_file_recvs = make(map[unsafe.Pointer]interface{})
	tox.cb_file_recv_chunks = make(map[unsafe.Pointer]interface{})
	tox.cb_file_chunk_requests = make(map[unsafe.Pointer]interface{})

	return tox, nil
}

func (this *Tox) Kill() {
	if this == nil || this.toxcore == nil {
		return
	}

	this.lock()
	defer this.unlock()

	cbUserDatas.del(this.toxcore)
	C.tox_kill(this.toxcore)
	this.toxcore = nil
}

// uint32_t tox_iteration_interval(Tox *tox);
func (this *Tox) IterationInterval() time.Duration {
	this.lock()
	defer this.unlock()

	return time.Duration(C.tox_iteration_interval(this.toxcore))
}

/* The main loop that needs to be run in intervals of tox_iteration_interval() ms. */
// void tox_iterate(Tox *tox);
// compatable with legacy version
func (this *Tox) Iterate() {
	this.lock()

	C.tox_iterate(
		this.toxcore, // Tox *tox
		nil,          // void *user_data
	)
	cbevts := this.cbevts
	this.cbevts = nil

	this.unlock()

	this.invokeCallbackEvents(cbevts)
}

// for toktok new method
func (this *Tox) Iterate2(userData interface{}) {
	this.lock()

	this.cb_iterate_data = userData
	C.tox_iterate(
		this.toxcore, // Tox *tox
		nil,          // void *user_data
	)
	this.cb_iterate_data = nil
	cbevts := this.cbevts
	this.cbevts = nil

	this.unlock()

	this.invokeCallbackEvents(cbevts)
}

func (this *Tox) invokeCallbackEvents(cbevts []func()) {
	for _, cbfn := range cbevts {
		cbfn()
	}
}

func (this *Tox) lock() {
	if this.opts.ThreadSafe {
		this.mu.Lock()
	}
}
func (this *Tox) unlock() {
	if this.opts.ThreadSafe {
		this.mu.Unlock()
	}
}

func (this *Tox) GetSavedataSize() int32 {
	return int32(C.tox_get_savedata_size(this.toxcore))
}

func (this *Tox) GetSavedata() []byte {
	var savedata = make([]byte, this.GetSavedataSize())

	C.tox_get_savedata(
		this.toxcore,               // const Tox *tox
		(*C.uint8_t)(&savedata[0]), // uint8_t *savedata
	)

	return savedata
}

/*
 * @param pubkey hex 64B length
 */
func (t *Tox) Bootstrap(addr string, port uint16, pubKey string) error {
	t.lock()
	defer t.unlock()

	if len(pubKey) != PUBLIC_KEY_SIZE*2 {
		return fmt.Errorf("invalid Public Key in size (%d): %d, %s", PUBLIC_KEY_SIZE*2, len(pubKey), pubKey)
	}
	pubkeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return fmt.Errorf("invalid Public Key: %s", pubKey)
	}

	var addrBytes = []byte(addr)
	var cerr C.TOX_ERR_BOOTSTRAP

	ok := bool(
		C.tox_bootstrap(
			t.toxcore, // Tox *tox
			(*C.char)(unsafe.Pointer(&addrBytes[0])), // const char *address
			C.uint16_t(port),                         // uint16_t port
			(*C.uint8_t)(&pubkeyBytes[0]),            // const uint8_t *public_key
			&cerr, // TOX_ERR_BOOTSTRAP *error
		),
	)

	switch cerr {
	case C.TOX_ERR_BOOTSTRAP_OK:
		assert(ok, "tox_bootstrap() return 'false' on success")

		return nil

	case C.TOX_ERR_BOOTSTRAP_NULL:
		return fmt.Errorf("one of the arguments to the function was false when it was not expected")

	case C.TOX_ERR_BOOTSTRAP_BAD_HOST:
		return fmt.Errorf("the address could not be resolved to an IP address, or the IP address passed was invalid: %s", addr)

	case C.TOX_ERR_BOOTSTRAP_BAD_PORT:
		return fmt.Errorf("the port passed was invalid. The valid port range is (1, 65535): %d", port)

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (t *Tox) AddTcpRelay(addr string, port uint16, pubKey string) error {
	t.lock()
	defer t.unlock()

	var csAddr = C.CString(addr)
	defer C.free(unsafe.Pointer(csAddr))

	pubkeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return err
	}
	if strings.ToUpper(hex.EncodeToString(pubkeyBytes)) != pubKey {
		return fmt.Errorf("wtf, hex enc/dec err")
	}

	var cerr C.TOX_ERR_BOOTSTRAP

	ok := bool(
		C.tox_add_tcp_relay(
			t.toxcore,                     // Tox *tox
			csAddr,                        // const char *address
			C.uint16_t(port),              // uint16_t port
			(*C.uint8_t)(&pubkeyBytes[0]), // const uint8_t *public_key
			&cerr, // TOX_ERR_BOOTSTRAP *error
		),
	)

	switch cerr {
	case C.TOX_ERR_BOOTSTRAP_OK:
		assert(ok, "tox_add_tcp_relay() return 'false' on success")

		return nil

	case C.TOX_ERR_BOOTSTRAP_NULL:
		return fmt.Errorf("one of the arguments to the function was false when it was not expected")

	case C.TOX_ERR_BOOTSTRAP_BAD_HOST:
		return fmt.Errorf("the address could not be resolved to an IP address, or the IP address passed was invalid: %s", addr)

	case C.TOX_ERR_BOOTSTRAP_BAD_PORT:
		return fmt.Errorf("the port passed was invalid. The valid port range is (1, 65535): %d", port)

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetAddress() string {
	var addr [ADDRESS_SIZE]byte

	C.tox_self_get_address(
		this.toxcore,                           // const Tox *tox
		(*C.uint8_t)(unsafe.Pointer(&addr[0])), // uint8_t *address
	)

	return strings.ToUpper(hex.EncodeToString(addr[:]))
}

func (this *Tox) SelfGetConnectionStatus() int {
	return int(C.tox_self_get_connection_status(this.toxcore))
}

func (this *Tox) FriendAdd(friendId string, message string) (uint32, error) {
	this.lock()
	defer this.unlock()

	if len(message) == 1 || len(message) > MAX_FRIEND_REQUEST_LENGTH {
		return 0, fmt.Errorf("friend request message must be in range [1, %d]: %d", MAX_FRIEND_REQUEST_LENGTH, len(message))
	}

	friendIDBytes, err := hex.DecodeString(friendId)
	if err != nil {
		return 0, err
	}

	// If more than INT32_MAX friends are added, this function causes undefined
	// behavior.
	// TODO: Check current friend list size, and return error when needed.

	messageBytes := []byte(message)
	var cerr C.TOX_ERR_FRIEND_ADD

	friendNumber := uint32(
		C.tox_friend_add(
			this.toxcore,                    // Tox *tox
			(*C.uint8_t)(&friendIDBytes[0]), // const uint8_t *address
			(*C.uint8_t)(&messageBytes[0]),  // const uint8_t *message
			C.size_t(len(message)),          // size_t length
			&cerr, // TOX_ERR_FRIEND_ADD *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_ADD_OK:
		if friendNumber == math.MaxUint32 {
			return 0, fmt.Errorf("failed on add friend")
		}

		return friendNumber, nil

	case C.TOX_ERR_FRIEND_ADD_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_ADD_TOO_LONG:
		return 0, fmt.Errorf("the length of the friend request message exceeded MaxFriendRequestLength(%d): %d", MAX_FRIEND_REQUEST_LENGTH, len(message))

	case C.TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return 0, fmt.Errorf("the friend request message was empty")

	case C.TOX_ERR_FRIEND_ADD_OWN_KEY:
		return 0, fmt.Errorf("the friend address belongs to the sending client")

	case C.TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return 0, fmt.Errorf("friend request has already been sent, or the address belongs to a friend that is already on the friend list: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return 0, fmt.Errorf("the friend address checksum failed: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return 0, fmt.Errorf("the friend was already there, but the nospam value was different: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_MALLOC:
		return 0, fmt.Errorf("memory allocation failed when trying to increase the friend list size")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendAddNorequest(friendId string) (uint32, error) {
	this.lock()
	defer this.unlock()

	if len(friendId) != PUBLIC_KEY_SIZE*2 {
		return 0, fmt.Errorf("invalid friendId in size (%d): %d, %s", PUBLIC_KEY_SIZE*2, len(friendId), friendId)
	}
	friendIDBytes, err := hex.DecodeString(friendId)
	if err != nil {
		return 0, err
	}

	var cerr C.TOX_ERR_FRIEND_ADD
	friendNumber := uint32(
		C.tox_friend_add_norequest(
			this.toxcore,                    // Tox *tox
			(*C.uint8_t)(&friendIDBytes[0]), // const uint8_t *public_key
			&cerr, // TOX_ERR_FRIEND_ADD *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_ADD_OK:
		if friendNumber == math.MaxUint32 {
			return 0, fmt.Errorf("failed on add friend")
		}

		return friendNumber, nil

	case C.TOX_ERR_FRIEND_ADD_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_ADD_TOO_LONG:
		return 0, fmt.Errorf("the length of the friend request message exceeded MaxFriendRequestLength(%d)", MAX_FRIEND_REQUEST_LENGTH)

	case C.TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return 0, fmt.Errorf("the friend request message was empty")

	case C.TOX_ERR_FRIEND_ADD_OWN_KEY:
		return 0, fmt.Errorf("the friend address belongs to the sending client")

	case C.TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return 0, fmt.Errorf("friend request has already been sent, or the address belongs to a friend that is already on the friend list: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return 0, fmt.Errorf("the friend address checksum failed: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return 0, fmt.Errorf("the friend was already there, but the nospam value was different: %s", friendId)

	case C.TOX_ERR_FRIEND_ADD_MALLOC:
		return 0, fmt.Errorf("memory allocation failed when trying to increase the friend list size")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendByPublicKey(pubkey string) (uint32, error) {
	if len(pubkey) != ADDRESS_SIZE*2 {
		return 0, fmt.Errorf("invalid Public Key in size (%d): %d, %s", ADDRESS_SIZE*2, len(pubkey), pubkey)
	}
	pubkey_b, err := hex.DecodeString(pubkey)
	if err != nil {
		return 0, err
	}

	var cerr C.TOX_ERR_FRIEND_BY_PUBLIC_KEY

	friendNumber := uint32(
		C.tox_friend_by_public_key(
			this.toxcore,               // const Tox *tox
			(*C.uint8_t)(&pubkey_b[0]), // const uint8_t *public_key
			&cerr, // TOX_ERR_FRIEND_BY_PUBLIC_KEY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK:
		if friendNumber == math.MaxUint32 {
			return 0, fmt.Errorf("failed on get friend number by public key")
		}

		return friendNumber, nil

	case C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND:
		return 0, fmt.Errorf("no friend with the given Public Key exists on the friend list")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendGetPublicKey(friendNumber uint32) (string, error) {
	pubkey_b := make([]byte, PUBLIC_KEY_SIZE)
	var cerr C.TOX_ERR_FRIEND_GET_PUBLIC_KEY

	ok := bool(
		C.tox_friend_get_public_key(
			this.toxcore,               // const Tox *tox
			C.uint32_t(friendNumber),   // uint32_t friend_number
			(*C.uint8_t)(&pubkey_b[0]), // uint8_t *public_key
			&cerr, // TOX_ERR_FRIEND_GET_PUBLIC_KEY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK:
		assert(ok, "tox_friend_get_public_key() return 'false' on success")

		return strings.ToUpper(hex.EncodeToString(pubkey_b)), nil

	case C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("no friend with the given number exists on the friend list")

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendDelete(friendNumber uint32) error {
	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_FRIEND_DELETE

	ok := bool(
		C.tox_friend_delete(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_DELETE *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_DELETE_OK:
		assert(ok, "tox_friend_delete() return 'false' on success")

		return nil

	case C.TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND:
		return fmt.Errorf("there was no friend with the given friend number. No friends were deleted")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendGetConnectionStatus(friendNumber uint32) (int, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	connStatus := int(
		C.tox_friend_get_connection_status(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return connStatus, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return CONNECTION_NONE, fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return CONNECTION_NONE, fmt.Errorf("friendNumber did not designate a valid friend: %d", friendNumber)

	default:
		return CONNECTION_NONE, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendExists(friendNumber uint32) bool {
	return bool(
		C.tox_friend_exists(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
		),
	)
}

func (this *Tox) friendSendMessage(friendNumber uint32, message string, msgType C.TOX_MESSAGE_TYPE) (messageID uint32, err error) {
	if len(message) > MAX_MESSAGE_LENGTH {
		return 0, fmt.Errorf("length of message over ranged (max: %d): %d", MAX_MESSAGE_LENGTH, len(message))
	}

	this.lock()
	defer this.unlock()

	var _message = []byte(message)

	var cerr C.TOX_ERR_FRIEND_SEND_MESSAGE

	messageID = uint32(
		C.tox_friend_send_message(
			this.toxcore,               // Tox *tox
			C.uint32_t(friendNumber),   // uint32_t friend_number
			msgType,                    // TOX_MESSAGE_TYPE type
			(*C.uint8_t)(&_message[0]), // const uint8_t *message
			C.size_t(len(message)),     // size_t length
			&cerr, // TOX_ERR_FRIEND_SEND_MESSAGE *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_SEND_MESSAGE_OK:
		return messageID, nil

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND:
		return 0, fmt.Errorf("friend number did not designate a valid friend")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED:
		return 0, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ:
		return 0, fmt.Errorf("allocation error occurred while increasing the send queue size")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG:
		return 0, fmt.Errorf("message length exceeded MaxMessageLength(%d): %d", MAX_MESSAGE_LENGTH, len(message))

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY:
		return 0, fmt.Errorf("attempted to send a zero-length message")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendSendMessage(friendNumber uint32, message string) (messageID uint32, err error) {
	return this.friendSendMessage(friendNumber, message, C.TOX_MESSAGE_TYPE_NORMAL)
}

func (this *Tox) FriendSendAction(friendNumber uint32, action string) (uint32, error) {
	return this.friendSendMessage(friendNumber, action, C.TOX_MESSAGE_TYPE_ACTION)
}

func (this *Tox) SelfSetName(name string) error {
	if len(name) > MAX_NAME_LENGTH {
		return fmt.Errorf("length nickname is over ranged (max: %d): %d", MAX_NAME_LENGTH, len(name))
	}

	this.lock()
	defer this.unlock()

	var _name = []byte(name)
	var cerr C.TOX_ERR_SET_INFO

	ok := bool(
		C.tox_self_set_name(
			this.toxcore,            // Tox *tox
			(*C.uint8_t)(&_name[0]), // const uint8_t *name
			C.size_t(len(name)),     // size_t length
			&cerr,                   // TOX_ERR_SET_INFO *error
		),
	)

	switch cerr {
	case C.TOX_ERR_SET_INFO_OK:
		assert(ok, "tox_self_set_name() return 'false' on success")

		return nil

	case C.TOX_ERR_SET_INFO_NULL:
		return fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_SET_INFO_TOO_LONG:
		return fmt.Errorf("length exceeded maximum permissible size(%d): %d", MAX_NAME_LENGTH, len(name))

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetName() string {
	nlen := C.tox_self_get_name_size(this.toxcore)
	_name := make([]byte, nlen)

	C.tox_self_get_name(
		this.toxcore,                 // const Tox *tox
		(*C.uint8_t)(safeptr(_name)), // uint8_t *name
	)

	return string(_name)
}

func (this *Tox) FriendGetName(friendNumber uint32) (string, error) {
	nameSize, err := this.FriendGetNameSize(friendNumber)
	if err != nil {
		return "", err
	}
	_name := make([]byte, nameSize)
	var cerr C.TOX_ERR_FRIEND_QUERY

	ok := bool(
		C.tox_friend_get_name(
			this.toxcore,                 // const Tox *tox
			C.uint32_t(friendNumber),     // uint32_t friend_number
			(*C.uint8_t)(safeptr(_name)), // uint8_t *name
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		assert(ok, "tox_friend_get_name() return 'false' on success")

		return string(_name), nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return "", fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendGetNameSize(friendNumber uint32) (int, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	nameSize := int(
		C.tox_friend_get_name_size(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return nameSize, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return 0, fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return 0, fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetNameSize() int {
	// TODO: thread safe is needed
	return int(C.tox_self_get_name_size(this.toxcore))
}

func (this *Tox) SelfSetStatusMessage(status string) error {
	if len(status) > MAX_STATUS_MESSAGE_LENGTH {
		return fmt.Errorf("status is over ranged (max: %d): %d", MAX_STATUS_MESSAGE_LENGTH, len(status))
	}

	this.lock()
	defer this.unlock()

	var _status = []byte(status)
	var cerr C.TOX_ERR_SET_INFO

	ok := bool(
		C.tox_self_set_status_message(
			this.toxcore,              // Tox *tox
			(*C.uint8_t)(&_status[0]), // const uint8_t *status_message
			C.size_t(len(status)),     // size_t length
			&cerr, // TOX_ERR_SET_INFO *error
		),
	)

	switch cerr {
	case C.TOX_ERR_SET_INFO_OK:
		assert(ok, "tox_self_set_status_message() return 'false' on success")

		return nil

	case C.TOX_ERR_SET_INFO_NULL:
		return fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_SET_INFO_TOO_LONG:
		return fmt.Errorf("length exceeded maximum permissible size(%d): %d", MAX_STATUS_MESSAGE_LENGTH, len(status))

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfSetStatus(status uint8) {
	C.tox_self_set_status(
		this.toxcore,              // Tox *tox
		C.TOX_USER_STATUS(status), // TOX_USER_STATUS status
	)
}

func (this *Tox) FriendGetStatusMessageSize(friendNumber uint32) (int, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	size := int(
		C.tox_friend_get_status_message_size(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		// TODO:
		//
		// if size == C.SIZE_MAX {
		//     return 0, fmt.Errorf("invalid friend number")
		// }

		return size, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return 0, fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return 0, fmt.Errorf("friendNumber did not designate a valid friend: %d", friendNumber)

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetStatusMessageSize() int {
	return int(C.tox_self_get_status_message_size(this.toxcore))
}

func (this *Tox) FriendGetStatusMessage(friendNumber uint32) (string, error) {
	size, err := this.FriendGetStatusMessageSize(friendNumber)
	if err != nil {
		return "", err
	}
	_buf := make([]byte, size)
	var cerr C.TOX_ERR_FRIEND_QUERY

	ok := bool(
		C.tox_friend_get_status_message(
			this.toxcore,                // const Tox *tox
			C.uint32_t(friendNumber),    // uint32_t friend_number
			(*C.uint8_t)(safeptr(_buf)), // uint8_t *status_message
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		assert(ok, "tox_friend_get_status_message() return 'false' on success")

		return string(_buf[:]), nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return "", fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetStatusMessage() (string, error) {
	var _buf = make([]byte, this.SelfGetStatusMessageSize())

	C.tox_self_get_status_message(
		this.toxcore,                // const Tox *tox
		(*C.uint8_t)(safeptr(_buf)), // uint8_t *status_message
	)

	return string(_buf[:]), nil
}

func (this *Tox) FriendGetStatus(friendNumber uint32) (int, error) {
	friendNumberInFriendList := func(list []uint32, value uint32) bool {
		for _, v := range list {
			if value == v {
				return true
			}
		}
		return false
	}(this.SelfGetFriendList(), friendNumber)
	if !friendNumberInFriendList {
		return USER_STATUS_NONE, fmt.Errorf("friend is not in friend list")
	}

	var cerr C.TOX_ERR_FRIEND_QUERY

	friendStatus := int(
		C.tox_friend_get_status(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return friendStatus, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return USER_STATUS_NONE, fmt.Errorf("pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return USER_STATUS_NONE, fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return USER_STATUS_NONE, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetStatus() int {
	return int(C.tox_self_get_status(this.toxcore))
}

func (this *Tox) FriendGetLastOnline(friendNumber uint32) (time.Time, error) {
	var nullTime time.Time // defined for returning a null-time.
	var cerr C.TOX_ERR_FRIEND_GET_LAST_ONLINE

	timestamp := uint64(
		C.tox_friend_get_last_online(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_GET_LAST_ONLINE *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK:
		if timestamp == math.MaxUint64 {
			return nullTime, fmt.Errorf("unexpected timestamp")
		}

		return time.Unix(int64(timestamp), 0), nil

	case C.TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND:
		return nullTime, fmt.Errorf("no friend with the given number exists on the friend list")

	default:
		return nullTime, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfSetTyping(friendNumber uint32, typing bool) error {
	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_SET_TYPING

	ok := bool(
		C.tox_self_set_typing(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			C._Bool(typing),          // bool typing
			&cerr,                    // TOX_ERR_SET_TYPING *error
		),
	)

	switch cerr {
	case C.TOX_ERR_SET_TYPING_OK:
		assert(ok, "tox_self_set_typing() return 'false' on success")

		return nil

	case C.TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND:
		return fmt.Errorf("friend number did not designate a valid friend")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendGetTyping(friendNumber uint32) (bool, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	typing := bool(
		C.tox_friend_get_typing(
			this.toxcore,             // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return typing, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return false, fmt.Errorf("pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return false, fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return false, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) SelfGetFriendListSize() int {
	return int(C.tox_self_get_friend_list_size(this.toxcore))
}

func (this *Tox) SelfGetFriendList() []uint32 {
	size := this.SelfGetFriendListSize()
	vec := make([]uint32, size)
	if size == 0 {
		return vec // fast return on zero-sized friend list
	}

	C.tox_self_get_friend_list(
		this.toxcore,                           // const Tox *tox
		(*C.uint32_t)(unsafe.Pointer(&vec[0])), // uint32_t *friend_list
	)

	return vec
}

func (this *Tox) SelfGetNospam() uint32 {
	return uint32(C.tox_self_get_nospam(this.toxcore))
}

func (t *Tox) SelfGetNoSpamString() string {
	return fmt.Sprintf("%X", t.SelfGetNospam())
}

func (this *Tox) SelfSetNospam(nospam uint32) {
	this.lock()
	defer this.unlock()

	C.tox_self_set_nospam(
		this.toxcore,       // Tox *tox
		C.uint32_t(nospam), // uint32_t nospam
	)
}

func (t *Tox) SelfSetNoSpamString(noSpam string) error {
	if len(noSpam) != 8 {
		return fmt.Errorf("invalid NoSpam format, which should be a 8-char hex string")
	}

	var noSpamNum uint32
	_, err := fmt.Sscanf(noSpam, "%8x", &noSpamNum)
	if err != nil {
		return err
	}

	t.SelfSetNospam(noSpamNum)

	return nil
}

func (this *Tox) SelfGetPublicKey() string {
	var _pubkey [PUBLIC_KEY_SIZE]byte

	C.tox_self_get_public_key(
		this.toxcore,              // const Tox *tox
		(*C.uint8_t)(&_pubkey[0]), // uint8_t *public_key
	)

	return strings.ToUpper(hex.EncodeToString(_pubkey[:]))
}

func (this *Tox) SelfGetSecretKey() string {
	var _seckey [SECRET_KEY_SIZE]byte

	C.tox_self_get_secret_key(
		this.toxcore,              // const Tox *tox
		(*C.uint8_t)(&_seckey[0]), // uint8_t *secret_key
	)

	return strings.ToUpper(hex.EncodeToString(_seckey[:]))
}

// tox_lossy_***

func (this *Tox) FriendSendLossyPacket(friendNumber uint32, data string) error {
	if len(data) > MAX_CUSTOM_PACKET_SIZE {
		return fmt.Errorf("length of data is out of range (max: %d): %d", MAX_CUSTOM_PACKET_SIZE, len(data))
	}

	this.lock()
	defer this.unlock()

	var _data = []byte(data)
	if 200 > _data[0] || _data[0] < 254 {
		return fmt.Errorf("the first byte of data must be in the range 200-254: %d", _data[0])
	}

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET

	ok := bool(
		C.tox_friend_send_lossy_packet(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			(*C.uint8_t)(&_data[0]),  // const uint8_t *data
			C.size_t(len(data)),      // size_t length
			&cerr,                    // TOX_ERR_FRIEND_CUSTOM_PACKET *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		assert(ok, "tox_friend_send_lossy_packet() return 'false' on success")

		return nil

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_FOUND:
		return fmt.Errorf("friend number did not designate a valid friend: %d", friendNumber)

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED:
		return fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_INVALID:
		return fmt.Errorf("the first byte of data was not in the specified range for the packet type")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_EMPTY:
		return fmt.Errorf("attempted to send an empty packet")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_TOO_LONG:
		return fmt.Errorf("packet data length exceeded MaxCustomPacketSize(%d): %d", MAX_CUSTOM_PACKET_SIZE, len(data))

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		return fmt.Errorf("packet queue is full")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FriendSendLosslessPacket(friendNumber uint32, data string) error {
	if len(data) > MAX_CUSTOM_PACKET_SIZE {
		return fmt.Errorf("length of data is out of range (max: %d): %d", MAX_CUSTOM_PACKET_SIZE, len(data))
	}

	this.lock()
	defer this.unlock()

	var _data = []byte(data)
	if 160 > _data[0] || _data[0] > 191 {
		return fmt.Errorf("the first byte of data must be in the range 160-191: %d", _data[0])
	}

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET

	ok := bool(
		C.tox_friend_send_lossless_packet(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			(*C.uint8_t)(&_data[0]),  // const uint8_t *data
			C.size_t(len(data)),      // size_t length
			&cerr,                    // TOX_ERR_FRIEND_CUSTOM_PACKET *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		assert(ok, "tox_friend_send_lossless_packet() return 'false' on success")

		return nil

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_FOUND:
		return fmt.Errorf("friend number did not designate a valid friend: %d", friendNumber)

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED:
		return fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_INVALID:
		return fmt.Errorf("the first byte of data was not in the specified range for the packet type")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_EMPTY:
		return fmt.Errorf("attempted to send an empty packet")

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_TOO_LONG:
		return fmt.Errorf("packet data length exceeded MaxCustomPacketSize(%d): %d", MAX_CUSTOM_PACKET_SIZE, len(data))

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		return fmt.Errorf("packet queue is full")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// tox_callback_avatar_**

func (this *Tox) Hash(data string) (bool, string) {
	_data := []byte(data)
	_hash := make([]byte, HASH_LENGTH)

	ok := bool(
		C.tox_hash(
			(*C.uint8_t)(&_hash[0]), // uint8_t *hash
			(*C.uint8_t)(&_data[0]), // const uint8_t *data
			C.size_t(len(data)),     // size_t length
		),
	)

	// If hash is NULL or data is NULL while length is not 0 the
	// function returns false, otherwise it returns true.
	assert(ok && len(data) > 0, "tox_hash() return 'false' on success")

	return ok, string(_hash)
}

// tox_callback_file_***

func (this *Tox) FileControl(friendNumber uint32, fileNumber uint32, control int) (bool, error) {
	var cerr C.TOX_ERR_FILE_CONTROL

	ok := bool(
		C.tox_file_control(
			this.toxcore,                // Tox *tox
			C.uint32_t(friendNumber),    // uint32_t friend_number
			C.uint32_t(fileNumber),      // uint32_t file_number
			C.TOX_FILE_CONTROL(control), // TOX_FILE_CONTROL control
			&cerr, // TOX_ERR_FILE_CONTROL *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_CONTROL_OK:
		assert(ok, "tox_file_control() return 'false' on success")

		return ok, nil

	case C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND:
		return false, fmt.Errorf("friend_number passed did not designate a valid friend")

	case C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED:
		return false, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FILE_CONTROL_NOT_FOUND:
		return false, fmt.Errorf("no file transfer with the given file number was found for the given friend")

	case C.TOX_ERR_FILE_CONTROL_NOT_PAUSED:
		return false, fmt.Errorf("RESUME control was sent, but the file transfer is running normally")

	case C.TOX_ERR_FILE_CONTROL_DENIED:
		return false, fmt.Errorf("RESUME control was sent, but the file transfer was paused by the other party. Only the party that paused the transfer can resume it")

	case C.TOX_ERR_FILE_CONTROL_ALREADY_PAUSED:
		return false, fmt.Errorf("PAUSE control was sent, but the file transfer was already paused")

	case C.TOX_ERR_FILE_CONTROL_SENDQ:
		return false, fmt.Errorf("packet queue is full")

	default:
		return false, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FileSend(friendNumber uint32, kind uint32, fileSize uint64, fileId string, fileName string) (uint32, error) {
	if len(fileName) > MAX_FILENAME_LENGTH {
		return 0, fmt.Errorf("fileName length over range (max: %d): %d", MAX_FILENAME_LENGTH, len(fileName))
	}

	this.lock()
	defer this.unlock()

	fileIDBytes := []byte(fileId)
	_fileName := []byte(fileName)
	var cerr C.TOX_ERR_FILE_SEND

	fileNumber := uint32(
		C.tox_file_send(
			this.toxcore,                  // Tox *tox
			C.uint32_t(friendNumber),      // uint32_t friend_number
			C.uint32_t(kind),              // uint32_t kind
			C.uint64_t(fileSize),          // uint64_t file_size
			(*C.uint8_t)(&fileIDBytes[0]), // const uint8_t *file_id
			(*C.uint8_t)(&_fileName[0]),   // const uint8_t *filename
			C.size_t(len(fileName)),       // size_t filename_length
			&cerr, // TOX_ERR_FILE_SEND *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_SEND_OK:
		if fileNumber == math.MaxUint32 {
			return 0, fmt.Errorf("file send error")
		}

		return fileNumber, nil

	case C.TOX_ERR_FILE_SEND_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND:
		return 0, fmt.Errorf("friend_number passed did not designate a valid friend")

	case C.TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED:
		return 0, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FILE_SEND_NAME_TOO_LONG:
		return 0, fmt.Errorf("filename length exceeded MaxFilenameLength (%d) bytes: %d", MAX_FILENAME_LENGTH, len(fileName))

	case C.TOX_ERR_FILE_SEND_TOO_MANY:
		return 0, fmt.Errorf("too many ongoing transfers. allowed in 256 per friend per direction (sending and receiving)")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FileSendChunk(friendNumber uint32, fileNumber uint32, position uint64, data []byte) (bool, error) {
	if data == nil || len(data) == 0 {
		return false, toxerr("empty data")
	}

	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_FILE_SEND_CHUNK

	ok := bool(
		C.tox_file_send_chunk(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			C.uint32_t(fileNumber),   // uint32_t file_number
			C.uint64_t(position),     // uint64_t position
			(*C.uint8_t)(&data[0]),   // const uint8_t *data
			C.size_t(len(data)),      // size_t length
			&cerr,                    // TOX_ERR_FILE_SEND_CHUNK *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_SEND_CHUNK_OK:
		assert(ok, "tox_file_send_chunk() return 'false' on success")

		return ok, nil

	case C.TOX_ERR_FILE_SEND_CHUNK_NULL:
		return false, fmt.Errorf("length parameter was non-zero, but data was NULL")

	case C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND:
		return false, fmt.Errorf("friend_number passed did not designate a valid friend")

	case C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED:
		return false, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND:
		return false, fmt.Errorf("no file transfer with the given file number was found for the given friend")

	case C.TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING:
		return false, fmt.Errorf("file transfer was found but isn't in a transferring state")

	case C.TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH:
		return false, fmt.Errorf("attempted to send more or less data than requested")

	case C.TOX_ERR_FILE_SEND_CHUNK_SENDQ:
		return false, fmt.Errorf("packet queue is full")

	case C.TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION:
		return false, fmt.Errorf("position parameter was wrong")

	default:
		return false, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FileSeek(friendNumber uint32, fileNumber uint32, position uint64) (bool, error) {
	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_FILE_SEEK

	ok := bool(
		C.tox_file_seek(
			this.toxcore,             // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			C.uint32_t(fileNumber),   // uint32_t file_number
			C.uint64_t(position),     // uint64_t position
			&cerr,                    // TOX_ERR_FILE_SEEK *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_SEEK_OK:
		assert(ok, "tox_file_seek() return 'false' on success")

		return ok, nil

	case C.TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND:
		return false, fmt.Errorf("friend_number passed did not designate a valid friend")

	case C.TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED:
		return false, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FILE_SEEK_NOT_FOUND:
		return false, fmt.Errorf("file transfer with the given file number was found for the given friend")

	case C.TOX_ERR_FILE_SEEK_DENIED:
		return false, fmt.Errorf("file was not in a state where it could be seed")

	case C.TOX_ERR_FILE_SEEK_INVALID_POSITION:
		return false, fmt.Errorf("seek position was invalid")

	case C.TOX_ERR_FILE_SEEK_SENDQ:
		return false, fmt.Errorf("packet queue is full")

	default:
		return false, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) FileGetFileId(friendNumber uint32, fileNumber uint32) (string, error) {
	var cerr C.TOX_ERR_FILE_GET
	var fileId_b = make([]byte, C.TOX_FILE_ID_LENGTH)

	ok := bool(
		C.tox_file_get_file_id(
			this.toxcore,               // const Tox *tox
			C.uint32_t(fileNumber),     // uint32_t friend_number
			C.uint32_t(fileNumber),     // uint32_t file_number
			(*C.uint8_t)(&fileId_b[0]), // uint8_t *file_id
			&cerr, // TOX_ERR_FILE_GET *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_GET_OK:
		assert(ok, "tox_file_get_file_id() return 'false' on success")

		return strings.ToUpper(hex.EncodeToString(fileId_b)), nil

	case C.TOX_ERR_FILE_GET_NULL:
		return "", fmt.Errorf("One of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FILE_GET_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("friend_number passed did not designate a valid friend")

	case C.TOX_ERR_FILE_GET_NOT_FOUND:
		return "", fmt.Errorf("No file transfer with the given file number was found for the given friend")

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

func (this *Tox) IsConnected() int {
	r := C.tox_self_get_connection_status(this.toxcore)
	return int(r)
}

func (this *Tox) putcbevts(f func()) { this.cbevts = append(this.cbevts, f) }
