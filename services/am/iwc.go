package am

import (
	"encoding/binary"

	"github.com/ca1e/gonx/nx/nxerrors"
	"github.com/ca1e/gonx/nx/nxtypes"
	"github.com/ca1e/gonx/services/ipc"
)

func IwcGetAppletResourceUserId() (nxtypes.ARUID, error) {
	if amInitializations <= 0 {
		return 0, nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(1)
	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // one uint64

	err := ipc.Send(iWindowController, &rq, &rs)
	if err != nil {
		return 0, err
	}

	return nxtypes.ARUID(binary.LittleEndian.Uint64(rs.RawData)), nil
}

func IwcAcquireForegroundRights() error {
	if amInitializations <= 0 {
		return nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(10)
	rs := ipc.ResponseFmt{}

	return ipc.Send(iWindowController, &rq, &rs)
}
