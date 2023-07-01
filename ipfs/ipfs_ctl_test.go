package ipfs

import (
    "testing"
)

func TestPullRunDefaultStartStop(t *testing.T) {

    
    err := PullDefault()
    if err != nil {
        t.Errorf("Pull Failed: %s",err) 
    }     

    runInfo,err2 := RunDefault()
    if err2 != nil {
        t.Errorf("RunDefault failed: %s",err2)
    }
    
    err3 := Stop(runInfo);
    if err3 != nil {
        t.Errorf("Stop failed: %s",err3)
    } 

    err4 := Start(runInfo);
    if err4 != nil {
        t.Errorf("(Re)Start failed: %s",err4);
    }

    err5 := Stop(runInfo);
    if (err5 != nil) {
        t.Errorf("Stop failed (after Restart): %s",err5);
    }


}
