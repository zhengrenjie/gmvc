package gmvc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGmvcBuilder(t *testing.T) {
	builder := CreateGmvcBuilder()
	assert.True(t, len(builder.resolverMap) == 1)
	assert.True(t, len(builder.responsor) == 3)
}

func TestRegisterTypedResolver(t *testing.T) {
	testcase1 := struct{}{}
	testcase2 := struct {
		Name string
	}{}
	typ1 := reflect.TypeOf(&testcase1)
	typ2 := reflect.TypeOf(&testcase2)

	mapp := make(map[reflect.Type]struct{})
	mapp[typ1] = struct{}{}
	mapp[typ2] = struct{}{}

	assert.True(t, len(mapp) == 2)
}
