package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
)

// @contract:state
type StatusMessage struct {
	Records *collections.LookupMap[string, string]
}

// @contract:init
func (s *StatusMessage) Init() {
	s.Records = collections.NewLookupMap[string, string]("r")
	env.LogString("StatusMessage contract initialized")
}

// @contract:mutating
func (s *StatusMessage) SetStatus(message string) {
	if s.Records == nil {
		s.Records = collections.NewLookupMap[string, string]("r")
	}

	accountId, _ := env.GetPredecessorAccountID()

	s.Records.Insert(accountId, message)

	env.LogString("Status stored for " + accountId)
}

// @contract:view
func (s *StatusMessage) GetStatus(accountId string) string {
	if s.Records == nil {
		return ""
	}

	val, err := s.Records.Get(accountId)
	if err != nil {
		return ""
	}

	return val
}
