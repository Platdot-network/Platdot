package substrate

import (
	"fmt"
	"math/big"
	"strconv"

	log "github.com/ChainSafe/log15"
	"github.com/Platdot-Network/substrate-go/expand"
	"github.com/Platdot-Network/substrate-go/expand/base"
	"github.com/Platdot-Network/substrate-go/models"
	"github.com/Platdot-network/Platdot/chains/chainset"
	"github.com/Platdot-Network/go-substrate-rpc-client/v3/types"
	"github.com/hacpy/go-ethereum/common"
	"github.com/rjman-ljm/go-substrate-crypto/ss58"
	"github.com/rjman-ljm/platdot-utils/msg"
)

const HexPrefix = "hex"

type multiSigTxId uint64
type BlockNumber int64
type OtherSignatories []string

type multiSigTx struct {
	Block BlockNumber
	TxId  multiSigTxId
}

type MultiSigAsMulti struct {
	OriginMsTx     multiSigTx
	Executed       bool
	Threshold      uint16
	Others         []OtherSignatories
	MaybeTimePoint expand.TimePointSafe32
	DestAddress    string
	DestAmount     string
	StoreCall      bool
	MaxWeight      uint64
	DepositNonce   msg.Nonce
	YesVote        []types.AccountID
}

func (l *listener) dealBlockTx(resp *models.BlockResponse, currentBlock int64) {
	for _, e := range resp.Extrinsic {
		if e.Status == "fail" {
			continue
		}

		// Current Extrinsic { Block, Index }
		l.curTx.Block = BlockNumber(currentBlock)
		l.curTx.TxId = multiSigTxId(e.ExtrinsicIndex)
		msTx := MultiSigAsMulti{
			DestAddress: e.MultiSigAsMulti.DestAddress,
			DestAmount:  e.MultiSigAsMulti.DestAmount,
		}
		fromAddressValid := l.checkFromAddress(e)
		toAddressValid := l.checkToAddress(e)

		if e.Type == base.AsMultiNew && fromAddressValid {
			//l.logInfo(FindNewMultiSigTx, currentBlock)
			/// Make a new multiSigTransfer record
			l.makeNewMultiSigRecord(e)
		}

		if e.Type == base.AsMultiApprove && fromAddressValid {
			//l.logInfo(FindApproveMultiSigTx, currentBlock)
			/// Vote(Approve) for the existed MultiSigTransfer record
			l.voteMultiSigRecord(msTx, e)
		}

		if e.Type == base.AsMultiExecuted && fromAddressValid {
			//l.logInfo(FindExecutedMultiSigTx, currentBlock)
			/// Vote and execute the existed MultiSigTransfer record
			l.voteMultiSigRecord(msTx, e)
			l.executeMultiSigRecord(msTx)
		}

		if e.Type == base.UtilityBatch && toAddressValid {
			//l.logInfo(FindBatchMultiSigTx, currentBlock)
			if l.findLostTxByAddress(currentBlock, e) {
				continue
			}

			sendAmount, ok := l.getSendAmount(e)
			/// if `chainId wrong`, `amount is negative` or `not cross-chain tx`
			if !ok || e.Recipient == "" {
				continue
			}

			destId, rId, recipient, err := l.parseRemark(e.Recipient)
			if err != nil {
				l.log.Error("parse remark error", "err", err)
				continue
			}
			if l.checkRemark(destId, rId, recipient) {
				depositNonce, _ := strconv.ParseInt(strconv.FormatInt(currentBlock, 10)+strconv.FormatInt(int64(e.ExtrinsicIndex), 10), 10, 64)

				m := msg.NewMultiSigTransfer(
					l.chainId,
					destId,
					msg.Nonce(depositNonce),
					sendAmount,
					rId,
					recipient[:],
				)
				l.logReadyToSend(sendAmount, recipient, e)
				l.submitMessage(m, nil)
			}
		}
	}
}

func (l *listener) parseRemark(res string) (msg.ChainId, msg.ResourceId, []byte, error) {
	offset1 := -1
	offset2 := -1
	for i, v := range res {
		if v == ',' && offset1 < 0 {
			offset1 = i
			continue
		}
		if v == ',' && offset1 > 0 {
			offset2 = i
			break
		}
	}
	if offset1 < 0 || offset2 < 0 {
		return msg.ChainId(0), msg.ResourceId{}, nil, fmt.Errorf("remark value err, didn't parse out destId, resourceId and address")
	}

	dest := res[:offset1]
	rId := chainset.ResourceIdPrefix + res[offset1+1:offset2]
	address := res[offset2+1:]

	destId, err := strconv.ParseInt(dest, 10, 64)
	if err != nil {
		fmt.Printf("parse remark_destId to int err, value is %v\n", destId)
		return msg.ChainId(0), msg.ResourceId{}, nil, fmt.Errorf("parse remark_destId to int err")
	}

	recipient := []byte(address)

	if address[:3] == HexPrefix {
		recipientAccount := types.NewAccountID(common.FromHex(address[3:]))
		recipient = recipientAccount[:]
	}

	return msg.ChainId(destId), l.chainCore.ConvertStringToResourceId(rId), recipient, nil
}
func (l *listener) checkRemark(destId msg.ChainId, rId msg.ResourceId, recipient []byte) bool {
	alayaPass := l.chainId == 1 && destId == 2 && rId == l.chainCore.ConvertStringToResourceId(chainset.ResourceIdAKSM)
	platonPass := l.chainId == 3 && destId == 4 && rId == l.chainCore.ConvertStringToResourceId(chainset.ResourceIdPDOT)

	if alayaPass || platonPass {
		log.Info("Parameter check passed", "Dest", destId, "RId", rId.Shorten(), "Recipient", l.chainCore.ShortenAddress(string(recipient)))
	}
	return alayaPass || platonPass
}

func (l *listener) findLostTxByAddress(currentBlock int64, e *models.ExtrinsicResponse) bool {
	sendPubAddress, _ := ss58.DecodeToPub(e.FromAddress)
	lostPubAddress, _ := ss58.DecodeToPub(l.lostAddress)

	if l.lostAddress != "" {
		/// Find the lost transaction
		if string(sendPubAddress) == string(lostPubAddress[:]) {
			l.logInfo(FindLostMultiSigTx, currentBlock)
		}
		return true
	} else {
		return false
	}
}

func (l *listener) getSendAmount(e *models.ExtrinsicResponse) (*big.Int, bool) {
	// Construct parameters of message
	amount, ok := big.NewInt(0).SetString(e.Amount, 10)
	if !ok || amount.Uint64() == 0 {
		fmt.Printf("parse transfer amount %v, amount.string %v\n", amount, amount.String())
		return nil, false
	}

	sendAmount, err := l.chainCore.GetAmountToEth(amount.Bytes(), e.AssetId)
	if err != nil {
		return nil, false
	}

	return sendAmount, true
}

func (l *listener) checkToAddress(e *models.ExtrinsicResponse) bool {
	/// Validate whether a cross-chain transaction
	toPubAddress, _ := ss58.DecodeToPub(e.ToAddress)
	toAddress := types.NewAddressFromAccountID(toPubAddress)
	return toAddress.AsAccountID == l.multiSigAddr
}

func (l *listener) checkFromAddress(e *models.ExtrinsicResponse) bool {
	fromPubAddress, _ := ss58.DecodeToPub(e.FromAddress)
	fromAddress := types.NewAddressFromAccountID(fromPubAddress)
	currentRelayer := types.NewAddressFromAccountID(l.relayer.kr.PublicKey)
	if currentRelayer.AsAccountID == fromAddress.AsAccountID {
		return true
	}
	for _, r := range l.relayer.otherSignatories {
		if types.AccountID(r) == fromAddress.AsAccountID {
			return true
		}
	}
	return false
}

func (l *listener) makeNewMultiSigRecord(e *models.ExtrinsicResponse) {
	msTx := MultiSigAsMulti{
		Executed:       false,
		Threshold:      e.MultiSigAsMulti.Threshold,
		MaybeTimePoint: e.MultiSigAsMulti.MaybeTimePoint,
		DestAddress:    e.MultiSigAsMulti.DestAddress,
		DestAmount:     e.MultiSigAsMulti.DestAmount,
		Others:         nil,
		StoreCall:      e.MultiSigAsMulti.StoreCall,
		MaxWeight:      e.MultiSigAsMulti.MaxWeight,
		OriginMsTx:     l.curTx,
	}
	/// Mark voted
	msTx.Others = append(msTx.Others, e.MultiSigAsMulti.OtherSignatories)
	l.asMulti[l.curTx] = msTx
}

func (l *listener) voteMultiSigRecord(msTx MultiSigAsMulti, e *models.ExtrinsicResponse) {
	for k, ms := range l.asMulti {
		if !ms.Executed && ms.DestAddress == msTx.DestAddress && ms.DestAmount == msTx.DestAmount {
			//l.log.Info("relayer succeed vote", "Address", e.FromAddress)
			voteMsTx := l.asMulti[k]
			voteMsTx.Others = append(voteMsTx.Others, e.MultiSigAsMulti.OtherSignatories)
			l.asMulti[k] = voteMsTx
		}
	}
}

func (l *listener) executeMultiSigRecord(msTx MultiSigAsMulti) {
	for k, ms := range l.asMulti {
		if !ms.Executed && ms.DestAddress == msTx.DestAddress && ms.DestAmount == msTx.DestAmount {
			exeMsTx := l.asMulti[k]
			exeMsTx.Executed = true
			l.asMulti[k] = exeMsTx
		}
	}
}
