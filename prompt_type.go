package typechat

var (
	promptUserRequest = promptType{"UserRequest"}
	promptProgram     = promptType{"Program"}
)

type promptType struct {
	name string
}

func (p promptType) String() string {
	return p.name
}
