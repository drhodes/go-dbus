package dbus

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAlign(t *testing.T) {
	if 4 != _Align(4, 1) {
		t.Error("#1: Failed")
	}
	if 8 != _Align(8, 3) {
		t.Error("#2: Failed")
	}
	if 24 != _Align(8, 17) {
		t.Error("#3: Failed")
	}

}

func checkAppendAlign(t *testing.T, input string, align int, expected string) {
	buff := bytes.NewBufferString(input)
	_AppendAlign(align, buff)
	if !bytes.Equal([]byte(expected), buff.Bytes()) {
		t.Error("Failed")
	}
}

func TestAppendAlign(t *testing.T) {
	checkAppendAlign(t, "string", 4, "string\x00\x00")
	checkAppendAlign(t, "str", 8, "str\x00\x00\x00\x00\x00")
	checkAppendAlign(t, "str", 1, "str")
}

func checkAppendString(t *testing.T, input []string, expected string) {
	buff := bytes.NewBuffer([]byte{})
	for _, str := range input {
		_AppendString(buff, str)
	}
	if !bytes.Equal([]byte(expected), buff.Bytes()) {
		t.Error("Failed:expected", []byte(expected), ", actual:", buff.Bytes())
	}
}

func TestAppendString(t *testing.T) {
	checkAppendString(t, []string{"test1"}, "\x05\x00\x00\x00test1\x00")
	checkAppendString(t, []string{"string", "test2"}, "\x06\x00\x00\x00string\x00\x00\x05\x00\x00\x00test2\x00")
}

func TestAppendByte(t *testing.T) {
	buff := bytes.NewBuffer([]byte{})
	_AppendByte(buff, 1)
	if !bytes.Equal([]byte("\x01"), buff.Bytes()) {
		t.Error("#1 Failed")
	}
	_AppendByte(buff, 2)
	if !bytes.Equal([]byte("\x01\x02"), buff.Bytes()) {
		t.Error("#2 Failed")
	}
}

func TestAppendUint32(t *testing.T) {
	buff := bytes.NewBuffer([]byte{})
	_AppendUint32(buff, 1)
	if !bytes.Equal([]byte("\x01\x00\x00\x00"), buff.Bytes()) {
		t.Error("#1 Failed")
	}
	_AppendByte(buff, 2)
	_AppendUint32(buff, 0xffffffff)
	if !bytes.Equal([]byte("\x01\x00\x00\x00\x02\x00\x00\x00\xff\xff\xff\xff"), buff.Bytes()) {
		t.Error("#2 Failed")
	}
}

func TestAppendInt32(t *testing.T) {
	buff := bytes.NewBuffer([]byte{})
	_AppendInt32(buff, int32(-1))
	if !bytes.Equal([]byte("\xff\xff\xff\xff"), buff.Bytes()) {
		t.Error("#1 Failed")
	}
}

/*func TestDoublePack(t *testing.T){
	vec := vector.New(0);
	vec.Push(1);
	vec.Push(vec);
}*/

func TestAppendArray(t *testing.T) {
	teststr := "\x01\x02\x03\x04\x05\x00\x00\x00\x05\x00\x00\x00\x00\x00\x00\x00\x02"

	buff := bytes.NewBuffer([]byte{})
	_AppendByte(buff, 1)
	_AppendByte(buff, 2)
	_AppendByte(buff, 3)
	_AppendByte(buff, 4)
	_AppendByte(buff, 5)

	_AppendArray(buff, 1,
		func(b *bytes.Buffer) {
			t.Log(b.Bytes())
			_AppendAlign(8, b)
			t.Log(b.Bytes())
			_AppendByte(b, 2)
			t.Log(b.Bytes())
		})

	if teststr != string(buff.Bytes()) {
		t.Error("#1 Failed\n", buff.Bytes(), []byte(teststr))
	}
}

func TestAppendValue(t *testing.T) {
	buff := bytes.NewBuffer([]byte{})

	_AppendValue(buff, "s", "string")
	_AppendValue(buff, "s", "test2")
	if !bytes.Equal([]byte("\x06\x00\x00\x00string\x00\x00\x05\x00\x00\x00test2\x00"), buff.Bytes()) {
		t.Error("#1 Failed")
	}
	buff.Reset()
	slice := make([]interface{}, 0)
	slice = append(slice, []interface{}{"test1", uint32(1)})
	slice = append(slice, []interface{}{"test2", uint32(2)})
	slice = append(slice, []interface{}{"test3", uint32(3)})
	_AppendValue(buff, "a(su)", slice)
	if !bytes.Equal([]byte("\x34\x00\x00\x00\x00\x00\x00\x00\x05\x00\x00\x00test1\x00\x00\x00\x01\x00\x00\x00\x05\x00\x00\x00test2\x00\x00\x00\x02\x00\x00\x00\x05\x00\x00\x00test3\x00\x00\x00\x03\x00\x00\x00"), buff.Bytes()) {
		t.Error("#2 Failed", buff.Bytes())
	}
}

func TestGetByte(t *testing.T) {
	if b, _ := _GetByte([]byte("\x00\x11"), 1); b != 0x11 {
		t.Errorf("#1 Failed 0x%X != 0x11", b)
	}
	if _, e := _GetByte([]byte("\x00\x11"), 2); e == nil {
		t.Errorf("#2 Failed")
	}
}

func TestGetBoolean(t *testing.T) {
	b, e := _GetBoolean([]byte("\x01\x00\x00\x00"), 0)
	if e != nil {
		t.Error("#1-1 Failed")
	}
	if true != b {
		t.Error("#1-2 Failed")
	}
	_, e = _GetBoolean([]byte("\x01\x00\x00\x00"), 1)
	if e == nil {
		t.Error("#2 Failed")
	}
}

func TestGetString(t *testing.T) {
	s, e := _GetString([]byte("\x00\x00test"), 2, 4)
	if e != nil || s != "test" {
		t.Error("#1 Failed")
	}
	s, e = _GetString([]byte("1234"), 3, 1)
	if e != nil || s != "4" {
		t.Error("#2 Failed")
	}
}

func TestGetStructSig(t *testing.T) {
	var str string
	var e error
	str, _ = _GetStructSig("(yyy)(yyy)", 0)
	if "yyy" != str {
		t.Error("#1 Failed:", str)
	}

	str, _ = _GetStructSig("(y(ppp))yy", 0)
	if "y(ppp)" != str {
		t.Error("#2 Failed:", str)
	}

	str, _ = _GetStructSig("((test))yy", 0)
	if "(test)" != str {
		t.Error("#3 Failed:", str)
	}

	str, _ = _GetStructSig("123((test))yy", 3)
	if "(test)" != str {
		t.Error("#4 Failed:", str)
	}

	_, e = _GetStructSig("((test)(test)", 0)
	if e == nil {
		t.Error("#5 Failed")
	}

	_, e = _GetStructSig("((test(test", 0)
	if e == nil {
		t.Error("#6 Failed")
	}

}

func TestGetSigBlock(t *testing.T) {
	var str string
	str, _ = _GetSigBlock("123a3", 3)
	if "a" != str {
		t.Error("#1 Failed:", str)
	}
	str, _ = _GetSigBlock("123(abc)", 3)
	if "(abc)" != str {
		t.Error("#2 Failed:", str)
	}

}

// sliceRef([1,2,3], 1) => 2
// sliceRef([[1,2],3], 0, 1) => 2
func sliceRef(s []interface{}, arg1 int, args ...int) interface{} {
	ret := s[arg1]

	if len(args) > 0 {
		ret1 := ret.([]interface{})
		i := 0
		for ; i < len(args)-1; i++ {
			ret1 = ret1[args[i]].([]interface{})
		}
		return ret1[args[i]]
	}
	return ret
}

func TestParse(t *testing.T) {
	ret, _, _ := Parse([]byte("\x01\x02"), "y", 0)
	if !reflect.DeepEqual([]interface{}{byte(1)}, ret) {
		t.Error("#1 Failed:", ret)
	}

	ret, _, _ = Parse([]byte("\x03\x00\x00\x00\x04\x00\x00\x00test\x00\x04"), "ysy", 0)
	if !reflect.DeepEqual([]interface{}{byte(3), "test", byte(4)}, ret) {
		t.Error("#1 Failed:", ret)
	}

	ret, _, _ = Parse([]byte("\x22\x00\x00\x00\x04\x00\x00\x00test\x00\x00\x00\x00\x05\x00\x00\x00test2\x00\x00\x00\x05\x00\x00\x00test3\x00\x01"), "asy", 0)
	//	if "test" != ret.At(0).(*vector.Vector).At(0).(string) { t.Error("#3-1 Failed:")}
	if "test" != sliceRef(ret, 0, 0).(string) {
		t.Error("#3-1 Failed:")
	}
	if "test2" != sliceRef(ret, 0, 1).(string) {
		t.Error("#3-2 Failed:")
	}
	if "test3" != sliceRef(ret, 0, 2).(string) {
		t.Error("#3-3 Failed:")
	}
	if byte(1) != sliceRef(ret, 1).(byte) {
		t.Error("#3-4 Failed:")
	}

	ret, _, e := Parse([]byte("\x22\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00true\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00\x00\x00false\x00"), "a(bs)", 0)
	if e != nil {
		t.Error(e.Error())
	}
	if true != sliceRef(ret, 0, 0, 0).(bool) {
		t.Error("#4-1 Failed:")
	}
	if "true" != sliceRef(ret, 0, 0, 1).(string) {
		t.Error("#4-2 Failed:", sliceRef(ret, 0, 0, 1).(string))
	}
	if false != sliceRef(ret, 0, 1, 0).(bool) {
		t.Error("#4-3 Failed:")
	}
	if "false" != sliceRef(ret, 0, 1, 1).(string) {
		t.Error("#4-4 Failed:")
	}

	ret, _, _ = Parse([]byte("l\x00\x00\x00\x00\x01\x00\x00test"), "yu", 0)
	if 'l' != sliceRef(ret, 0).(byte) {
		t.Error("#5-1 Failed:")
	}
	if 0x100 != sliceRef(ret, 1).(uint32) {
		t.Error("#5-2 Failed:")
	}
}

func TestGetVariant(t *testing.T) {
	val, index, _ := _GetVariant([]byte("\x00\x00\x01s\x00\x00\x00\x00\x04\x00\x00\x00test\x00"), 2)
	str, ok := val[0].(string)
	if !ok {
		t.Error("#1-1 Failed")
	}
	if "test" != str {
		t.Error("#1-2 Failed", str)
	}
	if 17 != index {
		t.Error("#1-3 Failed")
	}
}

func TestParseVariant(t *testing.T) {
	vec, _, e := Parse([]byte("\x01s\x00\x00\x04\x00\x00\x00test\x00\x01y\x00\x03\x01u\x00\x04\x00\x00\x00"), "vvv", 0)
	if nil != e {
		t.Error("#1 Failed")
	}
	if "test" != vec[0].(string) {
		t.Error("#2 Failed")
	}
	if 3 != vec[1].(byte) {
		t.Error("#3 Failed")
	}
	if 4 != vec[2].(uint32) {
		t.Error("#4 Failed", vec[2].(uint32))
	}
}

func TestParseNumber(t *testing.T) {
	vec, _, e := Parse([]byte("\x04\x00\x00\x00"), "u", 0)
	if nil != e {
		t.Error("#1 Failed")
	}
	if uint32(4) != sliceRef(vec, 0).(uint32) {
		t.Error("#1 Failed", sliceRef(vec, 0).(uint32))
	}
}

func TestGetUint32(t *testing.T) {
	u, e := _GetUint32([]byte("\x04\x00\x00\x00"), 0)
	if e != nil {
		t.Error("Failed", e.Error())
	}
	if uint32(4) != u {
		t.Error("#1 Failed", u)
	}
}

func TestGetInt32(t *testing.T) {
	i, e := _GetInt32([]byte("\x04\x00\x00\x00"), 0)
	if e != nil {
		t.Error("Failed")
	}
	if int32(4) != i {
		t.Error("#1 Failed", i)
	}
}
