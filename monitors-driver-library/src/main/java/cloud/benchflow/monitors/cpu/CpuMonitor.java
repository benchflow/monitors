package cloud.benchflow.monitors.cpu;

import cloud.benchflow.monitors.Monitor;
import cloud.benchflow.monitors.MonitorAPI;
import com.sun.faban.driver.transport.hc3.ApacheHC3Transport;

import java.util.Map;
import java.util.logging.Logger;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public class CpuMonitor extends Monitor {

    public CpuMonitor(ApacheHC3Transport http,
                      Map<String, String> params,
                      String endpoint,
                      MonitorAPI api) {
        this(http, params, endpoint, api, Logger.getLogger(CpuMonitor.class.getName()));
    }

    public CpuMonitor(ApacheHC3Transport http,
                      Map<String, String> params,
                      String endpoint,
                      MonitorAPI api,
                      Logger logger) {
        super(http, params, endpoint, api, logger);
    }

    @Override
    public void monitor() throws Exception {
        //this doesn't do anything for now
    }

    @Override
    protected void start() throws Exception {
        //this doesn't do anything for now
    }

    @Override
    protected void stop() throws Exception {
        //this doesn't do anything for now
    }

}