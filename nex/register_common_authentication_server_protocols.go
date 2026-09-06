package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonticketgranting "github.com/PretendoNetwork/nex-protocols-common-go/v2/ticket-granting"
	ticketgranting "github.com/PretendoNetwork/nex-protocols-go/v2/ticket-granting"
	"github.com/Protarium-Network/splatoon-testfire-nex/globals"
)

func registerCommonAuthenticationServerProtocols() {
	ticketGrantingProtocol := ticketgranting.NewProtocol()
	globals.AuthenticationEndpoint.RegisterServiceProtocol(ticketGrantingProtocol)
	commonTicketGrantingProtocol := commonticketgranting.NewCommonProtocol(ticketGrantingProtocol)

	port, _ := strconv.Atoi(os.Getenv("PN_GLOBAL_TESTFIRE_SECURE_SERVER_PORT"))

	secureStationURL := types.NewStationURL("")
	secureStationURL.SetURLType(constants.StationURLPRUDPS)
	secureStationURL.SetAddress(os.Getenv("PN_GLOBAL_TESTFIRE_SECURE_SERVER_HOST"))
	secureStationURL.SetPortNumber(uint16(port))
	secureStationURL.SetConnectionID(1)
	secureStationURL.SetPrincipalID(types.NewPID(2))
	secureStationURL.SetStreamID(1)
	secureStationURL.SetStreamType(constants.StreamTypeRVSecure)
	secureStationURL.SetType(uint8(constants.StationURLFlagPublic))

	commonTicketGrantingProtocol.SecureStationURL = secureStationURL
	// The RPX does not embed the historical server build string. Leave this
	// configurable so a captured value can be supplied without changing code.
	commonTicketGrantingProtocol.BuildName = types.NewString(os.Getenv("PN_GLOBAL_TESTFIRE_NEX_BUILD_NAME"))
	commonTicketGrantingProtocol.SecureServerAccount = globals.SecureServerAccount
	if globals.LocalAuthMode {
		commonTicketGrantingProtocol.ValidateLoginData = func(pid types.PID, loginData types.DataHolder) *nex.Error {
			// my friends always told me i was valid no matter what. surely my login data is the same
			return nil
		}
	} else {
		commonTicketGrantingProtocol.ConfigurePNValidation([]string{globals.GameServerID})
	}
}
