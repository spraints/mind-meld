package lmsdump

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spraints/mind-meld/lmsp"
)

func Dump(w io.Writer, proj lmsp.Project) error {
	for _, target := range proj.Targets {
		if _, err := fmt.Fprintf(w, "target: %s\n", target.Name); err != nil {
			return err
		}
		// TODO check for errors in renderTarget.
		renderTarget(indent(w), target)
	}
	return nil
}

// It's not perfect, but here's more or less the layout below:
// + visit* are the generic funcs for walking the data structure.
//   - These will eventually be able to change from renderX(w,target,block) to w.visitBlock(target,block).
// + render* are the visitor that writes pseudocode to a Writer.

func renderTarget(w io.Writer, target lmsp.ProjectTarget) {
	for _, id := range target.GetRootBlockIDs() {
		fmt.Fprintf(w, "----- %s -----\n", id)
		visitBlock(w, target, id)
		fmt.Fprintln(w)
	}
	first := true
	for _, id := range target.GetStandaloneCommentIDs() {
		renderComment(w, target, id)
		if first {
			fmt.Fprintln(w, "--------------------------------")
			first = false
		}
	}
	// todo - other fields of ProjectTarget.
}

func visitBlock(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectBlockID) {
	block := target.Blocks[id].(*lmsp.ProjectBlockObject)
	if block.Comment != "" {
		renderComment(w, target, block.Comment)
	}
	switch block.Opcode {
	case "flippercontrol_stopOtherStacks",
		"flipperdisplay_displayOff",
		"flippermoremove_moveDidMovement",
		"flippermove_stopMove",
		"flippersensors_motion",
		"flippersensors_orientation",
		"flippersensors_resetTimer",
		"flippersensors_resetYaw",
		"flippersensors_timer",
		"flippersound_stopSound",
		"sound_cleareffects",
		"sound_volume":
		renderAction(w, target, block)

	case "argument_reporter_string_number":
		renderFieldSelector(w, target, block, "VALUE")
	case "control_forever":
		renderForever(w, target, block)
	case "control_if":
		renderControl(w, target, block, "if")
	case "control_if_else":
		renderIfElse(w, target, block)
	case "control_repeat_until":
		renderControl(w, target, block, "until")
	case "control_wait_until":
		renderWaitUntil(w, target, block)
	case "control_wait":
		renderWait(w, target, block)
	case "data_changevariableby":
		renderChangeVariableBy(w, target, block)
	case "data_setvariableto":
		renderSetVariableTo(w, target, block)
	case "event_broadcast":
		renderAction(w, target, block, inputArg("BROADCAST_INPUT"))
	case "event_whenbroadcastreceived":
		renderWhenBroadcastReceived(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flippercontrol_stop":
		renderAction(w, target, block, fieldArg("STOP_OPTION"))
	case "flipperdisplay_centerButtonLight":
		renderAction(w, target, block, namedInputArg("COLOR"))
	case "flipperdisplay_color-selector-vertical":
		renderFieldSelector(w, target, block, "field_flipperdisplay_color-selector-vertical")
	case "flipperdisplay_custom-animate-matrix":
		renderFieldSelector(w, target, block, "field_flipperdisplay_custom-animate-matrix")
	case "flipperdisplay_custom-matrix":
		renderFieldSelector(w, target, block, "field_flipperdisplay_custom-matrix")
	case "flipperdisplay_ledAnimation":
		renderAction(w, target, block, namedInputArg("MATRIX"))
	case "flipperdisplay_ledImage":
		renderAction(w, target, block, namedInputArg("MATRIX"))
	case "flipperdisplay_ledImageFor":
		renderAction(w, target, block, namedInputArg("MATRIX"), namedInputArg("VALUE"))
	case "flipperevents_force-sensor-selector":
		renderFieldSelector(w, target, block, "field_flipperevents_force-sensor-selector")
	case "flipperevents_whenPressed":
		renderWhenPressed(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenProgramStarts":
		renderWhenProgramStarts(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flippermoremotor_motorDidMovement":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermoremotor_motorSetDegreeCounted":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("VALUE"))
	case "flippermoremotor_motorTurnForSpeed":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("SPEED"), fieldInputArg("UNIT", "VALUE"))
	case "flippermoremotor_multiple-port-selector":
		renderFieldSelector(w, target, block, "field_flippermoremotor_multiple-port-selector")
	case "flippermoremotor_position":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermoremotor_single-motor-selector":
		renderFieldSelector(w, target, block, "field_flippermoremotor_single-motor-selector")
	case "flippermoresensors_force-sensor-selector":
		renderFieldSelector(w, target, block, "field_flippermoresensors_force-sensor-selector")
	case "flippermoresensors_isPressed":
		renderAction(w, target, block, namedInputArg("PORT"), namedFieldArg("OPTION"))
	case "flippermotor_absolutePosition":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermotor_custom-angle":
		renderFieldSelector(w, target, block, "field_flippermotor_custom-angle")
	case "flippermotor_custom-icon-direction":
		renderFieldSelector(w, target, block, "field_flippermotor_custom-icon-direction")
	case "flippermotor_motorGoDirectionToPosition":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("POSITION"), namedFieldArg("DIRECTION"))
	case "flippermotor_motorSetSpeed":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("SPEED"))
	case "flippermotor_motorStartDirection":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("DIRECTION"))
	case "flippermotor_motorStop":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermotor_motorTurnForDirection":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("DIRECTION"), fieldInputArg("UNIT", "VALUE"))
	case "flippermotor_multiple-port-selector":
		renderFieldSelector(w, target, block, "field_flippermotor_multiple-port-selector")
	case "flippermotor_single-motor-selector":
		renderFieldSelector(w, target, block, "field_flippermotor_single-motor-selector")
	case "flippermotor_speed":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermove_custom-icon-direction":
		renderFieldSelector(w, target, block, "field_flippermove_custom-icon-direction")
	case "flippermove_move":
		renderAction(w, target, block, namedInputArg("DIRECTION"), fieldInputArg("UNIT", "VALUE"))
	case "flippermove_movementSpeed":
		renderAction(w, target, block, namedInputArg("SPEED"))
	case "flippermove_movement-port-selector":
		renderFieldSelector(w, target, block, "field_flippermove_movement-port-selector")
	case "flippermove_rotation-wheel":
		renderFieldSelector(w, target, block, "field_flippermove_rotation-wheel")
	case "flippermove_setMovementPair":
		renderAction(w, target, block, namedInputArg("PAIR"))
	case "flippermove_startSteer":
		renderAction(w, target, block, namedInputArg("STEERING"))
	case "flippermove_steer":
		renderAction(w, target, block, namedInputArg("STEERING"), fieldInputArg("UNIT", "VALUE"))
	case "flippersensors_color-sensor-selector":
		renderFieldSelector(w, target, block, "field_flippersensors_color-sensor-selector")
	case "flippersensors_isReflectivity":
		renderIsReflectivity(w, target, block)
	case "flippersensors_orientationAxis":
		renderAction(w, target, block, fieldArg("AXIS"))
	case "flippersound_beep":
		renderAction(w, target, block, namedInputArg("NOTE"))
	case "flippersound_custom-piano":
		renderFieldSelector(w, target, block, "field_flippersound_custom-piano")
	case "flippersound_playSound":
		renderAction(w, target, block, namedInputArg("SOUND"))
	case "flippersound_sound-selector":
		renderFieldSelector(w, target, block, "field_flippersound_sound-selector")
	case "operator_add":
		renderBinaryOperator(w, target, block, "+", "NUM1", "NUM2")
	case "operator_equals":
		renderBinaryOperator(w, target, block, "==", "OPERAND1", "OPERAND2")
	case "operator_gt":
		renderBinaryOperator(w, target, block, ">", "OPERAND1", "OPERAND2")
	case "operator_lt":
		renderBinaryOperator(w, target, block, "<", "OPERAND1", "OPERAND2")
	case "operator_multiply":
		renderBinaryOperator(w, target, block, "*", "NUM1", "NUM2")
	case "operator_not":
		renderUnaryOperator(w, target, block, "not", "OPERAND")
	case "operator_subtract":
		renderBinaryOperator(w, target, block, "-", "NUM1", "NUM2")
	case "procedures_call":
		renderProcedureCall(w, target, block)
	case "procedures_definition":
		renderProcedureDefinition(w, target, block)
		w = indent(w) // TODO - move this to a 'renderX' func.
	case "procedures_prototype":
		renderProcedurePrototype(w, target, block)
	default:
		visitOtherBlock(w, target, block)
	}
	if block.Next != nil {
		fmt.Fprintln(w) // TODO - move this to a 'renderX' func.
		visitBlock(w, target, *block.Next)
	}
}

func visitOtherBlock(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	type sortableArg struct {
		name string
		a    argFn
		code string
	}
	argSorter := make([]sortableArg, 0, len(block.Inputs)+len(block.Fields))
	for input := range block.Inputs {
		argSorter = append(argSorter, sortableArg{string(input), namedInputArg(input), fmt.Sprintf("namedInputArg(%q)", input)})
	}
	for field := range block.Fields {
		argSorter = append(argSorter, sortableArg{string(field), namedFieldArg(field), fmt.Sprintf("namedFieldArg(%q)", field)})
	}
	sort.Slice(argSorter, func(i, j int) bool { return argSorter[i].name < argSorter[j].name })
	args := make([]argFn, 0, len(argSorter))
	argCode := make([]string, 0, len(argSorter))
	for _, sa := range argSorter {
		args = append(args, sa.a)
		argCode = append(argCode, sa.code)
	}
	suggestf("  case %q:\n    renderAction(w, target, block, %s)\n",
		block.Opcode,
		strings.Join(argCode, ", "))
	renderAction(w, target, block, args...)
}

func renderComment(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectCommentID) {
	fmt.Fprintf(w, "/****\n  %s\n****/\n", target.Comments[id].Text)
}

func renderProcedureCall(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "%s %s(", block.Mutation.ProcCode, block.Mutation.ArgumentIDs)
	inputs := block.Inputs
	first := true
	for id := range inputs {
		if !first {
			fmt.Fprint(w, ", ")
		}
		first = false
		fmt.Fprintf(w, "%s: ", id)
		visitInput(w, target, block, id)
	}
	fmt.Fprintf(w, ")")
}

func renderProcedureDefinition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "def ")
	visitInput(w, target, block, "custom_block")
}

func renderProcedurePrototype(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "%s %s", block.Mutation.ProcCode, block.Mutation.ArgumentNames)
	// Inputs is redundant with argument names.
}

// This goes with the visit* funcs.
type argFn func(io.Writer, lmsp.ProjectTarget, *lmsp.ProjectBlockObject)

// This goes with the render* funcs.
var opcodeActions = map[lmsp.ProjectOpcode]string{
	"event_broadcast":                         "broadcast",
	"flippercontrol_stop":                     "stop",
	"flipperdisplay_centerButtonLight":        "setCenterButtonLight",
	"flipperdisplay_ledAnimation":             "startAnimation",
	"flipperdisplay_ledImage":                 "turnOnPixels",
	"flipperdisplay_ledImageFor":              "turnOnPixels",
	"flippermoremotor_motorSetDegreeCounted":  "setRelativePosition",
	"flippermoremotor_motorTurnForSpeed":      "runMotor",
	"flippermoremotor_position":               "relativePosition",
	"flippermoremove_moveDidMovement":         "wasMovementInterrupted",
	"flippermotor_absolutePosition":           "position",
	"flippermotor_motorGoDirectionToPosition": "goToPosition",
	"flippermotor_motorSetSpeed":              "setMotorSpeed",
	"flippermotor_motorStartDirection":        "motorStart",
	"flippermotor_motorStop":                  "stopMotor",
	"flippermotor_motorTurnForDirection":      "run",
	"flippermotor_speed":                      "motorSpeed",
	"flippermove_move":                        "move",
	"flippermove_movementSpeed":               "setMovementSpeed",
	"flippermove_setMovementPair":             "setMovementMotors",
	"flippermove_startSteer":                  "startMoving",
	"flippermove_steer":                       "move",
	"flippermove_stopMove":                    "stopMoving",
	"flippersensors_orientationAxis":          "angle",
	"flippersensors_resetYaw":                 "resetYaw",
	"flippersound_beep":                       "beep",
	"flippersound_playSound":                  "playSound",
	"flippersound_stopSound":                  "stopSound",
}

// renderAction visits a block that is like a function call. These may be script
// blocks, input blocks, or boolean blocks.
func renderAction(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, args ...argFn) {
	label, ok := opcodeActions[block.Opcode]
	if !ok {
		suggestf("  %q: \"todo\"\n", block.Opcode)
		label = string(block.Opcode)
	}
	fmt.Fprintf(w, "%s(", label)
	for i, a := range args {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		a(w, target, block)
	}
	fmt.Fprint(w, ")")
}

// This goes with the render* funcs.
func fieldArg(fieldName lmsp.ProjectFieldName) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprint(w, getField(block, fieldName))
	}
}

// This goes with the render* funcs.
func namedFieldArg(fieldName lmsp.ProjectFieldName) argFn {
	label := strings.ToLower(string(fieldName))
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprintf(w, "%s: %v", label, getField(block, fieldName))
	}
}

// This goes with the render* funcs.
func inputArg(inputName lmsp.ProjectInputID) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		visitInput(w, target, block, inputName)
	}
}

// This goes with the render* funcs.
var inputLabelOverrides = map[lmsp.ProjectOpcode]map[lmsp.ProjectInputID]string{
	"flipperdisplay_ledImageFor": {
		"VALUE": "seconds",
	},
	"flippermove_movementSpeed": {
		"SPEED": "percent",
	},
}

// This goes with the render* funcs.
func namedInputArg(inputName lmsp.ProjectInputID) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		label := inputLabel(block.Opcode, inputName)
		fmt.Fprintf(w, "%s: ", label)
		visitInput(w, target, block, inputName)
	}
}

// This goes with the render* funcs.
func inputLabel(opcode lmsp.ProjectOpcode, input lmsp.ProjectInputID) string {
	opcodeOverrides := inputLabelOverrides[opcode]
	if opcodeOverrides != nil {
		if label, ok := opcodeOverrides[input]; ok {
			return label
		}
	}
	return strings.ToLower(string(input))
}

// This goes with the render* funcs.
func fieldInputArg(fieldName lmsp.ProjectFieldName, inputName lmsp.ProjectInputID) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprintf(w, "%s: ", getField(block, fieldName))
		visitInput(w, target, block, inputName)
	}
}

func renderWhenBroadcastReceived(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when I receive %q:", getField(block, "BROADCAST_OPTION"))
}

func renderWhenProgramStarts(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "when program starts:")
}

func renderWhenPressed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "[port ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, "] when %s:", getField(block, "OPTION"))
}

func renderFieldSelector(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, field lmsp.ProjectFieldName) {
	fmt.Fprint(w, getField(block, field))
}

func renderIsReflectivity(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "(reflectivity(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, ") %s ", getField(block, "COMPARATOR"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w, ")")
}

func renderForever(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintln(w, "forever:")
	visitInput(indent(w), target, block, "SUBSTACK")
}

func renderControl(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, keyword string) {
	fmt.Fprintf(w, "%s ", keyword)
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w, ":")
	visitInput(indent(w), target, block, "SUBSTACK")
}

func renderIfElse(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	renderControl(w, target, block, "if")
	fmt.Fprintln(w, "\nelse:")
	visitInput(indent(w), target, block, "SUBSTACK2")
}

func renderWaitUntil(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "wait until ")
	visitInput(w, target, block, "CONDITION")
	fmt.Fprint(w)
}

func renderWait(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "wait(duration: ")
	visitInput(w, target, block, "DURATION")
	fmt.Fprint(w, ")")
}

func renderChangeVariableBy(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	f := getField(block, "VARIABLE")
	fmt.Fprintf(w, "[variable %s] = [variable %s] + ", f, f)
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w)
}

func renderSetVariableTo(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "[variable %s] = ", getField(block, "VARIABLE"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w)
}

// TODO - make 'op' a lookup
func renderBinaryOperator(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, op string, arg1, arg2 lmsp.ProjectInputID) {
	fmt.Fprint(w, "(")
	visitInput(w, target, block, arg1)
	fmt.Fprintf(w, " %s ", op)
	visitInput(w, target, block, arg2)
	fmt.Fprint(w, ")")
}

// TODO - make 'op' a lookup
func renderUnaryOperator(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, op string, arg lmsp.ProjectInputID) {
	fmt.Fprint(w, "not(")
	visitInput(w, target, block, arg)
	fmt.Fprint(w, ")")
}

func visitInput(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, inputName lmsp.ProjectInputID) {
	blockInput, ok := block.Inputs[inputName]
	if !ok {
		fmt.Fprintf(w, "[missing input: %q]", inputName) // TODO - move to a render func
		return
	}

	if blockInput == nil {
		fmt.Fprint(w, "[nil]") // TODO - move to a render func
		return
	}

	input := blockInput.([]interface{})

	// Input is [shadow, val, opt-val]
	// var shadowStrs = []string{"(unused)", "shadow", "no shadow", "shadow obscured"}
	// opt-val only shows up when shadow is 3 (shadow obscured).
	switch val := input[1].(type) {
	case string:
		visitBlock(w, target, lmsp.ProjectBlockID(val))
	case []interface{}:
		id := int(val[0].(float64))
		v := val[1].(string)
		switch id {
		case 4, 5, 6, 7, 8, 9:
			fmt.Fprint(w, v) // TODO - move this to a render* func
		case 10:
			fmt.Fprintf(w, "%q", v) // TODO - move this to a render* func
		case 11:
			fmt.Fprintf(w, "[broadcast %q]", v) // TODO - move this to a render* func
		case 12:
			fmt.Fprintf(w, "[variable %s]", v) // TODO - move this to a render* func
		case 13:
			fmt.Fprintf(w, "[list %q]", v) // TODO - move this to a render* func
		default:
			fmt.Fprintf(w, "???%#v???", val) // TODO - move this to a render* func
		}
	}
}

// This goes with the visit* funcs.
func getField(block *lmsp.ProjectBlockObject, name lmsp.ProjectFieldName) string {
	field := block.Fields[name].([]interface{})
	// field[0] could be a string, float64, maybe others?
	return fmt.Sprint(field[0])
}

var suggestionsEnabled = os.Getenv("SUGGEST") != ""

func suggestf(f string, arg ...interface{}) {
	if suggestionsEnabled {
		fmt.Fprintf(os.Stderr, f, arg...)
	}
}
