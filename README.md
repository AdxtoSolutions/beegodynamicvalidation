# Beego Dynamic Validation
Quite hacky solution for calling beego validate by string identifiers instead of explicitly or by struct tags.

Example Usage:

    import bdv "github.com/byrnedo/beegodynamicvalidation"

    ...

    dynV := bdv.DynamicValidation{}

    isValid, err := dynV.ValidByStrings("MyField", "Required;Alpha", "myvalue")

    ...
