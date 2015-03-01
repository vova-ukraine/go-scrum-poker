package libs

type PokerRoom struct {
	Admin string
	Members map[string]PokerMember
}

type PokerMember struct {
	Vote string
	Subscribtion chan<- string
}

func (p PokerRoom) GetJSONState(currentMember string) []map[string]string {
	cards := make([]map[string]string, len(p.Members))
	var i = 0
	var finished = true
	for _, member := range(p.Members) {
		if member.Vote == "" {
			finished = false
			break
		}
	}
	for user, member := range(p.Members) {
		cards[i] = make(map[string]string)
		cards[i]["value"] = member.Vote
		if member.Vote == "" {
			cards[i]["state"] = "waiting"
		}else if !finished && currentMember != user {
			cards[i]["state"] = "played"
		} else {
			cards[i]["state"] = "open"
		}
		i++
	}

	return cards
}

func (p PokerRoom) SetVote(user, vote string) {
	member, exists := p.Members[user]
	if !exists { return }
	member.Vote = vote
	p.Members[user] = member
	go p.NotifyMembers()
}

func (p PokerRoom) NotifyMembers() {
	for _, member := range(p.Members) {
		member.Subscribtion <- "update!!!"
	}
}

func (p PokerRoom) DeleteMember(user string) {
	_, ok := p.Members[user]
	if ok {
		delete(p.Members, user)
		go p.NotifyMembers()
	}
}
