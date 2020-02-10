/*
This module consist of cache implementation of cache
and has global variable that will be require to access
the access the cache of Routing server.
*/
package servercac

import  "github.com/sirupsen/logrus"


// This function permit us to initialize various cache related
// variable at the start of the go subroutine.
func InitServerCache()  {
	logrus.Debug("Initializing various variable related cache.")


	logrus.Debug("Completed Initialization of cache.")
}