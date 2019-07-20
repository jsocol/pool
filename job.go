package pool

type Job interface {
	Run() (interface{}, error)
}
