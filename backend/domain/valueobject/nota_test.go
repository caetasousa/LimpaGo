package valueobject

import "testing"

func TestNovaNota(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		valor   int
		want    int
		wantErr bool
	}{
		{"nota 0 valida", 0, 0, false},
		{"nota 1 valida", 1, 1, false},
		{"nota 2 valida", 2, 2, false},
		{"nota 3 valida", 3, 3, false},
		{"nota 4 valida", 4, 4, false},
		{"nota 5 valida", 5, 5, false},
		{"nota -1 invalida", -1, 0, true},
		{"nota 6 invalida", 6, 0, true},
		{"nota -100 invalida", -100, 0, true},
		{"nota 100 invalida", 100, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NovaNota(tt.valor)

			if (err != nil) != tt.wantErr {
				t.Errorf("NovaNota(%d) error = %v; wantErr %v", tt.valor, err, tt.wantErr)
				return
			}

			if !tt.wantErr && int(got) != tt.want {
				t.Errorf("NovaNota(%d) = %d; want %d", tt.valor, got, tt.want)
			}
		})
	}
}
