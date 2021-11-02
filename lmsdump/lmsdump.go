package lmsdump

import (
	"fmt"
	"io"

	"github.com/spraints/mind-meld/lmsp"
)

func Dump(w io.Writer, proj lmsp.Project) {
	for _, target := range proj.Targets {
		fmt.Fprintf(w, "target: %s\n", target.Name)
		visitTarget(w, target)
	}
}

func visitTarget(w io.Writer, target lmsp.ProjectTarget) {
	for _, id := range target.GetRootBlockIDs() {
		fmt.Fprintln(w, "-------------------")
		visitBlock(w, target, id)
	}
	// todo - other fields of ProjectTarget.
}

func visitBlock(w io.Writer, target lmsp.ProjectTarget, id lmsp.ProjectBlockID) {
	block := target.Blocks[id].(*lmsp.ProjectBlockObject)
	switch block.Opcode {
	case "flippermoremotor_motorSetDegreeCounted":
		visitMoreMotorSetDegreeCounted(w, target, block)
	case "flippermoremotor_multiple-port-selector":
		visitMoreMotorMultiplePortSelector(w, target, block)
	case "flippersensors_resetYaw":
		fmt.Fprintln(w, "resetYaw()")
	case "procedures_call":
		visitProcedureCall(w, target, block)
	case "procedures_definition":
		visitProcedureDefinition(w, target, block)
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
	inputs := block.Inputs.(map[string]interface{})
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

func visitMoreMotorMultiplePortSelector(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject) {
	fmt.Fprint(w, getField(block, "field_flippermoremotor_multiple-port-selector"))
}

func visitInput(w io.Writer, target lmsp.ProjectTarget, block *lmsp.ProjectBlockObject, inputName string) {
	inputs := block.Inputs.(map[string]interface{})
	input := inputs[inputName].([]interface{})
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
	inputs := block.Inputs.(map[string]interface{})
	switch block.Opcode {
	case "argument_reporter_string_number":
		fmt.Fprintf(w, "%s%q\n", prefix, getField(block, "VALUE"))
	case "control_repeat_until":
		fmt.Fprintf(w, "%s%s Repeat\n", prefix, id)
		fmt.Fprintf(w, "%s until:\n", p)
		dumpInput(w, target, inputs, "CONDITION", p+"   ? ")
		fmt.Fprintf(w, "%s do\n", p)
		dumpInputChain(w, target, inputs, "SUBSTACK", p+"  >> ", p+"   + ")
		fmt.Fprintf(w, "%s done\n", p)
	case "flippermoremotor_motorSetDegreeCounted":
		fmt.Fprintf(w, "%s%s Set Degree Counted\n", prefix, id)
		fmt.Fprintf(w, "%s value: %v\n", p, describeInput(inputs["VALUE"]))
		dumpInput(w, target, inputs, "PORT", p)
	case "flippermoremotor_multiple-port-selector":
		fmt.Fprintf(w, "%s port: %v\n", prefix, getField(block, "field_flippermoremotor_multiple-port-selector"))
	case "flippersensors_resetYaw":
		fmt.Fprintf(w, "%s%s Reset yaw\n", prefix, id)
	case "operator_lt":
		fmt.Fprintf(w, "%s%s Compare { %v < %v }\n", prefix, id, describeInput(inputs["OPERAND1"]), describeInput(inputs["OPERAND2"]))
	case "operator_subtract":
		fmt.Fprintf(w, "%s%s Subtract { %v - %v }\n", prefix, id, describeInput(inputs["NUM1"]), describeInput(inputs["NUM2"]))
	case "procedures_call":
		fmt.Fprintf(w, "%s%s Call %v\n", prefix, id, block.Mutation.ProcCode)
	case "procedures_definition":
		fmt.Fprintf(w, "%s%s Myblock\n", prefix, id)
		dumpInput(w, target, inputs, "custom_block", p)
	case "procedures_prototype":
		fmt.Fprintf(w, "%s%s Myblock prototype %q\n", prefix, id, block.Mutation.ProcCode)
		for id := range inputs {
			dumpInput(w, target, inputs, id, p+" -> "+id+" ")
		}
	default:
		fmt.Fprintf(w, "%s%s %s\n", prefix, id, block.Opcode)
		for name, input := range inputs {
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

func dumpInput(w io.Writer, target lmsp.ProjectTarget, inputs map[string]interface{}, name string, prefix string) {
	in := inputs[name].([]interface{})
	id := lmsp.ProjectBlockID(in[1].(string))
	dumpBlock(w, target, id, target.Blocks[id].(*lmsp.ProjectBlockObject), prefix)
}

func dumpInputChain(w io.Writer, target lmsp.ProjectTarget, inputs map[string]interface{}, name string, prefix, prefix2 string) {
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
