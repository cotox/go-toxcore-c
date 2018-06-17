package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

extern void toxCallbackLog(Tox*, TOX_LOG_LEVEL, char*, uint32_t, char*, char*);
*/
import "C"
import "unsafe"

// Type of savedata to create the Tox instance from.
const (
	// No savedata.
	SavedataTypeNone = int(C.TOX_SAVEDATA_TYPE_NONE)

	// Savedata is one that was obtained from ${savedata.get}.
	SavedataTypeToxSave = int(C.TOX_SAVEDATA_TYPE_TOX_SAVE)

	// Savedata is a secret key of length $SECRET_KEY_SIZE.
	SavedataTypeSecretKey = int(C.TOX_SAVEDATA_TYPE_SECRET_KEY)
)

// Type of proxy used to connect to TCP relays.
const (
	// Don't use a proxy.
	ProxyTypeNone = int(C.TOX_PROXY_TYPE_NONE)

	// HTTP proxy using CONNECT.
	ProxyTypeHTTP = int(C.TOX_PROXY_TYPE_HTTP)

	// SOCKS proxy for simple socket pipes.
	ProxyTypeSOCKS5 = int(C.TOX_PROXY_TYPE_SOCKS5)
)

// Severity level of log messages.
const (
	// Very detailed traces including all network activity.
	LogLevelTrace = int(C.TOX_LOG_LEVEL_TRACE)

	// Debug messages such as which port we bind to.
	LogLevelDebug = int(C.TOX_LOG_LEVEL_DEBUG)

	// Informational log messages such as video call status changes.
	LogLevelInfo = int(C.TOX_LOG_LEVEL_INFO)

	// Warnings about internal inconsistency or logic errors.
	LogLevelWarning = int(C.TOX_LOG_LEVEL_WARNING)

	// Severe unexpected errors caused by external or internal inconsistency.
	LogLevelError = int(C.TOX_LOG_LEVEL_ERROR)
)

// ToxOptions contains all the startup options for Tox. You must $new to
// allocate an object of this type.
//
// WARNING: Although this struct happens to be visible in the API, it is
// effectively private. Do not allocate this yourself or access members
// directly, as it *will* break binary compatibility frequently.
//
// @deprecated The memory layout of this struct (size, alignment, and field
// order) is not part of the ABI. To remain compatible, prefer to use $new to
// allocate the object and accessor functions to set the members. The struct
// will become opaque (i.e. the definition will become private) in v0.3.0.
type ToxOptions struct {

	// The type of socket to create.
	//
	// If this is set to false, an IPv4 socket is created, which subsequently
	// only allows IPv4 communication.
	// If it is set to true, an IPv6 socket is created, allowing both IPv4 and
	// IPv6 communication.
	IPv6Enabled bool

	// Enable the use of UDP communication when available.
	//
	// Setting this to false will force Tox to use TCP only. Communications will
	// need to be relayed through a TCP relay node, potentially slowing them down.
	// Disabling UDP support is necessary when using anonymous proxies or Tor.
	UDPEnabled bool

	// Pass communications through a proxy.
	ProxyType int32

	// The IP address or DNS name of the proxy to be used.
	//
	// If used, this must be non-NULL and be a valid DNS name. The name must not
	// exceed 255 characters, and be in a NUL-terminated C string format
	// (255 chars + 1 NUL byte).
	//
	// This member is ignored (it can be NULL) if proxy_type is TOX_PROXY_TYPE_NONE.
	//
	// The data pointed at by this member is owned by the user, so must
	// outlive the options object.
	ProxyHost string

	// The port to use to connect to the proxy server.
	//
	// Ports must be in the range (1, 65535). The value is ignored if
	// proxy_type is TOX_PROXY_TYPE_NONE.
	ProxyPort uint16

	// The type of savedata to load from.
	SavedataType int

	// The savedata.
	//
	// The data pointed at by this member is owned by the user, so must
	// outlive the options object.
	SavedataData []byte

	// The port to use for the TCP server (relay). If 0, the TCP server is
	// disabled.
	//
	// Enabling it is not required for Tox to function properly.
	//
	// When enabled, your Tox instance can act as a TCP relay for other Tox
	// instance. This leads to increased traffic, thus when writing a client
	// it is recommended to enable TCP server only if the user has an option
	// to disable it.
	TCPPort uint16

	// Enable local network peer discovery.
	//
	// Disabling this will cause Tox to not look for peers on the local network.
	LocalDiscoveryEnabled bool

	// The start port of the inclusive port range to attempt to use.
	//
	// If both start_port and end_port are 0, the default port range will be
	// used: [33445, 33545].
	//
	// If either start_port or end_port is 0 while the other is non-zero, the
	// non-zero port will be the only port in the range.
	//
	// Having start_port > end_port will yield the same behavior as if start_port
	// and end_port were swapped.
	StartPort uint16

	// The end port of the inclusive port range to attempt to use.
	EndPort uint16

	// Enables or disables UDP hole-punching in toxcore. (Default: enabled).
	HolePunchingEnabled bool

	ThreadSafe bool

	// Logging callback for the new tox instance.
	LogCallback func(_ *Tox, level int, file string, line uint32, fname string, msg string)
}

// NewToxOptions allocates a new ToxOptions object and initialises it with the default
// options. This function can be used to preserve long term ABI compatibility by
// giving the responsibility of allocation and deallocation to the Tox library.
//
// Objects returned from this function must be freed using the tox_options_free
// function.
//
// @return A new ToxOptions object with default options or NULL on failure.
func NewToxOptions() *ToxOptions {
	toxopts := C.tox_options_new(nil)
	defer C.tox_options_free(toxopts)

	opts := new(ToxOptions)
	opts.IPv6Enabled = bool(C.tox_options_get_ipv6_enabled(toxopts))
	opts.UDPEnabled = bool(C.tox_options_get_udp_enabled(toxopts))
	opts.ProxyType = int32(C.tox_options_get_proxy_type(toxopts))
	opts.ProxyPort = uint16(C.tox_options_get_proxy_port(toxopts))
	opts.TCPPort = uint16(C.tox_options_get_tcp_port(toxopts))
	opts.LocalDiscoveryEnabled = bool(C.tox_options_get_local_discovery_enabled(toxopts))
	opts.StartPort = uint16(C.tox_options_get_start_port(toxopts))
	opts.EndPort = uint16(C.tox_options_get_end_port(toxopts))
	opts.HolePunchingEnabled = bool(C.tox_options_get_hole_punching_enabled(toxopts))

	return opts
}

func (this *ToxOptions) toCToxOptions() *C.struct_Tox_Options {
	toxopts := C.tox_options_new(nil)
	C.tox_options_default(toxopts)
	C.tox_options_set_ipv6_enabled(toxopts, (C._Bool)(this.IPv6Enabled))
	C.tox_options_set_udp_enabled(toxopts, (C._Bool)(this.UDPEnabled))

	if this.SavedataData != nil {
		C.tox_options_set_savedata_data(toxopts, (*C.uint8_t)(&this.SavedataData[0]), C.size_t(len(this.SavedataData)))
		C.tox_options_set_savedata_type(toxopts, C.TOX_SAVEDATA_TYPE(this.SavedataType))
	}
	C.tox_options_set_tcp_port(toxopts, (C.uint16_t)(this.TCPPort))

	C.tox_options_set_proxy_type(toxopts, C.TOX_PROXY_TYPE(this.ProxyType))
	C.tox_options_set_proxy_port(toxopts, C.uint16_t(this.ProxyPort))
	if len(this.ProxyHost) > 0 {
		C.tox_options_set_proxy_host(toxopts, C.CString(this.ProxyHost))
	}

	C.tox_options_set_local_discovery_enabled(toxopts, C._Bool(this.LocalDiscoveryEnabled))
	C.tox_options_set_start_port(toxopts, C.uint16_t(this.StartPort))
	C.tox_options_set_end_port(toxopts, C.uint16_t(this.EndPort))
	C.tox_options_set_hole_punching_enabled(toxopts, C._Bool(this.HolePunchingEnabled))

	C.tox_options_set_log_callback(toxopts, (*C.tox_log_cb)((unsafe.Pointer)(C.toxCallbackLog)))

	return toxopts
}

//export toxCallbackLog
func toxCallbackLog(ctox *C.Tox, level C.TOX_LOG_LEVEL, file *C.char, line C.uint32_t, fname *C.char, msg *C.char) {
	t := cbUserDatas.get(ctox)
	if t != nil && t.opts != nil && t.opts.LogCallback != nil {
		t.opts.LogCallback(t, int(level), C.GoString(file), uint32(line), C.GoString(fname), C.GoString(msg))
	}
}

type BootNode struct {
	Addr   string
	Port   int
	Pubkey string
}
