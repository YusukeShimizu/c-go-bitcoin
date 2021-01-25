package ecc

import (
	"math/big"
	"testing"
)

func TestNewPoint(t *testing.T) {
	p := big.NewInt(223)
	type args struct {
		x *fieldElement
		y *fieldElement
		a *fieldElement
		b *fieldElement
	}
	tests := []struct {
		name    string
		args    args
		want    *point
		wantErr bool
	}{
		{
			name: "Ok if it called with valid arguments",
			args: args{
				x: &fieldElement{big.NewInt(-1), p},
				y: &fieldElement{big.NewInt(-1), p},
				a: &fieldElement{big.NewInt(5), p},
				b: &fieldElement{big.NewInt(7), p},
			},
			want: &point{
				x: &fieldElement{big.NewInt(-1), p},
				y: &fieldElement{big.NewInt(-1), p},
				a: &fieldElement{big.NewInt(5), p},
				b: &fieldElement{big.NewInt(7), p},
			},
		},
		{
			name: "Error if args are not on the curve",
			args: args{
				x: &fieldElement{big.NewInt(-1), p},
				y: &fieldElement{big.NewInt(-2), p},
				a: &fieldElement{big.NewInt(5), p},
				b: &fieldElement{big.NewInt(7), p},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPoint(tt.args.x, tt.args.y, tt.args.a, tt.args.b)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewEllipticCurve() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !tt.want.x.Equal(got.x) || !tt.want.y.Equal(got.y) || !tt.want.a.Equal(got.a) || !tt.want.b.Equal(got.b) {
				t.Errorf("NewEllipticCurve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_point_Eq(t *testing.T) {
	p := big.NewInt(223)
	type fields struct {
		x *fieldElement
		y *fieldElement
		a *fieldElement
		b *fieldElement
	}
	type args struct {
		other *point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "true if it called with same arguments",
			fields: fields{
				x: &fieldElement{big.NewInt(-1), p},
				y: &fieldElement{big.NewInt(-1), p},
				a: &fieldElement{big.NewInt(5), p},
				b: &fieldElement{big.NewInt(7), p},
			},
			args: args{
				&point{
					x: &fieldElement{big.NewInt(-1), p},
					y: &fieldElement{big.NewInt(-1), p},
					a: &fieldElement{big.NewInt(5), p},
					b: &fieldElement{big.NewInt(7), p},
				},
			},
			want: true,
		},
		{
			name: "Error if it called with different arguments",
			fields: fields{
				x: &fieldElement{big.NewInt(-1), p},
				y: &fieldElement{big.NewInt(-1), p},
				a: &fieldElement{big.NewInt(5), p},
				b: &fieldElement{big.NewInt(7), p},
			},
			args: args{
				&point{
					x: &fieldElement{big.NewInt(1), p},
					y: &fieldElement{big.NewInt(-1), p},
					a: &fieldElement{big.NewInt(5), p},
					b: &fieldElement{big.NewInt(7), p},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &point{
				y: tt.fields.y,
				x: tt.fields.x,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			if got := p.Eq(tt.args.other); got != tt.want {
				t.Errorf("point.Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_point_Add(t *testing.T) {
	p := big.NewInt(223)
	t.Run("Ok if it called with self.x as infinity", func(t *testing.T) {
		po, err := NewPoint(
			nil,
			nil,
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(2), p},
			&fieldElement{big.NewInt(-5), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		if !po.Eq(other) {
			t.Errorf("point.Add() = %v want %v", po, other)
		}
	})
	t.Run("Ok if it called with other.x as infinity", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(2), p},
			&fieldElement{big.NewInt(-5), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			nil,
			nil,
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		want, err := NewPoint(
			&fieldElement{big.NewInt(2), p},
			&fieldElement{big.NewInt(-5), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		if !po.Eq(want) {
			t.Errorf("point.Add() = %v want %v", po, other)
		}
	})
	t.Run("Ok if it called with different x case1", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(192), p},
			&fieldElement{big.NewInt(105), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(17), p},
			&fieldElement{big.NewInt(56), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		want, err := NewPoint(
			&fieldElement{big.NewInt(170), p},
			&fieldElement{big.NewInt(142), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		if !po.Eq(want) {
			t.Errorf("point.Add() = x:%v y:%v want x:%v y:%v", po.x.number, po.y.number, want.x.number, want.y.number)
		}
	})
	t.Run("Ok if it called with different x case2", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(47), p},
			&fieldElement{big.NewInt(71), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(117), p},
			&fieldElement{big.NewInt(141), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		want, err := NewPoint(
			&fieldElement{big.NewInt(60), p},
			&fieldElement{big.NewInt(139), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		if !po.Eq(want) {
			t.Errorf("point.Add() = x:%v y:%v want x:%v y:%v", po.x.number, po.y.number, want.x.number, want.y.number)
		}
	})

	t.Run("Ok if it called with infinity", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(2), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(2), p},
			&fieldElement{big.NewInt(-5), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		if po.x != nil || po.y != nil {
			t.Errorf("point.Add() = %v,%v want %v,%v", po.x.number, po.y.number, nil, nil)
		}
	})
	t.Run("Ok if it called with same point", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(192), p},
			&fieldElement{big.NewInt(105), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(192), p},
			&fieldElement{big.NewInt(105), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err != nil {
			t.Errorf("point.Add() error = %v", err)
			return
		}
		want, err := NewPoint(
			&fieldElement{big.NewInt(49), p},
			&fieldElement{big.NewInt(71), p},
			&fieldElement{big.NewInt(0), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		if !po.Eq(want) {
			t.Errorf("point.Add() = x:%v y:%v want x:%v y:%v", po.x.number, po.y.number, want.x.number, want.y.number)
		}
	})
	t.Run("Error if it called with different curve", func(t *testing.T) {
		po, err := NewPoint(
			&fieldElement{big.NewInt(-1), p},
			&fieldElement{big.NewInt(-1), p},
			&fieldElement{big.NewInt(5), p},
			&fieldElement{big.NewInt(7), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		other, err := NewPoint(
			&fieldElement{big.NewInt(1), p},
			&fieldElement{big.NewInt(1), p},
			&fieldElement{big.NewInt(1), p},
			&fieldElement{big.NewInt(-1), p},
		)
		if err != nil {
			t.Fatal(err)
		}
		err = po.Add(other)
		if err == nil {
			t.Errorf("point.Add() should return error but nil")
			return
		}
	})
}

func Test_point_RMul(t *testing.T) {
	p := big.NewInt(223)
	type fields struct {
		coefficient *big.Int
		x           *fieldElement
		y           *fieldElement
		a           *fieldElement
		b           *fieldElement
	}
	tests := []struct {
		name    string
		fields  fields
		want    *point
		wantErr bool
	}{
		{
			name: "Ok if it called with 2",
			fields: fields{
				coefficient: big.NewInt(2),
				x:           &fieldElement{big.NewInt(192), p},
				y:           &fieldElement{big.NewInt(105), p},
				a:           &fieldElement{big.NewInt(0), p},
				b:           &fieldElement{big.NewInt(7), p},
			},
			want: &point{
				&fieldElement{big.NewInt(49), p},
				&fieldElement{big.NewInt(71), p},
				&fieldElement{big.NewInt(0), p},
				&fieldElement{big.NewInt(7), p},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &point{
				x: tt.fields.x,
				y: tt.fields.y,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			if err := p.RMul(tt.fields.coefficient); (err != nil) != tt.wantErr {
				t.Errorf("point.RMul() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !p.Eq(tt.want) {
				t.Errorf("point.RMul() = x:%v y:%v want x:%v y:%v", p.x.number, p.y.number, tt.want.x.number, tt.want.y.number)
			}
		})
	}
}

func Test_point_FastRMul(t *testing.T) {
	p := big.NewInt(223)
	type fields struct {
		coefficient *big.Int
		x           *fieldElement
		y           *fieldElement
		a           *fieldElement
		b           *fieldElement
	}
	tests := []struct {
		name    string
		fields  fields
		want    *point
		wantErr bool
	}{
		{
			name: "Ok if it called with 2",
			fields: fields{
				coefficient: big.NewInt(2),
				x:           &fieldElement{big.NewInt(192), p},
				y:           &fieldElement{big.NewInt(105), p},
				a:           &fieldElement{big.NewInt(0), p},
				b:           &fieldElement{big.NewInt(7), p},
			},
			want: &point{
				&fieldElement{big.NewInt(49), p},
				&fieldElement{big.NewInt(71), p},
				&fieldElement{big.NewInt(0), p},
				&fieldElement{big.NewInt(7), p},
			},
		},
		{
			name: "Ok if it called with 8",
			fields: fields{
				coefficient: big.NewInt(8),
				x:           &fieldElement{big.NewInt(47), p},
				y:           &fieldElement{big.NewInt(71), p},
				a:           &fieldElement{big.NewInt(0), p},
				b:           &fieldElement{big.NewInt(7), p},
			},
			want: &point{
				&fieldElement{big.NewInt(116), p},
				&fieldElement{big.NewInt(55), p},
				&fieldElement{big.NewInt(0), p},
				&fieldElement{big.NewInt(7), p},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &point{
				x: tt.fields.x,
				y: tt.fields.y,
				a: tt.fields.a,
				b: tt.fields.b,
			}
			if err := p.FastRMul(tt.fields.coefficient); (err != nil) != tt.wantErr {
				t.Errorf("point.FastRMul() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !p.Eq(tt.want) {
				t.Errorf("point.FastRMul() = x:%v y:%v want x:%v y:%v", p.x.number, p.y.number, tt.want.x.number, tt.want.y.number)
			}
		})
	}
}
