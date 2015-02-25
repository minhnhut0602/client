package kex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/keybase/client/go/libkb"
	jsonw "github.com/keybase/go-jsonw"
	"github.com/ugorji/go/codec"
)

const (
	startkexMsg    = "startkex"
	startrevkexMsg = "startrevkex"
	helloMsg       = "hello"
	pleasesignMsg  = "pleasesign"
	doneMsg        = "done"
)

var ErrMACMismatch = errors.New("Computed HMAC doesn't match message HMAC")

var GlobalTimeout = 5 * time.Minute

var G = &libkb.G

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) StartKexSession(ctx *Context, id StrongID) error {
	mb := &Body{Name: startkexMsg, Args: MsgArgs{StrongID: id}}
	return s.post(ctx, mb)
}

func (s *Sender) StartReverseKexSession(ctx *Context) error {
	return nil
}

func (s *Sender) Hello(ctx *Context, devID libkb.DeviceID, devKeyID libkb.KID) error {
	mb := &Body{Name: helloMsg, Args: MsgArgs{DeviceID: devID, DevKeyID: devKeyID}}
	return s.post(ctx, mb)
}

func (s *Sender) PleaseSign(ctx *Context, eddsa libkb.NaclSigningKeyPublic, sig, devType, devDesc string) error {
	mb := &Body{Name: pleasesignMsg, Args: MsgArgs{SigningKey: eddsa, Sig: sig, DevType: devType, DevDesc: devDesc}}
	return s.post(ctx, mb)
}

func (s *Sender) Done(ctx *Context, mt libkb.MerkleTriple) error {
	mb := &Body{Name: doneMsg, Args: MsgArgs{MerkleTriple: mt}}
	return s.post(ctx, mb)
}

// XXX get rid of this when real client comm works
func (s *Sender) RegisterTestDevice(srv Handler, device libkb.DeviceID) error {
	return nil
}

func (s *Sender) post(ctx *Context, body *Body) error {
	msg := NewMsg(ctx, body)
	msg.Direction = 1
	mac, err := msg.MacSum()
	if err != nil {
		return err
	}
	msg.Mac = mac

	menc, err := msg.Body.Encode()
	if err != nil {
		return err
	}

	G.Log.Info("G.API: %+v", G.API)

	_, err = G.API.Post(libkb.ApiArg{
		Endpoint:    "kex/send",
		NeedSession: true,
		Args: libkb.HttpArgs{
			"dir":      libkb.I{Val: msg.Direction},
			"I":        libkb.S{Val: hex.EncodeToString(msg.StrongID[:])},
			"msg":      libkb.S{Val: menc},
			"receiver": libkb.S{Val: msg.Dst.String()},
			"sender":   libkb.S{Val: msg.Src.String()},
			"seqno":    libkb.I{Val: msg.Seqno},
			"w":        libkb.S{Val: hex.EncodeToString(msg.WeakID[:])},
		},
	})
	return err
}

type KexReceiver struct {
	handler Handler
	seqno   int
	pollDur time.Duration
}

func NewKexReceiver(handler Handler) *KexReceiver {
	return &KexReceiver{handler: handler, pollDur: 20 * time.Second}
}

func (r *KexReceiver) Receive(ctx *Context) error {
	msgs, err := r.get(ctx)
	if err != nil {
		return err
	}
	for _, m := range msgs {
		switch m.Name {
		case startkexMsg:
			return r.handler.StartKexSession(ctx, m.Args.StrongID)
		case startrevkexMsg:
			return r.handler.StartReverseKexSession(ctx)
		case helloMsg:
			return r.handler.Hello(ctx, m.Args.DeviceID, m.Args.DevKeyID)
		case pleasesignMsg:
			return r.handler.PleaseSign(ctx, m.Args.SigningKey, m.Args.Sig, m.Args.DevType, m.Args.DevDesc)
		case doneMsg:
			return r.handler.Done(ctx, m.Args.MerkleTriple)
		default:
			return fmt.Errorf("unhandled message name: %q", m.Name)
		}
	}
	return nil
}

func (r *KexReceiver) ReceiveFilter(name string) error {
	return nil
}

func (r *KexReceiver) get(ctx *Context) ([]*Msg, error) {
	res, err := G.API.Get(libkb.ApiArg{
		Endpoint:    "kex/receive",
		NeedSession: true,
		Args: libkb.HttpArgs{
			"w":    libkb.S{Val: hex.EncodeToString(ctx.WeakID[:])},
			"dir":  libkb.I{Val: 1},
			"low":  libkb.I{Val: 0},
			"poll": libkb.I{Val: int(r.pollDur / time.Second)},
		},
	})
	if err != nil {
		return nil, err
	}

	msgs := res.Body.AtKey("msgs")
	n, err := msgs.Len()
	if err != nil {
		return nil, err
	}
	messages := make([]*Msg, n)
	for i := 0; i < n; i++ {
		messages[i], err = MsgImport(msgs.AtIndex(i))
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

type Msg struct {
	Meta
	Body
}

func NewMsg(ctx *Context, body *Body) *Msg {
	return &Msg{
		Meta: ctx.Meta,
		Body: *body,
	}
}

func (m *Msg) CheckMAC() (bool, error) {
	sum, err := m.MacSum()
	if err != nil {
		return false, err
	}
	return hmac.Equal(sum, m.Mac), nil
}

func (m *Msg) MacSum() ([]byte, error) {
	t := m.Mac
	defer func() { m.Mac = t }()
	m.Mac = nil
	var buf bytes.Buffer
	var h codec.MsgpackHandle
	if err := codec.NewEncoder(&buf, &h).Encode(m); err != nil {
		return nil, err
	}
	return m.mac(buf.Bytes(), m.StrongID[:]), nil
}

func (m *Msg) mac(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func deviceID(w *jsonw.Wrapper) (libkb.DeviceID, error) {
	s, err := w.GetString()
	if err != nil {
		return libkb.DeviceID{}, err
	}
	d, err := libkb.ImportDeviceID(s)
	if err != nil {
		return libkb.DeviceID{}, err
	}
	return *d, nil
}

func MsgImport(w *jsonw.Wrapper) (*Msg, error) {
	r := &Msg{}
	u, err := libkb.GetUid(w.AtKey("uid"))
	if err != nil {
		return nil, err
	}
	r.UID = *u

	r.Src, err = deviceID(w.AtKey("sender"))
	if err != nil {
		return nil, err
	}
	r.Dst, err = deviceID(w.AtKey("receiver"))
	if err != nil {
		return nil, err
	}

	r.Seqno, err = w.AtKey("seqno").GetInt()
	if err != nil {
		return nil, err
	}
	r.Direction, err = w.AtKey("dir").GetInt()
	if err != nil {
		return nil, err
	}

	stID, err := w.AtKey("I").GetString()
	if err != nil {
		return nil, err
	}
	bstID, err := hex.DecodeString(stID)
	if err != nil {
		return nil, err
	}
	copy(r.StrongID[:], bstID)

	wkID, err := w.AtKey("w").GetString()
	if err != nil {
		return nil, err
	}
	bwkID, err := hex.DecodeString(wkID)
	if err != nil {
		return nil, err
	}
	copy(r.WeakID[:], bwkID)

	body, err := w.AtKey("msg").GetString()
	if err != nil {
		return nil, err
	}
	mb, err := BodyDecode(body)
	if err != nil {
		return nil, err
	}
	r.Body = *mb

	ok, err := r.CheckMAC()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrMACMismatch
	}

	return r, nil
}

// MsgArgs has optional fields in it, but there aren't that many,
// so just using the same struct for all msgs for simplicity.
type MsgArgs struct {
	StrongID     StrongID
	DeviceID     libkb.DeviceID
	DevKeyID     libkb.KID
	SigningKey   libkb.NaclSigningKeyPublic
	Sig          string
	DevType      string
	DevDesc      string
	MerkleTriple libkb.MerkleTriple
}

type Body struct {
	Name string
	Args MsgArgs
	Mac  []byte
}

func BodyDecode(data string) (*Body, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var h codec.MsgpackHandle
	var k Body
	err = codec.NewDecoderBytes(bytes, &h).Decode(&k)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (k *Body) Encode() (string, error) {
	var buf bytes.Buffer
	var h codec.MsgpackHandle
	err := codec.NewEncoder(&buf, &h).Encode(k)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
