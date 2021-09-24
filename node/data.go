package node

import (
	stdJSON "encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	KindActivation          = "activate_account"
	KindBallot              = "ballot"
	KindDelegation          = "delegation"
	KindDoubleBaking        = "double_baking_evidence"
	KindDoubleEndorsing     = "double_endorsement_evidence"
	KindEndorsement         = "endorsement"
	KindEndorsementWithSlot = "endorsement_with_slot"
	KindOrigination         = "origination"
	KindProposal            = "proposals"
	KindReveal              = "reveal"
	KindNonceRevelation     = "seed_nonce_revelation"
	KindTransaction         = "transaction"
)

// Errors
var (
	ErrUnknownKind = errors.New("Unknown operation kind")
)

// MempoolResponse -
type MempoolResponse struct {
	Applied       []Applied `json:"applied"`
	Refused       []Failed  `json:"refused"`
	BranchRefused []Failed  `json:"branch_refused"`
	BranchDelayed []Failed  `json:"branch_delayed"`
}

// Applied -
type Applied struct {
	Hash      string             `json:"hash"`
	Branch    string             `json:"branch"`
	Signature string             `json:"signature"`
	Contents  []Content          `json:"contents"`
	Raw       stdJSON.RawMessage `json:"raw"`
}

// UnmarshalJSON -
func (a *Applied) UnmarshalJSON(data []byte) error {
	type buf Applied
	if err := json.Unmarshal(data, (*buf)(a)); err != nil {
		return err
	}
	a.Raw = data
	return nil
}

// Failed -
type Failed struct {
	Hash      string
	Protocol  string             `json:"protocol"`
	Branch    string             `json:"branch"`
	Contents  []Content          `json:"contents"`
	Signature string             `json:"signature"`
	Error     stdJSON.RawMessage `json:"error,omitempty"`
	Raw       stdJSON.RawMessage `json:"raw"`
}

// UnmarshalJSON -
func (f *Failed) UnmarshalJSON(data []byte) error {
	var body []stdJSON.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return errors.Errorf("Invalid failed operation body %s", string(data))
	}
	if err := json.Unmarshal(body[0], &f.Hash); err != nil {
		return err
	}
	type buf Failed
	if err := json.Unmarshal(body[1], (*buf)(f)); err != nil {
		return err
	}
	f.Raw = data
	return nil
}

// Contents -
type Content struct {
	Kind string             `json:"kind"`
	Body stdJSON.RawMessage `json:"-"`
}

// UnmarshalJSON -
func (c *Content) UnmarshalJSON(data []byte) error {
	type buf Content
	if err := json.Unmarshal(data, (*buf)(c)); err != nil {
		return err
	}
	c.Body = data
	return nil
}

// Constants -
type Constants struct {
	ProofOfWorkNonceSize         int64            `json:"proof_of_work_nonce_size"`
	NonceLength                  int64            `json:"nonce_length"`
	MaxAnonOpsPerBlock           int64            `json:"max_anon_ops_per_block"`
	MaxOperationDataLength       int64            `json:"max_operation_data_length"`
	MaxProposalsPerDelegate      int64            `json:"max_proposals_per_delegate"`
	PreservedCycles              uint64           `json:"preserved_cycles"`
	BlocksPerCycle               uint64           `json:"blocks_per_cycle"`
	BlocksPerCommitment          int64            `json:"blocks_per_commitment"`
	BlocksPerRollSnapshot        int64            `json:"blocks_per_roll_snapshot"`
	BlocksPerVotingPeriod        int64            `json:"blocks_per_voting_period"`
	TimeBetweenBlocks            Int64StringSlice `json:"time_between_blocks"`
	EndorsersPerBlock            int64            `json:"endorsers_per_block"`
	HardGasLimitPerOperation     int64            `json:"hard_gas_limit_per_operation,string"`
	HardGasLimitPerBlock         int64            `json:"hard_gas_limit_per_block,string"`
	ProofOfWorkThreshold         int64            `json:"proof_of_work_threshold,string"`
	TokensPerRoll                int64            `json:"tokens_per_roll,string"`
	MichelsonMaximumTypeSize     int64            `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       int64            `json:"seed_nonce_revelation_tip,string"`
	OriginationSize              int64            `json:"origination_size"`
	BlockSecurityDeposit         int64            `json:"block_security_deposit,string"`
	EndorsementSecurityDeposit   int64            `json:"endorsement_security_deposit,string"`
	BakingRewardPerEndorsement   Int64StringSlice `json:"baking_reward_per_endorsement"`
	EndorsementReward            Int64StringSlice `json:"endorsement_reward"`
	CostPerByte                  int64            `json:"cost_per_byte,string"`
	HardStorageLimitPerOperation int64            `json:"hard_storage_limit_per_operation,string"`
	TestChainDuration            int64            `json:"test_chain_duration,string"`
	QuorumMin                    int64            `json:"quorum_min"`
	QuorumMax                    int64            `json:"quorum_max"`
	MinProposalQuorum            int64            `json:"min_proposal_quorum"`
	InitialEndorsers             int64            `json:"initial_endorsers"`
	DelayPerMissingEndorsement   int64            `json:"delay_per_missing_endorsement,string"`
}

// Int64StringSlice -
type Int64StringSlice []int64

// UnmarshalJSON -
func (slice *Int64StringSlice) UnmarshalJSON(data []byte) error {
	s := make([]string, 0)
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*slice = make([]int64, len(s))
	for i := range s {
		value, err := strconv.ParseInt(s[i], 10, 64)
		if err != nil {
			return err
		}
		(*slice)[i] = value
	}
	return nil
}

// Header -
type Header struct {
	Protocol         string    `json:"protocol"`
	ChainID          string    `json:"chain_id"`
	Hash             string    `json:"hash"`
	Level            uint64    `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	Signature        string    `json:"signature"`
}

// EndorsementWithSlot -
type EndorsementWithSlot struct {
	Endorsement Endorsement `json:"endorsement"`
	Slot        uint64      `json:"slot"`
}

// Endorsement -
type Endorsement struct {
	Branch    string               `json:"branch"`
	Operation EndorsementOperation `json:"operations"`
	Signature string               `json:"signature"`
}

// EndorsementOperation -
type EndorsementOperation struct {
	Level uint64 `json:"level"`
}

// HeadMetadata -
type HeadMetadata struct {
	Protocol        string `json:"protocol"`
	NextProtocol    string `json:"next_protocol"`
	TestChainStatus struct {
		Status string `json:"status"`
	} `json:"test_chain_status"`
	MaxOperationsTTL       uint64 `json:"max_operations_ttl"`
	MaxOperationDataLength uint64 `json:"max_operation_data_length"`
	MaxBlockHeaderLength   uint64 `json:"max_block_header_length"`
	MaxOperationListLength []struct {
		MaxSize uint64 `json:"max_size"`
		MaxOp   uint64 `json:"max_op,omitempty"`
	} `json:"max_operation_list_length"`
	Baker string `json:"baker"`
	Level struct {
		Level                uint64 `json:"level"`
		LevelPosition        uint64 `json:"level_position"`
		Cycle                uint64 `json:"cycle"`
		CyclePosition        uint64 `json:"cycle_position"`
		VotingPeriod         uint64 `json:"voting_period"`
		VotingPeriodPosition uint64 `json:"voting_period_position"`
		ExpectedCommitment   bool   `json:"expected_commitment"`
	} `json:"level"`
	LevelInfo struct {
		Level              uint64 `json:"level"`
		LevelPosition      uint64 `json:"level_position"`
		Cycle              uint64 `json:"cycle"`
		CyclePosition      uint64 `json:"cycle_position"`
		ExpectedCommitment bool   `json:"expected_commitment"`
	} `json:"level_info"`
	VotingPeriodKind string `json:"voting_period_kind"`
	VotingPeriodInfo struct {
		VotingPeriod struct {
			Index         uint64 `json:"index"`
			Kind          string `json:"kind"`
			StartPosition uint64 `json:"start_position"`
		} `json:"voting_period"`
		Position  int `json:"position"`
		Remaining int `json:"remaining"`
	} `json:"voting_period_info"`
	NonceHash      interface{}   `json:"nonce_hash"`
	ConsumedGas    string        `json:"consumed_gas"`
	Deactivated    []interface{} `json:"deactivated"`
	BalanceUpdates []struct {
		Kind     string `json:"kind"`
		Contract string `json:"contract,omitempty"`
		Change   string `json:"change"`
		Origin   string `json:"origin"`
		Category string `json:"category,omitempty"`
		Delegate string `json:"delegate,omitempty"`
		Cycle    uint64 `json:"cycle,omitempty"`
	} `json:"balance_updates"`
}

// IsManager -
func IsManager(kind string) bool {
	return kind == KindDelegation || kind == KindOrigination || kind == KindReveal || kind == KindTransaction
}
