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
		renderTarget(indentStartingNow(w), target)
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
	case "flipperdisplay_menu_ledMatrixIndex":
		renderMenu(w, target, block, "ledMatrixIndex")
	case "sensing_keyoptions":
		renderFieldSelector(w, target, block, "KEY_OPTION")
	case "flipperdisplay_color-selector-vertical",
		"flipperdisplay_custom-animate-matrix",
		"flipperdisplay_custom-icon-direction",
		"flipperdisplay_custom-matrix",
		"flipperdisplay_distance-sensor-selector",
		"flipperdisplay_led-selector",
		"flipperevents_color-selector",
		"flipperevents_color-sensor-selector",
		"flipperevents_distance-sensor-selector",
		"flipperevents_force-sensor-selector",
		"flippermoremotor_multiple-port-selector",
		"flippermoremotor_single-motor-selector",
		"flippermoremove_rotation-wheel",
		"flippermoresensors_color-sensor-selector",
		"flippermoresensors_force-sensor-selector",
		"flippermotor_custom-angle",
		"flippermotor_custom-icon-direction",
		"flippermotor_multiple-port-selector",
		"flippermotor_single-motor-selector",
		"flippermove_custom-icon-direction",
		"flippermove_movement-port-selector",
		"flippermove_rotation-wheel",
		"flippersensors_color-selector",
		"flippersensors_color-sensor-selector",
		"flippersensors_distance-sensor-selector",
		"flippersound_custom-piano",
		"flippersound_sound-selector",
		"radiobroadcast_broadcast-signal":
		renderFieldSelector(w, target, block, lmsp.ProjectFieldName("field_"+block.Opcode))

	case "control_forever":
		renderForever(w, target, block)
	case "control_if":
		renderControl(w, target, block, "if")
	case "control_if_else":
		renderIfElse(w, target, block)
	case "control_repeat":
		renderControlRepeat(w, target, block)
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
	case "event_broadcastandwait":
		renderAction(w, target, block, inputArg("BROADCAST_INPUT"))
	case "event_whenbroadcastreceived":
		renderWhenBroadcastReceived(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "event_whenkeypressed":
		renderWhenKeyPressed(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.

	case "flippercontrol_stop":
		renderAction(w, target, block, fieldArg("STOP_OPTION"))

	case "flipperdisplay_centerButtonLight":
		renderAction(w, target, block, namedInputArg("COLOR"))
	case "flipperdisplay_ledAnimation":
		renderAction(w, target, block, namedInputArg("MATRIX"))
	case "flipperdisplay_ledAnimationUntilDone":
		renderAction(w, target, block, namedInputArg("MATRIX"))
	case "flipperdisplay_ledImage":
		renderAction(w, target, block, namedInputArg("MATRIX"))
	case "flipperdisplay_ledImageFor":
		renderAction(w, target, block, namedInputArg("MATRIX"), namedInputArg("VALUE"))
	case "flipperdisplay_ledOn":
		renderAction(w, target, block, namedInputArg("BRIGHTNESS"), namedInputArg("X"), namedInputArg("Y"))
	case "flipperdisplay_ledRotateDirection":
		renderAction(w, target, block, namedInputArg("DIRECTION"))
	case "flipperdisplay_ledRotateOrientation":
		renderAction(w, target, block, namedInputArg("ORIENTATION"))
	case "flipperdisplay_ledSetBrightness":
		renderAction(w, target, block, namedInputArg("BRIGHTNESS"))
	case "flipperdisplay_ledText":
		renderAction(w, target, block, namedInputArg("TEXT"))
	case "flipperdisplay_menu_orientation":
		renderMenu(w, target, block, "orientation")
	case "flipperdisplay_ultrasonicLightUp":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("VALUE"))

	case "flipperevents_whenButton":
		renderWhenButton(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenColor":
		renderWhenColor(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenCondition":
		renderWhenCondition(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenDistance":
		renderWhenDistance(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenGesture":
		renderWhenGesture(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenOrientation":
		renderWhenOrientation(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenPressed":
		renderWhenPressed(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenProgramStarts":
		renderWhenProgramStarts(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.
	case "flipperevents_whenTimer":
		renderWhenTimer(w, target, block)
		w = indent(w) // TODO - move this to a renderX func.

	case "flippermoremotor_menu_acceleration":
		renderMenu(w, target, block, "acceleration")
	case "flippermoremotor_motorDidMovement":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermoremotor_motorGoToRelativePosition":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("POSITION"), namedInputArg("SPEED"))
	case "flippermoremotor_motorSetAcceleration":
		renderAction(w, target, block, namedInputArg("ACCELERATION"), namedInputArg("PORT"))
	case "flippermoremotor_motorSetDegreeCounted":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("VALUE"))
	case "flippermoremotor_motorSetStallDetection":
		renderAction(w, target, block, namedFieldArg("ENABLED"), namedInputArg("PORT"))
	case "flippermoremotor_motorSetStopMethod":
		renderAction(w, target, block, namedInputArg("PORT"), namedFieldArg("STOP"))
	case "flippermoremotor_motorStartPower":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("POWER"))
	case "flippermoremotor_motorStartSpeed":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("SPEED"))
	case "flippermoremotor_motorTurnForSpeed":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("SPEED"), fieldInputArg("UNIT", "VALUE"))
	case "flippermoremotor_position":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippermoremotor_power":
		renderAction(w, target, block, namedInputArg("PORT"))

	case "flippermoremove_menu_acceleration":
		renderMenu(w, target, block, "acceleration")
	case "flippermoremove_moveDistanceAtSpeed":
		renderAction(w, target, block, namedInputArg("LEFT"), namedInputArg("RIGHT"), namedInputArg("DISTANCE"), namedFieldArg("UNIT"))
	case "flippermoremove_movementSetAcceleration":
		renderAction(w, target, block, inputArg("ACCELERATION"))
	case "flippermoremove_movementSetStopMethod":
		renderAction(w, target, block, namedFieldArg("STOP"))
	case "flippermoremove_startDualPower":
		renderAction(w, target, block, namedInputArg("LEFT"), namedInputArg("RIGHT"))
	case "flippermoremove_startDualSpeed":
		renderAction(w, target, block, namedInputArg("LEFT"), namedInputArg("RIGHT"))
	case "flippermoremove_startSteerAtSpeed":
		renderAction(w, target, block, namedInputArg("STEERING"), namedInputArg("SPEED"))
	case "flippermoremove_steerDistanceAtSpeed":
		renderAction(w, target, block, namedInputArg("STEERING"), namedInputArg("SPEED"), namedInputArg("DISTANCE"), namedFieldArg("UNIT"))

	case "flippermoresensors_acceleration":
		renderAction(w, target, block, namedFieldArg("AXIS"))
	case "flippermoresensors_angularVelocity":
		renderAction(w, target, block, namedFieldArg("AXIS"))
	case "flippermoresensors_force":
		renderAction(w, target, block, namedInputArg("PORT"), namedFieldArg("UNIT"))
	case "flippermoresensors_isPressed":
		renderAction(w, target, block, namedInputArg("PORT"), namedFieldArg("OPTION"))
	case "flippermoresensors_rawColor":
		renderAction(w, target, block, namedFieldArg("COLOR"), namedInputArg("PORT"))

	case "flippermotor_absolutePosition":
		renderAction(w, target, block, namedInputArg("PORT"))
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
	case "flippermotor_speed":
		renderAction(w, target, block, namedInputArg("PORT"))

	case "flippermove_move":
		renderAction(w, target, block, namedInputArg("DIRECTION"), fieldInputArg("UNIT", "VALUE"))
	case "flippermove_movementSpeed":
		renderAction(w, target, block, namedInputArg("SPEED"))
	case "flippermove_setDistance":
		renderAction(w, target, block, namedInputArg("DISTANCE"), namedFieldArg("UNIT"))
	case "flippermove_setMovementPair":
		renderAction(w, target, block, namedInputArg("PAIR"))
	case "flippermove_startSteer":
		renderAction(w, target, block, namedInputArg("STEERING"))
	case "flippermove_steer":
		renderAction(w, target, block, namedInputArg("STEERING"), fieldInputArg("UNIT", "VALUE"))

	case "flipperoperator_isInBetween":
		renderBetween(w, target, block)

	case "flippersensors_buttonIsPressed":
		renderAction(w, target, block, namedFieldArg("BUTTON"), namedFieldArg("EVENT"))
	case "flippersensors_color":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippersensors_distance":
		renderAction(w, target, block, namedInputArg("PORT"), namedFieldArg("UNIT"))
	case "flippersensors_isColor":
		renderAction(w, target, block, namedInputArg("PORT"), namedInputArg("VALUE"))
	case "flippersensors_isDistance":
		renderAction(w, target, block, namedFieldArg("COMPARATOR"), namedInputArg("PORT"), namedFieldArg("UNIT"), namedInputArg("VALUE"))
	case "flippersensors_ismotion":
		renderAction(w, target, block, namedFieldArg("MOTION"))
	case "flippersensors_isorientation":
		renderAction(w, target, block, namedFieldArg("ORIENTATION"))
	case "flippersensors_isReflectivity":
		renderIsReflectivity(w, target, block)
	case "flippersensors_reflectivity":
		renderAction(w, target, block, namedInputArg("PORT"))
	case "flippersensors_orientationAxis":
		renderAction(w, target, block, fieldArg("AXIS"))

	case "flippersound_beep":
		renderAction(w, target, block, namedInputArg("NOTE"))
	case "flippersound_beepForTime":
		renderAction(w, target, block, namedInputArg("DURATION"), namedInputArg("NOTE"))
	case "flippersound_playSound":
		renderAction(w, target, block, namedInputArg("SOUND"))
	case "flippersound_playSoundUntilDone":
		renderAction(w, target, block, namedInputArg("SOUND"))
	case "sound_setvolumeto":
		renderAction(w, target, block, namedInputArg("VOLUME"))

	case "operator_add":
		renderBinaryOperator(w, target, block, "+", "NUM1", "NUM2")
	case "operator_and":
		renderBinaryOperator(w, target, block, "AND", "OPERAND1", "OPERAND2")
	case "operator_contains":
		renderAction(w, target, block, namedInputArg("STRING1"), namedInputArg("STRING2"))
	case "operator_divide":
		renderBinaryOperator(w, target, block, "/", "NUM1", "NUM2")
	case "operator_equals":
		renderBinaryOperator(w, target, block, "==", "OPERAND1", "OPERAND2")
	case "operator_gt":
		renderBinaryOperator(w, target, block, ">", "OPERAND1", "OPERAND2")
	case "operator_join":
		renderAction(w, target, block, inputArg("STRING1"), inputArg("STRING2"))
	case "operator_length":
		renderAction(w, target, block, namedInputArg("STRING"))
	case "operator_letter_of":
		renderAction(w, target, block, inputArg("STRING"), inputArg("LETTER"))
	case "operator_lt":
		renderBinaryOperator(w, target, block, "<", "OPERAND1", "OPERAND2")
	case "operator_mathop":
		renderMathOp(w, target, block)
	case "operator_mod":
		renderBinaryOperator(w, target, block, "mod", "NUM1", "NUM2")
	case "operator_multiply":
		renderBinaryOperator(w, target, block, "*", "NUM1", "NUM2")
	case "operator_not":
		renderUnaryOperator(w, target, block, "NOT", "OPERAND")
	case "operator_or":
		renderBinaryOperator(w, target, block, "OR", "OPERAND1", "OPERAND2")
	case "operator_random":
		renderAction(w, target, block, namedInputArg("FROM"), namedInputArg("TO"))
	case "operator_round":
		renderUnaryOperator(w, target, block, "round", "NUM")
	case "operator_subtract":
		renderBinaryOperator(w, target, block, "-", "NUM1", "NUM2")

	case "procedures_call":
		renderProcedureCall(w, target, block)
	case "procedures_definition":
		renderProcedureDefinition(w, target, block)
		w = indent(w) // TODO - move this to a 'renderX' func.
	case "procedures_prototype":
		renderProcedurePrototype(w, target, block)

	case "radiobroadcast_radioSignalReporter":
		renderAction(w, target, block, namedInputArg("SIGNAL"))
	case "radiobroadcast_broadcastRadioSignalWithValueCommand":
		renderAction(w, target, block, namedInputArg("SIGNAL"), namedInputArg("VALUE"))
	case "radiobroadcast_whenIReceiveRadioSignalHat":
		renderAction(w, target, block, namedInputArg("SIGNAL"))

	case "sensing_keypressed":
		renderAction(w, target, block, inputArg("KEY_OPTION"))

	case "sound_changeeffectby":
		renderAction(w, target, block, namedFieldArg("EFFECT"), namedInputArg("VALUE"))
	case "sound_changevolumeby":
		renderAction(w, target, block, namedInputArg("VOLUME"))
	case "sound_seteffectto":
		renderAction(w, target, block, namedFieldArg("EFFECT"), namedInputArg("VALUE"))

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
	"event_broadcast":        "broadcast",
	"event_broadcastandwait": "broadcastAndWait",

	"flippercontrol_stop":            "stop",
	"flippercontrol_stopOtherStacks": "stopOtherStacks",

	"flipperdisplay_centerButtonLight":     "setCenterButtonLight",
	"flipperdisplay_displayOff":            "turnOffPixels",
	"flipperdisplay_ledAnimation":          "startAnimation",
	"flipperdisplay_ledAnimationUntilDone": "playAnimationUntilDone",
	"flipperdisplay_ledImage":              "turnOnPixels",
	"flipperdisplay_ledImageFor":           "turnOnPixels",
	"flipperdisplay_ledOn":                 "setPixel",
	"flipperdisplay_ledRotateDirection":    "rotateDisplay",
	"flipperdisplay_ledRotateOrientation":  "setDisplayRotation",
	"flipperdisplay_ledSetBrightness":      "setPixelBrightness",
	"flipperdisplay_ledText":               "write",
	"flipperdisplay_ultrasonicLightUp":     "lightUpUltrasonicSensor",

	"flippermoremotor_motorDidMovement":          "wasMotorInterrupted",
	"flippermoremotor_motorGoToRelativePosition": "goToRelativePosition",
	"flippermoremotor_motorSetAcceleration":      "setAcceleration",
	"flippermoremotor_motorSetDegreeCounted":     "setRelativePosition",
	"flippermoremotor_motorSetStallDetection":    "setStallDetection",
	"flippermoremotor_motorSetStopMethod":        "setStopMethod",
	"flippermoremotor_motorStartPower":           "startMotor",
	"flippermoremotor_motorStartSpeed":           "startMotor",
	"flippermoremotor_motorTurnForSpeed":         "runMotor",
	"flippermoremotor_position":                  "relativePosition",
	"flippermoremotor_power":                     "motorPower",

	"flippermoremove_moveDidMovement":         "wasMovementInterrupted",
	"flippermoremove_moveDistanceAtSpeed":     "moveAtSpeed",
	"flippermoremove_movementSetAcceleration": "setMovementAcceleration",
	"flippermoremove_movementSetStopMethod":   "setMovementStopMethod",
	"flippermoremove_startDualPower":          "startMovingAtPower",
	"flippermoremove_startDualSpeed":          "startMovingAtSpeed",
	"flippermoremove_startSteerAtSpeed":       "startMovingAtSpeed",
	"flippermoremove_steerDistanceAtSpeed":    "move",

	"flippermoresensors_acceleration":    "acceleration",
	"flippermoresensors_angularVelocity": "angularVelocity",
	"flippermoresensors_force":           "pressure",
	"flippermoresensors_isPressed":       "isPressed",
	"flippermoresensors_rawColor":        "rawColor",

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
		suggestf("  %q: \"todo--%s\",\n", block.Opcode, block.Opcode)
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
		value := getField(block, fieldName)
		if block.Opcode == "flippermoremove_movementSetStopMethod" && fieldName == "STOP" {
			switch value {
			case "1":
				value = "brake"
			case "2":
				value = "hold position"
			case "3":
				value = "coast"
			}
		}
		fmt.Fprintf(w, "%s: %v", label, value)
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

func renderWhenKeyPressed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when %q key pressed:", getField(block, "KEY_OPTION"))
}

func renderWhenProgramStarts(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "when program starts:")
}

func renderWhenTimer(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "when timer > ")
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w, ":")
}

func renderWhenButton(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when %q button %q:", getField(block, "BUTTON"), getField(block, "EVENT"))
}

func renderWhenColor(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "[port ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, "] when color is ")
	visitInput(w, target, block, "OPTION")
	fmt.Fprint(w, ":")
}

func renderWhenCondition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "when ")
	visitInput(w, target, block, "CONDITION")
	fmt.Fprint(w, ":")
}

func renderWhenDistance(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "[port ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, "] when distance %s ", getField(block, "COMPARATOR"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintf(w, " %s:", getField(block, "UNIT"))
}

func renderWhenGesture(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when gesture %q occurs:", getField(block, "EVENT"))
}

func renderWhenOrientation(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when %q is up:", getField(block, "VALUE"))
}

func renderWhenPressed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "[port ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, "] when %s:", getField(block, "OPTION"))
}

func renderFieldSelector(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, field lmsp.ProjectFieldName) {
	fmt.Fprint(w, getField(block, field))
}

var menus = map[lmsp.ProjectOpcode]map[string]string{

	"flipperdisplay_menu_ledMatrixIndex": {
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
	},
	"flipperdisplay_menu_orientation": {
		"1": "upright",
		"2": "left",
		"3": "right",
		"4": "upside down",
	},
	"flippermoremotor_menu_acceleration": {
		"-1 -1":     "default",
		"100 100":   "fast",
		"350 350":   "balanced",
		"800 800":   "smooth",
		"1200 1200": "slow",
		"2000 2000": "very slow",
	},
	"flippermoremove_menu_acceleration": {
		"-1 -1":     "default",
		"100 100":   "fast",
		"350 350":   "balanced",
		"800 800":   "smooth",
		"1200 1200": "slow",
		"2000 2000": "very slow",
	},
}

func renderMenu(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, field lmsp.ProjectFieldName) {
	value := getField(block, field)

	menu, ok := menus[block.Opcode]
	if !ok {
		fmt.Fprintf(w, "[unknown menu type %s %q]", block.Opcode, value)
		return
	}

	menuValue, ok := menu[value]
	if !ok {
		fmt.Fprintf(w, "[%s: %q]", block.Opcode, value)
		return
	}

	fmt.Fprint(w, menuValue)
}

func renderBetween(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "(")
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w, " between [")
	visitInput(w, target, block, "LOW")
	fmt.Fprint(w, " and ")
	visitInput(w, target, block, "HIGH")
	fmt.Fprint(w, "])")
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
	visitInput(indentStartingNow(w), target, block, "SUBSTACK")
}

func renderControlRepeat(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "repeat ")
	visitInput(w, target, block, "TIMES")
	fmt.Fprintln(w, " times:")
	visitInput(indentStartingNow(w), target, block, "SUBSTACK")
}

func renderControl(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, keyword string) {
	fmt.Fprintf(w, "%s ", keyword)
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w, ":")
	visitInput(indentStartingNow(w), target, block, "SUBSTACK")
}

func renderIfElse(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	renderControl(w, target, block, "if")
	fmt.Fprintln(w, "\nelse:")
	visitInput(indentStartingNow(w), target, block, "SUBSTACK2")
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
	fmt.Fprintf(w, "%s(", op)
	visitInput(w, target, block, arg)
	fmt.Fprint(w, ")")
}

func renderMathOp(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "math.%s(", getField(block, "OPERATOR"))
	visitInput(w, target, block, "NUM")
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
			if v == "" {
				fmt.Fprint(w, "[unset number]") // TODO - move this to a render* func
			} else {
				fmt.Fprint(w, v) // TODO - move this to a render* func
			}
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
