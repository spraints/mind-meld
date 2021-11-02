package lmsdump

import (
	"bytes"
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
	case "flippermoremotor_motorSetDegreeCounted":
		visitMoreMotorSetDegreeCounted(w, target, block)
	case "flippermoremotor_multiple-port-selector":
		visitFieldSelector(w, target, block, "field_flippermoremotor_multiple-port-selector")
	case "flippermoremotor_position":
		visitMoreMotorPosition(w, target, block)
	case "flippermoremotor_single-motor-selector":
		visitFieldSelector(w, target, block, "field_flippermoremotor_single-motor-selector")
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
	case "control_repeat_until":
		visitControlRepeatUntil(w, target, block)
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
		fmt.Fprintf(w, "!!!TODO %v!!!\n", block.Opcode)
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

func visitControlRepeatUntil(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, "until ")
	visitInput(w, target, block, "CONDITION")
	fmt.Fprintln(w, ":")
	visitInput(indent(w), target, block, "SUBSTACK")
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
			fmt.Fprintf(w, "[variable %q]", v)
		case 13:
			fmt.Fprintf(w, "[list %q]", v)
		default:
			fmt.Fprintf(w, "???%#v???", val)
		}
	}
}

func dumpTarget(w io.Writer, target lmsp.ProjectTarget) {
	for _, id := range target.GetRootBlockIDs() {
		dumpBlockChain(w, target, id, " >> ", "  + ")
	}
	for id, comment := range target.Comments {
		blockRef := ""
		if comment.BlockID != nil {
			blockRef = fmt.Sprintf(" (%s)", *comment.BlockID)
		}
		fmt.Fprintf(w, " \"\" %s %s%s\n", comment.Text, id, blockRef)
	}
}

func dumpBlockChain(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectBlockID, startPrefix, restPrefix string) {
	prefix := startPrefix
	for {
		block := target.Blocks[id].(*lmsp.ProjectBlockObject)
		dumpBlock(w, target, id, block, prefix)
		if block.Next == nil {
			break
		}
		prefix = restPrefix
		id = *block.Next
	}
}

func dumpBlock(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectBlockID, block *lmsp.ProjectBlockObject, prefix string) {
	p := pad(prefix)
	switch block.Opcode {
	case "argument_reporter_string_number":
		fmt.Fprintf(w, "%s%q\n", prefix, getField(block, "VALUE"))
	case "control_repeat_until":
		fmt.Fprintf(w, "%s%s Repeat\n", prefix, id)
		fmt.Fprintf(w, "%s until:\n", p)
		dumpInput(w, target, block.Inputs, "CONDITION", p+"   ? ")
		fmt.Fprintf(w, "%s do\n", p)
		dumpInputChain(w, target, block.Inputs, "SUBSTACK", p+"  >> ", p+"   + ")
		fmt.Fprintf(w, "%s done\n", p)
	case "flippermoremotor_motorSetDegreeCounted":
		fmt.Fprintf(w, "%s%s Set Degree Counted\n", prefix, id)
		fmt.Fprintf(w, "%s value: %v\n", p, describeInput(block.Inputs["VALUE"]))
		dumpInput(w, target, block.Inputs, "PORT", p)
	case "flippermoremotor_multiple-port-selector":
		fmt.Fprintf(w, "%s port: %v\n", prefix, getField(block, "field_flippermoremotor_multiple-port-selector"))
	case "flippersensors_resetYaw":
		fmt.Fprintf(w, "%s%s Reset yaw\n", prefix, id)
	case "operator_lt":
		fmt.Fprintf(w, "%s%s Compare { %v < %v }\n", prefix, id, describeInput(block.Inputs["OPERAND1"]), describeInput(block.Inputs["OPERAND2"]))
	case "operator_subtract":
		fmt.Fprintf(w, "%s%s Subtract { %v - %v }\n", prefix, id, describeInput(block.Inputs["NUM1"]), describeInput(block.Inputs["NUM2"]))
	case "procedures_call":
		fmt.Fprintf(w, "%s%s Call %v\n", prefix, id, block.Mutation.ProcCode)
	case "procedures_definition":
		fmt.Fprintf(w, "%s%s Myblock\n", prefix, id)
		dumpInput(w, target, block.Inputs, "custom_block", p)
	case "procedures_prototype":
		fmt.Fprintf(w, "%s%s Myblock prototype %q\n", prefix, id, block.Mutation.ProcCode)
		for id := range block.Inputs {
			dumpInput(w, target, block.Inputs, id, p+" -> "+string(id)+" ")
		}
	default:
		fmt.Fprintf(w, "%s%s %s\n", prefix, id, block.Opcode)
		for name, input := range block.Inputs {
			fmt.Fprintf(w, "%s - %s = %s\n", p, name, describeInput(input))
		}
		for name, field := range block.Fields {
			fmt.Fprintf(w, "%s field %s = %#v\n", p, name, field)
		}
	}
}

func getField(block *lmsp.ProjectBlockObject, name lmsp.ProjectFieldName) string {
	field := block.Fields[name].([]interface{})
	return field[0].(string)
}

func dumpInput(w io.Writer, target lmsp.ProjectTarget, inputs map[lmsp.ProjectInputID]lmsp.TODO, name lmsp.ProjectInputID, prefix string) {
	in := inputs[name].([]interface{})
	id := lmsp.ProjectBlockID(in[1].(string))
	dumpBlock(w, target, id, target.Blocks[id].(*lmsp.ProjectBlockObject), prefix)
}

func dumpInputChain(w io.Writer, target lmsp.ProjectTarget, inputs map[lmsp.ProjectInputID]lmsp.TODO, name lmsp.ProjectInputID, prefix, prefix2 string) {
	in := inputs[name].([]interface{})
	id := lmsp.ProjectBlockID(in[1].(string))
	dumpBlockChain(w, target, id, prefix, prefix2)
}

var shadowStrs = []string{"(unused)", "shadow", "no shadow", "shadow obscured"}

func describeInput(input interface{}) string {
	v := input.([]interface{})
	shadow := int(v[0].(float64))
	val := describeValue(v[1])
	return fmt.Sprintf("[%d-%s] %s", shadow, shadowStrs[shadow], val)
}

func describeValue(val interface{}) string {
	switch val := val.(type) {
	case []interface{}:
		id := int(val[0].(float64))
		v := val[1].(string)
		switch id {
		case 4, 5, 6, 7, 8, 9:
			return v
		case 10:
			return fmt.Sprintf("%q", v)
		case 11:
			return fmt.Sprintf("broadcast %q", v)
		case 12:
			return fmt.Sprintf("variable %q", v)
		case 13:
			return fmt.Sprintf("list %q", v)
		default:
			return fmt.Sprintf("???%#v???", val)
		}
	case string:
		return fmt.Sprintf("TODO ID of input: %q", val)
	default:
		return fmt.Sprintf("??2??%#v??", val)
	}
}

func pad(s string) string {
	switch len(s) {
	case 4:
		return "    "
	default:
		return " " + pad(s[1:])
	}
}

func indent(w io.Writer) io.Writer {
	return &indented{w: w, i: true}
}

type indented struct {
	w io.Writer
	i bool
}

var indentPadding = []byte("  ")

func (i *indented) Write(p []byte) (int, error) {
	n := 0
	for len(p) > 0 {
		if i.i {
			_, err := i.w.Write(indentPadding)
			if err != nil {
				return n, err
			}
			i.i = false
		}
		x := bytes.IndexRune(p, '\n')
		if x == -1 {
			nn, err := i.w.Write(p)
			n += nn
			return n, err
		}
		nn, err := i.w.Write(p[:x+1])
		n += nn
		if err != nil {
			return n, err
		}
		i.i = true
		p = p[x+1:]
	}
	return n, nil
}
