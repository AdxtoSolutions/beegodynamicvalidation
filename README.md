# Beego Dynamic Validation
Quite hacky solution for calling beego validate by string identifiers instead of explicitly or by struct tags.

Example Usage:

    import "github.com/byrnedo/BeegoDynamicValidator/dynamicvalidation"

    ...

    dynV := DynamicValidation{}

    isValid, err := dynV.ValidByStrings("MyField", "Required;Alpha", "myvalue")

    ...
