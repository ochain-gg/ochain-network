package types

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
	Id         string `cbor:"1,keyasint"`
	UniverseId string `cbor:"2,keyasint"`
	Name       string `cbor:"3,keyasint"`

	Level   uint64                 `cbor:"4,keyasint"`
	Members []OChainAllianceMember `cbor:"4,keyasint"`
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
