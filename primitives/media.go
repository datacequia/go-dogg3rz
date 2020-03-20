package primitives

//	"reflect"

//"github.com/adpadilla/go-dogg3rz/errors"
//"github.com/fatih/structs"

//const TYPE_DOGG3RZ_MEDIA = "dogg3rz.media"

const TYPE_DOGG3RZ_MEDIA Dogg3rzObjectType = 1 << 4

//const D_ATTR_ENTRIES = "entries"

type dgrzMedia struct {
	name string

	parent string
}
