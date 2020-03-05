package main 
import "os"
import "fmt"


func main() {

   var buf []byte = make([]byte,10)

   var r int
   var err error
   for r,err = os.Stdin.Read(buf); err != nil {
      fmt.Print(buf[:r])

   } 
    
   fmt.Print("err",err)
 

}
