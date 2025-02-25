package toml

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTomlTreeConversionToString(t *testing.T) {
	toml, err := Load(`name = { first = "Tom", last = "Preston-Werner" }
points = { x = 1, y = 2 }`)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	reparsedTree, err := Load(toml.ToString())

	assertTree(t, reparsedTree, err, map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tom",
			"last":  "Preston-Werner",
		},
		"points": map[string]interface{}{
			"x": int64(1),
			"y": int64(2),
		},
	})
}

func TestTomlTreeConversionToStringKeysOrders(t *testing.T) {
	for i := 0; i < 100; i++ {
		tree, _ := Load(`
		foobar = true
		bar = "baz"
		foo = 1
		[qux]
		  foo = 1
		  bar = "baz2"`)

		stringRepr := tree.ToString()

		t.Log("Intermediate string representation:")
		t.Log(stringRepr)

		r := strings.NewReader(stringRepr)
		toml, err := LoadReader(r)

		if err != nil {
			t.Fatal("Unexpected error:", err)
		}

		assertTree(t, toml, err, map[string]interface{}{
			"foobar": true,
			"bar":    "baz",
			"foo":    1,
			"qux": map[string]interface{}{
				"foo": 1,
				"bar": "baz2",
			},
		})
	}
}

func testMaps(t *testing.T, actual, expected map[string]interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("trees aren't equal.\n", "Expected:\n", expected, "\nActual:\n", actual)
	}
}

func TestTomlTreeConversionToMapSimple(t *testing.T) {
	tree, _ := Load("a = 42\nb = 17")

	expected := map[string]interface{}{
		"a": int64(42),
		"b": int64(17),
	}

	testMaps(t, tree.ToMap(), expected)
}

func TestTomlTreeConversionToMapExampleFile(t *testing.T) {
	tree, _ := LoadFile("example.toml")
	expected := map[string]interface{}{
		"title": "TOML Example",
		"owner": map[string]interface{}{
			"name":         "Tom Preston-Werner",
			"organization": "GitHub",
			"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
			"dob":          time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
		},
		"database": map[string]interface{}{
			"server":         "192.168.1.1",
			"ports":          []interface{}{int64(8001), int64(8001), int64(8002)},
			"connection_max": int64(5000),
			"enabled":        true,
		},
		"servers": map[string]interface{}{
			"alpha": map[string]interface{}{
				"ip": "10.0.0.1",
				"dc": "eqdc10",
			},
			"beta": map[string]interface{}{
				"ip": "10.0.0.2",
				"dc": "eqdc10",
			},
		},
		"clients": map[string]interface{}{
			"data": []interface{}{
				[]interface{}{"gamma", "delta"},
				[]interface{}{int64(1), int64(2)},
			},
		},
	}
	testMaps(t, tree.ToMap(), expected)
}

func TestTomlTreeConversionToMapWithTablesInMultipleChunks(t *testing.T) {
	tree, _ := Load(`
	[[menu.main]]
        a = "menu 1"
        b = "menu 2"
        [[menu.main]]
        c = "menu 3"
        d = "menu 4"`)
	expected := map[string]interface{}{
		"menu": map[string]interface{}{
			"main": []interface{}{
				map[string]interface{}{"a": "menu 1", "b": "menu 2"},
				map[string]interface{}{"c": "menu 3", "d": "menu 4"},
			},
		},
	}
	treeMap := tree.ToMap()

	testMaps(t, treeMap, expected)
}

func TestTomlTreeConversionToMapWithArrayOfInlineTables(t *testing.T) {
	tree, _ := Load(`
    	[params]
	language_tabs = [
    		{ key = "shell", name = "Shell" },
    		{ key = "ruby", name = "Ruby" },
    		{ key = "python", name = "Python" }
	]`)

	expected := map[string]interface{}{
		"params": map[string]interface{}{
			"language_tabs": []interface{}{
				map[string]interface{}{
					"key":  "shell",
					"name": "Shell",
				},
				map[string]interface{}{
					"key":  "ruby",
					"name": "Ruby",
				},
				map[string]interface{}{
					"key":  "python",
					"name": "Python",
				},
			},
		},
	}

	treeMap := tree.ToMap()
	testMaps(t, treeMap, expected)
}
