package testItems

const (
	SERVICE_NAME = "自動測試"
)

var (
	TestType map[string]TestItemsInterface
)

func init() {

	TestType = make(map[string]TestItemsInterface)
	//記得新增測項目這邊也要加入map
	TestType[TaskType_DockerWatcher] = &DockerWatcher{}
}

func IsTestTypeExsit(taskType string) (isExsit bool) {
	_, isExsit = TestType[taskType]
	return
}

type TestItemsInterface interface {
	ValidScript(testName, testTag, testScript string) (err error)
	GetScriptFunc() (scriptFunc func(), err error)
}
