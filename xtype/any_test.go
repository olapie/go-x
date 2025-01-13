package xtype

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func jsonString(i any) string {
	b, _ := json.Marshal(i)
	return string(b)
}

func nextImage() *Image {
	return &Image{
		URL:    "https://www.image.com/" + fmt.Sprint(time.Now().Unix()),
		Width:  rand.Int31(),
		Height: rand.Int31(),
		Format: "png",
	}
}

func nextVideo() *Video {
	return &Video{
		URL:      "http://www.video.com/" + fmt.Sprint(time.Now().Unix()),
		Format:   "rmvb",
		Duration: rand.Int31(),
		Size:     rand.Int31(),
		Image:    nextImage(),
	}
}

func TestAny(t *testing.T) {
	t.Run("AliasType", func(t *testing.T) {
		type ID int
		RegisterAnyType(ID(0))

		v := NewAny(ID(10))
		b, err := json.Marshal(v)
		xtest.NoError(t, err)

		var vv *Any
		err = json.Unmarshal(b, &vv)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("String", func(t *testing.T) {
		v := NewAny("hello")
		b, err := json.Marshal(v)
		xtest.NoError(t, err)

		var vv *Any
		err = json.Unmarshal(b, &vv)
		xtest.NoError(t, err)
		xtest.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("Struct", func(t *testing.T) {
		v := NewAny(nextVideo())
		b, err := json.Marshal(v)
		xtest.NoError(t, err)

		var vv *Any
		err = json.Unmarshal(b, &vv)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(v), jsonString(vv))
	})

	t.Run("Array", func(t *testing.T) {
		var l []*Any
		l = append(l, NewAny("hello"))
		l = append(l, NewAny(nextImage()))
		l = append(l, NewAny(nextVideo()))
		b, err := json.Marshal(l)
		xtest.NoError(t, err)

		var ll []*Any
		err = json.Unmarshal(b, &ll)
		xtest.NoError(t, err)

		xtest.Equal(t, jsonString(l), jsonString(ll))
	})
}
