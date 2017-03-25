package GSSAPI

import (
	"errors"
	"fmt"
	"github.com/jcmturner/asn1"
)

/*
https://msdn.microsoft.com/en-us/library/ms995330.aspx

NegotiationToken ::= CHOICE {
  negTokenInit    [0] NegTokenInit,  This is the Negotiation token sent from the client to the server.
  negTokenResp    [1] NegTokenResp
}

NegTokenInit ::= SEQUENCE {
  mechTypes       [0] MechTypeList,
  reqFlags        [1] ContextFlags  OPTIONAL,
  -- inherited from RFC 2478 for backward compatibility,
  -- RECOMMENDED to be left out
  mechToken       [2] OCTET STRING  OPTIONAL,
  mechListMIC     [3] OCTET STRING  OPTIONAL,
  ...
}

NegTokenResp ::= SEQUENCE {
  negState       [0] ENUMERATED {
    accept-completed    (0),
    accept-incomplete   (1),
    reject              (2),
    request-mic         (3)
  }                                 OPTIONAL,
  -- REQUIRED in the first reply from the target
  supportedMech   [1] MechType      OPTIONAL,
  -- present only in the first reply from the target
  responseToken   [2] OCTET STRING  OPTIONAL,
  mechListMIC     [3] OCTET STRING  OPTIONAL,
  ...
}
*/

// Tag attribute of NegotiationToken will indicate type:
// 0xa0 (160) - negTokenInit
// 0xa1 (161) - negTokenResp
type NegotiationToken asn1.RawValue

type NegTokenInit struct {
	MechTypes    MechTypeList `asn1:"explicit,tag:0"`
	ReqFlags     ContextFlags `asn1:"explicit,optional,tag:1"`
	MechToken    []byte       `asn1:"explicit,optional,tag:2"`
	MechTokenMIC []byte       `asn1:"explicit,optional,tag:3"`
}

type NegTokenResp struct {
	NegState      asn1.Enumerated `asn1:"explicit,optional,tag:0"`
	SupportedMech MechType        `asn1:"explicit,optional,tag:1"`
	ResponseToken []byte          `asn1:"explicit,optional,tag:2"`
	MechListMIC   []byte          `asn1:"explicit,optional,tag:3"`
}

// Unmarshal and return either a NegTokenInit or a NegTokenResp.
//
// The boolean indicates if the reponse is a NegTokenInit.
// If error is nil and the boolean is false the response is a NegTokenResp.
func (n *NegotiationToken) Unmarshal(b []byte) (bool, interface{}, error) {
	_, err := asn1.Unmarshal(b, n)
	if err != nil {
		return false, nil, fmt.Errorf("Error unmarshalling NegotiationToken: %v", err)
	}
	var negToken interface{}
	var isInit bool
	switch n.Tag {
	case 0:
		negToken = &NegTokenInit{}
		isInit = true
	case 1:
		negToken = &NegTokenResp{}
	default:
		return false, nil, errors.New("Unknown choice type for NegotiationToken")
	}
	_, err = asn1.Unmarshal(n.Bytes, negToken)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling NegotiationToken type %d: %v", n.Tag, err)
	}
	return isInit, negToken, nil
}

// Returns marshalled bytes of a NegotiationToken rather than the NegTokenInit
func (n *NegTokenInit) Marshal() ([]byte, error) {
	b, err := asn1.Marshal(*n)
	if err != nil {
		return nil, err
	}
	nt := NegotiationToken{
		Tag:        0,
		Class:      2,
		IsCompound: true,
		Bytes:      b,
	}
	nb, err := asn1.Marshal(nt)
	if err != nil {
		return nil, err
	}
	return nb, nil
}

// Returns marshalled bytes of a NegotiationToken rather than the NegTokenResp
func (n *NegTokenResp) Marshal() ([]byte, error) {
	b, err := asn1.Marshal(*n)
	if err != nil {
		return nil, err
	}
	nt := NegotiationToken{
		Tag:        1,
		Class:      2,
		IsCompound: true,
		Bytes:      b,
	}
	nb, err := asn1.Marshal(nt)
	if err != nil {
		return nil, err
	}
	return nb, nil
}