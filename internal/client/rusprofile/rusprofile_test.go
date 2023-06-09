package rusprofile

import (
	"context"
	"github.com/stretchr/testify/assert"
	gRPC_task "rusprofile/proto"
	"testing"
	"time"
)

func TestRusProfileParse(t *testing.T) {
	t.Skip()
	type args struct {
		INN int64
	}
	tests := []struct {
		name    string
		args    args
		want    gRPC_task.CompanyInfo
		wantErr bool
	}{
		{
			name: "Existing company",
			args: args{INN: 7802836667},
			want: gRPC_task.CompanyInfo{
				INN:      7802836667,
				KPP:      507401001,
				Name:     "ОБЩЕСТВО С ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТЬЮ \"КОТИКИ\"",
				Director: "Хромов Дмитрий Михайлович",
			}, wantErr: false,
		},
		{
			name: "Company not found",
			args: args{INN: 78028366671},
			want: gRPC_task.CompanyInfo{
				INN:      0,
				KPP:      0,
				Name:     "",
				Director: "",
			}, wantErr: true,
		},
	}

	client := NewClient("www.rusprofile.ru", time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetProfile(context.Background(), tt.args.INN)
			if (err != nil) != tt.wantErr {
				t.Errorf("RusProfileParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)

		})
	}
}
