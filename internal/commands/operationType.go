package commands

type TodoAction string

const (
	TodoActionAdd    TodoAction = "add"
	TodoActionList   TodoAction = "list"
	TodoActionDone   TodoAction = "done"
	TodoActionUpdate TodoAction = "update"
)

func (a TodoAction) String() string {
	return string(a)
}
