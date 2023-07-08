package interfaces

import "context"

//go:generate mockgen -package mock_interfaces -destination mock_interfaces/iface4.go . Iface4

type Iface4 interface {
	Meth1(ctx context.Context, arg int) error
	Meth2(ctx context.Context, arg int) error
	Meth3(ctx context.Context, arg int) error
	Meth4(ctx context.Context, arg int) error
	Meth5(ctx context.Context, arg int) error
	Meth6(ctx context.Context, arg int) error
	Meth7(ctx context.Context, arg int) error
	Meth8(ctx context.Context, arg int) error
	Meth9(ctx context.Context, arg int) error
	Meth10(ctx context.Context, arg int) error
	Meth11(ctx context.Context, arg int) error
	Meth12(ctx context.Context, arg int) error
	Meth13(ctx context.Context, arg int) error
	Meth14(ctx context.Context, arg int) error
	Meth15(ctx context.Context, arg int) error
	Meth16(ctx context.Context, arg int) error
	Meth17(ctx context.Context, arg int) error
	Meth18(ctx context.Context, arg int) error
	Meth19(ctx context.Context, arg int) error
	Meth20(ctx context.Context, arg int) error
}
