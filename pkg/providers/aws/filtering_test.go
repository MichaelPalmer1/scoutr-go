package aws_test

import (
	"fmt"
	"testing"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func TestOperations(t *testing.T) {
	f := aws.NewFilter()

	operationMap := f.Operations()

	if len(operationMap) != 13 {
		t.Errorf("Expected 13 operations, but got %d", len(operationMap))
	}
}

func TestBuildInExpr(t *testing.T) {
	f := aws.NewFilter()

	var values []string
	for i := 0; i < 327; i++ {
		values = append(values, fmt.Sprintf("%d", i))
	}

	conditions := f.BuildInExpr("key", values, false)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(expr.Names()) != 1 {
		t.Errorf("Expected 1 name, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != len(values) {
		t.Errorf("Expected %d values, but got %d", len(values), len(expr.Values()))
	}

	if *expr.Filter() != "(((#0 IN (:0, :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22, :23, :24, :25, :26, :27, :28, :29, :30, :31, :32, :33, :34, :35, :36, :37, :38, :39, :40, :41, :42, :43, :44, :45, :46, :47, :48, :49, :50, :51, :52, :53, :54, :55, :56, :57, :58, :59, :60, :61, :62, :63, :64, :65, :66, :67, :68, :69, :70, :71, :72, :73, :74, :75, :76, :77, :78, :79, :80, :81, :82, :83, :84, :85, :86, :87, :88, :89, :90, :91, :92, :93, :94, :95, :96, :97, :98, :99)) OR (#0 IN (:100, :101, :102, :103, :104, :105, :106, :107, :108, :109, :110, :111, :112, :113, :114, :115, :116, :117, :118, :119, :120, :121, :122, :123, :124, :125, :126, :127, :128, :129, :130, :131, :132, :133, :134, :135, :136, :137, :138, :139, :140, :141, :142, :143, :144, :145, :146, :147, :148, :149, :150, :151, :152, :153, :154, :155, :156, :157, :158, :159, :160, :161, :162, :163, :164, :165, :166, :167, :168, :169, :170, :171, :172, :173, :174, :175, :176, :177, :178, :179, :180, :181, :182, :183, :184, :185, :186, :187, :188, :189, :190, :191, :192, :193, :194, :195, :196, :197, :198, :199))) OR (#0 IN (:200, :201, :202, :203, :204, :205, :206, :207, :208, :209, :210, :211, :212, :213, :214, :215, :216, :217, :218, :219, :220, :221, :222, :223, :224, :225, :226, :227, :228, :229, :230, :231, :232, :233, :234, :235, :236, :237, :238, :239, :240, :241, :242, :243, :244, :245, :246, :247, :248, :249, :250, :251, :252, :253, :254, :255, :256, :257, :258, :259, :260, :261, :262, :263, :264, :265, :266, :267, :268, :269, :270, :271, :272, :273, :274, :275, :276, :277, :278, :279, :280, :281, :282, :283, :284, :285, :286, :287, :288, :289, :290, :291, :292, :293, :294, :295, :296, :297, :298, :299))) OR (#0 IN (:300, :301, :302, :303, :304, :305, :306, :307, :308, :309, :310, :311, :312, :313, :314, :315, :316, :317, :318, :319, :320, :321, :322, :323, :324, :325, :326))" {
		t.Error("Filter expression is incorrect")
		t.Logf("Filter Expression: %+v", *expr.Filter())
	}

}

func TestBuildInExprSingleValue(t *testing.T) {
	f := aws.NewFilter()
	conditions := f.BuildInExpr("key", []string{"1"}, false)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected 1 value, but got %d value(s)", len(expr.Values()))
	}
}

func TestBuildInExprEmpty(t *testing.T) {
	f := aws.NewFilter()
	conditions := f.BuildInExpr("key", nil, false)

	if conditions.IsSet() {
		t.Error("Conditions should not be set")
	}
}

func TestBuildInExprNegated(t *testing.T) {
	f := aws.NewFilter()

	var values []string
	for i := 0; i < 327; i++ {
		values = append(values, fmt.Sprintf("%d", i))
	}

	conditions := f.BuildInExpr("key", values, true)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(expr.Names()) != 1 {
		t.Errorf("Expected 1 name, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != len(values) {
		t.Errorf("Expected %d values, but got %d", len(values), len(expr.Values()))
	}

	if *expr.Filter() != "(((NOT (#0 IN (:0, :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22, :23, :24, :25, :26, :27, :28, :29, :30, :31, :32, :33, :34, :35, :36, :37, :38, :39, :40, :41, :42, :43, :44, :45, :46, :47, :48, :49, :50, :51, :52, :53, :54, :55, :56, :57, :58, :59, :60, :61, :62, :63, :64, :65, :66, :67, :68, :69, :70, :71, :72, :73, :74, :75, :76, :77, :78, :79, :80, :81, :82, :83, :84, :85, :86, :87, :88, :89, :90, :91, :92, :93, :94, :95, :96, :97, :98, :99))) AND (NOT (#0 IN (:100, :101, :102, :103, :104, :105, :106, :107, :108, :109, :110, :111, :112, :113, :114, :115, :116, :117, :118, :119, :120, :121, :122, :123, :124, :125, :126, :127, :128, :129, :130, :131, :132, :133, :134, :135, :136, :137, :138, :139, :140, :141, :142, :143, :144, :145, :146, :147, :148, :149, :150, :151, :152, :153, :154, :155, :156, :157, :158, :159, :160, :161, :162, :163, :164, :165, :166, :167, :168, :169, :170, :171, :172, :173, :174, :175, :176, :177, :178, :179, :180, :181, :182, :183, :184, :185, :186, :187, :188, :189, :190, :191, :192, :193, :194, :195, :196, :197, :198, :199)))) AND (NOT (#0 IN (:200, :201, :202, :203, :204, :205, :206, :207, :208, :209, :210, :211, :212, :213, :214, :215, :216, :217, :218, :219, :220, :221, :222, :223, :224, :225, :226, :227, :228, :229, :230, :231, :232, :233, :234, :235, :236, :237, :238, :239, :240, :241, :242, :243, :244, :245, :246, :247, :248, :249, :250, :251, :252, :253, :254, :255, :256, :257, :258, :259, :260, :261, :262, :263, :264, :265, :266, :267, :268, :269, :270, :271, :272, :273, :274, :275, :276, :277, :278, :279, :280, :281, :282, :283, :284, :285, :286, :287, :288, :289, :290, :291, :292, :293, :294, :295, :296, :297, :298, :299)))) AND (NOT (#0 IN (:300, :301, :302, :303, :304, :305, :306, :307, :308, :309, :310, :311, :312, :313, :314, :315, :316, :317, :318, :319, :320, :321, :322, :323, :324, :325, :326)))" {
		t.Error("Filter expression is incorrect")
		t.Logf("Filter Expression: %+v", *expr.Filter())
	}
}

func TestBuildInExprSingleValueNegated(t *testing.T) {
	f := aws.NewFilter()
	conditions := f.BuildInExpr("key", []string{"1"}, true)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected 1 value, but got %d value(s)", len(expr.Values()))
	}
}

func TestBuildInExprEmptyNegated(t *testing.T) {
	f := aws.NewFilter()
	conditions := f.BuildInExpr("key", nil, true)

	if conditions.IsSet() {
		t.Error("Conditions should not be set")
	}
}

func TestAnd(t *testing.T) {
	f := aws.NewFilter()

	cond1 := expression.Name("key1").Equal(expression.Value("val1"))
	cond2 := expression.Name("key2").Equal(expression.Value("val2"))

	conditions := f.And(cond1, cond2).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "(#0 = :0) AND (#1 = :1)" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestAndCond1(t *testing.T) {
	f := aws.NewFilter()

	cond1 := expression.Name("key1").Equal(expression.Value("val1"))

	conditions := f.And(cond1, nil).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "#0 = :0" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestAndCond2(t *testing.T) {
	f := aws.NewFilter()

	cond2 := expression.Name("key2").Equal(expression.Value("val2"))

	conditions := f.And(nil, cond2).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "#0 = :0" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestAndCondsNil(t *testing.T) {
	f := aws.NewFilter()

	conditions := f.And(nil, nil).(expression.ConditionBuilder)

	if conditions.IsSet() {
		t.Error("Conditions should be not be set")
	}
}

func TestOr(t *testing.T) {
	f := aws.NewFilter()

	cond1 := expression.Name("key1").Equal(expression.Value("val1"))
	cond2 := expression.Name("key2").Equal(expression.Value("val2"))

	conditions := f.Or(cond1, cond2).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "(#0 = :0) OR (#1 = :1)" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestOrCond1(t *testing.T) {
	f := aws.NewFilter()

	cond1 := expression.Name("key1").Equal(expression.Value("val1"))

	conditions := f.Or(cond1, nil).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "#0 = :0" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestOrCond2(t *testing.T) {
	f := aws.NewFilter()

	cond2 := expression.Name("key2").Equal(expression.Value("val2"))

	conditions := f.Or(nil, cond2).(expression.ConditionBuilder)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "#0 = :0" {
		t.Errorf("Incorrect filter expression")
		t.Logf("Filter expression: %+v", *expr.Filter())
	}
}

func TestOrCondsNil(t *testing.T) {
	f := aws.NewFilter()

	conditions := f.Or(nil, nil).(expression.ConditionBuilder)

	if conditions.IsSet() {
		t.Error("Conditions should be not be set")
	}
}

func TestEquals(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Equals("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 = :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 = :0', got '%s'", *expr.Filter())
	}
}

func TestNotEqual(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.NotEqual("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 <> :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 <> :0', got '%s'", *expr.Filter())
	}
}

func TestGreaterThan(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.GreaterThan("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 > :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 > :0', got '%s'", *expr.Filter())
	}
}

func TestGreaterThanEqual(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.GreaterThanEqual("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 >= :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 >= :0', got '%s'", *expr.Filter())
	}
}

func TestLessThan(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.LessThan("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 < :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 < :0', got '%s'", *expr.Filter())
	}
}

func TestLessThanEqual(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.LessThanEqual("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 <= :0" {
		t.Errorf("Incorrect filter expression. Expected '#0 <= :0', got '%s'", *expr.Filter())
	}
}

func TestStartsWith(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.StartsWith("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "begins_with (#0, :0)" {
		t.Errorf("Incorrect filter expression. Expected 'begins_with (#0, :0)', got '%s'", *expr.Filter())
	}
}

func TestContains(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Contains("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "contains (#0, :0)" {
		t.Errorf("Incorrect filter expression. Expected 'contains (#0, :0)', got '%s'", *expr.Filter())
	}
}

func TestNotContains(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.NotContains("key", "value123")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "NOT (contains (#0, :0))" {
		t.Errorf("Incorrect filter expression. Expected 'NOT (contains (#0, :0))', got '%s'", *expr.Filter())
	}
}

func TestExistsTrue(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Exists("key", "true")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "attribute_exists (#0)" {
		t.Errorf("Incorrect filter expression. Expected 'attribute_exists (#0)', got '%s'", *expr.Filter())
	}
}

func TestExistsFalse(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Exists("key", "false")
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "attribute_not_exists (#0)" {
		t.Errorf("Incorrect filter expression. Expected 'attribute_not_exists (#0)', got '%s'", *expr.Filter())
	}
}

func TestExistsOther(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Exists("key", "blah")

	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	} else if err.Error() != "invalid value for Exists operation. Supported values are ['true'/'false']" {
		t.Error(err)
	}
}

func TestBetween(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Between("key", `["1", "2"]`)
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 BETWEEN :0 AND :1" {
		t.Errorf("Incorrect filter expression. Expected '#0 BETWEEN :0 AND :1', got '%s'", *expr.Filter())
	}
}

func TestBetweenInvalidValue(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Between("key", false)
	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	} else if err.Error() != "false could not be cast as a string" {
		t.Error(err)
	}
}

func TestBetweenUnmarshalError(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.Between("key", `{`)

	if conds != nil {
		t.Error("Conditions should be nil")
	}
	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestIn(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.In("key", `["1", "2"]`)
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "#0 IN (:0, :1)" {
		t.Errorf("Incorrect filter expression. Expected '#0 IN (:0, :1)', got '%s'", *expr.Filter())
	}
}

func TestInInvalidValue(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.In("key", false)

	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	} else if err.Error() != "false could not be cast as a string" {
		t.Error(err)
	}
}

func TestInJsonUnmarshalError(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.In("key", `[`)

	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestNotIn(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.NotIn("key", `["1", "2"]`)
	if err != nil {
		t.Error(err)
	}

	conditions := conds.(expression.ConditionBuilder)
	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if *expr.Filter() != "NOT (#0 IN (:0, :1))" {
		t.Errorf("Incorrect filter expression. Expected 'NOT (#0 IN (:0, :1))', got '%s'", *expr.Filter())
	}
}

func TestNotInInvalidValue(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.NotIn("key", false)

	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	} else if err.Error() != "false could not be cast as a string" {
		t.Error(err)
	}
}

func TestNotInJsonUnmarshalError(t *testing.T) {
	f := aws.NewFilter()
	conds, err := f.NotIn("key", `[`)

	if conds != nil {
		t.Error("Conditions should be nil")
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestFilter(t *testing.T) {
	f := aws.NewFilter()

	filters := map[string][]string{
		"key":      {"value"},
		"key2":     {"value1", "value2"},
		"key3__gt": {"value3"},
		"key4__lt": nil,
	}

	conditions, err := f.Filter(nil, filters, "")
	if err != nil {
		t.Error(err)
	}

	conds := conditions.(expression.ConditionBuilder)
	if !conds.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conds).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "((#0 = :0) AND ((#1 = :1) OR (#1 = :2))) AND (#2 > :3)" {
		t.Errorf("Invalid filter expression. Expected '((#0 = :0) AND ((#1 = :1) OR (#1 = :2))) AND (#2 > :3)' but got '%s'", *expr.Filter())
	}
}

func TestFilterWithUser(t *testing.T) {
	f := aws.NewFilter()

	filters := map[string][]string{
		"key":      {"value"},
		"key2":     {"value1", "value2"},
		"key3__gt": {"value3"},
		"key4__lt": nil,
	}

	user := &types.User{
		Permissions: types.Permissions{
			ReadFilters: []types.FilterField{
				{
					Field:    "key",
					Operator: base.OperationContains,
					Value:    "value4",
				},
			},
		},
	}

	conditions, err := f.Filter(user, filters, "")
	if err != nil {
		t.Error(err)
	}

	conds := conditions.(expression.ConditionBuilder)
	if !conds.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conds).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "(((contains (#0, :0)) AND (#0 = :1)) AND ((#1 = :2) OR (#1 = :3))) AND (#2 > :4)" {
		t.Errorf("Invalid filter expression. Expected '(((contains (#0, :0)) AND (#0 = :1)) AND ((#1 = :2) OR (#1 = :3))) AND (#2 > :4)' but got '%s'", *expr.Filter())
	}
}

func TestMultiFilter(t *testing.T) {
	f := aws.NewFilter()

	conditions, err := f.MultiFilter(nil, "key", []string{"value1", "value2", "value3"})
	if err != nil {
		t.Error(err)
	}

	conds := conditions.(expression.ConditionBuilder)
	if !conds.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conds).Build()
	if err != nil {
		t.Error(err)
	}

	if *expr.Filter() != "#0 IN (:0, :1, :2)" {
		t.Errorf("Invalid filter expression. Expected '#0 IN (:0, :1, :2)' but got '%s'", *expr.Filter())
	}
}
