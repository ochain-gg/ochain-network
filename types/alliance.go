package types

import "math"

type OChainAllianceMemberType uint64

const (
	OChainAllianceLeaderMember  OChainAllianceMemberType = 3
	OChainAllianceOfficerMember OChainAllianceMemberType = 2
	OChainAllianceSimpleMember  OChainAllianceMemberType = 1
	OChainAllianceNotMember     OChainAllianceMemberType = 0
)

type OChainAllianceMember struct {
	Role    OChainAllianceMemberType `cbor:"1,keyasint"`
	Address string                   `cbor:"2,keyasint"`
}

type OChainAlliance struct {
	Id          string `cbor:"1,keyasint"`
	UniverseId  string `cbor:"2,keyasint"`
	Tag         string `cbor:"3,keyasint"`
	Name        string `cbor:"4,keyasint"`
	Description string `cbor:"5,keyasint"`

	Level   uint64                 `cbor:"6,keyasint"`
	Members []OChainAllianceMember `cbor:"7,keyasint"`
	Deleted bool                   `cbor:"8,keyasint"`
}

type OChainAllianceJoinRequest struct {
	Id         string `cbor:"1,keyasint"`
	From       string `cbor:"2,keyasint"`
	AllianceId string `cbor:"3,keyasint"`

	Accepted    bool   `cbor:"4,keyasint"`
	Canceled    bool   `cbor:"5,keyasint"`
	RequestedAt int64  `cbor:"6,keyasint"`
	AnsweredAt  int64  `cbor:"7,keyasint"`
	AnsweredBy  string `cbor:"8,keyasint"`
}

func (alliance *OChainAlliance) GetMaxAllianceSize() uint64 {
	return uint64(math.Pow(2, float64(3+alliance.Level)))
}

func (alliance *OChainAlliance) IsMember(address string) (bool, OChainAllianceMemberType) {
	for _, member := range alliance.Members {
		if member.Address == address {
			return true, member.Role
		}
	}

	return false, OChainAllianceNotMember
}

func (alliance *OChainAlliance) AddMember(address string, role OChainAllianceMemberType) {
	alliance.Members = append(alliance.Members, OChainAllianceMember{
		Role:    role,
		Address: address,
	})
}

func (alliance *OChainAlliance) RemoveMember(address string) {
	newMemberList := []OChainAllianceMember{}
	for _, member := range alliance.Members {
		if address != member.Address {
			newMemberList = append(newMemberList, member)
		}
	}

	alliance.Members = newMemberList
}

func (alliance *OChainAlliance) ChangeRole(address string, role OChainAllianceMemberType) {
	newMemberList := []OChainAllianceMember{}
	for _, member := range alliance.Members {
		if address != member.Address {
			newMemberList = append(newMemberList, member)
		} else {
			newMemberList = append(newMemberList, OChainAllianceMember{
				Address: address,
				Role:    member.Role,
			})
		}
	}

	alliance.Members = newMemberList
}
