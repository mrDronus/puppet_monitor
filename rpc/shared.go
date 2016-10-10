package rpc

type PPTRpc string

type EmptyArgs struct{}

type RemoveNodeArgs struct {
	Host string
}

type GetStatusArgs struct {
	Errors bool
}
