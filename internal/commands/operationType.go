package commands

type Action string

const (
	TodoActionAdd    Action = "add"
	TodoActionList   Action = "list"
	TodoActionDone   Action = "done"
	TodoActionUpdate Action = "update"
)

func (a Action) String() string {
	return string(a)
}
