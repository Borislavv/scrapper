package livenessinterface

// Prober can handle one service/application, if you need observe few services, just create a new instance for each app.
type Prober interface {
	// Watch starts a goroutine which watches for new messages in the liveness channel and respond.
	// This method must be called on main service thread before the lock action will be called.
	// Example:
	//        cancelLivenessProbe := Probe.MonitorWorkness()
	//        wg.Add(1)
	//        go useful.Work(wg)
	//        // service is alive and probe will be report it
	//		  wg.Wait()
	//        cancelLivenessProbe()
	//        // service and probe are closed, probe will be report of failed probe
	Watch(s Service) (cancel func())
	// IsAlive checks wether is target service (a service which called Watch method) alive.
	IsAlive() bool
}
