package config

import (
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{name: "正常系",want: &Config{Env: "dev",Port: 8080,RedisHost: "127.0.0.1",RedisPort: 6379,RedisPassword: "P@ssword",RedisTLS: "127.0.0.1"},wantErr: false},
	}
	for _, tt := range tests {
		if strings.Compare(tt.name,"正常系") == 0 {
			t.Setenv("GAME_ENV",tt.want.Env)
			t.Setenv("PORT",strconv.Itoa(tt.want.Port))
			t.Setenv("GAME_REDIS_HOST",tt.want.RedisHost)
			t.Setenv("GAME_REDIS_PORT",strconv.Itoa(tt.want.RedisPort))
			t.Setenv("GAME_REDIS_PASSWORD",tt.want.RedisPassword)
			t.Setenv("GAME_REDIS_TLS_SERVER_NAME",tt.want.RedisTLS)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name,"正常系") {return}
			if got.Env != tt.want.Env || got.Port != tt.want.Port || got.RedisHost != tt.want.RedisHost || got.RedisPassword != tt.want.RedisPassword || got.RedisPort != tt.want.RedisPort || got.RedisTLS != tt.want.RedisTLS {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
