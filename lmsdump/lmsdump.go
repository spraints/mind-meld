package lmsdump

import (
	"fmt"
	"io"

	"github.com/spraints/mind-meld/lmsp"
)

func Dump(w io.Writer, proj lmsp.Project) {
	for _, target := range proj.Targets {
		fmt.Fprintf(w, "target: %s\n", target.Name)
		visitTarget(indent(w), target)
	}
}

func visitTarget(w io.Writer, target lmsp.ProjectTarget) {
	for _, id := range target.GetRootBlockIDs() {
		fmt.Fprintf(w, "----- %s -----\n", id)
		visitBlock(w, target, id)
	}
	// todo - other fields of ProjectTarget.
}

func visitBlock(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectBlockID) {
	block := target.Blocks[id].(*lmsp.ProjectBlockObject)
	switch block.Opcode {
	case "argument_reporter_string_number":
		visitFieldSelector(w, target, block, "VALUE")
	case "control_if":
		visitControl(w, target, block, "if")
	case "control_repeat_until":
		visitControl(w, target, block, "until")
	case "control_wait_until":
		visitWaitUntil(w, target, block)
	case "data_changevariableby":
		visitChangeVariableBy(w, target, block)
	case "data_setvariableto":
		visitSetVariableTo(w, target, block)
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
	case "flippermotor_motorGoDirectionToPosition":
		visitMotorGoDirectionToPosition(w, target, block)
	case "flippermotor_motorSetSpeed":
		visitMotorSetSpeed(w, target, block)
	case "flippermotor_motorStartDirection":
		visitMotorStartDirection(w, target, block)
	case "flippermotor_multiple-port-selector":
		visitFieldSelector(w, target, block, "field_flippermotor_multiple-port-selector")
	case "flippermove_movementSpeed":
		visitMovementSpeed(w, target, block)
	case "flippermove_rotation-wheel":
		visitFieldSelector(w, target, block, "field_flippermove_rotation-wheel")
	case "flippermove_startSteer":
		visitMoveStartSteer(w, target, block)
	case "flippermove_stopMove":
		visitMoveStopMove(w, target, block)
	case "flippersensors_orientationAxis":
		visitFieldSelector(w, target, block, "AXIS")
	case "flippersensors_resetYaw":
		fmt.Fprintln(w, "resetYaw()")
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
		fmt.Fprintf(w, "!!!TODO\n  case %q:\n    visitXYZ(w, target, block)\n!!!\n", block.Opcode)
		return
	}
	if block.Next != nil {
		visitBlock(w, target, *block.Next)
	}
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

func visitMoreMotorSetDegreeCounted(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "setDegreesCounted(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ", value: ")
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w, ")")
}

func visitMotorTurnForSpeed(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "motorTurnForSpeed(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ", speed: ")
	visitInput(w, target, block, "SPEED")
	fmt.Fprintf(w, ", %s: ", getField(block, "UNIT"))
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w, ")")
}

func visitMotorGoDirectionToPosition(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "goToPosition(port: ")
	visitInput(w, target, block, "PORT")
	fmt.Fprint(w, ", position: ")
	visitInput(w, target, block, "POSITION")
	fmt.Fprintf(w, ", direction: %s)\n", getField(block, "DIRECTION"))
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

func visitMoveStartSteer(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "startSteer(steering: ")
	visitInput(w, target, block, "STEERING")
	fmt.Fprintln(w, ")")
}

func visitMoveStopMove(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprintln(w, "stopMove()")
}

func visitControl(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, keyword string) {
	fmt.Fprintf(w, "%s ", keyword)
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w, ":")
	visitInput(indent(w), target, block, "SUBSTACK")
}

func visitWaitUntil(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "wait until: ")
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w)
}

func visitChangeVariableBy(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	f := getField(block, "VARIABLE")
	fmt.Fprintf(w, "%s = %s + ", f, f)
	visitInput(w, target, block, "VALUE")
	fmt.Fprintln(w)
}

func visitSetVariableTo(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	f := getField(block, "VARIABLE")
	fmt.Fprintf(w, "%s = ", f)
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
			fmt.Fprint(w, v)
		case 13:
			fmt.Fprintf(w, "[list %q]", v)
		default:
			fmt.Fprintf(w, "???%#v???", val)
		}
	}
}

func getField(block *lmsp.ProjectBlockObject, name lmsp.ProjectFieldName) string {
	field := block.Fields[name].([]interface{})
	return field[0].(string)
}
