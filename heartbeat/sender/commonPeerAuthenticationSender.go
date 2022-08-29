package sender

import (
	"bytes"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/batch"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/heartbeat"
)

type commonPeerAuthenticationSender struct {
	baseSender
	nodesCoordinator      heartbeat.NodesCoordinator
	peerSignatureHandler  crypto.PeerSignatureHandler
	hardforkTrigger       heartbeat.HardforkTrigger
	hardforkTriggerPubKey []byte
}

func (cpas *commonPeerAuthenticationSender) generateMessageBytes(
	pkBytes []byte,
	privateKey crypto.PrivateKey,
	p2pSkBytes []byte,
	pidBytes []byte,
) ([]byte, bool, error) {

	msg := &heartbeat.PeerAuthentication{
		Pid:    pidBytes,
		Pubkey: pkBytes,
	}

	hardforkPayload, isTriggered := cpas.getHardforkPayload()
	payload := &heartbeat.Payload{
		Timestamp:       time.Now().Unix(),
		HardforkMessage: string(hardforkPayload),
	}
	payloadBytes, err := cpas.marshaller.Marshal(payload)
	if err != nil {
		return nil, isTriggered, err
	}
	msg.Payload = payloadBytes

	if p2pSkBytes != nil {
		msg.PayloadSignature, err = cpas.messenger.SignWithPrivateKey(p2pSkBytes, payloadBytes)
		if err != nil {
			return nil, isTriggered, err
		}
	} else {
		msg.PayloadSignature, err = cpas.messenger.Sign(payloadBytes)
		if err != nil {
			return nil, isTriggered, err
		}
	}

	msg.Signature, err = cpas.peerSignatureHandler.GetPeerSignature(privateKey, msg.Pid)
	if err != nil {
		return nil, isTriggered, err
	}

	msgBytes, err := cpas.marshaller.Marshal(msg)
	if err != nil {
		return nil, isTriggered, err
	}

	b := &batch.Batch{
		Data: make([][]byte, 1),
	}
	b.Data[0] = msgBytes
	data, err := cpas.marshaller.Marshal(b)
	if err != nil {
		return nil, isTriggered, err
	}

	return data, isTriggered, nil
}

func (cpas *commonPeerAuthenticationSender) isValidator(pkBytes []byte) bool {
	_, _, err := cpas.nodesCoordinator.GetValidatorWithPublicKey(pkBytes)
	return err == nil
}

func (cpas *commonPeerAuthenticationSender) isHardforkSource(pkBytes []byte) bool {
	return bytes.Equal(pkBytes, cpas.hardforkTriggerPubKey)
}

func (cpas *commonPeerAuthenticationSender) getHardforkPayload() ([]byte, bool) {
	payload := make([]byte, 0)
	_, isTriggered := cpas.hardforkTrigger.RecordedTriggerMessage()
	if isTriggered {
		payload = cpas.hardforkTrigger.CreateData()
	}

	return payload, isTriggered
}
