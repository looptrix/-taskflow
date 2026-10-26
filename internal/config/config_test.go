package config

import "testing"

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Port:        "8080",
				DatabaseURL: "postgres://localhost/taskflow",
				JWTSecret:   "01234567890123456789012345678901",
			},
			wantErr: false,
		},
		{
			name: "missing database URL",
			config: Config{
				JWTSecret: "01234567890123456789012345678901",
			},
			wantErr: true,
		},
		{
			name: "missing JWT secret",
			config: Config{
				DatabaseURL: "postgres://localhost/taskflow",
			},
			wantErr: true,
		},
		{
			name: "short JWT secret",
			config: Config{
				DatabaseURL: "postgres://localhost/taskflow",
				JWTSecret:   "short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
