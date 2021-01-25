package ecc

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		number *big.Int
		prime  *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *fieldElement
		wantErr bool
	}{
		{
			name: "Ok if it called with proper values",
			args: args{
				number: big.NewInt(3),
				prime:  big.NewInt(5),
			},
			want: &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
		},
		{
			name: "NG if number is lager than prime",
			args: args{
				number: big.NewInt(5),
				prime:  big.NewInt(3),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFieldElement(tt.args.number, tt.args.prime)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldElement_Equal(t *testing.T) {
	tests := []struct {
		name  string
		self  *fieldElement
		other *fieldElement
		want  bool
	}{
		{
			name:  "Ok if it called with different number",
			self:  &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other: &fieldElement{number: big.NewInt(4), prime: big.NewInt(5)},
			want:  false,
		},
		{
			name:  "Ok if it called with same number and prime",
			self:  &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other: &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			want:  true,
		},
		{
			name:  "Ok if it called with nil other",
			self:  &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other: nil,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.self.Equal(tt.other); got != tt.want {
				t.Errorf("fieldElement.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldElement_Add(t *testing.T) {
	tests := []struct {
		name    string
		self    *fieldElement
		other   *fieldElement
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "Ok if it called with proper value",
			self:  &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other: &fieldElement{number: big.NewInt(2), prime: big.NewInt(5)},
			want:  big.NewInt(0),
		},
		{
			name:    "NG if different prime",
			self:    &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other:   &fieldElement{number: big.NewInt(2), prime: big.NewInt(4)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.self.Add(tt.self, tt.other)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fieldElement.Add() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.self.number.Cmp(tt.want) != 0 {
				t.Errorf("fieldElement.Add() result = %v, want %v", tt.self.number, tt.want)
			}
		})
	}
}

func Test_fieldElement_Sub(t *testing.T) {
	tests := []struct {
		name    string
		self    *fieldElement
		other   *fieldElement
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "Ok if it called with proper value",
			self:  &fieldElement{number: big.NewInt(4), prime: big.NewInt(5)},
			other: &fieldElement{number: big.NewInt(1), prime: big.NewInt(5)},
			want:  big.NewInt(3),
		},
		{
			name:    "NG if different prime",
			self:    &fieldElement{number: big.NewInt(3), prime: big.NewInt(5)},
			other:   &fieldElement{number: big.NewInt(3), prime: big.NewInt(4)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.self.Sub(tt.self, tt.other)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fieldElement.Sub() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.self.number.Cmp(tt.want) != 0 {
				t.Errorf("fieldElement.Sub() result = %v, want %v", tt.self.number, tt.want)
			}
		})
	}
}

func Test_fieldElement_Mul(t *testing.T) {
	tests := []struct {
		name    string
		self    *fieldElement
		other   *fieldElement
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "Ok if it called with proper value",
			self:  &fieldElement{number: big.NewInt(2), prime: big.NewInt(5)},
			other: &fieldElement{number: big.NewInt(4), prime: big.NewInt(5)},
			want:  big.NewInt(3),
		},
		{
			name:    "NG if different prime",
			self:    &fieldElement{number: big.NewInt(2), prime: big.NewInt(5)},
			other:   &fieldElement{number: big.NewInt(4), prime: big.NewInt(4)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.self.Mul(tt.self, tt.other)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fieldElement.Mul() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.self.number.Cmp(tt.want) != 0 {
				t.Errorf("fieldElement.Mul() result = %v, want %v", tt.self.number, tt.want)
			}
		})
	}
}

func Test_fieldElement_Pow(t *testing.T) {
	tests := []struct {
		name    string
		self    *fieldElement
		exp     *big.Int
		want    *big.Int
		wantErr bool
	}{
		{
			name: "Ok if it called with proper value",
			self: &fieldElement{number: big.NewInt(2), prime: big.NewInt(5)},
			exp:  big.NewInt(3),
			want: big.NewInt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.self.Pow(tt.self, tt.exp)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fieldElement.Pow() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.self.number.Cmp(tt.want) != 0 {
				t.Errorf("fieldElement.Pow() result = %v, want %v", tt.self.number, tt.want)
			}
		})
	}
}

func Test_fieldElement_Div(t *testing.T) {
	tests := []struct {
		name    string
		self    *fieldElement
		other   *fieldElement
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "Ok if it called with proper value",
			self:  &fieldElement{number: big.NewInt(3), prime: big.NewInt(31)},
			other: &fieldElement{number: big.NewInt(24), prime: big.NewInt(31)},
			want:  big.NewInt(4),
		},
		{
			name:    "NG if different prime",
			self:    &fieldElement{number: big.NewInt(7), prime: big.NewInt(19)},
			other:   &fieldElement{number: big.NewInt(5), prime: big.NewInt(20)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.self.Div(tt.self, tt.other)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fieldElement.Div() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.self.number.Cmp(tt.want) != 0 {
				t.Errorf("fieldElement.Div() result = %v, want %v", tt.self.number, tt.want)
			}
		})
	}
}
