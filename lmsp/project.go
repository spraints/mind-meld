package lmsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// Project is the deserialized verion of a scratch.sb3 project.json file.
//
// https://en.scratch-wiki.info/wiki/Scratch_File_Format
//
// TODO - move this to a package for scratch.sb3 types.
type Project struct {
	Targets    []ProjectTarget    `json:"targets"`
	Monitors   []ProjectMonitor   `json:"monitors"`
	Extensions []ProjectExtension `json:"extensions"`
	Meta       ProjectMeta        `json:"meta"`
}

// ProjectExtension is an identifier of a scratch extension used by this project.
type ProjectExtension string

type ProjectMonitor TODO
type ProjectMeta TODO

// ProjectTarget is the stage or a sprite.
type ProjectTarget struct {
	// True if this is the stage and false otherwise. Defaults to false.
	IsStage bool `json:"isStage"`

	// The name of the sprite. Always "Stage" for the stage. If not provided, the target will not be loaded.
	Name string `json:"name"`

	// An object associating IDs with arrays representing variables. The first element of the array is the variable name and the second is the value.
	Variables map[ProjectVariableID]ProjectVariable `json:"variables"`

	// An object associating IDs with arrays representing lists. The first element of the array is the list name and the second is the list as an array.
	Lists map[ProjectListID]ProjectList `json:"lists"`

	// An object associating IDs with broadcast names. Normally only present in the stage.
	Broadcasts map[ProjectBroadcastID]ProjectBroadcast `json:"broadcasts"`

	// An object associating IDs with blocks.
	Blocks ProjectBlocks `json:"blocks"`

	// An object associating IDs with comments.
	Comments map[ProjectCommentID]ProjectComment `json:"comments"`

	// The costume number.
	CurrentCostume ProjectCostumeID `json:"currentCostume"`

	// An array of costumes.
	Costumes []ProjectCostume `json:"costumes"`

	// An array of sounds.
	Sounds []ProjectSound `json:"sounds"`

	// The layer number.
	LayerOrder ProjectLayerNumber `json:"layerOrder,omitempty"` // this doesn't appear in scratch.sb3.

	// The volume.
	Volume ProjectVolume `json:"volume"`

	// The stage properties are pointers so that they can be omitted when the Target isn't a stage.

	// The tempo in BPM.
	Tempo *ProjectTempo `json:"tempo,omitempty"`

	// Possible values are "on", "off" and "on-flipped".[4] Determines if video is visible on the stage and if it is flipped. Has no effect if the project does not use an extension with video input.
	VideoState *ProjectVideoState `json:"videoState,omitempty"`

	// The video transparency. Defaults to 50. Has no effect if videoState is "off" or if the project does not use an extension with video input.
	VideoTransparency *ProjectVideoTransparency `json:"videoTransparency,omitempty"`

	// The language of the Text to Speech extension. Defaults to the editor language.
	// This is a double-pointer so that the stage can set it to null in JSON.
	// I can't get this to behave correctly (omit when not set, 'null' if set to nil, or a string) so this will be emitted more often than it's supposed to.
	TextToSpeechLanguage ProjectLanguage `json:"textToSpeechLanguage"`

	Visible       *bool                 `json:"visible,omitempty"`
	X             *int                  `json:"x,omitempty"`
	Y             *int                  `json:"y,omitempty"`
	Size          *int                  `json:"size,omitempty"`
	Direction     *int                  `json:"direction,omitempty"`
	Draggable     *bool                 `json:"draggable,omitempty"`
	RotationStyle *ProjectRotationStyle `json:"rotationStyle,omitempty"`
}

func (t ProjectTarget) GetRootBlockIDs() []ProjectBlockID {
	nodes := map[ProjectBlockID]bool{}
	for id, block := range t.Blocks {
		switch block := block.(type) {
		case *ProjectBlockObject:
			if block.Next != nil {
				nodes[*block.Next] = false
			}
			if block.Parent == nil {
				if _, ok := nodes[id]; !ok {
					nodes[id] = true
				}
			}
		default:
			panic(fmt.Sprintf("%q %T %#v", id, block, block))
		}
	}
	var roots []ProjectBlockID
	for id, isRoot := range nodes {
		if isRoot {
			roots = append(roots, id)
		}
	}
	sort.Slice(roots, func(a, b int) bool { return roots[a] < roots[b] })
	return roots
}

const (
	VideoState_On        ProjectVideoState = "on"
	VideoState_Off                         = "off"
	VideoState_OnFlipped                   = "on-flipped"
)

const (
	RotationStyle_AllAround ProjectRotationStyle = "all around"
	// others?
)

type ProjectListID string
type ProjectVariableID string
type ProjectBroadcastID string
type ProjectBroadcast string
type ProjectBlockID string
type ProjectCommentID string
type ProjectCostumeID uint32

type ProjectLayerNumber int // ??
type ProjectVolume int

type ProjectTempo int
type ProjectVideoTransparency int
type ProjectVideoState string
type ProjectLanguage string

func (l ProjectLanguage) MarshalJSON() ([]byte, error) {
	if l == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(l))
}

func (l *ProjectLanguage) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*l = ProjectLanguage(s)
	return nil
}

type ProjectOpcode string
type ProjectFieldName string
type ProjectRotationStyle string

type ProjectList struct {
	Name   string
	Values []interface{}
}

func (l ProjectList) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{l.Name, l.Values})
}

func (l *ProjectList) UnmarshalJSON(data []byte) error {
	var vals []interface{}
	err := json.Unmarshal(data, &vals)
	if err != nil {
		return err
	}
	if len(vals) != 2 {
		return errors.Errorf("expected name and value for a list but got %d fields", len(vals))
	}
	if name, ok := vals[0].(string); !ok {
		return errors.Errorf("expected list name but found %+v", vals[0])
	} else {
		l.Name = name
	}
	if values, ok := vals[1].([]interface{}); !ok {
		return errors.Errorf("expected list values but found %+v", vals[1])
	} else {
		l.Values = values
	}
	return nil
}

type ProjectVariable struct {
	Name  string
	Value interface{}
}

func (v ProjectVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{v.Name, v.Value})
}

func (v *ProjectVariable) UnmarshalJSON(data []byte) error {
	var vals []interface{}
	err := json.Unmarshal(data, &vals)
	if err != nil {
		return err
	}
	if len(vals) != 2 {
		return errors.Errorf("expected name and value for a variable but got %d fields", len(vals))
	}
	if name, ok := vals[0].(string); !ok {
		return errors.Errorf("expected variable name but found %+v", vals[0])
	} else {
		v.Name = name
	}
	v.Value = vals[1]
	return nil
}

type ProjectBlocks map[ProjectBlockID]ProjectBlock

func (b ProjectBlocks) MarshalJSON() ([]byte, error) {
	var blocks map[ProjectBlockID]ProjectBlock = b
	return json.Marshal(blocks)
}

func (b *ProjectBlocks) UnmarshalJSON(data []byte) error {
	var blocks map[ProjectBlockID]json.RawMessage
	if err := json.Unmarshal(data, &blocks); err != nil {
		return err
	}
	res := make(ProjectBlocks, len(blocks))
	for id, data := range blocks {
		if val, err := unmarshalProjectBlock(data); err != nil {
			return err
		} else {
			res[id] = val
		}
	}
	*b = res
	return nil
}

const (
	projectBlockNumber          int = 4
	projectBlockPositiveNumber      = 5
	projectBlockPositiveInteger     = 6
	projectBlockInteger             = 7
	projectBlockAngle               = 8
	projectBlockColor               = 9
	projectBlockString              = 10
	projectBlockBroadcast           = 11
	projectBlockVariable            = 12
	projectBlockList                = 13
)

func unmarshalProjectBlock(data json.RawMessage) (ProjectBlock, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil, errors.Errorf("expected json array or object but got %q", data)
	}
	switch delim {
	case '[':
		var vals []interface{}
		err := json.Unmarshal(data, &vals)
		if err != nil {
			return nil, err
		}
		if len(vals) < 1 {
			return nil, errors.Errorf("expected at least one element in array")
		}
		fcode, ok := vals[0].(float64)
		if !ok {
			return nil, errors.Errorf("expected first element (%+v, %T) to be a number in %q", vals[0], vals[0], data)
		}
		code := int(fcode)
		switch code {
		case projectBlockNumber, projectBlockPositiveNumber:
			return &ProjectBlockNumber{
				code:  code,
				Value: vals[1].(float64),
			}, nil
		case projectBlockPositiveInteger, projectBlockInteger:
			return &ProjectBlockInt{
				code:  code,
				Value: vals[1].(int64),
			}, nil
		case projectBlockAngle:
			return &ProjectBlockAngle{
				Value: vals[1].(float64),
			}, nil
		case projectBlockColor:
			return &ProjectBlockColor{
				Value: vals[1].(string),
			}, nil
		case projectBlockString:
			return &ProjectBlockString{
				Value: vals[1].(string),
			}, nil
		case projectBlockBroadcast:
			return &ProjectBlockBroadcast{
				Name: vals[1].(string),
				ID:   ProjectBroadcastID(vals[2].(string)),
			}, nil
		case projectBlockVariable:
			return &ProjectBlockVariable{
				Name:   vals[1].(string),
				ID:     ProjectVariableID(vals[2].(string)),
				Coords: vals[3:],
			}, nil
		case projectBlockList:
			return &ProjectBlockList{
				Name:   vals[1].(string),
				ID:     ProjectListID(vals[2].(string)),
				Coords: vals[3:],
			}, nil
		default:
			return nil, errors.Errorf("unrecognized block type in %q", data)
		}
	case '{':
		var res ProjectBlockObject
		err := json.Unmarshal(data, &res)
		return &res, err
	default:
		return nil, errors.Errorf("expected json array or object but got %q", data)
	}
}

type ProjectBlock interface {
	Description() string
}

type ProjectBlockNumber struct {
	code  int
	Value float64
}

func (n ProjectBlockNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{n.code, n.Value})
}

func (n ProjectBlockNumber) Description() string {
	return fmt.Sprintf("[number %f]", n.Value)
}

type ProjectBlockInt struct {
	code  int
	Value int64
}

func (n ProjectBlockInt) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{n.code, n.Value})
}

func (n ProjectBlockInt) Description() string {
	return fmt.Sprintf("[integer %d]", n.Value)
}

type ProjectBlockAngle struct {
	Value float64
}

func (n ProjectBlockAngle) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{projectBlockAngle, n.Value})
}

func (n ProjectBlockAngle) Description() string {
	return fmt.Sprintf("[angle %f]", n.Value)
}

type ProjectBlockColor struct {
	Value string
}

func (n ProjectBlockColor) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{projectBlockColor, n.Value})
}

func (n ProjectBlockColor) Description() string {
	return n.Value
}

type ProjectBlockString struct {
	Value string
}

func (n ProjectBlockString) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{projectBlockString, n.Value})
}

func (n ProjectBlockString) Description() string {
	return n.Value
}

type ProjectBlockBroadcast struct {
	Name string
	ID   ProjectBroadcastID
}

func (n ProjectBlockBroadcast) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{projectBlockBroadcast, n.Name, n.ID})
}

func (n ProjectBlockBroadcast) Description() string {
	return fmt.Sprintf("[broadcast %s/%s]", n.Name, n.ID)
}

type ProjectBlockVariable struct {
	Name   string
	ID     ProjectVariableID
	Coords []interface{}
}

func (n ProjectBlockVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(append([]interface{}{projectBlockVariable, n.Name, n.ID}, n.Coords...))
}

func (n ProjectBlockVariable) Description() string {
	return fmt.Sprintf("[variable %s/%s %+v]", n.Name, n.ID, n.Coords)
}

type ProjectBlockList struct {
	Name   string
	ID     ProjectListID
	Coords []interface{}
}

func (n ProjectBlockList) MarshalJSON() ([]byte, error) {
	return json.Marshal(append([]interface{}{projectBlockList, n.Name, n.ID}, n.Coords...))
}

func (n ProjectBlockList) Description() string {
	return fmt.Sprintf("[list %s/%s %+v]", n.Name, n.ID, n.Coords)
}

type ProjectBlockObject struct {
	// A string naming the block. The opcode of a "core" block may be found
	// in the Scratch source code here or here for shadows, and the opcode
	// of an extension's block may be found in the extension's source code
	// here.
	Opcode ProjectOpcode `json:"opcode"`

	// The ID of the following block or null.
	Next *ProjectBlockID `json:"next"`

	// If the block is a stack block and is preceded, this is the ID of the
	// preceding block. If the block is the first stack block in a C mouth,
	// this is the ID of the C block. If the block is an input to another
	// block, this is the ID of that other block. Otherwise it is null.
	Parent *ProjectBlockID `json:"parent"`

	// An object associating names with arrays representing inputs into
	// which other blocks may be dropped, including C mouths. The first
	// element of each array is 1 if the input is a shadow, 2 if there is
	// no shadow, and 3 if there is a shadow but it is obscured by the
	// input. The second is either the ID of the input or an array
	// representing it as described in the table below. If there is an
	// obscured shadow, the third element is its ID or an array
	// representing it.
	Inputs TODO `json:"inputs"`

	// An object associating names with arrays representing fields. The
	// first element of each array is the field's value. For certain
	// fields, such as variable and broadcast dropdown menus, there is also
	// a second element, which is the ID of the field's value.
	Fields map[ProjectFieldName]ProjectField `json:"fields"`

	// True if this is a shadow block and false otherwise.
	Shadow bool `json:"shadow"`

	// False if the block has a parent and true otherwise.
	TopLevel bool `json:"topLevel"`

	// A top-level block object also has the x- and y-coordinates of the
	// block in the code area as x and y.
	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`

	// A block with a comment attached has a comment property whose value
	// is the comment's ID.
	Comment ProjectCommentID `json:"comment,omitempty"`

	// A block with a mutation also has a mutation property whose value is
	// an object representing the mutation.
	Mutation *ProjectMutation `json:"mutation,omitempty"`
}

func (o ProjectBlockObject) Description() string {
	return fmt.Sprintf("[object %s]", o.Opcode)
}

type ProjectField TODO // [name, optional ID of field's value]

// Mutations are present on blocks where the opcode property is equal to
// "procedures_call" (i.e.  'custom block') or "procedures_prototype" (i.e. the
// inner part of 'define [custom block]'). Mutations have the following
// properties:
type ProjectMutation struct {
	// Always equal to "mutation".
	TagName string `json:"tagName"`

	// Seems to always be an empty array.
	Children TODO `json:"children"`

	// The name of the custom block, including inputs: %s for string/number
	// inputs and %b for boolean inputs.
	ProcCode string `json:"proccode"`

	// An array of the ids of the arguments; these can also be found in the
	// input property of the main block.
	ArgumentIDs TODO `json:"argumentids"`

	// An array of the names of the arguments. This is only present when
	// the block has an opcode of procedures_prototype.
	ArgumentNames TODO `json:"argumentnames"`

	// An array of the defaults of the arguments; for string/number
	// arguments, this is an empty string, and for boolean arguments it is
	// false. This is only present when the block has an opcode of
	// procedures_prototype.
	ArgumentDefaults TODO `json:"argumentdefaults"`

	// Whether to run the block without screen refresh or not.
	Warp string `json:"warp"`
}

type ProjectComment struct {
	Width     float64         `json:"width"`
	Height    float64         `json:"height"`
	Minimized bool            `json:"minimized"`
	Text      string          `json:"text"`
	BlockID   *ProjectBlockID `json:"blockId"`
	X         float64         `json:"x"`
	Y         float64         `json:"y"`
}

type ProjectCostume TODO

type ProjectSound TODO

type TODO interface{}
