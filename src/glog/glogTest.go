package glog

func TestWriteLog() {
	for i := 0; i < 100000000; i++ {
		PushLogicLog(Lerror, "PushLogicLogglog.  PushLogicLogglog.PushLogicLogglog.PushLogicLog glog.PushLogicLog ")
		PushDbLog(Linfo, "PushDblog glog.PushDblog  glog.PushDblogglog.PushDblog   glog.PushDblogglog.PushDblog")
	}
}
