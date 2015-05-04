/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package gabs implements a simplified wrapper around creating and parsing JSON.
package gabs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

/*---------------------------------------------------------------------------------------------------
 */

/*
Container - an internal structure that holds a reference to the core interface map of the parsed
json. Use this container to move context.
*/
type Container struct {
	object interface{}
}

/*
Data - Return the contained data as an interface{}.
*/
func (g *Container) Data() interface{} {
	return g.object
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Path - Search for a value using dot notation.
*/
func (g *Container) Path(path string) *Container {
	return g.Search(strings.Split(path, ".")...)
}

/*
Search - Attempt to find and return an object within the JSON structure by specifying the hierarchy
of field names to locate the target. If the search encounters an array and has not reached the end
target then it will iterate each object of the array for the target and return all of the results in
a JSON array.
*/
func (g *Container) Search(hierarchy ...string) *Container {
	var object interface{}

	object = g.object
	for target := 0; target < len(hierarchy); target++ {
		if mmap, ok := object.(map[string]interface{}); ok {
			object = mmap[hierarchy[target]]
		} else if marray, ok := object.([]interface{}); ok {
			tmpArray := []interface{}{}
			for _, val := range marray {
				tmpGabs := &Container{val}
				res := tmpGabs.Search(hierarchy[target:]...).Data()
				if res != nil {
					tmpArray = append(tmpArray, res)
				}
			}
			if len(tmpArray) == 0 {
				return &Container{nil}
			}
			return &Container{tmpArray}
		} else {
			return &Container{nil}
		}
	}
	return &Container{object}
}

/*
S - Shorthand method, does the same thing as Search.
*/
func (g *Container) S(hierarchy ...string) *Container {
	return g.Search(hierarchy...)
}

/*
Children - Return a slice of all the children of the array. This also works for objects, however, the
children returned for an object will NOT be in order and you lose the names of the returned objects
this way.
*/
func (g *Container) Children() ([]*Container, error) {
	if array, ok := g.Data().([]interface{}); ok {
		children := make([]*Container, len(array))
		for i := 0; i < len(array); i++ {
			children[i] = &Container{array[i]}
		}
		return children, nil
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		children := []*Container{}
		for _, obj := range mmap {
			children = append(children, &Container{obj})
		}
		return children, nil
	}
	return nil, errors.New("parent was not a valid JSON object or array")
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Set - Set the value of a field at a JSON path, any parts of the path that do not exist will be
constructed, and if a collision occurs with a non object type whilst iterating the path an error is
returned.
*/
func (g *Container) Set(value interface{}, path ...string) (*Container, error) {
	var object interface{}

	if g.object == nil {
		g.object = map[string]interface{}{}
	}
	object = g.object
	for target := 0; target < len(path); target++ {
		if mmap, ok := object.(map[string]interface{}); ok {
			if target == len(path)-1 {
				mmap[path[target]] = value
			} else if mmap[path[target]] == nil {
				mmap[path[target]] = map[string]interface{}{}
			}
			object = mmap[path[target]]
		} else {
			return &Container{nil}, errors.New("encountered object collision whilst building path")
		}
	}
	return &Container{object}, nil
}

/*
SetP - Does the same as Set, but using a dot notation JSON path.
*/
func (g *Container) SetP(value interface{}, path string) (*Container, error) {
	return g.Set(value, strings.Split(path, ".")...)
}

/*
Object - Create a new JSON object at a path. Returns an error if the path contains a collision with
a non object type.
*/
func (g *Container) Object(path ...string) (*Container, error) {
	return g.Set(map[string]interface{}{}, path...)
}

/*
ObjectP - Does the same as Object, but using a dot notation JSON path.
*/
func (g *Container) ObjectP(path string) (*Container, error) {
	return g.Object(strings.Split(path, ".")...)
}

/*
Array - Create a new JSON array at a path. Returns an error if the path contains a collision with
a non object type.
*/
func (g *Container) Array(path ...string) (*Container, error) {
	return g.Set([]interface{}{}, path...)
}

/*
ArrayP - Does the same as Array, but using a dot notation JSON path.
*/
func (g *Container) ArrayP(path string) (*Container, error) {
	return g.Array(strings.Split(path, ".")...)
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Array modification/search - Keeping these options simple right now, no need for anything more
complicated since you can just cast to []interface{}, modify and then reassign with Set.
*/

/*
ArrayAppend - Append a value onto a JSON array.
*/
func (g *Container) ArrayAppend(value interface{}, path ...string) error {
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return errors.New("target object was not an array")
	}
	array = append(array, value)
	_, err := g.Set(array, path...)
	return err
}

/*
ArrayAppendP - Append a value onto a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayAppendP(value interface{}, path string) error {
	return g.ArrayAppend(value, strings.Split(path, ".")...)
}

/*
ArrayRemove - Remove an element from a JSON array.
*/
func (g *Container) ArrayRemove(index int, path ...string) error {
	if index < 0 {
		return errors.New("target index out of bounds")
	}
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return errors.New("target object was not an array")
	}
	if index < len(array) {
		array = append(array[:index], array[index+1:]...)
	} else {
		return errors.New("target index was out of bounds of array")
	}
	_, err := g.Set(array, path...)
	return err
}

/*
ArrayRemoveP - Remove an element from a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayRemoveP(index int, path string) error {
	return g.ArrayRemove(index, strings.Split(path, ".")...)
}

/*
ArrayElement - Access an element from a JSON array.
*/
func (g *Container) ArrayElement(index int, path ...string) (*Container, error) {
	if index < 0 {
		return &Container{nil}, errors.New("target index out of bounds")
	}
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return &Container{nil}, errors.New("target object was not an array")
	}
	if index < len(array) {
		return &Container{array[index]}, nil
	}
	return &Container{nil}, errors.New("target index was out of bounds of array")
}

/*
ArrayElementP - Access an element from a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayElementP(index int, path string) (*Container, error) {
	return g.ArrayElement(index, strings.Split(path, ".")...)
}

/*
ArrayCount - Count the number of elements in a JSON array.
*/
func (g *Container) ArrayCount(path ...string) (int, error) {
	if array, ok := g.Search(path...).Data().([]interface{}); ok {
		return len(array), nil
	}
	return 0, errors.New("target object was not an array")
}

/*
ArrayCountP - Count the number of elements in a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayCountP(path string) (int, error) {
	return g.ArrayCount(strings.Split(path, ".")...)
}

/*---------------------------------------------------------------------------------------------------
 */

/*
String - Converts the contained object back to a JSON formatted string.
*/
func (g *Container) String() string {
	if g.object != nil {
		if bytes, err := json.Marshal(g.object); err == nil {
			return string(bytes)
		}
	}
	return "{}"
}

/*
StringIndent - Converts the contained object back to a JSON formatted string with prefix and indent.
*/
func (g *Container) StringIndent(prefix string, indent string) string {
	if g.object != nil {
		if bytes, err := json.MarshalIndent(g.object, prefix, indent); err == nil {
			return string(bytes)
		}
	}
	return "{}"
}

/*
New - Create a new gabs JSON object.
*/
func New() *Container {
	return &Container{map[string]interface{}{}}
}

/*
Consume - Gobble up an already converted JSON object, or a fresh map[string]interface{} object.
*/
func Consume(root interface{}) (*Container, error) {
	if _, ok := root.(map[string]interface{}); ok {
		return &Container{root}, nil
	}
	return nil, errors.New("root was not a valid JSON object")
}

/*
ParseJSON - Convert a string into a representation of the parsed JSON.
*/
func ParseJSON(sample []byte) (*Container, error) {
	var gabs Container

	if err := json.Unmarshal(sample, &gabs.object); err != nil {
		return nil, err
	}
	if _, ok := gabs.object.(map[string]interface{}); ok {
		return &gabs, nil
	}
	return nil, errors.New("json appears to contain no data")
}

/*
ParseJSONFile - Read a file and convert into a representation of the parsed JSON.
*/
func ParseJSONFile(path string) (*Container, error) {
	if len(path) > 0 {
		cBytes, err := ioutil.ReadFile(path)
		if err != nil {
			container, err := ParseJSON(cBytes)
			if err != nil {
				return container, nil
			}
			return nil, err
		}
		return nil, err
	}
	return nil, errors.New("file path was invalid")
}

/*---------------------------------------------------------------------------------------------------
 */

// DEPRECATED METHODS

/*
Push - DEPRECATED: Push a value onto a JSON array.
*/
func (g *Container) Push(target string, value interface{}) error {
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			mmap[target] = append(array, value)
		} else {
			return errors.New("target object was not an array")
		}
	} else {
		return errors.New("parent was not a valid JSON object")
	}
	return nil
}

/*
RemoveElement - DEPRECATED: Remove a value from a JSON array.
*/
func (g *Container) RemoveElement(target string, index int) error {
	if index < 0 {
		return errors.New("target index out of bounds")
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			if index < len(array) {
				mmap[target] = append(array[:index], array[index+1:]...)
			} else {
				return errors.New("target index was out of bounds of array")
			}
		} else {
			return errors.New("target object was not an array")
		}
	} else {
		return errors.New("parent was not a valid JSON object")
	}
	return nil
}

/*
GetElement - DEPRECATED: Get the desired element from a JSON array
*/
func (g *Container) GetElement(target string, index int) *Container {
	if index < 0 {
		return &Container{nil}
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			if index < len(array) {
				return &Container{array[index]}
			}
		}
	}
	return &Container{nil}
}

/*
CountElements - DEPRECATED: Count the elements of a JSON array, returns -1 if the target is not an
array
*/
func (g *Container) CountElements(target string) int {
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			return len(array)
		}
	}
	return -1
}

/*---------------------------------------------------------------------------------------------------
 */
