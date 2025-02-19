package sm

import (
	"fmt"

	"github.com/ca1e/gonx/nx/nxerrors"
	"github.com/ca1e/gonx/nx/nxtypes"
	"github.com/ca1e/gonx/services/ipc"
	"github.com/ca1e/gonx/svc"
)

var smInitializations = 0
var smObject ipc.Object

const (
	smServiceName = "sm:\x00"
	debugSm       = false
)

// str2u64 converts a string to uint64 representation
// used on SM service name
func str2u64(str string) uint64 {
	var b [8]byte

	for i := 0; i < 8; i++ {
		if len(str) <= i {
			break
		}
		b[i] = str[i]
	}

	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

// Init initializes the SM Service if needed
func Init() error {
	if debugSm {
		println("SM::Init")
	}
	if smInitializations > 0 {
		smInitializations++ // Already initialized, increment ref count
		return nil
	}

	smInitializations++

	smObject.ObjectID = -1
	handle := nxtypes.Handle(0)
	smName := []byte(smServiceName)
	r := svc.ConnectToNamedPort(&handle, &smName[0])
	if r != nxtypes.ResultOK {
		smInitializations--
		return nxerrors.IPCError{
			Message: "error initializing sm",
			Result:  r,
		}
	}

	smObject.SetSession(handle)

	// sm:#0 Initialize
	rq := ipc.MakeDefaultRequest(0)
	rq.SendPID = true
	rq.SetRawDataFromUint64(uint64(0))

	rs := ipc.ResponseFmt{}

	err := ipc.Send(smObject, &rq, &rs)
	if err != nil {
		if debugSm {
			fmt.Printf("error initializing sm: %s", err)
		}
		_ = ipc.Close(&smObject)
		smInitializations--
		return err
	}

	return nil
}

// Finalize closes the a initialized SM Service
func Finalize() {
	smInitializations--
	if smInitializations == 0 {
		smForceFinalize()
	}
}

func smForceFinalize() {
	_ = ipc.Close(&smObject)
	smInitializations = 0
}

func GetService(outObject *ipc.Object, name string) error {
	if debugSm {
		fmt.Printf("SM::GetService(%p, %s)\n", outObject, name)
	}
	if smObject.GetSession() == 0 {
		return nxerrors.SmNotInitialized
	}

	if len(name) > 8 {
		return nxerrors.SmServiceNameTooLong
	}

	serviceName := str2u64(name)
	outObject.ObjectID = -1
	outObject.IsBorrowed = false

	rq := ipc.MakeDefaultRequest(1)
	rq.SetRawDataFromUint64(serviceName)

	rs := ipc.ResponseFmt{}
	rs.MoveHandles = make([]nxtypes.Handle, 1)

	err := ipc.Send(smObject, &rq, &rs)
	if err != nil {
		return err
	}

	outObject.SetSession(rs.MoveHandles[0])

	return nil
}
