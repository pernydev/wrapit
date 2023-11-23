package datapack

type DataPackage struct {
	Account    Account
	Activities []Activity // /activity in package
	Messages   []Channel
	Servers    []Guild
}

type Account struct {
	User                   User
	Email                  string
	Verifited              bool
	HasMobile              bool
	NeedsEmailVerification bool
	PremiumUntil           string
	Flags                  int
	Phone                  string
	IP                     string
	Avatar                 string // Base64

	Applications  []Application
	Payments      []Payment
	Relationships []Relationship
}

type Payment struct {
	ID             string
	Timestamp      string
	Amount         int
	Currency       string
	AmountRefunded int
}

type Relationship struct {
	ID    string
	Type  int
	User  User
	Since string
}

type User struct {
	ID            string
	Username      string
	Discriminator int
	GlobalName    string
	Bot           bool

	Avatar string // Might sometimes be null, it is a Discord CDN key or something idk what it's called
}

type Application struct {
	ID          string
	Name        string
	Icon        string // Base64
	Description string
	Summary     string
	Bot         User
	BotPublic   bool
}

type Activity struct {
	Type      string
	ID        string
	Source    string
	Day       string
	Timestamp string

	Browser string
	Device  string
	OS      string

	ChannelID   string
	ChannelType string
	GuildID     int
}

type Channel struct {
	ID    string
	Type  int
	Name  string
	Guild Guild

	Messages []Message
}

type Message struct {
	ID          string
	Timestamp   string
	Content     string
	Attachments string // tf? how does this work lol
}

type Guild struct {
	ID   string
	Name string
}
