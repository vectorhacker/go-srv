package srv

import (
	"reflect"
	"sync"
	"testing"

	"google.golang.org/grpc/naming"
)

func Test_watcher_Next(t *testing.T) {
	type fields struct {
		target   string
		existing map[string]int
		previous map[string]int
		m        sync.RWMutex
		stopChan chan bool
		errChan  chan error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*naming.Update
		wantErr bool
	}{
		{
			fields: fields{
				target: "hello.service.consul",
				existing: map[string]int{
					"10.0.0.1": 1222,
					"10.0.0.2": 2432,
					"10.0.0.3": 2344,
				},
				previous: map[string]int{},
			},
			want: []*naming.Update{
				&naming.Update{
					Op:   naming.Add,
					Addr: "10.0.0.1:1222",
				},
				&naming.Update{
					Op:   naming.Add,
					Addr: "10.0.0.2:2432",
				},
				&naming.Update{
					Op:   naming.Add,
					Addr: "10.0.0.3:2344",
				},
			},
		},
		{
			fields: fields{
				target: "hello.service.consul",
				existing: map[string]int{
					"10.0.0.1": 1222,
					"10.0.0.2": 2432,
				},
				previous: map[string]int{
					"10.0.0.1": 1222,
					"10.0.0.3": 2344,
				},
			},
			want: []*naming.Update{
				&naming.Update{
					Addr: "10.0.0.1:1222",
				},
				&naming.Update{
					Op:   naming.Delete,
					Addr: "10.0.0.3:2344",
				},
				&naming.Update{
					Op:   naming.Add,
					Addr: "10.0.0.2:2432",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &watcher{
				target:   tt.fields.target,
				existing: tt.fields.existing,
				previous: tt.fields.previous,
				m:        tt.fields.m,
				stopChan: tt.fields.stopChan,
				errChan:  tt.fields.errChan,
			}
			got, err := w.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("watcher.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("watcher.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatAddress(t *testing.T) {
	type args struct {
		addr string
		port int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				addr: "10.0.0.2",
				port: 1234,
			},
			want: "10.0.0.2:1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatAddress(tt.args.addr, tt.args.port); got != tt.want {
				t.Errorf("formatAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
