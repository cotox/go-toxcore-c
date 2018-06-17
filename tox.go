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
static inline __attribute__((__unused__)) void fixnousetox() {
}

*/
import "C"
import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	// "sync"
	"unsafe"

	deadlock "github.com/sasha-s/go-deadlock"
)

// "reflect"
// "runtime"

//////////
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
		pubkey_b := C.GoBytes(unsafe.Pointer(a0), C.int(PublicKeySize))
		pubkey := hex.EncodeToString(pubkey_b)
		pubkey = strings.ToUpper(pubkey)
		message_b := C.GoBytes(unsafe.Pointer(a1), C.int(a2))
		message := string(message_b)
		cbfn := *(*cb_friend_request_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, pubkey, message, ud) })
	}
}

// CallbackFriendRequest sets event handler which is triggered when a friend request is received.
func (this *Tox) CallbackFriendRequest(cbfn cb_friend_request_ftype, userData interface{}) {
	this.callbackFriendRequestAdd(cbfn, userData)
}

func (this *Tox) callbackFriendRequestAdd(cbfn cb_friend_request_ftype, userData interface{}) {
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

// CallbackFriendMessage sets event handler which is triggered when a message from a friend is received.
func (this *Tox) CallbackFriendMessage(cbfn cb_friend_message_ftype, userData interface{}) {
	this.callbackFriendMessageAdd(cbfn, userData)
}

func (this *Tox) callbackFriendMessageAdd(cbfn cb_friend_message_ftype, userData interface{}) {
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

// CallbackFriendName sets event handler which is triggered when a friend changes their name.
func (this *Tox) CallbackFriendName(cbfn cb_friend_name_ftype, userData interface{}) {
	this.callbackFriendNameAdd(cbfn, userData)
}

func (this *Tox) callbackFriendNameAdd(cbfn cb_friend_name_ftype, userData interface{}) {
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

// CallbackFriendStatusMessage sets event handler which is triggered when a friend changes their status message.
func (this *Tox) CallbackFriendStatusMessage(cbfn cb_friend_status_message_ftype, userData interface{}) {
	this.callbackFriendStatusMessageAdd(cbfn, userData)
}

func (this *Tox) callbackFriendStatusMessageAdd(cbfn cb_friend_status_message_ftype, userData interface{}) {
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

// CallbackFriendStatus sets event handler which is triggered when a friend changes their user status.
func (this *Tox) CallbackFriendStatus(cbfn cb_friend_status_ftype, userData interface{}) {
	this.callbackFriendStatusAdd(cbfn, userData)
}

func (this *Tox) callbackFriendStatusAdd(cbfn cb_friend_status_ftype, userData interface{}) {
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

// CallbackFriendConnectionStatus sets event handler which is triggered when a friend goes offline after having been online, or when a friend goes online.
//
// The handler will not triggered while adding friends. It is assumed that when adding friends, their connection status is initially offline.
func (this *Tox) CallbackFriendConnectionStatus(cbfn cb_friend_connection_status_ftype, userData interface{}) {
	this.callbackFriendConnectionStatusAdd(cbfn, userData)
}

func (this *Tox) callbackFriendConnectionStatusAdd(cbfn cb_friend_connection_status_ftype, userData interface{}) {
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

// CallbackFriendTyping sets event handler which is triggered when a friend starts or stops typing.
func (this *Tox) CallbackFriendTyping(cbfn cb_friend_typing_ftype, userData interface{}) {
	this.callbackFriendTypingAdd(cbfn, userData)
}

func (this *Tox) callbackFriendTypingAdd(cbfn cb_friend_typing_ftype, userData interface{}) {
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

// CallbackFriendReadReceipt sets event handler which is triggered when the friend receives the message sent with tox_friend_send_message with the corresponding message ID.
func (this *Tox) CallbackFriendReadReceipt(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
	this.callbackFriendReadReceiptAdd(cbfn, userData)
}

func (this *Tox) callbackFriendReadReceiptAdd(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
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
	this.callbackFriendLossyPacketAdd(cbfn, userData)
}

func (this *Tox) callbackFriendLossyPacketAdd(cbfn cb_friend_lossy_packet_ftype, userData interface{}) {
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
	this.callbackFriendLosslessPacketAdd(cbfn, userData)
}

func (this *Tox) callbackFriendLosslessPacketAdd(cbfn cb_friend_lossless_packet_ftype, userData interface{}) {
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

// CallbackSelfConnectionStatus sets event handler which is triggered whenever there is a change in the DHT connection state. When disconnected, a client may choose to call tox_bootstrap again, to reconnect to the DHT. Note that this state may frequently change for short amounts of time. Clients should therefore not immediately bootstrap on receiving a disconnect.
func (this *Tox) CallbackSelfConnectionStatus(cbfn cb_self_connection_status_ftype, userData interface{}) {
	this.callbackSelfConnectionStatusAdd(cbfn, userData)
}

func (this *Tox) callbackSelfConnectionStatusAdd(cbfn cb_self_connection_status_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_self_connection_statuss[cbfnp]; ok {
		return
	}
	this.cb_self_connection_statuss[cbfnp] = userData

	C.tox_callback_self_connection_status(this.toxcore, (*C.tox_self_connection_status_cb)(C.callbackSelfConnectionStatusWrapperForC))
}

// 包内部函数
//export callbackFileRecvControlWrapperForC
func callbackFileRecvControlWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	control C.TOX_FILE_CONTROL, userData unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_file_recv_controls {
		cbfn := *(*cb_file_recv_control_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(friendNumber), uint32(fileNumber), int(control), ud) })
	}
}

// CallbackFileRecvControl sets event handler which is triggered when a file control command is received from a friend.
func (this *Tox) CallbackFileRecvControl(cbfn cb_file_recv_control_ftype, userData interface{}) {
	this.callbackFileRecvControlAdd(cbfn, userData)
}

func (this *Tox) callbackFileRecvControlAdd(cbfn cb_file_recv_control_ftype, userData interface{}) {
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

// CallbackFileRecv sets event handler which is triggered when a file transfer request is received.
func (this *Tox) CallbackFileRecv(cbfn cb_file_recv_ftype, userData interface{}) {
	this.callbackFileRecvAdd(cbfn, userData)
}

func (this *Tox) callbackFileRecvAdd(cbfn cb_file_recv_ftype, userData interface{}) {
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

// CallbackFileRecvChunk sets event handler which is first triggered when a file transfer request is received, and subsequently when a chunk of file data for an accepted request was received.
func (this *Tox) CallbackFileRecvChunk(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
	this.callbackFileRecvChunkAdd(cbfn, userData)
}

func (this *Tox) callbackFileRecvChunkAdd(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
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

// CallbackFileChunkRequest sets event handler which is triggered when Core is ready to send more file data.
func (this *Tox) CallbackFileChunkRequest(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	this.callbackFileChunkRequestAdd(cbfn, userData)
}

func (this *Tox) callbackFileChunkRequestAdd(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_file_chunk_requests[cbfnp]; ok {
		return
	}
	this.cb_file_chunk_requests[cbfnp] = userData

	C.tox_callback_file_chunk_request(this.toxcore, (*C.tox_file_chunk_request_cb)(C.callbackFileChunkRequestWrapperForC))
}

// NewTox creates and initializes a new Tox instance with the options passed.
// If the opt is nil, the default options are used.
func NewTox(opt *ToxOptions) *Tox {
	var tox = new(Tox)
	if opt != nil {
		tox.opts = opt
	} else {
		tox.opts = NewToxOptions()
	}
	toxopts := tox.opts.toCToxOptions()
	defer C.tox_options_free(toxopts)

	var cerr C.TOX_ERR_NEW
	var toxcore = C.tox_new(toxopts, &cerr)
	tox.toxcore = toxcore
	if toxcore == nil {
		log.Println(toxerr(cerr))
		return nil
	}
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

	return tox
}

// Kill releases all resources associated with the Tox instance and disconnects from the network.
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
func (this *Tox) IterationInterval() int {
	this.lock()
	defer this.unlock()

	r := C.tox_iteration_interval(this.toxcore)
	return int(r)
}

/* The main loop that needs to be run in intervals of tox_iteration_interval() ms. */
// void tox_iterate(Tox *tox);
// compatable with legacy version
func (this *Tox) Iterate() {
	this.lock()
	C.tox_iterate(this.toxcore, nil)
	cbevts := this.cbevts
	this.cbevts = nil
	this.unlock()

	this.invokeCallbackEvents(cbevts)
}

// for toktok new method
func (this *Tox) Iterate2(userData interface{}) {
	this.lock()
	this.cb_iterate_data = userData
	C.tox_iterate(this.toxcore, nil)
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
	r := C.tox_get_savedata_size(this.toxcore)
	return int32(r)
}

func (this *Tox) GetSavedata() []byte {
	r := C.tox_get_savedata_size(this.toxcore)
	var savedata = make([]byte, int(r))

	C.tox_get_savedata(this.toxcore, (*C.uint8_t)(&savedata[0]))
	return savedata
}

// Bootstrap Sends a "get nodes" request to the given bootstrap node with IP, port, and public key to setup connections.
//
// This function will attempt to connect to the node using UDP. You must use this function even if Tox_Options.udp_enabled was set to false.
func (this *Tox) Bootstrap(addr string, port uint16, pubkey string) (bool, error) {
	this.lock()
	defer this.unlock()

	b_pubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return false, toxerr("Invalid pubkey")
	}

	var _addr = []byte(addr)
	var _port = C.uint16_t(port)
	var _cpubkey = (*C.uint8_t)(&b_pubkey[0])

	var cerr C.TOX_ERR_BOOTSTRAP
	r := C.tox_bootstrap(this.toxcore, (*C.char)(unsafe.Pointer(&_addr[0])), _port, _cpubkey, &cerr)
	if cerr > 0 {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

// SelfGetAddress returns the Tox friend address of the client.
func (this *Tox) SelfGetAddress() string {
	var addr [AddressSize]byte
	var caddr = (*C.uint8_t)(unsafe.Pointer(&addr[0]))
	C.tox_self_get_address(this.toxcore, caddr)

	return strings.ToUpper(hex.EncodeToString(addr[:]))
}

// SelfGetConnectionStatus returns whether we are connected to the DHT. The return value is equal to the last value received through the `self_connection_status` callback.
//
// @deprecated This getter is deprecated. Use the event and store the status in the client state.
//
// TODO: remove and handle the status inside go-toxcore-c, and provides the status as an attribute.
// TODO: wrap TOX_CONNECTION as a Go type, and change the return value suite for it.
func (this *Tox) SelfGetConnectionStatus() int {
	r := C.tox_self_get_connection_status(this.toxcore)
	return int(r)
}

// FriendAdd adds a friend to the friend list, send a friend request and returns the friend number.
//
// The length of message should between 1 and TOX_MAX_FRIEND_REQUEST_LENGTH.
//
// Friend numbers are unique identifiers used in all functions that operate on friends. Once added, a friend number is stable for the lifetime of the Tox object. After saving the state and reloading it, the friend numbers may not be the same as before. Deleting a friend creates a gap in the friend number set, which is filled by the next adding of a friend. Any pattern in friend numbers should not be relied on.
//
// NOTE: If more than INT32_MAX friends are added, this function causes undefined behaviour.
func (this *Tox) FriendAdd(friendId string, message string) (uint32, error) {
	this.lock()
	defer this.unlock()

	friendId_b, err := hex.DecodeString(friendId)
	friendId_p := (*C.uint8_t)(&friendId_b[0])
	if err != nil {
		log.Panic(err)
	}

	cmessage := []byte(message)

	var cerr C.TOX_ERR_FRIEND_ADD
	r := C.tox_friend_add(this.toxcore, friendId_p,
		(*C.uint8_t)(&cmessage[0]), C.size_t(len(message)), &cerr)
	if cerr > 0 {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

// FriendAddNorequest adds a friend without sending a friend request and returns the friend number.
//
// This function is used to add a friend in response to a friend request. If the client receives a friend request, it can be reasonably sure that the other client added this client as a friend, eliminating the need for a friend request.
//
// This function is also useful in a situation where both instances are controlled by the same entity, so that this entity can perform the mutual friend adding. In this case, there is no need for a friend request, either.
func (this *Tox) FriendAddNorequest(friendId string) (uint32, error) {
	this.lock()
	defer this.unlock()

	friendId_b, err := hex.DecodeString(friendId)
	if err != nil {
		return 0, err
	}
	friendId_p := (*C.uint8_t)(&friendId_b[0])

	var cerr C.TOX_ERR_FRIEND_ADD
	r := C.tox_friend_add_norequest(this.toxcore, friendId_p, &cerr)
	if cerr > 0 {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

// FriendByPublicKey returns the friend number associated with that Public Key.
func (this *Tox) FriendByPublicKey(pubkey string) (uint32, error) {
	pubkey_b, err := hex.DecodeString(pubkey)
	if err != nil {
		return 0, err
	}
	var pubkey_p = (*C.uint8_t)(&pubkey_b[0])

	var cerr C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	r := C.tox_friend_by_public_key(this.toxcore, pubkey_p, &cerr)
	if cerr != C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

// FriendGetPublicKey returns the Public Key associated with a given friend number.
func (this *Tox) FriendGetPublicKey(friendNumber uint32) (string, error) {
	var _fn = C.uint32_t(friendNumber)
	var pubkey_b = make([]byte, PublicKeySize)
	var pubkey_p = (*C.uint8_t)(&pubkey_b[0])

	var cerr C.TOX_ERR_FRIEND_GET_PUBLIC_KEY
	r := C.tox_friend_get_public_key(this.toxcore, _fn, pubkey_p, &cerr)
	if cerr > 0 || bool(r) == false {
		// TOFIX: cerr is undefined when r is false.
		return "", toxerr(cerr)
	}
	pubkey_h := hex.EncodeToString(pubkey_b)
	pubkey_h = strings.ToUpper(pubkey_h)
	return pubkey_h, nil
}

// FriendDelete removes friend from the friend list and returns true on success.
//
// This does not notify the friend of their deletion. After calling this function, this client will appear offline to the friend and no communication can occur between the two.
func (this *Tox) FriendDelete(friendNumber uint32) (bool, error) {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_DELETE
	r := C.tox_friend_delete(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

// FriendGetConnectionStatus Check whether a friend is currently connected to this client.
//
// The result of this function is equal to the last value received by the `friend_connection_status` callback.
//
// @deprecated This getter is deprecated. Use the event and store the status in the client state.
//
// TODO: remove this func and implement it in recommend.
func (this *Tox) FriendGetConnectionStatus(friendNumber uint32) (int, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	r := C.tox_friend_get_connection_status(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return int(r), toxerr(cerr)
	}
	return int(r), nil
}

// FriendExists checks if a friend with the given friend number exists and returns true if it does.
func (this *Tox) FriendExists(friendNumber uint32) bool {
	var _fn = C.uint32_t(friendNumber)

	r := C.tox_friend_exists(this.toxcore, _fn)
	return bool(r)
}

// FriendSendMessage send text chat message to an online friend and returns the message ID.
//
// This function creates a chat message packet and pushes it into the send queue.
//
// The message length may not exceed TOX_MAX_MESSAGE_LENGTH. Larger messages must be split by the client and sent as separate messages. Other clients can then reassemble the fragments. Messages may not be empty.
//
// The return value of this function is the message ID. If a read receipt is received, the triggered `friend_read_receipt` event will be passed this message ID.
//
// Message IDs are unique per friend. The first message ID is 0. Message IDs are incremented by 1 each time a message is sent. If UINT32_MAX messages were sent, the next message ID is 0.
func (this *Tox) FriendSendMessage(friendNumber uint32, message string) (uint32, error) {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _message = []byte(message)
	var _length = C.size_t(len(message))

	var mtype C.TOX_MESSAGE_TYPE = C.TOX_MESSAGE_TYPE_NORMAL
	var cerr C.TOX_ERR_FRIEND_SEND_MESSAGE
	r := C.tox_friend_send_message(this.toxcore, _fn, mtype, (*C.uint8_t)(&_message[0]), _length, &cerr)
	if cerr != C.TOX_ERR_FRIEND_SEND_MESSAGE_OK {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

func (this *Tox) FriendSendAction(friendNumber uint32, action string) (uint32, error) {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _action = []byte(action)
	var _length = C.size_t(len(action))

	var mtype C.TOX_MESSAGE_TYPE = C.TOX_MESSAGE_TYPE_ACTION
	var cerr C.TOX_ERR_FRIEND_SEND_MESSAGE
	r := C.tox_friend_send_message(this.toxcore, _fn, mtype, (*C.uint8_t)(&_action[0]), _length, &cerr)
	if cerr > 0 {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

// SelfSetName sets nickname for the Tox client.
//
// Nickname length cannot exceed TOX_MAX_NAME_LENGTH. If length is 0, the name parameter is ignored (it can be NULL), and the nickname is set back to empty.
//
// TODO: tox_self_set_name() returns boolean value indicate status of set.
func (this *Tox) SelfSetName(name string) error {
	this.lock()
	defer this.unlock()

	var _name = []byte(name)
	var _length = C.size_t(len(name))

	var cerr C.TOX_ERR_SET_INFO
	C.tox_self_set_name(this.toxcore, (*C.uint8_t)(&_name[0]), _length, &cerr)
	if cerr > 0 {
		return toxerr(cerr)
	}
	return nil
}

// SelfGetName returns the nickname set by SelfSetName.
// If no nickname was set before calling this function, the name is empty, and this function has no effect.
func (this *Tox) SelfGetName() string {
	// TODO: tox_self_get_name_size() could return 0 if the nickname is not set. line below wrong?
	nlen := C.tox_self_get_name_size(this.toxcore) // TODO: to replace by SelfGetNameSize()
	_name := make([]byte, nlen)

	C.tox_self_get_name(this.toxcore, (*C.uint8_t)(safeptr(_name)))
	return string(_name)
}

// SelfGetNameSize returns the length of the current nickname as passed to tox_self_set_name.
//
// If no nickname was set before calling this function, the name is empty, and this function returns 0.
//
// @see threading for concurrency implications.
func (this *Tox) SelfGetNameSize() int {
	r := C.tox_self_get_name_size(this.toxcore)
	return int(r)
}

// FriendGetName returns the name of the friend designated by the given friend number.
//
// The returned value is equal to the data received by the last `friend_name` callback.
func (this *Tox) FriendGetName(friendNumber uint32) (string, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	// TODO: tox_friend_get_name_size() should return an unspecified value. need to fix?
	nlen := C.tox_friend_get_name_size(this.toxcore, _fn, &cerr)
	_name := make([]byte, nlen)

	r := C.tox_friend_get_name(this.toxcore, _fn, (*C.uint8_t)(safeptr(_name)), &cerr)
	if !bool(r) {
		return "", toxerr(cerr)
	}
	return string(_name), nil
}

// FriendGetNameSize returns the length of the friend's name. If the friend number is invalid, the return value is unspecified.
func (this *Tox) FriendGetNameSize(friendNumber uint32) (int, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	r := C.tox_friend_get_name_size(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return int(r), toxerr(cerr)
	}
	return int(r), nil
}

// SelfSetStatusMessage sets client's status message.
//
// Status message length cannot exceed TOX_MAX_STATUS_MESSAGE_LENGTH. If length is 0, the status parameter is ignored (it can be NULL), and the user status is set back to empty.
func (this *Tox) SelfSetStatusMessage(status string) (bool, error) {
	this.lock()
	defer this.unlock()

	var _status = []byte(status)
	var _length = C.size_t(len(status))

	var cerr C.TOX_ERR_SET_INFO
	r := C.tox_self_set_status_message(this.toxcore, (*C.uint8_t)(&_status[0]), _length, &cerr)
	if cerr > 0 {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

// SelfSetStatus sets client's user status.
func (this *Tox) SelfSetStatus(status uint8) {
	var _status = C.TOX_USER_STATUS(status)
	C.tox_self_set_status(this.toxcore, _status)
}

// FriendGetStatusMessageSize returns the length of the friend's status message. If the friend number is invalid, the return value is SIZE_MAX.
func (this *Tox) FriendGetStatusMessageSize(friendNumber uint32) (int, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	r := C.tox_friend_get_status_message_size(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return int(r), toxerr(cerr)
	}
	return int(r), nil
}

// SelfGetStatusMessageSize returns the length of the current status message as passed to tox_self_set_status_message.
//
// If no status message was set before calling this function, the status is empty, and this function returns 0.
func (this *Tox) SelfGetStatusMessageSize() int {
	r := C.tox_self_get_status_message_size(this.toxcore)
	return int(r)
}

// FriendGetStatusMessage returns the status message of the friend designated by the given friend number.
func (this *Tox) FriendGetStatusMessage(friendNumber uint32) (string, error) {
	var _fn = C.uint32_t(friendNumber)
	var cerr C.TOX_ERR_FRIEND_QUERY
	len := C.tox_friend_get_status_message_size(this.toxcore, _fn, &cerr) // TODO: to replace by FriendGetStatusMessageSize
	if cerr > 0 {
		return "", toxerr(cerr)
	}

	_buf := make([]byte, len)

	cerr = 0
	r := C.tox_friend_get_status_message(this.toxcore, _fn, (*C.uint8_t)(safeptr(_buf)), &cerr)
	if !bool(r) || cerr > 0 {
		return "", toxerr(cerr)
	}
	return string(_buf[:]), nil
}

// SelfGetStatusMessage returns the status message set by tox_self_set_status_message.
//
// If no status message was set before calling this function, the status is empty, and this function has no effect.
func (this *Tox) SelfGetStatusMessage() (string, error) {
	nlen := C.tox_self_get_status_message_size(this.toxcore) // TODO: replace by SelfGetStatusMessageSize
	var _buf = make([]byte, nlen)

	C.tox_self_get_status_message(this.toxcore, (*C.uint8_t)(safeptr(_buf)))
	return string(_buf[:]), nil
}

// FriendGetStatus returns the friend's user status (away/busy/...). If the friend number is invalid, the return value is unspecified.
//
// The status returned is equal to the last status received through the `friend_status` callback.
//
// @deprecated This getter is deprecated. Use the event and store the status in the client state.
//
// TODO: remove this func
func (this *Tox) FriendGetStatus(friendNumber uint32) (int, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	r := C.tox_friend_get_status(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return int(r), toxerr(cerr)
	}
	return int(r), nil
}

// SelfGetStatus returns client's user status.
func (this *Tox) SelfGetStatus() int {
	r := C.tox_self_get_status(this.toxcore)
	return int(r)
}

// FriendGetLastOnline returns a unix-time timestamp of the last time the friend associated with a given friend number was seen online. This function will return UINT64_MAX on error.
//
// TODO: change return value in type time.Time
func (this *Tox) FriendGetLastOnline(friendNumber uint32) (uint64, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_GET_LAST_ONLINE
	r := C.tox_friend_get_last_online(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return uint64(r), toxerr(cerr)
	}
	return uint64(r), nil
}

// SelfSetTyping sets client's typing status for a friend and returns true on success.
//
// The client is responsible for turning it on or off.
func (this *Tox) SelfSetTyping(friendNumber uint32, typing bool) (bool, error) {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _typing = C._Bool(typing)

	var cerr C.TOX_ERR_SET_TYPING
	r := C.tox_self_set_typing(this.toxcore, _fn, _typing, &cerr)
	if cerr > 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

// FriendGetTyping checks whether a friend is currently typing a message and returns true if the friend is typing, or false if the friend is not typing, or the friend number was invalid. Inspect the error code to determine which case it is.
//
// @deprecated This getter is deprecated. Use the event and store the status in the client state.
//
// TODO: remove this func
func (this *Tox) FriendGetTyping(friendNumber uint32) (bool, error) {
	var _fn = C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	r := C.tox_friend_get_typing(this.toxcore, _fn, &cerr)
	if cerr > 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

// SelfGetFriendListSize returns the number of friends on the friend list.
//
// This function can be used to determine how much memory to allocate for tox_self_get_friend_list.
func (this *Tox) SelfGetFriendListSize() uint32 {
	r := C.tox_self_get_friend_list_size(this.toxcore)
	return uint32(r)
}

// SelfGetFriendList returns a list of valid friend numbers.
func (this *Tox) SelfGetFriendList() []uint32 {
	sz := C.tox_self_get_friend_list_size(this.toxcore)
	vec := make([]uint32, sz)
	if sz == 0 {
		return vec
	}
	vec_p := unsafe.Pointer(&vec[0])
	C.tox_self_get_friend_list(this.toxcore, (*C.uint32_t)(vec_p))
	return vec
}

// tox_callback_***

// SelfGetNospam returns the 4-byte nospam part of the address. This value is returned in host byte order.
func (this *Tox) SelfGetNospam() uint32 {
	r := C.tox_self_get_nospam(this.toxcore)
	return uint32(r)
}

// SelfSetNospam sets the 4-byte nospam part of the address. This value is expected in host byte order. I.e. 0x12345678 will form the bytes [12, 34, 56, 78] in the nospam part of the Tox friend address.
func (this *Tox) SelfSetNospam(nospam uint32) {
	this.lock()
	defer this.unlock()

	var _nospam = C.uint32_t(nospam)

	C.tox_self_set_nospam(this.toxcore, _nospam)
}

// SelfGetPublicKey returns the Tox Public Key (long term) from the Tox object.
func (this *Tox) SelfGetPublicKey() string {
	var _pubkey [PublicKeySize]byte

	C.tox_self_get_public_key(this.toxcore, (*C.uint8_t)(&_pubkey[0]))

	return strings.ToUpper(hex.EncodeToString(_pubkey[:]))
}

// SelfGetSecretKey returns the Tox Secret Key from the Tox object.
func (this *Tox) SelfGetSecretKey() string {
	var _seckey [SecretKeySize]byte

	C.tox_self_get_secret_key(this.toxcore, (*C.uint8_t)(&_seckey[0]))

	return strings.ToUpper(hex.EncodeToString(_seckey[:]))
}

// tox_lossy_***

// FriendSendLossyPacket sends a custom lossy packet to a friend.
//
// The first byte of data must be in the range 200-254. Maximum length of a custom packet is TOX_MAX_CUSTOM_PACKET_SIZE.
//
// Lossy packets behave like UDP packets, meaning they might never reach the other side or might arrive more than once (if someone is messing with the connection) or might arrive in the wrong order.
//
// Unless latency is an issue, it is recommended that you use lossless custom packets instead.
func (this *Tox) FriendSendLossyPacket(friendNumber uint32, data string) error {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _data = []byte(data)
	var _length = C.size_t(len(data))

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	r := C.tox_friend_send_lossy_packet(this.toxcore, _fn, (*C.uint8_t)(&_data[0]), _length, &cerr)
	if !r || cerr != C.TOX_ERR_FRIEND_CUSTOM_PACKET_OK {
		return toxerr(cerr)
	}
	return nil
}

// FriendSendLosslessPacket sends a custom lossless packet to a friend and returns true on success.
//
// The first byte of data must be in the range 160-191. Maximum length of a custom packet is TOX_MAX_CUSTOM_PACKET_SIZE.
//
// Lossless packet behaviour is comparable to TCP (reliability, arrive in order) but with packets instead of a stream.
func (this *Tox) FriendSendLosslessPacket(friendNumber uint32, data string) error {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _data = []byte(data)
	var _length = C.size_t(len(data))

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	r := C.tox_friend_send_lossless_packet(this.toxcore, _fn, (*C.uint8_t)(&_data[0]), _length, &cerr)
	if !r || cerr != C.TOX_ERR_FRIEND_CUSTOM_PACKET_OK {
		return toxerr(cerr)
	}
	return nil
}

// tox_callback_avatar_**

// Hash generates a cryptographic hash of the given data.
//
// This function may be used by clients for any purpose, but is provided primarily for validating cached avatars. This use is highly recommended to avoid unnecessary avatar updates.
//
// If hash is NULL or data is NULL while length is not 0 the function returns false, otherwise it returns true.
//
// This function is a wrapper to internal message-digest functions.
func (this *Tox) Hash(data string, datalen uint32) (string, bool, error) {
	_data := []byte(data)
	_hash := make([]byte, C.TOX_HASH_LENGTH)
	var _datalen = C.size_t(datalen)

	r := C.tox_hash((*C.uint8_t)(&_hash[0]), (*C.uint8_t)(&_data[0]), _datalen)
	return string(_hash), bool(r), nil
}

// tox_callback_file_***

// FileControl sends a file control command to a friend for a given file transfer and returns true on success.
func (this *Tox) FileControl(friendNumber uint32, fileNumber uint32, control int) (bool, error) {
	var cerr C.TOX_ERR_FILE_CONTROL
	r := C.tox_file_control(this.toxcore, C.uint32_t(friendNumber), C.uint32_t(fileNumber),
		C.TOX_FILE_CONTROL(control), &cerr)
	if cerr > 0 {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

// FileSend sends a file transmission request and returns a file number used as an identifier in subsequent callbacks. This number is per friend. File numbers are reused after a transfer terminates. On failure, this function returns UINT32_MAX. Any pattern in file numbers should not be relied on.
//
// Maximum filename length is TOX_MAX_FILENAME_LENGTH bytes. The filename should generally just be a file name, not a path with directory names.
//
// If a non-UINT64_MAX file size is provided, it can be used by both sides to determine the sending progress. File size can be set to UINT64_MAX for streaming data of unknown size.
//
// File transmission occurs in chunks, which are requested through the `file_chunk_request` event.
//
// When a friend goes offline, all file transfers associated with the friend are purged from core.
//
// If the file contents change during a transfer, the behaviour is unspecified in general. What will actually happen depends on the mode in which the file was modified and how the client determines the file size.
//
// - If the file size was increased
//   - and sending mode was streaming (file_size = UINT64_MAX), the behaviour will be as expected.
//   - and sending mode was file (file_size != UINT64_MAX), the file_chunk_request callback will receive length = 0 when Core thinks the file transfer has finished. If the client remembers the file size as it was when sending the request, it will terminate the transfer normally. If the client re-reads the size, it will think the friend cancelled the transfer.
// - If the file size was decreased
//   - and sending mode was streaming, the behaviour is as expected.
//   - and sending mode was file, the callback will return 0 at the new (earlier) end-of-file, signalling to the friend that the transfer was cancelled.
// - If the file contents were modified
//   - at a position before the current read, the two files (local and remote) will differ after the transfer terminates.
//   - at a position after the current read, the file transfer will succeed as expected.
//   - In either case, both sides will regard the transfer as complete and successful.
func (this *Tox) FileSend(friendNumber uint32, kind uint32, fileSize uint64, fileId string, fileName string) (uint32, error) {
	this.lock()
	defer this.unlock()

	if len(fileId) != FileIDLength*2 {
	}

	_fileName := []byte(fileName)

	var cerr C.TOX_ERR_FILE_SEND
	r := C.tox_file_send(this.toxcore, C.uint32_t(friendNumber), C.uint32_t(kind), C.uint64_t(fileSize),
		nil, (*C.uint8_t)(&_fileName[0]), C.size_t(len(fileName)), &cerr)
	if cerr > 0 {
		return uint32(r), toxerr(cerr)
	}
	return uint32(r), nil
}

// FileSendChunk sends a chunk of file data to a friend and returns true on success.
//
// This function is called in response to the `file_chunk_request` callback. The length parameter should be equal to the one received though the callback. If it is zero, the transfer is assumed complete. For files with known size, Core will know that the transfer is complete after the last byte has been received, so it is not necessary (though not harmful) to send a zero-length chunk to terminate. For streams, core will know that the transfer is finished if a chunk with length less than the length requested in the callback is sent.
func (this *Tox) FileSendChunk(friendNumber uint32, fileNumber uint32, position uint64, data []byte) (bool, error) {
	this.lock()
	defer this.unlock()

	if data == nil || len(data) == 0 {
		return false, toxerr("empty data")
	}
	var cerr C.TOX_ERR_FILE_SEND_CHUNK
	r := C.tox_file_send_chunk(this.toxcore, C.uint32_t(friendNumber), C.uint32_t(fileNumber),
		C.uint64_t(position), (*C.uint8_t)(&data[0]), C.size_t(len(data)), &cerr)
	if cerr > 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

// FileSeek sends a file seek control command to a friend for a given file transfer.
//
// This function can only be called to resume a file transfer right before TOX_FILE_CONTROL_RESUME is sent.
func (this *Tox) FileSeek(friendNumber uint32, fileNumber uint32, position uint64) (bool, error) {
	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_FILE_SEEK
	r := C.tox_file_seek(this.toxcore, C.uint32_t(friendNumber), C.uint32_t(fileNumber),
		C.uint64_t(position), &cerr)
	if cerr > 0 {
		return false, toxerr(cerr)
	}
	return bool(r), nil
}

// FileGetFileId copies the file id associated to the file transfer and returns true on success.
func (this *Tox) FileGetFileId(friendNumber uint32, fileNumber uint32) (string, error) {
	var cerr C.TOX_ERR_FILE_GET
	var fileId_b = make([]byte, C.TOX_FILE_ID_LENGTH)

	r := C.tox_file_get_file_id(this.toxcore, C.uint32_t(fileNumber), C.uint32_t(fileNumber),
		(*C.uint8_t)(&fileId_b[0]), &cerr)
	if cerr > 0 || bool(r) == false {
		return "", toxerr(cerr)
	}

	var fileId_h = strings.ToUpper(hex.EncodeToString(fileId_b))
	return fileId_h, nil
}

// boostrap, see upper

// AddTcpRelay adds additional host:port pair as TCP relay and returns true on success.
//
// This function can be used to initiate TCP connections to different ports on the same bootstrap node, or to add TCP relays without using them as bootstrap nodes.
func (this *Tox) AddTcpRelay(addr string, port uint16, pubkey string) (bool, error) {
	this.lock()
	defer this.unlock()

	var _addr = C.CString(addr)
	defer C.free(unsafe.Pointer(_addr))
	var _port = C.uint16_t(port)
	b_pubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		log.Panic(err)
	}
	if strings.ToUpper(hex.EncodeToString(b_pubkey)) != pubkey {
		log.Panic("wtf, hex enc/dec err")
	}
	var _pubkey = (*C.uint8_t)(&b_pubkey[0])

	var cerr C.TOX_ERR_BOOTSTRAP
	r := C.tox_add_tcp_relay(this.toxcore, _addr, _port, _pubkey, &cerr)
	if cerr > 0 {
		return bool(r), toxerr(cerr)
	}
	return bool(r), nil
}

// IsConnected Return whether we are connected to the DHT. The return value is equal to the last value received through the `self_connection_status` callback.
//
// @deprecated This getter is deprecated. Use the event and store the status in the client state.
//
// TODO: remove this func
func (this *Tox) IsConnected() int {
	r := C.tox_self_get_connection_status(this.toxcore)
	return int(r)
}

func (this *Tox) putcbevts(f func()) { this.cbevts = append(this.cbevts, f) }

////////////
/*
原则说明：
所有需要public_key的地方，在go空间内是实际串的16进制字符串表示。

*/

////////////////////
func KeepPkg() {
}

func _dirty_init() {
	fmt.Println("ddddddddd")
}
