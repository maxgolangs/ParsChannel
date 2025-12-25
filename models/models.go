package models

type ParticipantInfo struct {
	ID           int64
	FirstName    string
	LastName     string
	Username     string
	IsPremium    bool
	IsBot        bool
	IsRestricted bool
}

type ParseResult struct {
	UsernamesFile string
	IDsFile       string
	FullCSVFile   string
	TotalUsers    int
	WithUsername  int
	ChannelName   string
	TotalMembers  int
}

