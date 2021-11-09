package lmsdump

import (
	"fmt"
	"io"
	"strings"

	"github.com/spraints/mind-meld/lmsp"
)

func Dump(w io.Writer, proj lmsp.Project) error {
	for _, target := range proj.Targets {
		if _, err := fmt.Fprintf(w, "target: %s\n", target.Name); err != nil {
			return err
		}
		// TODO check for errors in visitTarget.
		visitTarget(indent(w), target)
	}
	return nil
}

func visitTarget(w io.Writer, target lmsp.ProjectTarget) {
	for _, id := range target.GetRootBlockIDs() {
		fmt.Fprintf(w, "----- %s -----\n", id)
		fw := finishWithNewline(w)
		visitBlock(fw, target, id)
		fw.Finish()
	}
	first := true
	for _, id := range target.GetStandaloneCommentIDs() {
		visitComment(w, target, id)
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
		visitComment(w, target, block.Comment)
	}
	switch block.Opcode {
	case "argument_reporter_string_number":
		visitFieldSelector(w, target, block, "VALUE")
	case "control_forever":
		visitForever(w, target, block)
	case "control_if":
		visitControl(w, target, block, "if")
	case "control_if_else":
		visitIfElse(w, target, block)
	case "control_repeat_until":
		visitControl(w, target, block, "until")
	case "control_wait_until":
		visitWaitUntil(w, target, block)
	case "control_wait":
		visitWait(w, target, block)
	case "data_changevariableby":
		visitChangeVariableBy(w, target, block)
	case "data_setvariableto":
		visitSetVariableTo(w, target, block)
	case "event_broadcast":
		visitAction(w, target, block, "broadcast", inputArg("BROADCAST_INPUT"))
	case "event_whenbroadcastreceived":
		visitWhenBroadcastReceived(w, target, block)
		w = indent(w)
	case "flippercontrol_stop":
		visitStop(w, target, block)
	case "flipperdisplay_centerButtonLight":
		visitDisplayCenterButtonLight(w, target, block)
	case "flipperdisplay_color-selector-vertical":
		visitFieldSelector(w, target, block, "field_flipperdisplay_color-selector-vertical")
	case "flipperdisplay_custom-animate-matrix":
		visitFieldSelector(w, target, block, "field_flipperdisplay_custom-animate-matrix")
	case "flipperdisplay_custom-matrix":
		visitFieldSelector(w, target, block, "field_flipperdisplay_custom-matrix")
	case "flipperdisplay_ledAnimation":
		visitLEDAnimation(w, target, block)
	case "flipperdisplay_ledImage":
		visitLEDImage(w, target, block)
	case "flipperdisplay_ledImageFor":
		visitLEDImageFor(w, target, block)
	case "flipperevents_force-sensor-selector":
		visitFieldSelector(w, target, block, "field_flipperevents_force-sensor-selector")
	case "flipperevents_whenPressed":
		visitWhenPressed(w, target, block)
		w = indent(w)
	case "flipperevents_whenProgramStarts":
		fmt.Fprintln(w, "when program starts:")
		w = indent(w)
	case "flippermoremotor_motorSetDegreeCounted":
		visitMoreMotorSetDegreeCounted(w, target, block)
	case "flippermoremotor_motorTurnForSpeed":
		visitMotorTurnForSpeed(w, target, block)
	case "flippermoremotor_multiple-port-selector":
		visitFieldSelector(w, target, block, "field_flippermoremotor_multiple-port-selector")
	case "flippermoremotor_position":
		visitMoreMotorPosition(w, target, block)
	case "flippermoremotor_single-motor-selector":
		visitFieldSelector(w, target, block, "field_flippermoremotor_single-motor-selector")
	case "flippermotor_absolutePosition":
		visitMotorAbsolutePosition(w, target, block)
	case "flippermotor_custom-angle":
		visitFieldSelector(w, target, block, "field_flippermotor_custom-angle")
	case "flippermotor_custom-icon-direction":
		visitFieldSelector(w, target, block, "field_flippermotor_custom-icon-direction")
	case "flippermotor_motorGoDirectionToPosition":
		visitMotorGoDirectionToPosition(w, target, block)
	case "flippermotor_motorSetSpeed":
		visitMotorSetSpeed(w, target, block)
	case "flippermotor_motorStartDirection":
		visitMotorStartDirection(w, target, block)
	case "flippermotor_multiple-port-selector":
		visitFieldSelector(w, target, block, "field_flippermotor_multiple-port-selector")
	case "flippermotor_motorStop":
		visitMotorStop(w, target, block)
	case "flippermotor_motorTurnForDirection":
		visitMotorTurnForDirection(w, target, block)
	case "flippermotor_single-motor-selector":
		visitFieldSelector(w, target, block, "field_flippermotor_single-motor-selector")
	case "flippermotor_speed":
		visitMotorSpeed(w, target, block)
	case "flippermove_custom-icon-direction":
		visitFieldSelector(w, target, block, "field_flippermove_custom-icon-direction")
	case "flippermove_move":
		visitMove(w, target, block)
	case "flippermove_movementSpeed":
		visitMovementSpeed(w, target, block)
	case "flippermove_movement-port-selector":
		visitFieldSelector(w, target, block, "field_flippermove_movement-port-selector")
	case "flippermove_rotation-wheel":
		visitFieldSelector(w, target, block, "field_flippermove_rotation-wheel")
	case "flippermove_setMovementPair":
		visitSetMovementPair(w, target, block)
	case "flippermove_startSteer":
		visitMoveStartSteer(w, target, block)
	case "flippermove_steer":
		visitMoveSteer(w, target, block)
	case "flippermove_stopMove":
		visitMoveStopMove(w, target, block)
	case "flippersensors_color-sensor-selector":
		visitFieldSelector(w, target, block, "field_flippersensors_color-sensor-selector")
	case "flippersensors_isReflectivity":
		visitIsReflectivity(w, target, block)
	case "flippersensors_orientationAxis":
		visitOrientationAxis(w, target, block)
	case "flippersensors_resetYaw":
		fmt.Fprintln(w, "resetYaw()")
	case "flippersound_beep":
		visitPlayBeep(w, target, block)
	case "flippersound_custom-piano":
		visitFieldSelector(w, target, block, "field_flippersound_custom-piano")
	case "flippersound_playSound":
		visitPlaySound(w, target, block)
	case "flippersound_sound-selector":
		visitFieldSelector(w, target, block, "field_flippersound_sound-selector")
	case "flippersound_stopSound":
		fmt.Fprintln(w, "stopSound()")
	case "operator_add":
		visitBinaryOperator(w, target, block, "+", "NUM1", "NUM2")
	case "operator_equals":
		visitBinaryOperator(w, target, block, "==", "OPERAND1", "OPERAND2")
	case "operator_gt":
		visitBinaryOperator(w, target, block, ">", "OPERAND1", "OPERAND2")
	case "operator_lt":
		visitBinaryOperator(w, target, block, "<", "OPERAND1", "OPERAND2")
	case "operator_multiply":
		visitBinaryOperator(w, target, block, "*", "NUM1", "NUM2")
	case "operator_subtract":
		visitBinaryOperator(w, target, block, "-", "NUM1", "NUM2")
	case "procedures_call":
		visitProcedureCall(w, target, block)
	case "procedures_definition":
		visitProcedureDefinition(w, target, block)
		w = indent(w)
	case "procedures_prototype":
		visitProcedurePrototype(w, target, block)
	default:
		fmt.Fprintf(w, "!!!TODO\n  case %q:\n    visitXYZ(w, target, block)\n  %#v!!!\n", block.Opcode, block)
		return
	}
	if block.Next != nil {
		visitBlock(w, target, *block.Next)
	}
}

func visitComment(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectCommentID) {
	fmt.Fprintf(w, "/****\n  %s\n****/\n", target.Comments[id].Text)
}

func visitProcedureCall(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
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
	fmt.Fprintf(w, ")\n")
}

func visitProcedureDefinition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "def ")
	visitInput(w, target, block, "custom_block")
}

func visitProcedurePrototype(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "%s %s\n", block.Mutation.ProcCode, block.Mutation.ArgumentNames)
	// Inputs is redundant with argument names.
}

type argFn func(io.Writer, lmsp.ProjectTarget, *lmsp.ProjectBlockObject)

func visitAction(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, action string, args ...argFn) {
	fmt.Fprintf(w, "%s(", action)
	for i, a := range args {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		a(w, target, block)
	}
	fmt.Fprintln(w, ")")
}

func fieldArg(fieldName lmsp.ProjectFieldName) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprint(w, getField(block, fieldName))
	}
}

func namedFieldArg(fieldName lmsp.ProjectFieldName) argFn {
	label := strings.ToLower(string(fieldName))
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprintf(w, "%s: %v", label, getField(block, fieldName))
	}
}

func inputArg(inputName lmsp.ProjectInputID) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		visitInput(w, target, block, inputName)
	}
}

func namedInputArg(inputName lmsp.ProjectInputID) argFn {
	return namedInputArg2(inputName, strings.ToLower(string(inputName)))
}

func namedInputArg2(inputName lmsp.ProjectInputID, label string) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprintf(w, "%s: ", label)
		visitInput(w, target, block, inputName)
	}
}

func fieldInputArg(fieldName lmsp.ProjectFieldName, inputName lmsp.ProjectInputID) argFn {
	return func(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
		fmt.Fprintf(w, "%s: ", getField(block, fieldName))
		visitInput(w, target, block, inputName)
	}
}

func visitStop(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "stop", fieldArg("STOP_OPTION"))
}

func visitDisplayCenterButtonLight(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "setCenterButtonLight", namedInputArg("COLOR"))
}

func visitLEDAnimation(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "startAnimation", namedInputArg("MATRIX"))
}

func visitLEDImage(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "turnOnPixels", namedInputArg("MATRIX"))
}

func visitLEDImageFor(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "turnOnPixels", namedInputArg("MATRIX"), namedInputArg2("VALUE", "seconds"))
}

func visitWhenBroadcastReceived(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "when I receive %q:\n", getField(block, "BROADCAST_OPTION"))
}

func visitWhenPressed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "[port ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, "] when %s:\n", getField(block, "OPTION"))
}

func visitMoreMotorSetDegreeCounted(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "setRelativePosition", namedInputArg("PORT"), namedInputArg("VALUE"))
}

func visitMotorTurnForSpeed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "runMotor", namedInputArg("PORT"), namedInputArg("SPEED"), fieldInputArg("UNIT", "VALUE"))
}

func visitMotorGoDirectionToPosition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitAction(w, target, block, "goToPosition", namedInputArg("PORT"), namedInputArg("POSITION"), namedFieldArg("DIRECTION"))
}

func visitMotorSetSpeed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "setMotorSpeed(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ", speed: ")
	visitInput(w, target, block, "SPEED")
	fmt.Fprintln(w, ")")
}

func visitMotorStartDirection(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorStartDirection(direction: ")
	visitInput(w, target, block, "DIRECTION")
	fmt.Fprint(w, ", port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintln(w, ")")
}

func visitMotorStop(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "stopMotor(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintln(w, ")")
}

func visitMotorTurnForDirection(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorTurnForDirection(direction: ")
	visitInput(w, target, block, "DIRECTION")
	fmt.Fprint(w, ", port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, ", %s: ", getField(block, "UNIT"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w, ")")
}

func visitMotorSpeed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorSpeed(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ")")
}

func visitMove(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "move(direction: ")
	visitInput(w, target, block, "DIRECTION")
	fmt.Fprintf(w, ", %s: ", getField(block, "UNIT"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w, ")")
}

func visitMovementSpeed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "setMovementSpeed(speed: ")
	visitInput(w, target, block, "SPEED")
	fmt.Fprintln(w, ")")
}

func visitFieldSelector(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, field lmsp.ProjectFieldName) {
	fmt.Fprint(w, getField(block, field))
}

func visitMoreMotorPosition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorPosition(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ")")
}

func visitMotorAbsolutePosition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorAbsolutePosition(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ")")
}

func visitSetMovementPair(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "setMovementPair(")
	visitInput(w, target, block, "PAIR")
	fmt.Fprintln(w, ")")
}

func visitMoveStartSteer(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "startMoveSteer(steering: ")
	visitInput(w, target, block, "STEERING")
	fmt.Fprintln(w, ")")
}

func visitMoveSteer(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "moveSteer(steering: ")
	visitInput(w, target, block, "STEERING")
	fmt.Fprintf(w, ", %s: ", getField(block, "UNIT"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w, ")")
}

func visitMoveStopMove(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintln(w, "stopMove()")
}

func visitOrientationAxis(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "orientation(%s)", getField(block, "AXIS"))
}

func visitPlayBeep(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "beep(note: ")
	visitInput(w, target, block, "NOTE")
	fmt.Fprintln(w, ")")
}

func visitPlaySound(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "playSound(sound: ")
	visitInput(w, target, block, "SOUND")
	fmt.Fprintln(w, ")")
}

func visitIsReflectivity(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "(reflectivity(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprintf(w, ") %s ", getField(block, "COMPARATOR"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprint(w, ")")
}

func visitForever(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintln(w, "forever:")
	visitInput(indent(w), target, block, "SUBSTACK")
}

func visitControl(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, keyword string) {
	fmt.Fprintf(w, "%s ", keyword)
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w, ":")
	visitInput(indent(w), target, block, "SUBSTACK")
}

func visitIfElse(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	visitControl(w, target, block, "if")
	fmt.Fprintln(w, "else:")
	visitInput(indent(w), target, block, "SUBSTACK2")
}

func visitWaitUntil(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "wait until: ")
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w)
}

func visitWait(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "wait(duration: ")
	visitInput(w, target, block, "DURATION")
	fmt.Fprintln(w, ")")
}

func visitChangeVariableBy(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	f := getField(block, "VARIABLE")
	fmt.Fprintf(w, "%s = %s + ", f, f)
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w)
}

func visitSetVariableTo(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintf(w, "[variable %s] = ", getField(block, "VARIABLE"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w)
}

func visitBinaryOperator(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, op string, arg1, arg2 lmsp.ProjectInputID) {
	fmt.Fprint(w, "(")
	visitInput(w, target, block, arg1)
	fmt.Fprintf(w, " %s ", op)
	visitInput(w, target, block, arg2)
	fmt.Fprint(w, ")")
}

func visitInput(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, inputName lmsp.ProjectInputID) {
	input := block.Inputs[inputName].([]interface{})
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
			fmt.Fprint(w, v)
		case 10:
			fmt.Fprintf(w, "%q", v)
		case 11:
			fmt.Fprintf(w, "[broadcast %q]", v)
		case 12:
			fmt.Fprintf(w, "[variable %s]", v)
		case 13:
			fmt.Fprintf(w, "[list %q]", v)
		default:
			fmt.Fprintf(w, "???%#v???", val)
		}
	}
}

func getField(block *lmsp.ProjectBlockObject, name lmsp.ProjectFieldName) string {
	field := block.Fields[name].([]interface{})
	// field[0] could be a string, float64, maybe others?
	return fmt.Sprint(field[0])
}
