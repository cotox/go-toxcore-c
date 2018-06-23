package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

// typedef void tox_friend_request_cb(Tox *tox, const uint8_t *public_key, const uint8_t *message, size_t length, void *user_data);
void callbackFriendRequestWrapperForC(Tox *tox, uint8_t *public_key, uint8_t *message, uint16_t length, void *user_data);

// typedef void tox_friend_message_cb(Tox *tox, uint32_t friend_number, TOX_MESSAGE_TYPE type, const uint8_t *message, size_t length, void *user_data);
void callbackFriendMessageWrapperForC(Tox *tox, uint32_t friend_number, int message_type, uint8_t *message, uint32_t length, void *user_data);

// typedef void tox_friend_name_cb(Tox *tox, uint32_t friend_number, const uint8_t *name, size_t length, void *user_data);
void callbackFriendNameWrapperForC(Tox *tox, uint32_t friend_number, uint8_t *name, uint32_t length, void *user_data);

// typedef void tox_friend_status_message_cb(Tox *tox, uint32_t friend_number, const uint8_t *message, size_t length, void *user_data);
void callbackFriendStatusMessageWrapperForC(Tox *tox, uint32_t friend_number, uint8_t *message, uint32_t length, void *user_data);

// typedef void tox_friend_status_cb(Tox *tox, uint32_t friend_number, TOX_USER_STATUS status, void *user_data);
void callbackFriendStatusWrapperForC(Tox *tox, uint32_t friend_number, int user_status, void *user_data);

// typedef void tox_friend_connection_status_cb(Tox *tox, uint32_t friend_number, TOX_CONNECTION connection_status, void *user_data);
void callbackFriendConnectionStatusWrapperForC(Tox *tox, uint32_t friend_number, int connection_status, void *user_data);

// typedef void tox_friend_typing_cb(Tox *tox, uint32_t friend_number, bool is_typing, void *user_data);
void callbackFriendTypingWrapperForC(Tox *tox, uint32_t friend_number, uint8_t is_typing, void *user_data);

// typedef void tox_friend_read_receipt_cb(Tox *tox, uint32_t friend_number, uint32_t message_id, void *user_data);
void callbackFriendReadReceiptWrapperForC(Tox *tox, uint32_t friend_number, uint32_t message_id, void *user_data);

// typedef void tox_friend_lossy_packet_cb(Tox *tox, uint32_t friend_number, const uint8_t *data, size_t length, void *user_data);
void callbackFriendLossyPacketWrapperForC(Tox *tox, uint32_t friend_number, uint8_t *data, size_t length, void *user_data);

// typedef void tox_friend_lossless_packet_cb(Tox *tox, uint32_t friend_number, const uint8_t *data, size_t length, void *user_data);
void callbackFriendLosslessPacketWrapperForC(Tox *tox, uint32_t friend_number, uint8_t *data, size_t length, void *user_data);

// typedef void tox_self_connection_status_cb(Tox *tox, TOX_CONNECTION connection_status, void *user_data);
void callbackSelfConnectionStatusWrapperForC(Tox *tox, int connection_status, void *user_data);

// typedef void tox_file_recv_control_cb(Tox *tox, uint32_t friend_number, uint32_t file_number, TOX_FILE_CONTROL control, void *user_data);
void callbackFileRecvControlWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, TOX_FILE_CONTROL control, void *user_data);

// typedef void tox_file_recv_cb(Tox *tox, uint32_t friend_number, uint32_t file_number, uint32_t kind, uint64_t file_size, const uint8_t *filename, size_t filename_length, void *user_data);
void callbackFileRecvWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint32_t kind, uint64_t file_size, uint8_t *filename, size_t filename_length, void *user_data);

// typedef void tox_file_recv_chunk_cb(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position, const uint8_t *data, size_t length, void *user_data);
void callbackFileRecvChunkWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position, uint8_t *data, size_t length, void *user_data);

// typedef void tox_file_chunk_request_cb(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position, size_t length, void *user_data);
void callbackFileChunkRequestWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position, size_t length, void *user_data);


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
	// "sync"
	"unsafe"

	deadlock "github.com/sasha-s/go-deadlock"
)

// "reflect"
// "runtime"

//////////
// friend callback type
type cb_friend_request_ftype func(t *Tox, pubkey string, message string, userData interface{})
type cb_friend_message_ftype func(t *Tox, friendNumber uint32, message string, userData interface{})
type cb_friend_name_ftype func(t *Tox, friendNumber uint32, newName string, userData interface{})
type cb_friend_status_message_ftype func(t *Tox, friendNumber uint32, newStatus string, userData interface{})
type cb_friend_status_ftype func(t *Tox, friendNumber uint32, status int, userData interface{})
type cb_friend_connection_status_ftype func(t *Tox, friendNumber uint32, status int, userData interface{})
type cb_friend_typing_ftype func(t *Tox, friendNumber uint32, isTyping uint8, userData interface{})
type cb_friend_read_receipt_ftype func(t *Tox, friendNumber uint32, receipt uint32, userData interface{})
type cb_friend_lossy_packet_ftype func(t *Tox, friendNumber uint32, data string, userData interface{})
type cb_friend_lossless_packet_ftype func(t *Tox, friendNumber uint32, data string, userData interface{})

// self callback type
type cb_self_connection_status_ftype func(t *Tox, status int, userData interface{})

// file callback type
type cb_file_recv_control_ftype func(t *Tox, friendNumber uint32, fileNumber uint32,
	control int, userData interface{})
type cb_file_recv_ftype func(t *Tox, friendNumber uint32, fileNumber uint32, kind uint32, fileSize uint64,
	fileName string, userData interface{})
type cb_file_recv_chunk_ftype func(t *Tox, friendNumber uint32, fileNumber uint32, position uint64,
	data []byte, userData interface{})
type cb_file_chunk_request_ftype func(t *Tox, friend_number uint32, file_number uint32, position uint64,
	length int, user_data interface{})

// Tox
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
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_requests {
		pubkey_b := C.GoBytes(unsafe.Pointer(a0), C.int(PublicKeySize))
		pubkey := hex.EncodeToString(pubkey_b)
		pubkey = strings.ToUpper(pubkey)
		message_b := C.GoBytes(unsafe.Pointer(a1), C.int(a2))
		message := string(message_b)
		cbfn := *(*cb_friend_request_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, pubkey, message, ud) })
	}
}

func (t *Tox) CallbackFriendRequest(cbfn cb_friend_request_ftype, userData interface{}) {
	t.CallbackFriendRequestAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendRequestAdd(cbfn cb_friend_request_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_requests[cbfnp]; ok {
		return
	}
	t.cb_friend_requests[cbfnp] = userData

	C.tox_callback_friend_request(t.toxcore, (*C.tox_friend_request_cb)(C.callbackFriendRequestWrapperForC))
}

//export callbackFriendMessageWrapperForC
func callbackFriendMessageWrapperForC(m *C.Tox, a0 C.uint32_t, mtype C.int,
	a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_messages {
		message_ := C.GoStringN((*C.char)(unsafe.Pointer(a1)), (C.int)(a2))
		cbfn := *(*cb_friend_message_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), message_, ud) })
	}
}

func (t *Tox) CallbackFriendMessage(cbfn cb_friend_message_ftype, userData interface{}) {
	t.CallbackFriendMessageAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendMessageAdd(cbfn cb_friend_message_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_messages[cbfnp]; ok {
		return
	}
	t.cb_friend_messages[cbfnp] = userData

	C.tox_callback_friend_message(t.toxcore, (*C.tox_friend_message_cb)(C.callbackFriendMessageWrapperForC))
}

//export callbackFriendNameWrapperForC
func callbackFriendNameWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_names {
		name := C.GoStringN((*C.char)((unsafe.Pointer)(a1)), C.int(a2))
		cbfn := *(*cb_friend_name_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), name, ud) })
	}
}

func (t *Tox) CallbackFriendName(cbfn cb_friend_name_ftype, userData interface{}) {
	t.CallbackFriendNameAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendNameAdd(cbfn cb_friend_name_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_names[cbfnp]; ok {
		return
	}
	t.cb_friend_names[cbfnp] = userData

	C.tox_callback_friend_name(t.toxcore, (*C.tox_friend_name_cb)(C.callbackFriendNameWrapperForC))
}

//export callbackFriendStatusMessageWrapperForC
func callbackFriendStatusMessageWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, a3 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_status_messages {
		statusText := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(a2))
		cbfn := *(*cb_friend_status_message_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), statusText, ud) })
	}
}

func (t *Tox) CallbackFriendStatusMessage(cbfn cb_friend_status_message_ftype, userData interface{}) {
	t.CallbackFriendStatusMessageAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendStatusMessageAdd(cbfn cb_friend_status_message_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_status_messages[cbfnp]; ok {
		return
	}
	t.cb_friend_status_messages[cbfnp] = userData

	C.tox_callback_friend_status_message(t.toxcore, (*C.tox_friend_status_message_cb)(C.callbackFriendStatusMessageWrapperForC))
}

//export callbackFriendStatusWrapperForC
func callbackFriendStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.int, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_statuss {
		cbfn := *(*cb_friend_status_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), int(a1), ud) })
	}
}

func (t *Tox) CallbackFriendStatus(cbfn cb_friend_status_ftype, userData interface{}) {
	t.CallbackFriendStatusAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendStatusAdd(cbfn cb_friend_status_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_statuss[cbfnp]; ok {
		return
	}
	t.cb_friend_statuss[cbfnp] = userData

	C.tox_callback_friend_status(t.toxcore, (*C.tox_friend_status_cb)(C.callbackFriendStatusWrapperForC))
}

//export callbackFriendConnectionStatusWrapperForC
func callbackFriendConnectionStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.int, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_connection_statuss {
		cbfn := *(*cb_friend_connection_status_ftype)((unsafe.Pointer)(cbfni))
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), int(a1), ud) })
	}
}

func (t *Tox) CallbackFriendConnectionStatus(cbfn cb_friend_connection_status_ftype, userData interface{}) {
	t.CallbackFriendConnectionStatusAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendConnectionStatusAdd(cbfn cb_friend_connection_status_ftype, userData interface{}) {
	cbfnp := unsafe.Pointer(&cbfn)
	if _, ok := t.cb_friend_connection_statuss[cbfnp]; ok {
		return
	}
	t.cb_friend_connection_statuss[cbfnp] = userData

	C.tox_callback_friend_connection_status(t.toxcore, (*C.tox_friend_connection_status_cb)(C.callbackFriendConnectionStatusWrapperForC))
}

//export callbackFriendTypingWrapperForC
func callbackFriendTypingWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint8_t, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_typings {
		cbfn := *(*cb_friend_typing_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), uint8(a1), ud) })
	}
}

func (t *Tox) CallbackFriendTyping(cbfn cb_friend_typing_ftype, userData interface{}) {
	t.CallbackFriendTypingAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendTypingAdd(cbfn cb_friend_typing_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_typings[cbfnp]; ok {
		return
	}
	t.cb_friend_typings[cbfnp] = userData

	C.tox_callback_friend_typing(t.toxcore, (*C.tox_friend_typing_cb)(C.callbackFriendTypingWrapperForC))
}

//export callbackFriendReadReceiptWrapperForC
func callbackFriendReadReceiptWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_read_receipts {
		cbfn := *(*cb_friend_read_receipt_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), uint32(a1), ud) })
	}
}

func (t *Tox) CallbackFriendReadReceipt(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
	t.CallbackFriendReadReceiptAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendReadReceiptAdd(cbfn cb_friend_read_receipt_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_read_receipts[cbfnp]; ok {
		return
	}
	t.cb_friend_read_receipts[cbfnp] = userData

	C.tox_callback_friend_read_receipt(t.toxcore, (*C.tox_friend_read_receipt_cb)(C.callbackFriendReadReceiptWrapperForC))
}

//export callbackFriendLossyPacketWrapperForC
func callbackFriendLossyPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_lossy_packets {
		cbfn := *(*cb_friend_lossy_packet_ftype)(cbfni)
		msg := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(len))
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), msg, ud) })
	}
}

func (t *Tox) CallbackFriendLossyPacket(cbfn cb_friend_lossy_packet_ftype, userData interface{}) {
	t.CallbackFriendLossyPacketAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendLossyPacketAdd(cbfn cb_friend_lossy_packet_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_lossy_packets[cbfnp]; ok {
		return
	}
	t.cb_friend_lossy_packets[cbfnp] = userData

	C.tox_callback_friend_lossy_packet(t.toxcore, (*C.tox_friend_lossy_packet_cb)(C.callbackFriendLossyPacketWrapperForC))
}

//export callbackFriendLosslessPacketWrapperForC
func callbackFriendLosslessPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_friend_lossless_packets {
		cbfn := *(*cb_friend_lossless_packet_ftype)(cbfni)
		msg := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(len))
		cTox.putcbevts(func() { cbfn(cTox, uint32(a0), msg, ud) })
	}
}

func (t *Tox) CallbackFriendLosslessPacket(cbfn cb_friend_lossless_packet_ftype, userData interface{}) {
	t.CallbackFriendLosslessPacketAdd(cbfn, userData)
}
func (t *Tox) CallbackFriendLosslessPacketAdd(cbfn cb_friend_lossless_packet_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_friend_lossless_packets[cbfnp]; ok {
		return
	}
	t.cb_friend_lossless_packets[cbfnp] = userData

	C.tox_callback_friend_lossless_packet(t.toxcore, (*C.tox_friend_lossless_packet_cb)(C.callbackFriendLosslessPacketWrapperForC))
}

//export callbackSelfConnectionStatusWrapperForC
func callbackSelfConnectionStatusWrapperForC(m *C.Tox, status C.int, a2 unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_self_connection_statuss {
		cbfn := *(*cb_self_connection_status_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, int(status), ud) })
	}
}

func (t *Tox) CallbackSelfConnectionStatus(cbfn cb_self_connection_status_ftype, userData interface{}) {
	t.CallbackSelfConnectionStatusAdd(cbfn, userData)
}
func (t *Tox) CallbackSelfConnectionStatusAdd(cbfn cb_self_connection_status_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_self_connection_statuss[cbfnp]; ok {
		return
	}
	t.cb_self_connection_statuss[cbfnp] = userData

	C.tox_callback_self_connection_status(t.toxcore, (*C.tox_self_connection_status_cb)(C.callbackSelfConnectionStatusWrapperForC))
}

// 包内部函数
//export callbackFileRecvControlWrapperForC
func callbackFileRecvControlWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	control C.TOX_FILE_CONTROL, userData unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_file_recv_controls {
		cbfn := *(*cb_file_recv_control_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(friendNumber), uint32(fileNumber), int(control), ud) })
	}
}

func (t *Tox) CallbackFileRecvControl(cbfn cb_file_recv_control_ftype, userData interface{}) {
	t.CallbackFileRecvControlAdd(cbfn, userData)
}
func (t *Tox) CallbackFileRecvControlAdd(cbfn cb_file_recv_control_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_file_recv_controls[cbfnp]; ok {
		return
	}
	t.cb_file_recv_controls[cbfnp] = userData

	C.tox_callback_file_recv_control(t.toxcore, (*C.tox_file_recv_control_cb)(C.callbackFileRecvControlWrapperForC))
}

//export callbackFileRecvWrapperForC
func callbackFileRecvWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, kind C.uint32_t,
	fileSize C.uint64_t, fileName *C.uint8_t, fileNameLength C.size_t, userData unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_file_recvs {
		cbfn := *(*cb_file_recv_ftype)(cbfni)
		fileName_ := C.GoStringN((*C.char)(unsafe.Pointer(fileName)), C.int(fileNameLength))
		cTox.putcbevts(func() {
			cbfn(cTox, uint32(friendNumber), uint32(fileNumber), uint32(kind),
				uint64(fileSize), fileName_, ud)
		})
	}
}

func (t *Tox) CallbackFileRecv(cbfn cb_file_recv_ftype, userData interface{}) {
	t.CallbackFileRecvAdd(cbfn, userData)
}
func (t *Tox) CallbackFileRecvAdd(cbfn cb_file_recv_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_file_recvs[cbfnp]; ok {
		return
	}
	t.cb_file_recvs[cbfnp] = userData

	C.tox_callback_file_recv(t.toxcore, (*C.tox_file_recv_cb)(C.callbackFileRecvWrapperForC))
}

//export callbackFileRecvChunkWrapperForC
func callbackFileRecvChunkWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	position C.uint64_t, data *C.uint8_t, length C.size_t, userData unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_file_recv_chunks {
		cbfn := *(*cb_file_recv_chunk_ftype)(cbfni)
		data_ := C.GoBytes((unsafe.Pointer)(data), C.int(length))
		cTox.putcbevts(func() { cbfn(cTox, uint32(friendNumber), uint32(fileNumber), uint64(position), data_, ud) })
	}
}

func (t *Tox) CallbackFileRecvChunk(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
	t.CallbackFileRecvChunkAdd(cbfn, userData)
}
func (t *Tox) CallbackFileRecvChunkAdd(cbfn cb_file_recv_chunk_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_file_recv_chunks[cbfnp]; ok {
		return
	}
	t.cb_file_recv_chunks[cbfnp] = userData

	C.tox_callback_file_recv_chunk(t.toxcore, (*C.tox_file_recv_chunk_cb)(C.callbackFileRecvChunkWrapperForC))
}

//export callbackFileChunkRequestWrapperForC
func callbackFileChunkRequestWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t,
	position C.uint64_t, length C.size_t, userData unsafe.Pointer) {
	var cTox = cbUserDatas.get(m)
	for cbfni, ud := range cTox.cb_file_chunk_requests {
		cbfn := *(*cb_file_chunk_request_ftype)(cbfni)
		cTox.putcbevts(func() { cbfn(cTox, uint32(friendNumber), uint32(fileNumber), uint64(position), int(length), ud) })
	}
}

func (t *Tox) CallbackFileChunkRequest(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	t.CallbackFileChunkRequestAdd(cbfn, userData)
}
func (t *Tox) CallbackFileChunkRequestAdd(cbfn cb_file_chunk_request_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := t.cb_file_chunk_requests[cbfnp]; ok {
		return
	}
	t.cb_file_chunk_requests[cbfnp] = userData

	C.tox_callback_file_chunk_request(t.toxcore, (*C.tox_file_chunk_request_cb)(C.callbackFileChunkRequestWrapperForC))
}

// NewTox returns a new Tox instance initialized with specific ToxOptions. The
// default ToxOptions will be use if argument opts is nil.
//
// TODO: rename "NewTox" => "New"
func NewTox(opts *ToxOptions) (*Tox, error) {
	var tox = new(Tox)
	if opts != nil {
		tox.opts = opts
	} else {
		tox.opts = NewToxOptions()
	}

	cToxOpts := tox.opts.toCToxOptions()
	defer C.tox_options_free(cToxOpts)

	var cerr C.TOX_ERR_NEW

	var cTox = C.tox_new(cToxOpts, &cerr)

	switch cerr {
	case C.TOX_ERR_NEW_OK:
		assert(cTox != nil, "cTox != nil")

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

	tox.toxcore = cTox
	cbUserDatas.set(cTox, tox)

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

// Kill releases all resources associated with the Tox instance and disconnects
// from the network.
func (t *Tox) Kill() {
	if t == nil || t.toxcore == nil {
		return
	}

	t.lock()
	defer t.unlock()

	cbUserDatas.del(t.toxcore)
	C.tox_kill(t.toxcore)
	t.toxcore = nil
}

// IterationInterval returns the time in milliseconds before Iterate should be
// called again for optimal performance.
func (t *Tox) IterationInterval() time.Duration {
	t.lock()
	defer t.unlock()

	return time.Duration(C.tox_iteration_interval(t.toxcore))
}

// Iterate is the main loop that needs to be run in intervals of
// IterationInterval milliseconds.
//
// void tox_iterate(Tox *tox);
// compatable with legacy version
//
// TODO: Remove this function. Use Iterate2 instead.
func (t *Tox) Iterate() {
	t.lock()

	C.tox_iterate(
		t.toxcore, // Tox *tox
		nil,       // void *user_data
	)
	cbevts := t.cbevts
	t.cbevts = nil

	t.unlock()

	t.invokeCallbackEvents(cbevts)
}

// Iterate2 is the main loop that needs to be run in intervals of
// IterationInterval milliseconds.
//
// NOTE: for toktok new method
func (t *Tox) Iterate2(userData interface{}) {
	t.lock()

	t.cb_iterate_data = userData
	C.tox_iterate(
		t.toxcore, // Tox *tox
		nil,       // void *user_data
	)
	t.cb_iterate_data = nil
	cbevts := t.cbevts
	t.cbevts = nil

	t.unlock()

	t.invokeCallbackEvents(cbevts)
}

func (t *Tox) invokeCallbackEvents(cbevts []func()) {
	for _, cbfn := range cbevts {
		cbfn()
	}
}

func (t *Tox) lock() {
	if t.opts.ThreadSafe {
		t.mu.Lock()
	}
}
func (t *Tox) unlock() {
	if t.opts.ThreadSafe {
		t.mu.Unlock()
	}
}

// getSavedataSize returns the number of bytes required to store the tox
// instance with GetSavedata. This function cannot fail. The result is always
// greater than 0.
func (t *Tox) getSavedataSize() int32 {
	return int32(C.tox_get_savedata_size(t.toxcore))
}

// GetSavedata returns a copy of all information associated with the tox
// instance.
func (t *Tox) GetSavedata() []byte {
	size := t.getSavedataSize()
	var savedata = make([]byte, int(size))

	C.tox_get_savedata(
		t.toxcore,                  // const Tox *tox
		(*C.uint8_t)(&savedata[0]), // uint8_t *savedata
	)

	return savedata
}

// Bootstrap sends a "get nodes" request to the given bootstrap node with IP,
// port, and public key to setup connections.
//
// This function will attempt to connect to the node using UDP. You must use
// this function even if ToxOptions.UDPEnabled is false.
func (t *Tox) Bootstrap(addr string, port uint16, pubKey string) error {
	t.lock()
	defer t.unlock()

	if len(pubKey) != PublicKeySize {
		return fmt.Errorf("invalid pubKey")
	}
	pubkeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return fmt.Errorf("invalid pubKey")
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

// AddTCPRelay adds additional host:port pair as TCP relay.
//
// This function can be used to initiate TCP connections to different ports on
// the same bootstrap node, or to add TCP relays without using them as
// bootstrap nodes.
func (t *Tox) AddTCPRelay(addr string, port uint16, pubKey string) error {
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

// SelfGetAddress returns a copy of the Tox friend address of the client.
//
// TODO: rename "SelfGetAddress" => "GetSelfAddress"?
func (t *Tox) SelfGetAddress() string {
	var addr [AddressSize]byte

	C.tox_self_get_address(
		t.toxcore, // const Tox *tox
		(*C.uint8_t)(unsafe.Pointer(&addr[0])), // uint8_t *address
	)

	return strings.ToUpper(hex.EncodeToString(addr[:]))
}

// SelfGetConnectionStatus returns whether we are connected to the DHT.
//
// The return value is equal to the last value received through the
// OnSelfConnectionStatusChanged event.
//
// TODO: @deprecated This getter is deprecated. Use the event and store the
//       status in the client state.
// TODO: rename "SelfGetConnectionStatus" => "GetSelfConnectionStatus"?
func (t *Tox) SelfGetConnectionStatus() ConnectionType {
	return ConnectionType(C.tox_self_get_connection_status(t.toxcore))
}

// FriendAdd adds friend to the friend list by sending a friend request, and
// returns the friend number on success.
//
// Friend numbers are unique identifiers used in all functions that operate on
// friends. Once added, a friend number is stable for the lifetime of the Tox
// object. After saving the state and reloading it, the friend numbers may not
// be the same as before. Deleting a friend creates a gap in the friend number
// set, which is filled by the next adding of a friend. Any pattern in friend
// numbers should not be relied on.
//
// TODO: rename "FriendAdd" => "AddFriend"?
func (t *Tox) FriendAdd(friendAddr string, message string) (uint32, error) {
	t.lock()
	defer t.unlock()

	if len(message) == 1 || len(message) > MaxFriendRequestLength {
		return 0, fmt.Errorf("friend request message must be in range [1, %d]: %d", MaxFriendRequestLength, len(message))
	}

	friendIDBytes, err := hex.DecodeString(friendAddr)
	if err != nil {
		return 0, err
	}

	// If more than INT32_MAX friends are added, this function causes undefined behavior.
	// TODO: Check current friend list size, and return error when needed.

	messageBytes := []byte(message)
	var cerr C.TOX_ERR_FRIEND_ADD

	friendNumber := uint32(
		C.tox_friend_add(
			t.toxcore,                       // Tox *tox
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
		return 0, fmt.Errorf("the length of the friend request message exceeded MaxFriendRequestLength(%d): %d", MaxFriendRequestLength, len(message))

	case C.TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return 0, fmt.Errorf("the friend request message was empty")

	case C.TOX_ERR_FRIEND_ADD_OWN_KEY:
		return 0, fmt.Errorf("the friend address belongs to the sending client")

	case C.TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return 0, fmt.Errorf("friend request has already been sent, or the address belongs to a friend that is already on the friend list: %s", friendAddr)

	case C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return 0, fmt.Errorf("the friend address checksum failed: %s", friendAddr)

	case C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return 0, fmt.Errorf("the friend was already there, but the nospam value was different: %s", friendAddr)

	case C.TOX_ERR_FRIEND_ADD_MALLOC:
		return 0, fmt.Errorf("memory allocation failed when trying to increase the friend list size")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendAddNorequest adds friend without sending friend request.
//
// This function is used to add a friend in response to a friend request. If the
// client receives a friend request, it can be reasonably sure that the other
// client added this client as a friend, eliminating the need for a friend
// request.
//
// This function is also useful in a situation where both instances are
// controlled by the same entity, so that this entity can perform the mutual
// friend adding. In this case, there is no need for a friend request, either.
//
// TODO: The argument friendID need a suitable name. "friend id", "friend
//       address", and "(friend) public key" are messing. More tox docs check?
// TODO: rename "FriendAddNorequest" => "AddFriendWithoutRequest" ?
func (t *Tox) FriendAddNorequest(friendID string) (uint32, error) {
	t.lock()
	defer t.unlock()

	if len(friendID) != PublicKeySize {
		return 0, fmt.Errorf("invalid friendID")
	}
	friendIDBytes, err := hex.DecodeString(friendID)
	if err != nil {
		return 0, err
	}

	var cerr C.TOX_ERR_FRIEND_ADD
	friendNumber := uint32(
		C.tox_friend_add_norequest(
			t.toxcore,                       // Tox *tox
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
		return 0, fmt.Errorf("the length of the friend request message exceeded MaxFriendRequestLength(%d)", MaxFriendRequestLength)

	case C.TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return 0, fmt.Errorf("the friend request message was empty")

	case C.TOX_ERR_FRIEND_ADD_OWN_KEY:
		return 0, fmt.Errorf("the friend address belongs to the sending client")

	case C.TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return 0, fmt.Errorf("friend request has already been sent, or the address belongs to a friend that is already on the friend list: %s", friendID)

	case C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return 0, fmt.Errorf("the friend address checksum failed: %s", friendID)

	case C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return 0, fmt.Errorf("the friend was already there, but the nospam value was different: %s", friendID)

	case C.TOX_ERR_FRIEND_ADD_MALLOC:
		return 0, fmt.Errorf("memory allocation failed when trying to increase the friend list size")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendByPublicKey returns the friend number associated with the specific
// Public Key.
//
// TODO: rename "FriendByPublicKey" => "GetFriendNumberByPublicKey"?
func (t *Tox) FriendByPublicKey(pubKey string) (uint32, error) {
	if len(pubKey) != PublicKeySize {
		return 0, fmt.Errorf("invalid public key")
	}
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return 0, err
	}

	var cerr C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	friendNumber := uint32(
		C.tox_friend_by_public_key(
			t.toxcore,                     // const Tox *tox
			(*C.uint8_t)(&pubKeyBytes[0]), // const uint8_t *public_key
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

// FriendGetPublicKey returns a copy of the Public Key associated with a given
// friend number.
//
// TODO: rename "FriendGetPublicKey" => "GetPublicKeyByFriendNumber" or
//       "GetFriendPublicKeyByNumber"?
func (t *Tox) FriendGetPublicKey(friendNumber uint32) (string, error) {
	var pubKeyBytes = make([]byte, PublicKeySize)

	var cerr C.TOX_ERR_FRIEND_GET_PUBLIC_KEY

	ok := bool(
		C.tox_friend_get_public_key(
			t.toxcore,                     // const Tox *tox
			C.uint32_t(friendNumber),      // uint32_t friend_number
			(*C.uint8_t)(&pubKeyBytes[0]), // uint8_t *public_key
			&cerr, // TOX_ERR_FRIEND_GET_PUBLIC_KEY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK:
		assert(ok, "tox_friend_get_public_key() return 'false' on success")

		return strings.ToUpper(hex.EncodeToString(pubKeyBytes)), nil

	case C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("no friend with the given number exists on the friend list")

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendDelete removes friend from the friend list.
//
// This does not notify the friend of their deletion. After calling this
// function, this client will appear offline to the friend and no communication
// can occur between the two.
//
// TODO: rename "FriendDelete" = "DeleteFriend"?
func (t *Tox) FriendDelete(friendNumber uint32) error {
	t.lock()
	defer t.unlock()

	var cerr C.TOX_ERR_FRIEND_DELETE

	ok := bool(
		C.tox_friend_delete(
			t.toxcore,                // Tox *tox
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

// FriendGetConnectionStatus returns the friend's connection status as it was
// received through the friendConnectionStatus event.
//
// TODO: @deprecated This getter is deprecated. Use the event and store the
//       status in the client state.
// TODO: rename "FriendGetConnectionStatus" => "GetFriendConnectionStatus"?
func (t *Tox) FriendGetConnectionStatus(friendNumber uint32) (ConnectionType, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	connStatus := ConnectionType(
		C.tox_friend_get_connection_status(
			t.toxcore,                // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return connStatus, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return ConnectionNone, fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return ConnectionNone, fmt.Errorf("friendNumber did not designate a valid friend: %d", friendNumber)

	default:
		return ConnectionNone, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendExists checks if a friend with the given friend number exists and
// returns true if it does.
func (t *Tox) FriendExists(friendNumber uint32) bool {
	return bool(
		C.tox_friend_exists(
			t.toxcore,                // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
		),
	)
}

// friendSendMessage sends text chat message to an online friend.
// This function creates a chat message packet and pushes it into the
// send queue.
//
// The message length may not exceed TOX_MAX_MESSAGE_LENGTH. Larger messages
// must be split by the client and sent as separate messages. Other clients can
// then reassemble the fragments. Messages may not be empty.
//
// The return value of this function is the message ID. If a read receipt is
// received, the triggered OnFriendReadReceipt event will be passed this message
// ID.
//
// Message IDs are unique per friend. The first message ID is 0. Message IDs are
// incremented by 1 each time a message is sent. If math.MaxUint32 messages were
// sent, the next message ID is 0.
//
// TODO: rename "friendSendMessage" => "sendFriendMessage"?
func (t *Tox) friendSendMessage(friendNumber uint32, msg string, msgType C.TOX_MESSAGE_TYPE) (msgID uint32, err error) {
	if len(msg) > MaxMessageLength {
		return 0, fmt.Errorf("length of message over ranged (max: %d): %d", MaxMessageLength, len(msg))
	}

	t.lock()
	defer t.unlock()

	var messageBytes = []byte(msg)
	var cerr C.TOX_ERR_FRIEND_SEND_MESSAGE

	msgID = uint32(
		C.tox_friend_send_message(
			t.toxcore,                // Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			msgType,                  // TOX_MESSAGE_TYPE type
			(*C.uint8_t)(&messageBytes[0]), // const uint8_t *message
			C.size_t(len(msg)),             // size_t length
			&cerr,                          // TOX_ERR_FRIEND_SEND_MESSAGE *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_SEND_MESSAGE_OK:
		return msgID, nil

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_NULL:
		return 0, fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND:
		return 0, fmt.Errorf("friend number did not designate a valid friend")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED:
		return 0, fmt.Errorf("client is currently not connected to the friend")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ:
		return 0, fmt.Errorf("allocation error occurred while increasing the send queue size")

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG:
		return 0, fmt.Errorf("message length exceeded MaxMessageLength(%d): %d", MaxMessageLength, len(msg))

	case C.TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY:
		return 0, fmt.Errorf("attempted to send a zero-length message")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendSendMessage sends text chat message to an online friend.
// See friendSendMessage.
//
// TODO: rename "FriendSendMessage" => "SendFriendMessage"?
func (t *Tox) FriendSendMessage(friendNumber uint32, message string) (msgID uint32, err error) {
	return t.friendSendMessage(friendNumber, message, C.TOX_MESSAGE_TYPE_NORMAL)
}

// FriendSendAction sends user action to an online friend.
// See friendSendMessage.
//
// TODO: rename "FriendSendAction" => "SendFriendAction"?
func (t *Tox) FriendSendAction(friendNumber uint32, action string) (msgID uint32, err error) {
	return t.friendSendMessage(friendNumber, action, C.TOX_MESSAGE_TYPE_ACTION)
}

// SelfSetName sets nickname for the Tox client.
//
// TODO: rename "SelfSetName" => "SetSelfName" or "SetName"?
func (t *Tox) SelfSetName(name string) error {
	if len(name) > MaxNameLength {
		return fmt.Errorf("length nickname is over ranged (max: %d): %d", MaxNameLength, len(name))
	}

	t.lock()
	defer t.unlock()

	var nameBytes = []byte(name)

	var cerr C.TOX_ERR_SET_INFO
	ok := bool(
		C.tox_self_set_name(
			t.toxcore,                   // Tox *tox
			(*C.uint8_t)(&nameBytes[0]), // const uint8_t *name
			C.size_t(len(nameBytes)),    // size_t length
			&cerr, // TOX_ERR_SET_INFO *error
		),
	)

	switch cerr {
	case C.TOX_ERR_SET_INFO_OK:
		assert(ok, "tox_self_set_name() return 'false' on success")

		return nil

	case C.TOX_ERR_SET_INFO_NULL:
		return fmt.Errorf("one of the arguments to the function was NULL when it was not expected")

	case C.TOX_ERR_SET_INFO_TOO_LONG:
		return fmt.Errorf("length exceeded maximum permissible size(%d): %d", MaxNameLength, len(name))

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// selfGetNameSize returns the length of the current nickname.
//
// If no nickname was set before calling this function, the name is empty, and
// this function returns 0.
//
// TODO: rename "selfGetNameSize" => "getSelfNameSize" or "getNameSize"?
func (t *Tox) selfGetNameSize() int {
	// TODO: thread safe is needed
	return int(C.tox_self_get_name_size(t.toxcore))
}

// SelfGetName returns a copy of nickname.
//
// TODO: rename "SelfGetName" => "GetSelfName" or "GetName"?
func (t *Tox) SelfGetName() string {
	nameBytes := make([]byte, t.selfGetNameSize())

	C.tox_self_get_name(
		t.toxcore,                        // const Tox *tox
		(*C.uint8_t)(safeptr(nameBytes)), // uint8_t *name
	)

	return string(nameBytes)
}

// FriendGetName returns the friend name by given friend number.
//
// The return value is equal to the `name` argument received by the last
// OnFriendNameChanged event.
//
// TODO: rename "FriendGetName" => "GetFriendName"?
func (t *Tox) FriendGetName(friendNumber uint32) (string, error) {
	nameSize, err := t.friendGetNameSize(friendNumber)
	if err != nil {
		return "", err
	}
	nameBytes := make([]byte, nameSize)
	var cerr C.TOX_ERR_FRIEND_QUERY

	ok := bool(
		C.tox_friend_get_name(
			t.toxcore,                        // const Tox *tox
			C.uint32_t(friendNumber),         // uint32_t friend_number
			(*C.uint8_t)(safeptr(nameBytes)), // uint8_t *name
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		assert(ok, "tox_friend_get_name() return 'false' on success")

		return string(nameBytes), nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return "", fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// friendGetNameSize returns the length of the friend's name. If the friend
// number is invalid, the return value is unspecified.
//
// The return value is equal to the `length` argument received by the last
// OnFriendNameChanged event.
//
// TODO: rename "friendGetNameSize" => "getFriendNameSize"?
func (t *Tox) friendGetNameSize(friendNumber uint32) (int, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	nameSize := int(
		C.tox_friend_get_name_size(
			t.toxcore,                // const Tox *tox
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

// SelfSetStatusMessage sets client's status message.
//
// TODO: rename "SelfSetStatusMessage" => "SetSelfStatusMessage" or
//       "SetStatusMessage"?
func (t *Tox) SelfSetStatusMessage(status string) error {
	if len(status) > MaxStatusMessageLength {
		return fmt.Errorf("status is over ranged (max: %d): %d", MaxStatusMessageLength, len(status))
	}

	t.lock()
	defer t.unlock()

	var _status = []byte(status)

	var cerr C.TOX_ERR_SET_INFO
	ok := bool(
		C.tox_self_set_status_message(
			t.toxcore,                 // Tox *tox
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
		return fmt.Errorf("length exceeded maximum permissible size(%d): %d", MaxNameLength, len(status))

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// SelfSetStatus sets client's user status.
//
// TODO: rename "SelfSetStatus" => "SetSelfStatus" or "SetStatus"?
func (t *Tox) SelfSetStatus(status UserStatus) {
	C.tox_self_set_status(
		t.toxcore,                 // Tox *tox
		C.TOX_USER_STATUS(status), // TOX_USER_STATUS status
	)
}

// friendGetStatusMessageSize returns the length of the friend's status message.
//
// TODO: rename "friendGetStatusMessageSize" => "getFriendStatusMessageSize"?
func (t *Tox) friendGetStatusMessageSize(friendNumber uint32) (int, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	size := int(
		C.tox_friend_get_status_message_size(
			t.toxcore,                // const Tox *tox
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

// FriendGetStatusMessage returns status message of friend specified by given
// friend number.
//
// The data written to `status_message` is equal to the data received by the
// last OnFriendStatusMessageChanged event.
//
// TODO: rename "FriendGetStatusMessage" => "GetFriendStatusMessage"?
func (t *Tox) FriendGetStatusMessage(friendNumber uint32) (string, error) {
	size, err := t.friendGetStatusMessageSize(friendNumber)
	if err != nil {
		return "", err
	}
	buff := make([]byte, size)

	var cerr C.TOX_ERR_FRIEND_QUERY
	ok := bool(
		C.tox_friend_get_status_message(
			t.toxcore,                   // const Tox *tox
			C.uint32_t(friendNumber),    // uint32_t friend_number
			(*C.uint8_t)(safeptr(buff)), // uint8_t *status_message
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		assert(ok, "tox_friend_get_status_message() return 'false' on success")

		return string(buff[:]), nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return "", fmt.Errorf("the pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return "", fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return "", fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// selfGetStatusMessageSize returns the length of the current status message.
//
// TODO: rename "selfGetStatusMessageSize" => "getSelfStatusMessageSize" or
//       "getStatusMessageSize"?
func (t *Tox) selfGetStatusMessageSize() int {
	return int(C.tox_self_get_status_message_size(t.toxcore))
}

// SelfGetStatusMessage returns a copy of status message.
//
// TODO: rename "SelfGetStatusMessage" => "GetSelfStatusMessage" or
//       "GetStatusMessage"?
func (t *Tox) SelfGetStatusMessage() (string, error) {
	var buff = make([]byte, t.selfGetStatusMessageSize())

	C.tox_self_get_status_message(
		t.toxcore,                   // const Tox *tox
		(*C.uint8_t)(safeptr(buff)), // uint8_t *status_message
	)

	return string(buff[:]), nil
}

// FriendGetStatus returns the friend's user status (away/busy/...). If the
// friend number is invalid, the return value is unspecified.
//
// TODO: @deprecated This getter is deprecated. Use the event and store the
//       status in the client state.
func (t *Tox) FriendGetStatus(friendNumber uint32) (UserStatus, error) {
	friendNumberInFriendList := func(list []uint32, value uint32) bool {
		for _, v := range list {
			if value == v {
				return true
			}
		}
		return false
	}(t.SelfGetFriendList(), friendNumber)
	if !friendNumberInFriendList {
		return UserStatusNone, fmt.Errorf("friend is not in friend list")
	}

	var cerr C.TOX_ERR_FRIEND_QUERY

	userStatus := UserStatus(
		C.tox_friend_get_status(
			t.toxcore,                // const Tox *tox
			C.uint32_t(friendNumber), // uint32_t friend_number
			&cerr, // TOX_ERR_FRIEND_QUERY *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FRIEND_QUERY_OK:
		return userStatus, nil

	case C.TOX_ERR_FRIEND_QUERY_NULL:
		return UserStatusNone, fmt.Errorf("pointer parameter for storing the query result (name, message) was nil")

	case C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND:
		return UserStatusNone, fmt.Errorf("friendNumber did not designate a valid friend: %X", friendNumber)

	default:
		return UserStatusNone, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// SelfGetStatus returns the client's user status.
//
// TODO: rename "SelfGetStatus" => "GetSelfStatus" or "GetStatus"?
func (t *Tox) SelfGetStatus() UserStatus {
	return UserStatus(C.tox_self_get_status(t.toxcore))
}

// FriendGetLastOnline returns the time the friend, specified by given friend
// number, was seen online.
//
// TODO: rename "FriendGetLastOnline" => "GetFriendLastOnline"?
func (t *Tox) FriendGetLastOnline(friendNumber uint32) (time.Time, error) {
	var nullTime time.Time

	var cerr C.TOX_ERR_FRIEND_GET_LAST_ONLINE
	timestamp := uint64(
		C.tox_friend_get_last_online(
			t.toxcore,                // const Tox *tox
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

// SelfSetTyping sets client's typing status for a friend specified by
// friendNumber.
//
// TODO: rename "SelfSetTyping" => "SetSelfTyping" or "SetTyping"?
func (t *Tox) SelfSetTyping(friendNumber uint32, typing bool) error {
	t.lock()
	defer t.unlock()

	var cerr C.TOX_ERR_SET_TYPING

	ok := bool(
		C.tox_self_set_typing(
			t.toxcore,                // Tox *tox
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

// FriendGetTyping checks whether a friend is currently typing a message.
//
// TODO: @deprecated This getter is deprecated. Use the event and store the
//       status in the client state.
// TODO: rename "FriendGetTyping" => "GetFriendTyping"?
func (t *Tox) FriendGetTyping(friendNumber uint32) (bool, error) {
	var cerr C.TOX_ERR_FRIEND_QUERY

	typing := bool(
		C.tox_friend_get_typing(
			t.toxcore,                // const Tox *tox
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

// SelfGetFriendListSize returns the number of friends on the friend list.
//
// TODO: rename "SelfGetFriendListSize" => "GetSelfFriendListSize" or
//       "GetFriendListSize"?
func (t *Tox) SelfGetFriendListSize() int {
	return int(C.tox_self_get_friend_list_size(t.toxcore))
}

// SelfGetFriendList returns a copy of friend numbers on the friend list.
//
// TODO: rename "SelfGetFriendList" => "GetSelfFriendList" OR "GetFriendList"?
func (t *Tox) SelfGetFriendList() []uint32 {
	size := t.SelfGetFriendListSize()
	vec := make([]uint32, size)
	if size == 0 {
		return vec // fast return on zero length friend list
	}

	C.tox_self_get_friend_list(
		t.toxcore, // const Tox *tox
		(*C.uint32_t)(unsafe.Pointer(&vec[0])), // uint32_t *friend_list
	)

	return vec
}

// tox_callback_***

// SelfGetNoSpam returns the 4-byte nospam part of the address.
// The value is returned in host byte order.
//
// TODO: rename "SelfGetNoSpam" => "GetSelfNoSpam" or "GetNoSpam"?
func (t *Tox) SelfGetNoSpam() uint32 {
	return uint32(C.tox_self_get_nospam(t.toxcore))
}

// SelfGetNoSpamString returns the 4-byte nospam part of the address in string
// format. The value is upper cased.
//
// TODO: rename "SelfGetNoSpamString" => "GetSelfNoSpamString" or
//       "GetNoSpamString"?
func (t *Tox) SelfGetNoSpamString() string {
	return fmt.Sprintf("%X", t.SelfGetNoSpam())
}

// SelfSetNoSpam sets the 4-byte noSpam part of the address. The value is
// expected in host byte order. I.e. 0x12345678 will form the bytes
// [12, 34, 56, 78] in the noSpam part of the Tox friend address.
//
// TODO: rename "SelfSetNoSpam" => "SetSelfNoSpam" or "SetNoSpam"?
func (t *Tox) SelfSetNoSpam(noSpam uint32) {
	t.lock()
	defer t.unlock()

	C.tox_self_set_nospam(
		t.toxcore,          // Tox *tox
		C.uint32_t(noSpam), // uint32_t nospam
	)
}

// SelfSetNoSpamString
//
// TODO: rename "SelfSetNoSpamString" => "SetSelfNoSpamString" or
//       "SetNoSpamString"?
func (t *Tox) SelfSetNoSpamString(noSpam string) error {
	if len(noSpam) != 8 {
		return fmt.Errorf("invalid NoSpam format, which should be a 8-char hex string")
	}

	var noSpamNum uint32
	_, err := fmt.Sscanf(noSpam, "%8x", &noSpamNum)
	if err != nil {
		return err
	}

	t.SelfSetNoSpam(noSpamNum)

	return nil
}

// SelfGetPublicKey returns the Tox Public Key.
//
// TODO: rename "SelfGetPublicKey" => "GetSelfPublicKey" or "GetPublicKey"?
func (t *Tox) SelfGetPublicKey() string {
	var _pubkey [PublicKeySize]byte

	C.tox_self_get_public_key(
		t.toxcore,                 // const Tox *tox
		(*C.uint8_t)(&_pubkey[0]), // uint8_t *public_key
	)

	return strings.ToUpper(hex.EncodeToString(_pubkey[:]))
}

// SelfGetSecretKey returns the Tox Secret Key.
//
// TODO: rename "SelfGetSecretKey" => "GetSelfSecretKey" or "GetSecretKey"?
func (t *Tox) SelfGetSecretKey() string {
	var _seckey [SecretKeySize]byte

	C.tox_self_get_secret_key(
		t.toxcore,                 // const Tox *tox
		(*C.uint8_t)(&_seckey[0]), // uint8_t *secret_key
	)

	return strings.ToUpper(hex.EncodeToString(_seckey[:]))
}

// tox_lossy_***

// FriendSendLossyPacket sends a custom lossy packet to a friend.
//
// Lossy packets behave like UDP packets, meaning they might never reach the
// other side or might arrive more than once (if someone is messing with the
// connection) or might arrive in the wrong order.
//
// Unless latency is an issue, it is recommended that you use lossless custom
// packets instead.
//
// TODO: rename "FriendSendLossyPacket" => "SendFriendLossyPacket" or
//       "SendLossyPacket"?
func (t *Tox) FriendSendLossyPacket(friendNumber uint32, data string) error {
	if len(data) > MaxCustomPacketSize {
		return fmt.Errorf("length of data is out of range (max: %d): %d", MaxCustomPacketSize, len(data))
	}

	t.lock()
	defer t.unlock()

	var dataBytes = []byte(data)
	if 200 > dataBytes[0] || dataBytes[0] < 254 {
		return fmt.Errorf("the first byte of data must be in the range 200-254: %d", dataBytes[0])
	}

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET

	ok := bool(
		C.tox_friend_send_lossy_packet(
			t.toxcore,                   // Tox *tox
			C.uint32_t(friendNumber),    // uint32_t friend_number
			(*C.uint8_t)(&dataBytes[0]), // const uint8_t *data
			C.size_t(len(data)),         // size_t length
			&cerr,                       // TOX_ERR_FRIEND_CUSTOM_PACKET *error
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
		return fmt.Errorf("packet data length exceeded MaxCustomPacketSize(%d): %d", MaxCustomPacketSize, len(data))

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		return fmt.Errorf("packet queue is full")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FriendSendLosslessPacket sends a custom lossless packet to a friend.
//
// Lossless packet behavior is comparable to TCP (reliability, arrive in order)
// but with packets instead of a stream.
//
// TODO: rename "FriendSendLosslessPacket" => "SendFriendLosslessPacket" or
//       "SendLosslessPacket"?
func (t *Tox) FriendSendLosslessPacket(friendNumber uint32, data string) error {
	if len(data) > MaxCustomPacketSize {
		return fmt.Errorf("length of data is out of range (max: %d): %d", MaxCustomPacketSize, len(data))
	}

	t.lock()
	defer t.unlock()

	var dataBytes = []byte(data)
	if 160 > dataBytes[0] || dataBytes[0] > 191 {
		return fmt.Errorf("the first byte of data must be in the range 160-191: %d", dataBytes[0])
	}

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET

	ok := bool(
		C.tox_friend_send_lossless_packet(
			t.toxcore,                   // Tox *tox
			C.uint32_t(friendNumber),    // uint32_t friend_number
			(*C.uint8_t)(&dataBytes[0]), // const uint8_t *data
			C.size_t(len(data)),         // size_t length
			&cerr,                       // TOX_ERR_FRIEND_CUSTOM_PACKET *error
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
		return fmt.Errorf("packet data length exceeded MaxCustomPacketSize(%d): %d", MaxCustomPacketSize, len(data))

	case C.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		return fmt.Errorf("packet queue is full")

	default:
		return fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// tox_callback_avatar_**

// Hash generates a cryptographic hash of the given data.
//
// This function may be used by clients for any purpose, but is provided
// primarily for validating cached avatars. This use is highly recommended to
// avoid unnecessary avatar updates.
func (t *Tox) Hash(data string) (string, bool, error) {
	hashBytes := make([]byte, HashLength)
	dataBytes := []byte(data)

	ok := bool(
		C.tox_hash(
			(*C.uint8_t)(&hashBytes[0]), // uint8_t *hash
			(*C.uint8_t)(&dataBytes[0]), // const uint8_t *data
			C.size_t(len(data)),         // size_t length
		),
	)

	// If hash is NULL or data is NULL while length is not 0 the
	// function returns false, otherwise it returns true.
	assert(ok && len(data) > 0, "tox_hash() return 'false' on success")

	return string(hashBytes), ok, nil
}

// tox_callback_file_***

// FileControl sends a file control command to a friend for a given file
// transfer.
func (t *Tox) FileControl(friendNumber uint32, fileNumber uint32, control FileControlType) (bool, error) {
	var cerr C.TOX_ERR_FILE_CONTROL

	ok := bool(
		C.tox_file_control(
			t.toxcore,                   // Tox *tox
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

// FileSend sends a file transmission request.
//
// fileSize can be set to math.MaxUint64 for streaming data of unknown size.
//
// File transmission occurs in chunks, which are requested through the
// OnFileChunkRequest event.
//
// When a friend goes offline, all file transfers associated with the friend are
// purged from core.
//
// If the file contents change during a transfer, the behavior is unspecified in
// general. What will actually happen depends on the mode in which the file was
// modified and how the client determines the file size.
//
// - If the file size was increased
//   - and sending mode was streaming (fileSize = math.MaxUint64), the behavior
//     will be as expected.
//   - and sending mode was file (fileSize != math.MaxUint64), the
//     OnFileChunkRequest event will receive length = 0 when Core thinks the
//     file transfer has finished. If the client remembers the file size as it
//     was when sending the request, it will terminate the transfer normally. If
//     the client re-reads the size, it will think the friend cancelled the
//     transfer.
// - If the file size was decreased
//   - and sending mode was streaming, the behavior is as expected.
//   - and sending mode was file, the callback will return 0 at the new
//     (earlier) end-of-file, signalling to the friend that the transfer was
//     cancelled.
// - If the file contents were modified
//   - at a position before the current read, the two files (local and remote)
//     will differ after the transfer terminates.
//   - at a position after the current read, the file transfer will succeed as
//     expected.
//   - In either case, both sides will regard the transfer as complete and
//     successful.
//
// TODO: rename "FileSend" => "SendFile"?
func (t *Tox) FileSend(friendNumber uint32, kind FileKind, fileSize uint64, fileID, fileName string) (uint32, error) {
	t.lock()
	defer t.unlock()

	if len(fileName) > MaxFilenameLength {
		return 0, fmt.Errorf("fileName length over range (max: %d): %d", MaxFilenameLength, len(fileName))
	}

	fileIDBytes := []byte(fileID)
	fileNameBytes := []byte(fileName)
	var cerr C.TOX_ERR_FILE_SEND

	fileNumber := uint32(
		C.tox_file_send(
			t.toxcore,                       // Tox *tox
			C.uint32_t(friendNumber),        // uint32_t friend_number
			C.uint32_t(kind),                // uint32_t kind
			C.uint64_t(fileSize),            // uint64_t file_size
			(*C.uint8_t)(&fileIDBytes[0]),   // const uint8_t *file_id
			(*C.uint8_t)(&fileNameBytes[0]), // const uint8_t *filename
			C.size_t(len(fileName)),         // size_t filename_length
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
		return 0, fmt.Errorf("filename length exceeded MaxFilenameLength (%d) bytes: %d", MaxFilenameLength, len(fileName))

	case C.TOX_ERR_FILE_SEND_TOO_MANY:
		return 0, fmt.Errorf("too many ongoing transfers. allowed in 256 per friend per direction (sending and receiving)")

	default:
		return 0, fmt.Errorf("unknown error code, maybe go-toxcore-c is outdated from c-toxcore: %d", cerr)
	}
}

// FileSendChunk sends a chunk of file data to a friend.
//
// This function is called in response to the OnFileChunkRequest event. The
// length parameter should be equal to the one received though the callback. If
// it is zero, the transfer is assumed complete. For files with known size, Core
// will know that the transfer is complete after the last byte has been
// received, so it is not necessary (though not harmful) to send a zero-length
// chunk to terminate. For streams, core will know that the transfer is finished
// if a chunk with length less than the length requested in the callback is
// sent.
//
// TODO: rename "FileSendChunk" => "SendFileChunk"?
func (t *Tox) FileSendChunk(friendNumber uint32, fileNumber uint32, position uint64, data []byte) (bool, error) {
	if data == nil || len(data) == 0 {
		return false, toxerr("empty data")
	}

	t.lock()
	defer t.unlock()

	var cerr C.TOX_ERR_FILE_SEND_CHUNK

	ok := bool(
		C.tox_file_send_chunk(
			t.toxcore,                // Tox *tox
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

// FileSeek sends a file seek control command to a friend for a given file
// transfer.
//
// This function can only be called to resume a file transfer right before
// FileControlResume is sent.
func (t *Tox) FileSeek(friendNumber uint32, fileNumber uint32, position uint64) (bool, error) {
	t.lock()
	defer t.unlock()

	var cerr C.TOX_ERR_FILE_SEEK

	ok := bool(
		C.tox_file_seek(
			t.toxcore,                // Tox *tox
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

// FileGetFileID returns a copy of the file id associated to the file transfer.
//
// TODO: rename "FileGetFileID" => "GetFileID"
func (t *Tox) FileGetFileID(friendNumber uint32, fileNumber uint32) (string, error) {
	var fileIDBytes = make([]byte, FileIDLength)

	var cerr C.TOX_ERR_FILE_GET
	ok := bool(
		C.tox_file_get_file_id(
			t.toxcore,                     // const Tox *tox
			C.uint32_t(fileNumber),        // uint32_t friend_number
			C.uint32_t(fileNumber),        // uint32_t file_number
			(*C.uint8_t)(&fileIDBytes[0]), // uint8_t *file_id
			&cerr, // TOX_ERR_FILE_GET *error
		),
	)

	switch cerr {
	case C.TOX_ERR_FILE_GET_OK:
		assert(ok, "tox_file_get_file_id() return 'false' on success")

		return strings.ToUpper(hex.EncodeToString(fileIDBytes)), nil

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

func (t *Tox) IsConnected() bool {
	status := C.tox_self_get_connection_status(t.toxcore)

	return status == C.TOX_CONNECTION_TCP || status == C.TOX_CONNECTION_UDP
}

func (t *Tox) putcbevts(f func()) {
	t.cbevts = append(t.cbevts, f)
}
