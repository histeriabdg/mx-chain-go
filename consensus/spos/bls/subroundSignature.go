package bls

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/consensus"
	"github.com/ElrondNetwork/elrond-go/consensus/spos"
)

type subroundSignature struct {
	*spos.Subround

	appStatusHandler core.AppStatusHandler
}

// NewSubroundSignature creates a subroundSignature object
func NewSubroundSignature(
	baseSubround *spos.Subround,
	extend func(subroundId int),
	appStatusHandler core.AppStatusHandler,
) (*subroundSignature, error) {
	err := checkNewSubroundSignatureParams(
		baseSubround,
	)
	if err != nil {
		return nil, err
	}
	if check.IfNil(appStatusHandler) {
		return nil, spos.ErrNilAppStatusHandler
	}

	srSignature := subroundSignature{
		Subround:         baseSubround,
		appStatusHandler: appStatusHandler,
	}
	srSignature.Job = srSignature.doSignatureJob
	srSignature.Check = srSignature.doSignatureConsensusCheck
	srSignature.Extend = extend

	return &srSignature, nil
}

func checkNewSubroundSignatureParams(
	baseSubround *spos.Subround,
) error {
	if baseSubround == nil {
		return spos.ErrNilSubround
	}
	if baseSubround.ConsensusState == nil {
		return spos.ErrNilConsensusState
	}

	err := spos.ValidateConsensusCore(baseSubround.ConsensusCoreHandler)

	return err
}

// doSignatureJob method does the job of the subround Signature
func (sr *subroundSignature) doSignatureJob(_ context.Context) bool {
	if !sr.IsNodeInConsensusGroup(sr.SelfPubKey()) && !sr.IsMultiKeyInConsensusGroup() {
		return true
	}
	if !sr.CanDoSubroundJob(sr.Current()) {
		return false
	}
	if check.IfNil(sr.Header) {
		log.Error("doSignatureJob", "error", spos.ErrNilHeader)
		return false
	}

	selfIndex, err := sr.SelfConsensusGroupIndex()
	if err != nil {
		log.Debug("doSignatureJob.SelfConsensusGroupIndex: not in consensus group")
		return false
	}

	signatureShare, err := sr.SignatureHandler().CreateSignatureShare(sr.GetData(), uint16(selfIndex), sr.Header.GetEpoch())
	if err != nil {
		log.Debug("doSignatureJob.CreateSignatureShare", "error", err.Error())
		return false
	}
	isSelfLeader := sr.IsSelfLeaderInCurrentRound()

	if isSelfLeader || sr.IsNodeInConsensusGroup(sr.SelfPubKey()) {
		signatureShare, err := sr.MultiSigner().CreateSignatureShare(sr.GetData(), nil)
		if err != nil {
			log.Debug("doSignatureJob.CreateSignatureShare", "error", err.Error())
			return false
		}

		if !isSelfLeader {
			ok := sr.createAndSendSignatureMessage(signatureShare, []byte(sr.SelfPubKey()))
			if !ok {
				return false
			}
		}

		ok := sr.completeSignatureSubRound(sr.SelfPubKey(), isSelfLeader)
		if !ok {
			return false
		}
	}

	return sr.doSignatureJobForManagedKeys()
}

func (sr *subroundSignature) createAndSendSignatureMessage(signatureShare []byte, pkBytes []byte) bool {
	// TODO: Analyze it is possible to send message only to leader with O(1) instead of O(n)
	cnsMsg := consensus.NewConsensusMessage(
		sr.GetData(),
		signatureShare,
		nil,
		nil,
		pkBytes,
		nil,
		int(MtSignature),
		sr.RoundHandler().Index(),
		sr.ChainID(),
		nil,
		nil,
		nil,
		sr.GetAssociatedPid(pkBytes),
	)

	err := sr.BroadcastMessenger().BroadcastConsensusMessage(cnsMsg)
	if err != nil {
		log.Debug("createAndSendSignatureMessage.BroadcastConsensusMessage",
			"error", err.Error(), "pk", pkBytes)
		return false
	}

	log.Debug("step 2: signature has been sent", "pk", pkBytes)

	return true
}

func (sr *subroundSignature) completeSignatureSubRound(pk string, shouldWaitForAllSigsAsync bool) bool {
	err := sr.SetJobDone(pk, sr.Current(), true)
	if err != nil {
		log.Debug("doSignatureJob.SetSelfJobDone",
			"subround", sr.Name(),
			"error", err.Error(),
			"pk", []byte(pk),
		)
		return false
	}

	if shouldWaitForAllSigsAsync {
		go sr.waitAllSignatures()
	}

	return true
}

// receivedSignature method is called when a signature is received through the signature channel.
// If the signature is valid, than the jobDone map corresponding to the node which sent it,
// is set on true for the subround Signature
func (sr *subroundSignature) receivedSignature(_ context.Context, cnsDta *consensus.Message) bool {
	node := string(cnsDta.PubKey)
	pkForLogs := core.GetTrimmedPk(hex.EncodeToString(cnsDta.PubKey))

	if !sr.IsConsensusDataSet() {
		return false
	}

	if !sr.IsNodeInConsensusGroup(node) {
		sr.PeerHonestyHandler().ChangeScore(
			node,
			spos.GetConsensusTopicID(sr.ShardCoordinator()),
			spos.ValidatorPeerHonestyDecreaseFactor,
		)

		return false
	}

	if !sr.IsSelfLeaderInCurrentRound() && !sr.IsMultiKeyLeaderInCurrentRound() {
		return false
	}

	if !sr.IsConsensusDataEqual(cnsDta.BlockHeaderHash) {
		return false
	}

	if !sr.CanProcessReceivedMessage(cnsDta, sr.RoundHandler().Index(), sr.Current()) {
		return false
	}

	index, err := sr.ConsensusGroupIndex(node)
	if err != nil {
		log.Debug("receivedSignature.ConsensusGroupIndex",
			"node", pkForLogs,
			"error", err.Error())
		return false
	}

	if check.IfNil(sr.Header) {
		log.Error("receivedSignature", "error", spos.ErrNilHeader)
		return false
	}

	err = sr.SignatureHandler().VerifySignatureShare(uint16(index), cnsDta.SignatureShare, sr.GetData(), sr.Header.GetEpoch())
	if err != nil {
		log.Debug("receivedSignature.VerifySignatureShare",
			"node", pkForLogs,
			"index", index,
			"error", err.Error())
		return false
	}

	err = sr.SignatureHandler().StoreSignatureShare(uint16(index), cnsDta.SignatureShare)
	if err != nil {
		log.Debug("receivedSignature.StoreSignatureShare",
			"node", pkForLogs,
			"index", index,
			"error", err.Error())
		return false
	}

	err = sr.SetJobDone(node, sr.Current(), true)
	if err != nil {
		log.Debug("receivedSignature.SetJobDone",
			"node", pkForLogs,
			"subround", sr.Name(),
			"error", err.Error())
		return false
	}

	sr.PeerHonestyHandler().ChangeScore(
		node,
		spos.GetConsensusTopicID(sr.ShardCoordinator()),
		spos.ValidatorPeerHonestyIncreaseFactor,
	)

	sr.appStatusHandler.SetStringValue(common.MetricConsensusRoundState, "signed")
	return true
}

// doSignatureConsensusCheck method checks if the consensus in the subround Signature is achieved
func (sr *subroundSignature) doSignatureConsensusCheck() bool {
	if sr.RoundCanceled {
		return false
	}

	if sr.IsSubroundFinished(sr.Current()) {
		sr.appStatusHandler.SetStringValue(common.MetricConsensusRoundState, "signed")

		return true
	}

	isSelfLeader := sr.IsSelfLeaderInCurrentRound() || sr.IsMultiKeyLeaderInCurrentRound()
	isSelfInConsensusGroup := sr.IsNodeInConsensusGroup(sr.SelfPubKey()) || sr.IsMultiKeyInConsensusGroup()

	threshold := sr.Threshold(sr.Current())
	if sr.FallbackHeaderValidator().ShouldApplyFallbackValidation(sr.Header) {
		threshold = sr.FallbackThreshold(sr.Current())
		log.Warn("subroundSignature.doSignatureConsensusCheck: fallback validation has been applied",
			"minimum number of signatures required", threshold,
			"actual number of signatures received", sr.getNumOfSignaturesCollected(),
		)
	}

	areSignaturesCollected, numSigs := sr.areSignaturesCollected(threshold)
	areAllSignaturesCollected := numSigs == sr.ConsensusGroupSize()

	isJobDoneByLeader := isSelfLeader && (areAllSignaturesCollected || (areSignaturesCollected && sr.WaitingAllSignaturesTimeOut))

	selfJobDone := true
	if sr.IsNodeInConsensusGroup(sr.SelfPubKey()) {
		selfJobDone = sr.IsSelfJobDone(sr.Current())
	}
	multiKeyJobDone := true
	if sr.IsMultiKeyInConsensusGroup() {
		multiKeyJobDone = sr.IsMultiKeyJobDone(sr.Current())
	}
	isJobDoneByConsensusNode := !isSelfLeader && isSelfInConsensusGroup && selfJobDone && multiKeyJobDone

	isSubroundFinished := !isSelfInConsensusGroup || isJobDoneByConsensusNode || isJobDoneByLeader

	if isSubroundFinished {
		if isSelfLeader {
			log.Debug("step 2: signatures",
				"received", numSigs,
				"total", len(sr.ConsensusGroup()))
		}

		log.Debug("step 2: subround has been finished",
			"subround", sr.Name())
		sr.SetStatus(sr.Current(), spos.SsFinished)

		sr.appStatusHandler.SetStringValue(common.MetricConsensusRoundState, "signed")

		return true
	}

	return false
}

// areSignaturesCollected method checks if the signatures received from the nodes, belonging to the current
// jobDone group, are more than the necessary given threshold
func (sr *subroundSignature) areSignaturesCollected(threshold int) (bool, int) {
	n := sr.getNumOfSignaturesCollected()
	return n >= threshold, n
}

func (sr *subroundSignature) getNumOfSignaturesCollected() int {
	n := 0

	for i := 0; i < len(sr.ConsensusGroup()); i++ {
		node := sr.ConsensusGroup()[i]

		isSignJobDone, err := sr.JobDone(node, sr.Current())
		if err != nil {
			log.Debug("getNumOfSignaturesCollected.JobDone",
				"node", node,
				"subround", sr.Name(),
				"error", err.Error())
			continue
		}

		if isSignJobDone {
			n++
		}
	}

	return n
}

func (sr *subroundSignature) waitAllSignatures() {
	remainingTime := sr.remainingTime()
	time.Sleep(remainingTime)

	if sr.IsSubroundFinished(sr.Current()) {
		return
	}

	sr.WaitingAllSignaturesTimeOut = true

	select {
	case sr.ConsensusChannel() <- true:
	default:
	}
}

func (sr *subroundSignature) remainingTime() time.Duration {
	startTime := sr.RoundHandler().TimeStamp()
	maxTime := time.Duration(float64(sr.StartTime()) + float64(sr.EndTime()-sr.StartTime())*waitingAllSigsMaxTimeThreshold)
	remainigTime := sr.RoundHandler().RemainingTime(startTime, maxTime)

	return remainigTime
}

func (sr *subroundSignature) doSignatureJobForManagedKeys() bool {
	isMultiKeyLeader := sr.IsMultiKeyLeaderInCurrentRound()

	numMultiKeysSignaturesSent := 0
	for idx, pk := range sr.ConsensusGroup() {
		pkBytes := []byte(pk)
		if sr.IsJobDone(pk, sr.Current()) {
			continue
		}
		if !sr.IsKeyManagedByCurrentNode(pkBytes) {
			continue
		}

		managedPrivateKey := sr.GetMessageSigningPrivateKey(pkBytes)
		signatureShare, err := sr.MultiSigner().CreateAndAddSignatureShareForKey(sr.GetData(), managedPrivateKey, pkBytes)
		if err != nil {
			log.Debug("doSignatureJobForManagedKeys.CreateAndAddSignatureShareForKey", "error", err.Error())
			return false
		}

		if !isMultiKeyLeader {
			ok := sr.createAndSendSignatureMessage(signatureShare, pkBytes)
			if !ok {
				return false
			}

			numMultiKeysSignaturesSent++
		}

		isLeader := idx == spos.IndexOfLeaderInConsensusGroup
		ok := sr.completeSignatureSubRound(pk, isLeader)
		if !ok {
			return false
		}
	}

	if numMultiKeysSignaturesSent > 0 {
		log.Debug("step 2: multi keys signatures have been sent", "num", numMultiKeysSignaturesSent)
	}

	return true
}
