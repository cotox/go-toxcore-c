package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

extern void toxCallbackLog(Tox*, TOX_LOG_LEVEL, char*, uint32_t, char*, char*);

*/
import "C"
import "unsafe"

// TODO: define SavedataType for type-check at compile time.
//
// type SavedataType int
// const (
//     SavedataTypeNone      SavedataType = int(C.TOX_SAVEDATA_TYPE_NONE)
//     SavedataTypeToxSave   SavedataType = int(C.TOX_SAVEDATA_TYPE_TOX_SAVE)
//     SavedataTypeSecretKey SavedataType = int(C.TOX_SAVEDATA_TYPE_SECRET_KEY)
// )
const (
	SavedataTypeNone      = int(C.TOX_SAVEDATA_TYPE_NONE)
	SavedataTypeToxSave   = int(C.TOX_SAVEDATA_TYPE_TOX_SAVE)
	SavedataTypeSecretKey = int(C.TOX_SAVEDATA_TYPE_SECRET_KEY)
)

// TODO: define ProxyType for type-check at compile time.
//
// type ProxyType int
// const (
//     ProxyTypeNone   ProxyType = int(C.TOX_PROXY_TYPE_NONE)
//     ProxyTypeHTTP   ProxyType = int(C.TOX_PROXY_TYPE_HTTP)
//     ProxyTypeSOCKS5 ProxyType = int(C.TOX_PROXY_TYPE_SOCKS5)
// )
const (
	ProxyTypeNone   = int(C.TOX_PROXY_TYPE_NONE)
	ProxyTypeHTTP   = int(C.TOX_PROXY_TYPE_HTTP)
	ProxyTypeSOCKS5 = int(C.TOX_PROXY_TYPE_SOCKS5)
)

// TODO: define LogLevel for type-check at compile time.
//
// type LogLevel int
// const (
//     LogLevelTrace   LogLevel = int(C.TOX_LOG_LEVEL_TRACE)
//     LogLevelDebug   LogLevel = int(C.TOX_LOG_LEVEL_DEBUG)
//     LogLevelInfo    LogLevel = int(C.TOX_LOG_LEVEL_INFO)
//     LogLevelWarning LogLevel = int(C.TOX_LOG_LEVEL_WARNING)
//     LogLevelError   LogLevel = int(C.TOX_LOG_LEVEL_ERROR)
// )
const (
	LogLevelTrace   = int(C.TOX_LOG_LEVEL_TRACE)
	LogLevelDebug   = int(C.TOX_LOG_LEVEL_DEBUG)
	LogLevelInfo    = int(C.TOX_LOG_LEVEL_INFO)
	LogLevelWarning = int(C.TOX_LOG_LEVEL_WARNING)
	LogLevelError   = int(C.TOX_LOG_LEVEL_ERROR)
)

// TODO: rename "ToxOptions" => "Options"
type ToxOptions struct {
	IPv6Enabled           bool
	UDPEnabled            bool
	ProxyType             int32
	ProxyHost             string
	ProxyPort             uint16
	SavedataType          int
	SavedataData          []byte
	TCPPort               uint16
	LocalDiscoveryEnabled bool
	StartPort             uint16
	EndPort               uint16
	HolePunchingEnabled   bool
	ThreadSafe            bool
	LogCallback           func(_ *Tox, level int, file string, line uint32, fname string, msg string)
}

// TODO: rename "NewToxOptions()" => "NewDefaultOptions()"
func NewToxOptions() *ToxOptions {
	cToxOpts := C.tox_options_new(nil)
	defer C.tox_options_free(cToxOpts)

	opts := new(ToxOptions)
	opts.IPv6Enabled = bool(C.tox_options_get_ipv6_enabled(cToxOpts))
	opts.UDPEnabled = bool(C.tox_options_get_udp_enabled(cToxOpts))
	opts.ProxyType = int32(C.tox_options_get_proxy_type(cToxOpts))
	opts.ProxyPort = uint16(C.tox_options_get_proxy_port(cToxOpts))
	opts.TCPPort = uint16(C.tox_options_get_tcp_port(cToxOpts))
	opts.LocalDiscoveryEnabled = bool(C.tox_options_get_local_discovery_enabled(cToxOpts))
	opts.StartPort = uint16(C.tox_options_get_start_port(cToxOpts))
	opts.EndPort = uint16(C.tox_options_get_end_port(cToxOpts))
	opts.HolePunchingEnabled = bool(C.tox_options_get_hole_punching_enabled(cToxOpts))

	return opts
}

func (this *ToxOptions) toCToxOptions() *C.struct_Tox_Options {
	cToxOpts := C.tox_options_new(nil)
	C.tox_options_default(cToxOpts)
	C.tox_options_set_ipv6_enabled(cToxOpts, (C._Bool)(this.IPv6Enabled))
	C.tox_options_set_udp_enabled(cToxOpts, (C._Bool)(this.UDPEnabled))

	if this.SavedataData != nil {
		C.tox_options_set_savedata_data(cToxOpts, (*C.uint8_t)(&this.SavedataData[0]), C.size_t(len(this.SavedataData)))
		C.tox_options_set_savedata_type(cToxOpts, C.TOX_SAVEDATA_TYPE(this.SavedataType))
	}
	C.tox_options_set_tcp_port(cToxOpts, (C.uint16_t)(this.TCPPort))

	C.tox_options_set_proxy_type(cToxOpts, C.TOX_PROXY_TYPE(this.ProxyType))
	C.tox_options_set_proxy_port(cToxOpts, C.uint16_t(this.ProxyPort))
	if len(this.ProxyHost) > 0 {
		C.tox_options_set_proxy_host(cToxOpts, C.CString(this.ProxyHost))
	}

	C.tox_options_set_local_discovery_enabled(cToxOpts, C._Bool(this.LocalDiscoveryEnabled))
	C.tox_options_set_start_port(cToxOpts, C.uint16_t(this.StartPort))
	C.tox_options_set_end_port(cToxOpts, C.uint16_t(this.EndPort))
	C.tox_options_set_hole_punching_enabled(cToxOpts, C._Bool(this.HolePunchingEnabled))

	C.tox_options_set_log_callback(cToxOpts, (*C.tox_log_cb)((unsafe.Pointer)(C.toxCallbackLog)))

	return cToxOpts
}

//export toxCallbackLog
func toxCallbackLog(cTox *C.Tox, level C.TOX_LOG_LEVEL, file *C.char, line C.uint32_t, fname *C.char, msg *C.char) {
	t := cbUserDatas.get(cTox)
	if t != nil && t.opts != nil && t.opts.LogCallback != nil {
		t.opts.LogCallback(t, int(level), C.GoString(file), uint32(line), C.GoString(fname), C.GoString(msg))
	}
}

type BootNode struct {
	Addr   string
	Port   int
	Pubkey string
}
