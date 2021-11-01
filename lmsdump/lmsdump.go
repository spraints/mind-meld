package lmsdump

import (
	"fmt"
	"io"

	"github.com/spraints/mind-meld/lmsp"
)

func Dump(w io.Writer, proj lmsp.Project) {
	for _, target := range proj.Targets {
		fmt.Fprintf(w, "target: %s\n", target.Name)
		dumpTarget(w, target)
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
	case "operator_subtract":
		fmt.Fprintf(w, "%s%s Subtract %v - %v\n", prefix, id, describeInput(inputs["NUM1"]), describeInput(inputs["NUM2"]))
	case "procedures_definition":
		fmt.Fprintf(w, "%s%s Myblock\n", prefix, id)
		in := inputs["custom_block"].([]interface{})
		protoID := lmsp.ProjectBlockID(in[1].(string))
		dumpBlock(w, target, protoID, target.Blocks[protoID].(*lmsp.ProjectBlockObject), p)
	case "procedures_prototype":
		fmt.Fprintf(w, "%s%s Myblock prototype\n", prefix, id)
		// TODO
		fmt.Fprintf(w, "%s - inputs: %#v\n", p, inputs)
		fmt.Fprintf(w, "%s - mutation: %#v\n", p, block.Mutation)
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
		panic(len(s))
	}
}
