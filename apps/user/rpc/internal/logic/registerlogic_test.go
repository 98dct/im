package logic

import (
	"context"
	"im/apps/user/rpc/user"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {

	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{in: &user.RegisterReq{
				Phone:    "15109201825",
				Nickname: "dct",
				Password: "123456",
				Avatar:   "aaa",
				Sex:      1,
			}},
			wantPrint: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerLogic := NewRegisterLogic(context.Background(), svcCtx)
			got, err := registerLogic.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}
