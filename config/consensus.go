package config

import "time"

type ConsensusConf struct {
	MaxSteps                       uint16
	ProposerTheshold               float64
	ThesholdBa                     float64
	CommitteePercent               float64
	FinalCommitteeConsensusPercent float64
	WaitBlockDelay                 time.Duration
	WaitSortitionProofDelay        time.Duration
	EstimatedBaVariance            time.Duration
	WaitForStepDelay               time.Duration
}