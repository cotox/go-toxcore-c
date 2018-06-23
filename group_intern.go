package tox

/*
#include <stdint.h>
#include <tox/tox.h>

extern uint8_t *tox_conference_get_pubkey(Tox *tox, uint32_t conference_number, void *pkbuf);
extern uint8_t *tox_conference_get_identifier(Tox *tox, uint32_t conference_number, void *idbuf);
*/
import "C"
import (
	"encoding/hex"
	"strings"
	"unsafe"
)

func (t *Tox) ConferenceGetPubkey(groupNumber uint32) (string, error) {
	pkbuf := [PublicKeySize]byte{}
	C.tox_conference_get_pubkey(t.toxcore, C.uint32_t(groupNumber), (unsafe.Pointer)(&pkbuf[0]))
	pubkey := strings.ToUpper(hex.EncodeToString(pkbuf[:]))

	return pubkey, nil
}

func (t *Tox) ConferenceGetIdentifier(groupNumber uint32) (string, error) {
	idbuf := [1 + PublicKeySize]byte{}
	C.tox_conference_get_identifier(t.toxcore, C.uint32_t(groupNumber), (unsafe.Pointer)(&idbuf[0]))
	identifier := strings.ToUpper(hex.EncodeToString(idbuf[:]))
	identifier = identifier[2:] // 1B(type)+32B(identifier)

	return identifier, nil
}
