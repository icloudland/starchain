package service

import (
	. "github.com/icloudland/starchain/common"
)

type AccoutInfo struct {
	ProgramHash	string
	IsForze		bool
	Blanace		map[string]Fixed64
}

