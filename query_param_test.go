package queries

import (
	"testing"
)

func TestSnakeCasetoCamelCase(t *testing.T) {
	dts := map[string]struct {
		got  string
		want string
	}{
		"t1": {"", ""},
		"t2": {"Kevin", "Kevin"},
		"t3": {"Kevin", "Kevin"},
		"t4": {"hola_como_estas", "HolaComoEstas"},
		"t5": {"Hola_como_Estas", "HolaComoEstas"},
		"t6": {"h_o_l_a", "HOLA"},
		"t7": {"Hola3_4com_00", "Hola34com00"},
	}
	for name, dt := range dts {
		t.Run(name, func(t *testing.T) {
			r := snakeCasetoCamelCase(dt.got)
			if r != dt.want {
				t.Errorf("error se esperaba %s, se obtuvo %s", dt.want, r)
			}
		})
	}
}
func TestFieldName(t *testing.T) {
	dts := map[string]struct {
		got   string
		want1 string
		wnat2 string
	}{
		"t1":  {"", "", ""},
		"t2":  {"kevin", "kevin", "Kevin"},
		"t3":  {"kevin,Kevin", "kevin", "Kevin"},
		"t4":  {"kevin,Saucedo", "kevin", "Saucedo"},
		"t5":  {"kevin,Saucedo,Hola", "kevin", "Saucedo"},
		"t6":  {"kevin_saucedo", "kevin_saucedo", "KevinSaucedo"},
		"t7":  {"kevin_saucedo,Saucedo", "kevin_saucedo", "Saucedo"},
		"t8":  {",Hola", "", "Hola"},
		"t9":  {"kevin.saucedo,Kevin.Saucedo2,Hola", "kevin.saucedo", "Kevin.Saucedo2"},
		"t10": {"kevin_saucedo.hola_mundo", "kevin_saucedo.hola_mundo", "KevinSaucedo.HolaMundo"},
		"t11": {"kevin_saucedo.hola_mundo,Saucedo.HolaMundo", "kevin_saucedo.hola_mundo", "Saucedo.HolaMundo"},
		"t12": {",Kevin.Saucedo", "", "Kevin.Saucedo"},
	}
	for name, dt := range dts {
		t.Run(name, func(t *testing.T) {
			jsonname, modelname := fieldName(dt.got)
			if jsonname != dt.want1 {
				t.Errorf("error se esperaba jsonname %s, se obtuvo %s", dt.want1, jsonname)

			}
			if modelname != dt.wnat2 {
				t.Errorf("error se esperaba modelname %s, se obtuvo %s", dt.wnat2, modelname)
			}

		})
	}
}
