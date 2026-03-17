package publish

import "testing"

func TestSendMessage(t *testing.T) {
	type args struct {
		topic string
		msg   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				topic: "topic",
				msg:   "aaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMessage(tt.args.topic, tt.args.msg)
		})
	}
}
