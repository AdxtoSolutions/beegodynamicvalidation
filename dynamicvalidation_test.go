package beegodynamicvalidation

import (
	"testing"
)

func TestDynamicValidator(t *testing.T) {
	dynV := DynamicValidation{}

	isValid, err := dynV.ValidByStrings("MyField", "Required", "test")

	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if !isValid {
		t.Error("Should not have failed validation")
	}

	isValid, err = dynV.ValidByStrings("MyField", "Required", "")
	dynV.Clear()
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	//fmt.Printf("Validation error given: %+q", dynV.Validation.Errors)
	if isValid {
		t.Error("Empty required string should have failed validation")
	}
	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;Alpha", "alphaman")
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if !isValid {
		t.Errorf("Should not have failed validation: %+v", dynV.ErrorsMap)
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;Alpha", "4lph4Num3r1cman")
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if isValid {
		t.Error("Non alpha numeric required string should have failed validation")
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;Numeric", "99999")
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if !isValid {
		t.Errorf("Should not have failed validation: %+v", dynV.ErrorsMap)
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;Min(9998)", 9999)
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if !isValid {
		t.Errorf("Numeric greater than min should not have failed validation: %+v", dynV.ErrorsMap)
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;Min(9998)", 9997)
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if isValid {
		t.Errorf("Numeric less than min should not have passed validation: %+v", dynV.ErrorsMap)
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;MinSize(6)", "teststring")
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if !isValid {
		t.Errorf("String greater than min should not have failed validation: %+v", dynV.ErrorsMap)
	}

	dynV.Clear()
	isValid, err = dynV.ValidByStrings("MyField", "Required;MinSize(76)", "teststring")
	if err != nil {
		t.Errorf("Should not have thrown runtime error [%s]", err.Error())
	}
	if isValid {
		t.Errorf("String less than min should have failed validation: %+v", dynV.ErrorsMap)
	}

}
